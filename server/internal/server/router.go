package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"k8s.io/client-go/kubernetes"

	"k8_gui/internal/api"
	"k8_gui/internal/auth"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

// NewRouter creates application router
func NewRouter(clientset *kubernetes.Clientset, metricsClient *metricsclientset.Clientset) http.Handler {
	router := mux.NewRouter()

	//  Public routes
	router.HandleFunc("/api/login", auth.HandleLogin).Methods("POST")
	router.HandleFunc("/api/verify", api.VerifyToken).Methods("GET")

	// Protected subrouter
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(auth.ValidateJWTMiddleware)

	// Cluster endpoints
	protected.HandleFunc("/clusters", api.GetClusters(clientset)).Methods("GET")

	// Pod endpoints
	protected.HandleFunc("/pods", api.ListPods(clientset)).Methods("GET")
	protected.HandleFunc("/pods/{namespace}/{name}", api.GetPod(clientset)).Methods("GET")
	protected.HandleFunc("/pods/{namespace}/{name}", api.DeletePod(clientset)).Methods("DELETE")
	protected.HandleFunc("/pods/{namespace}/{name}/logs", api.GetPodLogs(clientset)).Methods("GET")

	// Deployment endpoints
	protected.HandleFunc("/deployments", api.ListDeployments(clientset)).Methods("GET")
	protected.HandleFunc("/deployments/{namespace}/{name}", api.GetDeployment(clientset)).Methods("GET")
	protected.HandleFunc("/deployments/{namespace}/{name}", api.DeleteDeployment(clientset)).Methods("DELETE")
	protected.HandleFunc("/deployments", api.CreateDeployment(clientset)).Methods("POST")
	protected.HandleFunc("/deployments/{namespace}/{name}", api.UpdateDeployment(clientset)).Methods("PUT")

	// Service endpoints
	protected.HandleFunc("/services", api.ListServices(clientset)).Methods("GET")
	protected.HandleFunc("/services/{namespace}/{name}", api.GetService(clientset)).Methods("GET")
	protected.HandleFunc("/services", api.CreateService(clientset)).Methods("POST")
	protected.HandleFunc("/services/{namespace}/{name}", api.DeleteService(clientset)).Methods("DELETE")

	protected.HandleFunc("/namespaces", api.ListNamespaces(clientset)).Methods("GET")
	protected.HandleFunc("/namespaces/{name}", api.GetNamespace(clientset)).Methods("GET")
	protected.HandleFunc("/namespaces", api.CreateNamespace(clientset)).Methods("POST")
	protected.HandleFunc("/namespaces/{name}", api.DeleteNamespace(clientset)).Methods("DELETE")

	protected.HandleFunc("/nodes", api.ListNodes(clientset)).Methods("GET")
	protected.HandleFunc("/nodes/{name}", api.GetNode(clientset)).Methods("GET")
	protected.HandleFunc("/nodes/{name}/metrics", api.GetNodeMetrics(clientset, metricsClient)).Methods("GET")
	protected.HandleFunc("/metrics/nodes", api.GetNodesMetrics(metricsClient)).Methods("GET")
	protected.HandleFunc("/metrics/pods", api.GetPodsMetrics(metricsClient)).Methods("GET")
	protected.HandleFunc("/metrics/pods/{namespace}", api.GetPodMetricsByNamespace(metricsClient)).Methods("GET")

	protected.HandleFunc("/events", api.ListEvents(clientset)).Methods("GET")
	protected.HandleFunc("/events/{namespace}", api.ListEventsByNamespace(clientset)).Methods("GET")

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(router)
}
