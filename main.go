package main

import (
	"fmt"

	g "github.com/AJMBrands/SoftwareThatMatters/graph"
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
	// 		a -- b
	a := g.NewGraphNode(3)
	b := g.NewGraphNode(4)
	a.AddNeighbor(b)
	b.AddNeighbor(a)

	// 		subgraph clusterC {
	clusterC := g.NewGraphNode(5)
	// 			C -- D
	C := g.NewGraphNode(6)
	D := g.NewGraphNode(7)
	C.AddNeighbor(D)
	D.AddNeighbor(C)
	// 	}

	// 	subgraph clusterB {
	clusterB := g.NewGraphNode(8)
	// 		d -- f
	d := g.NewGraphNode(9)
	f := g.NewGraphNode(10)
	d.AddNeighbor(f)
	f.AddNeighbor(d)
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
	clusterB.AddNeighbor(G)
	// }

	fmt.Println(a)

}
