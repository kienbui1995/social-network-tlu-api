package controllers

import (
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

var secret = []byte(configs.JWTSecretKey)

// AccountController controller
type AccountController struct {
	Service services.AccountServiceInterface
}

// SignUp func
func (controller AccountController) SignUp(c *gin.Context) {
	var user models.User
	// BadRequest (400)
	if errBindJSON := c.BindJSON(&user); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}
	// Valida Invalid Entity (422)
	if govalidator.IsByteLength(user.Username, 3, 15) == false {
		helpers.ResponseErrorJSON(c, helpers.NewErrorDetail(382, "Please enter a valid username."))
		return
	}
	if govalidator.IsEmail(user.Email) == false {
		helpers.ResponseErrorJSON(c, helpers.NewErrorDetail(385, "Please enter a valid email address."))
		return
	}

	if exist, _ := controller.Service.CheckExistUsername(user.Username); exist == true {
		helpers.ResponseErrorJSON(c, helpers.NewErrorDetail(376, "The login credential you provided belongs to an existing account"))
		return
	}

	if exist, _ := controller.Service.CheckExistEmail(user.Email); exist == true {
		helpers.ResponseErrorJSON(c, helpers.NewErrorDetail(371, "The email address you provided belongs to an existing account"))
		return
	}

	user.Status = 0
	user.Followers = 0
	user.Followings = 0
	user.Posts = 0
	if _, errValidateStruct := govalidator.ValidateStruct(user); errValidateStruct != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcUsersRegisterAddUserFailed, "ValidateStruct: "+errValidateStruct.Error())
		return
	}
	userID, errUser := controller.Service.Create(user)
	if errUser != nil {
		helpers.ResponseServerErrorJSON(c) //, 400, 387, "There was an error with your registration. Please try registering again: "+errUser.Error(), nil)
		fmt.Printf("Create service: %s\n", errUser.Error())
		return

	}
	if userID < 0 {
		helpers.ResponseServerErrorJSON(c) //, 400, 387, "There was an error with your registration. Please try registering again: "+errUser.Error(), nil)
		fmt.Printf("Create service: userid <0")
		return
	}

	// activecode := helpers.RandStringBytes(6)
	// go func() {
	// 	if err := controller.Service.CreateEmailActive(user.Email, activecode, user.ID); err != nil {
	// 		fmt.Printf("\nCreate Email Active faile: %s", err.Error())
	// 	}
	//
	// }()

	//pause func send mail active
	go func() {
		activecode := helpers.RandStringBytes(6)
		created, errCreateEmailActive := controller.Service.CreateEmailActive(user.Email, activecode, user.ID)
		if errCreateEmailActive != nil {
			fmt.Printf("\nCreateEmailActive: %s\n", errCreateEmailActive.Error())
		}

		if created != true {
			fmt.Printf("CreateEmailActive: don't create")
		}
		sender := helpers.NewSender(configs.MailAddress, configs.MailKey)
		var email []string
		email = append(email, user.Email)
		linkActive := "<a href='tlu.cloudapp.net:8080/activation?use_id=" + string(userID) + "&active_code=" + activecode + "'>Active</a>"
		sender.SendMail(email, fmt.Sprintf("Active user %s on TLSEN", user.Username), fmt.Sprintf("Content-Type: text/html; charset=UTF-8\n\ncode: %s OR active via link: %s", activecode, linkActive))
	}()

	helpers.ResponseCreatedJSON(c, 1, "Create user successful!", map[string]interface{}{"id": userID})
}

// Login func
func (controller AccountController) Login(c *gin.Context) {
	var json models.Account
	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	if _, errValidate := json.Validate(); errValidate != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Validate: "+errValidate.Error())
		return
	}
	var errLogin error
	json.ID, errLogin = controller.Service.Login(json)
	if errLogin != nil {
		if json.ID == -1 {
			helpers.ResponseAuthJSON(c, configs.EcAuthNoExistUser, errLogin.Error())
		} else {
			helpers.ResponseAuthJSON(c, configs.EcAuthWrongPassword, errLogin.Error())
		}
		return
	}

	token, errGenerateToken := helpers.GenerateToken(json, secret)
	if errGenerateToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GenerateToken helpers: %v\n", errGenerateToken.Error())
	}

	if _, errSaveToken := controller.Service.SaveToken(json, token); errSaveToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("SaveToken service: %v\n", errSaveToken.Error())
	}
	helpers.ResponseSuccessJSON(c, configs.EcSuccess, "Login successful", map[string]interface{}{"id": json.ID, "token": token})
}

