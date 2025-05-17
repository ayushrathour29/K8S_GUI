package api

import (
	"encoding/json"
	"k8_gui/internal/auth"
	"net/http"
	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ClusterInfo represents basic information about a K8s cluster
type ClusterInfo struct {
	Name    string `json:"name"`
	Context string `json:"context"`
	Server  string `json:"server"`
}

// HandleLogin authenticates a user and returns a token
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	auth.HandleLogin(w, r)
}

// VerifyToken checks if a token is valid
func VerifyToken(w http.ResponseWriter, r *http.Request) {
	auth.VerifyToken(w, r)
}

// GetClusters returns available Kubernetes clusters
func GetClusters(w http.ResponseWriter, r *http.Request) {
	clusters := []ClusterInfo{
		{
			Name:    "local-cluster",
			Context: "minikube",
			Server:  "https://kubernetes.default.svc",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusters)
}

// ListPods returns all pods in the cluster
func ListPods(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pods, err := clientset.CoreV1().Pods("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, "Failed to list pods: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pods)
	}
}

// GetPod gets details of a specific pod
func GetPod(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		pod, err := clientset.CoreV1().Pods(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			http.Error(w, "Failed to get pod: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pod)
	}
}

// DeletePod deletes a specific pod
func DeletePod(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		err := clientset.CoreV1().Pods(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			http.Error(w, "Failed to delete pod: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ListDeployments returns all deployments
func ListDeployments(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement deployment listing
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Deployment listing not yet implemented"}`))
	}
}

// GetDeployment gets details of a specific deployment
func GetDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement deployment details
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Deployment details not yet implemented"}`))
	}
}

// CreateDeployment creates a new deployment
func CreateDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement deployment creation
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Deployment creation not yet implemented"}`))
	}
}

// UpdateDeployment updates an existing deployment
func UpdateDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement deployment update
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Deployment update not yet implemented"}`))
	}
}

// DeleteDeployment deletes a deployment
func DeleteDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement deployment deletion
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Deployment deletion not yet implemented"}`))
	}
}

// ListServices returns all services
func ListServices(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement service listing
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Service listing not yet implemented"}`))
	}
}

// GetService gets details of a specific service
func GetService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement service details
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Service details not yet implemented"}`))
	}
}

// CreateService creates a new service
func CreateService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement service creation
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Service creation not yet implemented"}`))
	}
}

// DeleteService deletes a service
func DeleteService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement service deletion
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Service deletion not yet implemented"}`))
	}
}
