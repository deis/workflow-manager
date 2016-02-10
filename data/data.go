package data

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/types"
	"github.com/satori/go.uuid"
	"k8s.io/kubernetes/pkg/api"
	k8sErrors "k8s.io/kubernetes/pkg/api/errors"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	deisNamespace      = "deis"
	wfmSecretName      = "deis-workflow-manager"
	clusterIDSecretKey = "cluster-id"
)

var (
	clusterID         string
	availableVersions []types.ComponentVersion
)

// ClusterID is an interface for managing cluster ID data
type ClusterID interface {
	// will have a Get method to retrieve the cluster ID
	Get() (string, error)
}

// ClusterIDFromPersistentStorage fulfills the ClusterID interface
type ClusterIDFromPersistentStorage struct{}

// Get method for ClusterIDFromPersistentStorage
func (c ClusterIDFromPersistentStorage) Get() (string, error) {
	kubeClient, err := kcl.NewInCluster()
	if err != nil {
		log.Printf("Error getting kubernetes client [%s]", err)
		os.Exit(1)
	}
	deisSecrets := kubeClient.Secrets(deisNamespace)
	secret, err := deisSecrets.Get(wfmSecretName)
	if err != nil {
		log.Printf("Error getting secret [%s]", err)
		switch e := err.(type) {
		case *k8sErrors.StatusError:
			// If the error isn't a 404, we don't know how to deal with it
			if e.ErrStatus.Code != 404 {
				return "", err
			}
		default:
			return "", err
		}
	}
	// if we don't have secret data for the cluster ID we assume a new cluster
	// and create a new secret
	if secret.Data[clusterIDSecretKey] == nil {
		newSecret := new(api.Secret)
		newSecret.Name = wfmSecretName
		newSecret.Data = make(map[string][]byte)
		newSecret.Data[clusterIDSecretKey] = []byte(getNewID())
		fromAPI, err := deisSecrets.Create(newSecret)
		if err != nil {
			log.Printf("Error creating new ID [%s]", err)
			return "", err
		}
		secret = fromAPI
	}
	return string(secret.Data[clusterIDSecretKey]), nil
}

// AvailableVersions is an interface for managing available component version data
type AvailableVersions interface {
	// will have a Refresh method to retrieve the version data from the remote authority
	Refresh() ([]types.ComponentVersion, error)
	// will have a Store method to cache the version data in memory
	Store([]types.ComponentVersion)
}

// AvailableVersionsFromAPI fulfills the AvailableVersions interface
type AvailableVersionsFromAPI struct{}

// Refresh method for AvailableVersionsFromAPI
func (a AvailableVersionsFromAPI) Refresh() ([]types.ComponentVersion, error) {
	var versionsRoute = "/" + config.Spec.APIVersion + "/versions"
	url := config.Spec.VersionsAPIURL + versionsRoute
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []types.ComponentVersion{}, err
	}

	resp, err := getTLSClient().Do(req)
	if err != nil {
		return []types.ComponentVersion{}, err
	}
	defer resp.Body.Close()
	var availableVersions []types.ComponentVersion
	json.NewDecoder(resp.Body).Decode(&availableVersions)
	a.Store(availableVersions)
	return availableVersions, nil
}

// Store method for AvailableVersionsFromAPI
func (a AvailableVersionsFromAPI) Store(c []types.ComponentVersion) {
	availableVersions = c
}

// getClusterID returns the cached ID, or an error if it's not cached in memory
func getClusterID() (string, error) {
	if clusterID == "" {
		return "", fmt.Errorf("cluster ID not cached in memory")
	}
	return clusterID, nil
}

// getAvailableVersions returns the cached available version data, or an error
func getAvailableVersions() ([]types.ComponentVersion, error) {
	if len(availableVersions) == 0 {
		return nil, fmt.Errorf("no available versions data cached")
	}
	return availableVersions, nil
}

// GetID gets the cluster ID
func GetID(id ClusterID) (string, error) {
	// First, check to see if we have an in-memory copy
	data, err := getClusterID()
	// If we haven't yet cached the ID in memory, invoke the passed-in getter
	if err != nil {
		data, err = id.Get()
		if err != nil {
			log.Print(err)
			return "", err
		}
		clusterID = data
	}
	return data, nil
}

// GetAvailableVersions gets available component version data
func GetAvailableVersions(a AvailableVersions) ([]types.ComponentVersion, error) {
	// First, check to see if we have an in-memory copy
	data, err := getAvailableVersions()
	// If we don't have any cached data, get the data from the remote authority
	if err != nil {
		data, err = a.Refresh()
		if err != nil {
			log.Print(err)
			return nil, err
		}
		return data, nil
	}
	return data, nil
}

// getNewID returns a new Cluster ID string value
func getNewID() string {
	return uuid.NewV4().String()
}

// getTLSClient returns a TLS-enabled http.Client
func getTLSClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	return &http.Client{Transport: tr}
}
