package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"

	"github.com/AJMBrands/SoftwareThatMatters/customgraph"
	g "github.com/AJMBrands/SoftwareThatMatters/graph"
	"github.com/spf13/cobra"
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

	fmt.Println("Creating the graph. This may take a while!")
	if err != nil {
		panic(err)
	}

	//graph, packagesList, stringIDToNodeInfo, idToNodeInfo, nameToVersions := g.CreateGraph(path, isUsingMaven)
	graph, hashMap, idToNodeInfo, _ := g.CreateGraph(path, isUsingMaven)

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
				"Find the latest dependencies of a package between two timestamps",
				"Find the n most used packages by PageRank between two time stamps",
				"Find the n most used packages by PageRank (without considering time)",
				"Find the n most used packages by node betweenness",
				"Find the n most used packages by node betweenness (without considering time)",
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
			name := generateAndRunPackageNamePrompt("Please input the package name", idToNodeInfo)
			nodes := g.GetTransitiveDependenciesNode(graph, idToNodeInfo, hashMap, name)
			for _, node := range *nodes {
				fmt.Println(node)
			}

		case 2:
			fmt.Println("This should find all the possible dependencies of a package between two timestamps")
			nodes := findAllDependenciesOfAPackageBetweenTwoTimestamps(graph, hashMap, idToNodeInfo)
			for _, node := range *nodes {
				fmt.Println(node)
			}
		case 3:
			fmt.Println("This should find the latest dependencies of a package")
			nodes := findLatestDependenciesOfAPackage(graph, hashMap, idToNodeInfo)

			for _, node := range *nodes {
				fmt.Println(node)
			}
		case 4:
			findMostUsedPackages(graph, idToNodeInfo, true)

		case 5:
			findMostUsedPackages(graph, idToNodeInfo, false)

		case 6:
			//findMostUsedPackages(graph, idToNodeInfo, true)
			fmt.Println("TODO: Reimlpement the betweenes by time stamp")

		case 7:
			fmt.Println("This should find the n most used packages (by node betweenness)")
			fmt.Println("Running node betweenness algorithm")
			b := g.Betweenness(graph)
			keys := make([]int64, 0, len(b))
			for k := range b {
				keys = append(keys, k)
			}

			sort.SliceStable(keys, func(i, j int) bool {
				return b[keys[i]] > b[keys[j]]
			})

			count := generateAndRunNumberPrompt("Please select the number (n > 0) of highest-ranked packages you wish to see")
			for i := 0; i < count; i++ {
				fmt.Printf("The number %d highest-ranked node (%v) has rank %f \n", i, idToNodeInfo[keys[i]], b[keys[i]])
			}
		case 8:
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
		nodeTime := node.Timestamp
		if g.InInterval(nodeTime, beginTime, endTime) {
			nodesInInterval = append(nodesInInterval, node)
		}
	}

	return &nodesInInterval

}

func findMostUsedPackages(graph *customgraph.DirectedGraph, idToNodeInfo map[int64]g.NodeInfo, betweenTimestamps bool) {
	if betweenTimestamps {
		beginTime := generateAndRunDatePrompt("Please input the beginning date of the interval (DD-MM-YYYY)")
		endTime := generateAndRunDatePrompt("Please input the end date of the interval (DD-MM-YYYY)")
		fmt.Println("Getting the latest dependencies for packages. This will take a while")
		t1 := time.Now().Unix()
		g.FilterNoTraversal(graph, idToNodeInfo, beginTime, endTime)
		t2 := time.Now().Unix()
		fmt.Printf("Graph filtering took %d seconds\n", t2-t1)
	}

	fmt.Println("Running PageRank")
	pr := g.PageRank(graph)
	keys := make([]int64, 0, len(pr))
	aggregated := make(map[string]float64)

	for k, value := range pr {
		keys = append(keys, k)
		aggregated[idToNodeInfo[k].Name] += value
	}

	aggregatedKeys := make([]string, 0, len(aggregated))

	for k := range aggregated {
		aggregatedKeys = append(aggregatedKeys, k)
	}

	sort.SliceStable(aggregatedKeys, func(i, j int) bool {
		return aggregated[aggregatedKeys[i]] > aggregated[aggregatedKeys[j]]
	})

	sort.SliceStable(keys, func(i, j int) bool {
		return pr[keys[i]] > pr[keys[j]]
	})

	count := generateAndRunNumberPrompt("Please select the number (n > 0) of highest-ranked packages you wish to see")
	for i := 0; i < count; i++ {
		fmt.Printf("The %d-th highest-ranked node (%v) has rank %f \n", i, idToNodeInfo[keys[i]], pr[keys[i]])
	}

	fmt.Print("\n---------------------------------------------\n\n")

	for i := 0; i < count; i++ {
		fmt.Printf("The %d-th highest-ranked package (%v) has rank %f \n", i, aggregatedKeys[i], aggregated[aggregatedKeys[i]])
	}
}

func findAllDependenciesOfAPackageBetweenTwoTimestamps(graph *customgraph.DirectedGraph, hashMap map[uint64]int64, nodeMap map[int64]g.NodeInfo) *[]g.NodeInfo {
	beginTime := generateAndRunDatePrompt("Please input the beginning date of the interval (DD-MM-YYYY)")
	endTime := generateAndRunDatePrompt("Please input the end date of the interval (DD-MM-YYYY)")
	nodeStringId := generateAndRunPackageNamePrompt("Please select the name and the version of the package", nodeMap)
	g.FilterNoTraversal(graph, nodeMap, beginTime, endTime)
	return g.GetTransitiveDependenciesNode(graph, nodeMap, hashMap, nodeStringId)
}

func findLatestDependenciesOfAPackage(graph *customgraph.DirectedGraph, hashMap map[uint64]int64, nodeMap map[int64]g.NodeInfo) *[]g.NodeInfo {
	nodeStringId := generateAndRunPackageNamePrompt("Please select the name and the version of the package", nodeMap)
	g.LatestNoTraversal(graph, nodeMap, hashMap)
	return g.GetTransitiveDependenciesNode(graph, nodeMap, hashMap, nodeStringId)
}

func generateAndRunNumberPrompt(message string) int {
	validateNumber := func(input any) error {
		var num int
		if str, ok := input.(string); ok {
			if n, err := strconv.Atoi(str); err != nil {
				return errors.New("Input couldn't be parsed to integer")
			} else {
				num = n
			}
		} else {
			return errors.New("Input is not even a string")
		}
		if num <= 0 {
			return errors.New("Input must be a number larger than 0")
		} else {
			return nil
		}
	}

	numberPrompt := &survey.Input{Message: message}
	var number int = -1
	err := survey.AskOne(numberPrompt, &number, survey.WithValidator(validateNumber))

	if err != nil {
		panic(err)
	}
	return number
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

func generateAndRunPackageNamePrompt(message string, stringIDToNodeInfo map[int64]g.NodeInfo) string {
	names := make([]string, 0, len(stringIDToNodeInfo))
	for _, node := range stringIDToNodeInfo {
		name := fmt.Sprintf("%s-%s", node.Name, node.Version)
		names = append(names, name)
	}
	packagePrompt := &survey.Select{
		Message: message,
		Options: names,
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
