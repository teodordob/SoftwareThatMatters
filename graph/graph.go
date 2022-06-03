package graph

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/Masterminds/semver"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

type VersionInfo struct {
	Timestamp    string            `json:"timestamp"`
	Dependencies map[string]string `json:"dependencies"`
}

type PackageInfo struct {
	Name     string                 `json:"name"`
	Versions map[string]VersionInfo `json:"versions"`
}

// NodeInfo is a type structure for nodes. Name and Version can be removed if we find we don't use them often enough
type NodeInfo struct {
	id        int64
	stringID  string
	Name      string
	Version   string
	Timestamp string
}

// NewNodeInfo constructs a NodeInfo structure and automatically fills the stringID.
func NewNodeInfo(id int64, name string, version string, timestamp string) *NodeInfo {
	return &NodeInfo{
		id:        id,
		stringID:  fmt.Sprintf("%s-%s", name, version),
		Name:      name,
		Version:   version,
		Timestamp: timestamp}
}

func (nodeInfo NodeInfo) String() string {
	return fmt.Sprintf("Package: %v - Version: %v", nodeInfo.Name, nodeInfo.Version)
}

// CreateStringIDToNodeInfoMap takes a list of PackageInfo and a simple.DirectedGraph. For each of the packages,
// it creates a mapping of stringIDs to NodeInfo and also adds a node to the graph. The handling of the IDs is delegated
// to Gonum. These IDs are also included in the mapping for ease of access.
func CreateStringIDToNodeInfoMap(packagesInfo *[]PackageInfo, graph *simple.DirectedGraph) map[string]NodeInfo {
	stringIDToNodeInfoMap := make(map[string]NodeInfo, len(*packagesInfo))
	for _, packageInfo := range *packagesInfo {
		for packageVersion, versionInfo := range packageInfo.Versions {
			packageNameVersionString := fmt.Sprintf("%s-%s", packageInfo.Name, packageVersion)
			// Delegate the work of creating a unique ID to Gonum
			newNode := graph.NewNode()
			newId := newNode.ID()
			stringIDToNodeInfoMap[packageNameVersionString] = *NewNodeInfo(newId, packageInfo.Name, packageVersion, versionInfo.Timestamp)
			// idToNodeInfo[newId] =
			graph.AddNode(newNode)
		}
	}
	return stringIDToNodeInfoMap
}

// TODO: Maybe change to something like CreateIdToNodeInfoMap so it's not confusing for other people.

func CreateNodeIdToPackageMap(m map[string]NodeInfo) map[int64]NodeInfo {
	s := make(map[int64]NodeInfo, len(m))
	for _, val := range m {
		s[val.id] = val
	}
	return s
}

func CreateNameToVersionMap(m *[]PackageInfo) map[string][]string {
	newMap := make(map[string][]string, len(*m))
	for _, value := range *m {
		name := value.Name
		for k := range value.Versions {
			newMap[name] = append(newMap[name], k)
		}
	}
	return newMap
}

//Function to write the simple graph to a dot file so it could be visualized with GraphViz. This includes only Ids
func Visualization(graph *simple.DirectedGraph, name string) {
	result, _ := dot.Marshal(graph, name, "", "  ")

	file, err := os.Create(name + ".dot")

	if err != nil {
		log.Fatal("Error!", err)
	}
	defer file.Close()

	fmt.Fprint(file, string(result))

}

