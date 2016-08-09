package k8s

import (
	"testing"

	"github.com/arschles/assert"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/testapi"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/unversioned/testclient/simple"
)

const namespace = "deis"

func TestRunningK8sDataDaemonsets(t *testing.T) {
	c := getK8sClientForDaemonSets(t)
	deisK8sResources := NewResourceInterfaceNamespaced(c, namespace)
	runningK8sData := NewRunningK8sData(deisK8sResources)
	daemonsets, err := runningK8sData.DaemonSets()
	assert.NoErr(t, err)
	assert.True(t, len(daemonsets) == 2, "daemonsets response slice was not the expected length")
	assert.True(t, daemonsets[0].Data != daemonsets[1].Data, "daemonsets should not be identical")
}

func TestRunningK8sDataDeployments(t *testing.T) {
	c := getK8sClientForDeployments(t)
	deisK8sResources := NewResourceInterfaceNamespaced(c, namespace)
	runningK8sData := NewRunningK8sData(deisK8sResources)
	deployments, err := runningK8sData.Deployments()
	assert.NoErr(t, err)
	assert.True(t, len(deployments) == 2, "deployments response slice was not the expected length")
	assert.True(t, deployments[0].Data != deployments[1].Data, "deployments should not be identical")
}

func TestRunningK8sDataEvents(t *testing.T) {
	c := getK8sClientForEvents(t)
	deisK8sResources := NewResourceInterfaceNamespaced(c, namespace)
	runningK8sData := NewRunningK8sData(deisK8sResources)
	events, err := runningK8sData.Events()
	assert.NoErr(t, err)
	assert.True(t, len(events) == 2, "events response slice was not the expected length")
	assert.True(t, events[0].Data != events[1].Data, "events should not be identical")
}

func TestRunningK8sDataNodes(t *testing.T) {
	c := getK8sClientForNodes(t)
	deisK8sResources := NewResourceInterfaceNamespaced(c, namespace)
	runningK8sData := NewRunningK8sData(deisK8sResources)
	nodes, err := runningK8sData.Nodes()
	assert.NoErr(t, err)
	assert.True(t, len(nodes) == 3, "nodes response slice was not the expected length")
	assert.True(t, nodes[0].Data != nodes[1].Data, "nodes should not be identical")
}

func TestRunningK8sDataPods(t *testing.T) {
	c := getK8sClientForPods(t)
	deisK8sResources := NewResourceInterfaceNamespaced(c, namespace)
	runningK8sData := NewRunningK8sData(deisK8sResources)
	pods, err := runningK8sData.Pods()
	assert.NoErr(t, err)
	assert.True(t, len(pods) == 2, "pods response slice was not the expected length")
	assert.True(t, pods[0].Data != pods[1].Data, "pods should not be identical")
}

func TestRunningK8sDataReplicaSets(t *testing.T) {
	c := getK8sClientForReplicaSets(t)
	deisK8sResources := NewResourceInterfaceNamespaced(c, namespace)
	runningK8sData := NewRunningK8sData(deisK8sResources)
	replicaSets, err := runningK8sData.ReplicaSets()
	assert.NoErr(t, err)
	assert.True(t, len(replicaSets) == 2, "replica sets response slice was not the expected length")
	assert.True(t, replicaSets[0].Data != replicaSets[1].Data, "replica sets should not be identical")
}

func TestRunningK8sDataReplicationControllers(t *testing.T) {
	c := getK8sClientForReplicationControllers(t)
	deisK8sResources := NewResourceInterfaceNamespaced(c, namespace)
	runningK8sData := NewRunningK8sData(deisK8sResources)
	rcs, err := runningK8sData.ReplicationControllers()
	assert.NoErr(t, err)
	assert.True(t, len(rcs) == 2, "rc response slice was not the expected length")
	assert.True(t, rcs[0].Data != rcs[1].Data, "rcs should not be identical")
}

func TestRunningK8sDataServices(t *testing.T) {
	c := getK8sClientForServices(t)
	//_, _ = c.Setup(t).Services(namespace).List(api.ListOptions{})
	deisK8sResources := NewResourceInterfaceNamespaced(c, namespace)
	runningK8sData := NewRunningK8sData(deisK8sResources)
	services, err := runningK8sData.Services()
	assert.NoErr(t, err)
	assert.True(t, len(services) == 2, "services response slice was not the expected length")
	assert.True(t, services[0].Data != services[1].Data, "services should not be identical")
}

func getK8sClientForDaemonSets(t *testing.T) *simple.Client {
	c := &simple.Client{
		Request: simple.Request{
			Method: "GET",
			Path:   testapi.Extensions.ResourcePath("daemonsets", namespace, ""),
		},
		Response: simple.Response{StatusCode: 200,
			Body: &extensions.DaemonSetList{
				Items: []extensions.DaemonSet{
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-logger-fluentd",
						},
					},
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-monitor-telegraf",
						},
					},
				},
			},
		},
	}
	return c.Setup(t)
}

