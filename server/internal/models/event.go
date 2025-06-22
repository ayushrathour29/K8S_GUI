package models

// Event represents a simplified event view
type Event struct {
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	Type           string `json:"type"`
	InvolvedObject string `json:"involvedObject"`
	FirstTimestamp string `json:"firstTimestamp"`
	LastTimestamp  string `json:"lastTimestamp"`
	Count          int32  `json:"count"`
}

// EventListResponse represents event list response
type EventListResponse struct {
	Items []Event `json:"items"`
}