//Writes to dot file manually from the NodeInfoMap to include the Node info in the graphViz
//TODO: Optimize in the future since this is kind of barbaric probably there is a faster way.
func VisualizationNodeInfo(iDToNodeInfo *map[string]NodeInfo, graph *simple.DirectedGraph, name string) {
	file, err := os.Create(name + ".dot")
	d1 := []byte("strict digraph" + " " + name + " " + "{\n")
	d2 := []byte("}")
	lab := string("[label = \" ")
	edgIt := graph.Edges()

	fmt.Fprint(file, string(d1))

	for key, element := range *iDToNodeInfo {
		//fmt.Println("Key:", key, "=>", "Element:", element.id)
		fmt.Fprintf(file, fmt.Sprint(element.id)+lab+string(key)+` \n `+string(element.Version)+` \n `+string(element.Timestamp)+"\""+"];\n")

	}

	for edgIt.Next() {
		fmt.Fprintf(file, fmt.Sprint(edgIt.Edge().From().ID())+" -> "+fmt.Sprint(edgIt.Edge().To().ID())+";\n")
	}

	fmt.Fprint(file, string(d2))

	if err != nil {
		panic(err)
	}

	defer file.Close()

}

// CreateEdges takes a graph, a list of packages and their dependencies, a map of stringIDs to NodeInfo and
// a map of names to versions and creates directed edges between the dependent library and its dependencies.
// TODO: add documentation on how we use semver for edges
// TODO: Discuss removing pointers from maps since they are reference types without the need of using * : https://stackoverflow.com/questions/40680981/are-maps-passed-by-value-or-by-reference-in-go
func CreateEdges(graph *simple.DirectedGraph, inputList *[]PackageInfo, stringIDToNodeInfo map[string]NodeInfo, nameToVersionMap map[string][]string, isMaven bool) {
	r, _ := regexp.Compile("((?P<open>[\\(\\[])(?P<bothVer>((?P<firstVer>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)(?P<comma1>,)(?P<secondVer1>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)?)|((?P<comma2>,)?(?P<secondVer2>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)?))(?P<close>[\\)\\]]))|(?P<simplevers>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)")
	for id, packageInfo := range *inputList {
		for _, dependencyInfo := range packageInfo.Versions {
			for dependencyName, dependencyVersion := range dependencyInfo.Dependencies {
				finaldep := dependencyVersion
				if isMaven {
					finaldep = parseMultipleMavenSemVers(dependencyVersion, r)
				}
				constraint, err := semver.NewConstraint(finaldep)
				//c, err := semver2.ParseRange(dependencyVersion)
				if err != nil {
					continue
					//fmt.Println("sunt aici")
					//fmt.Println(finaldep)
					////log.Fatal(finaldep)
					//log.Fatal(err)
				}
				for _, v := range nameToVersionMap[dependencyName] {
					//newVersion, _ := semver2.Parse(v)
					newVersion, err := semver.NewVersion(v)
					if err != nil {
						//fmt.Println(v)
						//panic(err)
						continue
					}
					if constraint.Check(newVersion) {
						dependencyNameVersionString := fmt.Sprintf("%s-%s", dependencyName, v)
						dependencyNode := graph.Node(stringIDToNodeInfo[dependencyNameVersionString].id)
						packageNode := graph.Node(int64(id))
						// Ensure that we do not create edges to self because some packages do that...
						if dependencyNode != packageNode {
							graph.SetEdge(simple.Edge{F: packageNode, T: dependencyNode})
						}

					}
				}
			}
		}
	}
}

func parseMultipleMavenSemVers(s string, reg *regexp.Regexp) string {
	var finalResult string
	chars := []rune(s)
	openIndex := 0
	closeIndex := 0
	for i := 0; i < len(chars); i++ {
		char := string(chars[i])
		if char == "(" || char == "[" {
			openIndex = i
		}
		if char == ")" || char == "]" {
			closeIndex = i
			if i != len(chars)-1 {
				finalResult += translateMavenSemver(s[openIndex:closeIndex+1], reg) + " || "
			} else {
				finalResult += translateMavenSemver(s[openIndex:closeIndex+1], reg)
			}
		}

	}
	if closeIndex == 0 && openIndex == 0 {
		return translateMavenSemver(s, reg)
	}

	return finalResult
}

