package dag_go

import (
	"strings"
	"testing"
)

// edgePairsFromFuzzInput deterministically turns an arbitrary fuzzed string
// into a bounded sequence of (from, to) node-ID pairs, by splitting on
// whitespace and pairing consecutive tokens. This lets a single fuzzed
// string exercise arbitrary combinations of valid/invalid/reserved/
// duplicate/self-loop edges against AddEdge, rather than only the
// hand-picked cases the table tests cover.
func edgePairsFromFuzzInput(raw string) [][2]string {
	tokens := strings.Fields(raw)
	const maxPairs = 64 // bound graph size so a single fuzz iteration stays fast
	var pairs [][2]string
	for i := 0; i+1 < len(tokens) && len(pairs) < maxPairs; i += 2 {
		pairs = append(pairs, [2]string{tokens[i], tokens[i+1]})
	}
	return pairs
}

// FuzzAddEdge exercises AddEdge's validation and node/edge-creation logic
// (reserved-name rejection, self-loop rejection, duplicate-edge handling,
// auto node creation) against arbitrary generated input. The library's
// hand-picked table tests cover specific known edge cases; this targets the
// class of bug those can't: a panic on some malformed/unexpected combination
// of node IDs no one thought to write a table-test case for. AddEdge
// returning an error for most fuzzed input is expected and fine - the only
// failure mode this catches is a panic.
func FuzzAddEdge(f *testing.F) {
	f.Add("a b")
	f.Add("a a")
	f.Add("start_node a")
	f.Add("a end_node")
	f.Add("a b b c c a")
	f.Add("")
	f.Add("\x00 \x01")
	f.Add("a b a b")

	f.Fuzz(func(t *testing.T, raw string) {
		dag := NewDag()
		for _, pair := range edgePairsFromFuzzInput(raw) {
			_ = dag.AddEdge(pair[0], pair[1])
		}
	})
}

// FuzzFinishDag exercises validateTopology/detectCycle over arbitrary
// generated topologies (built the same way FuzzAddEdge does), catching
// panics in the cycle-detection/validation path itself that the hand-picked
// cycle/invalid-graph table tests structurally can't reach.
func FuzzFinishDag(f *testing.F) {
	f.Add("a b b c")
	f.Add("a b b a") // cycle
	f.Add("a b a c")
	f.Add("a a")
	f.Add("")

	f.Fuzz(func(t *testing.T, raw string) {
		dag := NewDag()
		for _, pair := range edgePairsFromFuzzInput(raw) {
			_ = dag.AddEdge(pair[0], pair[1])
		}
		_ = dag.FinishDag()
	})
}
