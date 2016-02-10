package mocks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/deis/workflow-manager/types"
)

const mainPackage = "workflow-manager"

func getMocksWd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	cwdSplit := strings.Split(cwd, "/")
	if (cwdSplit[len(cwdSplit)-1]) != mainPackage {
		cwdSplit = cwdSplit[:len(cwdSplit)-1] // strip last directory
		return strings.Join(cwdSplit, "/") + "/mocks/"
	}
	return cwd + "/mocks/"
}

// InstalledMockData mock data struct
type InstalledMockData struct{}

// Get method for InstalledMockData
func (g InstalledMockData) Get() ([]byte, error) {
	data, err := GetMockComponents()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return data, nil
}

// ClusterIDMockData mock data struct
type ClusterIDMockData struct{}

// Get method for ClusterIDMockData
func (c ClusterIDMockData) Get() (string, error) {
	data, err := GetMockClusterID()
	if err != nil {
		log.Print(err)
		return "", err
	}
	return data, nil
}

// LatestMockData mock data struct
type LatestMockData struct{}

// Get method for LatestMockData
func (c LatestMockData) Get(component string) (types.Version, error) {
	data, err := GetMockLatest(component)
	if err != nil {
		log.Print(err)
		return types.Version{}, err
	}
	return data, nil
}

// GetMockCluster returns a mock JSON cluster response
func GetMockCluster() ([]byte, error) {
	data, err := getJSON(getMocksWd() + "cluster.json")
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetMockComponents returns a mock JSON cluster response
func GetMockComponents() ([]byte, error) {
	data, err := getJSON(getMocksWd() + "components.json")
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetMockClusterPost returns mock JSON cluster POST data
func GetMockClusterPost() ([]byte, error) {
	data, err := getJSON(getMocksWd() + "cluster-post.json")
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetMockClusterID returns a mock JSON cluster response
func GetMockClusterID() (string, error) {
	data, err := getText(getMocksWd() + "id.txt")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetMockLatest returns a mock "latest component version" response
func GetMockLatest(c string) (types.Version, error) {
	data, err := getText(getMocksWd() + "latest-component-version-" + c + ".txt")
	if err != nil {
		return types.Version{}, err
	}
	return types.Version{Version: data}, nil
}

// GetMockComponentV2Beta returns a mock "latest component version" response, for v2-beta
func GetMockComponentV2Beta() ([]byte, error) {
	data, err := getJSON(getMocksWd() + "latest-component-version-v2-beta.json")
	if err != nil {
		return nil, err
	}
	return data, nil
}

// getJSON gets a JSON file from the local filesystem
func getJSON(filepath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("Error reading .json file: %#v\n", err)
		return nil, err
	}
	if !isJSON(data) {
		return nil, fmt.Errorf("data is not valid JSON")
	}
	return data, nil
}

// getText gets a text file from the local filesystem
func getText(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("Error reading .txt file: %#v\n", err)
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// isJSON checks for valid JSON
func isJSON(b []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(b, &js) == nil
}
