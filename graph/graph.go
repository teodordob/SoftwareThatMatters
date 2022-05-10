package graph

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type node struct {
	name         string
	version      string
	id           string  // unique id formatted as "name-version"
	dependencies []*node // nodes upon which this node depends. Could be seen as parents.
	dependents   []*node // nodes that depend on this node. Could be seen as children.
}

// Node - Constructor for a node. Used to create ids.
func Node(name, version string) *node {
	return &node{
		name:         name,
		version:      version,
		id:           fmt.Sprintf("%s-%s", name, version),
		dependencies: []*node{},
		dependents:   []*node{},
	}
}

type graph struct {
	nodes []*node
}

// Graph - Constructor for a graph. Allows future customizations of the graph structure.
func Graph() *graph {
	return &graph{
		nodes: []*node{},
	}
}

func (graph *graph) IsNodeAlreadyPresent(n *node) bool {
	for _, node := range graph.nodes {
		if node.id == n.id {
			return true
		}
	}
	return false
}

func (graph *graph) Show() {
	for _, node := range graph.nodes {
		fmt.Printf("%s(V%s):\n", node.name, node.version)
		fmt.Printf("\tDependencies:\n")
		for _, dependencies := range node.dependencies {
			fmt.Printf("\t\t%s(V%s)\n", dependencies.name, dependencies.version)
		}
		fmt.Printf("\tDependents:\n")
		for _, dependents := range node.dependents {
			fmt.Printf("\t%s(V%s)\n", dependents.name, dependents.version)
		}
	}
}

func (graph *graph) CreateAndAddNode(name, version string) {
	node := Node(name, version)
	graph.AddNode(node)
}

func (graph *graph) AddNode(node *node) {

	if !graph.IsNodeAlreadyPresent(node) {
		graph.nodes = append(graph.nodes, node)
	} else {
		log.Printf("Node %s(V%s) already present in the graph. Skipping...\n", node.name, node.version)
	}
}

func (graph *graph) AddDependency(dependent, dependency *node) {
	if !graph.IsNodeAlreadyPresent(dependent) {
		graph.nodes = append(graph.nodes, dependent)
	}
	if !graph.IsNodeAlreadyPresent(dependency) {
		graph.nodes = append(graph.nodes, dependency)
	}
	dependent.dependencies = append(dependent.dependencies, dependency)
	dependency.dependents = append(dependency.dependents, dependent)
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Test() {
	//graph := &graph{}
	//for i := 0; i < 3; i++ {
	//	graph.AddNode(fmt.Sprintf("node%d", i), fmt.Sprintf("%d.0.0", i))
	//}
	//graph.AddNode("node1", "1.0.0")
	//graph.Show()

	graph := Graph()
	node1 := Node("A", "1.1")
	node2 := Node("B", "1")

	graph.AddNode(node1)
	graph.AddNode(node2)
	graph.AddDependency(node1, node2)
	graph.Show()
}

//func CreateNodesFromCSV(csvLine []string) (*node, *node) {
//	dependentNode := Node(csvLine[0], csvLine[1])
//	dependencyNode := Node(csvLine[3], csvLine[4])
//}

func CreateGraphFromDependenciesCSV() {
	csvFile, err := os.Open("data/input/dependencies.csv")
	check(err)
	defer csvFile.Close()

	csvLines, readerErr := csv.NewReader(csvFile).ReadAll()
	check(readerErr)

	graph := Graph()

	for _, line := range csvLines {
		dependentNode := Node(line[0], line[1])
		dependencyNode := Node(line[3], line[4])
		graph.AddDependency(dependentNode, dependencyNode)
	}
	graph.Show()
}
