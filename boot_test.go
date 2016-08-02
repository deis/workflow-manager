package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/data"
	"github.com/deis/workflow-manager/handlers"
	"github.com/deis/workflow-manager/mocks"
	apiclient "github.com/deis/workflow-manager/pkg/swagger/client"
	"github.com/gorilla/mux"
)

func newServer(apiClient *apiclient.WorkflowManager) *httptest.Server {
	r := mux.NewRouter()
	compHdl := handlers.ComponentsHandler(
		mocks.InstalledMockData{},
		&mocks.ClusterIDMockData{},
		mocks.LatestMockData{},
	)
	r.Handle("/components", compHdl)
	idHdl := handlers.IDHandler(&mocks.ClusterIDMockData{})
	r.Handle("/id", idHdl)
	docHdl := handlers.DoctorHandler(
		mocks.InstalledMockData{},
		mocks.RunningK8sMockData{}, // TODO: mock k8s node data
		&mocks.ClusterIDMockData{},
		mocks.LatestMockData{},
		apiClient,
	)
	r.Handle("/doctor", docHdl).Methods("POST")
	return httptest.NewServer(r)
}

func TestGetComponents(t *testing.T) {
	const componentRoute = "/components"
	resp, apiServer, err := testGet(componentRoute)
	if apiServer != nil {
		apiServer.Close()
	}
	assert.NoErr(t, err)
	assert200(t, resp)
	respData, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	cluster, err := data.ParseJSONCluster(respData)
	assert.NoErr(t, err)
	mockData, err := mocks.GetMockCluster()
	assert.NoErr(t, err)
	mockCluster, err := data.ParseJSONCluster(mockData)
	assert.NoErr(t, err)
	assert.Equal(t, cluster.ID, mockCluster.ID, "cluster ID value")
	for i, component := range cluster.Components {
		assert.Equal(t, component.Component, mockCluster.Components[i].Component, "component type")
		assert.Equal(t, component.Version, mockCluster.Components[i].Version, "version type")
		_, err := mocks.GetMockLatest(component.Component.Name)
		assert.NoErr(t, err)
		// TODO add tests for UpdateAvailable field
	}
}

func TestPostDoctor(t *testing.T) {
	const doctorRoute = "/doctor"
	resp, apiServer, err := testPostNoBody(doctorRoute)
	if apiServer != nil {
		apiServer.Close()
	}
	assert.NoErr(t, err)
	assert200(t, resp)
}

func TestGetID(t *testing.T) {
	const idRoute = "/id"
	resp, apiServer, err := testGet(idRoute)
	if apiServer != nil {
		apiServer.Close()
	}
	assert.NoErr(t, err)
	assert200(t, resp)
	respData, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	mockData, err := mocks.GetMockClusterID()
	assert.NoErr(t, err)
	assert.Equal(t, string(respData), mockData, "id data response")
}

func testGet(route string) (*http.Response, *httptest.Server, error) {
	apiClient, apiServer, err := getWfmMockAPIClient([]byte(""))
	if err != nil {
		return nil, nil, err
	}
	server := newServer(apiClient)
	defer server.Close()
	resp, err := httpGet(server, route)
	if err != nil {
		return nil, nil, err
	}
	return resp, apiServer, nil
}

func testPostNoBody(route string) (*http.Response, *httptest.Server, error) {
	apiClient, apiServer, err := getWfmMockAPIClient([]byte(""))
	if err != nil {
		return nil, nil, err
	}
	server := newServer(apiClient)
	defer server.Close()
	resp, err := httpPost(server, route, "")
	if err != nil {
		return nil, nil, err
	}
	return resp, apiServer, nil
}

func getWfmMockAPIClient(respBody []byte) (*apiclient.WorkflowManager, *httptest.Server, error) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(respBody); err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
			return
		}
	}))
	apiClient, err := config.GetSwaggerClient(apiServer.URL)
	if err != nil {
		return nil, nil, err
	}
	return apiClient, apiServer, nil

}

func httpGet(s *httptest.Server, route string) (*http.Response, error) {
	return http.Get(s.URL + route)
}

func httpPost(s *httptest.Server, route string, json string) (*http.Response, error) {
	fullURL := s.URL + route
	const bodyType = "application/json"
	return http.Post(fullURL, bodyType, bytes.NewBuffer([]byte(json)))
}

func assert200(t *testing.T, resp *http.Response) {
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
