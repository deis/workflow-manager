package data

import (
	"sync"

	"github.com/deis/workflow-manager/k8s"
	"github.com/satori/go.uuid"
	"k8s.io/kubernetes/pkg/api"
	apierrors "k8s.io/kubernetes/pkg/api/errors"
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

// GetID gets the cluster ID from the cache. on a cache miss, uses the k8s API to get it
func GetID(id ClusterID) (string, error) {
	// First, check to see if we have an in-memory copy
	data := id.Cached()
	// If we haven't yet cached the ID in memory, invoke the passed-in getter
	if data == "" {
		d, err := id.Get()
		if err != nil {
			return "", err
		}
		data = d
	}
	return data, nil
}

type clusterIDFromPersistentStorage struct {
	rwm                 *sync.RWMutex
	cache               string
	secretGetterCreator k8s.KubeSecretGetterCreator
}

// NewClusterIDFromPersistentStorage returns a new ClusterID implementation that uses the kubernetes API to get its cluster information
func NewClusterIDFromPersistentStorage(sgc k8s.KubeSecretGetterCreator) ClusterID {
	return &clusterIDFromPersistentStorage{
		rwm:                 new(sync.RWMutex),
		cache:               "",
		secretGetterCreator: sgc,
	}
}

// Get is the ClusterID interface implementation
func (c clusterIDFromPersistentStorage) Get() (string, error) {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	secret, err := c.secretGetterCreator.Get(wfmSecretName)
	//If we don't have the secret we shouldn't be returning error and instead a create a new one
	if err != nil && !apierrors.IsNotFound(err) {
		return "", err
	}
	// if we don't have secret data for the cluster ID we assume a new cluster
	// and create a new secret
	if secret.Data[clusterIDSecretKey] == nil {
		newSecret := new(api.Secret)
		newSecret.Name = wfmSecretName
		newSecret.Data = make(map[string][]byte)
		newSecret.Data[clusterIDSecretKey] = []byte(uuid.NewV4().String())
		fromAPI, err := c.secretGetterCreator.Create(newSecret)
		if err != nil {
			return "", err
		}
		secret = fromAPI
	}
	return string(secret.Data[clusterIDSecretKey]), nil
}

// StoreInCache is the ClusterID interface implementation
func (c *clusterIDFromPersistentStorage) StoreInCache(cid string) {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	c.cache = cid
}

// Cached is the ClusterID interface implementation
func (c clusterIDFromPersistentStorage) Cached() string {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	return c.cache
}
