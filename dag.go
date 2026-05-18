package dag_go

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/HeaInSeo/utils"
	"golang.org/x/sync/errgroup"
)

// ==================== 상수 정의 ====================

// Create, Exist, Fault are the result codes returned by createEdge.
const (
	Create createEdgeErrorType = iota
	Exist
	Fault
)

// The status displayed when running the runner on each node.
const (
	Start runningStatus = iota
	Preflight
	PreflightFailed
	InFlight
	InFlightFailed
	PostFlight
	PostFlightFailed
	FlightEnd
	Failed
	Succeed
)

// StartNode and EndNode are the reserved IDs for the synthetic entry and exit nodes
// that are automatically created by StartDag and FinishDag respectively.
const (
	StartNode = "start_node"
	EndNode   = "end_node"
)

// It is the node ID when the condition that the node cannot be created.
const noNodeID = "-1"

// defaultErrorDrainFallback is the safe fallback drain timeout used when
// DagConfig.ErrorDrainTimeout is not set (zero or negative).
const defaultErrorDrainFallback = 5 * time.Second

// ==================== 타입 정의 ====================

// DagConfig holds tunable parameters for a Dag instance.
// Use DefaultDagConfig for production-ready defaults, or override individual
// fields before passing to NewDagWithConfig.
type DagConfig struct {
	// MinChannelBuffer is the buffer size for inter-node edge channels.  Larger
	// values reduce the chance of a parent blocking while writing to a slow child.
	// Default: 5.
	MinChannelBuffer int

	// MaxChannelBuffer is the buffer capacity for the NodesResult and Errors
	// aggregation channels.  Set this higher than the total number of nodes to
	// prevent back-pressure stalls.  Default: 100.
	MaxChannelBuffer int

	// StatusBuffer is reserved for future per-node status channel buffering.
	// Default: 10.
	StatusBuffer int

	// WorkerPoolSize caps the number of goroutines that execute nodes concurrently.
	// If the number of nodes is smaller than WorkerPoolSize, the pool is sized to
	// the node count instead (min of the two).  Default: 50.
	WorkerPoolSize int

	// DefaultTimeout is the implicit per-node execution timeout applied during
	// inFlight (RunE). It limits how long each node's user-supplied work may run.
	//
	// Zero (the default) means no implicit per-node execution timeout is applied;
	// only the caller's context deadline governs execution length.
	// Per-node overrides are set via Node.Timeout (when Timeout > 0).
	//
	// Dependency wait (preFlight) is never bounded by this value — it always
	// uses the caller's context only, so upstream work never causes false-negative
	// timeouts in downstream nodes.
	DefaultTimeout time.Duration

	// ErrorDrainTimeout is the maximum time collectErrors will wait to drain the
	// Errors channel.  Defaults to 5 s when left at zero (see DefaultDagConfig).
	ErrorDrainTimeout time.Duration

	// ExpectedNodeCount is a capacity hint for the internal nodes map and the
	// Edges slice.  When the final node count is known upfront, setting this
	// field avoids incremental map rehashing and slice growth during AddEdge
	// calls.  Zero means let the runtime decide the initial capacity.
	ExpectedNodeCount int

	// ErrorPolicy controls how downstream nodes react to upstream failures.
	// ErrorPolicyFailFast (default, zero value) skips children of a failed node.
	// ErrorPolicyContinueOnError allows children to run regardless of parent outcome;
	// errors are still recorded in the Errors channel and DroppedErrors counter.
	ErrorPolicy ErrorPolicy
}

// DagOption is a functional-option type for NewDagWithOptions.
// Use the provided With* constructors to build option values.
type DagOption func(*Dag)

// nodeTask bundles the execution arguments for a single node run.
// Passing a concrete struct to the worker queue eliminates the closure
// allocation that the previous func() approach incurred per submission —
// removing per-node heap pressure in large DAG executions.
type nodeTask struct {
	node *Node
	sc   *SafeChannel[*printStatus]
	ctx  context.Context //nolint:containedctx // transient task descriptor; ctx is consumed immediately by the worker
}

// DagWorkerPool manages a bounded pool of goroutines for concurrent node execution.
// Tasks are submitted via Submit; Close drains the queue and waits for all workers.
type DagWorkerPool struct {
	workerLimit int
	taskQueue   chan nodeTask
	wg          sync.WaitGroup
	closeOnce   sync.Once // prevents double-close panic on taskQueue
}

type (
	runningStatus int

	printStatus struct {
		rStatus runningStatus
		nodeID  string
	}
)

// createEdgeErrorType 0 if created, 1 if exists, 2 if error.
type createEdgeErrorType int

// Dag is a Directed Acyclic Graph execution engine.
//
// A Dag is created with InitDag (or NewDag / NewDagWithConfig), populated with
// AddEdge calls, sealed with FinishDag, and then executed via the lifecycle:
//
//	ConnectRunner → GetReady → Start → Wait
//
// A completed DAG can be re-executed by calling Reset followed by the same
// lifecycle.  Dag must always be handled as a pointer; value-copy is forbidden
// because it embeds sync.RWMutex.
type Dag struct {
	// ID is the unique identifier assigned at creation time (UUID v4).
	ID string

	// edges is the ordered list of directed edges in the graph.
	// External callers must use the Edges() method for read access.
	// Mutations must go through AddEdge / AddEdgeIfNodesExist before FinishDag.
	edges []*Edge

	nodes     map[string]*Node
	startNode *Node // synthetic entry node created by StartDag
	endNode   *Node // synthetic exit node created by FinishDag

	validated bool // true after FinishDag succeeds

	// NodesResult is the fan-in channel that collects printStatus events from
	// every node during execution.  Wait reads from this channel.
	NodesResult *SafeChannel[*printStatus]
	nodeResult  []*SafeChannel[*printStatus] // per-node result channels

	// errLogs collects structured errors for post-mortem inspection.
	errLogs []*systemError

	// Errors is the concurrency-safe error channel for runtime RunE failures.
	// Use reportError to write and collectErrors (or Errors.GetChannel()) to read.
	// The channel is recreated by Reset so it is valid for the next run.
	Errors *SafeChannel[error]

	// ── Timeout policy ──────────────────────────────────────────────────────
	//
	// dag-go has four timeout layers, applied in the following order of scope:
	//
	//  1. caller ctx         — outermost hard cap; passed to GetReady and Wait.
	//                          All phases (preFlight, inFlight, postFlight) respect it.
	//
	//  2. Dag.Timeout        — Wait-level deadline; set via WithTimeout.
	//                          Caps the entire dag.Wait call when bTimeout is true.
	//
	//  3. DagConfig.DefaultTimeout — implicit per-node execution timeout.
	//                          Applied to inFlight (RunE) only; zero = no implicit limit.
	//                          Dependency wait (preFlight) never consumes this budget.
	//
	//  4. Node.Timeout       — explicit per-node execution override.
	//                          When Node.Timeout > 0 it takes priority over DefaultTimeout.
	//
	// Key invariant: dependency wait (preFlight) always uses the caller ctx only.
	// Execution timeouts are applied by connectRunner immediately before inFlight.
	// ────────────────────────────────────────────────────────────────────────

	// Timeout is the DAG-level execution deadline applied when bTimeout is true.
	// Set via WithTimeout or by assigning directly before GetReady.
	Timeout  time.Duration
	bTimeout bool

	// ContainerCmd is the global default Runnable applied to every node that
	// has no per-node override and no resolver match.
	// Set via SetContainerCmd to ensure thread-safe mutation.
	ContainerCmd Runnable

	runnerResolver RunnerResolver // optional dynamic runner selector

	// Config holds the tunable parameters active for this DAG instance.
	Config DagConfig

	// execSem caps the number of nodes that may run inFlight (RunE) concurrently.
	// It is a buffered channel with capacity == WorkerPoolSize, initialised by
	// GetReady.  preFlight (dependency wait) never acquires a slot, so a small
	// cap cannot deadlock the graph.
	execSem chan struct{}
	// nodeWg tracks every goroutine launched by GetReady; Wait drains it.
	nodeWg sync.WaitGroup

	// started is set by StartE/Start to prevent the trigger signal from being
	// sent more than once per run.  Reset clears it.
	started atomic.Bool

	// running is true while GetReady goroutines are live (from GetReadyE until
	// the combined cleanup defer in Wait clears it).  Reset checks this flag so
	// it cannot replace live channels underneath running goroutines.
	running atomic.Bool

	nodeCount      int64        // total user-defined node count (atomic)
	completedCount int64        // nodes that called MarkCompleted (atomic)
	droppedErrors  int64        // errors dropped by reportError (atomic)
	mu             sync.RWMutex // guards nodes map, ContainerCmd, runnerResolver, nodeResult, startTrigger

	// startTrigger is the validated start-node trigger channel captured by
	// GetReadyE after structural invariants are confirmed.  StartE uses only
	// this field to fire the DAG — it never reads startNode.parentVertex while
	// goroutines are live.  Nil until GetReadyE succeeds; cleared by reset().
	startTrigger *SafeChannel[runningStatus]
}

// Edge represents a directed connection between two nodes.
// The embedded safeVertex channel carries runningStatus signals from the parent
// node to its child during execution.  Edges are created by AddEdge and reset
// by Reset; callers should not manipulate the fields directly.
type Edge struct {
	parentID   string
	childID    string
	safeVertex *SafeChannel[runningStatus]
}

// ==================== DAG 기본 및 옵션 함수 ====================

