package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"k8s.io/client-go/kubernetes"

	"k8_gui/internal/api"
)

// NewRouter creates and returns a new mux.Router with all routes and middleware configured
func NewRouter(clientset *kubernetes.Clientset) http.Handler {
	router := mux.NewRouter()

	// Auth endpoints
	router.HandleFunc("/api/login", api.HandleLogin).Methods("POST")
	router.HandleFunc("/api/verify", api.VerifyToken).Methods("GET")

	// Cluster endpoints
	router.HandleFunc("/api/clusters", api.GetClusters).Methods("GET")

	// Pod endpoints
	router.HandleFunc("/api/pods", api.ListPods(clientset)).Methods("GET")
	router.HandleFunc("/api/pods/{namespace}/{name}", api.GetPod(clientset)).Methods("GET")
	router.HandleFunc("/api/pods/{namespace}/{name}", api.DeletePod(clientset)).Methods("DELETE")

	// Deployment endpoints
	router.HandleFunc("/api/deployments", api.ListDeployments(clientset)).Methods("GET")
	router.HandleFunc("/api/deployments/{namespace}/{name}", api.GetDeployment(clientset)).Methods("GET")
	router.HandleFunc("/api/deployments", api.CreateDeployment(clientset)).Methods("POST")
	router.HandleFunc("/api/deployments/{namespace}/{name}", api.UpdateDeployment(clientset)).Methods("PUT")
	router.HandleFunc("/api/deployments/{namespace}/{name}", api.DeleteDeployment(clientset)).Methods("DELETE")

	// Service endpoints
	router.HandleFunc("/api/services", api.ListServices(clientset)).Methods("GET")
	router.HandleFunc("/api/services/{namespace}/{name}", api.GetService(clientset)).Methods("GET")
	router.HandleFunc("/api/services", api.CreateService(clientset)).Methods("POST")
	router.HandleFunc("/api/services/{namespace}/{name}", api.DeleteService(clientset)).Methods("DELETE")

	// CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(router)
}
