package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Event represents a simplified event view
type Event struct {
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	Type           string `json:"type"`
	InvolvedObject string `json:"involvedObject"`
	FirstTimestamp string `json:"firstTimestamp"`
	LastTimestamp  string `json:"lastTimestamp"`
	Count          int32  `json:"count"`
}

// EventListResponse represents event list response
type EventListResponse struct {
	Items []Event `json:"items"`
}

// ListEvents returns all events
func ListEvents(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := clientset.CoreV1().Events("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list events: %v", err)
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}

		response := EventListResponse{Items: make([]Event, 0, len(events.Items))}
		for _, e := range events.Items {
			response.Items = append(response.Items, Event{
				Name:           e.Name,
				Namespace:      e.Namespace,
				Reason:         e.Reason,
				Message:        e.Message,
				Type:           e.Type,
				InvolvedObject: e.InvolvedObject.Kind + "/" + e.InvolvedObject.Name,
				FirstTimestamp: e.FirstTimestamp.Time.Format(time.RFC3339),
				LastTimestamp:  e.LastTimestamp.Time.Format(time.RFC3339),
				Count:          e.Count,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// ListEventsByNamespace returns events for a specific namespace
func ListEventsByNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]

		events, err := clientset.CoreV1().Events(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list events for namespace %s: %v", namespace, err)
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}

		response := EventListResponse{Items: make([]Event, 0, len(events.Items))}
		for _, e := range events.Items {
			response.Items = append(response.Items, Event{
				Name:           e.Name,
				Namespace:      e.Namespace,
				Reason:         e.Reason,
				Message:        e.Message,
				Type:           e.Type,
				InvolvedObject: e.InvolvedObject.Kind + "/" + e.InvolvedObject.Name,
				FirstTimestamp: e.FirstTimestamp.Time.Format(time.RFC3339),
				LastTimestamp:  e.LastTimestamp.Time.Format(time.RFC3339),
				Count:          e.Count,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
