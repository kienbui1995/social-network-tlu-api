package models

// Subject struct
type Subject struct {
	ID        int64   `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name,omitempty"`
	Credits   int     `json:"credits,omitempty"`
	Factor    float32 `json:"factor,omitempty"`
	Status    int     `json:"status,omitempty"`
	CreatedAt int64   `json:"created_at,omitempty"`
	UpdatedAt int64   `json:"updated_at,omitempty"`
}