// Logout func
func (controller AccountController) Logout(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if len(token) == 0 {
		helpers.ResponseBadRequestJSON(c, configs.EcAuthMissingToken, "Missing token")
		return
	}
	if _, errValidateToken := helpers.ValidateToken(token, secret); errValidateToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcAuthInvalidToken, "Invalid token: "+errValidateToken.Error())
		return
	}
	claims, errExtractClaims := helpers.ExtractClaims(token, secret)
	if errExtractClaims != nil {
		helpers.ResponseAuthJSON(c, configs.EcAuthInvalidToken, "ExtractClaims: "+errExtractClaims.Error())
		return
	}
	accountid := int64(claims["userid"].(float64))
	if _, errCheckExistToken := controller.Service.CheckExistToken(accountid, token); errCheckExistToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcAuthNoExistToken, "CheckExistToken service: "+errCheckExistToken.Error())
		return
	}

	deleted, errDeleteToken := controller.Service.DeleteToken(accountid, token)
	if errDeleteToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("DeleteToken service: %v\n", errDeleteToken.Error())
	}

	if deleted == true {
		helpers.ResponseJSON(c, 200, 1, "Logout successful", nil)
	}
}

// LoginViaFacebook func is login or sign up via Facebook
func (controller AccountController) LoginViaFacebook(c *gin.Context) {
	type FacebookToken struct {
		ID          string `json:"id" valid:"required"`
		AccessToken string `json:"access_token" valid:"required"`
		Device      string `json:"device" valid:"required"`
	}
	var json FacebookToken
	errBindJSON := c.BindJSON(&json)
	if errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}
	if len(json.ID) == 0 || len(json.Device) == 0 || len(json.AccessToken) == 0 {
		helpers.ResponseBadRequestJSON(c, configs.EcParamMissingField, "Missing a few fields.")
		return
	}

	if id, errExist := controller.Service.CheckExistFacebookID(json.ID); errExist == nil && id != 0 {
		verify := helpers.VerifyFacebookID(json.ID, json.AccessToken)
		if verify == true {
			acc := models.Account{ID: id, Device: json.Device}
			tokenstring, errtoken := helpers.GenerateToken(acc, secret)
			if errtoken != nil {
				helpers.ResponseServerErrorJSON(c)
				fmt.Printf("GenerateToken helpers: %s\n", errtoken.Error())
				return
			}
			if _, errSaveToken := controller.Service.SaveToken(acc, tokenstring); errSaveToken != nil {
				helpers.ResponseBadRequestJSON(c, 1, "Don't save token"+errSaveToken.Error())
				return
			}
			data := map[string]interface{}{"token": tokenstring, "id": id}
			helpers.ResponseSuccessJSON(c, 1, "Login successful!", data)
			return
		}
		helpers.ResponseBadRequestJSON(c, configs.EcAuthInvalidFacebookToken, "Error in checking facebook access token.")
	} else {
		helpers.ResponseNotFoundJSON(c, configs.EcAuthNoExistFacebook, "No exist account with this facebook.")
		return
	}
	// libs.ResponseBadRequestJSON(c, -1, "Login Facebook fail")
}

// ForgotPassword func
func (controller AccountController) ForgotPassword(c *gin.Context) {
	var json struct {
		Email string `json:"email"`
	}

	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	if len(json.Email) == 0 {
		helpers.ResponseBadRequestJSON(c, 101, "Missing a few fields")
		c.Abort()
		return
	}
	existemail, errCheckExistEmail := controller.Service.CheckExistEmail(json.Email)
	if errCheckExistEmail != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistEmail service: %s\n", errCheckExistEmail.Error())
		return
	}

	if existemail == true { // send password via mail
		type RecoverPassword struct {
			Email        string `json:"email"`
			RecoveryCode string `json:"recovery_code"`
		}
		recoverpass := RecoverPassword{Email: json.Email, RecoveryCode: helpers.RandNumberBytes(6)}
		created, errCreateRecoverPassword := controller.Service.CreateRecoverPassword(recoverpass.Email, recoverpass.RecoveryCode)
		if errCreateRecoverPassword != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("CreateRecoverPassword service: %s\n", errCreateRecoverPassword.Error())
			return
		}
		if created != true {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("CreateRecoverPassword service: don't create\n")
			return
		}
		sender := helpers.NewSender(configs.MailAddress, configs.MailKey)
		var email []string
		email = append(email, recoverpass.Email)
		go sender.SendMail(email, fmt.Sprintf("Recover password on TLSEN"), fmt.Sprintf("\ncode: %s\n Please verify within 30 minutes.", recoverpass.RecoveryCode))
		helpers.ResponseJSON(c, 200, 1, "A email sent", nil)
	} else { // no exist email
		helpers.ResponseAuthJSON(c, configs.EcAuthNoExistEmail, "No exist email")
	}
}