func translateMavenSemver(s string, reg *regexp.Regexp) string {
	match := reg.FindStringSubmatch(s)
	result := make(map[string]string)
	var finalResult string
	for i, name := range reg.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
		//TODO: What is happening here?
		//fmt.Printf("by name: %s %s\n", result["singur"])
	}
	if len(result["close"]) > 0 {
		if len(result["secondVer2"]) > 0 {
			if len(result["comma1"]) > 0 || len(result["comma2"]) > 0 {
				switch result["close"] {
				case "]":
					finalResult = "<= " + result["secondVer2"]
				case ")":
					finalResult = "< " + result["secondVer2"]
				}
			} else {
				finalResult = "= " + result["secondVer2"]
			}
		} else {
			if len(result["firstVer"]) > 0 && len(result["secondVer1"]) > 0 {
				switch result["open"] {
				case "[":
					finalResult = ">= " + result["firstVer"] + ", "
				case "(":
					finalResult = "> " + result["firstVer"] + ", "
				}
				switch result["close"] {
				case "]":
					finalResult += "<= " + result["secondVer1"]
				case ")":
					finalResult += "< " + result["secondVer1"]
				}
			} else if len(result["firstVer"]) > 0 && len(result["secondVer1"]) == 0 {
				switch result["open"] {
				case "[":
					finalResult = ">= " + result["firstVer"]
				case "(":
					finalResult = "> " + result["firstVer"]
				}
			}
		}
	} else {
		finalResult = ">= " + result["simplevers"]
	}
	return finalResult

}

func ParseJSON(inPath string) *[]PackageInfo {
	// For NPM at least, about 2 million packages are expected, so we initialize so the array doesn't have to be re-allocated all the time
	const expectedAmount int = 2000000
	// An array for now since lists aren't type-safe, and they would overcomplicate things
	result := make([]PackageInfo, 0, expectedAmount)
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

func CreateGraph(inputPath string, isUsingMaven bool) (*simple.DirectedGraph, map[string]NodeInfo, map[int64]NodeInfo, map[string][]string) {
	packagesList := ParseJSON(inputPath)
	graph := simple.NewDirectedGraph()
	stringIDToNodeInfo := CreateStringIDToNodeInfoMap(packagesList, graph)
	idToNodeInfo := CreateNodeIdToPackageMap(stringIDToNodeInfo)
	nameToVersions := CreateNameToVersionMap(packagesList)
	CreateEdges(graph, packagesList, stringIDToNodeInfo, nameToVersions, isUsingMaven)
	// TODO: This might cause some issues but for now it saves it quite a lot of memory
	runtime.GC()
	return graph, stringIDToNodeInfo, idToNodeInfo, nameToVersions
}

// This function returns true when time t lies in the interval [begin, end], false otherwise
func InInterval(t, begin, end time.Time) bool {
	return t.Equal(begin) || t.Equal(end) || t.After(begin) && t.Before(end)
}

// This is a helper function used to initialize all required auxillary data structures for the graph traversal
func initializeTraversal(g *simple.DirectedGraph, nodeMap map[int64]NodeInfo, connected []*graph.Edge, withinInterval map[int64]bool, beginTime time.Time, endTime time.Time, w traverse.DepthFirst) {
	nodes := g.Nodes()
	for nodes.Next() { // Initialize withinInterval data structure
		n := nodes.Node()
		id := n.ID()
		publishTime, _ := time.Parse(time.RFC3339, nodeMap[id].Timestamp)
		if InInterval(publishTime, beginTime, endTime) {
			withinInterval[id] = true
		}
	}

	// TODO: Discuss if we should just leave packages free-floating if they haven't been visited even once
	w = traverse.DepthFirst{
		Traverse: func(e graph.Edge) bool { // The dependent / parent node
			var traverse bool
			fromId := e.From().ID()
			toId := e.To().ID()
			if withinInterval[toId] {
				fromTime, _ := time.Parse(time.RFC3339, nodeMap[fromId].Timestamp) // The dependent node's time stamp
				toTime, _ := time.Parse(time.RFC3339, nodeMap[toId].Timestamp)     // The dependency node's time stamp
				if traverse = fromTime.After(toTime); traverse {
					connected = append(connected, &e)
				} // If the dependency was released before the parent node, add this edge to the connected nodes
			}

			return traverse
		},
	}
}

func removeDisconnected(g *simple.DirectedGraph, connected []*graph.Edge) {
	edges := g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		for _, disconnectedEdge := range connected { // Found that it's connected, move on
			if edge == *disconnectedEdge {
				break
			} else {
				g.RemoveEdge(edge.From().ID(), edge.To().ID())
			}
		}
	}
}

