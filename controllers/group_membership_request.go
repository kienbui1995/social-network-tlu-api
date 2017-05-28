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

// GroupMembershipRequestController controller
type GroupMembershipRequestController struct {
	Service services.GroupMembershipRequestServiceInterface
}

// Create func
func (controller GroupMembershipRequestController) Create(c *gin.Context) {
	groupID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: group_id")
		return
	}

	// check exist group
	exist, errCheckExistGroup := services.NewGroupService().CheckExistGroup(groupID)
	if errCheckExistGroup != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistGroup service: %s\n", errCheckExistGroup.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist group")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := services.NewGroupService().CheckUserRole(groupID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.ICanRequest {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permission error")
		return
	}

	if role == configs.IAdmin {
		helpers.ResponseBadRequestJSON(c, configs.EcExisObject, "Exist membership")
		return
	}
	if role == configs.IMember {
		helpers.ResponseBadRequestJSON(c, configs.EcExisObject, "Exist membership")
		return
	}

	if role == configs.ICanJoin {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Can't create this group membership request")
		return
	}

	json := struct {
		Message string `json:"message,omitempty"`
	}{}
	if errBind := c.BindJSON(&json); errBind != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: %s\n"+errBind.Error())
		return
	}
	created, errCreate := controller.Service.Create(groupID, myUserID, json.Message)
	if errCreate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Create service: %s\n", errCreate.Error())
	}
	if created {
		helpers.ResponseJSON(c, 200, 1, "create group membership request successful", nil)
	}
}

// GetAll func
func (controller GroupMembershipRequestController) GetAll(c *gin.Context) {
	groupID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: group_id")
		return
	}

	// check exist group
	exist, errCheckExistGroup := services.NewGroupService().CheckExistGroup(groupID)
	if errCheckExistGroup != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistGroup service: %s\n", errCheckExistGroup.Error())
		return
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
	role, errCheckUserRole := services.NewGroupService().CheckUserRole(groupID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permission error")
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

	requests, errGetAll := controller.Service.GetAll(params, groupID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get group membership requests successful", requests, params, len(requests))
	return
}

// Get func
func (controller GroupMembershipRequestController) Get(c *gin.Context) {
	requestID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: group_membership_request_id")
		return
	}

	// check exist group_membership_request/get request
	request, errGet := controller.Service.Get(requestID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	if request.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist object")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := services.NewGroupService().CheckUserRole(request.Group.ID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin && myUserID != request.User.ID {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permission error")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Get group membership request successful", request)
}

// Update func
func (controller GroupMembershipRequestController) Update(c *gin.Context) {
	requestID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: group_membership_request_id")
		return
	}

	// check exist group_membership_request/get request
	request, errGet := controller.Service.Get(requestID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	if request.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist group membership request")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := services.NewGroupService().CheckUserRole(request.Group.ID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	json := models.GroupMembershipRequest{}
	if errBind := c.BindJSON(&json); errBind != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: %s\n"+errBind.Error())
		return
	}
	helpers.Replace(request, &json)
	updatedRequest, errUpdate := controller.Service.Update(requestID, json)
	if errUpdate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: %s\n", errUpdate.Error())
		return
	}
	if updatedRequest.IsEmpty() {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: don't update")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Update group membership request successful", updatedRequest)
}

// Delete func
func (controller GroupMembershipRequestController) Delete(c *gin.Context) {
	requestID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: group_membership_request_id")
		return
	}

	// check exist group/get group_membership_request
	request, errGet := controller.Service.Get(requestID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	if request.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist group_membership_request")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := services.NewGroupService().CheckUserRole(request.Group.ID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin && myUserID != request.User.ID {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	deleted, errDelete := controller.Service.Delete(requestID)
	if errDelete != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: %s\n", errDelete.Error())
		return
	}
	if deleted != true {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: don't delete")
		return
	}
	helpers.ResponseNoContentJSON(c)
}

// DeleteByUser func
func (controller GroupMembershipRequestController) DeleteByUser(c *gin.Context) {
	requestID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: group_membership_request_id")
		return
	}

	// check exist group/get group_membership_request
	request, errGet := controller.Service.Get(requestID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	if request.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist group_membership_request")
		return
	}

	//check permisson/role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := services.NewGroupService().CheckUserRole(request.Group.ID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin && myUserID != request.User.ID {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	deleted, errDelete := controller.Service.Delete(requestID)
	if errDelete != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: %s\n", errDelete.Error())
		return
	}
	if deleted != true {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: don't delete")
		return
	}
	helpers.ResponseNoContentJSON(c)
}
