package routes

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8_gui/internal/api"
)

func RegisterDeploymentRoutes(r *mux.Router, clientset *kubernetes.Clientset) {
	r.HandleFunc("/deployments", api.ListDeployments(clientset)).Methods("GET")
	r.HandleFunc("/deployments/{namespace}/{name}", api.GetDeployment(clientset)).Methods("GET")
	r.HandleFunc("/deployments/{namespace}/{name}", api.DeleteDeployment(clientset)).Methods("DELETE")
	r.HandleFunc("/deployments", api.CreateDeployment(clientset)).Methods("POST")
	r.HandleFunc("/deployments/{namespace}/{name}", api.UpdateDeployment(clientset)).Methods("PUT")
}
