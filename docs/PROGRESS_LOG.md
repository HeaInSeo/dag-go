# dag-go Project Progress Log

This file tracks the incremental improvement stages of the dag-go library.
Update this file at the start and end of every stage.

---

## Current Status: Stage 11 — Data Plane Performance Optimization (completed)

**Branch:** `main`
**Last updated:** 2026-02-27

---

## Completed Items

### Stage 1 — Analysis Report (read-only)
- File system cleanup candidates identified (-bak files, .travis.yml, ommit-test/).
- Library purity audit: ommit-test/ had no direct K8s deps but structural risk; recommended removal.
- Concurrency rule review: bare `chan error` for Dag.Errors, missing goleak, RunCommand dual-state.
- Top-3 TODOs prioritised: context propagation to inFlight, preFlight 30 s hardcode, RunCommand removal.

### Stage 2 — Core Refactoring
- Deleted: `ommit-test/`, `dag.go-bak`, `node.go-bak`, `.travis.yml`.
- Moved: `backup.md`, `atomic_Value.md` → `docs/`.
- `Runnable.RunE` signature extended: `RunE(ctx context.Context, a interface{}) error`.
- `Dag.Errors chan error` replaced with `*SafeChannel[error]` (double-close safe).
- `RunCommand Runnable` field removed; all runner access unified through `runnerVal atomic.Value`.
- `inFlight`, `Execute`, `execute` functions now accept and forward `context.Context`.
- `preFlight` timeout now respects priority: Node.Timeout > DagConfig.DefaultTimeout > caller deadline.
- `SetRunnerUnsafe` removed (bypassed atomic path).
- **Bug fixed:** `SetNodeRunner` caused re-entrant mutex self-deadlock; fixed by calling `runnerStore` directly instead of `SetRunner`.
- `goleak.VerifyNone(t)` added to `runner_test.go:Test_lateBinding`.
- All 21 tests passing; lint 0 issues. Commit: `d55ed59`.

### Stage 3 — Lint Config & Dependency Cleanup
- `.golangci.yml`: govet analyser name corrected `copylock` → `copylocks` (CI fix).
- `go mod tidy`: removed unused `github.com/dlsniper/debugger` direct dependency.
- Lint: 0 issues. All tests passing. Commit: `707bd29`.

### Stage 4 — Cycle Detection Completion
- **`errors.go`**: Added `ErrCycleDetected` sentinel error (`errors.New`).
- **`dag.go` (`FinishDag`)**: Cycle detection now returns `fmt.Errorf("FinishDag: %w", ErrCycleDetected)`,
  enabling callers to use `errors.Is(err, ErrCycleDetected)`.
- **`dag_test.go`**: Four cycle detection test cases added:
  - `TestDetectCycle_SimpleCycle`: start → A → B → A (2-node cycle via FinishDag).
  - `TestDetectCycle_ComplexCycle`: start → A → B → C → A (3-node cycle via FinishDag).
  - `TestDetectCycle_SelfLoop`: A → A (self-loop via direct graph construction + DetectCycle).
  - `TestDetectCycle_NoCycle`: diamond DAG — verifies FinishDag and DetectCycle both return clean.
- All tests passing; lint 0 issues.

### Stage 5 — Node Status Atomicity: CAS Pattern
- **`node.go`**: Added `isValidTransition(from, to NodeStatus) bool` state-machine helper.
- **`node.go`**: Added `TransitionStatus(from, to NodeStatus) bool` — mutex-protected CAS:
  - Validates the from→to edge against the state machine before acquiring the lock.
  - Acquires `n.mu.Lock()`, checks `n.status == from`, writes `n.status = to`, returns true.
  - Returns false (and leaves status unchanged) on invalid edge or wrong pre-condition.
- **`node.go` (`CheckParentsStatus`)**: Replaced `SetStatus(Skipped)` with `TransitionStatus(Pending, Skipped)`.
- **`node.go` (`inFlight`)**: Removed redundant `n.SetStatus(NodeStatusRunning)` (already set by `connectRunner`).
- **`dag.go` (`connectRunner`)**:
  - Removed redundant `n.SetStatus(NodeStatusSkipped)` after `!CheckParentsStatus()`.
  - All six `SetStatus` calls replaced with `TransitionStatus` + `Log.Warnf` on rejection.
  - State machine transitions enforced: Pending→Running, Running→{Succeeded,Failed}.
