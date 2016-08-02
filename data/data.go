package data

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/k8s"
	"github.com/deis/workflow-manager/pkg/swagger/models"
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
) (models.Cluster, error) {

	// Populate cluster object with installed components
	cluster, err := GetInstalled(c)
	if err != nil {
		return models.Cluster{}, err
	}
	if err := AddUpdateData(&cluster, v); err != nil {
		log.Printf("unable to decorate cluster data with available updates data: %#v", err)
	}
	// Get the cluster ID
	id, err := GetID(i)
	if err != nil {
		return cluster, err
	}
	// Attach the cluster ID to the components-populated cluster object
	cluster.ID = id
	return cluster, nil
}

// AddUpdateData adds UpdateAvailable field data to cluster components
// Any cluster object modifications are made "in-place"
func AddUpdateData(c *models.Cluster, v AvailableComponentVersion) error {
	// Determine if any components have an available update
	for i, component := range c.Components {
		installed := component.Version.Version
		latest, err := v.Get(component.Component.Name, *c)
		if err != nil {
			return err
		}
		newest := newestVersion(installed, latest.Version)
		if newest != installed {
			c.Components[i].UpdateAvailable = &newest
		}
	}
	return nil
}

// GetAvailableVersions gets available component version data from the cache. If there was a cache miss, gets the versions from the k8s and versions APIs
func GetAvailableVersions(a AvailableVersions, cluster models.Cluster) ([]models.ComponentVersion, error) {
	// First, check to see if we have an in-memory copy
	data := a.Cached()
	// If we don't have any cached data, get the data from the remote authority
	if len(data) == 0 {
		d, err := a.Refresh(cluster)
		if err != nil {
			return nil, err
		}
		data = d
	}
	return data, nil
}

// GetInstalled collects all installed components and returns a Cluster
func GetInstalled(g InstalledData) (models.Cluster, error) {
	installed, err := g.Get()
	if err != nil {
		return models.Cluster{}, err
	}
	var cluster models.Cluster
	cluster, err = ParseJSONCluster(installed)
	if err != nil {
		return models.Cluster{}, err
	}
	return cluster, nil
}

// GetLatestVersion returns the latest known version of a deis component
func GetLatestVersion(
	component string,
	cluster models.Cluster,
	availVsns AvailableVersions,
) (models.Version, error) {
	var latestVersion models.Version
	latestVersions, err := GetAvailableVersions(availVsns, cluster)
	if err != nil {
		return models.Version{}, err
	}
	for _, componentVersion := range latestVersions {
		if componentVersion.Component.Name == component {
			latestVersion = *componentVersion.Version
		}
	}
	if latestVersion.Version == "" {
		return models.Version{}, fmt.Errorf("latest version not available for %s", component)
	}
	return latestVersion, nil
}

// ParseJSONCluster converts a JSON representation of a cluster
// to a Cluster type
func ParseJSONCluster(rawJSON []byte) (models.Cluster, error) {
	var cluster models.Cluster
	if err := json.Unmarshal(rawJSON, &cluster); err != nil {
		return models.Cluster{}, err
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

// GetDoctorInfo collects doctor info and return DoctorInfo struct
func GetDoctorInfo(
	c InstalledData, // workflow cluster data
	k k8s.RunningK8sData, // k8s data
	i ClusterID,
	v AvailableComponentVersion,
) (models.DoctorInfo, error) {
	cluster, err := GetCluster(c, i, v)
	if err != nil {
		return models.DoctorInfo{}, err
	}
	nodes := getK8sNodes(k)
	namespaces := []*models.Namespace{getK8sDeisNamespace(k)}
	doctor := models.DoctorInfo{
		Workflow:   &cluster,
		Nodes:      nodes,
		Namespaces: namespaces,
	}
	return doctor, nil
}

// getK8sDeisNamespace is a helper function that returns data
// from the "deis" K8s namespace for RESTful consumption
func getK8sDeisNamespace(k k8s.RunningK8sData) *models.Namespace {
	pods, err := k8s.GetPodsModels(k)
	if err != nil {
		log.Printf("unable to get K8s pods data: %#v", err)
	}
	services, err := k8s.GetServicesModels(k)
	if err != nil {
		log.Printf("unable to get K8s services data: %#v", err)
	}
	replicationControllers, err := k8s.GetReplicationControllersModels(k)
	if err != nil {
		log.Printf("unable to get K8s RC data: %#v", err)
	}
	replicaSets, err := k8s.GetReplicaSetsModels(k)
	if err != nil {
		log.Printf("unable to get K8s replicaSets data: %#v", err)
	}
	daemonSets, err := k8s.GetDaemonSetsModels(k)
	if err != nil {
		log.Printf("unable to get K8s daemonSets data: %#v", err)
	}
	deployments, err := k8s.GetDeploymentsModels(k)
	if err != nil {
		log.Printf("unable to get K8s deployments data: %#v", err)
	}
	events, err := k8s.GetEventsModels(k)
	if err != nil {
		log.Printf("unable to get K8s events data: %#v", err)
	}
	return &models.Namespace{
		Name:                   config.Spec.DeisNamespace,
		DaemonSets:             daemonSets,
		Deployments:            deployments,
		Events:                 events,
		Pods:                   pods,
		ReplicaSets:            replicaSets,
		ReplicationControllers: replicationControllers,
		Services:               services,
	}
}

// getK8sNodes is a helper function that returns K8s nodes data for RESTful consumption
func getK8sNodes(k k8s.RunningK8sData) []*models.K8sResource {
	nodes, err := k8s.GetNodesModels(k)
	if err != nil {
		log.Println("unable to get K8s nodes data")
	}
	return nodes
}

// newestVersion is a temporary static implementation of a real "return newest version" function
func newestVersion(v1 string, v2 string) string {
	return v1
}
