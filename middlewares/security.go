package middlewares

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

var secret = []byte(configs.JWTSecretKey)

// SetToken func
func SetToken(c *gin.Context) {

	var account models.Account
	if errBindJSON := c.BindJSON(&account); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcAuthLogin, "BindJSON: "+errBindJSON.Error())
		return
	}

	// valid
	if _, errValidateStruct := valid.ValidateStruct(account); errValidateStruct != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcAuthLogin, "ValidateStruct: "+errValidateStruct.Error())
		return
	}

	service := services.NewAccountService()
	accountID, errLogin := service.Login(account)
	if errLogin != nil {
		if accountID >= 0 {
			helpers.ResponseAuthJSON(c, configs.EcAuthWrongPassword, errLogin.Error())
			return
		}
		if accountID < 0 {
			helpers.ResponseAuthJSON(c, configs.EcAuthNoExistUser, errLogin.Error())
			return
		}
	}
	account.ID = accountID
	tokenstring, errGenerateToken := helpers.GenerateToken(account, secret)
	if errGenerateToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GenerateToken helpers: %s\n", errGenerateToken.Error())
		return
	}
	if _, errSaveToken := service.SaveToken(account, tokenstring); errSaveToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("SaveToken helpers: %s\n", errSaveToken.Error())
		return
	}
	helpers.ResponseSuccessJSON(c, configs.EcSuccess, "Set token successful", map[string]interface{}{"id": accountID, "token": tokenstring})
}

// AuthRequired func to check token in header
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.Request.Header.Get("token")

		if tokenString == "" || len(tokenString) == 0 {
			helpers.ResponseAuthJSON(c, configs.EcAuthMissingToken, "Missing token")
			return
		}
		service := services.NewAccountService()
		valid, errValidateToken := helpers.ValidateToken(tokenString, secret)
		if errValidateToken != nil {
			helpers.ResponseAuthJSON(c, configs.EcAuthInvalidToken, "ValidateToken helpers: "+errValidateToken.Error())
			return
		} else if valid != true {
			helpers.ResponseAuthJSON(c, configs.EcAuthInvalidToken, "ValidateToken helpers: valid is false")
			return
		}
		claims, errExtractClaims := helpers.ExtractClaims(tokenString, secret)
		if errExtractClaims != nil {
			helpers.ResponseAuthJSON(c, configs.EcAuthInvalidToken, "ExtractClaims helpers: "+errExtractClaims.Error())
			return
		}
		accountID := int64(claims["userid"].(float64))
		fmt.Printf("accID: %v\n", accountID)
		existToken, errCheckExistToken := service.CheckExistToken(accountID, tokenString)
		if errCheckExistToken != nil {
			helpers.ResponseAuthJSON(c, configs.EcAuthNoExistToken, "CheckExistToken services: "+errCheckExistToken.Error())
			return
		}
		if existToken != true {
			helpers.ResponseAuthJSON(c, configs.EcAuthNoExistToken, "CheckExistToken services: NULL")
			return
		}
		c.Next()

	}
}
