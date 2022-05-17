package main

import (
	"fmt"
	g "github.com/AJMBrands/SoftwareThatMatters/graph"
)

func main() {
	// graph G {
	//G := g.NewGraphNode(0)
	//
	//// 	e
	//e := g.NewGraphNode(1)
	//
	//// 	subgraph clusterA {
	//// 		a -- b
	//a := g.NewGraphNode(3)
	//b := g.NewGraphNode(4)
	//a.AddNeighbor(b)
	//b.AddNeighbor(a)
	//
	//// 		subgraph clusterC {
	//clusterC := g.NewGraphNode(5)
	//// 			C -- D
	//C := g.NewGraphNode(6)
	//D := g.NewGraphNode(7)
	//C.AddNeighbor(D)
	//D.AddNeighbor(C)
	//// 	}
	//
	//// 	subgraph clusterB {
	//clusterB := g.NewGraphNode(8)
	//// 		d -- f
	//d := g.NewGraphNode(9)
	//f := g.NewGraphNode(10)
	//d.AddNeighbor(f)
	//f.AddNeighbor(d)
	//// 	}
	//
	//// 	d -- D
	//d.AddNeighbor(D)
	//D.AddNeighbor(d)
	//
	//// 	e -- clusterB
	//e.AddNeighbor(clusterB)
	//clusterB.AddNeighbor(e)
	//
	//// 	clusterC -- clusterB
	//clusterC.AddNeighbor(clusterB)
	//clusterB.AddNeighbor(clusterC)
	//clusterB.AddNeighbor(G)
	//// }

	x := g.DependenciesInfo{TimeStamp: "2020-10-11T15:19:50", Dependencies: map[string]string{"org.hamcrest:hamcrest-core": "1.3",
		"org.hamcrest:hamcrest-library": "1.3"}}
	y := g.DependenciesInfo{TimeStamp: "2007-05-14T21:34:53", Dependencies: map[string]string{}}
	z := g.DependenciesInfo{TimeStamp: "2007-05-14T21:34:18", Dependencies: map[string]string{"org.hamcrest:hamcrest-core": "1.3"}}
	as := g.DependenciesInfo{TimeStamp: "2021-10-11T15:19:50", Dependencies: map[string]string{"org.hamcrest:hamcrest-core": "1.3",
		"org.hamcrest:hamcrest-library": "1.3"}}
	var myarr = []g.DependenciesInfo{x, y, z}
	m := g.CreateMap(myarr)
	g.AddElementToMap(as, m)
	fmt.Println(m)
}
