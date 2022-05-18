package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/AJMBrands/SoftwareThatMatters/ingest"
)

const inPath string = "data/in/input.json"

const outPathTemplate string = "data/out/streamedout-%s.json"

var m1, m2 runtime.MemStats
var t1, t2 time.Time

//TODO: Make ingest process and file writing scalable
func main() {
	runtime.ReadMemStats(&m1) // Reading memory stats for debug purposes
	t1 = time.Now()
	//ingest.IngestFile(offline_file, outPath)
	// ingest.Ingest(limited_discovery_query, outPathTemplate, versionPath)

	amount := ingest.StreamParse(inPath, outPathTemplate)
	ingest.MergeJSON(outPathTemplate, amount)
	runtime.ReadMemStats(&m2)
	t2 = time.Now()
	fmt.Printf("Took %d ms\n", t2.UnixMilli()-t1.UnixMilli())
}
