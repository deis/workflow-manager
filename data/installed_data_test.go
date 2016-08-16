package data

import (
	"testing"

	"github.com/deis/workflow-manager/k8s"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/testapi"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/unversioned/testclient/simple"
)

const namespace = "deis"

func TestInstalledDeisData(t *testing.T) {
	client := getK8sClient(t)
	k := k8s.NewResourceInterfaceNamespaced(client, namespace)
	installedData := installedDeisData{
		k8sResources: k,
	}
	_, _ = installedData.Get()
	//TODO: we need to create a fake client interface that is a union of api+extension fake clients
}

func getK8sClient(t *testing.T) *simple.Client {
	c := &simple.Client{
		Request: simple.Request{
			Method: "GET",
			Path:   testapi.Extensions.ResourcePath("deployments", namespace, ""),
		},
		Response: simple.Response{StatusCode: 200,
			Body: &extensions.DeploymentList{
				Items: []extensions.Deployment{
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-builder",
						},
						Spec: extensions.DeploymentSpec{
							Template: api.PodTemplateSpec{
								Spec: api.PodSpec{
									Containers: []api.Container{
										{
											Image: "container-image",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return c.Setup(t)
}
