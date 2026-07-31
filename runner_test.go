package dag_go

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.uber.org/goleak"
)

type (
	echoOK   struct{ d time.Duration }
	echoFail struct{ d time.Duration }
)

func (r echoFail) RunE(ctx context.Context, a interface{}) error {
	id := "<unknown>"
	if n, ok := a.(*Node); ok && n != nil {
		id = n.ID
	}
	fmt.Printf("[RunE FAIL] node=%s sleep=%s\n", id, r.d)
	select {
	case <-time.After(r.d):
		return errors.New("mock failure")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r echoOK) RunE(ctx context.Context, a interface{}) error {
	id := "<unknown>"
	if n, ok := a.(*Node); ok && n != nil {
		id = n.ID
	}
	fmt.Printf("[RunE OK] node=%s sleep=%s\n", id, r.d)
	select {
	case <-time.After(r.d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestSetNodeRunner_AfterGetReadyRejected(t *testing.T) {
	defer goleak.VerifyNone(t)

	dag, _ := InitDag()

	// Global default runner: succeed after 100 ms.
	dag.SetContainerCmd(echoOK{d: 100 * time.Millisecond})

	// Build node graph: start -> A -> B.
	a := dag.CreateNode("A")
	b := dag.CreateNode("B")
	_ = dag.AddEdge(StartNode, a.ID)
	_ = dag.AddEdge(a.ID, b.ID)

	// Resolver: node "B" uses a longer delay runner by default.
	dag.SetRunnerResolver(func(n *Node) Runnable {
		if n.ID == "B" {
			return echoOK{d: 300 * time.Millisecond}
		}
		return nil
	})

	if err := dag.FinishDag(); err != nil {
		t.Fatalf("FinishDag: %v", err)
	}
	if !dag.ConnectRunner() {
		t.Fatal("ConnectRunner failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if !dag.GetReady(ctx) {
		t.Fatal("GetReady failed")
	}

	// The DAG is frozen after GetReady. Even if B is still Pending, changing its
	// runner here would make execution depend on scheduler timing.
	if dag.SetNodeRunner("B", echoFail{d: 150 * time.Millisecond}) {
		t.Fatal("SetNodeRunner after GetReady should be rejected")
	}

	if !dag.Start() {
		t.Fatal("Start failed")
	}
	if !dag.Wait(ctx) {
		t.Fatal("Wait failed")
	}
}
