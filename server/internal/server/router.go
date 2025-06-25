package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"k8s.io/client-go/kubernetes"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"

	"k8_gui/internal/auth"
	"k8_gui/internal/server/routes"
	"k8_gui/internal/utils"
)

// NewRouter creates application router
func NewRouter(clientset *kubernetes.Clientset, metricsClient *metricsclientset.Clientset) http.Handler {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/api/login", auth.HandleLogin).Methods("POST")
	router.HandleFunc("/api/validate-token", auth.VerifyToken).Methods("GET")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(auth.ValidateJWTMiddleware)

	if clientset != nil {
		// Register grouped routes
		routes.RegisterPodRoutes(protected, clientset)
		routes.RegisterDeploymentRoutes(protected, clientset)
		routes.RegisterServiceRoutes(protected, clientset)
		routes.RegisterEventRoutes(protected, clientset)
		routes.RegisterNodeRoutes(protected, clientset, metricsClient)
		routes.RegisterNamspaceRoutes(protected, clientset)

		if metricsClient != nil {
			routes.RegisterMetricsRoutes(protected, metricsClient)
		}
	} else {
		// Add mock routes if Kubernetes is unavailable
		protected.HandleFunc("/pods", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{"name": "nginx-pod", "namespace": "default", "status": "Running"},
					{"name": "redis-pod", "namespace": "default", "status": "Running"},
					{"name": "mysql-pod", "namespace": "database", "status": "Running"},
				},
			})
		}).Methods("GET")

		protected.HandleFunc("/deployments", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{"name": "nginx-deployment", "namespace": "default", "status": "Available"},
					{"name": "api-deployment", "namespace": "backend", "status": "Available"},
				},
			})
		}).Methods("GET")

		protected.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{"name": "nginx-service", "namespace": "default", "status": "Active"},
					{"name": "api-service", "namespace": "backend", "status": "Active"},
				},
			})
		}).Methods("GET")

		protected.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{"name": "node-1", "status": "Ready"},
					{"name": "node-2", "status": "Ready"},
				},
			})
		}).Methods("GET")

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
						"lastTimestamp": time.Now().Format(time.RFC3339),
					},
				},
			})
		}).Methods("GET")

		router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(utils.MsgServerRunningLimited))
		}).Methods("GET")
	}

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(router)
}
