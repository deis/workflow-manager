package data

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/deis/workflow-manager/types"
	"github.com/satori/go.uuid"
)

const (
	deisNamespace      = "deis"
	wfmSecretName      = "deis-workflow-manager"
	clusterIDSecretKey = "cluster-id"
)

var (
	clusterID         string
	availableVersions []types.ComponentVersion
)

// getClusterID returns the cached ID, or an error if it's not cached in memory
func getClusterID() (string, error) {
	if clusterID == "" {
		return "", fmt.Errorf("cluster ID not cached in memory")
	}
	return clusterID, nil
}

// getAvailableVersions returns the cached available version data, or an error
func getAvailableVersions() ([]types.ComponentVersion, error) {
	if len(availableVersions) == 0 {
		return nil, fmt.Errorf("no available versions data cached")
	}
	return availableVersions, nil
}

// GetID gets the cluster ID
func GetID(id ClusterID) (string, error) {
	// First, check to see if we have an in-memory copy
	data, err := getClusterID()
	// If we haven't yet cached the ID in memory, invoke the passed-in getter
	if err != nil {
		data, err = id.Get()
		if err != nil {
			log.Print(err)
			return "", err
		}
		clusterID = data
	}
	return data, nil
}

// GetAvailableVersions gets available component version data
func GetAvailableVersions(a AvailableVersions) ([]types.ComponentVersion, error) {
	// First, check to see if we have an in-memory copy
	data, err := getAvailableVersions()
	// If we don't have any cached data, get the data from the remote authority
	if err != nil {
		data, err = a.Refresh()
		if err != nil {
			log.Print(err)
			return nil, err
		}
		return data, nil
	}
	return data, nil
}

// getNewID returns a new Cluster ID string value
func getNewID() string {
	return uuid.NewV4().String()
}

// getTLSClient returns a TLS-enabled http.Client
func getTLSClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	return &http.Client{Transport: tr}
}
