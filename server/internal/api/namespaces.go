package api

import (
	"encoding/json"
	"k8_gui/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListNamespaces returns all namespaces
func ListNamespaces(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespaces, err := clientset.CoreV1().Namespaces().List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list namespaces: %v", err)
			http.Error(w, "Failed to list namespaces", http.StatusInternalServerError)
			return
		}

		response := models.NamespaceListResponse{Items: make([]models.Namespace, 0, len(namespaces.Items))}
		for _, ns := range namespaces.Items {
			response.Items = append(response.Items, models.Namespace{
				Name:      ns.Name,
				Status:    string(ns.Status.Phase),
				CreatedAt: ns.CreationTimestamp.Time.Format(time.RFC3339),
				Labels:    ns.Labels,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode namespaces list: %v", err)
		}
	}
}

// GetNamespace returns namespace details
func GetNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		namespace, err := clientset.CoreV1().Namespaces().Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get namespace: %v", err)
			http.Error(w, "Namespace not found", http.StatusNotFound)
			return
		}

		response := models.Namespace{
			Name:      namespace.Name,
			Status:    string(namespace.Status.Phase),
			CreatedAt: namespace.CreationTimestamp.Time.Format(time.RFC3339),
			Labels:    namespace.Labels,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode namespace: %v", err)
		}
	}
}

// CreateNamespace creates a new namespace
func CreateNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateNamespaceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:   req.Name,
				Labels: req.Labels,
			},
		}

		created, err := clientset.CoreV1().Namespaces().Create(r.Context(), namespace, metav1.CreateOptions{})
		if err != nil {
			log.Printf("Failed to create namespace: %v", err)
			http.Error(w, "Failed to create namespace", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(created); err != nil {
			log.Printf("Failed to encode created namespace: %v", err)
		}
	}
}

// DeleteNamespace deletes a namespace
func DeleteNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		err := clientset.CoreV1().Namespaces().Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete namespace: %v", err)
			http.Error(w, "Failed to delete namespace", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
