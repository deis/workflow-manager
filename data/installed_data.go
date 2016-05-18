package data

import (
	"encoding/json"
	"log"

	"github.com/deis/workflow-manager/types"
)

// InstalledData is an interface for managing installed cluster metadata
type InstalledData interface {
	// will have a Get method to retrieve installed data
	Get() ([]byte, error)
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
