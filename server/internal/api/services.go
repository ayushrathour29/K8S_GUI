package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Service represents a simplified service view
type Service struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Type      string        `json:"type"`
	ClusterIP string        `json:"clusterIP"`
	Ports     []ServicePort `json:"ports"`
	CreatedAt string        `json:"createdAt"`
}

// ServicePort represents service port information
type ServicePort struct {
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

// ServiceListResponse represents service list response
type ServiceListResponse struct {
	Items []Service `json:"items"`
}

// ListServices returns all services
func ListServices(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := clientset.CoreV1().Services("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list services: %v", err)
			http.Error(w, "Failed to list services", http.StatusInternalServerError)
			return
		}

		response := ServiceListResponse{Items: make([]Service, 0, len(services.Items))}
		for _, s := range services.Items {
			ports := make([]ServicePort, len(s.Spec.Ports))
			for i, p := range s.Spec.Ports {
				ports[i] = ServicePort{
					Port:     p.Port,
					Protocol: string(p.Protocol),
				}
			}

			response.Items = append(response.Items, Service{
				Name:      s.Name,
				Namespace: s.Namespace,
				Type:      string(s.Spec.Type),
				ClusterIP: s.Spec.ClusterIP,
				Ports:     ports,
				CreatedAt: s.CreationTimestamp.Time.Format(time.RFC3339),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetService returns service details
func GetService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		service, err := clientset.CoreV1().Services(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get service: %v", err)
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		ports := make([]ServicePort, len(service.Spec.Ports))
		for i, p := range service.Spec.Ports {
			ports[i] = ServicePort{
				Port:     p.Port,
				Protocol: string(p.Protocol),
			}
		}

		response := Service{
			Name:      service.Name,
			Namespace: service.Namespace,
			Type:      string(service.Spec.Type),
			ClusterIP: service.Spec.ClusterIP,
			Ports:     ports,
			CreatedAt: service.CreationTimestamp.Time.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateService creates a new service
func CreateService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Name       string `json:"name"`
			Namespace  string `json:"namespace"`
			Port       int32  `json:"port"`
			TargetPort int32  `json:"targetPort"`
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": req.Name,
				},
				Ports: []corev1.ServicePort{
					{
						Port:       req.Port,
						TargetPort: intstrFromInt(req.TargetPort),
						Protocol:   corev1.ProtocolTCP,
					},
				},
				Type: corev1.ServiceTypeClusterIP,
			},
		}

		created, err := clientset.CoreV1().Services(req.Namespace).Create(r.Context(), service, metav1.CreateOptions{})
		if err != nil {
			log.Printf("Failed to create service: %v", err)
			http.Error(w, "Failed to create service", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	}
}

// DeleteService deletes a service
func DeleteService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		err := clientset.CoreV1().Services(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete service: %v", err)
			http.Error(w, "Failed to delete service", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