// DefaultDagConfig returns a DagConfig populated with production-ready defaults:
//   - MinChannelBuffer:  5
//   - MaxChannelBuffer:  100
//   - StatusBuffer:      10
//   - WorkerPoolSize:    50
//   - DefaultTimeout:    0 (no implicit per-node execution timeout)
//   - ErrorDrainTimeout: 5 s
//
// DefaultTimeout is intentionally 0: individual node execution time is not
// bounded by default. Use WithDefaultTimeout or set Node.Timeout explicitly
// when per-node execution limits are required. DAG-wide time limits should
// be enforced via the caller context passed to GetReady and Wait.
func DefaultDagConfig() DagConfig {
	return DagConfig{
		MinChannelBuffer:  5,
		MaxChannelBuffer:  100,
		StatusBuffer:      10,
		WorkerPoolSize:    50,
		DefaultTimeout:    0,
		ErrorDrainTimeout: 5 * time.Second,
	}
}

// NewDag returns a new Dag with default configuration (see DefaultDagConfig).
// Call StartDag (or use InitDag) to add the synthetic start node before adding edges.
func NewDag() *Dag {
	return NewDagWithConfig(DefaultDagConfig())
}

// normalizeDagConfig clamps any out-of-range DagConfig fields to safe minimums.
// It is called by NewDagWithConfig and by WithWorkerPool / WithChannelBuffers so
// that the DAG always starts with a coherent configuration regardless of caller
// input.
func normalizeDagConfig(cfg DagConfig) DagConfig {
	if cfg.MinChannelBuffer < 1 {
		cfg.MinChannelBuffer = 1
	}
	if cfg.MaxChannelBuffer < 1 {
		cfg.MaxChannelBuffer = 1
	}
	if cfg.StatusBuffer < 1 {
		cfg.StatusBuffer = 1
	}
	if cfg.WorkerPoolSize < 1 {
		cfg.WorkerPoolSize = 1
	}
	if cfg.ErrorDrainTimeout <= 0 {
		cfg.ErrorDrainTimeout = 5 * time.Second
	}
	// ExpectedNodeCount and DefaultTimeout may legitimately be 0; no clamping.
	return cfg
}

// NewDagWithConfig returns a new Dag using the supplied DagConfig.
// Prefer InitDag for the common case; use this constructor when you need to
// customise buffer sizes, timeouts, or the worker pool size before adding nodes.
// Invalid config values (zero or negative buffers / pool size) are clamped to
// safe minimums by normalizeDagConfig.
func NewDagWithConfig(config DagConfig) *Dag {
	config = normalizeDagConfig(config)
	nodeCapacity := config.ExpectedNodeCount
	return &Dag{
		nodes:       make(map[string]*Node, nodeCapacity),
		edges:       make([]*Edge, 0, nodeCapacity),
		ID:          uuid.NewString(),
		NodesResult: NewSafeChannelGen[*printStatus](config.MaxChannelBuffer),
		Config:      config,
		Errors:      NewSafeChannelGen[error](config.MaxChannelBuffer),
	}
}

// NewDagWithOptions returns a new Dag with DefaultDagConfig, then applies each
// DagOption in order.  Functional options (e.g. WithTimeout, WithWorkerPool)
// are applied after default values, so later options can override earlier ones.
// The final config is normalised after all options are applied.
func NewDagWithOptions(options ...DagOption) *Dag {
	dag := NewDagWithConfig(DefaultDagConfig())

	for _, option := range options {
		option(dag)
	}
	// Re-normalise after options may have set invalid values (e.g. WithWorkerPool(0)).
	dag.Config = normalizeDagConfig(dag.Config)

	return dag
}

// InitDag creates a new DAG with default configuration and immediately calls
// StartDag to create the synthetic start node.  It is the recommended entry
// point for most users.
func InitDag() (*Dag, error) {
	dag := NewDag()
	if dag == nil {
		return nil, fmt.Errorf("failed to run NewDag")
	}
	return dag.StartDag()
}

// ==================== 전역 러너/리졸버 ====================

// SetContainerCmd sets the global default Runnable for all nodes in this DAG.
// It is safe to call concurrently.  Per-node overrides (SetNodeRunner) and the
// RunnerResolver take priority over this value; see runner priority in the README.
//
// Mutation policy: SetContainerCmd must be called before GetReady and must not
// be called again until Reset.  The DAG is considered frozen from the moment
// GetReadyE succeeds until reset() clears nodeResult — this covers both the
// running window (goroutines live) and the post-Wait / pre-Reset window.
// Mutating the global runner inside the frozen window would cause different nodes
// to execute with different runners depending on scheduling order, breaking
// reproducibility.
func (dag *Dag) SetContainerCmd(r Runnable) {
	dag.mu.Lock()
	defer dag.mu.Unlock()
	if dag.nodeResult != nil || dag.running.Load() {
		Log.Error("SetContainerCmd: DAG is frozen after GetReady; ignoring — call Reset first")
		return
	}
	dag.ContainerCmd = r
}

// loadDefaultRunnerAtomic 추가
/* func (dag *Dag) loadDefaultRunnerAtomic() Runnable {
	v := dag.defVal.Load()
	if v == nil {
		return nil
	}
	return v.(*runnerSlot).r
} */

// SetRunnerResolver installs a dynamic runner selector for this DAG.
// The resolver is called at execution time for each node, after the per-node
// atomic override is checked but before the global ContainerCmd fallback.
// Pass nil to clear a previously installed resolver.  Thread-safe.
//
// Mutation policy: same as SetContainerCmd — must be called before GetReady and
// not again until Reset.  The frozen window spans from GetReadyE success through
// Reset; changing the resolver inside it would break reproducibility.
func (dag *Dag) SetRunnerResolver(rr RunnerResolver) {
	dag.mu.Lock()
	defer dag.mu.Unlock()
	if dag.nodeResult != nil || dag.running.Load() {
		Log.Error("SetRunnerResolver: DAG is frozen after GetReady; ignoring — call Reset first")
		return
	}
	dag.runnerResolver = rr
}

// 원자적으로 Resolver 반환
/* func (dag *Dag) loadRunnerResolverAtomic() RunnerResolver {
	// rrVal은 단지 "초기화 여부"를 위한 타입 고정용이고,
	// 실제 rr는 락으로 보호된 dag.runnerResolver에서 읽는다.
	// TODO 완전히 락-프리로 하려면 rrVal에 rr 자체를 담는 별도 래퍼 타입을 써야 함.
	dag.mu.RLock()
	rr := dag.runnerResolver
	dag.mu.RUnlock()
	return rr
} */

// SetNodeRunner sets the runner for the node with the given id.
// Returns false if the node does not exist, is not in Pending status, or if the
// DAG is frozen (GetReady has been called and Reset has not yet been called).
//
// Mutation policy: same as SetContainerCmd — must be called before GetReady and
// not again until Reset.  The frozen window spans from GetReadyE success through
// Reset.
func (dag *Dag) SetNodeRunner(id string, r Runnable) bool {
	dag.mu.RLock()
	n := dag.nodes[id]
	frozen := dag.nodeResult != nil || dag.running.Load()
	dag.mu.RUnlock()

	if frozen {
		Log.Errorf("SetNodeRunner: DAG is frozen after GetReady; ignoring node %s — call Reset first", id)
		return false
	}
	if n == nil {
		return false
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	switch n.status {
	case NodeStatusPending:
		// We already hold n.mu — call runnerStore directly to avoid re-entrant lock.
		// (SetRunner would try to acquire n.mu again, causing a self-deadlock.)
		n.runnerStore(r)
		return true

	case NodeStatusRunning, NodeStatusSucceeded, NodeStatusFailed, NodeStatusSkipped:
		Log.Infof("SetNodeRunner ignored: node %s status=%v", n.ID, n.status)
		return false

	default:
		Log.Warnf("SetNodeRunner unknown status: node %s status=%v", n.ID, n.status)
		return false
	}
}

// SetNodeRunners bulk-sets runners from a map of node-id to Runnable.
// Returns the count of applied runners, and slices of missing/skipped node ids.
//
// Mutation policy: same as SetContainerCmd — must be called before GetReady and
// not again until Reset.  When the DAG is frozen all ids are returned in the
// skipped slice with applied == 0.
func (dag *Dag) SetNodeRunners(m map[string]Runnable) (applied int, missing, skipped []string) {
	dag.mu.RLock()
	defer dag.mu.RUnlock()

	if dag.nodeResult != nil || dag.running.Load() {
		Log.Error("SetNodeRunners: DAG is frozen after GetReady; ignoring all entries — call Reset first")
		for id := range m {
			skipped = append(skipped, id)
		}
		return 0, nil, skipped
	}

	for id, r := range m {
		// 실행 노드가 아닌 특수 노드 방어
		if id == StartNode || id == EndNode {
			skipped = append(skipped, id)
			continue
		}

		n := dag.nodes[id]
		if n == nil {
			missing = append(missing, id)
			continue
		}

		n.mu.Lock()
		switch n.status {
		case NodeStatusPending:
			// atomic.Value must always receive *runnerSlot to prevent a type-change panic.
			n.runnerStore(r)
			applied++

		case NodeStatusRunning, NodeStatusSucceeded, NodeStatusFailed, NodeStatusSkipped:
			// 이미 실행 중/완료/실패/스킵된 노드는 건너뜀 (원자성 보장)
			skipped = append(skipped, id)

		default:
			// 미래에 상태가 늘어나도 여기로 들어오면 “안전하게” 건너뜀
			skipped = append(skipped, id)
		}
		n.mu.Unlock()
	}
	return applied, missing, skipped
}

// InitDagWithOptions creates and initialises a new DAG, applying the supplied
// functional options before adding the synthetic start node.  It is the
// option-friendly equivalent of InitDag.
func InitDagWithOptions(options ...DagOption) (*Dag, error) {
	dag := NewDagWithOptions(options...)
	if dag == nil {
		return nil, fmt.Errorf("failed to run NewDag")
	}
	return dag.StartDag()
}

// WithTimeout sets the DAG-level execution deadline used by Wait.
// This is distinct from WithDefaultTimeout, which bounds individual node
// execution (inFlight).  Pass 0 or omit to rely on the caller context only.
func WithTimeout(timeout time.Duration) DagOption {
	return func(dag *Dag) {
		dag.Timeout = timeout
		dag.bTimeout = true
	}
}

// WithDefaultTimeout sets the implicit per-node execution timeout applied
// during inFlight (RunE).  This is distinct from WithTimeout, which caps the
// overall DAG run in Wait().
//
// Semantics:
//   - d == 0: no implicit per-node execution timeout (the default). Each
//     node's RunE runs until the caller context expires or the node returns.
//   - d  > 0: RunE is bounded by d. If RunE does not return within d, its
//     context is cancelled and the node is marked failed.
//
// Dependency wait (preFlight) is never bounded by this value; it always
// honours the caller context only.  This prevents long upstream execution
// chains from causing false-negative timeouts in downstream nodes.
//
// Per-node overrides: set Node.Timeout > 0 on individual nodes;
// that takes priority over this DAG-wide default.
func WithDefaultTimeout(d time.Duration) DagOption {
	return func(dag *Dag) {
		dag.Config.DefaultTimeout = d
	}
}

// WithChannelBuffers sets the channel buffer sizes used for edge signalling and
// result aggregation.  Larger buffers reduce back-pressure in high-fan-out DAGs.
func WithChannelBuffers(minBuffer, maxBuffer, statusBuffer int) DagOption {
	return func(dag *Dag) {
		dag.Config.MinChannelBuffer = minBuffer
		dag.Config.MaxChannelBuffer = maxBuffer
		dag.Config.StatusBuffer = statusBuffer
	}
}

// WithWorkerPool sets the maximum number of concurrent node-execution goroutines.
// If the DAG has fewer nodes than size, the pool is sized to the node count instead.
func WithWorkerPool(size int) DagOption {
	return func(dag *Dag) {
		dag.Config.WorkerPoolSize = size
	}
}

// WithErrorPolicy sets the error propagation policy for this DAG.
// ErrorPolicyFailFast (default) skips downstream nodes when a parent fails.
// ErrorPolicyContinueOnError allows all nodes to run regardless of parent failures;
// errors are still reported via the Errors channel.
func WithErrorPolicy(p ErrorPolicy) DagOption {
	return func(dag *Dag) {
		dag.Config.ErrorPolicy = p
	}
}

// ==================== DagWorkerPool 메서드 ====================

// NewDagWorkerPool creates a new worker pool with the given number of goroutines.
// The internal task queue is buffered to twice the worker count so that
// callers are not serialised behind goroutine startup latency.
func NewDagWorkerPool(limit int) *DagWorkerPool {
	pool := &DagWorkerPool{
		workerLimit: limit,
		taskQueue:   make(chan nodeTask, limit*2),
	}

	for range limit {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for task := range pool.taskQueue {
				// Skip execution if the caller's context is already done.
				select {
				case <-task.ctx.Done():
				default:
					task.node.runner(task.ctx, task.sc)
				}
			}
		}()
	}

	return pool
}

