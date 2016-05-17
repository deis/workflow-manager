package data

import (
	"encoding/json"
	"net/http"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/types"
)

// AvailableVersions is an interface for managing available component version data
type AvailableVersions interface {
	// will have a Refresh method to retrieve the version data from the remote authority
	Refresh() ([]types.ComponentVersion, error)
	// will have a Store method to cache the version data in memory
	Store([]types.ComponentVersion)
}

// AvailableVersionsFromAPI fulfills the AvailableVersions interface
type AvailableVersionsFromAPI struct{}

// Refresh method for AvailableVersionsFromAPI
func (a AvailableVersionsFromAPI) Refresh() ([]types.ComponentVersion, error) {
	var versionsRoute = "/" + config.Spec.APIVersion + "/versions"
	url := config.Spec.VersionsAPIURL + versionsRoute
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []types.ComponentVersion{}, err
	}

	resp, err := getTLSClient().Do(req)
	if err != nil {
		return []types.ComponentVersion{}, err
	}
	defer resp.Body.Close()
	var availableVersions []types.ComponentVersion
	json.NewDecoder(resp.Body).Decode(&availableVersions)
	a.Store(availableVersions)
	return availableVersions, nil
}

// Store method for AvailableVersionsFromAPI
func (a AvailableVersionsFromAPI) Store(c []types.ComponentVersion) {
	availableVersions = c
}
