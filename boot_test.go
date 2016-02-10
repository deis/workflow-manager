package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/handlers"
	"github.com/deis/workflow-manager/mocks"
	"github.com/gorilla/mux"
)

func newServer() *httptest.Server {
	r := mux.NewRouter()
	r.Handle("/components", handlers.ComponentsHandler(mocks.InstalledMockData{}, mocks.ClusterIDMockData{}, mocks.LatestMockData{}))
	r.Handle("/id", handlers.IDHandler(mocks.ClusterIDMockData{}))
	return httptest.NewServer(r)
}

func TestGetComponents(t *testing.T) {
	const componentRoute = "/components"
	resp, err := testGet(componentRoute)
	assert.NoErr(t, err)
	assert200(t, resp)
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	cluster, err := components.ParseJSONCluster(data)
	assert.NoErr(t, err)
	mockData, err := mocks.GetMockCluster()
	assert.NoErr(t, err)
	mockCluster, err := components.ParseJSONCluster(mockData)
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

func TestGetID(t *testing.T) {
	const idRoute = "/id"
	resp, err := testGet(idRoute)
	assert.NoErr(t, err)
	assert200(t, resp)
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	mockData, err := mocks.GetMockClusterID()
	assert.NoErr(t, err)
	assert.Equal(t, string(data), mockData, "id data response")
}

func testGet(route string) (*http.Response, error) {
	server := newServer()
	defer server.Close()
	resp, err := httpGet(server, route)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func httpGet(s *httptest.Server, route string) (*http.Response, error) {
	return http.Get(s.URL + route)
}

func assert200(t *testing.T, resp *http.Response) {
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
