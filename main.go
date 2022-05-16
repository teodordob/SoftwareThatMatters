package main

import (
	"fmt"

	g "github.com/AJMBrands/SoftwareThatMatters/graph"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/topo"
)

func main() {
	// This example shows the construction of the following graph
	// using the implicit graph type above.
	//
	// The visual representation of the graph can be seen at
	// https://graphviz.gitlab.io/_pages/Gallery/undirected/fdpclust.html
	//
	// graph G {
	// 	e
	// 	subgraph clusterA {
	// 		a -- b
	// 		subgraph clusterC {
	// 			C -- D
	// 		}
	// 	}
	// 	subgraph clusterB {
	// 		d -- f
	// 	}
	// 	d -- D
	// 	e -- clusterB
	// 	clusterC -- clusterB
	// }

	// graph G {
	G := g.NewGraphNode(0)
	// 	e
	e := g.NewGraphNode(1)

	// 	subgraph clusterA {
	clusterA := g.NewGraphNode(2)
	// 		a -- b
	a := g.NewGraphNode(3)
	b := g.NewGraphNode(4)
	a.AddNeighbor(b)
	b.AddNeighbor(a)
	clusterA.AddRoot(a)
	clusterA.AddRoot(b)

	// 		subgraph clusterC {
	clusterC := g.NewGraphNode(5)
	// 			C -- D
	C := g.NewGraphNode(6)
	D := g.NewGraphNode(7)
	C.AddNeighbor(D)
	D.AddNeighbor(C)

	clusterC.AddRoot(C)
	clusterC.AddRoot(D)
	// 		}
	clusterA.AddRoot(clusterC)
	// 	}

	// 	subgraph clusterB {
	clusterB := g.NewGraphNode(8)
	// 		d -- f
	d := g.NewGraphNode(9)
	f := g.NewGraphNode(10)
	d.AddNeighbor(f)
	f.AddNeighbor(d)
	clusterB.AddRoot(d)
	clusterB.AddRoot(f)
	// 	}

	// 	d -- D
	d.AddNeighbor(D)
	D.AddNeighbor(d)

	// 	e -- clusterB
	e.AddNeighbor(clusterB)
	clusterB.AddNeighbor(e)

	// 	clusterC -- clusterB
	clusterC.AddNeighbor(clusterB)
	clusterB.AddNeighbor(clusterC)

	G.AddRoot(e)
	G.AddRoot(clusterA)
	G.AddRoot(clusterB)
	// }

	if topo.IsPathIn(G, []graph.Node{C, D, d, f}) {
		fmt.Println("C--D--d--f is a path in G.")
	}

}
