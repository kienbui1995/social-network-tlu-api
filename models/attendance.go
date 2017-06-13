package models

// Attendance struct
type Attendance struct {
	Student *StudentObject `json:"student,omitempty"`
	Class   *ClassObject   `json:"class,omitempty"`
	Status  int            `json:"status,omitempty"`
	Message string         `json:"message,omitempty"`
}
