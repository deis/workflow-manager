package data

import (
	"sync"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/kubeapp/api/secret"
	"github.com/deis/workflow-manager/k8s"
	"github.com/satori/go.uuid"
	"k8s.io/kubernetes/pkg/api"
)

func TestClusterIDFromPersistentStorage(t *testing.T) {
	sec := api.Secret{}
	secretGetter := &secret.FakeGetter{
		Secret: &sec,
	}
	secretCreator := &secret.FakeCreator{
		CreateFunc: func(sec *api.Secret) (*api.Secret, error) {
			return sec, nil
		},
	}
	secrets := k8s.NewFakeKubeSecretGetterCreator(secretGetter, secretCreator)
	clusterID := clusterIDFromPersistentStorage{
		rwm:                 new(sync.RWMutex),
		cache:               "",
		secretGetterCreator: secrets,
	}
	resp, err := clusterID.Get()
	assert.NoErr(t, err)
	// verify that the secret is a UUID
	_, err = uuid.FromString(resp)
	assert.NoErr(t, err)
}
