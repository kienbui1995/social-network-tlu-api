package models

// Classroom struct
type Classroom struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name,omitempty"`
	Status int    `json:"status,omitempty"`
}
