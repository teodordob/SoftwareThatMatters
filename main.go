package main

import (
	"github.com/AJMBrands/SoftwareThatMatters/cmd"
)

func main() {
	//TODO: Move to graph.go; Integrate nicely with cli
	// To use the cli: go run main.go start.
	cmd.Execute()
	// graph1, _, nodeInfoMap, _ := graph.CreateGraph("data/input/processed-100k.json", false)

	// pr := network.PageRankSparse(graph1, 0.85, 0.001)

	//duration := 365 * 24 * time.Hour
	//beginTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z01:00")
	//endTime := beginTime.Add(duration)
	////
	//var nodeMap map[int64]g.NodeInfo
	//var graph1 *simple.DirectedGraph
	//var stringIDToNodeInfo map[string]g.NodeInfo
	//graph1, _, stringIDToNodeInfo, nodeMap, _ = g.CreateGraph("data/input/processed-10k.json", false)
	//x := g.GetTransitiveDependenciesNode(graph1, nodeMap, stringIDToNodeInfo, "1221-1.0.0")
	//fmt.Println(x)
	////
	//g.FilterGraph(graph1, nodeMap, beginTime, endTime)
	//fmt.Println(graph1)
	//_ = network.PageRank(graph1, 0.85, 0.00001)

	//var nodeMap map[int64]g.NodeInfo
	//var stringMap map[string]g.NodeInfo
	//var graph1 *simple.DirectedGraph
	//graph1, _, stringMap, nodeMap, _ = g.CreateGraph("data/input/test_data.json", true)
	//g.FilterGraph(graph1, nodeMap, beginTime, endTime)
	//
	//g.FilterNode(graph1, nodeMap, stringMap, "A-1.0.0", beginTime, endTime)
	//_ = network.PageRank(graph1, 0.85, 0.00001)

	//g.FilterNode(graph1, nodeMap, stringMap, "A-1.0.0", beginTime, endTime)
	//g.GetTransitiveDependenciesNode(graph1, nodeMap, stringMap, "A-1.0.0")
	//_ = network.PageRank(graph1, 0.85, 0.00001)

	//Uncomment this to create the visualization and use these commands in the dot file
	//Toggle Preview - ctrl+shift+v (Mac: cmd+shift+v)
	//Open Preview to the Side - ctrl+k v (Mac: cmd+k shift+v)
	// graph.Visualization(graph1, "OnlyIds")
	// graph.VisualizationNodeInfo(nodeInfoMap, graph1, "IDInfo")
}
