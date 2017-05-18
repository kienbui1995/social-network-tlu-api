package controllers

import (
	"errors"
	"fmt"

	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
	"github.com/kienbui1995/social-network-tlu-api/services"
	"github.com/maddevsio/fcm"
)

// GetUserIDFromToken func return userid to check permission
func GetUserIDFromToken(token string) (int64, error) {

	if len(token) == 0 {
		return -1, errors.New("NULL userid in token")
	}

	claims, errclaim := helpers.ExtractClaims(token, secret)
	if errclaim != nil {
		return -1, errclaim
	}
	return int64(claims["userid"].(float64)), nil
}

// PushTest func
func PushTest(notify models.Notification) error {

	push := fcm.NewFCM(configs.FCMToken)
	data := map[string]interface{}{
		"id":      notify.ObjectID,
		"type":    notify.ObjectType,
		"message": notify.Title + ": " + notify.Message,
	}
	clientList, errGetDeviceByUserID := services.NewAccountService().GetDeviceByUserID(notify.UserID)
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
