package data

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/types"
)

const mockClusterID = "faa31f63-d8dc-42e3-9568-405d20a3f755"
const mockComponentName = "component"
const mockComponentDescription = "mock component"
const mockComponentVersion = "v2-beta"

// Creating a novel mock struct that fulfills the ClusterID interface
type testClusterID struct {
	cache string
}

func (c testClusterID) Get() (string, error) {
	return mockClusterID, nil
}

func (c testClusterID) Cached() string {
	return c.cache
}

func (c *testClusterID) StoreInCache(cid string) {
	c.cache = cid
}

// Creating a novel mock struct that fulfills the AvailableVersions interface
type testAvailableVersions struct{}

func (a testAvailableVersions) Refresh() ([]types.ComponentVersion, error) {
	data := getMockComponentVersions()
	var componentVersions []types.ComponentVersion
	_ = json.Unmarshal(data, &componentVersions)
	return componentVersions, nil
}

func (a testAvailableVersions) Store(c []types.ComponentVersion) {
	return
}

func (a testAvailableVersions) Cached() []types.ComponentVersion {
	return nil
}

// Creating another mock struct that fulfills the AvailableVersions interface
type shouldBypassAvailableVersions struct{}

func (a shouldBypassAvailableVersions) Refresh() ([]types.ComponentVersion, error) {
	var componentVersions []types.ComponentVersion
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
	_ = json.Unmarshal(data, &componentVersions)
	return componentVersions, nil
}

func (a shouldBypassAvailableVersions) Store(c []types.ComponentVersion) {
	return
}

func (a shouldBypassAvailableVersions) Cached() []types.ComponentVersion {
	return nil
}

// Calls GetID twice, the first time we expect our passed-in struct w/ Get() method
// to be invoked, the 2nd time we expect to receive the same value back (cached in memory)
// and for the passed-in Get() method to be ignored
func TestGetID(t *testing.T) {
	cid := &testClusterID{}
	id, err := GetID(cid)
	assert.NoErr(t, err)
	assert.Equal(t, id, mockClusterID, "cluster ID value")
	cid.cache = "something else"
	id, err = GetID(cid)
	assert.NoErr(t, err)
	assert.Equal(t, id, "something else", "cluster ID value")
}

// Calls GetAvailableVersions twice, the first time we expect our passed-in struct w/ Refresh() method
// to be invoked, the 2nd time we expect to receive the same value back (cached in memory)
// and for the passed-in Refresh() method to be ignored
func TestGetAvailableVersions(t *testing.T) {
	mock := getMockComponentVersions()
	var mockVersions []types.ComponentVersion
	_ = json.Unmarshal(mock, &mockVersions)
	versions, err := GetAvailableVersions(testAvailableVersions{})
	assert.NoErr(t, err)
	assert.Equal(t, versions, mockVersions, "component versions data")
	versions, err = GetAvailableVersions(shouldBypassAvailableVersions{})
	assert.NoErr(t, err)
	assert.Equal(t, versions, mockVersions, "component versions data")
}

func getMockComponentVersions() []byte {
	return []byte(fmt.Sprintf(`[{
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
	}]`, mockComponentName, mockComponentDescription, mockComponentVersion))
}
