# dag-go Project Progress Log

This file tracks the incremental improvement stages of the dag-go library.
Update this file at the start and end of every stage.

---

## Current Status: Stage 4 — Cycle Detection Completion & Progress Tracking

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

### Stage 4 — Cycle Detection Completion (current)
- **`errors.go`**: Added `ErrCycleDetected` sentinel error (`errors.New`).
- **`dag.go` (`FinishDag`)**: Cycle detection now returns `fmt.Errorf("FinishDag: %w", ErrCycleDetected)`,
  enabling callers to use `errors.Is(err, ErrCycleDetected)`.
- **`dag_test.go`**: Four cycle detection test cases added:
  - `TestDetectCycle_SimpleCycle`: start → A → B → A (2-node cycle via FinishDag).
  - `TestDetectCycle_ComplexCycle`: start → A → B → C → A (3-node cycle via FinishDag).
  - `TestDetectCycle_SelfLoop`: A → A (self-loop via direct graph construction + DetectCycle).
  - `TestDetectCycle_NoCycle`: diamond DAG — verifies FinishDag and DetectCycle both return clean.
- All tests passing; lint 0 issues.

---

## Pending Items

### Stage 5 — Error Handling & Observability Hardening (planned)
- [ ] `collectErrors` 5 s hard-cap: make configurable via `DagConfig.ErrorDrainTimeout`.
- [ ] `ErrNoRunner` structured error type (currently `errors.New`); add `NodeID` field via `NodeError`.
- [ ] `reportError` drop-log: replace `Log.Printf` with a structured log entry (logrus fields).
- [ ] Add `Dag.ErrCount()` helper to expose current Errors channel depth without draining.
- [ ] Consider `ErrorPolicy` enum: `FailFast` (current) vs `ContinueOnError`.

### Stage 6 — Performance & Concurrency (planned)
- [ ] `DagWorkerPool.taskQueue` is a bare `chan func()` — evaluate replacing with `SafeChannel`.
- [ ] `fanIn` goroutine: each node spawns a goroutine; explore bounded semaphore approach.
- [ ] `merge` result channel (`mergeResult chan bool`) — evaluate promotion to `SafeChannel[bool]`.
- [ ] Profile `TestComplexDag` under race detector; verify no data races.
- [ ] `Progress()` uses two independent `atomic.LoadInt64` calls — not atomic as a pair; document limitation.

### Stage 7 — Public API Hardening (planned)
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