- **`node_test.go`**: Four new test functions (13 sub-tests + concurrent scenarios):
  - `TestTransitionStatus_ValidTransitions`: all 4 valid state-machine edges.
  - `TestTransitionStatus_InvalidTransitions`: 6 illegal transitions — all blocked.
  - `TestTransitionStatus_ConcurrentPendingToRunning`: 200 goroutines race; exactly 1 wins.
  - `TestTransitionStatus_ConcurrentFullLifecycle`: 3-phase race (Pending→Running→terminal).
- `go test -race ./...` — **0 data races detected**.
- `golangci-lint run ./...` — **0 issues**.

### Stage 6 — README & Error Handling Hardening
- **`README.md`**: Complete rewrite — Introduction, Key Features table, Quick Start code example,
  DAG Lifecycle diagram, Node State Machine (Mermaid), Configuration reference,
  Per-Node Runner Override section, Error Handling section, Development guide.
  Travis CI badge removed (`.travis.yml` was deleted in Stage 2).
  pkg.go.dev link promoted as the canonical API reference.
- **`dag.go` (`DagConfig`)**: Added `ErrorDrainTimeout time.Duration` field.
- **`dag.go` (`DefaultDagConfig`)**: Sets `ErrorDrainTimeout: 5 * time.Second` as default.
- **`dag.go` (`collectErrors`)**: Removed hardcoded `5 * time.Second` timeout.
  Now reads `dag.Config.ErrorDrainTimeout`; falls back to 5 s when field is zero.
- **`dag.go` (`reportError`)**: Replaced `Log.Printf` with structured logrus entry:
  `Log.WithField("dag_id", ...).WithError(err).Warn(...)` for log-aggregation compatibility.
- Remaining Stage-6 items deferred (ErrNoRunner structured type, ErrCount helper, ErrorPolicy enum).
- `go test -race ./...` — **0 data races detected**.
- `golangci-lint run ./...` — **0 issues**.

### Stage 7 — Performance Optimization & Internal Concurrency Refinement (current)
- **`dag.go` (`DagWorkerPool`)**: Added `closeOnce sync.Once` field.
  `Close()` now wraps `close(taskQueue)` in `closeOnce.Do(...)` — double-close panic eliminated.
  `Submit()` retains blocking send semantics (backpressure); converting to `SafeChannel` would
  silently drop tasks on a full queue and was therefore rejected.
- **`dag.go` (`fanIn`)**: Refactored from `sync.WaitGroup + atomic.int32` to `errgroup.WithContext`.
  Added `"golang.org/x/sync/errgroup"` import.  No `SetLimit` applied: relay goroutines are
  I/O-bound (blocked on channel reads) and must all run concurrently; serialising them would
  increase tail latency with no CPU benefit.
- **`dag.go` (`Wait`)**: Added documentation comment to `mergeResult` channel explaining its
  one-shot buffered safety: one writer, one reader, goroutine exits after writing — no close needed.
- **`dag.go` (`Progress`)**: Extended godoc to document that the two `atomic.LoadInt64` calls
  (nodeCount + completedCount) do not form an atomic pair.  The ratio may be slightly ahead of
  reality between loads; acceptable for observability but not for correctness decisions.
- **`dag.go` (`detectCycle` / `DetectCycle`)**: Split into two functions:
  - `detectCycle(dag *Dag) bool` — internal, no lock; callers must hold `dag.mu`.
  - `DetectCycle(dag *Dag) bool` — exported, acquires `dag.mu.RLock()` then delegates.
  - `FinishDag` updated to call `detectCycle(dag)` directly (already holds write lock),
    preventing a re-entrant lock attempt that would have deadlocked.
- `go test -race ./...` — **0 data races detected**.
- `golangci-lint run ./...` — **0 issues**.

### Stage 10 — Examples, Visualisation & v1.0.0 Release Preparation (completed)
- **`dag.go` (`ToMermaid()`)**: New exported method — generates a Mermaid `graph TD` flowchart
  string from `dag.Edges`.  Synthetic nodes use stadium shape (`([" "])`); per-node runner type
  appended to label when set via `SetNodeRunner`.  `mermaidSafeID()` helper sanitises node IDs.
  `"strings"` import added.
