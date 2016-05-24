package data

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/pkg/swagger/client/operations"
	"github.com/deis/workflow-manager/pkg/swagger/models"
)

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
	vsns := availableVersionsFromAPI{
		rwm:             new(sync.RWMutex),
		baseVersionsURL: ts.URL,
		apiClient:       config.GetSwaggerClient(strings.TrimPrefix(ts.URL, "http://")),
	}
	retCompVsns, err := vsns.Refresh(models.Cluster{})
	assert.NoErr(t, err)
	assert.Equal(t, len(retCompVsns), len(expectedCompVsns.Data), "number of component versions")
}
