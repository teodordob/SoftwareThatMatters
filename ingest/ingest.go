package ingest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type CreatedTime time.Time

func (ct *CreatedTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*ct = CreatedTime(t)
	return nil
}

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
	Number      string      `json:number`
	PublishedAt CreatedTime `json:published_at`
}

type PackageInfo struct {
	Name     string    `json:name`
	Versions []Version `json:versions`
}

type VersionDependencies struct {
	Name         string
	Version      string
	Dependencies []Dependency
}

type VersionInfo struct {
	Dependencies    map[string]interface{} `json:dependencies`
	DevDependencies map[string]interface{} `json:devDependencies`
}

type Dependency struct {
	Name            string `json:name`
	RequiredVersion string `json:requirements`
}

func Ingest(query string) *[]VersionDependencies {
	rawDataAddr, statusAddr := request(query)
	var arr []PackageInfo
	if err := json.Unmarshal(*rawDataAddr, &arr); err != nil {
		status := *statusAddr
		fmt.Println("Uh-oh, HTTP status was: ", status) // This will probably be a rate-limit status code
		panic(err)
	}
	fmt.Println("Got data from input query")
	fmt.Println("Processing...")
	//return &arr
	return process(arr)
	// fmt.Println(arr)
}

func request(req string) (*[]byte, *string) {
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
	return &body, &resp.Status
}

// TODO: Gracefully handle 404s and rate limits
func process(input []PackageInfo) *[]VersionDependencies {
	var result []VersionDependencies

	for packageIdx, _ := range input {
		p := &input[packageIdx]
		name, versionsAddr := p.Name, &p.Versions
		for verIdx, _ := range *versionsAddr {
			number := (*versionsAddr)[verIdx].Number
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
			versionDeps := VersionDependencies{name, number, allDependencies}
			fmt.Println(versionDeps)
			result = append(result, versionDeps)
		}

	}
	return &result
}

func testProcess() *[]VersionDependencies {
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
}
