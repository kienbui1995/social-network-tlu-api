package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
	"github.com/kienbui1995/social-network-tlu-api/services"
	"github.com/maddevsio/fcm"
)

// NotificationController controller
type NotificationController struct {
	Service services.NotificationServiceInterface
}

// GetAll func
func (controller NotificationController) GetAll(c *gin.Context) {

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Type = c.DefaultQuery("type", configs.SNotiAll)

	// params.Type = configs.SNotiAll
	//
	// if params.Type != configs.SNotiComment && params.Type != configs.SNotiFollow && params.Type != configs.SNotiLike && params.Type != configs.SNotiMention && params.Type != configs.SNotiPhoto && params.Type != configs.SNotiStatus && params.Type != configs.SPost {
	// 	helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: type")
	// 	return
	// }
	params.Sort = c.DefaultQuery("sort", "updated_at DESC")
	notifications, errGetAll := controller.Service.GetAll(params, myUserID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get notifications successful", notifications, params, len(notifications))
}

// UpdateSeenNotification func
func (controller NotificationController) UpdateSeenNotification(c *gin.Context) {
	notificationID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid notification id")
		return
	}

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	if seen, _ := controller.Service.CheckSeenNotification(notificationID, myUserID); seen == true {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid notification id")
		return
	}
	seen, errSeenNotification := controller.Service.SeenNotification(notificationID, myUserID)
	if errSeenNotification != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("SeenNotification service: %s\n", errSeenNotification.Error())
	}
	if seen == true {
		helpers.ResponseJSON(c, 200, 1, "Seen notification successful", nil)
		return
	}

}

