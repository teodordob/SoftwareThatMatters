package graph

import (
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"log"
	"os"
	"time"

	"github.com/AJMBrands/SoftwareThatMatters/customgraph"
	"github.com/Masterminds/semver"
	"github.com/mailru/easyjson"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/traverse"
)

type VersionInfo struct {
	Dependencies map[string]string `json:"dependencies"`
	Timestamp    time.Time         `json:"timestamp"`
}

type PackageInfo struct {
	Versions map[string]VersionInfo `json:"versions"`
	Name     string                 `json:"name"`
}

type Doc struct {
	Pkgs []PackageInfo `json:"pkgs"`
}

// NodeInfo is a type structure for nodes. Name and Version can be removed if we find we don't use them often enough
type NodeInfo struct {
	Timestamp time.Time
	Name      string
	Version   string
	id        int64
}

var crcTable *crc64.Table = crc64.MakeTable(crc64.ISO)

// NewNodeInfo constructs a NodeInfo structure and automatically fills the stringID.
func NewNodeInfo(id int64, name string, version string, timestamp time.Time) *NodeInfo {
	return &NodeInfo{
		id: id,

		Name:      name,
		Version:   version,
		Timestamp: timestamp}
}

func (nodeInfo NodeInfo) String() string {
	return fmt.Sprintf("Package: %v - Version: %v", nodeInfo.Name, nodeInfo.Version)
}

func CreateNodeIdToPackageMap(m map[string]NodeInfo) map[int64]NodeInfo {
	s := make(map[int64]NodeInfo, len(m))
	for _, val := range m {
		s[val.id] = val
	}
	return s
}

