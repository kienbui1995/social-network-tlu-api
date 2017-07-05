package models

// Violation struct
type Violation struct {
	ID        int64              `json:"id"`
	Owner     *StudentObject     `json:"owner,omitempty"`
	Catcher   *SupervisiorObject `json:"catcher,omitempty"`
	Message   string             `json:"message"`
	Photo     string             `json:"photo,omitempty"`
	TimeAt    int64              `json:"time_at,omitempty"`
	Place     string             `json:"place,omitempty"`
	UpdatedAt int64              `json:"updated_at,omitempty"`
	CreatedAt int64              `json:"created_at,omitempty"`
	Status    int                `json:"status,omitemtpy"`
}

// IsEmpty func
func (violation Violation) IsEmpty() bool {
	return violation == Violation{}
}
