package data

import (
	"encoding/json"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/types"
)

// InstalledData is an interface for managing installed cluster metadata
type InstalledData interface {
	// will have a Get method to retrieve installed data
	Get() ([]byte, error)
}

// InstalledDeisData fulfills the InstalledData interface
type installedDeisData struct {
	rcLister rc.Lister
}

// NewInstalledDeisData returns a new InstalledDeisData using rcl as the rc.Lister implementation
func NewInstalledDeisData(rcl rc.Lister) InstalledData {
	return &installedDeisData{rcLister: rcl}
}

// Get method for InstalledDeisData
func (g *installedDeisData) Get() ([]byte, error) {
	rcItems, err := getDeisRCItems(g.rcLister)
	if err != nil {
		return nil, err
	}
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
		return []byte{}, err
	}
	return js, nil
}
