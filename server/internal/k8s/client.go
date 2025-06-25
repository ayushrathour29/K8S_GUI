package k8s

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

// InitK8sClient initializes Kubernetes and Metrics clients
func InitK8sClient() (*kubernetes.Clientset, *metricsclientset.Clientset, error) {
	// Try in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return clientset, metricsClient, nil
}


// type K8sClient struct {
// 	Clientset     *kubernetes.Clientset
// 	MetricsClient *metricsclientset.Clientset
// }

// func InitK8sClient() (*K8sClient, error) {
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
// 		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	cs, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mc, err := metricsclientset.NewForConfig(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &K8sClient{
// 		Clientset:     cs,
// 		MetricsClient: mc,
// 	}, nil
// }