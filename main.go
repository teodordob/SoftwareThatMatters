package main

import (
	"fmt"

	g "github.com/AJMBrands/SoftwareThatMatters/graph"
)

func main() {
	//x := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.VersionInfo{"2020-10-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
	//	"org.hamcrest:hamcrest-library": "1.3"}}}
	//x := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.VersionInfo{"3.8.1": {"2020-10-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
	//	"org.hamcrest:hamcrest-library": "1.3"}}}}
	//y := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.VersionInfo{"3.8.2": {"2021-10-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
	//	"org.hamcrest:hamcrest-library": "1.3"}}}}
	//z := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.VersionInfo{"3.8.3": {"2021-11-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
	//	"org.hamcrest:hamcrest-library": "1.4"}}}}
	//as := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.VersionInfo{"4.0.2": {"2022-10-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
	//	"org.hamcrest:hamcrest-library": "2.0"}}}}
	//myarr := []g.PackageInfo{x, y, z}
	//m := g.CreateMap(&myarr)
	//g.AddElementToMap(as, m)
	//g1 := g.CreateGraph(m)
	//fmt.Println(g1)

	parsed := g.ParseJSON("data/input/test_data.json")
	m2 := g.CreateMap(parsed)
	nameToIdMap := g.CreateNameToIDMap(m2)
	nameToVersionMap := g.CreateNameToVersionMap(parsed)
	g2 := g.CreateGraph(m2)
	//fmt.Println(g2)
	g.CreateEdges(g2, parsed, nameToIdMap, nameToVersionMap)

	//Uncomment this to create the visualization and use these commands in the dot file
	//Toggle Preview - ctrl+shift+v (Mac: cmd+shift+v)
	//Open Preview to the Side - ctrl+k v (Mac: cmd+k shift+v)
	//g.Visualization(g2, "Test")

	fmt.Println(m2)
	fmt.Println(nameToIdMap)
	//fmt.Println(m2)
	//fmt.Println(nameToIdMap)
}
