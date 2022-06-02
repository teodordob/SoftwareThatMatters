package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
	"regexp"
	"strings"
	"time"

	g "github.com/AJMBrands/SoftwareThatMatters/graph"
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

	fileSelectionPrompt := &survey.Select{
		Message: "Please select the file you would like to use for the creation of the graph",
		Options: *fileNames,
	}
	file := ""
	err := survey.AskOne(fileSelectionPrompt, &file)
	if err != nil {
		panic(err)
	}
	path := "data/input/" + file

	isUsingMaven := false

	usingMavenPrompt := &survey.Confirm{
		Message: "Is the packages data coming from Maven?",
	}
	err = survey.AskOne(usingMavenPrompt, &isUsingMaven)

	fmt.Println("Creating the graph. This make take a while!")
	if err != nil {
		panic(err)
	}

	//graph, packagesList, stringIDToNodeInfo, idToNodeInfo, nameToVersions := g.CreateGraph(path, isUsingMaven)
	graph, _, stringIDToNodeInfo, idToNodeInfo, _ := g.CreateGraph(path, isUsingMaven)
	// TODO: remove this when we use the actual variables. It is here to get rid of the unused variables warning
	//_, _, _, _, _ = g.CreateGraph(path, isUsingMaven)

	//"View the graph", "View the packages list", "View the packages list with versions", "View the packages list with versions and dependencies"
	stop := false
	for !stop {
		operationIndex := 0
		processPrompt := &survey.Select{
			Message: "What would you like to do now?",
			Options: []string{
				"Find all packages between two timestamps",
				"Find all the possible dependencies of a package",
				"Find all the possible dependencies of a package between two timestamps",
				"Find the most used package",
				"Quit",
			},
		}
		err := survey.AskOne(processPrompt, &operationIndex)

		if err != nil {
			panic(err)
		}

		switch operationIndex {
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
			nodes := findAllDependenciesOfAPackageBetweenTwoTimestamps(graph, idToNodeInfo, stringIDToNodeInfo)
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

func findAllDependenciesOfAPackageBetweenTwoTimestamps(graph *simple.DirectedGraph, nodeMap map[int64]g.NodeInfo, stringIDToNodeInfo map[string]g.NodeInfo) *[]g.NodeInfo {
	beginTime := generateAndRunDatePrompt("Please input the beginning date of the interval (DD-MM-YYYY)")
	endTime := generateAndRunDatePrompt("Please input the end date of the interval (DD-MM-YYYY)")
	nodeStringId := generateAndRunPackageNamePrompt("Please select the name and the version of the package", stringIDToNodeInfo)
	g.FilterGraph(graph, nodeMap, beginTime, endTime)
	return g.GetTransitiveDependenciesNode(graph, nodeMap, stringIDToNodeInfo, nodeStringId)
}

func generateAndRunDatePrompt(message string) time.Time {
	validateDate := func(input interface{}) error {
		str, ok := input.(string)
		if !ok {
			return errors.New("input is not a string")
		}
		if len(str) == 0 {
			return errors.New("input cannot be empty")
		}
		matched, _ := regexp.MatchString("\\d{2}-\\d{2}-\\d{4}", str)
		if !matched {
			return errors.New("input must be in the format: DD-MM-YYYY")
		}
		if len(str) == 10 {
			_, err := time.Parse("02-01-2006", str)
			if err != nil {
				return errors.New("input must be a valid date")
			}
		}
		return nil
	}

	timePrompt := &survey.Input{
		Message: message,
	}
	timeString := ""
	err := survey.AskOne(timePrompt, &timeString, survey.WithValidator(validateDate))

	if err != nil {
		panic(err)
	}

	time, _ := time.Parse("02-01-2006", timeString)
	return time

}

func generateAndRunPackageNamePrompt(message string, stringIDToNodeInfo map[string]g.NodeInfo) string {
	keys := make([]string, 0, len(stringIDToNodeInfo))
	for key := range stringIDToNodeInfo {
		keys = append(keys, key)
	}
	packagePrompt := &survey.Select{
		Message: message,
		Options: keys,
	}

	//packagePrompt := &survey.Input{
	//	Message: message,
	//}
	packageID := ""
	err := survey.AskOne(packagePrompt, &packageID)

	if err != nil {
		panic(err)
	}

	return packageID
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
