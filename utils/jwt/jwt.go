package jwt

import (
	"auth/config"
	"auth/models"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func GenerateToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["code"] = user.Code
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["scope"] = "default"

	t, err := token.SignedString([]byte(config.C.Token.Secret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GenerateAuthToken(user *models.AuthUser) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["phone_number"] = user.Phone
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 180).Unix()
	claims["scope"] = "auth_token"

	t, err := token.SignedString([]byte(config.C.Token.AuthSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GenerateSetPasswordToken(user *models.AuthUser) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Phone
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["scope"] = "password"

	t, err := token.SignedString([]byte(config.C.Token.PasswordSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GenerateVerifyEmailToken(user *models.AuthUser) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["scope"] = "verify_email"

	t, err := token.SignedString([]byte(config.C.Token.VerifyEmailSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GenerateResetPassByEmailToken(user *models.AuthUser) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["scope"] = "verify_email"

	t, err := token.SignedString([]byte(config.C.Token.ResetPassByEmailSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GenerateClientToken(client *models.Client) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = client.ClientID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["scope"] = "client"

	t, err := token.SignedString([]byte(client.ClientSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GenerateAppToken(client *models.Client, payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = client.ClientID
	claims["payload"] = payload
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["scope"] = "app"
	t, err := token.SignedString([]byte(client.ClientSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GenerateMega1CMSToken(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["scope"] = "mega1_cms"
	t, err := token.SignedString([]byte(config.C.Token.Mega1CMSSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func DecodeToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	fmt.Println("\ndecode user: ", user)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println("\nclaims user: ", claims)
	id := claims["id"].(string)
	return id
}

func GetClaimsUnverified(tokenString string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, err
	}
}

func VerifiedToken(tokenString string, key string) error {
	_, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func GetJWTFromHeader(c echo.Context, authScheme string) (string, error) {
	auth := c.Request().Header.Get("Authorization")
	l := len(authScheme)
	if len(auth) > l+1 && auth[:l] == authScheme {
		return auth[l+1:], nil
	}
	return "", errors.New("missing or malformed jwt")
}
