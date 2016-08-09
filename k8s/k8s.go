package k8s

import (
	"github.com/deis/kubeapp/api/daemonset"
	"github.com/deis/kubeapp/api/deployment"
	"github.com/deis/kubeapp/api/event"
	"github.com/deis/kubeapp/api/node"
	"github.com/deis/kubeapp/api/pod"
	"github.com/deis/kubeapp/api/rc"
	"github.com/deis/kubeapp/api/replicaset"
	"github.com/deis/kubeapp/api/service"
	"github.com/deis/workflow-manager/pkg/swagger/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/labels"
)

// RunningK8sData is an interface for representing installed K8s data RESTFully
type RunningK8sData interface {
	// get DaemonSet model data for RESTful consumption
	DaemonSets() ([]*models.K8sResource, error)
	// get Deployment model data for RESTful consumption
	Deployments() ([]*models.K8sResource, error)
	// get Event model data for RESTful consumption
	Events() ([]*models.K8sResource, error)
	// get Node model data for RESTful consumption
	Nodes() ([]*models.K8sResource, error)
	// get Pod model data for RESTful consumption
	Pods() ([]*models.K8sResource, error)
	// get ReplicaSet model data for RESTful consumption
	ReplicaSets() ([]*models.K8sResource, error)
	// get ReplicationController model data for RESTful consumption
	ReplicationControllers() ([]*models.K8sResource, error)
	// get Service model data for RESTful consumption
	Services() ([]*models.K8sResource, error)
}

// runningK8sData fulfills the RunningK8sData interface
type runningK8sData struct {
	daemonSetLister  daemonset.Lister
	deploymentLister deployment.Lister
	eventLister      event.Lister
	nodeLister       node.Lister
	podLister        pod.Lister
	rcLister         rc.Lister
	replicaSetLister replicaset.Lister
	serviceLister    service.Lister
}

// NewRunningK8sData returns a new runningK8sData using rcl as the rc.Lister implementation
func NewRunningK8sData(r *ResourceInterfaceNamespaced) RunningK8sData {
	return &runningK8sData{
		daemonSetLister:  r.DaemonSets(),
		deploymentLister: r.Deployments(),
		eventLister:      r.Events(),
		nodeLister:       r.Nodes(),
		podLister:        r.Pods(),
		rcLister:         r.ReplicationControllers(),
		replicaSetLister: r.ReplicaSets(),
		serviceLister:    r.Services(),
	}
}

// DaemonSets method for runningK8sData
func (rkd *runningK8sData) DaemonSets() ([]*models.K8sResource, error) {
	ds, err := getDaemonSets(rkd.daemonSetLister)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.K8sResource, len(ds))
	for i, d := range ds {
		d2 := d // grab a value copy of "d" to enforce block scope heap reference, and avoid shadowing non-block scope "d"
		daemonSet := &models.K8sResource{
			Data: &d2,
		}
		ret[i] = daemonSet
	}
	return ret, nil
}

// Deployments method for runningK8sData
func (rkd *runningK8sData) Deployments() ([]*models.K8sResource, error) {
	ds, err := getDeployments(rkd.deploymentLister)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.K8sResource, len(ds))
	for i, d := range ds {
		d2 := d // grab a value copy of "d" to enforce block scope heap reference, and avoid shadowing non-block scope "d"
		dep := &models.K8sResource{
			Data: &d2,
		}
		ret[i] = dep
	}
	return ret, nil
}

// Events method for runningK8sData
func (rkd *runningK8sData) Events() ([]*models.K8sResource, error) {
	events, err := getEvents(rkd.eventLister)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.K8sResource, len(events))
	for i, e := range events {
		e2 := e // grab a value copy of "e" to enforce block scope heap reference, and avoid shadowing non-block scope "e"
		event := &models.K8sResource{
			Data: &e2,
		}
		ret[i] = event
	}
	return ret, nil
}

// Nodes method for runningK8sData
func (rkd *runningK8sData) Nodes() ([]*models.K8sResource, error) {
	nodes, err := getNodes(rkd.nodeLister)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.K8sResource, len(nodes))
	for i, n := range nodes {
		n2 := n // grab a value copy of "n" to enforce block scope heap reference, and avoid shadowing non-block scope "n"
		node := &models.K8sResource{
			Data: &n2,
		}
		ret[i] = node
	}
	return ret, nil
}

// Pods method for runningK8sData
func (rkd *runningK8sData) Pods() ([]*models.K8sResource, error) {
	pods, err := getPods(rkd.podLister)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.K8sResource, len(pods))
	for i, p := range pods {
		p2 := p // grab a value copy of "p" to enforce block scope heap reference, and avoid shadowing non-block scope "p"
		pod := &models.K8sResource{
			Data: &p2,
		}
		ret[i] = pod
	}
	return ret, nil
}

// ReplicaSets method for runningK8sData
func (rkd *runningK8sData) ReplicaSets() ([]*models.K8sResource, error) {
	rs, err := getReplicaSets(rkd.replicaSetLister)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.K8sResource, len(rs))
	for i, r := range rs {
		r2 := r // grab a value copy of "r" to enforce block scope heap reference, and avoid shadowing non-block scope "r"
		replicaSet := &models.K8sResource{
			Data: &r2,
		}
		ret[i] = replicaSet
	}
	return ret, nil
}

