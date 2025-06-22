package models

// Pod represents a simplified pod view
type Pod struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	Status       string            `json:"status"`
	RestartCount int32             `json:"restartCount"`
	CreatedAt    string            `json:"createdAt"`
	NodeName     string            `json:"nodeName"`
	PodIP        string            `json:"podIP"`
	Containers   []string          `json:"containers"`
	Labels       map[string]string `json:"labels,omitempty"`
}

// PodListResponse represents pod list response
type PodListResponse struct {
	Items []Pod `json:"items"`
} 