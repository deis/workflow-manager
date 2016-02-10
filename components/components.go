package components

// handler echoes the HTTP request.
import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/data"
	"github.com/deis/workflow-manager/types"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
)

var (
	memoVersions = make(map[string]types.Version) // cached in-memory store of latest deis component versions
)

// InstalledData is an interface for managing installed cluster metadata
type InstalledData interface {
	// will have a Get method to retrieve installed data
	Get() ([]byte, error)
}

// AvailableComponentVersion is an interface for managing component version data
type AvailableComponentVersion interface {
	// will have a Get method to retrieve available component version data
	Get(component string) (types.Version, error)
}

// InstalledDeisData fulfills the InstalledData interface
type InstalledDeisData struct{}

// Get method for InstalledDeisData
func (g InstalledDeisData) Get() ([]byte, error) {
	rcItems, err := getDeisRCItems()
	var cluster types.Cluster
	for _, rc := range rcItems {
		component := types.ComponentVersion{}
		component.Component.Name = rc.Name
		component.Component.Description = rc.Annotations["chart.helm.sh/description"]
		component.Version.Version = rc.Annotations["chart.helm.sh/version"]
		cluster.Components = append(cluster.Components, component)
	}
	js, err := json.Marshal(cluster)
	if err != nil {
		log.Print(err)
		return []byte{}, err
	}
	return js, nil
}

// LatestReleasedComponent fulfills the AvailableComponentVersion interface
type LatestReleasedComponent struct{}

// Get method for LatestReleasedComponent
func (c LatestReleasedComponent) Get(component string) (types.Version, error) {
	version, err := GetLatestVersion(component)
	if err != nil {
		return types.Version{}, err
	}
	return version, nil
}

// GetCluster collects all cluster metadata and returns a Cluster
func GetCluster(c InstalledData, i data.ClusterID, v AvailableComponentVersion) (types.Cluster, error) {
	// Populate cluster object with installed components
	cluster, err := GetInstalled(c)
	if err != nil {
		log.Print(err)
		return types.Cluster{}, err
	}
	err = AddUpdateData(&cluster, v)
	if err != nil {
		log.Print(err)
	}
	// Get the cluster ID
	id, err := data.GetID(i)
	if err != nil {
		log.Print(err)
		return cluster, err
	}
	// Attach the cluster ID to the components-populated cluster object
	cluster.ID = id
	return cluster, nil
}

// AddUpdateData adds UpdateAvailable field data to cluster components
// Any cluster object modifications are made "in-place"
func AddUpdateData(c *types.Cluster, v AvailableComponentVersion) error {
	// Determine if any components have an available update
	for i, component := range c.Components {
		installed := component.Version.Version
		latest, err := v.Get(component.Component.Name)
		if err != nil {
			return err
		}
		newest := newestVersion(installed, latest.Version)
		if newest != installed {
			c.Components[i].UpdateAvailable = newest
		}
	}
	return nil
}

// GetInstalled collects all installed components and returns a Cluster
func GetInstalled(g InstalledData) (types.Cluster, error) {
	data, err := g.Get()
	if err != nil {
		log.Print(err)
		return types.Cluster{}, err
	}
	var cluster types.Cluster
	cluster, err = ParseJSONCluster(data)
	if err != nil {
		log.Print(err)
		return types.Cluster{}, err
	}
	return cluster, nil
}

// ParseJSONCluster converts a JSON representation of a cluster
// to a Cluster type
func ParseJSONCluster(rawJSON []byte) (types.Cluster, error) {
	var cluster types.Cluster
	err := json.Unmarshal(rawJSON, &cluster)
	if err != nil {
		log.Print(err)
		return types.Cluster{}, err
	}
	return cluster, nil
}

// GetLatestVersion returns the latest known version of a deis component
func GetLatestVersion(component string) (types.Version, error) {
	var latestVersion types.Version
	latestVersions, err := data.GetAvailableVersions(data.AvailableVersionsFromAPI{})
	if err != nil {
		return types.Version{}, err
	}
	for _, componentVersion := range latestVersions {
		if componentVersion.Component.Name == component {
			latestVersion = componentVersion.Version
		}
	}
	if latestVersion.Version == "" {
		return types.Version{}, fmt.Errorf("latest version not available for %s", component)
	}
	return latestVersion, nil
}

// NewestSemVer returns the newest (largest) semver string
func NewestSemVer(v1 string, v2 string) (string, error) {
	v1Slice := strings.Split(v1, ".")
	v2Slice := strings.Split(v2, ".")
	for i, subVer1 := range v1Slice {
		if v2Slice[i] > subVer1 {
			return v2, nil
		} else if subVer1 > v2Slice[i] {
			return v1, nil
		}
	}
	return v1, nil
}

// getDeisRCItems is a helper function that returns a slice of
// ReplicationController objects in the "deis" namespace
func getDeisRCItems() ([]api.ReplicationController, error) {
	kubeClient, err := kcl.NewInCluster()
	if err != nil {
		log.Printf("Error getting kubernetes client [%s]", err)
		return []api.ReplicationController{}, err
	}
	deis, err := kubeClient.ReplicationControllers(config.Spec.DeisNamespace).List(labels.Everything())
	if err != nil {
		log.Println("unable to get ReplicationControllers() data from kube client")
		log.Print(err)
		return []api.ReplicationController{}, err
	}
	return deis.Items, nil
}

// newestVersion is a temporary static implementation of a real "return newest version" function
func newestVersion(v1 string, v2 string) string {
	return v1
}
