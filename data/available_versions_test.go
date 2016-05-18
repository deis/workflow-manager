package data

import (
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/arschles/testsrv"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/types"
)

func TestRefreshAvailableVersions(t *testing.T) {
	expectedCompVsns := ComponentVersionsJSONWrapper{
		Data: []types.ComponentVersion{
			types.ComponentVersion{
				Component:       types.Component{Name: "test1", Description: "this is test1"},
				Version:         types.Version{Train: "testTrain"},
				UpdateAvailable: "nothing",
			},
		},
	}
	hdl := func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(expectedCompVsns); err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
			return
		}
	}
	srv := testsrv.StartServer(http.HandlerFunc(hdl))
	defer srv.Close()
	vsns := availableVersionsFromAPI{
		rwm:             new(sync.RWMutex),
		baseVersionsURL: srv.URLStr(),
		clusterGetter: func() (types.Cluster, error) {
			return types.Cluster{ID: "testCluster", Components: expectedCompVsns.Data}, nil
		},
	}
	retCompVsns, err := vsns.Refresh()
	assert.NoErr(t, err)
	assert.Equal(t, len(retCompVsns), len(expectedCompVsns.Data), "number of component versions")
	recv := srv.AcceptN(1, 1*time.Millisecond)
	assert.Equal(t, len(recv), 1, "number of requests to the fake server")
	req := recv[0].Request
	assert.Equal(t, req.Method, "POST", "request method")
	assert.Equal(t, req.URL.Path, "/"+config.Spec.APIVersion+"/versions/latest", "request path")
	// TODO: re-enable this code when https://github.com/arschles/testsrv/issues/1 is fixed.
	// see https://github.com/deis/workflow-manager/issues/48 for more
	// var reqBody SparseComponentAndTrainInfoJSONWrapper
	// assert.NoErr(t, json.NewDecoder(req.Body).Decode(&reqBody))
	// assert.Equal(t, len(reqBody.Data), len(expectedCompVsns.Data), "num components in request body")
}
