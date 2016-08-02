package k8s

import (
	"github.com/deis/kubeapp/api/secret"
)

// KubeSecretGetterCreator is a composition of secret.Getter and secret.Creator. Please refer to the Godoc for those two interfaces (https://godoc.org/github.com/arschles/kubeapp/api/secret)
type KubeSecretGetterCreator interface {
	secret.Getter
	secret.Creator
}

// FakeKubeSecretGetterCreator is a composition of the secret.FakeGetter and secret.FakeCreator structs
type FakeKubeSecretGetterCreator struct {
	*secret.FakeGetter
	*secret.FakeCreator
}

// NewFakeKubeSecretGetterCreator creates a new FakeKubeSecretGetterCreator from the given fakeGetter and fakeCreator
func NewFakeKubeSecretGetterCreator(fakeGetter *secret.FakeGetter, fakeCreator *secret.FakeCreator) *FakeKubeSecretGetterCreator {
	return &FakeKubeSecretGetterCreator{FakeGetter: fakeGetter, FakeCreator: fakeCreator}
}
