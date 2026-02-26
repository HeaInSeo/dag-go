package dag_go

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
