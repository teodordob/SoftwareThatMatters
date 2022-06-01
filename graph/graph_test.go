package graph

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

// TODO: Test ParseJSON
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
	stringMap := CreateStringIDToNodeInfoMap(&simplePackageInfo, graph)
	nameVersion := CreateNameToVersionMap(&simplePackageInfo)
	CreateEdges(graph, &simplePackageInfo, stringMap, nameVersion, false)

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

func TestNodeCreationMediumComplexity(t *testing.T) {
	packageB := PackageInfo{
		Name: "B",
		Versions: map[string]VersionInfo{
			"1.0.0": {
				Timestamp: "2021-04-22T20:15:37",
				Dependencies: map[string]string{
					"A": ">=0.9.0",
					"C": "1.0.0",
				},
			},
		},
	}
	packageC := PackageInfo{
		Name: "C",
		Versions: map[string]VersionInfo{
			"1.0.0": {
				Timestamp: "2022-04-22T20:13:34",
				Dependencies: map[string]string{
					"A": "<1.0.0",
				},
			},
			"2.0.0": {
				Timestamp: "2022-05-28T21:22:23",
				Dependencies: map[string]string{
					"A": "<2.0.0",
				},
			},
		},
	}
	packageA := PackageInfo{
		Name: "A",
		Versions: map[string]VersionInfo{
			"0.9.0": {
				Timestamp:    "2020-04-01T20:15:37",
				Dependencies: map[string]string{},
			},
			"1.0.0-rc.1": {
				Timestamp:    "2020-05-01T20:15:37",
				Dependencies: map[string]string{},
			},
			"1.0.0": {
				Timestamp:    "2021-06-01T20:15:37",
				Dependencies: map[string]string{},
			},
			"1.1.0": {
				Timestamp:    "2021-07-01T20:15:37",
				Dependencies: map[string]string{},
			},
			"2.0.0": {
				Timestamp:    "2022-01-04T04:02:00",
				Dependencies: map[string]string{},
			},
		},
	}
	mediumPackageInfo := []PackageInfo{
		packageB,
		packageC,
		packageA,
	}

	graph := simple.NewDirectedGraph()
	stringNodeInfo := CreateStringIDToNodeInfoMap(&mediumPackageInfo, graph)
	nameVersion := CreateNameToVersionMap(&mediumPackageInfo)
	CreateEdges(graph, &mediumPackageInfo, stringNodeInfo, nameVersion, false)

	t.Run("Creates 8 nodes, one for every package version", func(t *testing.T) {

		if numNodes := graph.Edges().Len(); numNodes != 8 {
			t.Errorf("Expected 8 edges, got %d", numNodes)
		}

	})

	t.Run("Creates the 8 correct nodes", func(t *testing.T) {
		packageIDS := []string{
			"A-0.9.0",
			"A-1.0.0-rc.1",
			"A-1.0.0",
			"A-1.1.0",
			"A-2.0.0",
			"B-1.0.0",
			"C-1.0.0",
			"C-2.0.0",
		}

		testInfo := map[string]NodeInfo{
			"A-0.9.0":      createTestNodeInfo(packageA, "0.9.0"),
			"A-1.0.0-rc.1": createTestNodeInfo(packageA, "1.0.0-rc.1"),
			"A-1.0.0":      createTestNodeInfo(packageA, "1.0.0"),
			"A-1.1.0":      createTestNodeInfo(packageA, "1.1.0"),
			"A-2.0.0":      createTestNodeInfo(packageA, "2.0.0"),
			"B-1.0.0":      createTestNodeInfo(packageB, "1.0.0"),
			"C-1.0.0":      createTestNodeInfo(packageC, "1.0.0"),
			"C-2.0.0":      createTestNodeInfo(packageC, "2.0.0"),
		}

		for _, v := range packageIDS {
			if actual, ok := stringNodeInfo[v]; !ok {
				t.Errorf("Package version node %s not found", v)
			} else {
				expected := testInfo[v]
				if !nodeInfosEqual(expected, actual) {
					t.Errorf("Node info for %s was incorrect (expected: %s, actual %s)", v, fmt.Sprint(expected), fmt.Sprint(actual))
				}

				if graph.Node(actual.id) == nil {
					t.Errorf("Node %s was not actually in the graph", v)
				}
			}
		}
	})
}

