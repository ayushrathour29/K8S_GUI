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

// ListNodes returns all nodes
func ListNodes(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodes, err := clientset.CoreV1().Nodes().List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedListNodes, err)
			http.Error(w, utils.MsgFailedListNodes, http.StatusInternalServerError)
			return
		}

		response := models.NodeListResponse{Items: make([]models.Node, 0, len(nodes.Items))}
		for _, n := range nodes.Items {
			capacity := make(map[string]string)
			allocatable := make(map[string]string)

			for resourceName, quantity := range n.Status.Capacity {
				capacity[string(resourceName)] = quantity.String()
			}
			for resourceName, quantity := range n.Status.Allocatable {
				allocatable[string(resourceName)] = quantity.String()
			}

			response.Items = append(response.Items, models.Node{
				Name:        n.Name,
				Status:      string(n.Status.Phase),
				Version:     n.Status.NodeInfo.KubeletVersion,
				OSImage:     n.Status.NodeInfo.OSImage,
				Capacity:    capacity,
				Allocatable: allocatable,
				CreatedAt:   n.CreationTimestamp.Time.Format(time.RFC3339),
				Labels:      n.Labels,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodeNodesList, err)
		}
	}
}

// GetNode returns node details
func GetNode(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		node, err := clientset.CoreV1().Nodes().Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf(utils.LogFailedGetNode, err)
			http.Error(w, utils.MsgNodeNotFound, http.StatusNotFound)
			return
		}

		capacity := make(map[string]string)
		allocatable := make(map[string]string)

		for resourceName, quantity := range node.Status.Capacity {
			capacity[string(resourceName)] = quantity.String()
		}
		for resourceName, quantity := range node.Status.Allocatable {
			allocatable[string(resourceName)] = quantity.String()
		}

		response := models.Node{
			Name:        node.Name,
			Status:      string(node.Status.Phase),
			Version:     node.Status.NodeInfo.KubeletVersion,
			OSImage:     node.Status.NodeInfo.OSImage,
			Capacity:    capacity,
			Allocatable: allocatable,
			CreatedAt:   node.CreationTimestamp.Time.Format(time.RFC3339),
			Labels:      node.Labels,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodeNode, err)
		}
	}
}
