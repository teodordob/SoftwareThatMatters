package main

import (
	"time"

	"github.com/AJMBrands/SoftwareThatMatters/cmd"
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
	//TODO: Move to graph.go; Integrate nicely with cli
	// To use the cli: go run main.go start.
	cmd.Execute()
	duration := 365 * 24 * time.Hour
	beginTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z01:00")
	endTime := beginTime.Add(duration)

	var nodeMap map[int64]g.NodeInfo
	var graph1 *simple.DirectedGraph
	graph1, _, _, nodeMap, _ = g.CreateGraph("data/input/test_data.json", true)
	// This stores whether the package existed in the specified time range
	withinInterval := make(map[int64]bool, len(nodeMap))
	// This keeps track of which edges we've visited
	traversed := make([][]bool, len(nodeMap))
	for from := range traversed {
		traversed[from] = make([]bool, 0, len(nodeMap))
	}
	nodes := graph1.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		id := n.ID()
		publishTime, _ := time.Parse(time.RFC3339, nodeMap[id].Timestamp)
		if inInterval(publishTime, beginTime, endTime) {
			withinInterval[id] = true
		}
	}
	// TODO: Discuss if we should just leave packages free-floating if they haven't been visited even once
	w := traverse.DepthFirst{
		Traverse: func(e graph.Edge) bool { // The dependent / parent node
			var traverse bool
			fromId := e.From().ID()
			toId := e.To().ID()
			if withinInterval[toId] {
				fromTime, _ := time.Parse(time.RFC3339, nodeMap[fromId].Timestamp) // The dependent node's time stamp
				toTime, _ := time.Parse(time.RFC3339, nodeMap[toId].Timestamp)     // The dependency node's time stamp
				if traverse = fromTime.After(toTime); traverse {
					traversed[fromId][toId] = true
				} // If the dependency was released before the parent node, keep this edge connected
			}

			return traverse
		},
	}

	traverseAndRemove(nodes, graph1, withinInterval, w, traversed)

	_ = network.PageRank(graph1, 0.85, 0.00001)

	//Uncomment this to create the visualization and use these commands in the dot file
	//Toggle Preview - ctrl+shift+v (Mac: cmd+shift+v)
	//Open Preview to the Side - ctrl+k v (Mac: cmd+k shift+v)
	// g.Visualization(graph, "OnlyIds")
	// g.VisualizationNodeInfo(stringIDToNodeInfo, graph, "IDInfo")
}

func traverseAndRemove(nodes graph.Nodes, graph1 *simple.DirectedGraph, withinInterval map[int64]bool, w traverse.DepthFirst, traversed [][]bool) {
	nodes = graph1.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		if withinInterval[n.ID()] { // We'll only consider traversing this subtree if its root was within the specified time interval
			_ = w.Walk(graph1, n, nil) // Continue walking this subtree until we've visited everything we're allowed to according to Traverse
			w.Reset()                  // Clean up for the next iteration
		}
	}

	for from := range traversed {
		for to, val := range traversed[from] {
			if !val { // If this potential edge wasn't touched at all, remove it
				graph1.RemoveEdge(int64(from), int64(to)) // Leaves unconnected nodes free-floating; ignores non-existent edges
			}
		}
	}
}
