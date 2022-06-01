package main

import (
	"fmt"

	g "github.com/AJMBrands/SoftwareThatMatters/graph"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

func main() {
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
	x := w.Walk(graph1, graph1.Node(0), nil)
	pageranking := network.PageRank(graph1, 0.85, 0.00001)

	fmt.Println(x)
	fmt.Println(pageranking)
	//Uncomment this to create the visualization and use these commands in the dot file
	//Toggle Preview - ctrl+shift+v (Mac: cmd+shift+v)
	//Open Preview to the Side - ctrl+k v (Mac: cmd+k shift+v)
	// g.Visualization(graph, "OnlyIds")
	// g.VisualizationNodeInfo(stringIDToNodeInfo, graph, "IDInfo")
}
