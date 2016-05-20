package data

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/rest"
	"github.com/deis/workflow-manager/types"
)

// AvailableVersions is an interface for managing available component version data
type AvailableVersions interface {
	// Cached returns the internal cache of component versions. Returns the empty slice on a miss
	Cached() []types.ComponentVersion
	// Refresh gets the latest versions of each component listed in the given cluster
	Refresh(types.Cluster) ([]types.ComponentVersion, error)
	// Store stores the given slice of types.ComponentVersion in internal storage
	Store([]types.ComponentVersion)
}

type availableVersionsFromAPI struct {
	cache           []types.ComponentVersion
	rwm             *sync.RWMutex
	baseVersionsURL string
	restClient      rest.Client
}

// NewAvailableVersionsFromAPI returns a new AvailableVersions implementation that fetches its version information from a workflow manager API. It uses baseVersionsURL as the server address. If that parameter is passed as the empty string, uses config.Spec.VersionsAPIURL
func NewAvailableVersionsFromAPI(
	restCl rest.Client,
	baseVersionsURL string,
	secretGetterCreator KubeSecretGetterCreator,
	rcLister rc.Lister,
) AvailableVersions {
	if baseVersionsURL == "" {
		baseVersionsURL = config.Spec.VersionsAPIURL
	}
	return &availableVersionsFromAPI{
		rwm:             new(sync.RWMutex),
		cache:           nil,
		baseVersionsURL: baseVersionsURL,
		restClient:      restCl,
	}
}

// Refresh method for AvailableVersionsFromAPI
func (a availableVersionsFromAPI) Refresh(cluster types.Cluster) ([]types.ComponentVersion, error) {
	reqBody := SparseComponentAndTrainInfoJSONWrapper{}
	for _, component := range cluster.Components {
		sparseComponentAndTrainInfo := SparseComponentAndTrainInfo{}
		sparseComponentAndTrainInfo.Component.Name = component.Component.Name
		sparseComponentAndTrainInfo.Version.Train = component.Version.Train
		reqBody.Data = append(reqBody.Data, sparseComponentAndTrainInfo)
	}
	js, err := json.Marshal(reqBody)
	if err != nil {
		return []types.ComponentVersion{}, err
	}

	resp, err := a.restClient.Do("POST", rest.JSContentTypeHeader, bytes.NewBuffer(js), config.Spec.APIVersion, "versions", "latest")
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
