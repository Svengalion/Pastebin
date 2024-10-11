// internal/utils/token.go
package utils

import (
	"time"

	"github.com/Svengalion/Pastebin/internal/models"
	"github.com/golang-jwt/jwt"
)

func GenerateJWT(user *models.User, jwtSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"login": user.Login,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
