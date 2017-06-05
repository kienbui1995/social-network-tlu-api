package controllers

import (
	"errors"

	"github.com/kienbui1995/social-network-tlu-api/helpers"
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