// Create func
func (controller NotificationController) Create(actorID int64, action int, objectID int64) (bool, error) {
	// if action == configs.IActionLike {
	// 	actionLike, errGetActionLike := controller.GetActionLike(objectID)
	// 	if errGetActionLike != nil {
	// 		fmt.Printf("GetActionLike in create notification controller: %s\n", errGetActionLike.Error())
	// 		return false, errGetActionLike
	// 	}
	// 	notification, errUpdateLikeNotification := controller.Service.UpdateLikeNotification(actionLike)
	// 	if errUpdateLikeNotification != nil {
	// 		fmt.Printf("UpdateLikeNotification service in create notification controller: %s\n", errUpdateLikeNotification.Error())
	// 		return false, errUpdateLikeNotification
	// 	}
	// 	userIDs, errGetSubcriberNotification := controller.Service.GetSubcriberNotification(notification.ID)
	// 	if errGetSubcriberNotification != nil {
	// 		fmt.Printf("GetSubcriberNotification service in create notification controller: %s\n", errGetSubcriberNotification.Error())
	// 		return false, errGetSubcriberNotification
	// 	}
	// 	PushTest(notification, userIDs)
	// } else if action == configs.IActionComment {
	// 	actionComment, errGetActionComment := controller.GetActionComment(objectID)
	// 	if errGetActionComment != nil {
	// 		fmt.Printf("GetActionComment in create notification controller: %s\n", errGetActionComment.Error())
	// 		return false, errGetActionComment
	// 	}
	// 	notification, errUpdateCommentNotification := controller.Service.UpdateCommentNotification(actionComment)
	// 	if errUpdateCommentNotification != nil {
	// 		fmt.Printf("UpdateCommentNotification service in create notification controller: %s\n", errUpdateCommentNotification.Error())
	// 		return false, errUpdateCommentNotification
	// 	}
	// 	userIDs, errGetSubcriberNotification := controller.Service.GetSubcriberNotification(notification.ID)
	// 	if errGetSubcriberNotification != nil {
	// 		fmt.Printf("GetSubcriberNotification service in create notification controller: %s\n", errGetSubcriberNotification.Error())
	// 		return false, errGetSubcriberNotification
	// 	}
	// 	PushTest(notification, userIDs)
	// } else if action == configs.IActionFollow {
	// 	actionFollow, errGetActionFollow := controller.GetActionFollow(actorID)
	// 	if errGetActionFollow != nil {
	// 		fmt.Printf("GetActionFollow in create notification controller: %s\n", errGetActionFollow.Error())
	// 		return false, errGetActionFollow
	// 	}
	// 	notification, errUpdateFollowNotification := controller.Service.UpdateFollowNotification(actionFollow)
	// 	if errUpdateFollowNotification != nil {
	// 		fmt.Printf("UpdateFollowNotification service in create notification controller: %s\n", errUpdateFollowNotification.Error())
	// 		return false, errUpdateFollowNotification
	// 	}
	// 	userIDs, errGetSubcriberNotification := controller.Service.GetSubcriberNotification(notification.ID)
	// 	if errGetSubcriberNotification != nil {
	// 		fmt.Printf("GetSubcriberNotification service in create notification controller: %s\n", errGetSubcriberNotification.Error())
	// 		return false, errGetSubcriberNotification
	// 	}
	// 	PushTest(notification, userIDs)
	// } else if action == configs.IActionPostStatus {
	// 	actionPostStatus, errGetActionPostStatus := controller.GetActionPostStatus(objectID)
	// 	if errGetActionPostStatus != nil {
	// 		fmt.Printf("GetActionPostStatus in create notification controller: %s\n", errGetActionPostStatus.Error())
	// 		return false, errGetActionPostStatus
	// 	}
	// 	notification, errUpdateStatusNotification := controller.Service.UpdateStatusNotification(actionPostStatus)
	// 	if errUpdateStatusNotification != nil {
	// 		fmt.Printf("UpdateStatusNotification service in create notification controller: %s\n", errUpdateStatusNotification.Error())
	// 		return false, errUpdateStatusNotification
	// 	}
	// 	userIDs, errGetSubcriberNotification := controller.Service.GetSubcriberNotification(notification.ID)
	// 	if errGetSubcriberNotification != nil {
	// 		fmt.Printf("GetSubcriberNotification service in create notification controller: %s\n", errGetSubcriberNotification.Error())
	// 		return false, errGetSubcriberNotification
	// 	}
	// 	PushTest(notification, userIDs)
	// } else if action == configs.IActionPostPhoto {
	// 	actionPostPhoto, errGetActionPostPhoto := controller.GetActionPostPhoto(objectID)
	// 	if errGetActionPostPhoto != nil {
	// 		fmt.Printf("GetActionPostPhoto in create notification controller: %s\n", errGetActionPostPhoto.Error())
	// 		return false, errGetActionPostPhoto
	// 	}
	// 	notification, errUpdatePhotoNotification := controller.Service.UpdatePhotoNotification(actionPostPhoto)
	// 	if errUpdatePhotoNotification != nil {
	// 		fmt.Printf("UpdatePhotoNotification service in create notification controller: %s\n", errUpdatePhotoNotification.Error())
	// 		return false, errUpdatePhotoNotification
	// 	}
	// 	userIDs, errGetSubcriberNotification := controller.Service.GetSubcriberNotification(notification.ID)
	// 	if errGetSubcriberNotification != nil {
	// 		fmt.Printf("GetSubcriberNotification service in create notification controller: %s\n", errGetSubcriberNotification.Error())
	// 		return false, errGetSubcriberNotification
	// 	}
	// 	PushTest(notification, userIDs)
	// } else if action == configs.IActionMention {
	// 	actionMention, errGetActionMention := controller.GetActionMention(objectID, actorID)
	// 	if errGetActionMention != nil {
	// 		fmt.Printf("GetActionMention in create notification controller: %s\n", errGetActionMention.Error())
	// 		return false, errGetActionMention
	// 	}
	// 	notification, errUpdateMentionNotification := controller.Service.UpdateMentionNotification(actionMention)
	// 	if errUpdateMentionNotification != nil {
	// 		fmt.Printf("UpdateMentionNotification service in create notification controller: %s\n", errUpdateMentionNotification.Error())
	// 		return false, errUpdateMentionNotification
	// 	}
	// 	userIDs, errGetSubcriberNotification := controller.Service.GetSubcriberNotification(notification.ID)
	// 	if errGetSubcriberNotification != nil {
	// 		fmt.Printf("GetSubcriberNotification service in create notification controller: %s\n", errGetSubcriberNotification.Error())
	// 		return false, errGetSubcriberNotification
	// 	}
	// 	PushTest(notification, userIDs)
	// }
	return true, nil
}

