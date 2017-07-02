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

// ChannelNotificationController controller
type ChannelNotificationController struct {
	Service services.ChannelNotificationServiceInterface
}

// GetAll func
func (controller ChannelNotificationController) GetAll(c *gin.Context) {
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
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	channelNotifications, errGetAll := controller.Service.GetAll(params, myUserID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get channel notification list successful", channelNotifications, params, len(channelNotifications))
}

// Get func
func (controller ChannelNotificationController) Get(c *gin.Context) {
	channelNotificationID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid channel notification id")
		return
	}
	// check exists
	exist, errCheckExistChannelNotification := controller.Service.CheckExistChannelNotification(channelNotificationID)
	if errCheckExistChannelNotification != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistChannelNotification service: %s\n", errCheckExistChannelNotification.Error())
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist channel notification")
		return
	}

	// // get myuserID
	// myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	// if errGetUserIDFromToken != nil {
	// 	helpers.ResponseServerErrorJSON(c)
	// 	fmt.Printf("GetUserIDFromToken func: %s\n", errGetUserIDFromToken.Error())
	// }
	channelNotification, errGet := controller.Service.Get(channelNotificationID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
	}

	// Valida se existe a pessoa (404)
	if channelNotification.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found Channel notification")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Info of a channel notification", channelNotification)
}

// Delete func
func (controller ChannelNotificationController) Delete(c *gin.Context) {
	channelNotificationID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid channelNotificationID")
		return
	}
	// check exists
	exist, errCheckExistChannelNotification := controller.Service.CheckExistChannelNotification(channelNotificationID)
	if errCheckExistChannelNotification != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistChannelNotification service: %s\n", errCheckExistChannelNotification.Error())
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist channel notification")
		return
	}

	// get myUserID from token
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	// check permisson
	allowed, errCheckPermissionByNotificationID := controller.Service.CheckPermissionByNotificationID(channelNotificationID, myUserID)
	if errCheckPermissionByNotificationID != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckPermissionByNotificationID service: %s\n", errCheckPermissionByNotificationID.Error())
		return
	}
	if allowed == false {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	// delete action
	deleted, errDelete := controller.Service.Delete(channelNotificationID)
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
func (controller ChannelNotificationController) Create(c *gin.Context) {
	channelID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid channel id")
	}

	// get myUserID
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if myUserID < 0 || errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	// check permisson
	allowed, errCheckPermission := controller.Service.CheckPermission(channelID, myUserID)
	if errCheckPermission != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckPermission service: %s\n", errCheckPermission.Error())
		return
	}
	if allowed == false {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	// json struct
	json := struct {
		Title   string `json:"title"`
		Message string `json:"message"`
		Photo   string `json:"photo"`
		Time    string `json:"time"`
		Place   string `json:"place"`
		Status  int    `json:"status"`
	}{}

	// BadRequest (400)
	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	// validation
	if len(json.Title) == 0 {
		helpers.ResponseJSON(c, 400, 100, "Missing a few fields:  Title is NULL", nil)
		return
	}
	if len(json.Message) == 0 {
		helpers.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
		return
	}
	if json.Status == 0 {
		json.Status = 1
	}

	channelNotification := models.ChannelNotification{}
	helpers.Replace(json, &channelNotification)
	channelNotificationID, errCreate := controller.Service.Create(channelNotification, channelID)
	if errCreate == nil && channelNotificationID >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Create channel notification successful", map[string]interface{}{"id": channelNotificationID})

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
		fmt.Printf("Create services: Don't create channel notification")
	}
}

// Update func
func (controller ChannelNotificationController) Update(c *gin.Context) {
	channelNotificationID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	exist, errCheckExistChannelNotification := controller.Service.CheckExistChannelNotification(channelNotificationID)
	if errCheckExistChannelNotification != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistChannelNotification service: %s\n", errCheckExistChannelNotification.Error())
		return
	}
	if exist == false {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Don't exist channel notification")
		return
	}

	// get myUserID from token
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil || myUserID < 0 {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	olderChannelNotification, errGet := controller.Service.Get(channelNotificationID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}

	var newChannelNotification models.ChannelNotification
	if errBindJSON := c.BindJSON(&newChannelNotification); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	errReplace := helpers.Replace(olderChannelNotification, &newChannelNotification)

	if errReplace != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Replace helpers: %s\n", errReplace.Error())
	}

	channelNotification, errUpdate := controller.Service.Update(newChannelNotification)
	if errUpdate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: %s\n", errUpdate.Error())
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Update channel notification successful", channelNotification)
}

// GetAllOfChannel func
func (controller ChannelNotificationController) GetAllOfChannel(c *gin.Context) {
	channelID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid channel id")
		return
	}

	//check permisson
	// myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	// if errGetUserIDFromToken != nil {
	// 	helpers.ResponseAuthJSON(c, 200, "Permissions error")
	// 	return
	// }
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Type = c.DefaultQuery("type", configs.SPost)
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	channelNotifications, errGetAllOfChannel := controller.Service.GetAllOfChannel(params, channelID)
	if errGetAllOfChannel != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAllOfChannel service: %s\n", errGetAllOfChannel.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get channel notification list successful", channelNotifications, params, len(channelNotifications))
}
