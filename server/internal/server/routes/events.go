package routes

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8_gui/internal/api"
)

func RegisterEventRoutes(r *mux.Router, clientset *kubernetes.Clientset) {
	r.HandleFunc("/events", api.ListEvents(clientset)).Methods("GET")
	r.HandleFunc("/events/{namespace}", api.ListEventsByNamespace(clientset)).Methods("GET")
}
