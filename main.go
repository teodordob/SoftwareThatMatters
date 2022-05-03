package main

import (
	"runtime"

	"github.com/AJMBrands/SoftwareThatMatters/ingest"
)

const limited_discovery_query string = "https://libraries.io/api/search?api_key=3dc75447d3681ffc2d17517265765d23&page=1&per_page=1&platforms=NPM"

const discovery_query string = "https://libraries.io/api/search?api_key=3dc75447d3681ffc2d17517265765d23&platforms=NPM&per_page=20"

const offline_file string = "data/100packages.json"
const outPath string = "data/out/parsed_data.csv"

var m1, m2 runtime.MemStats

//TODO: Make ingest process and file writing scalable
func main() {
	runtime.ReadMemStats(&m1) // Reading memory stats for debug purposes
	//ingestResultAddr := ingest.Ingest(limited_discovery_query)
	ingest.IngestFile(offline_file, outPath)
	runtime.ReadMemStats(&m2)
}
