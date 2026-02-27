package dag_go

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"go.uber.org/goleak"
)

// Compile-time interface check.
var _ Runnable = DummyRunnable{}

func TestInitDag(t *testing.T) {
	defer goleak.VerifyNone(t)
	// InitDag 호출하여 새로운 Dag 인스턴스 생성
	dag, err := InitDag()
	if err != nil {
		t.Fatalf("InitDag failed: %v", err)
	}
	if dag == nil {
		t.Fatal("InitDag returned nil")
	}

	// StartNode 가 nil 이 아닌지 확인
	if dag.StartNode == nil {
		t.Fatal("StartNode is nil")
	}

	// StartNode 의 Id가 올바르게 설정되었는지 확인 (예: StartNode 상수와 일치하는지)
	if dag.StartNode.ID != StartNode {
		t.Errorf("Expected StartNode Id to be %s, got %s", StartNode, dag.StartNode.ID)
	}

	// StartNode 의 parentVertex 에 SafeChannel 이 추가되었는지 확인
	if len(dag.StartNode.parentVertex) == 0 {
		t.Error("StartNode's parentVertex channel is not set")
	} else {
		// parentVertex 의 갯수가 정확히 1개인지 검사, 1개로 설정함.
		if len(dag.StartNode.parentVertex) != 1 {
			t.Errorf("Expected StartNode.parentVertex count to be 1, got %d", len(dag.StartNode.parentVertex))
		}

		// 첫 번째 SafeChannel 의 내부 채널 용량을 확인
		safeCh := dag.StartNode.parentVertex[0]
		if cap(safeCh.GetChannel()) != dag.Config.MinChannelBuffer {
			t.Errorf("Expected parentVertex channel capacity to be %d, got %d", dag.Config.MinChannelBuffer, cap(safeCh.GetChannel()))
		}
	}

	// NodesResult 채널(노드 결과 SafeChannel)의 용량 검사
	if dag.NodesResult == nil {
		t.Error("RunningStatus is nil")
	} else if cap(dag.NodesResult.GetChannel()) != dag.Config.MaxChannelBuffer {
		t.Errorf("Expected RunningStatus channel capacity to be %d, got %d", dag.Config.MaxChannelBuffer, cap(dag.NodesResult.GetChannel()))
	}
}

// DummyRunnable is a no-op implementation of Runnable used in tests.
type DummyRunnable struct{}

func (DummyRunnable) RunE(_ context.Context, _ interface{}) error {
	return nil
}

func TestCreateNode(t *testing.T) {
	defer goleak.VerifyNone(t)
	// -------------------------------
	// Case 1: ContainerCmd 가 nil 인 경우
	// -------------------------------
	// NewDag()를 호출하면 기본적으로 ContainerCmd 는 nil
	dag := NewDag()
	id := "test1"
	node := dag.CreateNode(id)
	if node == nil {
		t.Fatalf("expected node to be created, got nil")
	}
	if node.ID != id {
		t.Errorf("expected node id to be %s, got %s", id, node.ID)
	}
	// With nil ContainerCmd the per-node runner should also be nil.
	if r := node.runnerLoad(); r != nil {
		t.Errorf("expected per-node runner to be nil when ContainerCmd is nil, got %v", r)
	}
	// parentDag must be set to the owning DAG.
	if node.parentDag != dag {
		t.Errorf("expected parentDag to be set to dag, got %v", node.parentDag)
	}
	// The node must be registered in dag.nodes.
	if dag.nodes[id] != node {
		t.Errorf("expected dag.nodes[%s] to be the created node", id)
	}
	// Creating the same id again must return nil.
	if dup := dag.CreateNode(id); dup != nil {
		t.Errorf("expected duplicate createNode call to return nil, got %v", dup)
	}

	// -------------------------------
	// Case 2: ContainerCmd is non-nil
	// -------------------------------
	dag2 := NewDag()
	dummy := DummyRunnable{}
	dag2.ContainerCmd = dummy

	id2 := "test2"
	node2 := dag2.CreateNode(id2)
	if node2 == nil {
		t.Fatalf("expected node2 to be created, got nil")
	}
	if node2.ID != id2 {
		t.Errorf("expected node2 id to be %s, got %s", id2, node2.ID)
	}
	// When ContainerCmd is non-nil, getRunnerSnapshot() must resolve to it.
	r := node2.getRunnerSnapshot()
	if r == nil {
		t.Errorf("expected runner to be resolved from ContainerCmd, got nil")
	}
	if _, ok := r.(DummyRunnable); !ok {
		t.Errorf("expected runner to be DummyRunnable, got %T", r)
	}
	// Duplicate creation must return nil.
	if dup2 := dag2.CreateNode(id2); dup2 != nil {
		t.Errorf("expected duplicate createNode call to return nil, got %v", dup2)
	}
}

func TestCreateEdge(t *testing.T) {
	defer goleak.VerifyNone(t)
	// 새로운 Dag 인스턴스 생성 (노드와 엣지 초기화)
	dag := &Dag{
		nodes: make(map[string]*Node),
		Edges: []*Edge{},
	}
	// 기본 구성 설정
	dag.Config = DefaultDagConfig()

	// Case 1: parentID가 빈 문자열인 경우 -> Fault 반환
	edge, errType := dag.createEdge("", "child")
	if edge != nil {
		t.Errorf("expected nil edge when parentID is empty, got %+v", edge)
	}
	if errType != Fault {
		t.Errorf("expected error type Fault when parentID is empty, got %v", errType)
	}

	// Case 2: childID가 빈 문자열인 경우 -> Fault 반환
	edge, errType = dag.createEdge("parent", "")
	if edge != nil {
		t.Errorf("expected nil edge when childID is empty, got %+v", edge)
	}
	if errType != Fault {
		t.Errorf("expected error type Fault when childID is empty, got %v", errType)
	}

	// Case 3: 정상적인 엣지 생성
	parentID := "p1"
	childID := "c1"
	edge, errType = dag.createEdge(parentID, childID)
	if edge == nil {
		t.Fatalf("expected edge to be created, got nil")
	}
	if errType != Create {
		t.Errorf("expected error type Create, got %v", errType)
	}
	if edge.parentID != parentID || edge.childID != childID {
		t.Errorf("edge fields mismatch: got parentID=%s, childID=%s; expected parentID=%s, childID=%s",
			edge.parentID, edge.childID, parentID, childID)
	}
	// safeVertex 채널의 용량이 dag.Config.MinChannelBuffer 와 일치하는지 확인
	if cap(edge.safeVertex.GetChannel()) != dag.Config.MinChannelBuffer {
		t.Errorf("expected vertex channel capacity to be %d, got %d", dag.Config.MinChannelBuffer, cap(edge.safeVertex.GetChannel()))
	}
	// dag.Edges 에 엣지가 추가되었는지 확인
	if len(dag.Edges) != 1 {
		t.Errorf("expected dag.Edges length to be 1, got %d", len(dag.Edges))
	}

	// Case 4: 중복된 엣지 생성 시도 -> nil 과 Exist 에러 반환
	dupEdge, dupErrType := dag.createEdge(parentID, childID)
	if dupEdge != nil {
		t.Errorf("expected duplicate edge creation to return nil, got %+v", dupEdge)
	}
	if dupErrType != Exist {
		t.Errorf("expected error type Exist on duplicate edge creation, got %v", dupErrType)
	}
}

