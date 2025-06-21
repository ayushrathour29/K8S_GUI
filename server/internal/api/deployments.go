package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// Deployment represents a simplified deployment view
type Deployment struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Replicas          int32             `json:"replicas"`
	AvailableReplicas int32             `json:"availableReplicas"`
	CreatedAt         string            `json:"createdAt"`
	Strategy          string            `json:"strategy"`
	Labels            map[string]string `json:"labels,omitempty"`
}

// DeploymentListResponse represents deployment list response
type DeploymentListResponse struct {
	Items []Deployment `json:"items"`
}

// ListDeployments returns all deployments
func ListDeployments(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deployments, err := clientset.AppsV1().Deployments("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list deployments: %v", err)
			http.Error(w, "Failed to list deployments", http.StatusInternalServerError)
			return
		}

		response := DeploymentListResponse{Items: make([]Deployment, 0, len(deployments.Items))}
		for _, d := range deployments.Items {
			response.Items = append(response.Items, Deployment{
				Name:              d.Name,
				Namespace:         d.Namespace,
				Replicas:          *d.Spec.Replicas,
				AvailableReplicas: d.Status.AvailableReplicas,
				CreatedAt:         d.CreationTimestamp.Time.Format(time.RFC3339),
				Strategy:          string(d.Spec.Strategy.Type),
				Labels:            d.Labels,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetDeployment returns deployment details
func GetDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		deployment, err := clientset.AppsV1().Deployments(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get deployment: %v", err)
			http.Error(w, "Deployment not found", http.StatusNotFound)
			return
		}

		response := Deployment{
			Name:              deployment.Name,
			Namespace:         deployment.Namespace,
			Replicas:          *deployment.Spec.Replicas,
			AvailableReplicas: deployment.Status.AvailableReplicas,
			CreatedAt:         deployment.CreationTimestamp.Time.Format(time.RFC3339),
			Strategy:          string(deployment.Spec.Strategy.Type),
			Labels:            deployment.Labels,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateDeployment creates a new deployment
func CreateDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
			Image     string `json:"image"`
			Replicas  int32  `json:"replicas"`
			Port      int32  `json:"port"`
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &req.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": req.Name,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": req.Name,
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  req.Name,
								Image: req.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: req.Port,
									},
								},
							},
						},
					},
				},
			},
		}

		created, err := clientset.AppsV1().Deployments(req.Namespace).Create(r.Context(), deployment, metav1.CreateOptions{})
		if err != nil {
			log.Printf("Failed to create deployment: %v", err)
			http.Error(w, "Failed to create deployment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	}
}

// UpdateDeployment updates an existing deployment
func UpdateDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		type Request struct {
			Image    string `json:"image"`
			Replicas int32  `json:"replicas"`
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		deployment, err := clientset.AppsV1().Deployments(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get deployment: %v", err)
			http.Error(w, "Deployment not found", http.StatusNotFound)
			return
		}

		// Update image and replicas
		if req.Image != "" {
			deployment.Spec.Template.Spec.Containers[0].Image = req.Image
		}
		if req.Replicas > 0 {
			deployment.Spec.Replicas = &req.Replicas
		}

		updated, err := clientset.AppsV1().Deployments(namespace).Update(r.Context(), deployment, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("Failed to update deployment: %v", err)
			http.Error(w, "Failed to update deployment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	}
}

// DeleteDeployment deletes a deployment
func DeleteDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		err := clientset.AppsV1().Deployments(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete deployment: %v", err)
			http.Error(w, "Failed to delete deployment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Helper function to convert int32 to IntOrString
func intstrFromInt(i int32) intstr.IntOrString {
	return intstr.FromInt(int(i))
}
