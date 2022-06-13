package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/AJMBrands/SoftwareThatMatters/ingest"
	"github.com/spf13/cobra"
)

const outPathTemplate string = "data/out/out-%s.json"

var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Starts the ingest stage for the application and preprocesses the json input",
	Long:  "Starts the ingest stage for the application and preprocesses the json input",
	Run: func(cmd *cobra.Command, args []string) {
		ingestFile()
	},
}

func ingestFile() {
	_, err := os.Stat("data/in/input.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting ingest process...")
	t1 := time.Now().UnixMilli()
	numPackages := ingest.StreamParse("data/in/input.json", outPathTemplate)
	ingest.MergeJSON(outPathTemplate, numPackages)
	t2 := time.Now().UnixMilli()
	fmt.Printf("Done! Took %d ms", t2-t1)
}

func init() {
	rootCmd.AddCommand(ingestCmd)
}
