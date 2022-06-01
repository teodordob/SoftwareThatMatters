package main

import (
	//"fmt"

	"time"

	g "github.com/AJMBrands/SoftwareThatMatters/graph"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

// This function returns true when time t lies in the interval [begin, end], false otherwise
func inInterval(t, begin, end time.Time) bool {
	return t.Equal(begin) || t.Equal(end) || t.After(begin) && t.Before(end)
}

func main() {
	duration := 365 * 24 * time.Hour
	beginTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z01:00")
	endTime := beginTime.Add(duration)

	var nodeMap map[int64]g.NodeInfo
	var graph1 *simple.DirectedGraph
	graph1, _, _, nodeMap, _ = g.CreateGraph("data/input/test_data.json", true)
	// This stores whether the package existed in the specified time range
	withinInterval := make(map[int64]bool, len(nodeMap))
	visited := make(map[int64]bool, len(nodeMap))
	nodes := graph1.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		id := n.ID()
		publishTime, _ := time.Parse(time.RFC3339, nodeMap[id].Timestamp)
		if inInterval(publishTime, beginTime, endTime) {
			withinInterval[id] = true
		}
	}
	//TODO: Use w.WalkAll and use 'Visit' function to mark visited nodes.
	w := traverse.DepthFirst{
		Visit: func(n graph.Node) {
			visited[n.ID()] = true // Mark this node as visited
		},
		Traverse: func(e graph.Edge) bool { // The dependent / parent node
			fromId := e.From().ID()
			if withinInterval[fromId] { // Only if the parent node is within the specified time frame will we even consider its child.
				toId := e.To().ID()
				if withinInterval[toId] {
					fromTime, _ := time.Parse(time.RFC3339, nodeMap[fromId].Timestamp) // The dependent node's time stamp
					toTime, _ := time.Parse(time.RFC3339, nodeMap[toId].Timestamp)     // The dependency node's time stamp
					return fromTime.After(toTime)                                      // If the dependency was released before the parent node, keep this edge connected
				}
			}
			return false
		},
	}
	_ = w.Walk(graph1, graph1.Node(0), nil)
	_ = network.PageRank(graph1, 0.85, 0.00001)
	//Uncomment this to create the visualization and use these commands in the dot file
	//Toggle Preview - ctrl+shift+v (Mac: cmd+shift+v)
	//Open Preview to the Side - ctrl+k v (Mac: cmd+k shift+v)
	// g.Visualization(graph, "OnlyIds")
	// g.VisualizationNodeInfo(stringIDToNodeInfo, graph, "IDInfo")
}
