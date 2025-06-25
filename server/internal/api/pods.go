package api

import (
	"encoding/json"
	"fmt"
	"k8_gui/internal/models"
	"k8_gui/internal/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListPods returns all pods
func ListPods(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pods, err := clientset.CoreV1().Pods("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedListPods, err)
			http.Error(w, utils.MsgFailedListPods, http.StatusInternalServerError)
			return
		}

		response := models.PodListResponse{Items: make([]models.Pod, 0, len(pods.Items))}
		for _, p := range pods.Items {
			restartCount := int32(0)
			for _, cs := range p.Status.ContainerStatuses {
				restartCount += cs.RestartCount
			}

			containers := make([]string, len(p.Spec.Containers))
			for i, c := range p.Spec.Containers {
				containers[i] = c.Name
			}

			response.Items = append(response.Items, models.Pod{
				Name:         p.Name,
				Namespace:    p.Namespace,
				Status:       string(p.Status.Phase),
				RestartCount: restartCount,
				CreatedAt:    p.CreationTimestamp.Time.Format(time.RFC3339),
				NodeName:     p.Spec.NodeName,
				PodIP:        p.Status.PodIP,
				Containers:   containers,
				Labels:       p.Labels,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodePodsList, err)
		}
	}
}

// GetPod returns pod details
func GetPod(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		pod, err := clientset.CoreV1().Pods(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf(utils.LogFailedGetPod, err)
			http.Error(w, utils.MsgPodNotFound, http.StatusNotFound)
			return
		}

		restartCount := int32(0)
		for _, cs := range pod.Status.ContainerStatuses {
			restartCount += cs.RestartCount
		}

		containers := make([]string, len(pod.Spec.Containers))
		for i, c := range pod.Spec.Containers {
			containers[i] = c.Name
		}

		response := models.Pod{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Status:       string(pod.Status.Phase),
			RestartCount: restartCount,
			CreatedAt:    pod.CreationTimestamp.Time.Format(time.RFC3339),
			NodeName:     pod.Spec.NodeName,
			PodIP:        pod.Status.PodIP,
			Containers:   containers,
			Labels:       pod.Labels,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodePod, err)
		}
	}
}

// DeletePod deletes a pod
func DeletePod(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		err := clientset.CoreV1().Pods(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf(utils.LogFailedDeletePod, err)
			http.Error(w, utils.MsgFailedDeletePod, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetPodLogs returns pod logs
func GetPodLogs(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		// Get query parameters for log options
		tailLines := int64(100) // Default to last 100 lines
		if tail := r.URL.Query().Get("tail"); tail != "" {
			parsed, err := parseTailLines(tail)
			if err != nil {
				http.Error(w, utils.MsgInvalidTailParameter, http.StatusBadRequest)
				return
			}
			tailLines = parsed
		}

		logs, err := clientset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{
			TailLines: &tailLines,
		}).Do(r.Context()).Raw()
		if err != nil {
			log.Printf(utils.LogFailedGetPodLogs, err)
			http.Error(w, utils.MsgFailedGetPodLogs, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write(logs); err != nil {
			log.Printf(utils.LogFailedWritePodLogs, err)
		}
	}
}

// Helper function to parse tail lines parameter
func parseTailLines(tail string) (int64, error) {
	if tail == "" {
		return 0, nil
	}
	result, err := strconv.ParseInt(tail, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid value for tail: %s", tail)
	}
	return result, nil
}