// UpdateLikeNotification func
func (controller NotificationController) UpdateLikeNotification(postID int64, userID int64) error {
	notify1, errUpdateLikeNotification := controller.Service.UpdateLikeNotification(postID)
	if errUpdateLikeNotification != nil {
		return errUpdateLikeNotification
	}
	userIDs1, errGetNotificationSubcriber1 := controller.Service.GetNotificationSubcriber(notify1.ID)
	if errGetNotificationSubcriber1 != nil {
		return errGetNotificationSubcriber1
	}

	errPushTest1 := controller.PushTest(notify1, userIDs1)

	if errPushTest1 != nil {
		return errPushTest1
	}
	_, errUpdateLikedPostNotification := controller.Service.UpdateLikedPostNotification(userID)

	if errUpdateLikedPostNotification != nil {
		fmt.Printf("Pushnoti1:%s\n", errUpdateLikedPostNotification.Error())
		return errUpdateLikedPostNotification
	}

	// Only creeate new noti, no push notification on device

	// userIDs2, errGetNotificationSubcriber2 := controller.Service.GetNotificationSubcriber(notify2.ID)
	// if errGetNotificationSubcriber2 != nil {
	// 	return errGetNotificationSubcriber2
	// }
	// fmt.Printf("Pushnoti2:%s\n", errGetNotificationSubcriber2.Error())
	// errPushTest2 := controller.PushTest(notify2, userIDs2)
	// if errPushTest2 != nil {
	// 	return errPushTest2
	// }
	return nil
}

// UpdateCommentNotification func
func (controller NotificationController) UpdateCommentNotification(postID int64, userID int64) error {
	notify1, errUpdateCommentNotification := controller.Service.UpdateCommentNotification(postID)
	if errUpdateCommentNotification != nil {
		fmt.Printf("errUpdateCommentNotification: %s\n", errUpdateCommentNotification.Error())
		return errUpdateCommentNotification
	}
	userIDs1, errGetNotificationSubcriber1 := controller.Service.GetNotificationSubcriber(notify1.ID)
	if errGetNotificationSubcriber1 != nil {
		return errGetNotificationSubcriber1
	}
	errPushTest1 := controller.PushTest(notify1, userIDs1)
	if errPushTest1 != nil {
		return errPushTest1
	}
	_, errUpdateCommentedPostNotification := controller.Service.UpdateCommentedPostNotification(userID)
	if errUpdateCommentedPostNotification != nil {
		fmt.Printf("errUpdateCommentedPostNotification: %s\n", errUpdateCommentedPostNotification.Error())
		return errUpdateCommentedPostNotification
	}
	// Only creeate new noti, no push notification on device

	// userIDs2, errGetNotificationSubcriber2 := controller.Service.GetNotificationSubcriber(notify2.ID)
	// if errGetNotificationSubcriber2 != nil {
	// 	return errGetNotificationSubcriber2
	// }
	// errPushTest2 := controller.PushTest(notify2, userIDs2)
	// if errPushTest2 != nil {
	// 	return errPushTest2
	// }
	return nil
}

// UpdateMentionNotification func
func (controller NotificationController) UpdateMentionNotification(postID int64, userID int64, commentID int64) error {
	notify1, errUpdateCommentNotification := controller.Service.UpdateCommentNotification(postID)
	if errUpdateCommentNotification != nil {
		return errUpdateCommentNotification
	}
	userIDs1, errGetNotificationSubcriber1 := controller.Service.GetNotificationSubcriber(notify1.ID)
	if errGetNotificationSubcriber1 != nil {
		return errGetNotificationSubcriber1
	}
	errPushTest1 := controller.PushTest(notify1, userIDs1)
	if errPushTest1 != nil {
		return errPushTest1
	}
	notify2, errUpdateCommentedPostNotification := controller.Service.UpdateCommentedPostNotification(userID)
	if errUpdateCommentedPostNotification != nil {
		return errUpdateCommentedPostNotification
	}
	userIDs2, errGetNotificationSubcriber2 := controller.Service.GetNotificationSubcriber(notify2.ID)
	if errGetNotificationSubcriber2 != nil {
		return errGetNotificationSubcriber2
	}
	errPushTest2 := controller.PushTest(notify2, userIDs2)
	if errPushTest2 != nil {
		return errPushTest2
	}
	return nil
}

