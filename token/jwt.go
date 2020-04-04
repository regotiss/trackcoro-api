package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	secretKey = "TRACKCORO_KEY"
)

var (
	jwtKey []byte = []byte(secretKey)
)

type JWTBody struct {
	MobileNumber string `json:"mobile_number"`
}

type Claims struct {
	JWTBody
	jwt.StandardClaims
}

func GetToken(body JWTBody) (string, time.Time, error) {
	expirationTime := time.Now().Add(10 * time.Second)
	claims := &Claims{
		JWTBody: body,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logrus.Error("GetToken: ", err)
		return "", time.Time{}, fmt.Errorf("not able to sign token")
	}
	return tokenString, expirationTime, err
}

func ReadToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logrus.Error("ReadToken: ", err)
			return nil, fmt.Errorf("signature invalid")
		}
		logrus.Error("ReadToken: ", err)
		return nil, fmt.Errorf("bad request")
	}
	if !tkn.Valid {
		logrus.Error("ReadToken: ", fmt.Errorf("invalid token"))
		return nil, fmt.Errorf("not valid token")
	}
	return claims, nil
}

func RefreshToken(token string) (string, time.Time, error) {
	claims, err := ReadToken(token)
	if err != nil {
		logrus.Error("RefreshToken: ", err)
		return "", time.Time{}, fmt.Errorf("could not read token")
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		logrus.Error("RefreshToken: ", err)
		return "", time.Time{}, fmt.Errorf("current token is not expired")
	}
	newToken, expiryTime, err := GetToken(claims.JWTBody)
	if err != nil {
		logrus.Error("RefreshToken: ", err)
		return "", time.Time{}, fmt.Errorf("could not get token")
	}
	return newToken, expiryTime, nil
}