func nodeInfosEqual(expected, actual NodeInfo) bool {
	return expected.Name == actual.Name && expected.Version == actual.Version && expected.Timestamp == actual.Timestamp
}

func createTestNodeInfo(pi PackageInfo, version string) NodeInfo {
	return NodeInfo{
		id:        -1,
		stringID:  "invalid",
		Name:      pi.Name,
		Version:   version,
		Timestamp: pi.Versions[version].Timestamp,
	}
}

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
	stringIDToNodeInfo := CreateStringIDToNodeInfoMap(&simplePackagesInfo, graph)
	nameToVersions := CreateNameToVersionMap(&simplePackagesInfo)
	CreateEdges(graph, &simplePackagesInfo, stringIDToNodeInfo, nameToVersions, false)

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

func TestCreateEdgesMediumComplexityGraph(t *testing.T) {
	packagesInfo := []PackageInfo{
		{
			Name: "B",
			Versions: map[string]VersionInfo{
				"1.0.0": {
					Timestamp: "2021-04-22T20:15:37",
					Dependencies: map[string]string{
						"A": "> 0.9.0",
						"C": "1.0.0",
					},
				},
			},
		},
		{
			Name: "C",
			Versions: map[string]VersionInfo{
				"1.0.0": {
					Timestamp: "2022-04-22T20:13:34",
					Dependencies: map[string]string{
						"A": "< 1.0.0",
					},
				},
			},
		},
		{
			Name: "A",
			Versions: map[string]VersionInfo{
				"0.9.0": {
					Timestamp:    "2020-04-01T20:15:37",
					Dependencies: map[string]string{},
				},
				"1.0.0-rc.1": {
					Timestamp:    "2020-05-01T20:15:37",
					Dependencies: map[string]string{},
				},
				"1.0.0": {
					Timestamp:    "2021-06-01T20:15:37",
					Dependencies: map[string]string{},
				},
				"1.1.0": {
					Timestamp:    "2021-07-01T20:15:37",
					Dependencies: map[string]string{},
				},
			},
		},
	}

	graph := simple.NewDirectedGraph()
	stringIDToNodeInfo := CreateStringIDToNodeInfoMap(&packagesInfo, graph)
	nameToVersions := CreateNameToVersionMap(&packagesInfo)
	CreateEdges(graph, &packagesInfo, stringIDToNodeInfo, nameToVersions, false)
	t.Run("Creates 6 edges when there are 6 possible dependencies", func(t *testing.T) {
		if graph.Edges().Len() != 6 {
			t.Errorf("Expected 6 edges, got %d", graph.Edges().Len())
		}
	})
	t.Run("Creates edges to the correct dependencies for Node B-1.0.0", func(t *testing.T) {
		if graph.From(stringIDToNodeInfo["B-1.0.0"].id).Len() != 4 {
			t.Errorf("Expected 4 possible dependencies for Node B-1.0.0, got %d", graph.From(stringIDToNodeInfo["B-1.0.0"].id).Len())
		}
		nodesIterator := graph.From(stringIDToNodeInfo["B-1.0.0"].id)
		counter := 0
		for nodesIterator.Next() {
			currentNode := nodesIterator.Node()
			if currentNode.ID() == graph.Node(stringIDToNodeInfo["C-1.0.0"].id).ID() {
				counter++
			}
		}
		if counter > 1 {
			t.Errorf("Expected exactly 1 edge towards dependency C-1.0.0, got %d", counter)
		}

	})
	t.Run("Creates no edges from Node A-1.0.0 (it has no dependencies)", func(t *testing.T) {
		if graph.From(stringIDToNodeInfo["A-1.0.0"].id).Len() != 0 {
			t.Errorf("Expected 0 dependencies for Node A-1.0.0, got %d", graph.From(stringIDToNodeInfo["A-1.0.0"].id).Len())
		}
	})
}
