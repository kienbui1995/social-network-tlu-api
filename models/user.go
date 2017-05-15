package models

import "github.com/asaskevich/govalidator"

// User struct
type User struct {
	ID            int64  `json:"id,omitempty"`
	Username      string `json:"username,omitempty" valid:"required,length(8|32)"`
	Password      string `json:"password,omitempty" valid:"length(7|32)"`
	FullName      string `json:"full_name,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	MiddleName    string `json:"middle_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Birthday      string `json:"birthday,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	Cover         string `json:"cover,omitempty"`
	About         string `json:"about,omitempty"`
	Gender        string `json:"gender,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty" valid:"email"`
	FacebookID    string `json:"facebook_id,omitempty"`
	FacebookToken string `json:"facebook_token,omitempty"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
	IsVertified   bool   `json:"is_vertified,omitempty"`
	Status        int    `json:"status"`
	Posts         int    `json:"posts"`
	Followers     int    `json:"followers"`
	Followings    int    `json:"followings" `
}

// Users list
type Users []User

// IsEmpty func to check entity empty
func (u User) IsEmpty() bool {
	return u == User{}
}

// Validate to Validate struct
func (u User) Validate() (bool, error) {
	return govalidator.ValidateStruct(u)
}
