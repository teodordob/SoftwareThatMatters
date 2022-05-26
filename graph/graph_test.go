package graph

import (
	"gonum.org/v1/gonum/graph/simple"
	"testing"
)

// TODO: Test ParseJSON

func TestCreateEdgesBasicGraph(t *testing.T) {
	simplePackagesInfo := []PackageInfo{
		{
			Name: "B",
			Versions: map[string]VersionInfo{
				"1.0.0": {
					Timestamp: "2021-04-22T20:15:37",
					Dependencies: map[string]string{
						"A": "1.0.0",
					},
				},
			},
		},
		{
			Name: "A",
			Versions: map[string]VersionInfo{
				"1.0.0": {
					Timestamp:    "2021-04-01T20:15:37",
					Dependencies: map[string]string{},
				},
			},
		},
	}

	graph := simple.NewDirectedGraph()
	stringIDToNodeInfo := *CreateStringIDToNodeInfoMap(&simplePackagesInfo, graph)
	nameToVersions := *CreateNameToVersionMap(&simplePackagesInfo)
	CreateEdges(graph, &simplePackagesInfo, &stringIDToNodeInfo, &nameToVersions)

	t.Run("Creates one edge when there is one dependency", func(t *testing.T) {

		if graph.Edges().Len() != 1 {
			t.Errorf("Expected 1 edge, got %d", graph.Edges().Len())
		}
	})
	t.Run("Creates the edge with the correct direction (dependent -> dependency)", func(t *testing.T) {
		fromID := stringIDToNodeInfo["B-1.0.0"].id
		toID := stringIDToNodeInfo["A-1.0.0"].id
		if graph.Edge(fromID, toID) == nil {
			if graph.Edge(toID, fromID) != nil {
				t.Error("Expected the correct direction but got a reversed edge. Please check if the edge " +
					"creation happens in the correct direction.")
			} else {
				t.Error("Expected the correct direction but something went wrong.")
			}
		}
	})
}
