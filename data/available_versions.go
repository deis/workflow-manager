package data

import (
	"sync"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/config"
	apiclient "github.com/deis/workflow-manager/pkg/swagger/client"
	"github.com/deis/workflow-manager/pkg/swagger/client/operations"
	"github.com/deis/workflow-manager/pkg/swagger/models"
)

// AvailableVersions is an interface for managing available component version data
type AvailableVersions interface {
	// Cached returns the internal cache of component versions. Returns the empty slice on a miss
	Cached() []models.ComponentVersion
	// Refresh gets the latest versions of each component listed in the given cluster
	Refresh(models.Cluster) ([]models.ComponentVersion, error)
	// Store stores the given slice of models.ComponentVersion in internal storage
	Store([]models.ComponentVersion)
}

type availableVersionsFromAPI struct {
	cache           []models.ComponentVersion
	rwm             *sync.RWMutex
	baseVersionsURL string
	apiClient       *apiclient.WorkflowManager
}

// NewAvailableVersionsFromAPI returns a new AvailableVersions implementation that fetches its version information from a workflow manager API. It uses baseVersionsURL as the server address. If that parameter is passed as the empty string, uses config.Spec.VersionsAPIURL
func NewAvailableVersionsFromAPI(
	apiClient *apiclient.WorkflowManager,
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
		apiClient:       apiClient,
	}
}

// Refresh method for AvailableVersionsFromAPI
func (a availableVersionsFromAPI) Refresh(cluster models.Cluster) ([]models.ComponentVersion, error) {
	reqBody := operations.GetComponentsByLatestReleaseBody{}
	for _, component := range cluster.Components {
		cv := new(models.ComponentVersion)
		cv.Component = &models.Component{}
		cv.Version = &models.Version{}
		cv.Component.Name = component.Component.Name
		cv.Version.Train = "stable"
		reqBody.Data = append(reqBody.Data, cv)
	}

	resp, err := a.apiClient.Operations.GetComponentsByLatestRelease(&operations.GetComponentsByLatestReleaseParams{Body: reqBody})
	if err != nil {
		return []models.ComponentVersion{}, err
	}
	ret := []models.ComponentVersion{}
	for _, cv := range resp.Payload.Data {
		ret = append(ret, *cv)
	}
	a.Store(ret)
	return ret, nil
}

// Cached is the AvailableVersions interface implementation
func (a availableVersionsFromAPI) Cached() []models.ComponentVersion {
	a.rwm.RLock()
	defer a.rwm.RUnlock()
	return a.cache
}

// Store is the AvailableVersions interface implementation
func (a *availableVersionsFromAPI) Store(c []models.ComponentVersion) {
	a.rwm.Lock()
	defer a.rwm.Unlock()
	a.cache = c
}
