package models

// Student struct
type Student struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Photo  string `json:"photo,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Status int    `json:"status,omitempty"`
}
