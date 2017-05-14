package helpers

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// GenerateToken func to set token for a new login
func GenerateToken(account models.Account, secret []byte) (string, error) {
	// init token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": account.ID,
		"device": account.Device,
		"role":   account.Role,
		"exp":    time.Now().Add(time.Hour * 720).Unix(),
	})

	tokenstring, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

// ValidateToken func to authen token
func ValidateToken(tokenstring string, secret []byte) (bool, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err := claims.Valid()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil

}

// ExtractClaims func to get map claims
func ExtractClaims(tokenstring string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Don't extract claim")
}