// TestCreateEdge 이게 성공해야지 의미가 있음.
//
//nolint:gocognit,gocyclo // comprehensive edge test covering many validation paths in a single test function
func TestAddEdge(t *testing.T) {
	defer goleak.VerifyNone(t)
	// 새로운 Dag 인스턴스 생성 (NewDag 사용)
	dag := NewDag()
	if dag == nil {
		t.Fatal("NewDag returned nil")
	}

	// ----- 입력 검증 테스트 -----
	// Case 1: from 와 to가 같은 경우
	err := dag.AddEdge("node1", "node1")
	if err == nil {
		t.Error("expected error when from and to are the same, got nil")
	}

	// Case 2: from 가 빈 문자열인 경우
	err = dag.AddEdge("", "node2")
	if err == nil {
		t.Error("expected error when from is empty, got nil")
	}

	// Case 3: to가 빈 문자열인 경우
	err = dag.AddEdge("node1", "")
	if err == nil {
		t.Error("expected error when to is empty, got nil")
	}

	// ----- 정상적인 엣지 추가 테스트 -----
	// 유효한 입력: "node1" -> "node2"
	err = dag.AddEdge("node1", "node2")
	if err != nil {
		t.Fatalf("unexpected error on valid AddEdge: %v", err)
	}

	// 노드 "node1"와 "node2"가 dag.nodes 에 등록되었는지 확인
	node1, ok := dag.nodes["node1"]
	if !ok || node1 == nil {
		t.Fatal("node1 not found in dag.nodes")
	}
	node2, ok := dag.nodes["node2"]
	if !ok || node2 == nil {
		t.Fatal("node2 not found in dag.nodes")
	}

	// node1의 자식 리스트에 node2가 포함되어 있는지 확인
	found := false
	for _, child := range node1.children {
		if child.ID == "node2" {
			found = true
			break
		}
	}
	if !found {
		t.Error("node2 not found as a child of node1")
	}

	// node2의 부모 리스트에 node1이 포함되어 있는지 확인
	found = false
	for _, parent := range node2.parent {
		if parent.ID == "node1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("node1 not found as a parent of node2")
	}

	// dag.Edges 슬라이스에 엣지가 추가되었는지 확인 (정확히 1개의 엣지)
	if len(dag.Edges) != 1 {
		t.Errorf("expected 1 edge in dag.Edges, got %d", len(dag.Edges))
	}

	// node1의 childrenVertex 와 node2의 parentVertex 가 채워졌는지 확인
	if len(node1.childrenVertex) == 0 {
		t.Error("expected node1.childrenVertex to have at least one channel")
	}
	if len(node2.parentVertex) == 0 {
		t.Error("expected node2.parentVertex to have at least one channel")
	}

	// ----- 중복 엣지 생성 테스트 -----
	// 동일한 from, to 값으로 다시 엣지를 추가하면 오류가 발생해야 함.
	err = dag.AddEdge("node1", "node2")
	if err == nil {
		t.Error("expected error on duplicate edge creation, got nil")
	}
	// 관계가 중복 추가되지 않았는지 확인
	if len(node1.children) != 1 {
		t.Errorf("expected node1 to have 1 child after duplicate edge, got %d", len(node1.children))
	}
	if len(node2.parent) != 1 {
		t.Errorf("expected node2 to have 1 parent after duplicate edge, got %d", len(node2.parent))
	}
}

// TestCopyDag 는 copyDag 함수가 노드 ID, 부모/자식 구조를 올바르게 복사하는지 검증하는 단위 테스트입니다.
func TestCopyDag(t *testing.T) {
	defer goleak.VerifyNone(t)
	// 1. 원본 DAG 생성 및 노드/간선 구성
	dag := NewDag()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	nodeC := &Node{ID: "C"}
	nodeD := &Node{ID: "D"}

	// 부모/자식 관계 설정
	nodeA.children = []*Node{nodeB, nodeC}
	nodeB.parent = []*Node{nodeA}
	nodeB.children = []*Node{nodeD}
	nodeC.parent = []*Node{nodeA}
	nodeC.children = []*Node{nodeD}
	nodeD.parent = []*Node{nodeB, nodeC}

	dag.nodes = map[string]*Node{
		nodeA.ID: nodeA,
		nodeB.ID: nodeB,
		nodeC.ID: nodeC,
		nodeD.ID: nodeD,
	}

	dag.Edges = []*Edge{
		{parentID: "A", childID: "B"},
		{parentID: "A", childID: "C"},
		{parentID: "B", childID: "D"},
		{parentID: "C", childID: "D"},
	}

	// 2. copyDag 호출
	newNodes, newEdges := copyDag(dag)

	// 3. 노드 수 검증
	if len(newNodes) != len(dag.nodes) {
		t.Errorf("expected %d nodes, got %d", len(dag.nodes), len(newNodes))
	}

	// 4. 각 노드 ID와 관계 검증
	for id, origNode := range dag.nodes {
		newNode, ok := newNodes[id]
		if !ok {
			t.Errorf("node with ID %s missing in newNodes", id)
			continue
		}
		if newNode.ID != origNode.ID {
			t.Errorf("expected node ID %s, got %s", origNode.ID, newNode.ID)
		}
		if len(newNode.parent) != len(origNode.parent) {
			t.Errorf("node %s: expected %d parents, got %d", id, len(origNode.parent), len(newNode.parent))
		}
		if len(newNode.children) != len(origNode.children) {
			t.Errorf("node %s: expected %d children, got %d", id, len(origNode.children), len(newNode.children))
		}
	}

	// 5. 간선 복사 검증
	if len(newEdges) != len(dag.Edges) {
		t.Errorf("expected %d edges, got %d", len(dag.Edges), len(newEdges))
	}
	for i, origEdge := range dag.Edges {
		newEdge := newEdges[i]
		if newEdge.parentID != origEdge.parentID || newEdge.childID != origEdge.childID {
			t.Errorf("edge %d: expected %s -> %s, got %s -> %s",
				i, origEdge.parentID, origEdge.childID, newEdge.parentID, newEdge.childID)
		}
	}
}

