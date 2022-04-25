package main

import (
	"github.com/AJMBrands/SoftwareThatMatters/ingest"
)

var packagesAndVersions []ingest.DiscoveryResponse

func main() {
	packagesAndVersions = ingest.Ingest()
}
