package ingest

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	ccsv "github.com/tsak/concurrent-csv-writer"
)

type OutputVersion struct {
	TimeStamp    time.Time         `json:"timestamp"`
	Dependencies map[string]string `json:"dependencies"`
}

type OutputFormat struct {
	Name     string                   `json:"name"`
	Versions map[string]OutputVersion `json:"versions"`
}

type VersionData struct {
	Version         string            `json:"version"`
	DevDependencies map[string]string `json:"devDependencies"`
	Dependencies    map[string]string `json:"dependencies"`
}

type Doc struct {
	Name     string                 `json:"name"`
	Versions map[string]VersionData `json:"versions"`
	Time     map[string]CreatedTime `json:"time"`
}

type Entry struct {
	Doc Doc `json:"doc"`
}

//TODO: Add method to put resolved dependencies back into JSON and output to file

// Type alias so we can create a custom parser for time since it wasn't parsed correctly natively
type CreatedTime time.Time

// Function required to implement the JSON parser interface for CreatedTime
func (ct *CreatedTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*ct = CreatedTime(t)
	return nil
}

// Function required to implement the JSON parser interface for CreatedTime
func (ct CreatedTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct)
}

// This function forces the JSON Unmarshaler to use the CreatedTime Unmarshaler
func (v *Version) UnmarshalJSON(b []byte) error {
	var dat map[string]interface{}

	if err := json.Unmarshal(b, &dat); err != nil {
		return err
	}
	date_string := "\"" + dat["published_at"].(string) + "\""
	date_json := []byte(date_string)
	var date CreatedTime

	if err := json.Unmarshal(date_json, &date); err != nil {
		return err
	}

	*v = Version{dat["number"].(string), date}
	return nil
}

func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

func (ct CreatedTime) String() string {
	return time.Time(ct).Format(time.RFC3339Nano)
}

type Version struct {
	Number      string      `json:"number"`
	PublishedAt CreatedTime `json:"published_at"`
}

type PackageInfo struct {
	Name     string    `json:"name"`
	Versions []Version `json:"versions"`
}

type VersionDependencies struct {
	Name           string
	Version        string
	VersionCreated time.Time
	Dependencies   []Dependency
}

type VersionInfo struct {
	Dependencies    map[string]interface{} `json:"dependencies"`
	DevDependencies map[string]interface{} `json:"devDependencies"`
}

type Dependency struct {
	Name            string
	RequiredVersion string
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s:%s", d.Name, d.RequiredVersion)
}

// Ingest live data
func Ingest(query string, outPathTemplate, versionPath string) {
	rawDataAddr, requestAddr := request(query)

	if statusCode := requestAddr.StatusCode; statusCode != 200 {
		log.Fatalln("Uh-oh, HTTP status was: ", requestAddr.Status)
	}
	ingestInternal(*rawDataAddr, outPathTemplate, versionPath)
}

// Ingest (partially) offline data
func IngestFile(file string, outPathTemplate, versionPath string) {
	inputBytes, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println("Something went wrong with reading the file:")
		panic(err)
	}

	ingestInternal(inputBytes, outPathTemplate, versionPath)
}

func ingestInternal(inputBytes []byte, outPathTemplate, versionPath string) {
	var wg sync.WaitGroup

	var arr []PackageInfo
	if err := json.Unmarshal(inputBytes, &arr); err != nil {
		fmt.Println("JSON parsing went wrong:")
		panic(err)
	}

	fmt.Println("Got data from input")
	versionPrinter(&arr, versionPath)
	fmt.Println("Processing...")

	// result := make(chan *[]VersionDependencies)
	length := len(arr)
	// TODO: Find smarter way to divide input into threads?
	for i := 0; i < length; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			path := fmt.Sprintf(outPathTemplate, i)
			process(arr[i:i+1], path)
		}(i)
	}

	wg.Wait() // Wait for all goroutines to be done
}

func versionPrinter(input *[]PackageInfo, versionPath string) {
	var wg sync.WaitGroup

	csv, err := ccsv.NewCsvWriter(versionPath)

	defer csv.Close()
	defer wg.Wait()

	if err != nil {
		panic(fmt.Sprintln("Couldn't open ", versionPath))
	}

	for i, p := range *input {
		wg.Add(1)
		go func(i int, p PackageInfo) {
			defer wg.Done()
			printSinglePackage(&p, csv)
		}(i, p)
	}

}

