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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

// GetNodeMetrics returns metrics for a specific node
func GetNodeMetrics(clientset *kubernetes.Clientset, metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nodeName := vars["name"]

		// Get node to check capacity
		node, err := clientset.CoreV1().Nodes().Get(r.Context(), nodeName, metav1.GetOptions{})
		if err != nil {
			log.Printf(utils.LogFailedGetNodeMetrics, nodeName, err)
			http.Error(w, utils.MsgNodeNotFound, http.StatusNotFound)
			return
		}

		// Get node metrics
		metrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(r.Context(), nodeName, metav1.GetOptions{})
		if err != nil {
			log.Printf(utils.LogFailedGetNodeMetricsAPI, nodeName, err)
			response := models.NodeMetrics{
				NodeName:  nodeName,
				Available: false,
				Message:   "Metrics not available",
				Timestamp: time.Now(),
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Printf(utils.LogFailedEncodeNodeMetricsError, err)
			}
			return
		}

		// Calculate CPU metrics
		cpuUsage := metrics.Usage[corev1.ResourceCPU]
		cpuCapacity := node.Status.Capacity[corev1.ResourceCPU]
		cpuPercentage := float64(cpuUsage.MilliValue()) / float64(cpuCapacity.MilliValue()) * 100

		// Calculate Memory metrics
		memoryUsage := metrics.Usage[corev1.ResourceMemory]
		memoryCapacity := node.Status.Capacity[corev1.ResourceMemory]
		memoryPercentage := float64(memoryUsage.Value()) / float64(memoryCapacity.Value()) * 100

		response := models.NodeMetrics{
			NodeName: nodeName,
			CPU: models.NodeResourceMetric{
				Value:      cpuUsage.String(),
				Quantity:   cpuUsage.MilliValue(),
				Percentage: cpuPercentage,
				Unit:       "m",
			},
			Memory: models.NodeResourceMetric{
				Value:      memoryUsage.String(),
				Quantity:   memoryUsage.Value(),
				Percentage: memoryPercentage,
				Unit:       "bytes",
			},
			Timestamp: metrics.Timestamp.Time,
			Window:    metrics.Window.Duration.String(),
			Available: true,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodeNodeMetrics, err)
		}
	}
}

// GetNodesMetrics returns metrics for all nodes
func GetNodesMetrics(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := metricsClient.MetricsV1beta1().NodeMetricses().List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedGetNodeMetricsList, err)
			http.Error(w, utils.MsgFailedGetNodeMetrics, http.StatusInternalServerError)
			return
		}

		response := models.NodeMetricsListResponse{Items: make([]models.NodeMetrics, 0, len(metrics.Items))}
		for _, m := range metrics.Items {
			cpuUsage := m.Usage[corev1.ResourceCPU]
			memoryUsage := m.Usage[corev1.ResourceMemory]

			response.Items = append(response.Items, models.NodeMetrics{
				NodeName: m.Name,
				CPU: models.NodeResourceMetric{
					Value:    (&cpuUsage).String(),
					Quantity: (&cpuUsage).MilliValue(),
					Unit:     "m",
				},
				Memory: models.NodeResourceMetric{
					Value:    (&memoryUsage).String(),
					Quantity: (&memoryUsage).Value(),
					Unit:     "bytes",
				},
				Timestamp: m.Timestamp.Time,
				Window:    m.Window.Duration.String(),
				Available: true,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodeNodeMetricsList, err)
		}
	}
}

// GetPodsMetrics returns metrics for all pods
func GetPodsMetrics(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := metricsClient.MetricsV1beta1().PodMetricses("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedGetPodMetrics, err)
			http.Error(w, utils.MsgFailedGetPodMetrics, http.StatusInternalServerError)
			return
		}

		response := models.PodMetricsListResponse{Items: make([]models.PodMetrics, 0, len(metrics.Items))}
		for _, m := range metrics.Items {
			containers := make([]models.ContainerMetrics, len(m.Containers))
			totalCPU := resource.NewQuantity(0, resource.DecimalSI)
			totalMemory := resource.NewQuantity(0, resource.BinarySI)

			for i, c := range m.Containers {
				cpu := c.Usage[corev1.ResourceCPU]
				memory := c.Usage[corev1.ResourceMemory]
				totalCPU.Add(cpu)
				totalMemory.Add(memory)

				containers[i] = models.ContainerMetrics{
					Name: c.Name,
					CPU: models.NodeResourceMetric{
						Value:    cpu.String(),
						Quantity: cpu.MilliValue(),
						Unit:     "m",
					},
					Memory: models.NodeResourceMetric{
						Value:    memory.String(),
						Quantity: memory.Value(),
						Unit:     "bytes",
					},
				}
			}

			response.Items = append(response.Items, models.PodMetrics{
				PodName:    m.Name,
				Namespace:  m.Namespace,
				Containers: containers,
				CPU: models.NodeResourceMetric{
					Value:    totalCPU.String(),
					Quantity: totalCPU.MilliValue(),
					Unit:     "m",
				},
				Memory: models.NodeResourceMetric{
					Value:    totalMemory.String(),
					Quantity: totalMemory.Value(),
					Unit:     "bytes",
				},
				Timestamp: m.Timestamp.Time,
				Window:    m.Window.Duration.String(),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodePodMetricsList, err)
		}
	}
}

// GetPodMetricsByNamespace returns metrics for pods in a specific namespace
func GetPodMetricsByNamespace(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]

		metrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf(utils.LogFailedGetPodMetricsNamespace, namespace, err)
			http.Error(w, utils.MsgFailedGetPodMetrics, http.StatusInternalServerError)
			return
		}

		response := models.PodMetricsListResponse{Items: make([]models.PodMetrics, 0, len(metrics.Items))}
		for _, m := range metrics.Items {
			containers := make([]models.ContainerMetrics, len(m.Containers))
			totalCPU := resource.NewQuantity(0, resource.DecimalSI)
			totalMemory := resource.NewQuantity(0, resource.BinarySI)

			for i, c := range m.Containers {
				cpu := c.Usage[corev1.ResourceCPU]
				memory := c.Usage[corev1.ResourceMemory]
				totalCPU.Add(cpu)
				totalMemory.Add(memory)

				containers[i] = models.ContainerMetrics{
					Name: c.Name,
					CPU: models.NodeResourceMetric{
						Value:    cpu.String(),
						Quantity: cpu.MilliValue(),
						Unit:     "m",
					},
					Memory: models.NodeResourceMetric{
						Value:    memory.String(),
						Quantity: memory.Value(),
						Unit:     "bytes",
					},
				}
			}

			response.Items = append(response.Items, models.PodMetrics{
				PodName:    m.Name,
				Namespace:  m.Namespace,
				Containers: containers,
				CPU: models.NodeResourceMetric{
					Value:    totalCPU.String(),
					Quantity: totalCPU.MilliValue(),
					Unit:     "m",
				},
				Memory: models.NodeResourceMetric{
					Value:    totalMemory.String(),
					Quantity: totalMemory.Value(),
					Unit:     "bytes",
				},
				Timestamp: m.Timestamp.Time,
				Window:    m.Window.Duration.String(),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf(utils.LogFailedEncodePodMetricsListNamespace, err)
		}
	}
}
