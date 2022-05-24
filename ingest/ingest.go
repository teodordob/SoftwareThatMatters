package ingest

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
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

type VersionDependencies struct {
	Name           string
	Version        string
	VersionCreated time.Time
	Dependencies   []Dependency
}

type Dependency struct {
	Name            string
	RequiredVersion string
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s:%s", d.Name, d.RequiredVersion)
}

func StreamParse(inPath string, jsonOutPathTemplate string) int {
	fmt.Println("Starting input JSON parser...")
	var wg sync.WaitGroup
	maxGoRoutines := 5000
	guard := make(chan int, maxGoRoutines)
	f, _ := os.Open(inPath)
	dec := json.NewDecoder(f)

	// Read opening bracket
	if _, err := dec.Token(); err != nil {
		log.Fatal(err)
	}
	// While the decoder says there is more to parse, parse a JSON entries and print them one-by-one
	i := 0
	for dec.More() {
		var e Entry

		if err := dec.Decode(&e); err != nil {
			log.Println(err)
			continue // Just move on and skip this entry
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
		}
		jsonPath := fmt.Sprintf(jsonOutPathTemplate, fmt.Sprint(i)) // Append a number to filePath
		wg.Add(1)                                                   // The waitGroup needs to wait for one more
		guard <- 0                                                  // Add one thread to the amount of running threads
		go func(vds *[]VersionDependencies, jsonPath string) {
			defer wg.Done() // Tell the WaitGroup this task is done after the function below is done
			writeToFileJSON(vds, jsonPath)
			<-guard // One thread was freed, now another can start
		}(&vds, jsonPath)

		i++
	}
	// Read closing bracket
	if _, err := dec.Token(); err != nil {
		log.Fatal(err)
	}
	wg.Wait() // Wait for all subroutines to be done
	close(guard)
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
	const comma string = ","
	fmt.Println("Starting file merge process")
	outFile, err := os.OpenFile(fmt.Sprintf(inPathTemplate, "merged"), os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()

	outFile.WriteString("[\n")

	for i := 0; i < amount; i++ {
		currentPath := fmt.Sprintf(inPathTemplate, fmt.Sprint(i))
		currentData, err := os.ReadFile(currentPath)

		if err != nil {
			log.Println(err) // Just skip if we couldn't find the file
			continue
		}

		// If the input file was empty, move on
		if len(currentData) < 1 {
			fmt.Printf("\tFile %d didn't contain versions object\n", i)
			os.Remove(currentPath)
			continue
		}

		finalData := currentData

		if i < amount-1 { // Append a comma after entry if we're not on the last entry
			finalData = append(finalData, comma...)
		}

		if _, err := outFile.Write(finalData); err != nil {
			log.Fatal("Couldn't write sequence to file")
		} else {
			os.Remove(currentPath) // Remove successfully merged file
		}
	}

	outFile.WriteString("]\n")
	fmt.Println("Merged JSON files")
}