func printSinglePackage(packageAddr *PackageInfo, writer *ccsv.CsvWriter) {
	p := *packageAddr
	name, versions := p.Name, p.Versions

	for _, ver := range versions {
		num, date := ver.Number, ver.PublishedAt
		writer.Write([]string{name, num, time.Time(date).Format(time.RFC3339Nano)})
	}
}

func request(req string) (*[]byte, *http.Response) {
	// fmt.Println("Starting request...")
	resp, err := http.Get(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	return &body, resp
}

func process(input []PackageInfo, outPath string) *[]VersionDependencies {
	var result []VersionDependencies
	inputLength := len(input)

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	w := csv.NewWriter(file)

	for packageIdx := range input {
		p := &input[packageIdx]
		name, versionsAddr := p.Name, &p.Versions
		for verIdx := range *versionsAddr {
			version := (*versionsAddr)[verIdx]
			number, date := version.Number, version.PublishedAt
			currentURL := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, number)

			rawDataAddr, responseAddr := request(currentURL)
			var parsed VersionInfo

			if err := json.Unmarshal(*rawDataAddr, &parsed); err != nil {
				statusCode := responseAddr.StatusCode
				if statusCode == 404 { // This package's dependencies were not found, so try the next one
					fmt.Printf("The following package's dependencies weren't found: \"%s\" version \"%s\"\n", name, number)
					continue
				} else {
					status := responseAddr.Status
					fmt.Println("Http status code was: ", status) // This will probably be a rate-limit status code
					panic(err)
				}
			}

			deps, devDeps := parsed.Dependencies, parsed.DevDependencies
			allDependencies := make([]Dependency, 0, len(deps)+len(devDeps))
			for k, v := range deps {
				allDependencies = append(allDependencies, Dependency{k, v.(string)})
			}
			for k, v := range devDeps {
				allDependencies = append(allDependencies, Dependency{k, v.(string)})
			}
			versionDeps := VersionDependencies{name, number, time.Time(date), allDependencies}
			//fmt.Println(versionDeps)
			result = append(result, versionDeps)
			writeOneToFile(&versionDeps, w)

			if verIdx%10 == 0 { // Flush writer every 10 entries
				w.Flush()
			}
		}
		w.Flush() // Flush at the end to make sure there's no data left
		fmt.Printf("Package dependencies of %s (%d of %d) fully resolved \n", name, packageIdx+1, inputLength)
	}
	return &result
}

func writeOneToFile(input *VersionDependencies, csvWriter *csv.Writer) {
	name, version, date, deps := (*input).Name, (*input).Version, (*input).VersionCreated, (*input).Dependencies

	var depsString string

	l := len(deps)

	for i, dep := range deps {
		depsString += fmt.Sprint(dep)
		if i < l-1 {
			depsString += ";"
		}
	}
	depsString = fmt.Sprintf("[%s]", depsString)

	if err := csvWriter.Write([]string{name, version, date.Format(time.RFC3339Nano), depsString}); err != nil {
		log.Fatal(err)
	}
}

// Resolve semantic versions in parsed data CSV files using date and semantic version constraints
func ResolveVersions(versionPath string, parsedDepsPathTemplate string, outPathTemplate string) {
	//TODO: Find version that satisfies both of these requirements: Dependency satisfies semver constraints; Dependency was released before package
}

//TODO: output a JSON file per package
func StreamParse(inPath string, jsonOutPathTemplate string) int {
	fmt.Println("Starting input JSON parser...")
	var wg sync.WaitGroup
	f, _ := os.Open(inPath)
	dec := json.NewDecoder(f)

	// versionPath := strings.Replace(outPath, ".", ".versions.", 1)
	// versionFile, err := os.OpenFile(versionPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// versionWriter := csv.NewWriter(versionFile)

	// Read opening bracket
	if _, err := dec.Token(); err != nil {
		log.Fatal(err)
	}
	// While the decoder says there is more to parse, parse a JSON entries and print them one-by-one
	i := 0
	for dec.More() {
		var e Entry

		if err := dec.Decode(&e); err != nil {
			log.Fatal(err)
		}
		timeStamps := e.Doc.Time

		var vds []VersionDependencies = make([]VersionDependencies, 0, len(e.Doc.Versions))

		for number, vd := range e.Doc.Versions {
			t := time.Time(timeStamps[number])
			deps, devDeps := vd.Dependencies, vd.DevDependencies
			allDependencies := make([]Dependency, 0, len(deps)+len(devDeps))
			for k, v := range deps {
				allDependencies = append(allDependencies, Dependency{k, v})
			}
			for k, v := range devDeps {
				allDependencies = append(allDependencies, Dependency{k, v})
			}

			vd := VersionDependencies{e.Doc.Name, number, t, allDependencies}
			vds = append(vds, vd)
			// versionWriter.Write([]string{e.Doc.Name, number}) // Write version to separate file
		}
		jsonPath := fmt.Sprintf(jsonOutPathTemplate, fmt.Sprint(i)) // Append a number to filePath
		wg.Add(1)                                                   // Tell the WaitGroup it needs to wait for one more
		go func(vds *[]VersionDependencies, jsonPath string) {
			defer wg.Done() // Tell the WaitGroup this task is done after the function below is done
			writeToFileJSON(vds, jsonPath)
		}(&vds, jsonPath)

		// versionWriter.Flush()
		i++
	}
	// Read closing bracket
	if _, err := dec.Token(); err != nil {
		log.Fatal(err)
	}
	wg.Wait() // Wait for all subroutines to be done
	fmt.Println("JSON parsing done")
	return i
}

