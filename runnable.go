package dag_go

// Runnable is the interface that wraps the RunE method.
// Implementations are executed by a Node during the inFlight phase.
type Runnable interface {
	RunE(a interface{}) error
}