func CreateHashedVersionMap(pi *[]PackageInfo) map[uint32][]string {
	result := make(map[uint32][]string, len(*pi))
	for _, pkg := range *pi {
		hashedName := hashPackageName(pkg.Name)
		result[hashedName] = make([]string, 0, len(pkg.Versions))
		for ver := range pkg.Versions {
			result[hashedName] = append(result[hashedName], ver)
		}
	}
	return result
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
func Visualization(graph *customgraph.DirectedGraph, name string) {
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
func VisualizationNodeInfo(iDToNodeInfo map[int64]NodeInfo, graph *customgraph.DirectedGraph, name string) {
	file, err := os.Create(name + ".dot")
	d1 := []byte("strict digraph" + " " + name + " " + "{\n")
	d2 := []byte("}")
	lab := string("[label = \" ")
	edgIt := graph.Edges()

	fmt.Fprint(file, string(d1))

	for _, element := range iDToNodeInfo {
		//fmt.Println("Key:", key, "=>", "Element:", element.id)
		fmt.Fprintf(file, fmt.Sprint(element.id)+lab+element.Name+` \n `+string(element.Version)+` \n `+element.Timestamp.String()+"\""+"];\n")

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
func CreateEdges(graph *customgraph.DirectedGraph, inputList *[]PackageInfo, hashToNodeId map[uint64]int64, nodeInfoMap map[int64]NodeInfo, hashToVersionMap map[uint32][]string, isMaven bool) int {
	// r, _ := regexp.Compile("((?P<open>[\\(\\[])(?P<bothVer>((?P<firstVer>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)(?P<comma1>,)(?P<secondVer1>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)?)|((?P<comma2>,)?(?P<secondVer2>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)?))(?P<close>[\\)\\]]))|(?P<simplevers>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)")
	n := len(*inputList)
	numEdges := 0
	for id, packageInfo := range *inputList {
		for version, dependencyInfo := range packageInfo.Versions {
			for dependencyName, dependencyVersion := range dependencyInfo.Dependencies {
				finaldep := dependencyVersion
				constraint, err := semver.NewConstraint(finaldep)
				if err != nil {
					//log.Println(err)
					continue
				}
				for _, v := range LookupVersions(dependencyName, hashToVersionMap) {
					//newVersion, _ := semver2.Parse(v)
					newVersion, err := semver.NewVersion(v)
					if err != nil {
						//fmt.Println(v)
						//panic(err)
						continue
					}
					if constraint.Check(newVersion) {
						dependencyStringId := fmt.Sprintf("%s-%s", dependencyName, v)
						dependencyGoId := LookupByStringId(dependencyStringId, hashToNodeId)

						packageStringId := fmt.Sprintf("%s-%s", packageInfo.Name, version)
						packageGoId := LookupByStringId(packageStringId, hashToNodeId)

						// Ensure that we do not create edges to self because some packages do that...
						if dependencyGoId != packageGoId {
							packageNode := graph.Node(packageGoId)
							dependencyNode := graph.Node(dependencyGoId)
							graph.SetEdge(customgraph.Edge{F: packageNode, T: dependencyNode})
							numEdges++
						}

					}
				}
			}
		}
		fmt.Printf("\u001b[1A \u001b[2K \r") // Clear the last line
		fmt.Printf("%.2f%% done (%d / %d packages connected to their dependencies)\n", float32(id)/float32(n)*100, id, n)
	}
	return numEdges
}

func ParseJSON(inPath string) []PackageInfo {

	f, err := os.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var result Doc
	err = easyjson.UnmarshalFromReader(f, &result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d packages\n", len(result.Pkgs))

	return result.Pkgs
}

func CreateMaps(packageList *[]PackageInfo, graph *customgraph.DirectedGraph) (map[uint64]int64, map[int64]NodeInfo, int) {
	hashToNodeId := make(map[uint64]int64, len(*packageList)*10)
	idToNodeInfo := make(map[int64]NodeInfo, len(*packageList)*10)
	numNodes := 0
	for _, packageInfo := range *packageList {
		for packageVersion, versionInfo := range packageInfo.Versions {
			stringID := fmt.Sprintf("%s-%s", packageInfo.Name, packageVersion)
			hashed := hashStringId(stringID)
			// Delegate the work of creating a unique ID to Gonum
			newNode := graph.NewNode()
			newId := newNode.ID()
			hashToNodeId[hashed] = newId
			idToNodeInfo[newId] = *NewNodeInfo(newId, packageInfo.Name, packageVersion, versionInfo.Timestamp)
			graph.AddNode(newNode)
			numNodes++
		}
	}
	return hashToNodeId, idToNodeInfo, numNodes
}

func hashStringId(stringID string) uint64 {
	hashed := crc64.Checksum([]byte(stringID), crcTable)
	return hashed
}

func hashPackageName(packageName string) uint32 {
	hashed := crc32.ChecksumIEEE([]byte(packageName))
	return hashed
}

func LookupVersions(packageName string, versionMap map[uint32][]string) []string {
	hash := hashPackageName(packageName)
	return versionMap[hash]
}

func LookupByStringId(stringId string, hashTable map[uint64]int64) int64 {
	hash := hashStringId(stringId)
	goId := hashTable[hash]
	return goId
}

func CreateGraph(inputPath string, isUsingMaven bool) (*customgraph.DirectedGraph, map[uint64]int64, map[int64]NodeInfo, map[uint32][]string) {
	fmt.Println("Parsing input")
	packagesList := ParseJSON(inputPath)
	// runtime.GC()
	graph := customgraph.NewDirectedGraph()
	// stringIDToNodeInfo := CreateStringIDToNodeInfoMap(packagesList, graph)
	// idToNodeInfo := CreateNodeIdToPackageMap(stringIDToNodeInfo)
	fmt.Println("Adding nodes and creating indices")
	hashToNodeId, idToNodeInfo, numNodes := CreateMaps(&packagesList, graph)
	// nameToVersions := CreateNameToVersionMap(&packagesList)
	hashToVersions := CreateHashedVersionMap(&packagesList)
	fmt.Println("Creating edges")
	fmt.Println()
	numEdges := CreateEdges(graph, &packagesList, hashToNodeId, idToNodeInfo, hashToVersions, isUsingMaven)
	//CreateEdgesConcurrent(graph, &packagesList, hashToNodeId, idToNodeInfo, nameToVersions, isUsingMaven)
	fmt.Println("Done!")
	// TODO: This might cause some issues but for now it saves it quite a lot of memory
	fmt.Printf("Nodes: %d, Edges: %d\n", numNodes, numEdges)
	return graph, hashToNodeId, idToNodeInfo, hashToVersions
}

// This function returns true when time t lies in the interval [begin, end], false otherwise
func InInterval(t, begin, end time.Time) bool {
	return t.Equal(begin) || t.Equal(end) || t.After(begin) && t.Before(end)
}

// initializeTraversal is a helper function used to initialize all required auxiliary data structures for the graph traversal
func initializeTraversal(g *customgraph.DirectedGraph, nodeMap map[int64]NodeInfo, withinInterval map[int64]bool, beginTime time.Time, endTime time.Time) {
	nodes := g.Nodes()
	for nodes.Next() { // Initialize withinInterval data structure
		n := nodes.Node()
		id := n.ID()
		publishTime := nodeMap[id].Timestamp
		if InInterval(publishTime, beginTime, endTime) {
			withinInterval[id] = true
		}
	}
}

func removeDisconnected(g *customgraph.DirectedGraph, connected []*graph.Edge) {
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
func traverseAndRemoveEdges(g *customgraph.DirectedGraph, nodeMap map[int64]NodeInfo, withinInterval map[int64]bool) {
	nodes := g.Nodes()
	// This keeps track of which edges we've connected
	connected := make([]*graph.Edge, 0, len(nodeMap)*2)

	t := traverse.BreadthFirst{
		Traverse: func(e graph.Edge) bool { // The dependent / parent node
			var traversal bool
			fromId := e.From().ID()
			toId := e.To().ID()
			if withinInterval[toId] {
				fromTime := nodeMap[fromId].Timestamp // The dependent node's time stamp
				toTime := nodeMap[toId].Timestamp     // The dependency node's time stamp
				if traversal = fromTime.After(toTime); traversal {
					connected = append(connected, &e)
				} // If the dependency was released before the parent node, add this edge to the connected nodes
			}

			return traversal
		},
	}

	nodesAmount := len(nodeMap)

	i := 0

	for nodes.Next() {
		n := nodes.Node()
		if withinInterval[n.ID()] { // We'll only consider traversing this subtree if its root was within the specified time interval
			_ = t.Walk(g, n, nil) // Continue walking this subtree until we've visited everything we're allowed to according to Traverse
			t.Reset()             // Clean up for the next iteration
			i++

			fmt.Printf("%d / %d subtrees walked \n", i, nodesAmount)
			fmt.Printf("\u001b[1A \u001b[2K \r") // Clear the last line
		}
	}

	removeDisconnected(g, connected)

}

func traverseOneNode(g *customgraph.DirectedGraph, nodeMap map[int64]NodeInfo, withinInterval map[int64]bool, nodeId int64) {
	connected := make([]*graph.Edge, 0, len(nodeMap)*2)

	t := traverse.BreadthFirst{
		Traverse: func(e graph.Edge) bool { // The dependent / parent node
			var traversal bool
			fromId := e.From().ID()
			toId := e.To().ID()
			if withinInterval[toId] {
				fromTime := nodeMap[fromId].Timestamp // The dependent node's time stamp
				toTime := nodeMap[toId].Timestamp     // The dependency node's time stamp

				if traversal = fromTime.After(toTime); traversal {
					connected = append(connected, &e)
				} // If the dependency was released before the parent node, add this edge to the connected nodes
			}

			return traversal
		},
	}

	_ = t.Walk(g, g.Node(nodeId), nil)
	removeDisconnected(g, connected)
}

func filterGraph(g *customgraph.DirectedGraph, nodeMap map[int64]NodeInfo, beginTime, endTime time.Time) {
	// This stores whether the package existed in the specified time range
	withinInterval := make(map[int64]bool, len(nodeMap))

	initializeTraversal(g, nodeMap, withinInterval, beginTime, endTime) // Initialize all auxillary data structures for the traversal

	traverseAndRemoveEdges(g, nodeMap, withinInterval) // Traverse the graph and remove stale edges
}

func FilterGraph(g *customgraph.DirectedGraph, nodeMap map[int64]NodeInfo, beginTime, endTime time.Time) {
	filterGraph(g, nodeMap, beginTime, endTime)
}

func findNode(hashMap map[uint64]int64, idToNodeInfo map[int64]NodeInfo, stringId string) (int64, bool) {
	var nodeId int64
	var correctOk bool
	if info, ok := idToNodeInfo[LookupByStringId(stringId, hashMap)]; ok {
		nodeId = info.id
		correctOk = true
	} else {
		log.Printf("String id %s was not found \n", stringId)
		correctOk = false
	}
	return nodeId, correctOk
}

func FilterNode(g *customgraph.DirectedGraph, hashMap map[uint64]int64, nodeMap map[int64]NodeInfo, stringId string, beginTime, endTime time.Time) {

	var nodeId int64
	if id, ok := findNode(hashMap, nodeMap, stringId); ok {
		nodeId = id
	} else {
		return // This function is a no-op if we don't have a correct string id
	}

	// This stores whether the package existed in the specified time range
	withinInterval := make(map[int64]bool, len(nodeMap))

	initializeTraversal(g, nodeMap, withinInterval, beginTime, endTime) // Initialize all auxillary data structures for the traversal

	traverseOneNode(g, nodeMap, withinInterval, nodeId)
}

// This function returns the specified node and its dependencies
func GetTransitiveDependenciesNode(g *customgraph.DirectedGraph, nodeMap map[int64]NodeInfo, hashMap map[uint64]int64, stringId string) *[]NodeInfo {
	var nodeId int64
	result := make([]NodeInfo, 0, len(nodeMap)/2)
	if id, ok := findNode(hashMap, nodeMap, stringId); ok {
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

// Get the latest dependencies matching the node's version constraints. If you want this within a specific time frame, use filterNode first
func GetLatestTransitiveDependenciesNode(g *customgraph.DirectedGraph, nodeMap map[int64]NodeInfo, hashMap map[uint64]int64, stringId string) *[]NodeInfo {
	var rootNode NodeInfo
	allDeps := GetTransitiveDependenciesNode(g, nodeMap, hashMap, stringId)
	result := make([]NodeInfo, 0, len(*allDeps)/2)
	if len(*allDeps) > 1 {
		rootNode = (*allDeps)[0]
	} else {
		return &result // No-op if no dependencies were found for whatever reason
	}

	newestPackageVersion := make(map[uint32]NodeInfo, len(*allDeps)/2)

	result = append(result, rootNode)

	// This for loop does the actual filtering
	for _, current := range *allDeps {

		if current.id == rootNode.id {
			continue
		}

		hash := hashPackageName(current.Name)
		currentDate := current.Timestamp
		if latest, ok := newestPackageVersion[hash]; ok {
			latestDate := latest.Timestamp
			if currentDate.After(latestDate) { // If the key exists, and current date is later than the one stored
				newestPackageVersion[hash] = current // Set to the current package
			} else if currentDate.Equal(latestDate) { // If the dates are somehow equal, compare version numbers
				currentversion, _ := semver.NewVersion(current.Version)
				latestVersion, _ := semver.NewVersion(latest.Version)

				if currentversion.GreaterThan(latestVersion) {
					newestPackageVersion[hash] = current
				}
			}
		} else { // If the key doesn't exist yet
			newestPackageVersion[hash] = current
		}
	}

	for _, v := range newestPackageVersion { // Add all latest package versions to the result
		result = append(result, v)
	}

	return &result
}

func keepSelectedNodes(g *customgraph.DirectedGraph, removeIDs map[int64]struct{}) {
	edges := g.Edges()
	for edges.Next() {
		e := edges.Edge()
		fid := e.From().ID()
		tid := e.To().ID()

		if _, ok := removeIDs[fid]; ok {
			g.RemoveEdge(fid, tid)
		}
		if _, ok := removeIDs[tid]; ok {
			g.RemoveEdge(fid, tid)
		}
	}

	for id := range removeIDs {
		g.RemoveNode(id)
	}
}

// Filter the graph between the two given time stamps and then only keep the latest dependencies
func FilterLatestDepsGraph(g *customgraph.DirectedGraph, nodeMap map[int64]NodeInfo, hashMap map[uint64]int64, beginTime, endTime time.Time) {
	filterGraph(g, nodeMap, beginTime, endTime)
	length := g.Nodes().Len() / 2

	keepIDs := make(map[int64]struct{}, length)
	removeIDs := make(map[int64]struct{}, length)
	newestPackageVersion := make(map[uint32]NodeInfo, length)
	v := traverse.DepthFirst{
		Visit: func(n graph.Node) {
			current := nodeMap[n.ID()]
			currentDate := current.Timestamp
			hash := hashPackageName(current.Name)

			if latest, ok := newestPackageVersion[hash]; ok {
				latestDate := latest.Timestamp
				if currentDate.After(latestDate) { // If the key exists, and current date is later than the one stored
					newestPackageVersion[hash] = current // Set to the current package
				} else if currentDate.Equal(latestDate) { // If the dates are somehow equal, compare version numbers
					currentversion, _ := semver.NewVersion(current.Version)
					latestVersion, _ := semver.NewVersion(latest.Version)

					if currentversion.GreaterThan(latestVersion) {
						newestPackageVersion[hash] = current
					}
				}
			} else { // If the key doesn't exist yet
				newestPackageVersion[hash] = current
			}
		},
	}
	numNodes := len(hashMap)
	nodes := g.Nodes()
	i := 0
	fmt.Println()
	for nodes.Next() {
		n := nodes.Node()
		_ = v.Walk(g, n, nil)
		v.Reset()
		i++
		fmt.Printf("\u001b[1A \u001b[2K \r") // Clear the last line
		fmt.Printf("%d / %d subtrees walked \n", i, numNodes)
	}

	for _, v := range newestPackageVersion {
		keepIDs[v.id] = struct{}{}
	}

	for id := range nodeMap {
		if _, ok := keepIDs[id]; !ok { // If the node id was not on the list, kick it out
			removeIDs[id] = struct{}{}
		}
	}

	keepSelectedNodes(g, removeIDs)

}

// This uses the sparse page rank algorithm to find the Page ranks of all nodes
func PageRank(graph *customgraph.DirectedGraph) map[int64]float64 {
	pr := network.PageRankSparse(graph, 0.85, 0.01)
	return pr
}

func Betweenness(graph *customgraph.DirectedGraph) map[int64]float64 {
	return network.Betweenness(graph)
}
