package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// --- Actor commands ---

// CommandKind identifies the type of command sent to an Actor's mailbox.
type CommandKind uint8

// CmdRun, CmdCancel, CmdRunContainer are the supported command kinds.
const (
	CmdRun CommandKind = iota
	CmdCancel
	CmdRunContainer // (예시용) Router/Driver 데모
)

// Command is the message type sent to an Actor's mailbox.
type Command struct {
	Kind    CommandKind
	RunReq  RunReq
	SpawnID string
	Sink    EventSink

	RunID  string
	NodeID string
}

// ===== Router & Driver (스텁) =====

// Driver abstracts the container execution backend (Podman, K8s, etc.).
type Driver interface {
	RunContainer(ctx context.Context, c Command) error
}

type noopDriver struct{}

// RunContainer is a no-op implementation that simulates a short delay.
func (d *noopDriver) RunContainer(ctx context.Context, c Command) error {
	// 실제 실행 대신 짧은 지연으로 처리 흉내
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Millisecond):
		return nil
	}
}

// Router dispatches commands to the appropriate Driver based on CommandKind.
type Router struct{ drv Driver }

// Handle routes the command to the configured Driver.
func (r *Router) Handle(ctx context.Context, c Command) error {
	switch c.Kind {
	case CmdRunContainer:
		return r.drv.RunContainer(ctx, c)
	default:
		return errors.New("unknown command kind")
	}
}

// ===== Event sink (simulating gRPC server stream) =====

// EventSink is the consumer interface for run progress events.
type EventSink interface {
	Emit(ev RunRespEvent) error
}

// StdoutSink writes run events to the standard logger.
type StdoutSink struct{}

// Emit logs the event fields to the standard logger.
func (StdoutSink) Emit(ev RunRespEvent) error {
	log.Printf("[%s] %s run=%s node=%s spawn=%s msg=%s",
		ev.Level, ev.State, ev.RunID, ev.NodeID, ev.SpawnID, ev.Msg)
	return nil
}

// ===== Fake driver (simulates Podman/K8s driver) =====

// FakeDriver simulates a Podman/K8s container driver for testing.
type FakeDriver struct{}

// Run simulates a container lifecycle (pull → exec → upload) with progress callbacks.
func (d *FakeDriver) Run(ctx context.Context, req RunReq, onProgress func(string)) error {
	onProgress("resolving image digest @sha256:…")
	time.Sleep(100 * time.Millisecond)
	onProgress("creating pod + container")
	time.Sleep(100 * time.Millisecond)
	onProgress("pulling inputs to /in (RO), preparing /work, /out")
	time.Sleep(100 * time.Millisecond)
	onProgress("executing userscript …")
	// Simulate work chunks
	for i := 1; i <= 3; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
			onProgress(fmt.Sprintf("step %d/3 ok", i))
		}
	}
	onProgress("uploading results to S3 (CAS)")
	return nil
}

// ===== Actor =====

// Actor processes commands sequentially from its mailbox goroutine.
type Actor struct {
	spawnID   string
	mbox      chan Command
	closing   atomic.Bool
	curMu     sync.Mutex
	curCancel context.CancelFunc

	router     *Router
	idleTTL    time.Duration
	onFinalize func()
}

// NewActor creates a minimal Actor with a small mailbox buffer.
func NewActor(spawnID string) *Actor {
	return &Actor{
		spawnID: spawnID,
		mbox:    make(chan Command, 64), // simple buffer
	}
}

// NewActor1 creates an Actor with a router, idle TTL, and finalization hook.
func NewActor1(spawnID string, router *Router, idleTTL time.Duration, onFinalize func()) *Actor {
	return &Actor{
		spawnID:    spawnID,
		mbox:       make(chan Command, 1024),
		router:     router,
		idleTTL:    idleTTL,
		onFinalize: onFinalize,
	}
}

// Start launches the Actor's internal event loop in a goroutine.
func (a *Actor) Start() { go a.loop() }

func (a *Actor) enqueue(c Command) bool {
	if a.closing.Load() {
		return false
	}
	select {
	case a.mbox <- c:
		return true
	default:
		return false
	}
}

func (a *Actor) loop() {
	defer func() {
		if a.onFinalize != nil {
			a.onFinalize()
		}
	}()
	// (간단 구현) idle TTL 생략. 필요시 타이머 추가 가능.
	for c := range a.mbox {
		switch c.Kind {
		case CmdRun:
			a.handleRun(c)
		case CmdCancel:
			// 최소 취소: 현재 실행 컨텍스트 취소
			a.curMu.Lock()
			if a.curCancel != nil {
				a.curCancel()
			}
			a.curMu.Unlock()
		case CmdRunContainer:
			// Router/Driver 경로 데모
			if a.router != nil {
				_ = a.router.Handle(context.Background(), c)
			}
		}
	}
}

func (a *Actor) handleRun(c Command) {
	ctx, cancel := context.WithCancel(context.Background())
	// 현재 실행 취소 핸들 저장
	a.curMu.Lock()
	a.curCancel = cancel
	a.curMu.Unlock()
	defer func() {
		a.curMu.Lock()
		a.curCancel = nil
		a.curMu.Unlock()
		cancel()
	}()

	emit := func(level, state, msg string) {
		if c.Sink != nil {
			_ = c.Sink.Emit(RunRespEvent{
				When:    time.Now(),
				Level:   level,
				State:   state,
				Msg:     msg,
				RunID:   c.RunReq.RunID,
				NodeID:  c.RunReq.NodeID,
				SpawnID: c.SpawnID,
			})
		}
	}

	emit("INFO", "Created", "actor accepted run")

	drv := &FakeDriver{}
	if err := drv.Run(ctx, c.RunReq, func(e string) { emit("INFO", "Running", e) }); err != nil {
		emit("ERROR", "Failed", err.Error())
		return
	}
	emit("INFO", "Succeeded", "work done")
}
