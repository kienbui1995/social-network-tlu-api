package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// GroupController controller
type GroupController struct {
	Service services.GroupServiceInterface
}

// GetAll func
func (controller GroupController) GetAll(c *gin.Context) {

}

// Get func
func (controller GroupController) Get(c *gin.Context) {

}

// Delete func
func (controller GroupController) Delete(c *gin.Context) {

}

// Create func
func (controller GroupController) Create(c *gin.Context) {

}

// Update func
func (controller GroupController) Update(c *gin.Context) {

}

// GetMembers func
func (controller GroupController) GetMembers(c *gin.Context) {

}

// GetRequests func
func (controller GroupController) GetRequests(c *gin.Context) {

}

// CreateRequest func
func (controller GroupController) CreateRequest(c *gin.Context) {

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

// CreateMember func
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
