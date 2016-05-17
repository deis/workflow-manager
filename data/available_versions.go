package data

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/types"
)

// AvailableVersions is an interface for managing available component version data
type AvailableVersions interface {
	// Cached returns the internal cache of component versions. Returns the empty slice on a miss
	Cached() []types.ComponentVersion
	// will have a Refresh method to retrieve the version data from the remote authority
	Refresh() ([]types.ComponentVersion, error)
	// will have a Store method to cache the version data in memory
	Store([]types.ComponentVersion)
}

// AvailableVersionsFromAPI fulfills the AvailableVersions interface
type AvailableVersionsFromAPI struct {
	cache []types.ComponentVersion
	rwm   *sync.RWMutex
}

// NewAvailableVersionsFromAPI returns a new AvailableVersions implementation that uses the workflow manager API to get its version information
func NewAvailableVersionsFromAPI() *AvailableVersionsFromAPI {
	return &AvailableVersionsFromAPI{rwm: new(sync.RWMutex), cache: nil}
}

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

// Cached is the AvailableVersions interface implementation
func (a AvailableVersionsFromAPI) Cached() []types.ComponentVersion {
	a.rwm.RLock()
	defer a.rwm.RUnlock()
	return a.cache
}

// Store is the AvailableVersions interface implementation
func (a *AvailableVersionsFromAPI) Store(c []types.ComponentVersion) {
	a.rwm.Lock()
	defer a.rwm.Unlock()
	a.cache = c
}
