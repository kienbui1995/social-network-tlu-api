package models

// Class struct
type Class struct {
	ID          int64         `json:"id"`
	Code        string        `json:"code"`
	Name        string        `json:"name,omitempty"`
	StudyShifts []*StudyShift `json:"study_shift"`
	Status      int           `json:"status,omitempty"`
	Teacher     *Teacher      `json:"teacher,omitempty"`
	Subject     *Subject      `json:"subject,omitempty"`
	Semester    *Semester     `json:"semester,omitempty"`
}

// StudyShift struct
type StudyShift struct {
	ID      int64      `json:"id"`
	Day     string     `json:"day"`
	StartAt int        `json:"start_at"`
	EndAt   int        `json:"end_at"`
	Room    *Classroom `json:"room,omitempty"`
	Status  int        `json:"status,omitempty"`
}