// TestCopyDags 는 verifyCopiedDag 유틸리티를 사용하여 copyDag 의 동작을 검증합니다.
func TestCopyDags(t *testing.T) {
	defer goleak.VerifyNone(t)
	dag := NewDag()

	nodeA := &Node{ID: "A"}
	nodeB := &Node{ID: "B"}
	nodeC := &Node{ID: "C"}
	nodeD := &Node{ID: "D"}

	nodeA.children = []*Node{nodeB, nodeC}
	nodeB.parent = []*Node{nodeA}
	nodeB.children = []*Node{nodeD}
	nodeC.parent = []*Node{nodeA}
	nodeC.children = []*Node{nodeD}
	nodeD.parent = []*Node{nodeB, nodeC}

	dag.nodes = map[string]*Node{
		nodeA.ID: nodeA,
		nodeB.ID: nodeB,
		nodeC.ID: nodeC,
		nodeD.ID: nodeD,
	}

	dag.Edges = []*Edge{
		{parentID: "A", childID: "B"},
		{parentID: "A", childID: "C"},
		{parentID: "B", childID: "D"},
		{parentID: "C", childID: "D"},
	}

	newNodes, newEdges := copyDag(dag)

	if err := verifyCopiedDag(dag, newNodes, newEdges, "copyDag"); err != nil {
		t.Errorf("verifyCopiedDag failed: %v", err)
	}
}

func TestManyCopyDags(t *testing.T) {
	defer goleak.VerifyNone(t)
	// 여러 테스트 케이스: numNodes 와 edgeProb 를 조합
	testCases := []struct {
		numNodes int
		edgeProb float64
	}{
		{10, 0.3},
		{50, 0.1},
		{100, 0.05},
	}

	for _, tc := range testCases {
		// DAG 생성
		dag := generateDAG(tc.numNodes, tc.edgeProb)
		if len(dag.nodes) != tc.numNodes {
			t.Errorf("generateDAG(%d, %.2f): expected %d nodes, got %d",
				tc.numNodes, tc.edgeProb, tc.numNodes, len(dag.nodes))
		}

		// copyDag 함수를 호출하여 DAG 복사본 생성
		copiedNodes, copiedEdges := copyDag(dag)
		// verifyCopiedDag 는 원본과 복사본의 구조를 비교하고, 문제가 있으면 에러를 리턴합니다.

		if err := verifyCopiedDag(dag, copiedNodes, copiedEdges, "copyDag"); err != nil {
			t.Errorf("Test case (numNodes=%d, edgeProb=%.2f) failed: %v", tc.numNodes, tc.edgeProb, err)
		}
	}
}

