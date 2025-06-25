package api

import (
	"encoding/json"
	"k8_gui/internal/models"
	"k8_gui/internal/utils"
	"log"
	"net/http"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetClusters returns cluster info
func GetClusters(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := clientset.ServerVersion()
		if err != nil {
			log.Printf(utils.LogFailedGetServerVersion, err)
			http.Error(w, utils.MsgFailedGetClusterInfo, http.StatusInternalServerError)
			return
		}

		nodes, err := clientset.CoreV1().Nodes().List(r.Context(), v1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedListNodes, err)
			http.Error(w, utils.MsgFailedGetClusterInfo, http.StatusInternalServerError)
			return
		}

		clusterInfo := models.ClusterInfo{
			Name:    "default-cluster",
			Version: version.GitVersion,
			Nodes:   len(nodes.Items),
			Healthy: true,
			Status:  "Healthy",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(clusterInfo); err != nil {
			log.Printf(utils.LogFailedEncodeClusterInfo, err)
		}
	}
}

// GetClusterHealth returns cluster health status
func GetClusterHealth(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check nodes health
		nodes, err := clientset.CoreV1().Nodes().List(r.Context(), v1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedListNodes, err)
			http.Error(w, utils.MsgFailedGetClusterHealth, http.StatusInternalServerError)
			return
		}

		healthyNodes := 0
		for _, node := range nodes.Items {
			for _, condition := range node.Status.Conditions {
				if condition.Type == "Ready" && condition.Status == "True" {
					healthyNodes++
					break
				}
			}
		}

		// Check pods health
		pods, err := clientset.CoreV1().Pods("").List(r.Context(), v1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedListPods, err)
			http.Error(w, utils.MsgFailedGetClusterHealth, http.StatusInternalServerError)
			return
		}

		runningPods := 0
		for _, pod := range pods.Items {
			if pod.Status.Phase == "Running" {
				runningPods++
			}
		}

		healthStatus := map[string]interface{}{
			"nodes": map[string]interface{}{
				"total":   len(nodes.Items),
				"healthy": healthyNodes,
			},
			"pods": map[string]interface{}{
				"total":   len(pods.Items),
				"running": runningPods,
				"failed":  len(pods.Items) - runningPods,
			},
			"overall": "Healthy",
		}

		if healthyNodes < len(nodes.Items) {
			healthStatus["overall"] = "Degraded"
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(healthStatus); err != nil {
			log.Printf(utils.LogFailedEncodeClusterHealth, err)
		}
	}
}

// GetClusterVersion returns cluster version information
func GetClusterVersion(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := clientset.ServerVersion()
		if err != nil {
			log.Printf(utils.LogFailedGetServerVersion, err)
			http.Error(w, utils.MsgFailedGetClusterVersion, http.StatusInternalServerError)
			return
		}

		versionInfo := map[string]string{
			"gitVersion":   version.GitVersion,
			"gitCommit":    version.GitCommit,
			"gitTreeState": version.GitTreeState,
			"buildDate":    version.BuildDate,
			"goVersion":    version.GoVersion,
			"compiler":     version.Compiler,
			"platform":     version.Platform,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(versionInfo); err != nil {
			log.Printf(utils.LogFailedEncodeClusterVersion, err)
		}
	}
}
