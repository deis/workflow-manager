package data

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/types"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

const (
	wfmSecretName      = "deis-workflow-manager"
	clusterIDSecretKey = "cluster-id"
)

// GetCluster collects all cluster metadata and returns a Cluster
func GetCluster(
	c InstalledData,
	i ClusterID,
	v AvailableComponentVersion,
	secretGetterCreator KubeSecretGetterCreator,
) (types.Cluster, error) {

	// Populate cluster object with installed components
	cluster, err := GetInstalled(c)
	if err != nil {
		log.Print(err)
		return types.Cluster{}, err
	}
	err = AddUpdateData(&cluster, v, secretGetterCreator)
	if err != nil {
		log.Print(err)
	}
	// Get the cluster ID
	id, err := GetID(i)
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
func AddUpdateData(c *types.Cluster, v AvailableComponentVersion, secretGetterCreator KubeSecretGetterCreator) error {
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

// GetAvailableVersions gets available component version data from the cache. If there was a cache miss, gets the versions from the k8s and versions APIs
func GetAvailableVersions(a AvailableVersions) ([]types.ComponentVersion, error) {
	// First, check to see if we have an in-memory copy
	data := a.Cached()
	// If we don't have any cached data, get the data from the remote authority
	if len(data) == 0 {
		d, err := a.Refresh()
		if err != nil {
			return nil, err
		}
		data = d
	}
	return data, nil
}

// GetInstalled collects all installed components and returns a Cluster
func GetInstalled(g InstalledData) (types.Cluster, error) {
	installed, err := g.Get()
	if err != nil {
		log.Print(err)
		return types.Cluster{}, err
	}
	var cluster types.Cluster
	cluster, err = ParseJSONCluster(installed)
	if err != nil {
		log.Print(err)
		return types.Cluster{}, err
	}
	return cluster, nil
}

// GetLatestVersion returns the latest known version of a deis component
func GetLatestVersion(
	component string,
	secretGetterCreator KubeSecretGetterCreator,
	rcLister rc.Lister,
) (types.Version, error) {
	var latestVersion types.Version
	latestVersions, err := GetAvailableVersions(NewAvailableVersionsFromAPI("", secretGetterCreator, rcLister))
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
func getDeisRCItems(rcLister rc.Lister) ([]api.ReplicationController, error) {
	rcs, err := rcLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []api.ReplicationController{}, err
	}
	return rcs.Items, nil
}

// getTLSClient returns a TLS-enabled http.Client
func getTLSClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	return &http.Client{Transport: tr}
}

// newestVersion is a temporary static implementation of a real "return newest version" function
func newestVersion(v1 string, v2 string) string {
	return v1
}