// This function removes stale edges from the specified graph by doing a DFS with all packages as the root node in O(n^2)
func traverseAndRemoveEdges(g *simple.DirectedGraph, withinInterval map[int64]bool, w traverse.DepthFirst, connected []*graph.Edge) {
	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		if withinInterval[n.ID()] { // We'll only consider traversing this subtree if its root was within the specified time interval
			_ = w.Walk(g, n, nil) // Continue walking this subtree until we've visited everything we're allowed to according to Traverse
			w.Reset()             // Clean up for the next iteration
		}
	}

	removeDisconnected(g, connected)

}

func traverseOneNode(g *simple.DirectedGraph, nodeId int64, withinInterval map[int64]bool, w traverse.DepthFirst, connected []*graph.Edge) {
	_ = w.Walk(g, g.Node(nodeId), nil)
	removeDisconnected(g, connected)
}

func FilterGraph(g *simple.DirectedGraph, nodeMap map[int64]NodeInfo, beginTime, endTime time.Time) {
	// This stores whether the package existed in the specified time range
	withinInterval := make(map[int64]bool, len(nodeMap))
	// This keeps track of which edges we've connected
	connected := make([]*graph.Edge, 0, len(nodeMap)*2)
	var w traverse.DepthFirst
	initializeTraversal(g, nodeMap, connected, withinInterval, beginTime, endTime, w) // Initialize all auxillary data structures for the traversal

	traverseAndRemoveEdges(g, withinInterval, w, connected) // Traverse the graph and remove stale edges

}

func findNode(stringMap map[string]NodeInfo, stringId string) (int64, bool) {
	var nodeId int64
	var correctOk bool
	if info, ok := stringMap[stringId]; ok {
		nodeId = info.id
		correctOk = true
	} else {
		log.Printf("String id %s was not found \n", stringId)
		correctOk = false
	}
	return nodeId, correctOk
}

func FilterNode(g *simple.DirectedGraph, nodeMap map[int64]NodeInfo, stringMap map[string]NodeInfo, stringId string, beginTime, endTime time.Time) {

	var nodeId int64
	if id, ok := findNode(stringMap, stringId); ok {
		nodeId = id
	} else {
		return // This function is a no-op if we don't have a correct string id
	}

	// This stores whether the package existed in the specified time range
	withinInterval := make(map[int64]bool, len(nodeMap))
	// This keeps track of which edges we've connected
	connected := make([]*graph.Edge, 0, len(nodeMap)*2)
	var w traverse.DepthFirst
	initializeTraversal(g, nodeMap, connected, withinInterval, beginTime, endTime, w) // Initialize all auxillary data structures for the traversal

	traverseOneNode(g, nodeId, withinInterval, w, connected)
}

// This function returns the specified node and its dependencies
func GetTransitiveDependenciesNode(g *simple.DirectedGraph, nodeMap map[int64]NodeInfo, stringMap map[string]NodeInfo, stringId string) *[]NodeInfo {
	var nodeId int64
	result := make([]NodeInfo, 0, len(nodeMap)/2)
	if id, ok := findNode(stringMap, stringId); ok {
		nodeId = id
	} else {
		return &result // This function is a no-op if we don't have a correct string id
	}

	w := traverse.DepthFirst{
		Visit: func(n graph.Node) {
			result = append(result, nodeMap[n.ID()])
		},
	}

	_ = w.Walk(g, g.Node(nodeId), nil)
	return &result
}
