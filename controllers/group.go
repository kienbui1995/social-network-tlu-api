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

// GroupController controller
type GroupController struct {
	Service services.GroupServiceInterface
}

// GetAll func
func (controller GroupController) GetAll(c *gin.Context) {

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

	groups, errGetAll := controller.Service.GetAll(params, myUserID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get groups successful", groups, params, len(groups))
	return
}

// Get func
func (controller GroupController) Get(c *gin.Context) {
	groupID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid group_id")
	}
	// check exists
	exist, errCheckExistGroup := controller.Service.CheckExistGroup(groupID)
	if errCheckExistGroup != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistGroup service: %s\n", errCheckExistGroup.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist group")
		return
	}

	// get myuserID
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken func: %s\n", errGetUserIDFromToken.Error())
		return
	}
	group, errGet := controller.Service.Get(groupID, myUserID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
	}

	// Valida se existe a pessoa (404)
	if group.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found Group")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Get group successful", group)
}

// Delete func
func (controller GroupController) Delete(c *gin.Context) {
	groupID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid group_id")
		return
	}
	// check exists
	exist, errCheckExistGroup := controller.Service.CheckExistGroup(groupID)
	if errCheckExistGroup != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistGroup service: %s\n", errCheckExistGroup.Error())
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist group")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := controller.Service.CheckUserRole(groupID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	// delete action
	deleted, errDelete := controller.Service.Delete(groupID)
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
func (controller GroupController) Create(c *gin.Context) {
	// Check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	json := models.Group{}
	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	if json.Status == 0 {
		json.Status = 1
	}
	if json.Privacy == 0 {
		json.Privacy = 1
	}

	groupID, errCreate := controller.Service.Create(json, myUserID)
	if errCreate == nil && groupID >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Create group successful", map[string]interface{}{"id": groupID})
		return
	}
	helpers.ResponseServerErrorJSON(c)
	if errCreate != nil {
		fmt.Printf("Create services: %s\n", errCreate.Error())
	} else {
		fmt.Printf("Create services: Don't create group")
	}
}

// Update func
func (controller GroupController) Update(c *gin.Context) {
	groupID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: group_id")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}

	// check exist membership
	group, errGet := controller.Service.Get(groupID, myUserID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	if group.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist group")
		return
	}

	role, errCheckUserRole := services.NewGroupService().CheckUserRole(groupID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	json := models.InfoGroup{}
	if errBind := c.BindJSON(&json); errBind != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: %s\n"+errBind.Error())
		return
	}
	helpers.Replace(group, &json)
	updatedGroup, errUpdate := controller.Service.Update(groupID, json)
	if errUpdate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: %s\n", errUpdate.Error())
		return
	}
	if updatedGroup.IsEmpty() {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: don't update")
		return
	}
	helpers.ResponseSuccessJSON(c, 1, "Update group successful", updatedGroup)
}

// GetReports func
func (controller GroupController) GetReports(c *gin.Context) {

}

// CreateReport func
func (controller GroupController) CreateReport(c *gin.Context) {

}

// GetPosts func
func (controller GroupController) GetPosts(c *gin.Context) {

}

// CreatePost func
func (controller GroupController) CreatePost(c *gin.Context) {

}

// CreateMember func
func (controller GroupController) CreateMember(c *gin.Context) {

}

// GetUsers func
func (controller GroupController) GetUsers(c *gin.Context) {
	groupID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: group_id")
		return
	}

	exist, errCheckExistGroup := controller.Service.CheckExistGroup(groupID)
	if errCheckExistGroup != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistGroup service: %s\n", errCheckExistGroup.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist object")
		return
	}

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := controller.Service.CheckUserRole(groupID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
	}
	if role == configs.IBlocked || role == configs.IDeclined {
		helpers.ResponseForbiddenJSON(c, configs.EcPermissionGroup, "Group not visible")
		return
	}
	params := helpers.ParamsGetAll{}
	params.Type = configs.SCanMention
	params.Type = c.DefaultQuery("type", params.Type)
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)

	// get members
	if params.Type == configs.SMember {
		users, errGetMembers := controller.Service.GetMembers(params, groupID)
		if errGetMembers != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("GetMembers service: %s\n", errGetMembers.Error())
			return
		}
		helpers.ResponseEntityListJSON(c, 1, "Get members successful", users, params, len(users))
		return
	}

	// get pending requested users
	if params.Type == configs.SPending {
		users, errGetPendingUsers := controller.Service.GetPendingUsers(params, groupID)
		if errGetPendingUsers != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("GetMentionedUsers service: %s\n", errGetPendingUsers.Error())
			return
		}
		helpers.ResponseEntityListJSON(c, 1, "Get mentioned users successful", users, params, len(users))
		return
	}

	// get liked users
	if params.Type == configs.SBlocked {
		users, errGetBlockedUsers := controller.Service.GetBlockedUsers(params, groupID)
		if errGetBlockedUsers != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("GetLikedUsers service: %s\n", errGetBlockedUsers.Error())
			return
		}
		helpers.ResponseEntityListJSON(c, 1, "Get liked users successful", users, params, len(users))
		return
	}

}

// GetJoinedGroup func
func (controller GroupController) GetJoinedGroup(c *gin.Context) {
	userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user_id")
		return
	}

	// check exist membership
	exist, errCheckExistUser := services.NewUserService().CheckExistUser(userID)
	if errCheckExistUser != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("errCheckExistUser service: %s\n", errCheckExistUser.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcAuthNoExistUser, "No exist user")
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

	groups, errGetJoinedGroup := controller.Service.GetJoinedGroup(params, userID, myUserID)
	if errGetJoinedGroup != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetJoinedGroup service: %s\n", errGetJoinedGroup.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get joined groups successful", groups, params, len(groups))
	return
}
