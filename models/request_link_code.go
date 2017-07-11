package models

// RequestLinkCode struct
type RequestLinkCode struct {
	ID               int64          `json:"id"`
	CreatedAt        int64          `json:"created_at,omitempty"`
	UpdatedAt        int64          `json:"updated_at,omitempty"`
	User             *UserObject    `json:"user,omitempty"`
	Student          *StudentObject `json:"student,omitempty"`
	Status           int            `json:"status,omitempty"`
	Type             int            `json:"type,omitempty"`
	Email            string         `json:"email,omitempty"`
	VerificationCode string         `json:"verification_code,omitempty"`
	Code             string         `json:"code,omitempty"`
	FullName         string         `json:"full_name,omitempty"`
	Photo            string         `json:"photo,omitempty"`
	CanEdit          bool           `json:"can_edit,omitempty"`
	CanDelete        bool           `json:"can_delete,omitempty"`
}

// IsEmpty func to check RequestLinkCode is null
func (requestLinkCode RequestLinkCode) IsEmpty() bool {
	return requestLinkCode == RequestLinkCode{}
}
