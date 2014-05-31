
package graph

import (
    "bufio"
    "errors"
    "fmt"
    "log"
    "math"
    "math/rand"
    "os"
    "runtime"
    "strings"
    "strconv"
    "time"
)

func init() {
    switch os.Getenv("MCDBG") {
    case "":
        rand.Seed(time.Now().UnixNano())
    default:
        // do nothing
    }
}

type Node struct {
    vertex int
    edges []int
}

func (n *Node) setVertex(v int) {
    n.vertex = v
}

// Create a new node from a space-separated list of ints
func NewNode(s string) (n Node, err error) {
    words := strings.Fields(s)
    if len(words) < 2 {
        err = errors.New("NewNode: Not a real node")
    }
    for i, v := range words {
        value, err := strconv.Atoi(v)
        if err != nil {
            log.Fatal(err)
        }
        if i == 0 {
            n.setVertex(value)
        } else {
            n.edges = append(n.edges, value)
        }
    }
    return
}

type Graph map[int]Node

// Create a new container of nodes
func NewGraph(lines []string) Graph {
    g := make(Graph)
    for _, v := range lines {
        n, err := NewNode(v)
        if err == nil {
            g[n.vertex] = n
        }
    }
    return g
}

// Retrieve adjacency list strings from string literal
func FromString(s string) Graph {
    lines := strings.Split(s, "\n")
    return NewGraph(lines)
}

// Retrieve adjacency list strings from a file
func FromFile(f string) Graph {
    file, err := os.Open(f)
    if err != nil {
        log.Fatal(err)
    }
    lines := []string{}
    reader := bufio.NewReader(file)
    for err == nil {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        lines = append(lines, line)
    }
    return NewGraph(lines)
}

// Remove an edge from the grpah, collapsing the two nodes
// on the edge into a single super-node
func CollapseEdge(a, b Node, removed map[int]int) (sN Node, err error) {
    // The returned node takes the place of A.
    sN.setVertex(a.vertex)
    // Update any references to b.vertex in the Removed map
    for i, w := range removed {
        // if i has been removed, and mapped to v,
        // update the i reference to point to the a.vertex.
        if w == b.vertex {
            removed[i] = a.vertex
        }
    }
    // Remove any self-cycles
    for _, v := range append(a.edges, b.edges...) {
        switch {
        case v == a.vertex:
            continue
        case v == b.vertex:
            continue
        case removed[v] == a.vertex:
            continue
        case removed[v] == b.vertex:
            log.Print("wtf - should be handled in first loop")
            continue
        default:
            sN.edges = append(sN.edges, v)
        }
    }
    return
}

// Collapse the graph into two super-nodes. Return number of edges between
// the remaining two nodes.
func (g Graph) Cut(c chan int) {
    superNodes := make(Graph)
    goneNodes := make(map[int]int)
    // randFactor := rand.Int()
    for i := 0; i < len(g) - 2; i++ {
        // Select a random node in g
        // TODO: Try this using the 'fact' that golang range
        // for a map is random
        randN := rand.Intn(len(g) - 1) + 1
        sN, super := superNodes[randN]
        ref, gone := goneNodes[randN]
        var n Node
        var message string
        var ok bool
        switch {
        case super:
            if sN.vertex == randN {
                ok = true
            } else {
                message = fmt.Sprintf("%v != %v", sN.vertex, randN + 1)
            }
            n = sN
        case gone:
            n, ok = superNodes[ref]
            message += fmt.Sprintf("%v key is 'gone'. goneN: %v, gone: %v",
                                  ref, goneNodes[ref], goneNodes)
        default:
            n, ok = g[randN]
            message = fmt.Sprintf("Could not find %v in g: %v",
                                  randN, g)
        }
        if !ok {
            log.Fatal(message)
        }
        other := n
        // Find a non-n vertex. This probably isn't a super-effective
        // algorithm for searching...
        for other.vertex == n.vertex {
            length := len(n.edges)
            if length < 1 {
                log.Print(n, super, gone, superNodes, goneNodes)
            }
            randomOther := rand.Intn(length)
            o := n.edges[randomOther] // Provides a key for the graph
            sN, super := superNodes[o]
            ref, gone := goneNodes[o]
            ok = false
            switch {
            case super:
                if sN.vertex == o {
                    ok = true
                } else {
                    message = fmt.Sprintf("%v != %v", sN.vertex, o)
                }
                other = sN
            case gone:
                other, ok = superNodes[ref]
                message += fmt.Sprintf("%v key is 'gone'. goneN: %v, gone: %v",
                                      ref, goneNodes[ref], goneNodes)
            default:
                other, ok = g[o]
                message = fmt.Sprintf("Could not find %v in g: %v",
                                      o, g)
            }
            if !ok {
                log.Fatal(message)
            }
        }
        sNode, err := CollapseEdge(n, other, goneNodes)
        goneNodes[other.vertex] = n.vertex
        if err != nil {
            log.Print(err)
        }
        delete(superNodes, other.vertex)
        superNodes[n.vertex] = sNode
    }
    /*
    if len(superNodes) != 2 {
        message := fmt.Sprintf("Should have 2 supernodes. Counted %v",
                              len(superNodes))
        err = errors.New(message)
    }
    */
    var count int
    for _, lastNode := range superNodes {
        count = len(lastNode.edges)
        // break
    }
    c <-count
}

func MinCut(g Graph) (min int) {
    ncpu := runtime.NumCPU()
    runtime.GOMAXPROCS(ncpu)
    c := make(chan int)
    min = len(g)
    n := float64(len(g))
    exp := 2.0
    trials := int(math.Pow(n, exp) * math.Log(n)) + 1
    log.Printf("Trying %v times", trials)
    for i := 0; i < ncpu; i++ {
        go g.Cut(c)
    }
    for i := 0; i < trials; i++ {
        cut := <-c
        if min > cut {
            min = cut
        }
        go g.Cut(c)
    }
    /*
    for i := 0; i < trials; {
        select {
        case err = <- e:
            log.Print(err)
        case cut = <- c:
            i++
            if min > cut {
                min = cut
            }
            go g.Cut(c, e)
        default:
            // Do nothing
        }
    }
    */
    return
}