// UpdateFollowNotification func
func (controller NotificationController) UpdateFollowNotification(userID int64, objectID int64) error {
	notify1, errUpdateFollowNotification := controller.Service.UpdateFollowNotification(userID, objectID)
	if errUpdateFollowNotification != nil {
		return errUpdateFollowNotification
	}
	userIDs1, errGetNotificationSubcriber1 := controller.Service.GetNotificationSubcriber(notify1.ID)
	if errGetNotificationSubcriber1 != nil {
		return errGetNotificationSubcriber1
	}
	errPushTest1 := controller.PushTest(notify1, userIDs1)
	if errPushTest1 != nil {
		return errPushTest1
	}

	return nil
}

// UpdatePostNotification func
func (controller NotificationController) UpdatePostNotification(userID int64) error {
	notify1, errUpdateStatusNotification := controller.Service.UpdateStatusNotification(userID)
	if errUpdateStatusNotification != nil {
		return errUpdateStatusNotification
	}
	userIDs1, errGetNotificationSubcriber1 := controller.Service.GetNotificationSubcriber(notify1.ID)
	if errGetNotificationSubcriber1 != nil {
		return errGetNotificationSubcriber1
	}
	errPushTest1 := controller.PushTest(notify1, userIDs1)
	if errPushTest1 != nil {
		return errPushTest1
	}

	return nil
}

// UpdateStatusNotification func
func (controller NotificationController) UpdateStatusNotification(userID int64) error {
	notify1, errUpdateStatusNotification := controller.Service.UpdateStatusNotification(userID)
	if errUpdateStatusNotification != nil {
		return errUpdateStatusNotification
	}
	userIDs1, errGetNotificationSubcriber1 := controller.Service.GetNotificationSubcriber(notify1.ID)
	if errGetNotificationSubcriber1 != nil {
		return errGetNotificationSubcriber1
	}
	errPushTest1 := controller.PushTest(notify1, userIDs1)
	if errPushTest1 != nil {
		return errPushTest1
	}

	return nil
}

// UpdatePhotoNotification func
func (controller NotificationController) UpdatePhotoNotification(userID int64) error {
	notify1, errUpdatePhotoNotification := controller.Service.UpdatePhotoNotification(userID)
	if errUpdatePhotoNotification != nil {
		return errUpdatePhotoNotification
	}
	userIDs1, errGetNotificationSubcriber1 := controller.Service.GetNotificationSubcriber(notify1.ID)
	if errGetNotificationSubcriber1 != nil {
		return errGetNotificationSubcriber1
	}
	errPushTest1 := controller.PushTest(notify1, userIDs1)
	if errPushTest1 != nil {
		return errPushTest1
	}

	return nil
}

// END UPDATE

// PushTest func
func (controller NotificationController) PushTest(notify models.Notification, userIDs []int64) error {

	push := fcm.NewFCM(configs.FCMToken)
	data := notify
	clientList, errGetDeviceByUserID := services.NewAccountService().GetDeviceByUserIDs(userIDs)
	if errGetDeviceByUserID != nil {
		return errGetDeviceByUserID
	}
	if len(clientList) == 0 {
		return nil
	}
	response, errSend := push.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  clientList,
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
	})
	if errSend != nil {
		fmt.Println("Status Code   :", response.StatusCode)
		fmt.Println("Success       :", response.Success)
		fmt.Println("Fail          :", response.Fail)
		fmt.Println("Canonical_ids :", response.CanonicalIDs)
		fmt.Println("Topic MsgId   :", response.MsgID)
		return errSend
	}
	return nil
}
