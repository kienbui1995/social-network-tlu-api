package controllers

import (
	"fmt"

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
