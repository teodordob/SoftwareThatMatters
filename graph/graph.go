package graph

import (
	"encoding/json"
	"fmt"
	semver "github.com/Masterminds/semver"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"log"
	"os"
	"regexp"
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

// CreateStringIDToNodeInfoMap takes a list of PackageInfo and a simple.DirectedGraph. For each of the packages,
// it creates a mapping of stringIDs to NodeInfo and also adds a node to the graph. The handling of the IDs is delegated
// to Gonum. These IDs are also included in the mapping for ease of access.
func CreateStringIDToNodeInfoMap(packagesInfo *[]PackageInfo, graph *simple.DirectedGraph) *map[string]NodeInfo {
	stringIDToNodeInfoMap := make(map[string]NodeInfo, len(*packagesInfo))
	for _, packageInfo := range *packagesInfo {
		for packageVersion, versionInfo := range packageInfo.Versions {
			packageNameVersionString := fmt.Sprintf("%s-%s", packageInfo.Name, packageVersion)
			// Delegate the work of creating a unique ID to Gonum
			newNode := graph.NewNode()
			stringIDToNodeInfoMap[packageNameVersionString] = *NewNodeInfo(newNode.ID(), packageInfo.Name, packageVersion, versionInfo.Timestamp)
			graph.AddNode(newNode)
		}
	}
	return &stringIDToNodeInfoMap
}

func CreateNameToVersionMap(m *[]PackageInfo) *map[string][]string {
	newMap := make(map[string][]string, len(*m))
	for _, value := range *m {
		name := value.Name
		for k := range value.Versions {
			newMap[name] = append(newMap[name], k)
		}
	}
	return &newMap
}

//Function to write the simple graph to a dot file so it could be visualized with GraphViz. This includes only Ids
func Visualization(graph *simple.DirectedGraph, name string) {
	result, _ := dot.Marshal(graph, name, "", "  ")

	file, err := os.Create(name + ".dot")

	if err != nil {
		log.Fatal("Error!", err)
	}
	defer file.Close()

	fmt.Fprintf(file, string(result))

}

//Writes to dot file manually from the NodeInfoMap to include the Node info in the graphViz
//TODO: Optimize in the future since this is kind of barbaric probably there is a faster way.
func VisualizationNodeInfo(iDToNodeInfo *map[string]NodeInfo, graph *simple.DirectedGraph, name string) {
	file, err := os.Create(name + ".dot")
	d1 := []byte("strict digraph" + " " + name + " " + "{\n")
	d2 := []byte("}")
	lab := string("[label = \" ")
	edgIt := graph.Edges()

	fmt.Fprintf(file, string(d1))

	for key, element := range *iDToNodeInfo {
		//fmt.Println("Key:", key, "=>", "Element:", element.id)
		fmt.Fprintf(file, fmt.Sprint(element.id)+lab+string(key)+` \n `+string(element.Version)+` \n `+string(element.Timestamp)+"\""+"];\n")

	}

	for edgIt.Next() {
		fmt.Fprintf(file, fmt.Sprint(edgIt.Edge().From().ID())+" -> "+fmt.Sprint(edgIt.Edge().To().ID())+";\n")
	}

	fmt.Fprintf(file, string(d2))

	if err != nil {
		panic(err)
	}

	defer file.Close()

}

// CreateEdges takes a graph, a list of packages and their dependencies, a map of stringIDs to NodeInfo and
// a map of names to versions and creates directed edges between the dependent library and its dependencies.
// TODO: add documentation on how we use semver for edges
// TODO: Discuss removing pointers from maps since they are reference types without the need of using * : https://stackoverflow.com/questions/40680981/are-maps-passed-by-value-or-by-reference-in-go
func CreateEdges(graph *simple.DirectedGraph, inputList *[]PackageInfo, stringIDToNodeInfo *map[string]NodeInfo, nameToVersionMap *map[string][]string) {
	packagesInfo := *inputList // Dereferencing here results in copying the whole list. Maybe we can just use the dereferencing without the assigning as to avoid copying things
	nameToVersion := *nameToVersionMap
	r, _ := regexp.Compile("((?P<open>[\\(\\[])(?P<bothVer>((?P<firstVer>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)(?P<comma>,)(?P<secondVer>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)?)|((?P<comma2>,)?(?P<secondVer2>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+)))?)?))(?P<closing>[\\)\\]]))|(?P<simplvers>(0|[1-9]+)(\\.(0|[1-9]+)(\\.(0|[1-9]+))?)?)")
	for id, packageInfo := range packagesInfo {
		for _, dependencyInfo := range packageInfo.Versions {
			for dependencyName, dependencyVersion := range dependencyInfo.Dependencies {
				mvndep := translateMavenSemver(dependencyVersion, r)
				fmt.Println(mvndep)
				constraint, err := semver.NewConstraint(dependencyVersion)
				//c, err := semver2.ParseRange(dependencyVersion)
				if err != nil {
					log.Fatal(err)
				}
				for _, v := range nameToVersion[dependencyName] {
					//newVersion, _ := semver2.Parse(v)
					newVersion, _ := semver.NewVersion(v)
					if constraint.Check(newVersion) {
						dependencyNameVersionString := fmt.Sprintf("%s-%s", dependencyName, v)
						dependencyNode := graph.Node((*stringIDToNodeInfo)[dependencyNameVersionString].id)
						packageNode := graph.Node(int64(id))
						graph.SetEdge(simple.Edge{F: packageNode, T: dependencyNode})
					}
				}
			}
		}
	}
}

func translateMavenSemver(s string, reg *regexp.Regexp) string {
	match := reg.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
		fmt.Printf("by name: %s %s\n", result["singur"])
	}

	return result["primavers"]

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