// Submit enqueues a nodeTask for execution by the worker pool.
func (p *DagWorkerPool) Submit(task nodeTask) {
	p.taskQueue <- task
}

// Close 워커 풀을 종료.
// sync.Once 를 통해 taskQueue 를 한 번만 닫으므로 이중 호출 시 패닉이 발생하지 않는다.
func (p *DagWorkerPool) Close() {
	p.closeOnce.Do(func() { close(p.taskQueue) })
	p.wg.Wait()
}

// ==================== Dag 메서드 ====================

// StartDag creates the synthetic start node and its trigger channel.
// It is called automatically by InitDag; call it directly only when building a
// DAG with NewDag or NewDagWithConfig.
//
// Returns the receiver so calls can be chained, or an error if the start node
// could not be created (e.g. the DAG is nil or the node was already created).
func (dag *Dag) StartDag() (*Dag, error) {
	dag.mu.Lock()
	defer dag.mu.Unlock()

	// logErr 헬퍼 함수 정의: 에러 로그에 추가하고 reportError 호출
	logErr := func(err error) error {
		dag.errLogs = append(dag.errLogs, &systemError{StartDag, err}) // 필요에 따라 타입을 수정
		dag.reportError(err)
		return err
	}

	// Reject double-call: overwriting startNode with nil would crash Start().
	if dag.startNode != nil {
		return nil, logErr(fmt.Errorf("StartDag already called"))
	}
	node := dag.createNode(StartNode)
	if node == nil {
		return nil, logErr(fmt.Errorf("failed to create start node"))
	}
	dag.startNode = node

	// 새 제네릭 SafeChannel 생성 후, 시작 노드의 parentVertex에 추가
	safeChan := NewSafeChannelGen[runningStatus](dag.Config.MinChannelBuffer)
	dag.startNode.parentVertex = append(dag.startNode.parentVertex, safeChan)
	return dag, nil
}

// reportError delivers err to the error channel in a non-blocking fashion.
// When the channel is full or closed, the error is not silently discarded:
// the droppedErrors counter is incremented (observable via DroppedErrors) and
// a structured log entry is emitted so that log-aggregation pipelines can alert
// on the event.  The cumulative drop count is included in the log field
// "dropped_total" to make SLO alerting straightforward.
func (dag *Dag) reportError(err error) {
	if !dag.Errors.Send(err) {
		n := atomic.AddInt64(&dag.droppedErrors, 1)
		Log.WithField("dag_id", dag.ID).
			WithField("dropped_total", n).
			WithError(err).
			Warn("error channel full or closed; dropping error")
	}
}

// collectErrors drains the error channel until it is empty, or until
// DagConfig.ErrorDrainTimeout (default 5 s) or ctx fires — whichever is first.
//
// The implementation blocks on channel receive without polling: it returns as
// soon as the Errors channel is closed (all errors have been flushed by
// closeChannels), the drain timeout fires, or ctx is cancelled.
func (dag *Dag) collectErrors(ctx context.Context) []error {
	var errs []error

	drainTimeout := dag.Config.ErrorDrainTimeout
	if drainTimeout <= 0 {
		drainTimeout = defaultErrorDrainFallback
	}
	timeout := time.After(drainTimeout)
	ch := dag.Errors.GetChannel()

	for {
		select {
		case err, ok := <-ch:
			if !ok {
				// Errors channel closed — all errors have been collected.
				return errs
			}
			errs = append(errs, err)
		case <-timeout:
			return errs
		case <-ctx.Done():
			return errs
		}
	}
}

// createEdge creates an Edge with safety mechanisms.
func (dag *Dag) createEdge(parentID, childID string) (*Edge, createEdgeErrorType) {
	if utils.IsEmptyString(parentID) || utils.IsEmptyString(childID) {
		return nil, Fault
	}

	// 이미 존재하는 엣지 확인
	if edgeExists(dag.edges, parentID, childID) {
		return nil, Exist
	}

	edge := &Edge{
		parentID:   parentID,
		childID:    childID,
		safeVertex: NewSafeChannelGen[runningStatus](dag.Config.MinChannelBuffer), // 제네릭 SafeChannel 을 사용하여 안전한 채널 생성
	}

	dag.edges = append(dag.edges, edge)
	return edge, Create
}

// closeSafeChannel closes sc and logs the result. It is a no-op when sc is nil.
func closeSafeChannel[T any](sc *SafeChannel[T], label string) {
	if sc == nil {
		return
	}
	if err := sc.Close(); err != nil {
		Log.Warnf("Failed to close %s channel: %v", label, err)
	} else {
		Log.Infof("Closed %s channel", label)
	}
}

// closeChannels safely closes all channels in the DAG.
func (dag *Dag) closeChannels() {
	for _, edge := range dag.edges {
		closeSafeChannel(edge.safeVertex, fmt.Sprintf("edge [%s -> %s]", edge.parentID, edge.childID))
	}
	closeSafeChannel(dag.NodesResult, "NodesResult")
	for i, sc := range dag.nodeResult {
		closeSafeChannel(sc, fmt.Sprintf("nodeResult[%d]", i))
	}
	closeSafeChannel(dag.Errors, "Errors")
}

// getSafeVertex returns the SafeChannel for the edge between parentID and childID,
// or nil if no such edge exists.
func (dag *Dag) getSafeVertex(parentID, childID string) *SafeChannel[runningStatus] {
	for _, v := range dag.edges {
		if v.parentID == parentID && v.childID == childID {
			return v.safeVertex
		}
	}
	return nil
}

// Progress returns the DAG execution completion ratio in [0.0, 1.0].
//
// NOTE: nodeCount and completedCount are read in two separate atomic operations;
// they do not form an atomic pair.  completedCount may be incremented between the
// two reads, making the returned ratio momentarily slightly ahead of reality.
// This is acceptable for progress-bar or observability purposes, but must NOT
// be used for correctness decisions (e.g. deciding whether all nodes finished).
func (dag *Dag) Progress() float64 {
	nodeCount := atomic.LoadInt64(&dag.nodeCount)
	if nodeCount == 0 {
		return 0.0
	}
	completedCount := atomic.LoadInt64(&dag.completedCount)
	return float64(completedCount) / float64(nodeCount)
}

// DroppedErrors returns the number of errors that reportError could not deliver
// to the Errors channel (channel full or closed) since the DAG was created or
// last Reset.  A non-zero value indicates that DagConfig.MaxChannelBuffer is
// too small or that the Errors channel consumer is not draining fast enough.
func (dag *Dag) DroppedErrors() int64 {
	return atomic.LoadInt64(&dag.droppedErrors)
}

// ErrCount returns the number of errors currently buffered in the Errors channel.
// The count is a point-in-time snapshot and may change immediately after the call.
// Use DroppedErrors to check for errors that overflowed the buffer capacity.
func (dag *Dag) ErrCount() int {
	return len(dag.Errors.GetChannel())
}

