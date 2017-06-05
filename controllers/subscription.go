package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// SubscriptionController controller
type SubscriptionController struct {
	Service services.SubscriptionServiceInterface
}

// CreateSubscription func
func (controller SubscriptionController) CreateSubscription(c *gin.Context) {
	toID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
		return
	}

	existToID, errCheckExistUser := services.NewUserService().CheckExistUser(toID)
	if errCheckExistUser != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistUser service: %s\n", errCheckExistUser.Error())
		return
	}
	if existToID != true {
		helpers.ResponseBadRequestJSON(c, configs.EcAuthNoExistUser, "Don't exist user")
		return
	}
	//check permisson
	fromID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	// check exist subscriber
	existSub, errCheckExistSubscription := controller.Service.CheckExistSubscription(fromID, toID)
	if errCheckExistSubscription != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistSubscription service: %s\n", errCheckExistSubscription.Error())
		return
	}
	if existSub == true {
		helpers.ResponseBadRequestJSON(c, configs.EcExisObject, "Exist this subscription")
		return
	}

	subscriptionID, errCreateSubscription := controller.Service.CreateSubscription(fromID, toID)
	if errCreateSubscription != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CreateSubscription service: %s\n", errCreateSubscription)
		return
	}

	helpers.ResponseCreatedJSON(c, 1, "Create subscriber successful!", subscriptionID)

	// push noti
	go func() {
		Notification := NotificationController{Service: services.NewNotificationService()}
		Notification.Create(fromID, int(configs.IActionFollow), toID)
	}()
}

// DeleteSubscription func
func (controller SubscriptionController) DeleteSubscription(c *gin.Context) {
	toID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: id")
		return
	}

	//check permission
	fromID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
	}

	exist, errCheckExistSubscription := controller.Service.CheckExistSubscription(fromID, toID)
	if errCheckExistSubscription != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistSubscription service: %s\n", errCheckExistSubscription.Error())
		return
	}

	if exist != true {
		helpers.ResponseNotFoundJSON(c, 2, "No exist this subscription")
		return
	}

	deleted, errDeleteSubcription := controller.Service.DeleteSubcription(fromID, toID)
	if errDeleteSubcription == nil && deleted == true {

		// auto Decrease Followers And Followings
		go func() {
			Notification := NotificationController{Service: services.NewNotificationService()}
			Notification.Create(fromID, int(configs.IActionFollow), toID)
		}()

		helpers.ResponseJSON(c, 200, 1, "Delete subscriber successful", nil)
		return

	}
	if errDeleteSubcription != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("DeleteUserSubscriber: %s\n", errDeleteSubcription.Error())
	}
}

// GetFollowers func
func (controller SubscriptionController) GetFollowers(c *gin.Context) {
	userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
		return
	}

	exist, errCheckExistUser := services.NewUserService().CheckExistUser(userID)
	if errCheckExistUser != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistUser service: %s\n", errCheckExistUser.Error())
		return
	}
	if exist != true {
		helpers.ResponseBadRequestJSON(c, configs.EcAuthNoExistUser, "Don't exist user")
		return
	}

	userList, errGetFollowers := controller.Service.GetFollowers(userID)
	if errGetFollowers != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetFollowers service: %s\n", errGetFollowers.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get followers list successful", userList, nil, len(userList))
}

// GetSubcriptions func
func (controller SubscriptionController) GetSubcriptions(c *gin.Context) {
	userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
		return
	}

	exist, errCheckExistUser := services.NewUserService().CheckExistUser(userID)
	if errCheckExistUser != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistUser service: %s\n", errCheckExistUser.Error())
		return
	}
	if exist != true {
		helpers.ResponseBadRequestJSON(c, configs.EcAuthNoExistUser, "Don't exist user")
		return
	}

	userList, errGetFollowers := controller.Service.GetSubscriptions(userID)
	if errGetFollowers != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetFollowers service: %s\n", errGetFollowers.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get followers list successful", userList, nil, len(userList))
}
