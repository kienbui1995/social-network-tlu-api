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

// ChannelController controller
type ChannelController struct {
	Service services.ChannelServiceInterface
}

// GetAll func
func (controller ChannelController) GetAll(c *gin.Context) {

	// get my userid
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}

	// ParamsGetAll
	params := helpers.ParamsGetAll{}
	// params.Type = configs.SCanMention
	// params.Type = c.DefaultQuery("type", params.Type)
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)

	channels, errGetAll := controller.Service.GetAll(params, myUserID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get channels successful", channels, params, len(channels))
	return
}

// Get func
func (controller ChannelController) Get(c *gin.Context) {
	channelID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid channel_id")
	}
	// check exists
	exist, errCheckExistChannel := controller.Service.CheckExistChannel(channelID)
	if errCheckExistChannel != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistChannel service: %s\n", errCheckExistChannel.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist channel")
		return
	}

	// get myuserID
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken func: %s\n", errGetUserIDFromToken.Error())
		return
	}
	channel, errGet := controller.Service.Get(channelID, myUserID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
	}

	// check exist channel (404)
	if channel.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found channel")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Get channel successful", channel)
}

// Delete func
func (controller ChannelController) Delete(c *gin.Context) {
	channelID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid channel_id")
		return
	}
	// check exists
	exist, errCheckExistChannel := controller.Service.CheckExistChannel(channelID)
	if errCheckExistChannel != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistChannel service: %s\n", errCheckExistChannel.Error())
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist channel")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := controller.Service.CheckUserRole(channelID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IIsAdminChannel {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	// delete action
	deleted, errDelete := controller.Service.Delete(channelID)
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
func (controller ChannelController) Create(c *gin.Context) {
	// Check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	json := models.Channel{}
	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	if json.Status == 0 {
		json.Status = 1
	}

	channelID, errCreate := controller.Service.Create(json, myUserID)
	if errCreate == nil && channelID >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Create channel successful", map[string]interface{}{"id": channelID})
		return
	}
	helpers.ResponseServerErrorJSON(c)
	if errCreate != nil {
		fmt.Printf("Create services: %s\n", errCreate.Error())
	} else {
		fmt.Printf("Create services: Don't create channel")
	}
}

// Update func
func (controller ChannelController) Update(c *gin.Context) {
	channelID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: channel_id")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}

	// check exist channel
	channel, errGet := controller.Service.Get(channelID, myUserID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	if channel.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist channel")
		return
	}

	role, errCheckUserRole := controller.Service.CheckUserRole(channelID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IIsAdminChannel {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	json := models.InfoChannel{}
	if errBind := c.BindJSON(&json); errBind != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: %s\n"+errBind.Error())
		return
	}
	helpers.Replace(channel, &json)
	updated, errUpdate := controller.Service.Update(channelID, json)
	if errUpdate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: %s\n", errUpdate.Error())
		return
	}
	if updated.IsEmpty() {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: don't update")
		return
	}
	helpers.ResponseSuccessJSON(c, 1, "Update channel successful", updated)
}

// // GetUsers func
// func (controller ChannelController) GetUsers(c *gin.Context) {
// 	channelID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if errParseInt != nil {
// 		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: channelID")
// 		return
// 	}
//
// 	exist, errCheckExistChannel := controller.Service.CheckExistChannel(channelID)
// 	if errCheckExistChannel != nil {
// 		helpers.ResponseServerErrorJSON(c)
// 		fmt.Printf("CheckExistChannel service: %s\n", errCheckExistChannel.Error())
// 		return
// 	}
// 	if exist != true {
// 		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist object")
// 		return
// 	}
//
// 	//check permisson
// 	// myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
// 	// if errGetUserIDFromToken != nil {
// 	// 	helpers.ResponseServerErrorJSON(c)
// 	// 	fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
// 	// 	return
// 	// }
// 	// role, errCheckUserRole := controller.Service.CheckUserRole(channelID, myUserID)
// 	// if errCheckUserRole != nil {
// 	// 	helpers.ResponseServerErrorJSON(c)
// 	// 	fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
// 	// }
// 	// if role == configs.IBlocked {
// 	// 	helpers.ResponseForbiddenJSON(c, configs.EcPermissionGroup, "Group not visible")
// 	// 	return
// 	// }
// 	params := helpers.ParamsGetAll{}
// 	// params.Type = configs.SCanMention
// 	// params.Type = c.DefaultQuery("type", params.Type)
// 	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
// 	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
// 	params.Sort = c.DefaultQuery("sort", configs.SSort)
// 	params.Sort, _ = helpers.ConvertSort(params.Sort)
//
// 	// get members
// 	if params.Type == configs.SMember {
// 		users, errGetMembers := controller.Service.GetMembers(params, groupID)
// 		if errGetMembers != nil {
// 			helpers.ResponseServerErrorJSON(c)
// 			fmt.Printf("GetMembers service: %s\n", errGetMembers.Error())
// 			return
// 		}
// 		helpers.ResponseEntityListJSON(c, 1, "Get members successful", users, params, len(users))
// 		return
// 	}
//
// 	// get pending requested users
// 	if params.Type == configs.SPending {
// 		users, errGetPendingUsers := controller.Service.GetPendingUsers(params, groupID)
// 		if errGetPendingUsers != nil {
// 			helpers.ResponseServerErrorJSON(c)
// 			fmt.Printf("GetMentionedUsers service: %s\n", errGetPendingUsers.Error())
// 			return
// 		}
// 		helpers.ResponseEntityListJSON(c, 1, "Get mentioned users successful", users, params, len(users))
// 		return
// 	}
//
// 	// get liked users
// 	if params.Type == configs.SBlocked {
// 		users, errGetBlockedUsers := controller.Service.GetBlockedUsers(params, groupID)
// 		if errGetBlockedUsers != nil {
// 			helpers.ResponseServerErrorJSON(c)
// 			fmt.Printf("GetLikedUsers service: %s\n", errGetBlockedUsers.Error())
// 			return
// 		}
// 		helpers.ResponseEntityListJSON(c, 1, "Get liked users successful", users, params, len(users))
// 		return
// 	}
//
// }

// GetFollowers func
func (controller ChannelController) GetFollowers(c *gin.Context) {
	channelID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid channelID")
		return
	}

	// check exist channel
	exist, errCheckExistChannel := controller.Service.CheckExistChannel(channelID)
	if errCheckExistChannel != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistChannel service: %s\n", errCheckExistChannel.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcAuthNoExistUser, "No exist channel")
		return
	}

	// get my userid
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}

	// ParamsGetAll
	params := helpers.ParamsGetAll{}
	// params.Type = configs.SCanMention
	// params.Type = c.DefaultQuery("type", params.Type)
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)

	followers, errGetFollowers := controller.Service.GetFollowers(params, channelID, myUserID)
	if errGetFollowers != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetFollowers service: %s\n", errGetFollowers.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get followers of channel successful", followers, params, len(followers))
	return
}
