package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/types"
	"github.com/gorilla/mux"
)

const mockInstalledComponentName = "component"
const mockInstalledComponentDescription = "mock component"
const mockInstalledComponentVersion = "1.2.3"

// Creating a novel mock struct that fulfills the components.InstalledData interface
type mockInstalledComponents struct{}

func (g mockInstalledComponents) Get() ([]byte, error) {
	return []byte(fmt.Sprintf(`{
	  "components": [
	    {
	      "component": {
	        "name": "%s",
	        "description": "%s"
	      },
	      "version": {
	        "version": "%s"
	      }
	    }
	  ]
	}`, mockInstalledComponentName, mockInstalledComponentDescription, mockInstalledComponentVersion)), nil
}

const mockID = "faa31f63-d8dc-42e3-9568-405d20a3f755"

// Creating a novel mock struct that fulfills the data.ClusterID interface
type mockClusterID struct{}

func (c mockClusterID) Get() (string, error) {
	return mockID, nil
}

const mockAvailableComponentName = "component"
const mockAvailableComponentDescription = "mock component"
const mockAvailableComponentVersion = "3.2.1"

// Creating a novel mock struct that fulfills the components.AvailableComponentVersion interface
type mockAvailableVersion struct{}

func (c mockAvailableVersion) Get(component string) (types.Version, error) {
	if component == "component" {
		return types.Version{Version: "v2-beta"}, nil
	}
	return types.Version{}, fmt.Errorf("mock getter only accepts 'component' arg")
}

type genericJSON struct {
	Foo string `json:"foo"`
}

func TestComponentsHandler(t *testing.T) {
	componentsHandler := ComponentsHandler(mockInstalledComponents{}, mockClusterID{}, mockAvailableVersion{})
	resp, err := getTestHandlerResponse(componentsHandler)
	assert.NoErr(t, err)
	assert200(t, resp)
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	cluster, err := components.ParseJSONCluster(data)
	assert.NoErr(t, err)
	assert.Equal(t, cluster.ID, mockID, "ID value")
	assert.Equal(t, cluster.Components[0].Component.Name, mockInstalledComponentName, "Name value")
	assert.Equal(t, cluster.Components[0].Component.Description, mockInstalledComponentDescription, "Description value")
	assert.Equal(t, cluster.Components[0].Version.Version, mockInstalledComponentVersion, "Version value")
	//TODO
	//assert.Equal(t, cluster.Components[0].UpdateAvailable, mockAvailableComponentVersion, "available Version value")
}

func TestIDHandler(t *testing.T) {
	idHandler := IDHandler(mockClusterID{})
	resp, err := getTestHandlerResponse(idHandler)
	assert.NoErr(t, err)
	assert200(t, resp)
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	assert.Equal(t, string(data), mockID, "ID value")
}

func TestWritePlainText(t *testing.T) {
	const text = "foo"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writePlainText(text, w)
	})
	resp, err := getTestHandlerResponse(handler)
	assert.NoErr(t, err)
	assert.Equal(t, resp.Header.Get("Content-Type"), "text/plain", "Content-Type value")
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	assert.Equal(t, string(data), text, "text response")
}

func getTestHandlerResponse(handler http.Handler) (*http.Response, error) {
	r := mux.NewRouter()
	r.Handle("/", handler)
	server := httptest.NewServer(r)
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func assert200(t *testing.T, resp *http.Response) {
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