- **`examples/etl_pipeline/main.go`**: Fan-out / fan-in ETL pipeline.  3 parallel ingest nodes →
  validate → transform → load → report.  Per-node `SetNodeRunners` assignment.  Scenario 1 (all
  succeed) and Scenario 2 (ingest_api fails → downstream Skipped) in a single binary.
- **`examples/build_system/main.go`**: Multi-level build dependency graph.  3 parallel lib compiles →
  compile_app → {unit_test, integ_test} → package_app → deploy.  Bulk `SetNodeRunners` with
  typed `CompileRunner`/`TestRunner`/`PackageRunner`/`DeployRunner`.  `ToMermaid()` output shown.
- **`examples/reset_retry/main.go`**: Transient-failure retry strategy.  `TransientRunner` with
  atomic counter fails on attempt #1, succeeds on #2.  Demonstrates the full
  `Reset() → ConnectRunner → GetReady → Start → Wait` reuse lifecycle.
- **`docs/RELEASE_v1.0.0.md`**: Full release note covering API surface, lifecycle, state machine,
  configuration reference, runner priority, examples, and known limitations.
- `go build ./...` — **0 errors** (all examples compile as part of the same module).
- `go test -race ./...` — **0 data races, all tests pass**.

### Stage 11 — Data Plane Performance Optimization (completed)
- **`dag.go` (`detectCycle`)**: Eliminated `copyDag` call — now uses `dfsStatePool` (`sync.Pool` of
  `dfsState{visited, recStack map[string]bool}`) and traverses `dag.nodes` directly.
  `clear()` (Go 1.21+) resets maps without freeing backing memory.
  Caller lock guarantee documented: `FinishDag` holds `Lock()`; `DetectCycle` holds `RLock()`.
- **`dag.go` (`dfsState` / `dfsStatePool`)**: New pooled DFS traversal state added before
  `detectCycleDFS`.  Both maps pre-allocated at pool-object creation; reused via `clear()`.
- **`dag.go` (`connectRunner` — `copyStatus`)**: Changed from direct `&printStatus{...}` heap
  allocation to `newPrintStatus(ps.rStatus, ps.nodeID)` — now routes through `statusPool`.
- **`dag.go` (`Wait`)**: Added `releasePrintStatus(c)` in all consuming paths of the
  `NodesResult` case.  `c.rStatus` saved to local variable before release to avoid use-after-free.
  Non-EndNode statuses also released, closing the pool lifecycle loop.
- **`dag.go` (`DagConfig`)**: Added `ExpectedNodeCount int` field with godoc.
- **`dag.go` (`NewDagWithConfig`)**: Uses `ExpectedNodeCount` as capacity hint for
  `make(map[string]*Node, n)` and `make([]*Edge, 0, n)` — avoids rehash growth for known graphs.
- **Benchmark results** (before → after on Intel Xeon E5-2683 v4 @ 2.10GHz):

| Benchmark | Before (ns/op) | After (ns/op) | Before (allocs/op) | After (allocs/op) |
|---|---|---|---|---|
| DetectCycle/Small  |  7 239 |  1 896 | 50    | **0** |
| DetectCycle/Medium | 206 041 |  53 976 | 1 279 | **0** |
| DetectCycle/Large  | 5 822 780 | 1 558 478 | 28 011 | **0** |

- `go test -race ./...` — **0 data races, all tests pass**.
- `golangci-lint run ./...` — **0 issues**.

---

## Pending Items

### Stage 8 — Error Handling Continuation (planned)
- [ ] `ErrNoRunner` structured error type (currently `errors.New`); add `NodeID` field via `NodeError`.
- [ ] Add `Dag.ErrCount()` helper to expose current Errors channel depth without draining.
- [ ] Consider `ErrorPolicy` enum: `FailFast` (current) vs `ContinueOnError`.