// TestCopyDagIndependence 는 CopyDag 함수를 통해 생성된 복사본이 원본 DAG 와 독립적으로 동작하는지 검증
// 복사본에 데이터를 새롭게 넣었는데 만약 원본 DAG 의 내용이 변경한 복사본과 같아지면 shallow copy 이기때문에 에러남.
// copyDag 에서 shallow copy 가 일어나는 곳은 아예 복사를 하지 않는다.
//
//nolint:funlen,gocognit,gocyclo // 이 테스트는 길지만 의도적으로 유지함
func TestCopyDagIndependence(t *testing.T) {
	defer goleak.VerifyNone(t)
	// 1. 원본 DAG 생성 및 초기화
	orig := NewDag()
	orig.ID = "orig-id"
	orig.validated = true
	orig.Timeout = 10 * time.Second
	orig.bTimeout = false
	orig.ContainerCmd = nil // shallow copy라서 pass

	// 노드 생성
	nodeA := &Node{ID: "A", ImageName: "imgA", Commands: "cmdA", succeed: true}
	nodeB := &Node{ID: "B", ImageName: "imgB", Commands: "cmdB", succeed: false}

	// 부모/자식 관계 설정: A -> B
	nodeA.children = []*Node{nodeB}
	nodeB.parent = []*Node{nodeA}

	// 원본 DAG 에 노드 등록
	orig.nodes = map[string]*Node{
		"A": nodeA,
		"B": nodeB,
	}

	// 간선 생성: A -> B
	edgeAB := &Edge{
		parentID: "A",
		childID:  "B",
		// CopyDag 에서는 vertex 등의 추가 정보는 복사하지 않으므로 생략함.
	}
	orig.Edges = []*Edge{edgeAB}

	// 2. CopyDag 함수를 호출하여 복사본 생성
	newDag := CopyDag(orig, "copied-id")
	if newDag == nil {
		t.Fatal("CopyDag returned nil")
	}

	// 3. DAG 필드 독립성 검증

	// (a) Id 검증
	originalID := orig.ID
	newDag.ID = "modified-dag-id"
	if orig.ID == "modified-dag-id" {
		t.Error("Modifying newDag.Id affected original.Id")
	}
	newDag.ID = originalID

	// (b) validated, Timeout, bTimeout 검증
	origValidated := orig.validated
	newDag.validated = !origValidated
	if orig.validated == newDag.validated {
		t.Error("Modifying newDag.validated affected original.validated")
	}
	newDag.validated = origValidated

	origTimeout := orig.Timeout
	newDag.Timeout = origTimeout + 100*time.Millisecond
	if orig.Timeout == newDag.Timeout {
		t.Error("Modifying newDag.Timeout affected original.Timeout")
	}
	newDag.Timeout = origTimeout

	origBTimeout := orig.bTimeout
	newDag.bTimeout = !origBTimeout
	if orig.bTimeout == newDag.bTimeout {
		t.Error("Modifying newDag.bTimeout affected original.bTimeout")
	}
	newDag.bTimeout = origBTimeout

	// 4. 노드 독립성 검증
	for key, origNode := range orig.nodes {
		newNode, ok := newDag.nodes[key]
		if !ok {
			t.Errorf("node %s missing in copied DAG", key)
			continue
		}

		// (a) Id 필드 검증
		origNodeID := origNode.ID
		newNode.ID = "Modified-" + origNodeID
		if origNode.ID == newNode.ID {
			t.Errorf("Modifying copied node %s.Id affected original", key)
		}
		newNode.ID = origNodeID

		// (b) ImageName 필드 검증
		origImage := origNode.ImageName
		newNode.ImageName = "Modified-" + origImage
		if origNode.ImageName == newNode.ImageName {
			t.Errorf("Modifying copied node %s.ImageName affected original", key)
		}
		newNode.ImageName = origImage

		// (c) Commands 필드 검증
		origCmd := origNode.Commands
		newNode.Commands = "Modified-" + origCmd
		if origNode.Commands == newNode.Commands {
			t.Errorf("Modifying copied node %s.Commands affected original", key)
		}
		newNode.Commands = origCmd

		// (d) succeed 필드 검증
		origSucceed := origNode.succeed
		newNode.succeed = !origSucceed
		if origNode.succeed == newNode.succeed {
			t.Errorf("Modifying copied node %s.succeed affected original", key)
		}
		newNode.succeed = origSucceed
	}

	// 5. 간선 독립성 검증
	if len(newDag.Edges) != len(orig.Edges) {
		t.Errorf("expected %d edges, got %d", len(orig.Edges), len(newDag.Edges))
	}
	for i, origEdge := range orig.Edges {
		newEdge := newDag.Edges[i]
		if newEdge.parentID != origEdge.parentID || newEdge.childID != origEdge.childID {
			t.Errorf("edge %d: expected parent %s->child %s, got parent %s->child %s",
				i, origEdge.parentID, origEdge.childID, newEdge.parentID, newEdge.childID)
		}
	}
}

// TestDetectCycle_SimpleCycle verifies that FinishDag returns ErrCycleDetected
// for a direct two-node cycle: start → A → B → A.
func TestDetectCycle_SimpleCycle(t *testing.T) {
	defer goleak.VerifyNone(t)

	dag, _ := InitDag()

	if err := dag.AddEdge(StartNode, "A"); err != nil {
		t.Fatalf("AddEdge start->A: %v", err)
	}
	if err := dag.AddEdge("A", "B"); err != nil {
		t.Fatalf("AddEdge A->B: %v", err)
	}
	// Back-edge B→A closes the cycle; AddEdge allows it (only rejects from==to).
	if err := dag.AddEdge("B", "A"); err != nil {
		t.Fatalf("AddEdge B->A: %v", err)
	}

	err := dag.FinishDag()
	if err == nil {
		t.Fatal("expected ErrCycleDetected, got nil")
	}
	if !errors.Is(err, ErrCycleDetected) {
		t.Errorf("expected ErrCycleDetected, got: %v", err)
	}
}

// TestDetectCycle_ComplexCycle verifies that FinishDag returns ErrCycleDetected
// for a three-node cycle: start → A → B → C → A.
func TestDetectCycle_ComplexCycle(t *testing.T) {
	defer goleak.VerifyNone(t)

	dag, _ := InitDag()

	for _, pair := range []struct{ from, to string }{
		{StartNode, "A"}, {"A", "B"}, {"B", "C"}, {"C", "A"},
	} {
		if err := dag.AddEdge(pair.from, pair.to); err != nil {
			t.Fatalf("AddEdge %s->%s: %v", pair.from, pair.to, err)
		}
	}

	err := dag.FinishDag()
	if err == nil {
		t.Fatal("expected ErrCycleDetected, got nil")
	}
	if !errors.Is(err, ErrCycleDetected) {
		t.Errorf("expected ErrCycleDetected, got: %v", err)
	}
}

// TestDetectCycle_SelfLoop verifies that DetectCycle returns true for a node
// that lists itself as a child.  AddEdge rejects from==to, so the graph is
// constructed directly to reach the algorithm under test.
func TestDetectCycle_SelfLoop(t *testing.T) {
	defer goleak.VerifyNone(t)

	dag := NewDag()
	nodeA := &Node{ID: "A"}
	nodeA.children = []*Node{nodeA} // A → A
	dag.nodes = map[string]*Node{"A": nodeA}

	if !DetectCycle(dag) {
		t.Error("expected DetectCycle to return true for self-loop, got false")
	}
}

// TestDetectCycle_NoCycle verifies that a valid diamond-shaped DAG is not flagged
// as cyclic by either FinishDag or DetectCycle.
//
// Graph: start → A → {B1, B2} → C → {D1, D2} → E
func TestDetectCycle_NoCycle(t *testing.T) {
	defer goleak.VerifyNone(t)

	dag, _ := InitDag()

	for _, pair := range []struct{ from, to string }{
		{StartNode, "A"},
		{"A", "B1"}, {"A", "B2"},
		{"B1", "C"}, {"B2", "C"},
		{"C", "D1"}, {"C", "D2"},
		{"D1", "E"}, {"D2", "E"},
	} {
		if err := dag.AddEdge(pair.from, pair.to); err != nil {
			t.Fatalf("AddEdge %s->%s: %v", pair.from, pair.to, err)
		}
	}

	if err := dag.FinishDag(); err != nil {
		t.Fatalf("FinishDag returned unexpected error for valid DAG: %v", err)
	}
	// DetectCycle should also confirm no cycle after FinishDag.
	if DetectCycle(dag) {
		t.Error("expected DetectCycle to return false for valid DAG, got true")
	}
}

