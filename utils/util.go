package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trackcoro/constants"
	"trackcoro/token"
)

func AddTokenInHeader(ctx *gin.Context, mobileNumber string, role string) {
	tokenBody := token.UserInfo{MobileNumber: mobileNumber, Role: role}
	generatedToken, generatedTime, err := token.GenerateToken(tokenBody)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Header("Token", generatedToken)
	ctx.Header("Generated-At", generatedTime.String())
}

func GetMobileNumber(ctx *gin.Context) string {
	mobileNumber, _ := ctx.Get(constants.MobileNumber)
	return mobileNumber.(string)
}
