package data

import (
	"log"
	"os"
	"sync"

	"k8s.io/kubernetes/pkg/api"
	k8sErrors "k8s.io/kubernetes/pkg/api/errors"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// ClusterID is an interface for managing cluster ID data
type ClusterID interface {
	// will have a Get method to retrieve the cluster ID
	Get() (string, error)
	// Cached returns the internal cache of the cluster ID. returns the empty string on a miss
	Cached() string
	// StoreInCache stores the given string in the internal cluster ID cache
	StoreInCache(string)
}

// ClusterIDFromPersistentStorage fulfills the ClusterID interface
type ClusterIDFromPersistentStorage struct {
	rwm   *sync.RWMutex
	cache string
}

// NewClusterIDFromPersistentStorage returns a new ClusterID implementation that uses the kubernetes API to get its cluster information
func NewClusterIDFromPersistentStorage() *ClusterIDFromPersistentStorage {
	return &ClusterIDFromPersistentStorage{rwm: new(sync.RWMutex), cache: ""}
}

// Get is the ClusterID interface implementation
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

// StoreInCache is the ClusterID interface implementation
func (c *ClusterIDFromPersistentStorage) StoreInCache(cid string) {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	c.cache = cid
}

// Cached is the ClusterID interface implementation
func (c ClusterIDFromPersistentStorage) Cached() string {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	return c.cache
}
