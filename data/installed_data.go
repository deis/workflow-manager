package data

import (
	"encoding/json"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/pkg/swagger/models"
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
	rcItems, err := getRCItems(g.rcLister)
	if err != nil {
		return nil, err
	}
	var cluster models.Cluster
	for _, rc := range rcItems {
		component := models.ComponentVersion{}
		component.Component = &models.Component{}
		component.Version = &models.Version{}
		component.Component.Name = rc.Name
		desc := rc.Annotations["chart.helm.sh/description"]
		component.Component.Description = &desc
		component.Version.Version = rc.Annotations["chart.helm.sh/version"]
		cluster.Components = append(cluster.Components, &component)
	}
	js, err := json.Marshal(cluster)
	if err != nil {
		return []byte{}, err
	}
	return js, nil
}
