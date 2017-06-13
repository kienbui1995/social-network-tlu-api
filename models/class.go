package models

// Class struct
type Class struct {
	Teacher   *TeacherObject `json:"teacher,omitempty"`
	Subject   *SubjectObject `json:"subject,omitempty"`
	Room      *RoomObject    `json:"room,omitempty"`
	ID        int64          `json:"id"`
	Code      string         `json:"code"`
	Name      string         `json:"name,omitempty"`
	Symbol    string         `json:"symbol,omitempty"`
	Day       string         `json:"day,omitempty"`
	StartAt   string         `json:"start_at,omitempty"`
	FinishAt  string         `json:"finish_at,omitempty"`
	Status    int            `json:"status,omitempty"`
	CreatedAt int64          `json:"created_at,omitempty"`
	UpdatedAt int64          `json:"updated_at,omitempty"`
}
