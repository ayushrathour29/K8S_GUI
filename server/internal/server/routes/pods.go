package routes

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/metrics/pkg/client/clientset/versioned"
	"k8_gui/internal/api"
)

func RegisterPodRoutes(r *mux.Router, clientset *kubernetes.Clientset) {
	r.HandleFunc("/pods", api.ListPods(clientset)).Methods("GET")
	r.HandleFunc("/pods/{namespace}/{name}", api.GetPod(clientset)).Methods("GET")
	r.HandleFunc("/pods/{namespace}/{name}", api.DeletePod(clientset)).Methods("DELETE")
	r.HandleFunc("/pods/{namespace}/{name}/logs", api.GetPodLogs(clientset)).Methods("GET")
}
