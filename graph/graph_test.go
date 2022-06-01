package graph

import (
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

func TestNodeCreationBasicGraph(t *testing.T) {
	simplePackageInfo := []PackageInfo{
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
	stringMap := *CreateStringIDToNodeInfoMap(&simplePackageInfo, graph)
	nameVersion := *CreateNameToVersionMap(&simplePackageInfo)
	CreateEdges(graph, &simplePackageInfo, &stringMap, &nameVersion, false)

	t.Run("Create two nodes because we specified two packages", func(t *testing.T) {

		numNodes := graph.Nodes().Len()

		if numNodes != 2 {
			t.Errorf("Expected two nodes, got %d", numNodes)
		}
	})

	t.Run("Create the two unique, correct nodes", func(t *testing.T) {
		var idA, idB int64
		if a, check := stringMap["A-1.0.0"]; check && graph.Node(idA) != nil {
			idA = a.id
		} else {
			t.Error("Node A-1.0.0 didn't exist")
		}

		if b, check := stringMap["B-1.0.0"]; check && graph.Node(idB) != nil {
			idB = b.id
		} else {
			t.Error("Node B-1.0.0 didn't exist")
		}

		if idA == idB {
			t.Errorf("Node IDs were equal (%d == %d)", idA, idB)
		}

	})
}