// Edges returns a shallow copy of the DAG's edge list.  The slice is safe to
// read but mutations to the returned slice do not affect the DAG.
func (dag *Dag) Edges() []*Edge {
	dag.mu.RLock()
	defer dag.mu.RUnlock()
	cp := make([]*Edge, len(dag.edges))
	copy(cp, dag.edges)
	return cp
}

// StartNodeID returns the ID of the synthetic entry node ("start_node"), or an
// empty string if StartDag has not been called yet.
func (dag *Dag) StartNodeID() string {
	dag.mu.RLock()
	defer dag.mu.RUnlock()
	if dag.startNode == nil {
		return ""
	}
	return dag.startNode.ID
}

// EndNodeID returns the ID of the synthetic exit node ("end_node"), or an empty
// string if FinishDag has not been called yet.
func (dag *Dag) EndNodeID() string {
	dag.mu.RLock()
	defer dag.mu.RUnlock()
	if dag.endNode == nil {
		return ""
	}
	return dag.endNode.ID
}

// CreateNode creates a pointer to a new node with thread safety.
// Returns nil if the DAG has already been finalized by FinishDag, or if id is
// a reserved synthetic node name (StartNode / EndNode).
func (dag *Dag) CreateNode(id string) *Node {
	dag.mu.Lock()
	defer dag.mu.Unlock()

	if dag.validated {
		return nil
	}
	if id == StartNode || id == EndNode {
		Log.Warnf("CreateNode: id %q is reserved and cannot be created by the caller", id)
		return nil
	}
	return dag.createNode(id)
}

// CreateNodeWithTimeOut creates a node that applies a per-node timeout when bTimeOut is true.
// Returns nil if the DAG has already been finalized by FinishDag, or if id is
// a reserved synthetic node name (StartNode / EndNode).
func (dag *Dag) CreateNodeWithTimeOut(id string, bTimeOut bool, ti time.Duration) *Node {
	dag.mu.Lock()
	defer dag.mu.Unlock()

	if dag.validated {
		return nil
	}
	if id == StartNode || id == EndNode {
		Log.Warnf("CreateNodeWithTimeOut: id %q is reserved and cannot be created by the caller", id)
		return nil
	}
	if bTimeOut {
		return dag.createNodeWithTimeOut(id, ti)
	}
	return dag.createNode(id)
}

// createNode is the internal implementation of CreateNode.
func (dag *Dag) createNode(id string) *Node {
	// 이미 해당 id의 노드가 존재하면 nil 반환
	if _, exists := dag.nodes[id]; exists {
		return nil
	}

	var node *Node
	if dag.ContainerCmd != nil {
		node = createNode(id, dag.ContainerCmd)
	} else {
		node = createNodeWithID(id)
	}
	// node.runnerStore(dag.ContainerCmd)
	// 추가 초기 스토어: 기본 러너가 없어도 &runnerSlot{}로 non-nil 보장
	node.runnerStore(nil)
	node.parentDag = dag
	dag.nodes[id] = node

	// StartNode 나 EndNode 가 아닌 경우에만 노드 카운트 증가
	if id != StartNode && id != EndNode {
		atomic.AddInt64(&dag.nodeCount, 1)
	}

	return node
}

func (dag *Dag) createNodeWithTimeOut(id string, ti time.Duration) *Node {
	// 이미 해당 id의 노드가 존재하면 nil 반환
	if _, exists := dag.nodes[id]; exists {
		return nil
	}

	var node *Node
	if dag.ContainerCmd != nil {
		// dag.ContainerCmd 이건 모든 노드에 적용되게 하는 옵션이라서 이렇게 함.
		node = createNode(id, dag.ContainerCmd)
	} else {
		node = createNodeWithID(id)
	}

	// 실행 시점 반영 기본: nil 저장
	node.runnerStore(nil)

	node.Timeout = ti
	node.parentDag = dag
	dag.nodes[id] = node

	// StartNode 나 EndNode 가 아닌 경우에만 노드 카운트 증가
	if id != StartNode && id != EndNode {
		atomic.AddInt64(&dag.nodeCount, 1)
	}

	return node
}

// validateEdgeArgs performs argument-only validation that can run before the
// DAG lock is acquired.  dag.reportError is safe to call without the lock
// (SafeChannel.Send has its own lock); dag.errLogs is NOT touched here.
func validateEdgeArgs(from, to string) error {
	if from == to {
		return fmt.Errorf("from-node and to-node are same")
	}
	if utils.IsEmptyString(from) {
		return fmt.Errorf("from-node is empty string")
	}
	if utils.IsEmptyString(to) {
		return fmt.Errorf("to-node is empty string")
	}
	return nil
}

// AddEdge adds an edge between two nodes with improved error handling.
// Returns an error if the DAG has already been finalized by FinishDag.
func (dag *Dag) AddEdge(from, to string) error {
	if err := validateEdgeArgs(from, to); err != nil {
		dag.reportError(err)
		return err
	}

	dag.mu.Lock()
	defer dag.mu.Unlock()

	// logErr appends to errLogs (safe: under dag.mu) and reports to Errors channel.
	logErr := func(err error) error {
		dag.errLogs = append(dag.errLogs, &systemError{AddEdge, err})
		dag.reportError(err)
		return err
	}

	// Topology is frozen once GetReadyE has succeeded (nodeResult is non-nil).
	// Mutations are not permitted until Reset is called.
	if dag.nodeResult != nil {
		return logErr(fmt.Errorf("AddEdge: DAG topology is frozen after GetReady; call Reset before modifying the graph"))
	}
	if dag.validated {
		return logErr(fmt.Errorf("DAG is already finalized: AddEdge is not allowed after FinishDag"))
	}

	// Reserved edge rules (enforced under the lock to be thread-safe):
	//   to   == StartNode → rejected: StartNode is a synthetic entry node; it
	//                        must never have user-defined parents.
	//   from == EndNode   → rejected: EndNode is a synthetic sink managed
	//                        exclusively by FinishDag.
	//   to   == EndNode   → rejected: same reason.
	if to == StartNode {
		return logErr(fmt.Errorf("start_node cannot be an edge target"))
	}
	if from == EndNode || to == EndNode {
		return logErr(fmt.Errorf("end_node is managed internally by FinishDag"))
	}

	// 노드를 가져오거나 생성하는 클로저 함수
	getOrCreateNode := func(id string) (*Node, error) {
		if node := dag.nodes[id]; node != nil {
			return node, nil
		}
		node := dag.createNode(id)
		if node == nil {
			return nil, logErr(fmt.Errorf("%s: createNode returned nil", id))
		}
		return node, nil
	}

	fromNode, err := getOrCreateNode(from)
	if err != nil {
		return err
	}
	toNode, err := getOrCreateNode(to)
	if err != nil {
		return err
	}

	// 엣지 생성 및 검증
	edge, check := dag.createEdge(fromNode.ID, toNode.ID)
	if check == Fault || check == Exist {
		return logErr(fmt.Errorf("edge cannot be created"))
	}
	if edge == nil {
		return logErr(fmt.Errorf("vertex is nil"))
	}

	// 자식과 부모 관계 설정은 엣지 생성이 성공한 후에 수행
	fromNode.children = append(fromNode.children, toNode)
	toNode.parent = append(toNode.parent, fromNode)

	fromNode.childrenVertex = append(fromNode.childrenVertex, edge.safeVertex)
	toNode.parentVertex = append(toNode.parentVertex, edge.safeVertex)
	return nil
}

// AddEdgeIfNodesExist adds an edge only if both nodes already exist.
func (dag *Dag) AddEdgeIfNodesExist(from, to string) error {
	if err := validateEdgeArgs(from, to); err != nil {
		dag.reportError(err)
		return err
	}

	dag.mu.Lock()
	defer dag.mu.Unlock()

	// logErr appends to errLogs (safe: under dag.mu) and reports to Errors channel.
	logErr := func(err error) error {
		dag.errLogs = append(dag.errLogs, &systemError{AddEdgeIfNodesExist, err})
		dag.reportError(err)
		return err
	}

	// Topology is frozen once GetReadyE has succeeded (nodeResult is non-nil).
	// Mutations are not permitted until Reset is called.
	if dag.nodeResult != nil {
		return logErr(fmt.Errorf("AddEdgeIfNodesExist: DAG topology is frozen after GetReady; call Reset before modifying the graph"))
	}
	if dag.validated {
		return logErr(fmt.Errorf("DAG is already finalized: AddEdgeIfNodesExist is not allowed after FinishDag"))
	}

	// Same reserved edge rules as AddEdge.
	if to == StartNode {
		return logErr(fmt.Errorf("start_node cannot be an edge target"))
	}
	if from == EndNode || to == EndNode {
		return logErr(fmt.Errorf("end_node is managed internally by FinishDag"))
	}

	// 노드를 가져오는 클로저 함수: 노드가 없으면 에러 리턴
	getNode := func(id string) (*Node, error) {
		if node := dag.nodes[id]; node != nil {
			return node, nil
		}
		return nil, logErr(fmt.Errorf("%s: node does not exist", id))
	}

	fromNode, err := getNode(from)
	if err != nil {
		return err
	}
	toNode, err := getNode(to)
	if err != nil {
		return err
	}

	// 엣지 생성 및 검증
	edge, check := dag.createEdge(fromNode.ID, toNode.ID)
	if check == Fault || check == Exist {
		return logErr(fmt.Errorf("edge cannot be created"))
	}
	if edge == nil {
		return logErr(fmt.Errorf("vertex is nil"))
	}

	// 자식과 부모 관계 설정은 엣지 생성이 성공한 후에 수행
	fromNode.children = append(fromNode.children, toNode)
	toNode.parent = append(toNode.parent, fromNode)

	fromNode.childrenVertex = append(fromNode.childrenVertex, edge.safeVertex)
	toNode.parentVertex = append(toNode.parentVertex, edge.safeVertex)
	return nil
}

