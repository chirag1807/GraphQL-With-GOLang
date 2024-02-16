package utils

import (
	"articlewithgraphql/config"
	"articlewithgraphql/constants"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

//it's not necessary to refresh token be jwt, it could be random 64 bit number (UUID or something)

func CreateAccessToken(time time.Time, id int64, isAdmin bool) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      strconv.FormatInt(id, 10),
		"isadmin": isAdmin,
		"exp":     time.Unix(),
	})
	
	tokenString, err := token.SignedString([]byte(config.JWtSecretConfig.SecretKey))
	if err != nil {
		return constants.EMPTY_STRING, err
	}

	return tokenString, nil
}

// niche na function ma always accesstoken avse request header mathi.
func VerifyToken(tokenString string) (int64, bool, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWtSecretConfig.SecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return 0, false, err
	}
	if !token.Valid {
		fmt.Println(constants.INVALID_TOKEN)
		return 0, false, errors.New(constants.INVALID_TOKEN)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println(constants.INVALID_CLAIMS)
		return 0, false, errors.New(constants.INVALID_CLAIMS)
	}

	idString, _ := claims["id"].(string)
	idInt, _ := strconv.ParseInt(idString, 10, 64)
	isadmin, _ := claims["isadmin"].(bool)

	return idInt, isadmin, nil
}
