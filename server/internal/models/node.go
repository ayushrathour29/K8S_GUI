package models

// Node represents a simplified node view
type Node struct {
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Version     string            `json:"version"`
	OSImage     string            `json:"osImage"`
	Capacity    map[string]string `json:"capacity"`
	Allocatable map[string]string `json:"allocatable"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// NodeListResponse represents node list response
type NodeListResponse struct {
	Items []Node `json:"items"`
} 