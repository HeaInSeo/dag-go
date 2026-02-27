# dag-go Project Progress Log

This file tracks the incremental improvement stages of the dag-go library.
Update this file at the start and end of every stage.

---

## Current Status: Stage 6 — README & Error Handling Hardening

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

### Stage 6 — README & Error Handling Hardening (current)
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

---

## Pending Items

### Stage 7 — Error Handling Continuation (planned)
- [ ] `ErrNoRunner` structured error type (currently `errors.New`); add `NodeID` field via `NodeError`.
- [ ] Add `Dag.ErrCount()` helper to expose current Errors channel depth without draining.
- [ ] Consider `ErrorPolicy` enum: `FailFast` (current) vs `ContinueOnError`.

### Stage 8 — Performance & Concurrency (planned)
- [ ] `DagWorkerPool.taskQueue` is a bare `chan func()` — evaluate replacing with `SafeChannel`.
- [ ] `fanIn` goroutine: each node spawns a goroutine; explore bounded semaphore approach.
- [ ] `merge` result channel (`mergeResult chan bool`) — evaluate promotion to `SafeChannel[bool]`.
- [ ] Profile `TestComplexDag` under race detector; verify no data races.
- [ ] `Progress()` uses two independent `atomic.LoadInt64` calls — not atomic as a pair; document limitation.

### Stage 9 — Public API Hardening (planned)
- [ ] `DetectCycle` currently takes a `*Dag` with no lock; document that it must be called only when DAG mutation is complete (post-FinishDag) or add internal locking.
- [ ] `CopyDag` does not copy `Config` or `workerPool`; document this or fix.
- [ ] Godoc pass: all exported symbols must have comments (revive:exported rule already enforced).
- [ ] Evaluate adding `Dag.Reset()` to allow DAG reuse after `Wait` completes.

---

## Architecture Notes

| Component | Pattern | Key Invariant |
|---|---|---|
| `Node.runnerVal` | `atomic.Value` wrapping `*runnerSlot` | Store-type must never change; always use `runnerStore`/`runnerLoad` |
| `Dag.Errors` | `*SafeChannel[error]` | Closed centrally in `closeChannels()`; `Send` is non-blocking |
| Lock order | `Dag.mu` → `Node.mu` | Never invert; `SetNodeRunner` calls `runnerStore` directly (not `SetRunner`) to avoid re-entrant lock |
| Runner priority | `Node.runnerVal` > `runnerResolver` > `ContainerCmd` | Resolved at execution time in `getRunnerSnapshot` |
| Cycle detection | DFS + recStack (white/gray/black) | Called inside `FinishDag`; operates on `copyDag` snapshot |
| Status transitions | `TransitionStatus(from, to)` CAS | Validates state-machine edge before lock; `SetStatus` kept for unconditional override only |
| Error observability | `reportError` uses logrus fields | `dag_id` + `error` fields emitted on drop; `collectErrors` uses `DagConfig.ErrorDrainTimeout` |
