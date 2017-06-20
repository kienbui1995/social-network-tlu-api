package models

// ExamSchedule struct
type ExamSchedule struct {
	Subject   *SubjectObject   `json:"subject,omitempty"`
	Room      *RoomObject      `json:"room,omitempty"`
	Semester  *SemesterObject  `json:"semester,omitempty"`
	Students  []*StudentObject `json:"students,omitempty"`
	ID        int64            `json:"id"`
	Day       string           `json:"day,omitempty"`
	ExamTime  string           `json:"exam_time,omitempty"`
	Status    int              `json:"status,omitempty"`
	CreatedAt int64            `json:"created_at,omitempty"`
	UpdatedAt int64            `json:"updated_at,omitempty"`
}
