package graph

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

// GraphNode is a node in an implicit graph.
type GraphNode struct {
	id        int64
	Neighbors []graph.Node
}

type DependenciesInfo struct {
	TimeStamp    string            `json:"timestamp"`
	Dependencies map[string]string `json:"dependencies"`
}

type PackageInfo struct {
	Name     string                      `json:"name"`
	Versions map[string]DependenciesInfo `json:"versions"`
}

// NewGraphNode returns a new GraphNode.
func NewGraphNode(id int64) *GraphNode {
	return &GraphNode{id: id}
}

// Node allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) Node(id int64) graph.Node {
	if id == g.id {
		return g
	}

	seen := map[int64]struct{}{g.id: {}}

	for _, n := range g.Neighbors {
		if n.ID() == id {
			return n
		}

		if gn, ok := n.(*GraphNode); ok {
			if gn.Has(seen, id) {
				return gn
			}
		}
	}

	return nil
}

func (g *GraphNode) Has(seen map[int64]struct{}, id int64) bool {

	for _, n := range g.Neighbors {
		if _, ok := seen[n.ID()]; ok {
			continue
		}

		seen[n.ID()] = struct{}{}
		if n.ID() == id {
			return true
		}

		if gn, ok := n.(*GraphNode); ok {
			if gn.Has(seen, id) {
				return true
			}
		}
	}

	return false
}

// Nodes allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) Nodes() graph.Nodes {
	nodes := []graph.Node{g}
	seen := map[int64]struct{}{g.id: {}}

	for _, n := range g.Neighbors {
		nodes = append(nodes, n)
		seen[n.ID()] = struct{}{}

		if gn, ok := n.(*GraphNode); ok {
			nodes = gn.nodes(nodes, seen)
		}
	}

	return iterator.NewOrderedNodes(nodes)
}

func (g *GraphNode) nodes(dst []graph.Node, seen map[int64]struct{}) []graph.Node {

	for _, n := range g.Neighbors {
		if _, ok := seen[n.ID()]; ok {
			continue
		}

		dst = append(dst, n)
		if gn, ok := n.(*GraphNode); ok {
			dst = gn.nodes(dst, seen)
		}
	}

	return dst
}

// From allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) From(id int64) graph.Nodes {
	if id == g.ID() {
		return iterator.NewOrderedNodes(g.Neighbors)
	}

	seen := map[int64]struct{}{g.id: {}}

	for _, n := range g.Neighbors {
		seen[n.ID()] = struct{}{}

		if gn, ok := n.(*GraphNode); ok {
			if result := gn.FindNeighbors(id, seen); result != nil {
				return iterator.NewOrderedNodes(result)
			}
		}
	}

	return nil
}

func (g *GraphNode) FindNeighbors(id int64, seen map[int64]struct{}) []graph.Node {
	if id == g.ID() {
		return g.Neighbors
	}

	for _, n := range g.Neighbors {
		if _, ok := seen[n.ID()]; ok {
			continue
		}
		seen[n.ID()] = struct{}{}

		if gn, ok := n.(*GraphNode); ok {
			if result := gn.FindNeighbors(id, seen); result != nil {
				return result
			}
		}
	}

	return nil
}

// HasEdgeBetween allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) HasEdgeBetween(uid, vid int64) bool {
	return g.EdgeBetween(uid, vid) != nil
}

// Edge allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) Edge(uid, vid int64) graph.Edge {
	return g.EdgeBetween(uid, vid)
}

// EdgeBetween allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) EdgeBetween(uid, vid int64) graph.Edge {
	if uid == g.id || vid == g.id {
		for _, n := range g.Neighbors {
			if n.ID() == uid || n.ID() == vid {
				return simple.Edge{F: g, T: n}
			}

		}
		return nil
	}

	seen := map[int64]struct{}{g.id: {}}

	for _, n := range g.Neighbors {
		seen[n.ID()] = struct{}{}
		if gn, ok := n.(*GraphNode); ok {
			if result := gn.edgeBetween(uid, vid, seen); result != nil {
				return result
			}
		}
	}

	return nil
}

