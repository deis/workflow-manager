package data

import (
	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/types"
)

// AvailableComponentVersion is an interface for managing component version data
type AvailableComponentVersion interface {
	// will have a Get method to retrieve available component version data
	Get(component string) (types.Version, error)
}

// LatestReleasedComponent fulfills the AvailableComponentVersion interface
type LatestReleasedComponent struct {
	secretGetterCreator KubeSecretGetterCreator
	rcLister            rc.Lister
}

// NewLatestReleasedComponent creates a new LatestReleasedComponent using sgc as the implementation to get and create secrets
func NewLatestReleasedComponent(sgc KubeSecretGetterCreator, rcl rc.Lister) *LatestReleasedComponent {
	return &LatestReleasedComponent{secretGetterCreator: sgc, rcLister: rcl}
}

// Get method for LatestReleasedComponent
func (c LatestReleasedComponent) Get(component string) (types.Version, error) {
	version, err := GetLatestVersion(component, c.secretGetterCreator, c.rcLister)
	if err != nil {
		return types.Version{}, err
	}
	return version, nil
}