// addEndNode adds an edge to the end node.
func (dag *Dag) addEndNode(fromNode, toNode *Node) error {
	// logErr 헬퍼 함수: 에러 발생 시 errLogs에 기록하고, reportError 호출
	logErr := func(err error) error {
		dag.errLogs = append(dag.errLogs, &systemError{addEndNode, err})
		dag.reportError(err)
		return err
	}

	// 입력 노드 검증
	if fromNode == nil {
		return logErr(fmt.Errorf("fromNode is nil"))
	}
	if toNode == nil {
		return logErr(fmt.Errorf("toNode is nil"))
	}

	// 엣지 생성 및 체크
	edge, check := dag.createEdge(fromNode.ID, toNode.ID)
	if check == Fault || check == Exist {
		return logErr(fmt.Errorf("edge cannot be created"))
	}
	if edge == nil {
		return logErr(fmt.Errorf("vertex is nil"))
	}

	// 부모-자식 관계 설정은 엣지 생성이 성공한 후에 수행
	fromNode.children = append(fromNode.children, toNode)
	toNode.parent = append(toNode.parent, fromNode)

	// 엣지의 vertex를 양쪽 노드에 추가
	fromNode.childrenVertex = append(fromNode.childrenVertex, edge.safeVertex)
	toNode.parentVertex = append(toNode.parentVertex, edge.safeVertex)

	return nil
}

// checkIsolatedNodes returns an error if any node has neither parents nor children.
// A single-node DAG is valid only if that node is StartNode.
// dag.mu must be held by the caller.
func checkIsolatedNodes(nodes []*Node) error {
	for _, n := range nodes {
		if len(n.children) > 0 || len(n.parent) > 0 {
			continue
		}
		if len(nodes) == 1 && n.ID == StartNode {
			continue
		}
		if len(nodes) == 1 {
			return fmt.Errorf("invalid node: only node is not the start node")
		}
		return fmt.Errorf("node '%s' has no parent and no children", n.ID)
	}
	return nil
}

// checkReachable returns an error if any non-StartNode is not reachable from StartNode.
// dag.mu must be held by the caller.
func checkReachable(dag *Dag, nodes []*Node) error {
	reachable := reachableFromStart(dag)
	for _, n := range nodes {
		if n.ID == StartNode {
			continue
		}
		if !reachable[n.ID] {
			return fmt.Errorf("node %q is not reachable from start_node", n.ID)
		}
	}
	return nil
}

// validateTopology performs all read-only structural checks before FinishDag commits
// any graph mutation. dag.mu must be held by the caller.
func (dag *Dag) validateTopology(nodes []*Node) error {
	if dag.startNode == nil {
		return fmt.Errorf("start node is missing; call StartDag or InitDag before FinishDag")
	}
	if err := checkIsolatedNodes(nodes); err != nil {
		return err
	}
	if err := checkReachable(dag, nodes); err != nil {
		return err
	}
	// 사이클 검사: EndNode 추가 전에 실행하여 실패 시 그래프를 변경하지 않는다.
	if detectCycle(dag) {
		return fmt.Errorf("FinishDag: %w", ErrCycleDetected)
	}
	return nil
}

// FinishDag finalizes the DAG by connecting end nodes and validating the structure.
func (dag *Dag) FinishDag() error {
	dag.mu.Lock()
	defer dag.mu.Unlock()

	// 에러 로그를 기록하고 반환하는 클로저 함수 (finishDag 에러 타입 사용)
	logErr := func(err error) error {
		dag.errLogs = append(dag.errLogs, &systemError{FinishDag, err})
		dag.reportError(err)
		return err
	}

	// Topology is frozen once GetReadyE has succeeded (nodeResult is non-nil).
	// Mutations are not permitted until Reset is called.
	if dag.nodeResult != nil {
		return logErr(fmt.Errorf("FinishDag: DAG topology is frozen after GetReady; call Reset before modifying the graph"))
	}
	if dag.validated {
		return logErr(fmt.Errorf("validated is already set to true"))
	}
	if len(dag.nodes) == 0 {
		return logErr(fmt.Errorf("no node"))
	}

	// 안전한 반복을 위해 노드 슬라이스를 생성 (맵의 구조 변경 방지를 위해)
	nodes := make([]*Node, 0, len(dag.nodes))
	for _, n := range dag.nodes {
		nodes = append(nodes, n)
	}

	// ── Phase 1: validate only — no graph mutation ────────────────────────────
	// All checks that can return an error run before any structural changes so
	// that a failed FinishDag leaves the DAG in exactly the state it was in
	// before the call.
	if err := dag.validateTopology(nodes); err != nil {
		return logErr(err)
	}

	// ── Phase 2: commit — graph mutation only after all checks pass ───────────

	// 종료 노드 생성 및 초기화
	dag.endNode = dag.createNode(EndNode)
	if dag.endNode == nil {
		return logErr(fmt.Errorf("failed to create end node"))
	}
	// TODO(endnode-cleanup): EndNode is unconditionally marked succeeded.
	// If a future runner type needs to release resources at DAG completion,
	// a cleanup hook should be invoked here before SetSucceed.
	dag.endNode.SetSucceed(true)

	// 자식이 없는 노드를 종료 노드와 연결
	for _, n := range nodes {
		if n.ID != EndNode && len(n.children) == 0 {
			if err := dag.addEndNode(n, dag.endNode); err != nil {
				return logErr(fmt.Errorf("addEndNode failed for node '%s': %w", n.ID, err))
			}
		}
	}

	dag.validated = true
	return nil
}

// reset is the internal implementation shared by Reset and ResetE.
// Caller must hold dag.mu.Lock() and must have verified dag.running == false.
func (dag *Dag) reset() {
	atomic.StoreInt64(&dag.completedCount, 0)
	atomic.StoreInt64(&dag.droppedErrors, 0)

	// Reset every node to Pending and clear its per-execution state.
	// childrenVertex / parentVertex are cleared here and rebuilt from dag.edges
	// below so they point to the newly created SafeChannels.
	for _, n := range dag.nodes {
		n.mu.Lock()
		n.status = NodeStatusPending
		n.succeed = false
		n.runner = nil
		n.childrenVertex = nil
		n.parentVertex = nil
		n.mu.Unlock()
	}

	// Re-attach the start-node trigger channel.  This channel is not backed by
	// an Edge entry; StartE() writes to parentVertex[0] to fire the first node.
	if dag.startNode != nil {
		dag.startNode.mu.Lock()
		dag.startNode.parentVertex = []*SafeChannel[runningStatus]{
			NewSafeChannelGen[runningStatus](dag.Config.MinChannelBuffer),
		}
		dag.startNode.mu.Unlock()
	}

	// Recreate each edge's channel and rewire it into the owning nodes' vertex
	// slices.  The lock order (dag.mu → n.mu) is consistent with the rest of
	// the codebase and prevents deadlocks.
	for _, edge := range dag.edges {
		edge.safeVertex = NewSafeChannelGen[runningStatus](dag.Config.MinChannelBuffer)
		if parent := dag.nodes[edge.parentID]; parent != nil {
			parent.mu.Lock()
			parent.childrenVertex = append(parent.childrenVertex, edge.safeVertex)
			parent.mu.Unlock()
		}
		if child := dag.nodes[edge.childID]; child != nil {
			child.mu.Lock()
			child.parentVertex = append(child.parentVertex, edge.safeVertex)
			child.mu.Unlock()
		}
	}

	// Recreate the aggregated result and error channels that Wait → closeChannels
	// closed at the end of the previous run.
	dag.NodesResult = NewSafeChannelGen[*printStatus](dag.Config.MaxChannelBuffer)
	dag.Errors = NewSafeChannelGen[error](dag.Config.MaxChannelBuffer)

	// Clear per-run state so GetReady re-initialises fresh on the next call.
	dag.nodeResult = nil
	dag.execSem = nil
	// nodeWg is a value type; its zero value is ready to use after all previous
	// goroutines have exited (guaranteed by Wait's defer nodeWg.Wait()).
	dag.errLogs = nil

	// Clear the one-per-run atomic flags so the next run can proceed.
	dag.started.Store(false)
	// startTrigger is re-captured by the next GetReadyE call after ConnectRunner
	// rewires the fresh parentVertex channels created above.
	dag.startTrigger = nil
}

// Reset reinitialises a completed DAG so it can be executed again without
// rebuilding the graph from scratch.
//
// Reset MUST be called only after Wait returns.  Calling it while the DAG is
// still running is a no-op (the call is logged and returns immediately) to
// prevent live goroutines from racing against freshly created channels.
// Use ResetE if you need an error return instead of a silent guard.
//
// After Reset, follow the standard execution lifecycle:
//
//	dag.Reset()
//	dag.ConnectRunner()
//	dag.GetReady(ctx)
//	dag.Start()
//	dag.Wait(ctx)
//
// The graph topology (nodes, edges, Config, ContainerCmd, runners) is preserved;
// only execution state is reset.
func (dag *Dag) Reset() {
	dag.mu.Lock()
	defer dag.mu.Unlock()
	if dag.running.Load() {
		Log.Error("Reset: called while DAG is running; ignoring — call Wait first")
		return
	}
	dag.reset()
}

// ResetE is the error-returning variant of Reset.  It returns an error when
// called while the DAG is still running (i.e. Wait has not yet returned).
func (dag *Dag) ResetE() error {
	dag.mu.Lock()
	defer dag.mu.Unlock()
	if dag.running.Load() {
		return fmt.Errorf("Reset: cannot reset while DAG is running; call Wait first")
	}
	dag.reset()
	return nil
}

