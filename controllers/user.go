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

// UserController controller
type UserController struct {
	Service services.UserServiceInterface
}

// GetAll func
func (controller UserController) GetAll(c *gin.Context) {
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	users, errGetAll := controller.Service.GetAll(params)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get user list successful", users, params, len(users))
}

// Get func
func (controller UserController) Get(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, _ := controller.Service.Get(userID)
	// Valida se existe a pessoa (404)
	if user.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found User")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Info of user", user)
}

// Delete func
func (controller UserController) Delete(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, _ := controller.Service.Get(userID)
	// Valida se existe a pessoa que será excluida (404)
	if user.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found User")
		return
	}
	// Valida se deu erro ao tentar excluir (500)
	if _, err := controller.Service.Delete(userID); err != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: %s\n", err)
		return
	}
	helpers.ResponseNoContentJSON(c)
}

// Create func
func (controller UserController) Create(c *gin.Context) {
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

	helpers.ResponseCreatedJSON(c, 1, "Create user successful!", map[string]interface{}{"id": userID})
}

// Update func
func (controller UserController) Update(c *gin.Context) {
	userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
		return
	}
	exist, errCheckExistUser := controller.Service.CheckExistUser(userID)
	if errCheckExistUser != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistUser service: %s\n", errCheckExistUser.Error())
		return
	}
	if exist == false {
		helpers.ResponseNotFoundJSON(c, configs.EcAuthNoExistUser, "Don't exist user")
		return
	}

	if myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token")); myUserID != userID || errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permission error")
		return
	}

	userOld, errGet := controller.Service.Get(userID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	var newUser models.InfoUser
	// Valida BadRequest (400)
	if errBindJSON := c.BindJSON(&newUser); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	errReplace := helpers.Replace(userOld, &newUser)
	if errReplace != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Replace helpers: %s\n", errReplace.Error())
	}

	// Valida se deu erro ao inserir (500)
	userUpdate, errUpdate := controller.Service.Update(userID, newUser)
	if errUpdate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: %s\n", errUpdate.Error())
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Update user successful", userUpdate)
}

// RequestLinkCode func
func (controller UserController) RequestLinkCode(c *gin.Context) {
	userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
		return
	}
	exist, errCheckExistUser := controller.Service.CheckExistUser(userID)
	if errCheckExistUser != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistUser service: %s\n", errCheckExistUser.Error())
		return
	}
	if exist == false {
		helpers.ResponseNotFoundJSON(c, configs.EcAuthNoExistUser, "Don't exist user")
		return
	}

	if myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token")); myUserID != userID || errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permission error")
		return
	}

	json := struct {
		FullName string `json:"full_name,omitempty"`
		Code     string `json:"code,omitempty"`
		Photo    string `json:"photo,omitempty"`
		Email    string `json:"email,omitempty"`
	}{}
	// Valida BadRequest (400)
	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}
	request := models.RequestLinkCode{}
	errReplace := helpers.Replace(json, &request)
	if errReplace != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Replace helpers: %s\n", errReplace.Error())
	}
	activecode := helpers.RandStringBytes(6)
	request.VerificationCode = activecode
	requestID, errCreateRequestLinkCode := controller.Service.CreateRequestLinkCode(request, userID)
	if errCreateRequestLinkCode != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CreateRequestLinkCode service: %s\n", errCreateRequestLinkCode.Error())
		return
	}
	if requestID < 0 {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CreateRequestLinkCode service: Don't create request link code\n")
		return
	}

	if len(request.Email) > 0 {
		go func() {
			sender := helpers.NewSender(configs.MailAddress, configs.MailKey)
			var email []string
			email = append(email, request.Email)
			// linkActive := "<a href='tlu.cloudapp.net:8080/activation?use_id=" + string(userID) + "&active_code=" + activecode + "'>Active</a>"
			sender.SendMail(email, fmt.Sprintf("link code user on TLSEN"), fmt.Sprintf("Content-Type: text/html; charset=UTF-8\n\ncode: %s for request id %d", request.VerificationCode, requestID))
		}()
	}
	helpers.ResponseSuccessJSON(c, 1, "create request link code successful", map[string]interface{}{"id": requestID})
}

// AcceptLinkCode func
func (controller UserController) AcceptLinkCode(c *gin.Context) {
	requestID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid request id")
		return
	}
	// exist, errCheckExistUser := controller.Service.CheckExistUser(userID)
	// if errCheckExistUser != nil {
	// 	helpers.ResponseServerErrorJSON(c)
	// 	fmt.Printf("CheckExistUser service: %s\n", errCheckExistUser.Error())
	// 	return
	// }
	// if exist == false {
	// 	helpers.ResponseNotFoundJSON(c, configs.EcAuthNoExistUser, "Don't exist user")
	// 	return
	// }
	code := c.Query("code")
	if len(code) > 0 {
		accepted, errAcceptLinkCode := controller.Service.AcceptLinkCodeByEmail(requestID, code)
		if errAcceptLinkCode != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("AcceptLinkCode service: %s\n", errAcceptLinkCode.Error())
			return
		}
		if accepted {
			helpers.ResponseSuccessJSON(c, 1, "accept request link code successful", nil)
		}
		helpers.ResponseBadRequestJSON(c, configs.EcAuthWrongRecoveryCode, "wrong code")
		return
	}
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permission error")
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errGetRoleFromUserID := GetRoleFromUserID(myUserID)
	if errGetRoleFromUserID != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permission error")
		fmt.Printf("GetRoleFromUserID controller: %s\n", errGetRoleFromUserID.Error())
		return
	}

	if role != configs.IAdminRole {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permission error")
		return
	}

	accepted, errAcceptLinkCode := controller.Service.AcceptLinkCode(requestID)
	if errAcceptLinkCode != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("AcceptLinkCode service: %s\n", errAcceptLinkCode.Error())
		return
	}
	if accepted {
		helpers.ResponseSuccessJSON(c, 1, "accept request link code successful", nil)
	}
}

// DeleteRequestLinkCode func
func (controller UserController) DeleteRequestLinkCode(c *gin.Context) {
	requestID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	//user, _ := controller.Service.Get
	// Valida se existe a pessoa que será excluida (404)
	// if user.IsEmpty() {
	// 	helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found User")
	// 	return
	// }
	// Valida se deu erro ao tentar excluir (500)

	exist, errCheckExistRequestLinkCode := controller.Service.CheckExistRequestLinkCode(requestID)
	if errCheckExistRequestLinkCode != nil {
		helpers.ResponseServerErrorJSON(c)
		return
	}
	if exist == false {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not found request link code")
		return
	}
	if _, err := controller.Service.DeleteRequestLinkCode(requestID); err != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: %s\n", err)
		return
	}
	helpers.ResponseNoContentJSON(c)
}

// GetAllRequestLinkCode func
func (controller UserController) GetAllRequestLinkCode(c *gin.Context) {
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	requests, errGetAll := controller.Service.GetAllRequestsLinkCode(params)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAllRequestLinkCode service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get request link code list successful", requests, params, len(requests))
}
