package api

import (
	"encoding/json"
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

// NodeMetrics represents node metrics
type NodeMetrics struct {
	NodeName  string             `json:"nodeName"`
	CPU       NodeResourceMetric `json:"cpu"`
	Memory    NodeResourceMetric `json:"memory"`
	Timestamp time.Time          `json:"timestamp"`
	Window    string             `json:"window"`
	Available bool               `json:"available"`
	Message   string             `json:"message,omitempty"`
}

// NodeResourceMetric represents resource metric
type NodeResourceMetric struct {
	Value      string  `json:"value"`
	Quantity   int64   `json:"quantity"`
	Percentage float64 `json:"percentage,omitempty"`
	Unit       string  `json:"unit"`
}

// PodMetrics represents pod metrics
type PodMetrics struct {
	PodName    string             `json:"podName"`
	Namespace  string             `json:"namespace"`
	Containers []ContainerMetrics `json:"containers"`
	CPU        NodeResourceMetric `json:"cpu"`
	Memory     NodeResourceMetric `json:"memory"`
	Timestamp  time.Time          `json:"timestamp"`
	Window     string             `json:"window"`
}

// ContainerMetrics represents container metrics
type ContainerMetrics struct {
	Name   string             `json:"name"`
	CPU    NodeResourceMetric `json:"cpu"`
	Memory NodeResourceMetric `json:"memory"`
}

// NodeMetricsListResponse represents node metrics list response
type NodeMetricsListResponse struct {
	Items []NodeMetrics `json:"items"`
}

// PodMetricsListResponse represents pod metrics list response
type PodMetricsListResponse struct {
	Items []PodMetrics `json:"items"`
}

// GetNodeMetrics returns metrics for a specific node
func GetNodeMetrics(clientset *kubernetes.Clientset, metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nodeName := vars["name"]

		// Get node to check capacity
		node, err := clientset.CoreV1().Nodes().Get(r.Context(), nodeName, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get node %s: %v", nodeName, err)
			http.Error(w, "Node not found", http.StatusNotFound)
			return
		}

		// Get node metrics
		metrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(r.Context(), nodeName, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get metrics for node %s: %v", nodeName, err)
			response := NodeMetrics{
				NodeName:  nodeName,
				Available: false,
				Message:   "Metrics not available",
				Timestamp: time.Now(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
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

		response := NodeMetrics{
			NodeName: nodeName,
			CPU: NodeResourceMetric{
				Value:      cpuUsage.String(),
				Quantity:   cpuUsage.MilliValue(),
				Percentage: cpuPercentage,
				Unit:       "m",
			},
			Memory: NodeResourceMetric{
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
		json.NewEncoder(w).Encode(response)
	}
}

// GetNodesMetrics returns metrics for all nodes
func GetNodesMetrics(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := metricsClient.MetricsV1beta1().NodeMetricses().List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to get node metrics: %v", err)
			http.Error(w, "Failed to get node metrics", http.StatusInternalServerError)
			return
		}

		response := NodeMetricsListResponse{Items: make([]NodeMetrics, 0, len(metrics.Items))}
		for _, m := range metrics.Items {
			cpuUsage := m.Usage[corev1.ResourceCPU]
			memoryUsage := m.Usage[corev1.ResourceMemory]

			response.Items = append(response.Items, NodeMetrics{
				NodeName: m.Name,
				CPU: NodeResourceMetric{
					Value:    (&cpuUsage).String(),
					Quantity: (&cpuUsage).MilliValue(),
					Unit:     "m",
				},
				Memory: NodeResourceMetric{
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
		json.NewEncoder(w).Encode(response)
	}
}

// GetPodsMetrics returns metrics for all pods
func GetPodsMetrics(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := metricsClient.MetricsV1beta1().PodMetricses("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to get pod metrics: %v", err)
			http.Error(w, "Failed to get pod metrics", http.StatusInternalServerError)
			return
		}

		response := PodMetricsListResponse{Items: make([]PodMetrics, 0, len(metrics.Items))}
		for _, m := range metrics.Items {
			containers := make([]ContainerMetrics, len(m.Containers))
			totalCPU := resource.NewQuantity(0, resource.DecimalSI)
			totalMemory := resource.NewQuantity(0, resource.BinarySI)

			for i, c := range m.Containers {
				cpu := c.Usage[corev1.ResourceCPU]
				memory := c.Usage[corev1.ResourceMemory]
				totalCPU.Add(cpu)
				totalMemory.Add(memory)

				containers[i] = ContainerMetrics{
					Name: c.Name,
					CPU: NodeResourceMetric{
						Value:    cpu.String(),
						Quantity: cpu.MilliValue(),
						Unit:     "m",
					},
					Memory: NodeResourceMetric{
						Value:    memory.String(),
						Quantity: memory.Value(),
						Unit:     "bytes",
					},
				}
			}

			response.Items = append(response.Items, PodMetrics{
				PodName:    m.Name,
				Namespace:  m.Namespace,
				Containers: containers,
				CPU: NodeResourceMetric{
					Value:    totalCPU.String(),
					Quantity: totalCPU.MilliValue(),
					Unit:     "m",
				},
				Memory: NodeResourceMetric{
					Value:    totalMemory.String(),
					Quantity: totalMemory.Value(),
					Unit:     "bytes",
				},
				Timestamp: m.Timestamp.Time,
				Window:    m.Window.Duration.String(),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetPodMetricsByNamespace returns metrics for pods in a specific namespace
func GetPodMetricsByNamespace(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]

		metrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to get pod metrics for namespace %s: %v", namespace, err)
			http.Error(w, "Failed to get pod metrics", http.StatusInternalServerError)
			return
		}

		response := PodMetricsListResponse{Items: make([]PodMetrics, 0, len(metrics.Items))}
		for _, m := range metrics.Items {
			containers := make([]ContainerMetrics, len(m.Containers))
			totalCPU := resource.NewQuantity(0, resource.DecimalSI)
			totalMemory := resource.NewQuantity(0, resource.BinarySI)

			for i, c := range m.Containers {
				cpu := c.Usage[corev1.ResourceCPU]
				memory := c.Usage[corev1.ResourceMemory]
				totalCPU.Add(cpu)
				totalMemory.Add(memory)

				containers[i] = ContainerMetrics{
					Name: c.Name,
					CPU: NodeResourceMetric{
						Value:    cpu.String(),
						Quantity: cpu.MilliValue(),
						Unit:     "m",
					},
					Memory: NodeResourceMetric{
						Value:    memory.String(),
						Quantity: memory.Value(),
						Unit:     "bytes",
					},
				}
			}

			response.Items = append(response.Items, PodMetrics{
				PodName:    m.Name,
				Namespace:  m.Namespace,
				Containers: containers,
				CPU: NodeResourceMetric{
					Value:    totalCPU.String(),
					Quantity: totalCPU.MilliValue(),
					Unit:     "m",
				},
				Memory: NodeResourceMetric{
					Value:    totalMemory.String(),
					Quantity: totalMemory.Value(),
					Unit:     "bytes",
				},
				Timestamp: m.Timestamp.Time,
				Window:    m.Window.Duration.String(),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
