package routes

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	"k8_gui/internal/api"
)

func RegisterNodeRoutes(r *mux.Router, clientset *kubernetes.Clientset, metricsClient *metricsclientset.Clientset) {
	r.HandleFunc("/nodes", api.ListNodes(clientset)).Methods("GET")
	r.HandleFunc("/nodes/{name}", api.GetNode(clientset)).Methods("GET")

	if metricsClient != nil {
		r.HandleFunc("/nodes/{name}/metrics", api.GetNodeMetrics(clientset, metricsClient)).Methods("GET")
	}
}
