package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/pkg/swagger/client/operations"
	"github.com/deis/workflow-manager/pkg/swagger/models"
)

// Creating a novel mock struct that fulfills the AvailableVersions interface
type testAvailableVersions struct{}

func (a testAvailableVersions) Refresh(cluster models.Cluster) ([]models.ComponentVersion, error) {
	data := getMockComponentVersions()
	var componentVersions []models.ComponentVersion
	if err := json.Unmarshal(data, &componentVersions); err != nil {
		return nil, err
	}
	return componentVersions, nil
}

func (a testAvailableVersions) Store(c []models.ComponentVersion) {
	return
}

func (a testAvailableVersions) Cached() []models.ComponentVersion {
	return nil
}

// Creating another mock struct that fulfills the AvailableVersions interface
type shouldBypassAvailableVersions struct{}

func (a shouldBypassAvailableVersions) Refresh(cluster models.Cluster) ([]models.ComponentVersion, error) {
	var componentVersions []models.ComponentVersion
	data := []byte(fmt.Sprintf(`[{
	  "components": [
	    {
	      "component": {
	        "name": "bypass me",
	        "description": "bypass me"
	      },
	      "version": {
	        "version": "v2-bypass"
	      }
	    }
	  ]
	}]`))
	if err := json.Unmarshal(data, &componentVersions); err != nil {
		return nil, err
	}
	return componentVersions, nil
}

func (a shouldBypassAvailableVersions) Store(c []models.ComponentVersion) {
	return
}

func (a shouldBypassAvailableVersions) Cached() []models.ComponentVersion {
	return nil
}

// Calls GetAvailableVersions twice, the first time we expect our passed-in struct w/ Refresh() method
// to be invoked, the 2nd time we expect to receive the same value back (cached in memory)
// and for the passed-in Refresh() method to be ignored
func TestGetAvailableVersions(t *testing.T) {
	mock := getMockComponentVersions()
	var mockVersions []models.ComponentVersion
	assert.NoErr(t, json.Unmarshal(mock, &mockVersions))
	versions, err := GetAvailableVersions(testAvailableVersions{}, models.Cluster{})
	assert.NoErr(t, err)
	assert.Equal(t, versions, mockVersions, "component versions data")
	versions, err = GetAvailableVersions(shouldBypassAvailableVersions{}, models.Cluster{})
	assert.NoErr(t, err)
	assert.Equal(t, versions, mockVersions, "component versions data")
}

func TestRefreshAvailableVersions(t *testing.T) {
	desc := "this is test1"
	updAvail := "nothing"
	expectedCompVsns := operations.GetComponentsByLatestReleaseOKBodyBody{
		Data: []*models.ComponentVersion{
			&models.ComponentVersion{
				Component:       &models.Component{Name: "test1", Description: &desc},
				Version:         &models.Version{Train: "testTrain"},
				UpdateAvailable: &updAvail,
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(expectedCompVsns); err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
			return
		}
	}))
	defer ts.Close()
	apiclient, err := config.GetSwaggerClient(ts.URL)
	assert.NoErr(t, err)
	vsns := availableVersionsFromAPI{
		rwm:             new(sync.RWMutex),
		baseVersionsURL: ts.URL,
		apiClient:       apiclient,
	}
	retCompVsns, err := vsns.Refresh(models.Cluster{})
	assert.NoErr(t, err)
	assert.Equal(t, len(retCompVsns), len(expectedCompVsns.Data), "number of component versions")
}
