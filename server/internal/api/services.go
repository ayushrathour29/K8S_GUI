package api

import (
	"encoding/json"
	"k8_gui/internal/models"
	"k8_gui/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListServices returns all services
func ListServices(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := clientset.CoreV1().Services("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedListServices, err)
			http.Error(w, utils.MsgFailedListServices, http.StatusInternalServerError)
			return
		}

		response := models.ServiceListResponse{Items: make([]models.Service, 0, len(services.Items))}
		for _, s := range services.Items {
			ports := make([]models.ServicePort, len(s.Spec.Ports))
			for i, p := range s.Spec.Ports {
				ports[i] = models.ServicePort{
					Port:     p.Port,
					Protocol: string(p.Protocol),
				}
			}

			response.Items = append(response.Items, models.Service{
				Name:      s.Name,
				Namespace: s.Namespace,
				Type:      string(s.Spec.Type),
				ClusterIP: s.Spec.ClusterIP,
				Ports:     ports,
				CreatedAt: s.CreationTimestamp.Time.Format(time.RFC3339),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodeServicesList, err)
		}
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
			log.Printf(utils.LogFailedGetService, err)
			http.Error(w, utils.MsgServiceNotFound, http.StatusNotFound)
			return
		}

		ports := make([]models.ServicePort, len(service.Spec.Ports))
		for i, p := range service.Spec.Ports {
			ports[i] = models.ServicePort{
				Port:     p.Port,
				Protocol: string(p.Protocol),
			}
		}

		response := models.Service{
			Name:      service.Name,
			Namespace: service.Namespace,
			Type:      string(service.Spec.Type),
			ClusterIP: service.Spec.ClusterIP,
			Ports:     ports,
			CreatedAt: service.CreationTimestamp.Time.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodeService, err)
		}
	}
}

// CreateService creates a new service
func CreateService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateServiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, utils.MsgInvalidRequestBody, http.StatusBadRequest)
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
						TargetPort: utils.IntstrFromInt(req.TargetPort),
						Protocol:   corev1.ProtocolTCP,
					},
				},
				Type: corev1.ServiceTypeClusterIP,
			},
		}

		created, err := clientset.CoreV1().Services(req.Namespace).Create(r.Context(), service, metav1.CreateOptions{})
		if err != nil {
			log.Printf(utils.LogFailedCreateService, err)
			http.Error(w, utils.MsgFailedCreateService, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(created); err != nil {
			log.Printf(utils.LogFailedEncodeCreatedService, err)
		}
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
			log.Printf(utils.LogFailedDeleteService, err)
			http.Error(w, utils.MsgFailedDeleteService, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