// ConnectRunnerE attaches the three-phase runner closure (preFlight / inFlight /
// postFlight) to every node.  Call it after FinishDag and before GetReady; call
// it again after Reset so the closures reference the freshly created channels.
// Returns an error if the DAG has no nodes or if GetReady has already been called
// (runner closures must not be replaced while goroutines are live).
func (dag *Dag) ConnectRunnerE() error {
	dag.mu.Lock()
	defer dag.mu.Unlock()

	if dag.nodeResult != nil || dag.running.Load() {
		return fmt.Errorf("ConnectRunner: cannot reconnect runners after GetReady; call Reset first")
	}
	if len(dag.nodes) < 1 {
		return fmt.Errorf("ConnectRunner: DAG has no nodes")
	}
	for _, v := range dag.nodes {
		connectRunner(v)
	}
	return nil
}

// ConnectRunner is the bool-returning variant of ConnectRunnerE.
// It returns false if the DAG has no nodes.
func (dag *Dag) ConnectRunner() bool {
	if err := dag.ConnectRunnerE(); err != nil {
		Log.Warnf("%v", err)
		return false
	}
	return true
}

// GetReadyE initialises the execution semaphore and launches one goroutine per
// node.  Each goroutine runs preFlight (dependency wait) without holding an
// execution slot; the slot is acquired inside connectRunner just before inFlight
// (RunE) so that a small WorkerPoolSize cannot deadlock the graph.
//
// Prerequisites (returns a descriptive error if violated):
//   - ctx must be non-nil
//   - FinishDag must have been called (dag.validated == true)
//   - ConnectRunner must have been called (every node.runner must be non-nil)
//   - GetReady / GetReadyE must not have been called already
//
// ctx is forwarded to each node's runner; cancel it to abort the entire execution.
func (dag *Dag) GetReadyE(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("GetReady: ctx is nil")
	}

	dag.mu.Lock()
	defer dag.mu.Unlock()

	n := len(dag.nodes)
	if n < 1 {
		return fmt.Errorf("GetReady: DAG has no nodes")
	}
	if !dag.validated {
		return fmt.Errorf("GetReady: FinishDag has not been called")
	}
	if dag.startNode == nil || dag.endNode == nil {
		return fmt.Errorf("GetReady: start or end node is missing")
	}
	// Guard against double-call: a second call must not spawn duplicate goroutines.
	if dag.nodeResult != nil {
		return fmt.Errorf("GetReady: already called; call Reset before re-using the DAG")
	}
	// Every node must have a runner attached by ConnectRunner.  A nil runner
	// causes a nil function call panic inside the goroutine.
	for _, v := range dag.nodes {
		if v.runner == nil {
			return fmt.Errorf("GetReady: node %s has no runner; call ConnectRunner before GetReady", v.ID)
		}
	}

	// Validate startNode topology before launching any goroutines.
	// Structural invariants must be confirmed while the DAG is quiescent —
	// once goroutines are live, parentVertex / childrenVertex are immutable.
	// If this check fails the caller can fix the DAG state and retry;
	// no goroutines have been started yet so there is nothing to clean up.
	if len(dag.startNode.parentVertex) != 1 {
		return fmt.Errorf("GetReady: startNode has unexpected parentVertex length %d (want 1)",
			len(dag.startNode.parentVertex))
	}
	dag.startTrigger = dag.startNode.parentVertex[0]
	dag.execSem = make(chan struct{}, dag.Config.WorkerPoolSize)
	// Set running=true before launching goroutines so that Reset is blocked
	// from the moment the first goroutine is spawned.
	dag.running.Store(true)

	safeChs := make([]*SafeChannel[*printStatus], 0, n)

	for _, v := range dag.nodes {
		sc := NewSafeChannelGen[*printStatus](dag.Config.MinChannelBuffer)
		safeChs = append(safeChs, sc)
		dag.nodeWg.Add(1)
		go func(node *Node, sc *SafeChannel[*printStatus]) {
			defer dag.nodeWg.Done()
			node.runner(ctx, sc)
		}(v, sc)
	}

	dag.nodeResult = safeChs
	return nil
}

// GetReady is the bool-returning variant of GetReadyE.
func (dag *Dag) GetReady(ctx context.Context) bool {
	if err := dag.GetReadyE(ctx); err != nil {
		Log.Warnf("%v", err)
		return false
	}
	return true
}

// StartE fires the DAG by sending a trigger signal to the start-node trigger
// channel captured by GetReadyE.  All node goroutines are already waiting;
// this single send unblocks the start node and cascades through the graph.
//
// Call StartE exactly once after GetReady.  Returns an error if GetReady was
// not called, if Start was already called for this run, or if the trigger
// send fails unexpectedly.
//
// StartE never reads startNode.parentVertex directly — it uses the
// dag.startTrigger channel that GetReadyE captured and validated before any
// goroutines were launched.  This guarantees no concurrent access to the
// parentVertex slice while goroutines are live.
func (dag *Dag) StartE() error {
	// Snapshot nodeResult and startTrigger under read lock so we observe a
	// consistent view of GetReadyE's writes.  The lock is released before
	// Send so that dag.mu is never held during a channel operation.
	dag.mu.RLock()
	nodeResult := dag.nodeResult
	trigger := dag.startTrigger
	dag.mu.RUnlock()

	if nodeResult == nil || trigger == nil {
		return fmt.Errorf("Start: call GetReady before Start")
	}
	// Reject duplicate Start calls for the same run.
	if !dag.started.CompareAndSwap(false, true) {
		return fmt.Errorf("Start: already called; Start must be called exactly once per run")
	}
	if !trigger.Send(Start) {
		dag.started.Store(false)
		dag.startNode.SetSucceed(false)
		return fmt.Errorf("Start: failed to send trigger to start node")
	}
	dag.startNode.SetSucceed(true)
	return nil
}

// Start is the bool-returning variant of StartE.
func (dag *Dag) Start() bool {
	if err := dag.StartE(); err != nil {
		Log.Warnf("%v", err)
		return false
	}
	return true
}

