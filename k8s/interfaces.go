package k8s

import kcl "k8s.io/kubernetes/pkg/client/unversioned"

// ResourceInterface is an interface for k8s resources
type ResourceInterface interface {
	kcl.DaemonSetsNamespacer
	kcl.DeploymentsNamespacer
	kcl.EventNamespacer
	kcl.NodesInterface
	kcl.PodsNamespacer
	kcl.ReplicaSetsNamespacer
	kcl.ReplicationControllersNamespacer
	kcl.SecretsNamespacer
	kcl.ServicesNamespacer
}

// ResourceInterfaceNamespaced is a "union" of ResourceInterface+namespace
type ResourceInterfaceNamespaced struct {
	ri        ResourceInterface
	namespace string
}

// NewResourceInterfaceNamespaced constructs an instance of ResourceInterfaceNamespaced
func NewResourceInterfaceNamespaced(ri ResourceInterface, ns string) *ResourceInterfaceNamespaced {
	return &ResourceInterfaceNamespaced{ri: ri, namespace: ns}
}

// DaemonSets implementation
func (r *ResourceInterfaceNamespaced) DaemonSets() kcl.DaemonSetInterface {
	return r.ri.DaemonSets(r.namespace)
}

// Deployments implementation
func (r *ResourceInterfaceNamespaced) Deployments() kcl.DeploymentInterface {
	return r.ri.Deployments(r.namespace)
}

// Events implementation
func (r *ResourceInterfaceNamespaced) Events() kcl.EventInterface {
	return r.ri.Events(r.namespace)
}

// Nodes implementation
func (r *ResourceInterfaceNamespaced) Nodes() kcl.NodeInterface {
	return r.ri.Nodes()
}

// Pods implementation
func (r *ResourceInterfaceNamespaced) Pods() kcl.PodInterface {
	return r.ri.Pods(r.namespace)
}

// ReplicaSets implementation
func (r *ResourceInterfaceNamespaced) ReplicaSets() kcl.ReplicaSetInterface {
	return r.ri.ReplicaSets(r.namespace)
}

// ReplicationControllers implementation
func (r *ResourceInterfaceNamespaced) ReplicationControllers() kcl.ReplicationControllerInterface {
	return r.ri.ReplicationControllers(r.namespace)
}

// Services implementation
func (r *ResourceInterfaceNamespaced) Services() kcl.ServiceInterface {
	return r.ri.Services(r.namespace)
}

// Secrets implementation
func (r *ResourceInterfaceNamespaced) Secrets() kcl.SecretsInterface {
	return r.ri.Secrets(r.namespace)
}
