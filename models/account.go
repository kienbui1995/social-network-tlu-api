package models

import (
	"github.com/asaskevich/govalidator"
)

//Account struct
type Account struct {
	ID       int64  `json:"id"`
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
	Role     int    `json:"role,omitempty"`
	Code     string `json:"code,omitempty"`
	Device   string `json:"device" valid:"required"`
}

//Accounts list
type Accounts []Account

// IsEmpty func to check entity empty
func (a Account) IsEmpty() bool {
	return a == Account{}
}

// Validate func
func (a Account) Validate() (bool, error) {
	return govalidator.ValidateStruct(a)
}
