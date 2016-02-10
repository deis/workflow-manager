package components

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/mocks"
	"github.com/deis/workflow-manager/types"
)

const mockComponentName = "component"
const mockComponentDescription = "mock component"
const mockComponentVersion = "1.2.3"

// Creating a novel mock struct that fulfills the InstalledData interface
type mockInstalledComponents struct {
}

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
	}`, mockComponentName, mockComponentDescription, mockComponentVersion)), nil
}

// Creating a novel mock struct that fulfills the InstalledComponentVersion interface
type mockInstalledComponent struct{}

// Get method for InstalledComponent
func (c mockInstalledComponent) Get(component string) ([]byte, error) {
	if component == mockComponentName {
		return []byte(fmt.Sprintf(`{
		  "component": {
		    "name": "%s",
		    "description": "%s"
		  },
		  "version": {
		    "version": "%s"
		  }
		}`, mockComponentName, mockComponentDescription, mockComponentVersion)), nil
	}
	return []byte{}, fmt.Errorf("mock getter only accepts %s arg", mockComponentName)
}

func TestGetCluster(t *testing.T) {
	mockCluster := getMockCluster(t)
	cluster, err := GetCluster(mocks.InstalledMockData{}, mocks.ClusterIDMockData{}, mocks.LatestMockData{})
	assert.NoErr(t, err)
	assert.Equal(t, cluster.ID, mockCluster.ID, "ID value")
	for i, component := range cluster.Components {
		// verify that GetCluster returned the expected mock Component data
		assert.Equal(t, component.Component, mockCluster.Components[i].Component, "Component value")
		// verify that GetCluster returned the expected mock Version data
		assert.Equal(t, component.Version, mockCluster.Components[i].Version, "Version property")
		// TODO
		/*latestComponent := getMockLatest(component.Component.Name, t)
		// verify that the expected "UpdateAvailable" field values were added
		assertUpdateAvailableAdded(latestComponent, cluster.Components[i], t)*/
	}
}

func TestAddUpdateData(t *testing.T) {
	mockCluster := getMockCluster(t)
	// AddUpdateData should add an "UpdateAvailable" field to any components whose versions are out-of-date
	err := AddUpdateData(&mockCluster, mocks.LatestMockData{})
	assert.NoErr(t, err)
	//TODO when newestVersion is implemented, actually test for the addition of "UpdateAvailable" fields
}

func TestGetInstalled(t *testing.T) {
	cluster, err := GetInstalled(mockInstalledComponents{})
	assert.NoErr(t, err)
	assert.Equal(t, cluster.Components[0].Component.Name, mockComponentName, "Name value")
	assert.Equal(t, cluster.Components[0].Component.Description, mockComponentDescription, "Description value")
	assert.Equal(t, cluster.Components[0].Version.Version, mockComponentVersion, "Version value")
}

func TestParseJSONCluster(t *testing.T) {
	const id = "f91378a6-a815-4c20-9b0d-77b205cd3ba4"
	const name = "component"
	const description = "test component"
	const version = "1.0.0"
	raw := []byte(fmt.Sprintf(`{
	  "id": "%s",
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
	}`, id, name, description, version))
	cluster, err := ParseJSONCluster(raw)
	assert.NoErr(t, err)

	assert.Equal(t, cluster.ID, id, "ID value")
	assert.Equal(t, cluster.Components[0].Component.Name, name, "Name value")
	assert.Equal(t, cluster.Components[0].Component.Description, description, "Description value")
	assert.Equal(t, cluster.Components[0].Version.Version, version, "Version value")
}

func TestNewestSemVer(t *testing.T) {
	// Verify that NewestSemVer returns correct semver string for larger major, minor, and patch substrings
	const v1Lower = "2.0.0"
	v2s := [3]string{"3.0.0", "2.1.0", "2.0.1"}
	for _, v2 := range v2s {
		newest, err := NewestSemVer(v1Lower, v2)
		assert.NoErr(t, err)
		if newest != v2 {
			fmt.Printf("expected %s to be greater than %s\n", v2, v1Lower)
			t.Fatal("semver comparison failure")
		}
	}
	// Verify that NewestSemVer returns correct semver string for smaller major, minor, and patch substrings
	const v1Higher = "2.4.5"
	v2s = [3]string{"1.99.23", "2.3.99", "2.4.4"}
	for _, v2 := range v2s {
		newest, err := NewestSemVer(v1Higher, v2)
		assert.NoErr(t, err)
		if newest != v1Higher {
			fmt.Printf("expected %s to be greater than %s\n", v1Higher, v2)
			t.Fatal("semver comparison failure")
		}
	}
	// Verify that NewestSemVer returns correct semver string for comparing equal strings
	const v1Equal = "1.0.0"
	v2 := v1Equal
	newest, err := NewestSemVer(v1Equal, v2)
	assert.NoErr(t, err)
	if newest != v1Equal && newest != v2 {
		fmt.Printf("expected %s to be equal to %s and %s\n", newest, v1Equal, v2)
		t.Fatal("semver comparison failure")
	}
}

func getMockCluster(t *testing.T) types.Cluster {
	mockData, err := mocks.GetMockCluster()
	assert.NoErr(t, err)
	mockCluster, err := ParseJSONCluster(mockData)
	assert.NoErr(t, err)
	return mockCluster
}

func getMockLatest(name string, t *testing.T) types.Version {
	version, err := mocks.GetMockLatest(name)
	assert.NoErr(t, err)
	return version
}

func assertUpdateAvailableAdded(latestComponent types.ComponentVersion, component types.ComponentVersion, t *testing.T) {
	if latestComponent.Version.Version != component.Version.Version {
		if component.UpdateAvailable != latestComponent.Version.Version {
			t.Fatal("failed to get back UpdateAvailable property for a component that has a newer version")
		}
	} else if component.UpdateAvailable != "" {
		t.Fatal("expected a nil string type for UpdateAvailable property at component ", component.Component.Name)
	}
}