// Wait blocks until the DAG finishes execution, ctx expires, or a fatal node
// failure is detected on NodesResult.
//
// It returns true only when the end node emits a FlightEnd status — meaning
// every node in the graph reached a terminal state (Succeeded or Skipped).
// It returns false on any of:
//   - context cancellation or timeout
//   - NodesResult channel closed unexpectedly
//   - end node reporting a PreflightFailed / InFlightFailed / PostFlightFailed
//
// Wait closes all channels (closeChannels) and shuts down the worker pool when
// it returns, regardless of whether execution succeeded.  Do NOT use the DAG
// after Wait returns without first calling Reset.
//
//nolint:gocognit // fan-in select loop must handle merge result, node status stream, and context cancellation simultaneously
func (dag *Dag) Wait(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	// Guard against calling Wait before GetReady (nodeResult not initialised).
	if dag.nodeResult == nil {
		return false
	}

	// Cleanup: three sequential steps inside one defer; order is significant.
	//   1. nodeWg.Wait    — block until every per-node goroutine has exited.
	//   2. closeChannels  — release all channel resources; must run before
	//                       running is cleared so Reset (which gates on
	//                       running==false) cannot create new channels while
	//                       closeChannels is still releasing old ones.
	//   3. running=false  — signal that Reset is now safe.
	//                       nodeResult stays non-nil until reset() clears it,
	//                       so the DAG remains frozen (setter guards reject
	//                       changes) until the caller explicitly calls Reset.
	defer func() {
		dag.nodeWg.Wait()
		dag.closeChannels()
		dag.running.Store(false)
	}()

	// 컨텍스트에서 타임아웃 설정
	var waitCtx context.Context
	var cancel context.CancelFunc
	if dag.bTimeout {
		waitCtx, cancel = context.WithTimeout(ctx, dag.Timeout)
	} else {
		waitCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	// mergeResult is a one-shot buffered channel: the goroutine writes exactly
	// once and then exits, so no close is ever needed.  The buffer of 1 prevents
	// the goroutine from blocking if Wait returns early (e.g. on ctx.Done).
	mergeResult := make(chan bool, 1)
	go func() {
		mergeResult <- dag.merge(waitCtx)
	}()

	for {
		select {
		case ok := <-mergeResult:
			// merge 함수가 완료되어 결과를 반환한 경우,
			// false 이면 병합 작업에 실패한 것이므로 false 리턴.
			if !ok {
				return false
			}
			// merge 함수 결과가 true 면 계속 진행한다.
		case c, ok := <-dag.NodesResult.GetChannel():
			if !ok {
				// 채널이 종료되면 실패 처리.
				return false
			}
			// EndNode 에 대한 상태만 체크함.
			if c.nodeID == EndNode {
				// Save fields before returning c to the pool, because
				// releasePrintStatus zeroes the struct fields.
				rStatus := c.rStatus
				releasePrintStatus(c)
				if rStatus == PreflightFailed ||
					rStatus == InFlightFailed ||
					rStatus == PostFlightFailed {
					return false
				}
				if rStatus == FlightEnd {
					return true
				}
			} else {
				releasePrintStatus(c)
			}
		case <-waitCtx.Done():
			Log.Printf("DAG execution timed out or canceled: %v", waitCtx.Err())
			return false
		}
	}
}

// WaitE is the error-returning variant of Wait.  It returns nil when the DAG
// completes successfully (FlightEnd from end_node), or a descriptive error on
// context cancellation, timeout, or node failure.
func (dag *Dag) WaitE(ctx context.Context) error {
	if !dag.Wait(ctx) {
		if ctx.Err() != nil {
			return fmt.Errorf("DAG execution cancelled: %w", ctx.Err())
		}
		return fmt.Errorf("DAG execution failed; inspect dag.Errors for node-level failures")
	}
	return nil
}

// dfsState holds the per-traversal book-keeping maps used by detectCycleDFS.
// Both maps are reused across calls via dfsStatePool to avoid per-call heap
// allocations; clear() resets entries without freeing the backing memory.
type dfsState struct {
	visited  map[string]bool
	recStack map[string]bool
}

// dfsStatePool provides reusable dfsState objects for detectCycle.
var dfsStatePool = sync.Pool{
	New: func() any {
		return &dfsState{
			visited:  make(map[string]bool),
			recStack: make(map[string]bool),
		}
	},
}

// reachableFromStart performs a forward DFS from dag.startNode and returns the
// set of reachable node IDs.  The caller must hold dag.mu in any mode.
// FinishDag calls this directly because it already holds dag.mu.Lock().
func reachableFromStart(dag *Dag) map[string]bool {
	visited := make(map[string]bool, len(dag.nodes))
	var dfs func(*Node)
	dfs = func(n *Node) {
		if n == nil || visited[n.ID] {
			return
		}
		visited[n.ID] = true
		for _, child := range n.children {
			dfs(child)
		}
	}
	dfs(dag.startNode)
	return visited
}

// detectCycleDFS detects cycles using DFS.
func detectCycleDFS(node *Node, visited, recStack map[string]bool) bool {
	if recStack[node.ID] {
		return true
	}
	if visited[node.ID] {
		return false
	}
	visited[node.ID] = true
	recStack[node.ID] = true

	for _, child := range node.children {
		if detectCycleDFS(child, visited, recStack) {
			return true
		}
	}

	recStack[node.ID] = false
	return false
}

// detectCycle checks if the DAG contains a cycle without acquiring a lock.
// The caller must guarantee that dag.nodes is not mutated during this call
// (e.g. by holding dag.mu in any mode).  FinishDag calls this directly because
// it already holds dag.mu.Lock(); the exported DetectCycle acquires RLock first.
func detectCycle(dag *Dag) bool {
	// Acquire a pooled dfsState and reset its maps without freeing backing
	// memory (clear is O(N) but allocation-free, unlike make).
	state := dfsStatePool.Get().(*dfsState)
	clear(state.visited)
	clear(state.recStack)
	defer dfsStatePool.Put(state)

	// Traverse dag.nodes directly — no copyDag needed because the caller
	// guarantees dag.mu is held (FinishDag holds Lock; DetectCycle holds RLock).
	for _, node := range dag.nodes {
		if !state.visited[node.ID] {
			if detectCycleDFS(node, state.visited, state.recStack) {
				return true
			}
		}
	}
	return false
}

// DetectCycle checks if the DAG contains a directed cycle.  It is safe to call
// concurrently: it acquires a read lock before inspecting the graph.
//
// Internal callers that already hold dag.mu (e.g. FinishDag) must call
// detectCycle directly to avoid a re-entrant lock attempt.
func DetectCycle(dag *Dag) bool {
	dag.mu.RLock()
	defer dag.mu.RUnlock()
	return detectCycle(dag)
}

// connectRunner connects a runner function to a node.
//
//nolint:gocognit,gocyclo // three-phase flight + semaphore; splitting would obscure the node lifecycle
func connectRunner(n *Node) {
	n.runner = func(ctx context.Context, result *SafeChannel[*printStatus]) {
		// sendResult delivers a copy of ps to the per-node result channel.
		// SendBlocking ensures the monitoring event is never silently dropped:
		// if the result buffer is momentarily full the send waits for a consumer.
		// If ctx is cancelled the pool-acquired copy is returned to prevent leak.
		sendResult := func(ps *printStatus) {
			copied := newPrintStatus(ps.rStatus, ps.nodeID)
			if !result.SendBlocking(ctx, copied) {
				releasePrintStatus(copied)
			}
		}

		// With FailFast (default), skip this node if any parent is already Failed.
		// CheckParentsStatus internally transitions the node to Skipped via CAS.
		// With ContinueOnError, all nodes run regardless of parent outcome.
		policy := ErrorPolicyFailFast
		if n.parentDag != nil {
			policy = n.parentDag.Config.ErrorPolicy
		}
		if policy == ErrorPolicyFailFast && !n.CheckParentsStatus() {
			ps := newPrintStatus(PostFlightFailed, n.ID)
			sendResult(ps)
			n.notifyChildren(ctx, Failed)
			releasePrintStatus(ps)
			return
		}

		// Pending → Running: guarded CAS — rejects the transition if the node is
		// not Pending (e.g. already Skipped by a concurrent parent failure).
		if ok := n.TransitionStatus(NodeStatusPending, NodeStatusRunning); !ok {
			Log.Warnf("connectRunner: Pending→Running rejected for node %s (status=%v)", n.ID, n.GetStatus())
		}

		// preFlight phase — waits for all parent signals using the caller ctx only.
		// No execution budget is consumed here: dependency wait must not reduce
		// the time available for the node's own work (inFlight).
		ps := preFlight(ctx, n)
		sendResult(ps)
		if ps.rStatus == PreflightFailed {
			if ok := n.TransitionStatus(NodeStatusRunning, NodeStatusFailed); !ok {
				Log.Warnf("connectRunner: Running→Failed rejected for node %s (status=%v)", n.ID, n.GetStatus())
			}
			n.notifyChildren(ctx, Failed)
			releasePrintStatus(ps)
			return
		}
		releasePrintStatus(ps)

		// Acquire the execution semaphore slot.  Only inFlight (RunE) holds the
		// slot; preFlight runs without it, so a WorkerPoolSize smaller than the
		// node count can never cause a deadlock.
		// The timeout context is built AFTER acquiring the slot so the budget
		// measures actual RunE time, not semaphore wait time.
		if pd := n.parentDag; pd != nil && pd.execSem != nil {
			select {
			case pd.execSem <- struct{}{}:
				defer func() { <-pd.execSem }()
			case <-ctx.Done():
				semPs := newPrintStatus(InFlightFailed, n.ID)
				sendResult(semPs)
				if ok := n.TransitionStatus(NodeStatusRunning, NodeStatusFailed); !ok {
					Log.Warnf("connectRunner: Running→Failed (semaphore ctx cancel) rejected for node %s (status=%v)", n.ID, n.GetStatus())
				}
				n.notifyChildren(ctx, Failed)
				releasePrintStatus(semPs)
				return
			}
		}

		// Resolve execution context for inFlight.
		// Timeout priority (highest first):
		//   1. Node.Timeout  — when Timeout > 0 (explicit per-node override)
		//   2. Dag.Config.DefaultTimeout — when positive (implicit DAG-wide default)
		//   3. Caller ctx only — no additional deadline
		//
		// DefaultTimeout == 0 means no implicit per-node execution timeout is
		// applied; DAG-wide limits are enforced via the caller context.
		// Snapshot Node.Timeout under the read lock to avoid a data race with
		// callers that set Timeout before ConnectRunner (e.g. CreateNodeWithTimeOut).
		n.mu.RLock()
		nodeTimeout := n.Timeout
		n.mu.RUnlock()

		var execCtx context.Context
		var execCancel context.CancelFunc
		switch {
		case nodeTimeout > 0:
			execCtx, execCancel = context.WithTimeout(ctx, nodeTimeout)
		case n.parentDag != nil && n.parentDag.Config.DefaultTimeout > 0:
			execCtx, execCancel = context.WithTimeout(ctx, n.parentDag.Config.DefaultTimeout)
		default:
			execCtx, execCancel = context.WithCancel(ctx)
		}
		defer execCancel()

		// inFlight phase — runs node-owned work (RunE) with the execution context.
		ps = inFlight(execCtx, n)
		sendResult(ps)
		if ps.rStatus == InFlightFailed {
			if ok := n.TransitionStatus(NodeStatusRunning, NodeStatusFailed); !ok {
				Log.Warnf("connectRunner: Running→Failed rejected for node %s (status=%v)", n.ID, n.GetStatus())
			}
			n.notifyChildren(ctx, Failed)
			releasePrintStatus(ps)
			return
		}
		releasePrintStatus(ps)

		// postFlight phase — ctx forwarded so SendBlocking can abort on cancel.
		ps = postFlight(ctx, n)
		sendResult(ps)
		if ps.rStatus == PostFlightFailed {
			if ok := n.TransitionStatus(NodeStatusRunning, NodeStatusFailed); !ok {
				Log.Warnf("connectRunner: Running→Failed rejected for node %s (status=%v)", n.ID, n.GetStatus())
			}
		} else {
			if ok := n.TransitionStatus(NodeStatusRunning, NodeStatusSucceeded); !ok {
				Log.Warnf("connectRunner: Running→Succeeded rejected for node %s (status=%v)", n.ID, n.GetStatus())
			}
		}
		releasePrintStatus(ps)
	}
}

// merge merges all status channels using fan-in pattern.
func (dag *Dag) merge(ctx context.Context) bool {
	// 입력 SafeChannel 슬라이스가 비어 있으면 그대로 종료함.
	if len(dag.nodeResult) < 1 {
		return false
	}

	return fanIn(ctx, dag.nodeResult, dag.NodesResult)
}

// fanIn merges multiple channels into one using errgroup for structured cancellation.
//
// No concurrency limit is applied: the relay goroutines are I/O-bound (blocked
// on channel reads) and must all run concurrently.  Adding errgroup.SetLimit
// would serialize completions and increase tail latency without bounding CPU.
func fanIn(ctx context.Context, channels []*SafeChannel[*printStatus], merged *SafeChannel[*printStatus]) bool {
	eg, egCtx := errgroup.WithContext(ctx)

	for _, sc := range channels {
		eg.Go(func() error {
			for val := range sc.GetChannel() {
				if !merged.SendBlocking(egCtx, val) {
					return egCtx.Err()
				}
			}
			return nil
		})
	}

	return eg.Wait() == nil
}

// min returns the minimum of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ToMermaid generates a Mermaid flowchart string that represents the DAG topology.
//
// The output uses the "graph TD" (top-down) direction.  Synthetic nodes
// (start_node / end_node) are rendered with a stadium shape to distinguish
// infrastructure nodes from user-defined ones.  If a per-node Runnable has been
// registered via SetNodeRunner, its concrete Go type is appended to the node
// label after a line-break so the diagram shows which executor is bound to each
// step — useful for debugging or documentation.
//
// Node IDs are sanitised for Mermaid syntax by replacing any character that is
// not an ASCII letter, digit, or underscore with an underscore (see
// mermaidSafeID).  This prevents parser errors for IDs that contain hyphens,
// dots, or spaces.
//
// Example output:
//
//	graph TD
//	    start_node(["start_node"])
//	    A["A\n*main.MyRunner"]
//	    B["B"]
//	    end_node(["end_node"])
//	    start_node --> A
//	    A --> B
//	    B --> end_node
//
// ToMermaid acquires a read-lock and is safe to call concurrently with
// Progress() and other read-only observers.  It must be called after
// FinishDag so that dag.edges is complete; calling it before FinishDag will
// produce a diagram that is missing the edges to the synthetic end node.
func (dag *Dag) ToMermaid() string {
	dag.mu.RLock()
	defer dag.mu.RUnlock()

	// Build a collision-free ID map covering every node in the graph.
	allIDs := make([]string, 0, len(dag.nodes))
	for id := range dag.nodes {
		allIDs = append(allIDs, id)
	}
	idMap := buildMermaidIDMap(allIDs)

	var sb strings.Builder
	sb.WriteString("graph TD\n")

	// First pass — emit a labelled node definition for every node that appears
	// in at least one edge.
	defined := make(map[string]bool, len(dag.nodes))
	for _, edge := range dag.edges {
		dag.writeMermaidNode(&sb, edge.parentID, defined, idMap)
		dag.writeMermaidNode(&sb, edge.childID, defined, idMap)
	}

	// Second pass — emit directed edges using the collision-free IDs.
	for _, edge := range dag.edges {
		fmt.Fprintf(&sb, "    %s --> %s\n", idMap[edge.parentID], idMap[edge.childID])
	}

	return sb.String()
}

// writeMermaidNode emits a single Mermaid node definition into sb unless the
// node has already been defined (tracked via the defined map).
// Synthetic nodes (start_node / end_node) use the stadium shape; all others
// use the default rectangle.  Caller must hold dag.mu at least for reading.
func (dag *Dag) writeMermaidNode(sb *strings.Builder, nodeID string, defined map[string]bool, idMap map[string]string) {
	if defined[nodeID] {
		return
	}
	defined[nodeID] = true
	mID := idMap[nodeID]
	switch nodeID {
	case StartNode, EndNode:
		fmt.Fprintf(sb, "    %s([\"%s\"])\n", mID, mermaidEscapeLabel(nodeID))
	default:
		fmt.Fprintf(sb, "    %s[\"%s\"]\n", mID, mermaidNodeLabel(nodeID, dag.nodes[nodeID]))
	}
}

// mermaidNodeLabel returns the display label for a Mermaid node box.
// When a per-node Runnable has been registered via SetNodeRunner, its concrete
// Go type is appended after a line-break as a diagnostic hint visible in the
// rendered diagram.
func mermaidNodeLabel(nodeID string, n *Node) string {
	if n != nil {
		if r := n.runnerLoad(); r != nil {
			return fmt.Sprintf("%s\\n%T", mermaidEscapeLabel(nodeID), r)
		}
	}
	return mermaidEscapeLabel(nodeID)
}

// buildMermaidIDMap returns a map from each original node ID to a unique,
// Mermaid-safe identifier.  IDs that would produce the same safe base are
// disambiguated by appending a numeric suffix (_0, _1, …).
func buildMermaidIDMap(ids []string) map[string]string {
	// First pass: compute the naive safe ID for each original ID and count
	// how many originals share each safe base.
	bases := make(map[string]string, len(ids))  // original → base safe ID
	baseCount := make(map[string]int, len(ids)) // base safe ID → occurrence count
	for _, id := range ids {
		b := mermaidSafeBase(id)
		bases[id] = b
		baseCount[b]++
	}

	// Second pass: assign final IDs, adding a counter suffix only for
	// bases that are shared by more than one original.
	result := make(map[string]string, len(ids))
	counter := make(map[string]int, len(ids))
	for _, id := range ids {
		b := bases[id]
		if baseCount[b] > 1 {
			result[id] = fmt.Sprintf("%s_%d", b, counter[b])
			counter[b]++
		} else {
			result[id] = b
		}
	}
	return result
}

// mermaidSafeBase converts an arbitrary node ID to a Mermaid-safe base
// identifier by replacing any character that is not an ASCII letter, digit,
// or underscore with an underscore.
func mermaidSafeBase(id string) string {
	var b strings.Builder
	b.Grow(len(id))
	for _, c := range id {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			b.WriteRune(c)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

// mermaidEscapeLabel escapes characters that would break the Mermaid quoted
// label syntax: double-quotes become #quot; (the Mermaid HTML entity),
// backslashes are doubled, and carriage-returns are stripped.
// Embedded newlines are kept as the Mermaid literal "\n" sequence so the
// renderer can display multi-line labels.
func mermaidEscapeLabel(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "#quot;")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

// CopyDag creates a fully executable copy of original with a new ID.
//
// The copy preserves the graph topology (nodes, edges, parent/child
// relationships) and all DagConfig values.  Fresh execution channels
// (NodesResult, Errors, edge safeVertices, start-node trigger) are created so
// running the copy cannot affect the original.
//
// Items NOT copied: per-node runners (Node.runnerVal), execSem, nodeResult,
// startTrigger, errLogs, and runnerResolver.  startTrigger is intentionally
// omitted — it is nil until GetReadyE validates and captures it; the copy must
// go through GetReady independently.  Per-node runners are excluded because
// they are stateful closures bound to the original DAG's channels.
// The caller must wire runners and follow the standard lifecycle on the copy:
//
//	ConnectRunner → GetReady → Start → Wait
//
// Returns nil if original is nil or newID is empty.
func CopyDag(original *Dag, newID string) *Dag {
	if original == nil {
		return nil
	}
	if utils.IsEmptyString(newID) {
		return nil
	}

	copied := &Dag{
		ID:           newID,
		Timeout:      original.Timeout,
		bTimeout:     original.bTimeout,
		ContainerCmd: original.ContainerCmd,
		validated:    original.validated,
		Config:       original.Config,
		NodesResult:  NewSafeChannelGen[*printStatus](original.Config.MaxChannelBuffer),
		Errors:       NewSafeChannelGen[error](original.Config.MaxChannelBuffer),
	}

	// Copy graph topology (node IDs + parent/child pointers).
	newNodes, newEdges := copyDag(original)
	copied.nodes = newNodes
	copied.edges = newEdges

	// Set parentDag and identify synthetic nodes.
	for _, node := range newNodes {
		node.parentDag = copied
		switch node.ID {
		case StartNode:
			copied.startNode = node
		case EndNode:
			copied.endNode = node
		}
	}

	// Preserve user-node count so Progress() and MarkCompleted() work correctly.
	atomic.StoreInt64(&copied.nodeCount, atomic.LoadInt64(&original.nodeCount))

	// Wire execution channels — mirrors the channel-wiring half of Reset().
	// Without this step the copy has no parentVertex / childrenVertex entries:
	// preFlight immediately passes for every node (no channels to wait on) and
	// Start() returns false (startNode.parentVertex is empty).
	if copied.startNode != nil {
		copied.startNode.parentVertex = []*SafeChannel[runningStatus]{
			NewSafeChannelGen[runningStatus](copied.Config.MinChannelBuffer),
		}
	}
	for _, edge := range copied.edges {
		edge.safeVertex = NewSafeChannelGen[runningStatus](copied.Config.MinChannelBuffer)
		if parent := copied.nodes[edge.parentID]; parent != nil {
			parent.childrenVertex = append(parent.childrenVertex, edge.safeVertex)
		}
		if child := copied.nodes[edge.childID]; child != nil {
			child.parentVertex = append(child.parentVertex, edge.safeVertex)
		}
	}

	return copied
}

// copyDag dag 를 복사함. 필요한것만 복사함.
func copyDag(original *Dag) (map[string]*Node, []*Edge) {
	// 원본에 노드가 없으면 nil 반환
	if len(original.nodes) == 0 {
		return nil, nil
	}

	// 1. 노드의 기본 정보(Id)만 복사한 새 맵 생성
	newNodes := make(map[string]*Node, len(original.nodes))
	for _, n := range original.nodes {
		// 필요한 최소한의 정보만 복사
		newNode := &Node{
			ID: n.ID, // Node 구조체가 Id로 되어 있다면 그대로 유지
			// 기타 필드는 cycle 검증에 필요하지 않으므로 생략
		}
		newNodes[newNode.ID] = newNode
	}

	// 2. 원본 노드의 부모/자식 관계를 이용하여 새 노드들의 포인터 연결
	for _, n := range original.nodes {
		newNode := newNodes[n.ID]
		// Pre-allocate slices to the exact capacity needed so each append
		// below copies at most one node pointer without triggering a grow.
		newNode.parent = make([]*Node, 0, len(n.parent))
		newNode.children = make([]*Node, 0, len(n.children))
		// 부모 노드 연결
		for _, parent := range n.parent {
			if copiedParent, ok := newNodes[parent.ID]; ok {
				newNode.parent = append(newNode.parent, copiedParent)
			}
		}
		// 자식 노드 연결
		for _, child := range n.children {
			if copiedChild, ok := newNodes[child.ID]; ok {
				newNode.children = append(newNode.children, copiedChild)
			}
		}
	}

	// 3. 간선(Edge) 복사: parentID, childID만 복사
	newEdges := make([]*Edge, len(original.edges))
	for i, e := range original.edges {
		newEdges[i] = &Edge{
			parentID: e.parentID,
			childID:  e.childID,
			// vertex 등 기타 정보는 cycle 검증에 필요하지 않으므로 생략
		}
	}

	return newNodes, newEdges
}

func edgeExists(edges []*Edge, parentID, childID string) bool {
	for _, edge := range edges {
		if edge.parentID == parentID && edge.childID == childID {
			return true
		}
	}
	return false
}