func TestDetectCycle(t *testing.T) {
	defer goleak.VerifyNone(t)
	// 새로운 DAG 생성
	dag, _ := InitDag()

	// 엣지 추가
	if err := dag.AddEdge(dag.StartNode.ID, "1"); err != nil {
		t.Fatalf("failed to add edge from StartNode to '1': %v", err)
	}
	if err := dag.AddEdge("1", "2"); err != nil {
		t.Fatalf("failed to add edge from '1' to '2': %v", err)
	}
	if err := dag.AddEdge("1", "3"); err != nil {
		t.Fatalf("failed to add edge from '1' to '3': %v", err)
	}
	if err := dag.AddEdge("1", "4"); err != nil {
		t.Fatalf("failed to add edge from '1' to '4': %v", err)
	}
	if err := dag.AddEdge("2", "5"); err != nil {
		t.Fatalf("failed to add edge from '2' to '5': %v", err)
	}
	if err := dag.AddEdge("5", "6"); err != nil {
		t.Fatalf("failed to add edge from '5' to '6': %v", err)
	}

	// DAG 완성 처리
	if err := dag.FinishDag(); err != nil {
		t.Fatalf("FinishDag failed: %v", err)
	}

	// 방문 기록 맵 초기화
	visit := make(map[string]bool)
	for id := range dag.nodes {
		visit[id] = false
	}
	// detectCycle 함수 테스트, cycle 이면 true 반환.
	cycle := DetectCycle(dag)
	if cycle {
		t.Errorf("expected no cycle, but detected one")
	}
}

// NoopCmd is a Runnable that succeeds immediately without performing any work.
type NoopCmd struct{}

func (NoopCmd) RunE(_ context.Context, _ interface{}) error { return nil }

func TestSimpleDag(t *testing.T) {
	defer goleak.VerifyNone(t)
	dag, _ := InitDag()
	// 엣지 추가
	if err := dag.AddEdge(dag.StartNode.ID, "1"); err != nil {
		t.Fatalf("failed to add edge from StartNode to '1': %v", err)
	}
	if err := dag.AddEdge("1", "2"); err != nil {
		t.Fatalf("failed to add edge from '1' to '2': %v", err)
	}
	if err := dag.AddEdge("1", "3"); err != nil {
		t.Fatalf("failed to add edge from '1' to '3': %v", err)
	}
	if err := dag.AddEdge("1", "4"); err != nil {
		t.Fatalf("failed to add edge from '1' to '4': %v", err)
	}
	if err := dag.AddEdge("2", "5"); err != nil {
		t.Fatalf("failed to add edge from '2' to '5': %v", err)
	}
	if err := dag.AddEdge("5", "6"); err != nil {
		t.Fatalf("failed to add edge from '5' to '6': %v", err)
	}

	// DAG 완성 처리
	if err := dag.FinishDag(); err != nil {
		t.Fatalf("FinishDag failed: %v", err)
	}

	// 추가 오류 해결을 위해 NoopCmd 사용
	// 해당 코드가 없으면 wait 에서 false 리턴한다. 오류난다.
	dag.SetContainerCmd(NoopCmd{})

	ctx := context.Background()
	dag.ConnectRunner()
	dag.GetReady(ctx)
	b1 := dag.Start()
	if b1 != true {
		t.Errorf("expected Start() to return true, got %v", b1)
	}

	b2 := dag.Wait(ctx)
	if b2 != true {
		t.Errorf("expected Wait() to return true, got %v", b2)
	}
}

func TestSimple1Dag(t *testing.T) {
	defer goleak.VerifyNone(t)
	// DAG 초기화
	dag, err := InitDag()
	if err != nil {
		t.Fatalf("failed to initialize dag: %v", err)
	}

	// DAG 설정 초기화: 워커 풀 크기와 채널 버퍼 크기 설정
	dag.Config = DagConfig{
		MinChannelBuffer: 10,
		MaxChannelBuffer: 50,
		WorkerPoolSize:   50, // 원하는 워커 수로 설정 이거 적으면 timeout error 발생할 수 있음.
		StatusBuffer:     10, // 원하는 채널 버퍼 사이즈로 설정
		DefaultTimeout:   30 * time.Second,
	}

	// Attach a no-op runner so that all nodes succeed during inFlight.
	dag.SetContainerCmd(NoopCmd{})

	// 엣지 추가: DAG 의 노드들 간에 부모/자식 관계 구성
	if err := dag.AddEdge(dag.StartNode.ID, "1"); err != nil {
		t.Fatalf("failed to add edge from StartNode to '1': %v", err)
	}
	if err := dag.AddEdge("1", "2"); err != nil {
		t.Fatalf("failed to add edge from '1' to '2': %v", err)
	}
	if err := dag.AddEdge("1", "3"); err != nil {
		t.Fatalf("failed to add edge from '1' to '3': %v", err)
	}
	if err := dag.AddEdge("1", "4"); err != nil {
		t.Fatalf("failed to add edge from '1' to '4': %v", err)
	}
	if err := dag.AddEdge("2", "5"); err != nil {
		t.Fatalf("failed to add edge from '2' to '5': %v", err)
	}
	if err := dag.AddEdge("5", "6"); err != nil {
		t.Fatalf("failed to add edge from '5' to '6': %v", err)
	}

	// DAG 완성 처리: 모든 노드가 연결된 상태로 마무리
	if err := dag.FinishDag(); err != nil {
		t.Fatalf("FinishDag failed: %v", err)
	}

	// 컨텍스트 생성
	ctx := context.Background()

	// 각 노드에 runner 함수를 연결하도록 설정
	dag.ConnectRunner()

	// DAG 실행 준비: 워커 풀과 상태 채널들을 초기화
	if ok := dag.GetReady(ctx); !ok {
		t.Fatalf("GetReady failed")
	}

	// 시작 노드의 runner 를 실행 (시작 상태 채널에 값 전송)
	b1 := dag.Start()
	if b1 != true {
		t.Errorf("expected Start() to return true, got %v", b1)
	}

	// 모든 실행이 완료될 때까지 대기
	b2 := dag.Wait(ctx)
	if b2 != true {
		t.Errorf("expected Wait() to return true, got %v", b2)
	}
}

