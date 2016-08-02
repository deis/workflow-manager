package data

import (
	"github.com/deis/workflow-manager/k8s"
	"github.com/deis/workflow-manager/pkg/swagger/models"
)

// AvailableComponentVersion is an interface for managing component version data
type AvailableComponentVersion interface {
	// will have a Get method to retrieve available component version data
	Get(component string, cluster models.Cluster) (models.Version, error)
}

// latestReleasedComponent fulfills the AvailableComponentVersion interface
type latestReleasedComponent struct {
	k8sResources *k8s.ResourceInterfaceNamespaced
	availableVersions AvailableVersions
}

// NewLatestReleasedComponent creates a new AvailableComponentVersion that gets the latest released component using sgc as the implementation to get and create secrets
func NewLatestReleasedComponent(
	ri *k8s.ResourceInterfaceNamespaced,
	availableVersions AvailableVersions,
) AvailableComponentVersion {
	return &latestReleasedComponent{
		k8sResources: ri,
		availableVersions: availableVersions,
	}
}

// Get method for LatestReleasedComponent
func (c *latestReleasedComponent) Get(component string, cluster models.Cluster) (models.Version, error) {
	version, err := GetLatestVersion(
		component,
		cluster,
		c.availableVersions,
	)
	if err != nil {
		return models.Version{}, err
	}
	return version, nil
}
