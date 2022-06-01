package main

import (
	"github.com/AJMBrands/SoftwareThatMatters/cmd"
	g "github.com/AJMBrands/SoftwareThatMatters/graph"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

func main() {
	// To use the cli: go run main.go start.
	cmd.Execute()
	parsed := g.ParseJSON("data/input/test_data.json")
	graph1 := simple.NewDirectedGraph()
	stringIDToNodeInfo := g.CreateStringIDToNodeInfoMap(parsed, graph1)
	nameToVersions := g.CreateNameToVersionMap(parsed)
	//idToPackageMap := g.CreateNodeIdToPackageMap(stringIDToNodeInfo)
	g.CreateEdges(graph1, parsed, stringIDToNodeInfo, nameToVersions, true)
	//g.Visualization(graph, "graph2")
	//fmt.Println(stringIDToNodeInfo)
	w := traverse.DepthFirst{
		Traverse: func(e graph.Edge) bool {
			return false
		},
	}
	// x := w.Walk(graph1, graph1.Node(0), nil)
	// pageranking := network.PageRank(graph1, 0.85, 0.00001)
}
