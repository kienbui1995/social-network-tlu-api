package models

// Attendance struct
type Attendance struct {
	Student    *Student    `json:"student,omitempty"`
	StudyShift *StudyShift `json:"study_shift,omitempty"`
	Status     int         `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
}