func getK8sClientForDeployments(t *testing.T) *simple.Client {
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
					},
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-router",
						},
					},
				},
			},
		},
	}
	return c.Setup(t)
}

func getK8sClientForEvents(t *testing.T) *simple.Client {
	obj1Reference := &api.ObjectReference{
		Kind:            "Pod",
		Namespace:       namespace,
		Name:            "deis-builder-900960817-h9zmm",
		UID:             "uid",
		APIVersion:      "v1",
		ResourceVersion: "1477748",
	}
	obj2Reference := &api.ObjectReference{
		Kind:            "Pod",
		Namespace:       namespace,
		Name:            "deis-controller-139932026-oltfd",
		UID:             "uid",
		APIVersion:      "v1",
		ResourceVersion: "1477762",
	}
	timeStamp := unversioned.Now()
	timeStamp2 := unversioned.Now()
	eventList := &api.EventList{
		Items: []api.Event{
			{
				InvolvedObject: *obj1Reference,
				FirstTimestamp: timeStamp,
				LastTimestamp:  timeStamp,
				Count:          1,
				Type:           api.EventTypeNormal,
			},
			{
				InvolvedObject: *obj2Reference,
				FirstTimestamp: timeStamp2,
				LastTimestamp:  timeStamp2,
				Count:          1,
				Type:           api.EventTypeNormal,
			},
		},
	}
	c := &simple.Client{
		Request: simple.Request{
			Method: "GET",
			Path:   testapi.Default.ResourcePath("events", namespace, ""),
			Body:   nil,
		},
		Response: simple.Response{StatusCode: 200, Body: eventList},
	}
	return c.Setup(t)
}

func getK8sClientForNodes(t *testing.T) *simple.Client {
	c := &simple.Client{
		Request: simple.Request{
			Method: "GET",
			Path:   testapi.Default.ResourcePath("nodes", "", ""),
		},
		Response: simple.Response{
			StatusCode: 200, Body: &api.NodeList{
				Items: []api.Node{
					{
						ObjectMeta: api.ObjectMeta{
							Labels: map[string]string{
								"kubernetes.io/hostname": "gke-test-default-pool-0e781ee9-1xz6",
							},
							Name: "gke-test-default-pool-0e781ee9-1xz6",
						},
					},
					{
						ObjectMeta: api.ObjectMeta{
							Labels: map[string]string{
								"kubernetes.io/hostname": "gke-test-default-pool-0e781ee9-9j80",
							},
							Name: "gke-test-default-pool-0e781ee9-9j80",
						},
					},
					{
						ObjectMeta: api.ObjectMeta{
							Labels: map[string]string{
								"kubernetes.io/hostname": "gke-francis-default-pool-0e781ee9-z02l",
							},
							Name: "gke-francis-default-pool-0e781ee9-z02l",
						},
					},
				},
			},
		},
	}
	return c.Setup(t)
}

func getK8sClientForPods(t *testing.T) *simple.Client {
	c := &simple.Client{
		Request: simple.Request{
			Method: "GET",
			Path:   testapi.Default.ResourcePath("pods", namespace, ""),
		},
		Response: simple.Response{StatusCode: 200,
			Body: &api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-builder-900960817-h9zmm",
						},
					},
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-controller-139932026-oltfd",
						},
					},
				},
			},
		},
	}
	return c.Setup(t)
}

func getK8sClientForReplicaSets(t *testing.T) *simple.Client {
	c := &simple.Client{
		Request: simple.Request{
			Method: "GET",
			Path:   testapi.Extensions.ResourcePath("replicasets", namespace, ""),
		},
		Response: simple.Response{StatusCode: 200,
			Body: &extensions.ReplicaSetList{
				Items: []extensions.ReplicaSet{
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-builder-900960817",
						},
					},
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-controller-139932026",
						},
					},
				},
			},
		},
	}
	return c.Setup(t)
}

func getK8sClientForReplicationControllers(t *testing.T) *simple.Client {
	c := &simple.Client{
		Request: simple.Request{
			Method: "GET",
			Path:   testapi.Default.ResourcePath("replicationcontrollers", namespace, ""),
		},
		Response: simple.Response{StatusCode: 200,
			Body: &api.ReplicationControllerList{
				Items: []api.ReplicationController{
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-builder",
						},
					},
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-controller",
						},
					},
				},
			},
		},
	}
	return c.Setup(t)
}

func getK8sClientForServices(t *testing.T) *simple.Client {
	c := &simple.Client{
		Request: simple.Request{
			Method: "GET",
			Path:   testapi.Default.ResourcePath("services", namespace, ""),
		},
		Response: simple.Response{StatusCode: 200,
			Body: &api.ServiceList{
				Items: []api.Service{
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-builder",
						},
					},
					{
						ObjectMeta: api.ObjectMeta{
							Name: "deis-controller",
						},
					},
				},
			},
		},
	}
	return c.Setup(t)
}
