package models

// ClusterInfo represents cluster information
type ClusterInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Nodes   int    `json:"nodes"`
	Healthy bool   `json:"healthy"`
	Status  string `json:"status"`
}