func writeToFileJSON(vdAddr *[]VersionDependencies, outPath string) {
	outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0644)

	vds := *vdAddr

	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	enc := json.NewEncoder(outFile)
	// Only when vds is non-empty
	if len(vds) > 0 {
		name := vds[0].Name
		versionMap := make(map[string]OutputVersion, len(vds))

		for _, vd := range vds {
			number := vd.Version
			timestamp := vd.VersionCreated
			deps := vd.Dependencies

			depMap := make(map[string]string, len(deps))

			for _, ver := range deps {
				depMap[ver.Name] = ver.RequiredVersion
			}

			outVersion := OutputVersion{timestamp, depMap}
			versionMap[number] = outVersion
		}

		out := OutputFormat{name, versionMap}

		// Error handling for encoding
		if err := enc.Encode(out); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Wrote dependencies of %s to file \n", name)
	}
}

func MergeJSON(inPathTemplate string, amount int) {
	var wg sync.WaitGroup
	fmt.Println("Starting file merge process")
	var result []OutputFormat = make([]OutputFormat, 0, amount)
	outFile, err := os.OpenFile(fmt.Sprintf(inPathTemplate, "merged"), os.O_CREATE|os.O_WRONLY, 0644)
	enc := json.NewEncoder(outFile)

	if err != nil {
		log.Fatal(err)
	}

	resultChannel := make(chan OutputFormat, amount)
	for i := 0; i < amount; i++ {
		currentPath := fmt.Sprintf(inPathTemplate, fmt.Sprint(i))
		currentData, err := os.ReadFile(currentPath)

		// If the input file was empty, move on
		if len(currentData) < 1 {
			fmt.Printf("\tFile %d was empty\n", i)
			continue
		}

		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go func(currentDataAddr *[]byte, channel chan OutputFormat) {
			defer wg.Done()
			mergeJSONInternal(currentDataAddr, channel)
		}(&currentData, resultChannel)

		//os.Remove(fmt.Sprintf(inPathTemplate, fmt.Sprint(i)))
	}
	wg.Wait()
	close(resultChannel)
	for elem := range resultChannel {
		result = append(result, elem)
	}

	enc.Encode(result)
	fmt.Println("Merged JSON files")
}

func mergeJSONInternal(input *[]byte, channel chan OutputFormat) {
	var out OutputFormat
	if err := json.Unmarshal(*input, &out); err != nil {
		log.Fatal(err)
	}
	channel <- out
}

/** func testProcess() *[]VersionDependencies {
	var result []VersionDependencies
	name, number := "babel", "0.0.1"
	currentURL := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, number)

	rawDataAddr, statusAddr := request(currentURL)
	var parsed VersionInfo

	if err := json.Unmarshal(*rawDataAddr, &parsed); err != nil {
		status := *statusAddr
		fmt.Println("Http status code was: ", status) // This will probably be a rate-limit status code
		panic(err)
	}

	deps, devDeps := parsed.Dependencies, parsed.DevDependencies
	allDependencies := make([]Dependency, 0, len(deps)+len(devDeps))
	for k, v := range deps {
		allDependencies = append(allDependencies, Dependency{k, v.(string)})
	}
	for k, v := range devDeps {
		allDependencies = append(allDependencies, Dependency{k, v.(string)})
	}

	result = append(result, VersionDependencies{name, number, allDependencies})
	return &result
} **/