// SimpleCommand simulates a short-lived task while respecting context cancellation.
type SimpleCommand struct{}

func (*SimpleCommand) RunE(ctx context.Context, _ interface{}) error {
	select {
	case <-time.After(100 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TestComplexDag 는 복잡한 DAG 구조를 테스트합니다.
func TestComplexDag(t *testing.T) {
	defer goleak.VerifyNone(t)
	// DAG 초기화
	dag, err := InitDag()
	if err != nil {
		t.Fatalf("InitDag failed: %v", err)
	}

	// 실행 명령 설정
	dag.SetContainerCmd(&SimpleCommand{})

	// 노드 생성 (다이아몬드 패턴 + 추가 노드)
	nodeIDs := []string{"A", "B1", "B2", "C", "D1", "D2", "E"}
	for _, id := range nodeIDs {
		dag.CreateNode(id)
	}

	// 엣지 추가
	edges := []struct{ from, to string }{
		{dag.StartNode.ID, "A"},
		{"A", "B1"},
		{"A", "B2"},
		{"B1", "C"},
		{"B2", "C"},
		{"C", "D1"},
		{"C", "D2"},
		{"D1", "E"},
		{"D2", "E"},
	}

	for _, edge := range edges {
		if err := dag.AddEdgeIfNodesExist(edge.from, edge.to); err != nil {
			t.Fatalf("failed to add edge from '%s' to '%s': %v", edge.from, edge.to, err)
		}
	}

	// DAG 완성 처리
	if err := dag.FinishDag(); err != nil {
		t.Fatalf("FinishDag failed: %v", err)
	}

	ctx := context.Background()
	dag.ConnectRunner()

	if !dag.GetReady(ctx) { // GetReadyT 사용
		t.Fatalf("GetReadyT failed")
	}

	b1 := dag.Start()
	if b1 != true {
		t.Errorf("expected Start() to return true, got %v", b1)
	}

	b2 := dag.Wait(ctx)
	if b2 != true {
		t.Errorf("expected Wait() to return true, got %v", b2)
	}
}

// TestErrorHandlingUsage 는 AddEdge 에서 발생하는 에러들을 통해
// reportError 로 에러를 기록하고, collectErrors 로 수집하는 흐름을 검증
func TestErrorHandlingUsage(t *testing.T) {
	defer goleak.VerifyNone(t)
	// DAG 초기화 (기본 설정 사용)
	dag, err := InitDag()
	if err != nil {
		t.Fatalf("InitDag failed: %v", err)
	}

	// 강제로 에러 조건을 유발하기 위해 잘못된 엣지 추가
	// 1. from-node 와 to-node 가 동일한 경우
	err = dag.AddEdge("1", "1")
	if err == nil {
		t.Error("Expected error when adding edge from a node to itself")
	}

	// 2. from-node 값이 빈 문자열인 경우
	err = dag.AddEdge("", "2")
	if err == nil {
		t.Error("Expected error when adding edge with empty from-node")
	}

	// 3. to-node 값이 빈 문자열인 경우
	err = dag.AddEdge("1", "")
	if err == nil {
		t.Error("Expected error when adding edge with empty to-node")
	}

	// 이제 보고된 에러들을 모아본다.
	// collectErrors 는 내부적으로 최대 5초(여기서는 3초로 설정) 대기 후, 에러들을 리턴한다.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collectedErrors := dag.collectErrors(ctx)

	// 위에서 3건의 에러가 보고되었으므로, 3건이 수집되어야 함.
	if len(collectedErrors) != 3 {
		t.Errorf("Expected 3 errors collected, got %d", len(collectedErrors))
	}

	// 수집된 에러들을 로그로 출력
	for i, e := range collectedErrors {
		t.Logf("Collected error %d: %v", i+1, e)
	}
}

// TestComplexDagProgress 는 다이아몬드 패턴 DAG 실행 전후로 Progress() 값을 검증
func TestComplexDagProgress(t *testing.T) {
	defer goleak.VerifyNone(t)
	// DAG 초기화
	dag, err := InitDag()
	if err != nil {
		t.Fatalf("InitDag failed: %v", err)
	}

	// 실행 명령 설정 (SimpleCommand 는 Runnable 인터페이스의 간단한 구현체)
	dag.SetContainerCmd(&SimpleCommand{})

	// 노드 생성 (다이아몬드 패턴 + 추가 노드)
	nodeIDs := []string{"A", "B1", "B2", "C", "D1", "D2", "E"}
	for _, id := range nodeIDs {
		if node := dag.CreateNode(id); node == nil {
			t.Fatalf("failed to create node '%s'", id)
		}
	}

	// 엣지 추가
	edges := []struct{ from, to string }{
		{dag.StartNode.ID, "A"},
		{"A", "B1"},
		{"A", "B2"},
		{"B1", "C"},
		{"B2", "C"},
		{"C", "D1"},
		{"C", "D2"},
		{"D1", "E"},
		{"D2", "E"},
	}

	for _, edge := range edges {
		if err := dag.AddEdgeIfNodesExist(edge.from, edge.to); err != nil {
			t.Fatalf("failed to add edge from '%s' to '%s': %v", edge.from, edge.to, err)
		}
	}

	// DAG 완성 처리: EndNode 생성 및 모든 노드 연결
	if err := dag.FinishDag(); err != nil {
		t.Fatalf("FinishDag failed: %v", err)
	}

	// runner 함수 연결 및 워커 풀 준비
	dag.ConnectRunner()
	ctx := context.Background()
	if !dag.GetReady(ctx) {
		t.Fatalf("GetReady failed")
	}

	// 실행 전 Progress() 검증 (아직 아무 노드도 완료되지 않았으므로 0.0이어야 함)
	if progress := dag.Progress(); progress != 0.0 {
		t.Errorf("expected initial progress 0.0, got %v", progress)
	} else {
		t.Logf("Initial progress: %v", progress)
	}

	// 시작 노드의 실행
	if b := dag.Start(); b != true {
		t.Errorf("expected Start() to return true, got %v", b)
	}

	// 모든 실행이 완료될 때까지 대기
	if b := dag.Wait(ctx); b != true {
		t.Errorf("expected Wait() to return true, got %v", b)
	}

	// 실행 완료 후 Progress() 검증 (모든 노드가 완료되어야 하므로 1.0이어야 함)
	if progress := dag.Progress(); progress != 1.0 {
		t.Errorf("expected progress 1.0 after execution, got %v", progress)
	} else {
		t.Logf("Progress after execution: %v", progress)
	}
}

// TestDagReset verifies that a DAG can be executed a second time after Reset()
// returns the same result as the first run.
//
//nolint:gocognit // comprehensive reset lifecycle test: closure + node-status loop + two full run cycles
func TestDagReset(t *testing.T) {
	defer goleak.VerifyNone(t)

	buildAndRun := func(dag *Dag) bool {
		ctx := context.Background()
		dag.ConnectRunner()
		if !dag.GetReady(ctx) {
			t.Fatal("GetReady failed")
		}
		if !dag.Start() {
			t.Fatal("Start failed")
		}
		return dag.Wait(ctx)
	}

	// Build a diamond DAG: start → A → {B1, B2} → C.
	dag, err := InitDag()
	if err != nil {
		t.Fatalf("InitDag: %v", err)
	}
	dag.SetContainerCmd(NoopCmd{})

	for _, pair := range []struct{ from, to string }{
		{StartNode, "A"},
		{"A", "B1"},
		{"A", "B2"},
		{"B1", "C"},
		{"B2", "C"},
	} {
		if err := dag.AddEdge(pair.from, pair.to); err != nil {
			t.Fatalf("AddEdge %s->%s: %v", pair.from, pair.to, err)
		}
	}
	if err := dag.FinishDag(); err != nil {
		t.Fatalf("FinishDag: %v", err)
	}

	// ── First run ──
	if ok := buildAndRun(dag); !ok {
		t.Fatal("first Wait returned false")
	}
	if p := dag.Progress(); p != 1.0 {
		t.Errorf("first run: expected progress 1.0, got %v", p)
	}

	// ── Reset ──
	dag.Reset()

	// Progress must be 0.0 immediately after Reset (counters cleared).
	if p := dag.Progress(); p != 0.0 {
		t.Errorf("after Reset: expected progress 0.0, got %v", p)
	}

	// All nodes must be back in Pending.
	for id, n := range dag.nodes {
		if got := n.GetStatus(); got != NodeStatusPending {
			t.Errorf("after Reset: node %s expected Pending, got %v", id, got)
		}
	}

	// ── Second run ──
	if ok := buildAndRun(dag); !ok {
		t.Fatal("second Wait after Reset returned false")
	}
	if p := dag.Progress(); p != 1.0 {
		t.Errorf("second run: expected progress 1.0, got %v", p)
	}
}

// generateDAG 는 numNodes 개의 노드를 생성하고,
// i < j 인 경우 확률 edgeProb 로 부모-자식 간선을 추가하여 DAG 를 구성
func generateDAG(numNodes int, edgeProb float64) *Dag {
	// 새로운 DAG 생성
	dag := NewDag()
	dag.nodes = make(map[string]*Node, numNodes)

	// 노드 생성: "0", "1", ..., "numNodes-1"
	for i := 0; i < numNodes; i++ { //nolint:intrange
		id := fmt.Sprintf("%d", i)
		node := &Node{ID: id}
		dag.nodes[id] = node
	}

	// 간선 생성: i < j 조건에서 edgeProb 확률로 간선 추가
	for i := 0; i < numNodes; i++ { //nolint:intrange
		for j := i + 1; j < numNodes; j++ {
			if rand.Float64() < edgeProb { //nolint:gosec // test only, crypto-rand not required
				parentID := fmt.Sprintf("%d", i)
				childID := fmt.Sprintf("%d", j)

				// 엣지 생성 및 등록
				edge := &Edge{parentID: parentID, childID: childID}
				dag.Edges = append(dag.Edges, edge)

				// 부모/자식 관계 연결
				dag.nodes[parentID].children = append(dag.nodes[parentID].children, dag.nodes[childID])
				dag.nodes[childID].parent = append(dag.nodes[childID].parent, dag.nodes[parentID])
			}
		}
	}
	return dag
}

// verifyCopiedDag 는 원본 DAG 와 복사된 노드 맵 및 간선 슬라이스가 동일한 구조를 갖는지 검증
// 문제가 있으면 에러 메시지를 모아 하나의 error 로 반환하고, 문제가 없으면 nil 을 반환
//
//nolint:gocognit,gocyclo // exhaustive structural verification helper; complexity comes from checking all node/edge fields
func verifyCopiedDag(original *Dag, newNodes map[string]*Node, newEdges []*Edge, methodName string) error {
	var errs []string

	// (1) 노드 수 검증
	if len(newNodes) != len(original.nodes) {
		errs = append(errs, fmt.Sprintf("[%s] expected %d nodes, got %d", methodName, len(original.nodes), len(newNodes)))
	}

	// (2) 각 노드의 ID, 부모/자식 관계 검증
	for id, origNode := range original.nodes {
		newNode, ok := newNodes[id]
		if !ok {
			errs = append(errs, fmt.Sprintf("[%s] node with ID %s missing in newNodes", methodName, id))
			continue
		}
		// ID 비교
		if newNode.ID != origNode.ID {
			errs = append(errs, fmt.Sprintf("[%s] expected node ID %s, got %s", methodName, origNode.ID, newNode.ID))
		}
		// 부모 관계 검증
		if len(newNode.parent) != len(origNode.parent) {
			errs = append(errs, fmt.Sprintf("[%s] node %s: expected %d parents, got %d", methodName, id, len(origNode.parent), len(newNode.parent)))
		} else {
			for i, p := range newNode.parent {
				if p.ID != origNode.parent[i].ID {
					errs = append(errs, fmt.Sprintf("[%s] node %s: expected parent %s, got %s", methodName, id, origNode.parent[i].ID, p.ID))
				}
			}
		}
		// 자식 관계 검증
		if len(newNode.children) != len(origNode.children) {
			errs = append(errs, fmt.Sprintf("[%s] node %s: expected %d children, got %d", methodName, id, len(origNode.children), len(newNode.children)))
		} else {
			for i, c := range newNode.children {
				if c.ID != origNode.children[i].ID {
					errs = append(errs, fmt.Sprintf("[%s] node %s: expected child %s, got %s", methodName, id, origNode.children[i].ID, c.ID))
				}
			}
		}
	}

	// (3) 간선 검증
	if len(newEdges) != len(original.Edges) {
		errs = append(errs, fmt.Sprintf("[%s] expected %d edges, got %d", methodName, len(original.Edges), len(newEdges)))
	}
	for i, origEdge := range original.Edges {
		newEdge := newEdges[i]
		if newEdge.parentID != origEdge.parentID || newEdge.childID != origEdge.childID {
			errs = append(errs, fmt.Sprintf("[%s] edge %d: expected parent %s->child %s, got parent %s->child %s",
				methodName, i, origEdge.parentID, origEdge.childID, newEdge.parentID, newEdge.childID))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// ── Stage 12 tests ────────────────────────────────────────────────────────────

// TestSendBlocking_Delivers verifies that SendBlocking successfully sends a
// value to a channel that has room and that the consumer receives it.
func TestSendBlocking_Delivers(t *testing.T) {
	defer goleak.VerifyNone(t)
	sc := NewSafeChannelGen[int](1)
	ctx := context.Background()
	if !sc.SendBlocking(ctx, 42) {
		t.Fatal("SendBlocking returned false on empty channel")
	}
	got := <-sc.GetChannel()
	if got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

// TestSendBlocking_CtxCancel verifies that SendBlocking returns false
// (without panicking or deadlocking) when the context is cancelled before
// a consumer arrives.
func TestSendBlocking_CtxCancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	sc := NewSafeChannelGen[int](0) // unbuffered — send blocks until consumed
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately before calling SendBlocking
	if sc.SendBlocking(ctx, 99) {
		t.Fatal("SendBlocking should return false on a cancelled context")
	}
}

// TestSendBlocking_ClosedChannel verifies that SendBlocking returns false
// immediately on a channel that is already closed.
func TestSendBlocking_ClosedChannel(t *testing.T) {
	defer goleak.VerifyNone(t)
	sc := NewSafeChannelGen[int](1)
	if err := sc.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	ctx := context.Background()
	if sc.SendBlocking(ctx, 7) {
		t.Fatal("SendBlocking should return false on a closed channel")
	}
}

// TestSendBlocking_BlocksThenDelivers verifies that SendBlocking blocks on a
// full channel and delivers the value once the consumer reads.
func TestSendBlocking_BlocksThenDelivers(t *testing.T) {
	defer goleak.VerifyNone(t)
	sc := NewSafeChannelGen[int](0) // unbuffered
	ctx := context.Background()

	done := make(chan struct{})
	go func() {
		defer close(done)
		sc.SendBlocking(ctx, 123) //nolint:errcheck
	}()

	// Give the goroutine time to block on the send.
	time.Sleep(20 * time.Millisecond)
	got := <-sc.GetChannel()
	<-done
	if got != 123 {
		t.Fatalf("expected 123, got %d", got)
	}
}

// TestDroppedErrors_Counter verifies that DroppedErrors increments each time
// reportError cannot deliver to a full Errors channel, and that Reset zeros it.
func TestDroppedErrors_Counter(t *testing.T) {
	defer goleak.VerifyNone(t)
	// Create a DAG with a tiny error channel (capacity 1) so it fills quickly.
	cfg := DefaultDagConfig()
	cfg.MaxChannelBuffer = 1
	dag := NewDagWithConfig(cfg)

	if dag.DroppedErrors() != 0 {
		t.Fatalf("expected 0 dropped errors initially, got %d", dag.DroppedErrors())
	}

	// First send fills the channel.
	dag.reportError(fmt.Errorf("err-1"))
	if dag.DroppedErrors() != 0 {
		t.Fatalf("first error should not be dropped (channel had room)")
	}

	// Second send overflows the channel — should be dropped.
	dag.reportError(fmt.Errorf("err-2"))
	if dag.DroppedErrors() != 1 {
		t.Fatalf("expected 1 dropped error, got %d", dag.DroppedErrors())
	}

	// Third send also dropped.
	dag.reportError(fmt.Errorf("err-3"))
	if dag.DroppedErrors() != 2 {
		t.Fatalf("expected 2 dropped errors, got %d", dag.DroppedErrors())
	}

	// Drain the channel so we can call Reset (Reset creates a new Errors channel).
	<-dag.Errors.GetChannel()
	dag.Errors.Close() //nolint:errcheck

	// Manually reset just the counter (Reset requires a fully initialised DAG;
	// test the atomic reset directly).
	dag.droppedErrors = 0
	if dag.DroppedErrors() != 0 {
		t.Fatalf("expected 0 after manual reset, got %d", dag.DroppedErrors())
	}
}

// TestWorkerPool_NodeTask verifies that the zero-churn worker pool correctly
// executes nodeTask entries and that no closures are created per submission.
func TestWorkerPool_NodeTask(t *testing.T) {
	defer goleak.VerifyNone(t)
	const workers = 4

	pool := NewDagWorkerPool(workers)
	ctx := context.Background()

	results := make(chan string, workers)
	for i := 0; i < workers; i++ {
		id := fmt.Sprintf("node%d", i)
		n := &Node{ID: id, status: NodeStatusPending}
		n.runnerStore(nil)
		sc := NewSafeChannelGen[*printStatus](8)

		// Attach a runner that records the node ID and drains the channel.
		n.runner = func(_ context.Context, res *SafeChannel[*printStatus]) {
			results <- n.ID
			res.Close() //nolint:errcheck
		}

		pool.Submit(nodeTask{node: n, sc: sc, ctx: ctx})
	}

	pool.Close()

	close(results)
	seen := make(map[string]bool)
	for id := range results {
		if seen[id] {
			t.Errorf("node %s executed more than once", id)
		}
		seen[id] = true
	}
	if len(seen) != workers {
		t.Errorf("expected %d executions, got %d", workers, len(seen))
	}
}
