package server

import (
	"encoding/json"
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
	router.HandleFunc("/api/validate-token", auth.VerifyToken).Methods("GET")

	// Only add protected routes if Kubernetes client is available
	if clientset != nil {
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

		if metricsClient != nil {
			protected.HandleFunc("/nodes/{name}/metrics", api.GetNodeMetrics(clientset, metricsClient)).Methods("GET")
			protected.HandleFunc("/metrics/nodes", api.GetNodesMetrics(metricsClient)).Methods("GET")
			protected.HandleFunc("/metrics/pods", api.GetPodsMetrics(metricsClient)).Methods("GET")
			protected.HandleFunc("/metrics/pods/{namespace}", api.GetPodMetricsByNamespace(metricsClient)).Methods("GET")
		}

		protected.HandleFunc("/events", api.ListEvents(clientset)).Methods("GET")
		protected.HandleFunc("/events/{namespace}", api.ListEventsByNamespace(clientset)).Methods("GET")
	} else {
		// Add mock data endpoints when Kubernetes is not available
		protected := router.PathPrefix("/api").Subrouter()
		protected.Use(auth.ValidateJWTMiddleware)

		// Mock pods endpoint
		protected.HandleFunc("/pods", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"name":      "nginx-pod",
						"namespace": "default",
						"status":    "Running",
					},
					{
						"name":      "redis-pod",
						"namespace": "default",
						"status":    "Running",
					},
					{
						"name":      "mysql-pod",
						"namespace": "database",
						"status":    "Running",
					},
				},
			})
		}).Methods("GET")

		// Mock deployments endpoint
		protected.HandleFunc("/deployments", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"name":      "nginx-deployment",
						"namespace": "default",
						"status":    "Available",
					},
					{
						"name":      "api-deployment",
						"namespace": "backend",
						"status":    "Available",
					},
				},
			})
		}).Methods("GET")

		// Mock services endpoint
		protected.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"name":      "nginx-service",
						"namespace": "default",
						"status":    "Active",
					},
					{
						"name":      "api-service",
						"namespace": "backend",
						"status":    "Active",
					},
				},
			})
		}).Methods("GET")

		// Mock nodes endpoint
		protected.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"name":   "node-1",
						"status": "Ready",
					},
					{
						"name":   "node-2",
						"status": "Ready",
					},
				},
			})
		}).Methods("GET")

		// Mock events endpoint
		protected.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"name":          "pod-created",
						"namespace":     "default",
						"reason":        "Created",
						"message":       "Pod nginx-pod created successfully",
						"type":          "Normal",
						"lastTimestamp": "2024-01-01T10:00:00Z",
					},
				},
			})
		}).Methods("GET")

		// Add a simple health check endpoint when Kubernetes is not available
		router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok", "message": "Server running (Kubernetes not available)"}`))
		}).Methods("GET")
	}

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(router)
}
