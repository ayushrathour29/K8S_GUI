package api

import (
	"encoding/json"
	"k8_gui/internal/models"
	"k8_gui/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListEvents returns all events
func ListEvents(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := clientset.CoreV1().Events("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedListEvents, err)
			http.Error(w, utils.MsgFailedListEvents, http.StatusInternalServerError)
			return
		}

		response := models.EventListResponse{Items: make([]models.Event, 0, len(events.Items))}
		for _, e := range events.Items {
			response.Items = append(response.Items, models.Event{
				Name:           e.Name,
				Namespace:      e.Namespace,
				Reason:         e.Reason,
				Message:        e.Message,
				Type:           e.Type,
				InvolvedObject: e.InvolvedObject.Kind + "/" + e.InvolvedObject.Name,
				FirstTimestamp: e.FirstTimestamp.Time.Format(time.RFC3339),
				LastTimestamp: func() string {
					if !e.LastTimestamp.IsZero() {
						return e.LastTimestamp.Time.Format(time.RFC3339)
					}
					if !e.FirstTimestamp.IsZero() {
						return e.FirstTimestamp.Time.Format(time.RFC3339)
					}
					return ""
				}(),
				Count: e.Count,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodeEventsList, err)
		}
	}
}

// ListEventsByNamespace returns events for a specific namespace
func ListEventsByNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]

		events, err := clientset.CoreV1().Events(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedListEventsNamespace, namespace, err)
			http.Error(w, utils.MsgFailedListEvents, http.StatusInternalServerError)
			return
		}

		response := models.EventListResponse{Items: make([]models.Event, 0, len(events.Items))}
		for _, e := range events.Items {
			response.Items = append(response.Items, models.Event{
				Name:           e.Name,
				Namespace:      e.Namespace,
				Reason:         e.Reason,
				Message:        e.Message,
				Type:           e.Type,
				InvolvedObject: e.InvolvedObject.Kind + "/" + e.InvolvedObject.Name,
				FirstTimestamp: e.FirstTimestamp.Time.Format(time.RFC3339),
				LastTimestamp: func() string {
					if !e.LastTimestamp.IsZero() {
						return e.LastTimestamp.Time.Format(time.RFC3339)
					}
					if !e.FirstTimestamp.IsZero() {
						return e.FirstTimestamp.Time.Format(time.RFC3339)
					}
					return ""
				}(),
				Count: e.Count,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodeEventsListNamespace, err)
		}
	}
}
