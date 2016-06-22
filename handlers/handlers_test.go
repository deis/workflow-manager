package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/data"
	"github.com/deis/workflow-manager/pkg/swagger/models"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

const mockInstalledComponentName = "component"
const mockInstalledComponentDescription = "mock component"
const mockInstalledComponentVersion = "1.2.3"

// Creating a novel mock struct that fulfills the data.InstalledData interface
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
type mockClusterID struct {
	cached string
}

func (c mockClusterID) Get() (string, error) {
	return mockID, nil
}

func (c mockClusterID) Cached() string {
	return c.cached
}

func (c *mockClusterID) StoreInCache(cid string) {
	c.cached = cid
}

const mockAvailableComponentName = "component"
const mockAvailableComponentDescription = "mock component"
const mockAvailableComponentVersion = "3.2.1"

// Creating a novel mock struct that fulfills the data.AvailableComponentVersion interface
type mockAvailableVersion struct{}

func (c mockAvailableVersion) Get(component string, cluster models.Cluster) (models.Version, error) {
	if component == "component" {
		return models.Version{Version: "v2-beta"}, nil
	}
	return models.Version{}, fmt.Errorf("mock getter only accepts 'component' arg")
}

type genericJSON struct {
	Foo string `json:"foo"`
}

func TestComponentsHandler(t *testing.T) {
	componentsHandler := ComponentsHandler(
		mockInstalledComponents{},
		&mockClusterID{},
		mockAvailableVersion{},
		data.NewFakeKubeSecretGetterCreator(nil, nil),
	)
	resp, err := getTestHandlerResponse(componentsHandler)
	assert.NoErr(t, err)
	assert200(t, resp)
	respData, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	cluster, err := data.ParseJSONCluster(respData)
	assert.NoErr(t, err)
	assert.Equal(t, cluster.ID, mockID, "ID value")
	assert.Equal(t, cluster.Components[0].Component.Name, mockInstalledComponentName, "Name value")
	assert.Equal(t, *cluster.Components[0].Component.Description, mockInstalledComponentDescription, "Description value")
	assert.Equal(t, cluster.Components[0].Version.Version, mockInstalledComponentVersion, "Version value")
	//TODO
	//assert.Equal(t, cluster.Components[0].UpdateAvailable, mockAvailableComponentVersion, "available Version value")
}

func TestDoctorHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode([]byte("")); err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
			return
		}
	}))
	defer ts.Close()
	apiClient, err := config.GetSwaggerClient(ts.URL)
	doctorHandler := DoctorHandler(
		mockInstalledComponents{},
		&mockClusterID{},
		mockAvailableVersion{},
		data.NewFakeKubeSecretGetterCreator(nil, nil),
		apiClient,
	)
	resp, err := getTestHandlerResponse(doctorHandler)
	assert.NoErr(t, err)
	assert200(t, resp)
	respData, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	// verify that the return data is a uuid string
	doctorID1, err := uuid.FromString(string(respData[:]))
	assert.NoErr(t, err)
	// invoke the handler a 2nd time to ensure that unique IDs are created for
	// each request
	resp2, err := getTestHandlerResponse(doctorHandler)
	assert.NoErr(t, err)
	assert200(t, resp2)
	respData2, err := ioutil.ReadAll(resp2.Body)
	assert.NoErr(t, err)
	doctorID2, err := uuid.FromString(string(respData2[:]))
	assert.NoErr(t, err)
	if doctorID1 == doctorID2 {
		t.Error("DoctorHandler should return a unique ID for every invocation")
	}
}

func TestIDHandler(t *testing.T) {
	idHandler := IDHandler(&mockClusterID{})
	resp, err := getTestHandlerResponse(idHandler)
	assert.NoErr(t, err)
	assert200(t, resp)
	respData, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	assert.Equal(t, string(respData), mockID, "ID value")
}

func TestWritePlainText(t *testing.T) {
	const text = "foo"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writePlainText(text, w)
	})
	resp, err := getTestHandlerResponse(handler)
	assert.NoErr(t, err)
	assert.Equal(t, resp.Header.Get("Content-Type"), "text/plain", "Content-Type value")
	respData, err := ioutil.ReadAll(resp.Body)
	assert.NoErr(t, err)
	assert.Equal(t, string(respData), text, "text response")
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
