# Changelog

All notable changes to dag-go are documented in this file.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v1.1.0] — 2026-05-18

### Summary
DAG execution lifecycle stabilization and error handling hardening.
Two tracks of work: structured error types / error policy (Stage 14) and
a comprehensive lifecycle guardrail pass (Stage 15).

### Added

#### Error handling (Stage 14)
- `ErrorPolicy` type with `ErrorPolicyFailFast` (default) and `ErrorPolicyContinueOnError` constants.
- `WithErrorPolicy(p ErrorPolicy) DagOption` functional option.
- `Dag.ErrCount() int` — point-in-time snapshot of buffered error count.
- `ErrNoRunner` sentinel now wrapped in `*NodeError{NodeID, Phase}` so callers can use both
  `errors.Is(err, ErrNoRunner)` and `errors.As(err, &ne)`.

#### Lifecycle guardrails (Stage 15)
- `Dag.startTrigger` field: captures the start-node trigger channel in `GetReadyE` before any
  goroutine is launched; `StartE` uses only this field — never reads `parentVertex` while goroutines are live.
- Topology-frozen guards on `AddEdge`, `AddEdgeIfNodesExist`, and `FinishDag`: return an error
  when called after `GetReadyE` succeeds and before `Reset` is called.
- Runner/config mutation guards on `SetContainerCmd`, `SetRunnerResolver`, `SetNodeRunner`, and
  `SetNodeRunners`: no-op / return false when the DAG is frozen (`nodeResult != nil || running`),
  covering both the running window and the post-`Wait` / pre-`Reset` window.
- Godoc mutation policy documented on all four setter functions.

### Changed

#### Lifecycle / concurrency fixes (Stage 15)
- `GetReadyE`: structural invariant check (`startNode.parentVertex` length) moved before goroutine
  launch; `running.Store(true)` moved before goroutine loop (Reset blocked from first spawn).
- `ConnectRunnerE`: upgraded from `dag.mu.RLock()` to `dag.mu.Lock()`; returns error if called
  after `GetReadyE` (prevents runner-closure data race while goroutines are live).
- `Wait` defer order: `nodeWg.Wait → closeChannels → running=false`.
  `closeChannels` now completes before `running` is cleared, closing the Reset/closeChannels race window.
- `Reset` / `ResetE`: `dag.mu.Lock()` acquired before `running.Load()` check, eliminating the
  TOCTOU window between the check and the lock.
- `SetNodeRunner`: frozen check moved inside `dag.mu.RLock()` so `nodeResult` is read under the lock.
- `SetNodeRunners`: frozen check moved inside existing `dag.mu.RLock()`.
- `Wait` defer comment rewritten: "LIFO" framing removed; each step's sequencing rationale documented.

### Tests added (Stage 14)
- `TestErrNoRunner_Structured`
- `TestErrCount_Basic`
- `TestErrorPolicy_ContinueOnError`
- `TestErrorPolicy_FailFast_Default`

### Tests added (Stage 15)
- `TestGetReadyE_UnexpectedStartParentVertex` (replaces `TestStartE_UnexpectedParentVertex`)
- `TestStartE_TriggerSendFails` (uses `d.startTrigger.Close()`)
- `TestStart_BeforeGetReady` (renames `TestStart_ClosedChannel`)
- `TestAddEdge_AfterGetReady`
- `TestCopyDag_HasNilStartTrigger`
- `TestReset_ClearsStartTrigger`
- `TestConnectRunnerE_AfterGetReady_ReturnsError`
- `TestSetContainerCmd_FrozenAfterGetReady`
- `TestSetRunnerResolver_FrozenAfterGetReady`
- `TestSetNodeRunner_FrozenAfterGetReady`
- `TestSetNodeRunners_FrozenAfterGetReady`
- `TestSetters_UnfrozenAfterReset`

### Test results
```
go test ./...       ok  github.com/seoyhaein/dag-go
go test -race ./... ok  github.com/seoyhaein/dag-go   (0 data races)
```

---

## [v1.0.2] — prior release

See `docs/PROGRESS_LOG.md` Stage 11–13 for details.

## [v1.0.1] — prior release

See `docs/PROGRESS_LOG.md` Stage 7–10 for details.

## [v1.0.0] — 2026-02-27

See `docs/RELEASE_v1.0.0.md` for full release notes.
Initial stable release: SafeChannel, TransitionStatus CAS, ErrCycleDetected,
fanIn errgroup, DagWorkerPool, Reset, ToMermaid, examples.
