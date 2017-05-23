package models

import "github.com/asaskevich/govalidator"

// PublicUser struct
type PublicUser struct {
	ID          int64  `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	FullName    string `json:"full_name,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	MiddleName  string `json:"middle_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	BirthDay    string `json:"birthday,omitempty"`
	LargeAvatar string `json:"large_avatar,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
	About       string `json:"about,omitempty"`
	Gender      int    `json:"gender,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty" valid:"email"`
	FacebookID  string `json:"facebook_id,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
	Status      int    `json:"status"`
	Posts       int    `json:"posts"`
	Followers   int    `json:"followers"`
	Followings  int    `json:"followings"`
}

// PublicUsers list
type PublicUsers []PublicUser

// IsEmpty func to check entity empty
func (u PublicUser) IsEmpty() bool {
	return u == PublicUser{}
}

// Validate to Validate struct
func (u PublicUser) Validate() (bool, error) {
	return govalidator.ValidateStruct(u)
}
