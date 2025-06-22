package models

// Namespace represents a simplified namespace view
type Namespace struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	CreatedAt string            `json:"createdAt"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// NamespaceListResponse represents namespace list response
type NamespaceListResponse struct {
	Items []Namespace `json:"items"`
}

// CreateNamespaceRequest represents the request body for creating a namespace
type CreateNamespaceRequest struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
} 