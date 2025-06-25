package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodService defines pod-related methods
type PodService interface {
	ListPods(namespace string) ([]corev1.Pod, error)
	GetPod(namespace, name string) (*corev1.Pod, error)
	DeletePod(namespace, name string) error
}

// DeploymentService defines deployment-related methods
type DeploymentService interface {
	ListDeployments(namespace string) ([]appsv1.Deployment, error)
	GetDeployment(namespace, name string) (*appsv1.Deployment, error)
	DeleteDeployment(namespace, name string) error
}

// ClusterService defines cluster-related methods
// (e.g., get cluster info, health, version)
type ClusterService interface {
	GetClusterInfo() (interface{}, error)
	GetClusterHealth() (interface{}, error)
	GetClusterVersion() (interface{}, error)
}

// EventService defines event-related methods
type EventService interface {
	ListEvents(namespace string) ([]corev1.Event, error)
	GetEvent(namespace, name string) (*corev1.Event, error)
	DeleteEvent(namespace, name string) error
}

// MetricsService defines metrics-related methods
type MetricsService interface {
	GetNodeMetrics(name string) (interface{}, error)
	GetNodesMetrics() (interface{}, error)
	GetPodsMetrics() (interface{}, error)
	GetPodMetricsByNamespace(namespace string) (interface{}, error)
}

// NamespaceService defines namespace-related methods
type NamespaceService interface {
	ListNamespaces() ([]corev1.Namespace, error)
	GetNamespace(name string) (*corev1.Namespace, error)
	CreateNamespace(ns *corev1.Namespace) (*corev1.Namespace, error)
	DeleteNamespace(name string) error
}

// NodeService defines node-related methods
type NodeService interface {
	ListNodes() ([]corev1.Node, error)
	GetNode(name string) (*corev1.Node, error)
	DeleteNode(name string) error
}

// ServiceService defines service-related methods
type ServiceService interface {
	ListServices(namespace string) ([]corev1.Service, error)
	GetService(namespace, name string) (*corev1.Service, error)
	CreateService(namespace string, service *corev1.Service) (*corev1.Service, error)
	DeleteService(namespace, name string) error
}

// You can add more interfaces later: ServiceService, NodeService, etc.
