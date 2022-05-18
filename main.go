package main

import (
	"fmt"

	g "github.com/AJMBrands/SoftwareThatMatters/graph"
)

func main() {
	//x := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.DependenciesInfo{"2020-10-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
	//	"org.hamcrest:hamcrest-library": "1.3"}}}
	x := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.DependenciesInfo{"3.8.1": {"2020-10-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
		"org.hamcrest:hamcrest-library": "1.3"}}}}
	y := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.DependenciesInfo{"3.8.2": {"2021-10-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
		"org.hamcrest:hamcrest-library": "1.3"}}}}
	z := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.DependenciesInfo{"3.8.3": {"2021-11-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
		"org.hamcrest:hamcrest-library": "1.4"}}}}
	as := g.PackageInfo{Name: "junit:junit", Versions: map[string]g.DependenciesInfo{"4.0.2": {"2022-10-11T15:19:50", map[string]string{"org.hamcrest:hamcrest-core": "1.3",
		"org.hamcrest:hamcrest-library": "2.0"}}}}
	myarr := []g.PackageInfo{x, y, z}
	m := g.CreateMap(&myarr)
	g.AddElementToMap(as, m)
	g1 := g.CreateGraph(m)
	fmt.Println(g1)
	//fmt.Println(m)
}
