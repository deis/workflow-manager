package types

import "time"

// Component type definition
type Component struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Version type definition
type Version struct {
	Train    string                 `json:"train"` // e.g., "beta", "stable"
	Version  string                 `json:"version"`
	Released string                 `json:"released,omitempty"`
	Data     map[string]interface{} `json:"data"`
}

// ComponentVersion type definition
type ComponentVersion struct {
	Component       Component `json:"component"`
	Version         Version   `json:"version"`
	UpdateAvailable string    `json:"updateAvailable,omitempty"`
}

// Cluster type definition
type Cluster struct {
	ID string `json:"id"`
	// FirstSeen and/or LastSeen suggests a Cluster object in a lifecycle context,
	// i.e., for use in business logic which needs to determine a cluster's "freshness" or "staleness"
	// example use case: we omit these properties when submitting information to the versions API
	// another example use case: we populate these properties when gathering lifecycle statistics from the API
	FirstSeen  time.Time          `json:"firstSeen,omitempty"`
	LastSeen   time.Time          `json:"lastSeen,omitempty"`
	Components []ComponentVersion `json:"components"`
}
