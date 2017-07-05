package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// ViolationController controller
type ViolationController struct {
	Service services.ViolationServiceInterface
}

// GetAll func
func (controller ViolationController) GetAll(c *gin.Context) {
	// userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	// if errParseInt != nil {
	// 	helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
	// 	return
	// }

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	role, errGetRoleFromUserID := GetRoleFromUserID(myUserID)
	if errGetRoleFromUserID != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		fmt.Printf("GetRoleFromUserID service: %s\n", errGetRoleFromUserID.Error())
		return
	}
	if role != configs.IAdminRole && role != configs.ISupervisorRole {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	violations, errGetAll := controller.Service.GetAll(params)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get violations successful", violations, params, len(violations))
}

// Get func
func (controller ViolationController) Get(c *gin.Context) {
	violationID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid channel violation id")
		return
	}
	// check exists
	exist, errCheckExistViolation := controller.Service.CheckExistViolation(violationID)
	if errCheckExistViolation != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistViolation service: %s\n", errCheckExistViolation.Error())
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist violation")
		return
	}

	// // get myuserID
	// myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	// if errGetUserIDFromToken != nil {
	// 	helpers.ResponseServerErrorJSON(c)
	// 	fmt.Printf("GetUserIDFromToken func: %s\n", errGetUserIDFromToken.Error())
	// }
	violation, errGet := controller.Service.Get(violationID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
	}

	// Valida se existe a pessoa (404)
	if violation.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found violation")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Info of a violation", violation)
}

// Delete func
func (controller ViolationController) Delete(c *gin.Context) {
	violationID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid violationID")
		return
	}
	// check exists
	exist, errCheckExistViolation := controller.Service.CheckExistViolation(violationID)
	if errCheckExistViolation != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistViolation service: %s\n", errCheckExistViolation.Error())
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist violation")
		return
	}

	// get myUserID from token
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	// check permisson
	allowed, errCheckPermission := controller.Service.CheckPermission(violationID, myUserID)
	if errCheckPermission != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckPermission service: %s\n", errCheckPermission.Error())
		return
	}
	if allowed == false {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	// delete action
	deleted, errDelete := controller.Service.Delete(violationID)
	if errDelete == nil && deleted == true {
		helpers.ResponseNoContentJSON(c)
		return
	}
	if errDelete != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: %s\n", errDelete)
	}
}

// Create func
func (controller ViolationController) Create(c *gin.Context) {
	studentCode := c.Param("id")

	// get myUserID
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if myUserID < 0 || errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	supervisiorID, errGetSupervisiorCode := controller.Service.GetSupervisiorID(myUserID)
	if errGetSupervisiorCode != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		fmt.Printf("GetSupervisiorCode service: %s\n", errGetSupervisiorCode.Error())
		return
	}
	if supervisiorID < 0 {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	// json struct
	json := struct {
		Message string `json:"message"`
		Photo   string `json:"photo"`
		TimeAt  int64  `json:"time_at"`
		Place   string `json:"place"`
		Status  int    `json:"status"`
	}{}

	// BadRequest (400)
	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	// validation
	if len(json.Message) == 0 {
		helpers.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
		return
	}
	if json.Status == 0 {
		json.Status = 1
	}

	violation := models.Violation{}
	helpers.Replace(json, &violation)
	violationID, errCreate := controller.Service.Create(violation, studentCode, supervisiorID)
	if errCreate == nil && violationID >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Create violation successful", map[string]interface{}{"id": violationID})
		return
		// // auto noti
		// go func() {
		//
		// }()

		// push noti
		// go func() {
		// 	Notification := NotificationController{Service: services.NewNotificationService()}
		// 	err := Notification.UpdatePostNotification(userID)
		// 	if err != nil {
		// 		fmt.Printf("UpdateLikeNotification: %s\n", err.Error())
		// 	}
		// }()
		// return
	}
	helpers.ResponseServerErrorJSON(c)
	if errCreate != nil {
		fmt.Printf("Create services: %s\n", errCreate.Error())
	} else {
		fmt.Printf("Create services: Don't create violation")
	}
}

// Update func
func (controller ViolationController) Update(c *gin.Context) {
	violationID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	exist, errCheckExistViolation := controller.Service.CheckExistViolation(violationID)
	if errCheckExistViolation != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistViolation service: %s\n", errCheckExistViolation.Error())
		return
	}
	if exist == false {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Don't exist violation")
		return
	}

	// get myUserID from token
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil || myUserID < 0 {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	role, errGetRoleFromUserID := GetRoleFromUserID(myUserID)
	if errGetRoleFromUserID != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		fmt.Printf("GetRoleFromUserID service: %s\n", errGetRoleFromUserID.Error())
		return
	}
	if role != configs.IAdminRole && role != configs.ISupervisorRole {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	olderViolation, errGet := controller.Service.Get(violationID)

	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	if olderViolation.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found Violation")
		return
	}
	var newViolation models.Violation
	if errBindJSON := c.BindJSON(&newViolation); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	errReplace := helpers.Replace(olderViolation, &newViolation)

	if errReplace != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Replace helpers: %s\n", errReplace.Error())
	}

	violation, errUpdate := controller.Service.Update(newViolation)
	if errUpdate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: %s\n", errUpdate.Error())
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Update violationsuccessful", violation)
}

// GetAllOfStudent func
func (controller ViolationController) GetAllOfStudent(c *gin.Context) {
	studentCode := c.Param("id")

	//check permisson
	// myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	// if errGetUserIDFromToken != nil {
	// 	helpers.ResponseAuthJSON(c, 200, "Permissions error")
	// 	return
	// }
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	violations, errGetAllOfStudent := controller.Service.GetAllOfStudent(params, studentCode)
	if errGetAllOfStudent != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAllOfStudent service: %s\n", errGetAllOfStudent.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get violation list successful", violations, params, len(violations))
}

// GetAllOfSupervisior func
func (controller ViolationController) GetAllOfSupervisior(c *gin.Context) {
	supervisiorCode := c.Param("id")

	//check permisson
	// myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	// if errGetUserIDFromToken != nil {
	// 	helpers.ResponseAuthJSON(c, 200, "Permissions error")
	// 	return
	// }
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	violations, errGetAllOfSupervisior := controller.Service.GetAllOfSupervisior(params, supervisiorCode)
	if errGetAllOfSupervisior != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAllOfSupervisior service: %s\n", errGetAllOfSupervisior.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get violation list successful", violations, params, len(violations))
}