// VerifyRecoveryCode func
func (controller AccountController) VerifyRecoveryCode(c *gin.Context) {
	var json struct {
		Email        string `json:"email"`
		RecoveryCode string `json:"recovery_code"`
	}
	errBindJSON := c.BindJSON(&json)
	if errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	if len(json.Email) == 0 || len(json.RecoveryCode) == 0 {
		helpers.ResponseBadRequestJSON(c, configs.EcParamMissingField, "Missing a few fields")
		return
	}
	// check exist user with this email
	exist, errCheckExistEmail := controller.Service.CheckExistEmail(json.Email)
	if errCheckExistEmail != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistEmail: %s\n" + errCheckExistEmail.Error())

	} else if exist == true {
		id, errVerifyRecoveryCode := controller.Service.VerifyRecoveryCode(json.Email, json.RecoveryCode)
		if errVerifyRecoveryCode != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("Error in verify recovery code: %s\n" + errVerifyRecoveryCode.Error())

		}
		if id >= 0 { // generate a key
			key := helpers.RandStringBytes(6)

			helpers.ResponseSuccessJSON(c, 1, "ID user and key to create new password", map[string]interface{}{"id": id, "recovery_key": key})
			go func() {
				errAddUserRecoveryKey := controller.Service.AddUserRecoveryKey(id, key)
				if errAddUserRecoveryKey != nil {
					helpers.ResponseServerErrorJSON(c)
					fmt.Printf("errAddUserRecoveryKey service: %s\n", errAddUserRecoveryKey.Error())
				}
			}()

		} else {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("VerifyRecoveryCode:  NULL")
		}
	} else { // user with this emmail don't exist
		helpers.ResponseNotFoundJSON(c, 409, "No exist user.")
	}
}

// RenewPassword func
func (controller AccountController) RenewPassword(c *gin.Context) {

	var json struct {
		ID          int64  `json:"id"`
		RecoveryKey string `json:"recovery_key"`
		NewPassword string `json:"new_password"`
	}
	errBindJSON := c.BindJSON(&json)
	if errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter")
		return
	}

	if json.ID == 0 || len(json.RecoveryKey) == 0 || len(json.NewPassword) == 0 {
		helpers.ResponseBadRequestJSON(c, configs.EcParamMissingField, "Missing a few fields")
		return
	}

	if renew, errRenewPassword := controller.Service.RenewPassword(json.ID, json.RecoveryKey, json.NewPassword); errRenewPassword == nil && renew == true {
		helpers.ResponseJSON(c, 200, 1, "Renew password successful", nil)

		go func() {
			deleted, errDeleteRecoveryProperty := controller.Service.DeleteRecoveryProperty(json.ID)
			if errDeleteRecoveryProperty != nil {
				fmt.Printf("DeleteRecoveryProperty service: " + errDeleteRecoveryProperty.Error())
			}
			if deleted != true {
				fmt.Printf("DeleteRecoveryProperty service: Don't delete")
			}
		}()

	} else if errRenewPassword != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("RenewPassword service: %s\n", errRenewPassword.Error())
	} else if renew != true {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("RenewPassword service: renew false")
	}

}

// ActiveByEmail func
func (controller AccountController) ActiveByEmail(c *gin.Context) {
	accountID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	activeCode := c.Query("active_code")

	if accountID < 0 || len(activeCode) == 0 {
		helpers.ResponseBadRequestJSON(c, configs.EcParamMissingField, "Missing a few fields")
		return
	}
	if exist, _ := services.NewUserService().CheckExistUser(accountID); exist != true {
		helpers.ResponseBadRequestJSON(c, configs.EcAuthNoExistUser, "No exist user")
		return
	}
	actived, errActiveByEmail := controller.Service.ActiveByEmail(accountID, activeCode)
	if errActiveByEmail == nil && actived == true {
		helpers.ResponseJSON(c, 200, 1, "Active account successful", nil)

		go func() {
			deleted, errDeleteActiveCode := controller.Service.DeleteActiveCode(accountID)
			if errDeleteActiveCode != nil {
				fmt.Printf("DeleteActiveCode service: " + errDeleteActiveCode.Error())
			}
			if deleted != true {
				fmt.Printf("DeleteActiveCode service: Don't delete")
			}
		}()

	} else if errActiveByEmail != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("ActiveByEmail service: %s\n", errActiveByEmail.Error())
	} else if actived != true {
		helpers.ResponseJSON(c, 200, 1, "active false", nil)
		return
	}

}