### Stage 9 — Public API Hardening (completed)
- [x] `CopyDag` fixed: now copies `Config` (all DagConfig fields) and allocates a new `Errors` channel so the copy is fully independent from the original.
- [x] Godoc pass: all exported types, fields, and functions in `dag.go` now have comprehensive godoc comments (construction order, threading guarantees, usage caveats).
- [x] `Dag.Reset()` implemented — resets all nodes to Pending, recreates edge channels, rebuilds vertex wiring, and recreates the aggregation channels; allows full DAG reuse after `Wait` without graph rebuilding.
- [x] `TestDagReset`: two-run integration test verifying `Progress()` → 1.0 first run, `Progress()` → 0.0 after Reset, all nodes Pending, `Progress()` → 1.0 second run.
- `workerPool` intentionally NOT copied by `CopyDag` (documented in godoc); each copy must call its own `GetReady`.

### Stage 10 — Profiling & Advanced Optimisation (deferred to v1.1.0)
- [ ] Profile `TestComplexDag` under race detector; identify hot paths.
- [ ] Evaluate bounded semaphore for `preFlight` errgroup (currently hard-coded limit of 10).

> Stage 10 (this session) was renamed to "Examples, Visualisation & v1.0.0 Release Preparation"
> and is now **completed** (see Completed Items above).  The profiling work is deferred.

---

## Architecture Notes

| Component | Pattern | Key Invariant |
|---|---|---|
| `Node.runnerVal` | `atomic.Value` wrapping `*runnerSlot` | Store-type must never change; always use `runnerStore`/`runnerLoad` |
| `Dag.Errors` | `*SafeChannel[error]` | Closed centrally in `closeChannels()`; `Send` is non-blocking |
| Lock order | `Dag.mu` → `Node.mu` | Never invert; `SetNodeRunner` calls `runnerStore` directly (not `SetRunner`) to avoid re-entrant lock |
| Runner priority | `Node.runnerVal` > `runnerResolver` > `ContainerCmd` | Resolved at execution time in `getRunnerSnapshot` |
| Cycle detection | DFS + `dfsStatePool` (zero-alloc) | `detectCycle` iterates `dag.nodes` directly; `dfsStatePool` provides reusable `visited`/`recStack` maps; `clear()` resets without GC pressure |
| Status transitions | `TransitionStatus(from, to)` CAS | Validates state-machine edge before lock; `SetStatus` kept for unconditional override only |
| Error observability | `reportError` uses logrus fields | `dag_id` + `error` fields emitted on drop; `collectErrors` uses `DagConfig.ErrorDrainTimeout` |
| `DagWorkerPool.Close` | `sync.Once` | `closeOnce.Do(close)` prevents double-close panic; `Submit` keeps blocking send for backpressure |
| `fanIn` | `errgroup.WithContext` | No `SetLimit`: I/O-bound relay goroutines must run concurrently; cancellation via `egCtx.Done()` |
| `DetectCycle` / `detectCycle` | Lock-split pattern | `detectCycle` (internal, no lock) called by `FinishDag`; `DetectCycle` (exported) acquires `RLock` |
| `Progress()` | Two independent `atomic.LoadInt64` | Not an atomic pair; ratio may be slightly ahead between reads; acceptable for observability only |
| `Dag.Reset()` | Channel recreation + node state wipe | Must be called only after `Wait` returns; rebuilds edge channels from `dag.Edges` slice so topology is preserved |
| `CopyDag` | Structural copy | Copies `Config` and allocates new `NodesResult`/`Errors` channels; `workerPool`/`nodeResult`/`errLogs` intentionally not copied |
| `Dag.ToMermaid()` | Read-only observer | Acquires `dag.mu.RLock()`; emits `graph TD` Mermaid; sanitises node IDs via `mermaidSafeID()`; appends `%T` runner hint when per-node runner is set |
| `printStatus` pool | `statusPool` (`sync.Pool`) | `newPrintStatus` → `copyStatus` (connectRunner) → `result.Send` → `releasePrintStatus` (Wait); full lifecycle closed in Stage 11 |
| `dfsStatePool` | `sync.Pool` of `dfsState` | Provides zero-alloc DFS traversal for `detectCycle`; `clear()` resets maps per call |
| `DagConfig.ExpectedNodeCount` | Capacity hint | Pre-allocates `nodes` map and `Edges` slice in `NewDagWithConfig`; zero = runtime default |
