package models

// Semester struct
type Semester struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name,omitempty"`
	Group  int    `json:"group,omitempty"`
	Year   int    `json:"year,omitempty"`
	Status int    `json:"status,omitempty"`
}
