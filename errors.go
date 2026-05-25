package dag_go

import "errors"

// ErrCycleDetected is returned by FinishDag when the graph contains a directed cycle.
// Callers can test for this condition with errors.Is(err, ErrCycleDetected).
var ErrCycleDetected = errors.New("cycle detected in DAG")

// ErrDependencyBlocked is returned by parentReceiverFunc when a parent channel
// delivers a Failed signal under FailFast policy.  It is distinct from
// infrastructure errors (nil channel, context cancellation) so that connectRunner
// can transition a blocked node to NodeStatusSkipped rather than NodeStatusFailed:
// the node never ran, so it is not "failed" — it was simply blocked by an upstream
// dependency.  Use errors.Is to test for this condition.
var ErrDependencyBlocked = errors.New("dependency blocked")

// ErrorType identifies the DAG operation that produced a systemError.
type (
	ErrorType int

	systemError struct {
		errorType ErrorType
		reason    error
	}
)

// AddEdge, StartDag, AddEdgeIfNodesExist, addEndNode, FinishDag are the
// ErrorType values that identify which DAG operation recorded an error.
const (
	AddEdge ErrorType = iota
	StartDag
	AddEdgeIfNodesExist
	addEndNode
	FinishDag
)

// ErrorPolicy controls how downstream nodes react to upstream failures.
type ErrorPolicy int

const (
	// ErrorPolicyFailFast causes downstream nodes to be skipped when any parent
	// node fails. This is the default and preserves the invariant that no node
	// runs after a dependency failure.
	ErrorPolicyFailFast ErrorPolicy = iota

	// ErrorPolicyContinueOnError allows downstream nodes to execute even when a
	// parent has failed. Nodes proceed through all three flight phases regardless
	// of parent outcome; per-node errors are still collected via Errors channel.
	ErrorPolicyContinueOnError
)
