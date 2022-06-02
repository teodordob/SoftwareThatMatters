package main

import (
	"time"

	"github.com/AJMBrands/SoftwareThatMatters/cmd"
	g "github.com/AJMBrands/SoftwareThatMatters/graph"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	//TODO: Move to graph.go; Integrate nicely with cli
	// To use the cli: go run main.go start.
	cmd.Execute()
	duration := 365 * 24 * time.Hour
	beginTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z01:00")
	endTime := beginTime.Add(duration)

	var nodeMap map[int64]g.NodeInfo
	var stringMap map[string]g.NodeInfo
	var graph1 *simple.DirectedGraph
	graph1, _, stringMap, nodeMap, _ = g.CreateGraph("data/input/test_data.json", true)
	g.FilterGraph(graph1, nodeMap, beginTime, endTime)

	g.FilterNode(graph1, nodeMap, stringMap, "A-1.0.0", beginTime, endTime)
	g.GetTransitiveDependenciesNode(graph1, nodeMap, stringMap, "A-1.0.0")
	_ = network.PageRank(graph1, 0.85, 0.00001)

	//Uncomment this to create the visualization and use these commands in the dot file
	//Toggle Preview - ctrl+shift+v (Mac: cmd+shift+v)
	//Open Preview to the Side - ctrl+k v (Mac: cmd+k shift+v)
	// g.Visualization(graph, "OnlyIds")
	// g.VisualizationNodeInfo(stringIDToNodeInfo, graph, "IDInfo")
}
