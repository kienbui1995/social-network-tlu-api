package models

import "github.com/asaskevich/govalidator"

// Comment struct
type Comment struct {
	Owner            *UserObject     `json:"owner"`
	Post             *PostObject     `json:"post,omitempty"`
	ID               int64           `json:"id"`
	Message          string          `json:"message"`
	Mentions         []MentionObject `json:"mentions,omitempty"`
	CreatedAt        int64           `json:"created_at,omitempty"`
	UpdatedAt        int64           `json:"updated_at,omitempty"`
	Status           int             `json:"status,omitempty"`
	CanReportToAdmin bool            `json:"can_report_to_admin,omitempty"`
	CanReport        bool            `json:"can_report"`
	CanEdit          bool            `json:"can_edit"`
	CanDelete        bool            `json:"can_delete"`
}

// Comments type is comment list
type Comments []Comment

// Validate to Validate struct
func (c Comment) Validate() (bool, error) {
	return govalidator.ValidateStruct(c)
}
