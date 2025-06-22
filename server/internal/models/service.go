package models

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

// CreateServiceRequest represents the request body for creating a service
type CreateServiceRequest struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
}
