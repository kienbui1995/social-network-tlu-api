package models

import (
	"github.com/asaskevich/govalidator"
)

//Account struct
type Account struct {
	ID       int64  `json:"id"`
	Username string `json:"username" valid:"required,length(6|32)"`
	Password string `json:"password" valid:"required,length(8|32)"`
	Role     int    `json:"role"`
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
