package data

import (
	"crypto/tls"
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

// GetID gets the cluster ID
func GetID(id ClusterID) (string, error) {
	// First, check to see if we have an in-memory copy
	data := id.Cached()
	// If we haven't yet cached the ID in memory, invoke the passed-in getter
	if data == "" {
		d, err := id.Get()
		if err != nil {
			log.Print(err)
			return "", err
		}
		data = d
	}
	return data, nil
}

// GetAvailableVersions gets available component version data
func GetAvailableVersions(a AvailableVersions) ([]types.ComponentVersion, error) {
	// First, check to see if we have an in-memory copy
	data := a.Cached()
	// If we don't have any cached data, get the data from the remote authority
	if len(data) == 0 {
		d, err := a.Refresh()
		if err != nil {
			return nil, err
		}
		data = d
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
