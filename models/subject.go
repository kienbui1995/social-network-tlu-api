package models

// Subject struct
type Subject struct {
	ID      int64  `json:"id"`
	Code    string `json:"code"`
	Name    string `json:"name,omitempty"`
	Credits int    `json:"credits,omitempty"`
	Status  int    `json:"status,omitempty"`
}
