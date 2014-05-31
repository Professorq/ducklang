
package graph

import (
    "testing"
)

var ring =  `1 6 2
             2 1 3
             3 2 4
             4 3 5
             5 4 6
             6 5 1`

func TestNewGraph(t *testing.T) {
    graph := FromString(ring)
    if graph[1].vertex != 1 {
        t.Log("First vertex should be 1; is %v", graph[0].vertex)
        t.Fail()
    }
    if graph[6].vertex != 6 {
        t.Log("Last vertex should be 6; is %v", graph[5].vertex)
        t.Fail()
    }
    if graph[6].edges[0] != 5 || graph[6].edges[1] != 1 {
        t.Log("Edges are incorect: %v", graph)
        t.Fail()
    }
}

func TestEdgeCollapse(t *testing.T) {
    g := FromString(ring)
    const a, b = 1, 2
    sN, err := CollapseEdge(g[a], g[b], make(map[int]int))
    // Collapse entails a few behaviors:
    // 1) One of the nodes cease to exist
    //      a) Implementing the graph using a Map where the key is the vertex
    //
    //      b) It doesn't matter which ceases to exist. The second node 
    //           randomly selected will be collapsed into the second node
    //      
    // 2) Edge list of new node extends to include all edges of old node
    //        including duplicates.
    // 3) References to old node are converted to point to new node
    // 4) Self-loops are pruned (i.e. the edges between the two collapsed nodes)
    ExpectEdges := []int{6, 3}
    if len(ExpectEdges) != len(sN.edges) {
        t.Log(sN)
        t.Log("Received an edge list with the wrong length")
        t.Fail()
    } else {
        for i, v := range sN.edges {
            expected := ExpectEdges[i]
            if expected != v {
                t.Logf("%v != %v", expected, v)
                t.Fail()
            }
        }
    }
    if err != nil {
        t.Error(err)
    }
}

func TestAnyCutInRing(t *testing.T) {
    // c := make(chan int)
    // e := make(chan error)
    graph := FromString(ring)
    c := make(chan int)
    go graph.Cut(c)
    cut := <-c
    if cut != 2 {
        t.Log("Ring should always be cut in 2. Instead: ", cut)
        t.Fail()
    }
}

func TestGraphTraverse(t *testing.T) {
    source := "vikram_graph.txt"
    graph := FromFile(source)
    for _, node := range graph {
        for _, edge := range node.edges {
            AinB := false
            for _, reciprical := range graph[edge].edges {
                if reciprical == node.vertex {
                    AinB = true
                }
            }
            if !AinB {
                t.Logf("edge %v not reciprocated for %v", edge, node.vertex)
                t.Fail()
            }
        }
    }
}

func TestMinCut(t *testing.T) {
    const lines, source, expected = 40, "vikram_graph.txt", 3
    graph := FromFile(source)
    if len(graph) != lines {
        t.Logf("%v lines in %v. %v in graph", lines, source, len(graph))
        t.Fail()
    }
    min := MinCut(graph)
    if min != expected {
        t.Logf("%v MinCut expected. %v received", expected, min)
        t.Fail()
    }
}

func TestEasyMin(t *testing.T) {
    graph := FromString(ring)
    min := MinCut(graph)
    if min != 2 {
        t.Logf("Ring Mincut found %v, should have been 2", min)
        t.Fail()
    }
}
