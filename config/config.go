package config

import "github.com/kelseyhightower/envconfig"

// Specification config struct
type Specification struct {
	Port           string `default:"8080"`
	Polling        int    `default:"43200"` // 43200 seconds = 12 hours
	VersionsAPIURL string
	DoctorAPIURL   string
	APIVersion     string
	CheckVersions  bool   `default:"true"`
	DeisNamespace  string `default:"deis"`
}

// Spec is an exportable variable that contains workflow manager config data
var Spec Specification

func init() {
	envconfig.Process("workflow_manager", &Spec)
}
