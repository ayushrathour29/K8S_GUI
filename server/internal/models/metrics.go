package models

import "time"

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
