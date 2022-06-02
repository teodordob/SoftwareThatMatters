package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	g "github.com/AJMBrands/SoftwareThatMatters/graph"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/graph/simple"
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
	graph, _, stringIDToNodeInfo, idToNodeInfo, _ := g.CreateGraph(path, isUsingMaven)
	// TODO: remove this when we use the actual variables. It is here to get rid of the unused variables warning
	//_, _, _, _, _ = g.CreateGraph(path, isUsingMaven)

	//"View the graph", "View the packages list", "View the packages list with versions", "View the packages list with versions and dependencies"
	stop := false
	for !stop {
		processPrompt := promptui.Select{
			Label: "What would you like to do now?",
			Items: []string{
				"Find all packages between two timestamps",
				"Find all the possible dependencies of a package",
				"Find all the possible dependencies of a package between two timestamps",
				"Find the most used package",
				"Quit",
			},
		}
		processIndex, _, err := processPrompt.Run()
		if err != nil {
			fmt.Printf("Something went wrong!%v\n", err)
		}

		switch processIndex {
		case 0:
			fmt.Println("This should find all the packages between two timestamps")
			nodes := findAllPackagesBetweenTwoTimestamps(idToNodeInfo)
			for _, node := range *nodes {
				fmt.Println(node)
			}
		case 1:
			fmt.Println("This should find all the possible dependencies of a package")
			name := generateAndRunPackageNamePrompt("Please input the package name", stringIDToNodeInfo)
			nodes := g.GetTransitiveDependenciesNode(graph, idToNodeInfo, stringIDToNodeInfo, name)
			for _, node := range *nodes {
				fmt.Println(node)
			}

		case 2:
			fmt.Println("This should find all the possible dependencies of a package between two timestamps")
			nodes := findAllDepedenciesOfAPackageBetweenTwoTimestamps(graph, idToNodeInfo, stringIDToNodeInfo)
			for _, node := range *nodes {
				fmt.Println(node)
			}
		case 3:
			fmt.Println("This should find the most used package")
		case 4:
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

func findAllPackagesBetweenTwoTimestamps(idToNodeInfo map[int64]g.NodeInfo) *[]g.NodeInfo {
	//// TODO: Discuss if we should create a copy or not. My idea is that we should create a copy of the graph and then
	//// TODO: use the copy to find the packages. This way we can use the original graph for other operations.
	//graphCopy := *graph

	beginTime := generateAndRunDatePrompt("Please input the beginning date of the interval (DD-MM-YYYY)")
	endTime := generateAndRunDatePrompt("Please input the end date of the interval (DD-MM-YYYY)")

	var nodesInInterval []g.NodeInfo

	for _, node := range idToNodeInfo {
		//TODO: We need a way of properly parsing multiple times
		nodeTime, err := time.Parse(time.RFC3339, node.Timestamp)
		if err != nil {
			fmt.Println("There was an error parsing the timestamps in the nodes!")
			panic(err)
		}
		if g.InInterval(nodeTime, beginTime, endTime) {
			nodesInInterval = append(nodesInInterval, node)
		}
	}

	return &nodesInInterval

}

func findAllDepedenciesOfAPackageBetweenTwoTimestamps(graph *simple.DirectedGraph, nodeMap map[int64]g.NodeInfo, stringIDToNodeInfo map[string]g.NodeInfo) *[]g.NodeInfo {
	beginTime := generateAndRunDatePrompt("Please input the beginning date of the interval (DD-MM-YYYY)")
	endTime := generateAndRunDatePrompt("Please input the end date of the interval (DD-MM-YYYY)")
	nodeStringId := generateAndRunPackageNamePrompt("Please input the name and the version of the package (name-version)", stringIDToNodeInfo)
	g.FilterGraph(graph, nodeMap, beginTime, endTime)
	return g.GetTransitiveDependenciesNode(graph, nodeMap, stringIDToNodeInfo, nodeStringId)
}

func generateAndRunDatePrompt(label string) time.Time {
	validateDate := func(input string) error {
		if len(input) == 0 {
			return errors.New("input cannot be empty")
		}
		matched, _ := regexp.MatchString("\\d{2}-\\d{2}-\\d{4}", input)
		if !matched {
			return errors.New("input must be in the format: DD-MM-YYYY")
		}
		if len(input) == 10 {
			_, err := time.Parse("02-01-2006", input)
			if err != nil {
				return errors.New("input must be a valid date")
			}
		}
		return nil
	}

	timePrompt := promptui.Prompt{
		Label:    label,
		Validate: validateDate,
	}

	timeString, err := timePrompt.Run()
	if err != nil {
		panic(err)
	}

	time, _ := time.Parse("02-01-2006", timeString)
	return time

}

func generateAndRunPackageNamePrompt(label string, stringIDToNodeInfo map[string]g.NodeInfo) string {
	validateString := func(input string) error {
		if len(input) == 0 {
			return errors.New("input cannot be empty")
		}
		if _, ok := stringIDToNodeInfo[input]; ok {
			return nil
		} else {
			return errors.New("String id was not found \n")

		}
	}

	packagePrompt := promptui.Prompt{
		Label:    label,
		Validate: validateString,
	}
	packageId, err := packagePrompt.Run()
	if err != nil {
		panic(err)
	}
	return packageId
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
