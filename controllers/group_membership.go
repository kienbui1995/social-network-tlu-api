package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// GroupMembershipController controller
type GroupMembershipController struct {
	Service services.GroupMembershipServiceInterface
}

// GetAll func
func (controller GroupMembershipController) GetAll(c *gin.Context) {
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
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist object")
		return
	}

	//check permisson/ role
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
	if role == configs.IBlocked || role == configs.IDeclined || role == configs.IPending {
		helpers.ResponseForbiddenJSON(c, configs.EcPermissionGroup, "Group members not visible")
		return
	}

	// ParamsGetAll
	params := helpers.ParamsGetAll{}
	params.Type = configs.SCanMention
	params.Type = c.DefaultQuery("type", params.Type)
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)

	memberships, errGetAll := controller.Service.GetAll(params, groupID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get memberships successful", memberships, params, len(memberships))
	return
}

// Get func
func (controller GroupMembershipController) Get(c *gin.Context) {
	// membershipID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	// if errParseInt != nil {
	// 	helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: membership_id")
	// 	return
	// }
	//
	// // // check exist group
	// // exist, errCheckExistGroup := controller.Service.CheckExistGroupMembership(groupID, userID)
	// // if errCheckExistGroup != nil {
	// // 	helpers.ResponseServerErrorJSON(c)
	// // 	fmt.Printf("CheckExistGroup service: %s\n", errCheckExistGroup.Error())
	// // 	return
	// // }
	// // if exist != true {
	// // 	helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist object")
	// // 	return
	// // }
	//
	// //check permisson/ role
	// myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	// if errGetUserIDFromToken != nil {
	// 	helpers.ResponseServerErrorJSON(c)
	// 	fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
	// 	return
	// }
	// role, errCheckUserRole := services.NewGroupService().CheckUserRole(groupID, myUserID)
	// if errCheckUserRole != nil {
	// 	helpers.ResponseServerErrorJSON(c)
	// 	fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
	// }
	// if role == configs.IBlocked || role == configs.IDeclined || role == configs.IPending {
	// 	helpers.ResponseForbiddenJSON(c, configs.EcPermissionGroup, "Group members not visible")
	// 	return
	// }
}

// Create func
func (controller GroupMembershipController) Create(c *gin.Context) {
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

	//check permisson/ role
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
	if role == configs.IBlocked || role == configs.IDeclined || role == configs.IPending {
		helpers.ResponseForbiddenJSON(c, configs.EcPermissionGroup, "Group not visible")
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

	membershipID, errCreate := controller.Service.Create(groupID, myUserID)
	if errCreate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Create service: %s\n", errCreate.Error())
	}
	if membershipID >= 0 {
		helpers.ResponseJSON(c, 200, 1, "create membership successful", nil)
	}
}

// Update func
func (controller GroupMembershipController) Update(c *gin.Context) {
	membershipID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: membership_id")
		return
	}

	// check exist membership
	membership, errGet := controller.Service.Get(membershipID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistGroup service: %s\n", errGet.Error())
		return
	}
	if membership.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist membership")
		return
	}

	//check permisson/ role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := services.NewGroupService().CheckUserRole(membership.Group.ID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	json := struct {
		Status int `json:"status"`
	}{}
	if errBind := c.BindJSON(&json); errBind != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: %s\n"+errBind.Error())
		return
	}
	membership.Status = json.Status
	updatedMembership, errUpdate := controller.Service.Update(membership)
	if errUpdate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: %s\n", errUpdate.Error())
		return
	}
	if updatedMembership.IsEmpty() {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: don't update")
		return
	}
	helpers.ResponseSuccessJSON(c, 1, "Update membership successful", updatedMembership)
}

// Delete func
func (controller GroupMembershipController) Delete(c *gin.Context) {
	membershipID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: membership_id")
		return
	}

	// check exist group
	membership, errGet := controller.Service.Get(membershipID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
		return
	}
	if membership.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist membership")
		return
	}

	//check permisson/ role
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	role, errCheckUserRole := services.NewGroupService().CheckUserRole(membership.Group.ID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}
	if role != configs.IAdmin {
		helpers.ResponseForbiddenJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	deleted, errDelete := controller.Service.Delete(membershipID)
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
func (controller GroupMembershipController) DeleteByUser(c *gin.Context) {
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

	//check permisson/ role
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
	if role != configs.IMember {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist membership")
		return
	}

	deleted, errDeleteByUser := controller.Service.DeleteByUser(groupID, myUserID)
	if errDeleteByUser != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("DeleteByUser service: %s\n", errDeleteByUser.Error())
		return
	}
	if deleted != true {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("DeleteByUser service: don't delete")
		return
	}
	helpers.ResponseNoContentJSON(c)
}
