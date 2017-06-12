package controllers

import (
	"errors"
	"fmt"

	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/services"
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

// GetRoleFromUserID func return role
func GetRoleFromUserID(userID int64) (int, error) {
	fmt.Printf("vao get role controller: %d\n", userID)
	return services.NewAccountService().GetRoleFromUserID(userID)
}
