package data

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/types"
)

// AvailableVersions is an interface for managing available component version data
type AvailableVersions interface {
	// Cached returns the internal cache of component versions. Returns the empty slice on a miss
	Cached() []types.ComponentVersion
	// Refresh gets the latest versions of each component in the internally stored cluster from the versions API
	Refresh() ([]types.ComponentVersion, error)
	// Store stores the given slice of types.ComponentVersion in internal storage
	Store([]types.ComponentVersion)
}

type availableVersionsFromAPI struct {
	cache           []types.ComponentVersion
	rwm             *sync.RWMutex
	clusterGetter   func() (types.Cluster, error)
	baseVersionsURL string
}

// NewAvailableVersionsFromAPI returns a new AvailableVersions implementation that fetches its version information from a workflow manager API. It uses baseVersionsURL as the server address. If that parameter is passed as the empty string, uses config.Spec.VersionsAPIURL
func NewAvailableVersionsFromAPI(baseVersionsURL string) AvailableVersions {
	if baseVersionsURL == "" {
		baseVersionsURL = config.Spec.VersionsAPIURL
	}
	return &availableVersionsFromAPI{
		rwm:             new(sync.RWMutex),
		cache:           nil,
		baseVersionsURL: baseVersionsURL,
		clusterGetter: func() (types.Cluster, error) {
			installedData := InstalledDeisData{}
			clusterID := NewClusterIDFromPersistentStorage()
			compVsn := LatestReleasedComponent{}
			return GetCluster(installedData, clusterID, compVsn)
		},
	}
}

// Refresh method for AvailableVersionsFromAPI
func (a availableVersionsFromAPI) Refresh() ([]types.ComponentVersion, error) {
	cluster, err := a.clusterGetter()
	if err != nil {
		return []types.ComponentVersion{}, err
	}
	reqBody := SparseComponentAndTrainInfoJSONWrapper{}
	for _, component := range cluster.Components {
		sparseComponentAndTrainInfo := SparseComponentAndTrainInfo{}
		sparseComponentAndTrainInfo.Component.Name = component.Component.Name
		sparseComponentAndTrainInfo.Version.Train = component.Version.Train
		reqBody.Data = append(reqBody.Data, sparseComponentAndTrainInfo)
	}
	js, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("error making a JSON representation of cluster data")
		return []types.ComponentVersion{}, err
	}
	var versionsLatestRoute = "/" + config.Spec.APIVersion + "/versions/latest"
	url := a.baseVersionsURL + versionsLatestRoute
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if err != nil {
		return []types.ComponentVersion{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := getTLSClient().Do(req)
	if err != nil {
		return []types.ComponentVersion{}, err
	}
	defer resp.Body.Close()
	var ret ComponentVersionsJSONWrapper
	json.NewDecoder(resp.Body).Decode(&ret)
	a.Store(ret.Data)
	return ret.Data, nil
}

// Cached is the AvailableVersions interface implementation
func (a availableVersionsFromAPI) Cached() []types.ComponentVersion {
	a.rwm.RLock()
	defer a.rwm.RUnlock()
	return a.cache
}

// Store is the AvailableVersions interface implementation
func (a *availableVersionsFromAPI) Store(c []types.ComponentVersion) {
	a.rwm.Lock()
	defer a.rwm.Unlock()
	a.cache = c
}
