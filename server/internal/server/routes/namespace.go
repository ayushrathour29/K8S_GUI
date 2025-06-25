	package routes

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	// metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	"k8_gui/internal/api"
)

func RegisterNamspaceRoutes(r *mux.Router, clientset *kubernetes.Clientset) {
	  r.HandleFunc("/namespaces", api.ListNamespaces(clientset)).Methods("GET")
		r.HandleFunc("/namespaces/{name}", api.GetNamespace(clientset)).Methods("GET")
		r.HandleFunc("/namespaces", api.CreateNamespace(clientset)).Methods("POST")
		r.HandleFunc("/namespaces/{name}", api.DeleteNamespace(clientset)).Methods("DELETE")
}	
		
		
		
		
	  