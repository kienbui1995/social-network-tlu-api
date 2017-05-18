package models

import "github.com/asaskevich/govalidator"

// Comment struct
type Comment struct {
	UserObject
	ID        int64  `json:"id"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
	Status    int    `json:"status,omitempty"`
	CanEdit   bool   `json:"can_edit"`
	CanDelete bool   `json:"can_delete"`
}

// Comments type is comment list
type Comments []Comment

// IsEmpty funt to check zezo-value
func (c Comment) IsEmpty() bool {
	return c == Comment{}
}

// Validate to Validate struct
func (c Comment) Validate() (bool, error) {
	return govalidator.ValidateStruct(c)
}
