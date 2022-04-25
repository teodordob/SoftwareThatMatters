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
	Name         string `json:name`
	Version      string `json:dependencies_for_version`
	Dependencies []Dependency
}

type Dependency struct {
	Name            string `json:name`
	RequiredVersion string `json:requirements`
}

func Ingest(query string) *[]PackageInfo {
	rawData := *request(query)
	var arr []PackageInfo
	if err := json.Unmarshal(rawData, &arr); err != nil {
		panic(err)
	}
	return &arr
	// fmt.Println(arr)
}

func request(req string) *[]byte {
	fmt.Println("Starting request...")
	resp, err := http.Get(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
	defer fmt.Println("Processing...")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	return &body
}

func process(inputAddr *[]PackageInfo, platform string) *[]VersionDependencies {
	var result []VersionDependencies

	for i, p := range *inputAddr {
		i += 1
		name, versionsAddr := p.Name, &p.Versions
		for j, ver := range *versionsAddr {
			j += 1
			number := ver.Number
			currentURL := fmt.Sprintf("https://libraries.io/api/%s/%s/%s/dependencies?api_key=3dc75447d3681ffc2d17517265765d23", platform, name, number)

			rawData := *request(currentURL)
			var parsed VersionDependencies

			if err := json.Unmarshal(rawData, &parsed); err != nil {
				panic(err)
			}
			result[i*j] = parsed
		}

	}
	return &result
}
