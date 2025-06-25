package routes

import (
	"github.com/gorilla/mux"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	"k8_gui/internal/api"
)

func RegisterMetricsRoutes(r *mux.Router, metricsClient *metricsclientset.Clientset) {
	r.HandleFunc("/metrics/nodes", api.GetNodesMetrics(metricsClient)).Methods("GET")
	r.HandleFunc("/metrics/pods", api.GetPodsMetrics(metricsClient)).Methods("GET")
	r.HandleFunc("/metrics/pods/{namespace}", api.GetPodMetricsByNamespace(metricsClient)).Methods("GET")
}
