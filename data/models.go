package data

import (
	"github.com/deis/workflow-manager/pkg/swagger/models"
)

// SparseComponentInfo is the JSON compatible struct that holds limited data about a component
type SparseComponentInfo struct {
	Name string `json:"name"`
}

// SparseVersionInfo is the JSON compatible struct that holds limited data about a
// component version
type SparseVersionInfo struct {
	Train string `json:"train"`
}

// SparseComponentAndTrainInfo is the JSON compatible struct that holds a
// SparseComponentInfo and SparseVersionInfo
type SparseComponentAndTrainInfo struct {
	Component SparseComponentInfo `json:"component"`
	Version   SparseVersionInfo   `json:"version"`
}

// SparseComponentAndTrainInfoJSONWrapper is the JSON compatible struct that holds a slice of
// SparseComponentAndTrainInfo structs
type SparseComponentAndTrainInfoJSONWrapper struct {
	Data []SparseComponentAndTrainInfo `json:"data"`
}

// ComponentVersionsJSONWrapper is the JSON compatible struct that holds a slice of
// types.ComponentVersion structs
type ComponentVersionsJSONWrapper struct {
	Data []models.ComponentVersion `json:"data"`
}
