package routes

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8_gui/internal/api"
)

func RegisterServiceRoutes(r *mux.Router, clientset *kubernetes.Clientset) {
	r.HandleFunc("/services", api.ListServices(clientset)).Methods("GET")
	r.HandleFunc("/services/{namespace}/{name}", api.GetService(clientset)).Methods("GET")
	r.HandleFunc("/services", api.CreateService(clientset)).Methods("POST")
	r.HandleFunc("/services/{namespace}/{name}", api.DeleteService(clientset)).Methods("DELETE")
}
