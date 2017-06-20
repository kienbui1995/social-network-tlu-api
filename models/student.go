package models

// Student struct
type Student struct {
	ID        int64  `json:"id"`
	Code      string `json:"code"`
	FirstName string `json:"first_name,omitempty"`
	Birthday  string `json:"birthday,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Status    int    `json:"status,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}