func (g *GraphNode) edgeBetween(uid, vid int64, seen map[int64]struct{}) graph.Edge {
	if uid == g.id || vid == g.id {
		for _, n := range g.Neighbors {
			if n.ID() == uid || n.ID() == vid {
				return simple.Edge{F: g, T: n}
			}
		}
		return nil
	}

	for _, n := range g.Neighbors {
		if _, ok := seen[n.ID()]; ok {
			continue
		}

		seen[n.ID()] = struct{}{}
		if gn, ok := n.(*GraphNode); ok {
			if result := gn.edgeBetween(uid, vid, seen); result != nil {
				return result
			}
		}
	}

	return nil
}

// ID allows GraphNode to satisfy the graph.Node interface.
func (g *GraphNode) ID() int64 {
	return g.id
}

// AddMeighbor adds an edge between g and n.
func (g *GraphNode) AddNeighbor(n *GraphNode) {
	g.Neighbors = append(g.Neighbors, graph.Node(n))
}

func CreateMap(in *[]PackageInfo) *map[int64]PackageInfo {
	var id int64 = 0
	arr := *in
	m := make(map[int64]PackageInfo, len(arr))
	for n := range arr {
		m[id] = arr[n]
		id++
	}
	return &m
}
func AddElementToMap(x PackageInfo, inputMap *map[int64]PackageInfo) {
	m := *inputMap
	m[int64(len(m))] = x
}

func CreateNameToIDMap(m *map[int64]PackageInfo) *map[string]int64 {
	newMap := make(map[string]int64)
	for id, key := range *m {
		for versions, _ := range key.Versions {
			newKey := fmt.Sprintf("%s-%s", key.Name, versions)
			newMap[newKey] = id
		}
	}
	return &newMap
}

//func CreateHelperMap(m map[int64]PackageInfo) map[string]int64 {
//	n := make(map[string]int64)
//	for x, y := range m {
//
//	}
//	return n
//}

func CreateGraph(inputMap *map[int64]PackageInfo) *simple.DirectedGraph {
	m := *inputMap
	graph := simple.NewDirectedGraph()
	for x, _ := range m {
		graph.AddNode(NewGraphNode(x))
	}
	return graph
}

// CreateEdges takes a graph, a list of packages and their dependencies and a map of package names to package IDs
// and creates directed edges between the dependent library and its dependencies.
func CreateEdges(graph *simple.DirectedGraph, inputMap *map[int64]PackageInfo, nameToIDMap *map[string]int64) {
	packageInfo := *inputMap
	nameToID := *nameToIDMap
	for id, packageInfo := range packageInfo {
		for _, dependencyInfo := range packageInfo.Versions {
			for dependencyName, dependencyVersion := range dependencyInfo.Dependencies {
				dependencyNameVersionString := fmt.Sprintf("%s-%s", dependencyName, dependencyVersion)
				dependencyNode := graph.Node(nameToID[dependencyNameVersionString])
				packageNode := graph.Node(id)
				newEdge := graph.NewEdge(packageNode, dependencyNode)
				graph.SetEdge(newEdge)
			}
		}
	}
}

func ParseJSON(inPath string) *[]PackageInfo {
	var result []PackageInfo = make([]PackageInfo, 0, 10000)
	f, err := os.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(f)

	//Read opening bracket
	if _, err := dec.Token(); err != nil {
		log.Fatal(err)
	}

	for dec.More() {
		var packageInfo PackageInfo

		if err := dec.Decode(&packageInfo); err != nil {
			log.Fatal(err)
		}
		result = append(result, packageInfo)
	}

	//Read closing bracket
	if _, err := dec.Token(); err != nil {
		log.Fatal(err)
	}
	return &result
}
