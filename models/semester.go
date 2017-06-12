package models

// Semester struct
type Semester struct {
	ID        int64  `json:"id"`
	Code      string `json:"code"`
	Symbol    string `json:"symbol,omitempty"`
	Name      string `json:"name,omitempty"`
	Group     int    `json:"group,omitempty"`
	Year      string `json:"year,omitempty"`
	StartAt   string `json:"start_at,omitempty"`
	FinishAt  string `json:"finish_at,omitempty"`
	Status    int    `json:"status,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}
