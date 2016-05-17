package data

import (
	"log"
	"os"

	"k8s.io/kubernetes/pkg/api"
	k8sErrors "k8s.io/kubernetes/pkg/api/errors"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
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
