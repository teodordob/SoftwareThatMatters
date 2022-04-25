package ingest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var raw_data []byte

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

// func (v Version) String() string {
// 	return fmt.Sprintf("{%v %v}", v.Number, time.Time(v.PublishedAt.t))
// }

type Version struct {
	Number      string      `json:number`
	PublishedAt CreatedTime `json:published_at`
}

type DiscoveryResponse struct {
	Name     string    `json:name`
	Versions []Version `json:versions`
}

func Ingest(query string) []DiscoveryResponse {
	raw_data = request(query)
	var arr []DiscoveryResponse
	if err := json.Unmarshal(raw_data, &arr); err != nil {
		panic(err)
	}
	return arr
	// fmt.Println(arr)
}

func request(req string) []byte {
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
	return body
}
