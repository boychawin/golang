package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"login_jwt/errors"
	"login_jwt/models"
)

func CreateTokens(user models.Users) (string, string, error) {
	var SECRET = []byte(viper.GetString("jwt_ACCESS_TOKEN_SECRET"))

	if user.FirstName == "" {
		return "", "", errors.NewNotFoundError("ไม่พบข้อมูล")
	}

	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		Subject:   user.FirstName,
		Id:        strconv.Itoa(int(user.ID)),
		Audience:  user.Role,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessString, err := accessToken.SignedString([]byte(SECRET))
	if err != nil {
		return "", "", errors.NewNotFoundError(err.Error())
	}

	refreshClaims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(SECRET))
	if err != nil {
		return "", "", errors.NewNotFoundError(err.Error())
	}

	return accessString, refreshString, nil
}