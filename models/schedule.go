package models

// Schedule struct
type Schedule struct {
	Student  *Student  `json:"student,omitempty"`
	Classes  []*Class  `json:"classes,omitempty"`
	Status   int       `json:"status,omitempty"`
	Semester *Semester `json:"semester,omitempty"`
}