// ReplicationControllers method for runningK8sData
func (rkd *runningK8sData) ReplicationControllers() ([]*models.K8sResource, error) {
	rcs, err := GetReplicationControllers(rkd.rcLister)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.K8sResource, len(rcs))
	for i, rc := range rcs {
		rc2 := rc // grab a value copy of "rc" to enforce block scope heap reference, and avoid shadowing non-block scope "rc"
		replicationController := &models.K8sResource{
			Data: &rc2,
		}
		ret[i] = replicationController
	}
	return ret, nil
}

// Services method for runningK8sData
func (rkd *runningK8sData) Services() ([]*models.K8sResource, error) {
	services, err := getServices(rkd.serviceLister)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.K8sResource, len(services))
	for i, s := range services {
		s2 := s // grab a value copy of "s" to enforce block scope heap reference, and avoid shadowing non-block scope "s"
		service := &models.K8sResource{
			Data: &s2,
		}
		ret[i] = service
	}
	return ret, nil
}

// GetNodesModels gets k8s node model data for RESTful consumption
func GetNodesModels(k RunningK8sData) ([]*models.K8sResource, error) {
	nodes, err := k.Nodes()
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetPodsModels gets k8s pod model data for RESTful consumption
func GetPodsModels(k RunningK8sData) ([]*models.K8sResource, error) {
	pods, err := k.Pods()
	if err != nil {
		return nil, err
	}
	return pods, nil
}

// GetDeploymentsModels gets k8s deployment model data for RESTful consumption
func GetDeploymentsModels(k RunningK8sData) ([]*models.K8sResource, error) {
	ds, err := k.Deployments()
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// GetEventsModels gets k8s event model data for RESTful consumption
func GetEventsModels(k RunningK8sData) ([]*models.K8sResource, error) {
	events, err := k.Events()
	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetDaemonSetsModels gets k8s daemonSet model data for RESTful consumption
func GetDaemonSetsModels(k RunningK8sData) ([]*models.K8sResource, error) {
	ds, err := k.DaemonSets()
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// GetReplicaSetsModels gets k8s replicaSet model data for RESTful consumption
func GetReplicaSetsModels(k RunningK8sData) ([]*models.K8sResource, error) {
	rs, err := k.ReplicaSets()
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// GetReplicationControllersModels gets k8s rc model data for RESTful consumption
func GetReplicationControllersModels(k RunningK8sData) ([]*models.K8sResource, error) {
	rcs, err := k.ReplicationControllers()
	if err != nil {
		return nil, err
	}
	return rcs, nil
}

// GetServicesModels gets k8s pod model data for RESTful consumption
func GetServicesModels(k RunningK8sData) ([]*models.K8sResource, error) {
	services, err := k.Services()
	if err != nil {
		return []*models.K8sResource{}, err
	}
	return services, nil
}

// GetReplicationControllers is a helper function that returns a slice of
// ReplicationController objects given a rc.Lister interface
func GetReplicationControllers(rcLister rc.Lister) ([]api.ReplicationController, error) {
	rcs, err := rcLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []api.ReplicationController{}, err
	}
	return rcs.Items, nil
}

// getNodes is a helper function that returns a slice of
// Node objects given a node.Lister interface
func getNodes(nodeLister node.Lister) ([]api.Node, error) {
	nodes, err := nodeLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []api.Node{}, err
	}
	return nodes.Items, nil
}

// getPods is a helper function that returns a slice of
// Pod objects given a pod.Lister interface
func getPods(podLister pod.Lister) ([]api.Pod, error) {
	pods, err := podLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []api.Pod{}, err
	}
	return pods.Items, nil
}

// getDaemonSets is a helper function that returns a slice of
// DaemonSet objects given a daemonset.Lister interface
func getDaemonSets(dsLister daemonset.Lister) ([]extensions.DaemonSet, error) {
	daemonSets, err := dsLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []extensions.DaemonSet{}, err
	}
	return daemonSets.Items, nil
}

// getDeployments is a helper function that returns a slice of
// Deployment objects given a deployment.Lister interface
func getDeployments(dLister deployment.Lister) ([]extensions.Deployment, error) {
	deployments, err := dLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []extensions.Deployment{}, err
	}
	return deployments.Items, nil
}

// getEvents is a helper function that returns a slice of
// Event objects given an event.Lister interface
func getEvents(eLister event.Lister) ([]api.Event, error) {
	events, err := eLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []api.Event{}, err
	}
	return events.Items, nil
}

// getReplicaSets is a helper function that returns a slice of
// ReplicaSet objects given a replicaset.Lister interface
func getReplicaSets(rsLister replicaset.Lister) ([]extensions.ReplicaSet, error) {
	replicaSets, err := rsLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []extensions.ReplicaSet{}, err
	}
	return replicaSets.Items, nil
}

// getServices is a helper function that returns a slice of
// Service objects given a service.Lister interface
func getServices(serviceLister service.Lister) ([]api.Service, error) {
	services, err := serviceLister.List(api.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		return []api.Service{}, err
	}
	return services.Items, nil
}
