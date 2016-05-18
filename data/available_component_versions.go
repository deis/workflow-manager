package data

import (
	"github.com/deis/workflow-manager/types"
)

// AvailableComponentVersion is an interface for managing component version data
type AvailableComponentVersion interface {
	// will have a Get method to retrieve available component version data
	Get(component string) (types.Version, error)
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
