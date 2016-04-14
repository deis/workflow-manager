package types

import "time"

// Component type definition
type Component struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Version type definition
type Version struct {
	Train    string `json:"train"` // e.g., "beta", "stable"
	Version  string `json:"version"`
	Released string `json:"released,omitempty"`
	Data     []byte `json:"data"`
}

// ComponentVersion type definition
type ComponentVersion struct {
	Component       Component `json:"component"`
	Version         Version   `json:"version"`
	UpdateAvailable string    `json:"updateAvailable,omitempty"`
}

// Cluster type definition
type Cluster struct {
	ID         string             `json:"id"`
	FirstSeen  time.Time          `json:"firstSeen"`
	LastSeen   time.Time          `json:"lastSeen"`
	Components []ComponentVersion `json:"components"`
}
