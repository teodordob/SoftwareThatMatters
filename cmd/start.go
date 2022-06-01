package cmd

import (
	"fmt"
	g "github.com/AJMBrands/SoftwareThatMatters/graph"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the application and ask guides you through the process of generating a graph",
	Long:  `Starts the application and ask guides you through the process of generating a graph`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

// start is the main function that starts the application. It asks the user for the data file and then generates the graph.
// After the graph is generated, it asks the user how they want to proceed. The loop is done to allow the user to run
// multiple requests on the same graph. This means that the graph can be generated once, and then it can be processed
// multiple times.
func start() {

	//validate := func(input string) error {
	//	if len(input) == 0 {
	//		return errors.New("input cannot be empty")
	//	}
	//	return nil
	//}

	fileNames := getJSONFilesFromDataFolder()
	if len(*fileNames) == 0 {
		fmt.Println("No JSON files found in data folder! Make sure there is at least one file in the data/input folder.")
		return
	}

	fileSelectionPrompt := promptui.Select{
		Label: "Please select the file you would like to use for the creation of the graph",
		Items: *fileNames,
	}
	_, file, err := fileSelectionPrompt.Run()
	path := "data/input/" + file

	if err != nil {
		fmt.Printf("Something went wrong!%v\n", err)
	}

	usingMavenPrompt := promptui.Select{
		Label: "Is the packages data coming from Maven?",
		Items: []string{"Yes", "No"},
	}
	_, isUsingMavenString, err := usingMavenPrompt.Run()
	if err != nil {
		fmt.Printf("Something went wrong!%v\n", err)
	}
	isUsingMaven := strings.ToLower(isUsingMavenString) == "yes"
	fmt.Println("Creating the graph. This make take a while!")

	//graph, packagesList, stringIDToNodeInfo, idToNodeInfo, nameToVersions := g.CreateGraph(path, isUsingMaven)
	// TODO: remove this when we use the actual variables. It is here to get rid of the unused variables warning
	_, _, _, _, _ = g.CreateGraph(path, isUsingMaven)

	//"View the graph", "View the packages list", "View the packages list with versions", "View the packages list with versions and dependencies"
	stop := false
	for !stop {
		processPrompt := promptui.Select{
			Label: "What would you like to do now?",
			Items: []string{"Find all the possible dependencies of a package",
				"Find all the possible dependencies of a package between two timestamps",
				"Find the most used package",
				"Quit"},
		}
		processIndex, _, err := processPrompt.Run()
		if err != nil {
			fmt.Printf("Something went wrong!%v\n", err)
		}

		switch processIndex {
		//Find all possible dependencies of a package
		case 0:
			fmt.Println("This should find all the possible dependencies of a package")
		case 1:
			fmt.Println("This should find all the possible dependencies of a package between two timestamps")
		case 2:
			fmt.Println("This should find the most used package")
		case 3:
			fmt.Println("Stopping the program...")
			stop = true
		}

	}

}

// getJSONFilesFromDataFolder returns a slice of strings with the names of the JSON files in the data folder. It can
// return an empty slice if there are no JSON files in the data folder so a check should be done after using this
func getJSONFilesFromDataFolder() *[]string {

	dir, err := os.Open("data/input")
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}
	var fileNames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			fileNames = append(fileNames, file.Name())
		}

	}
	return &fileNames
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
