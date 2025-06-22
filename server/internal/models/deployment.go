package models

// Deployment represents a simplified deployment view
type Deployment struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Replicas          int32             `json:"replicas"`
	AvailableReplicas int32             `json:"availableReplicas"`
	CreatedAt         string            `json:"createdAt"`
	Strategy          string            `json:"strategy"`
	Labels            map[string]string `json:"labels,omitempty"`
}

// DeploymentListResponse represents deployment list response
type DeploymentListResponse struct {
	Items []Deployment `json:"items"`
}

// CreateDeploymentRequest represents the request body for creating a deployment
type CreateDeploymentRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
	Replicas  int32  `json:"replicas"`
	Port      int32  `json:"port"`
}

// UpdateDeploymentRequest represents the request body for updating a deployment
type UpdateDeploymentRequest struct {
	Image    string `json:"image"`
	Replicas int32  `json:"replicas"`
}
