package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"trackcoro/constants"
	models2 "trackcoro/models"
	"trackcoro/token"
)

type Controller interface {
	Verify(ctx *gin.Context)
	Add(ctx *gin.Context)
	AddSO(ctx *gin.Context)
	GetSOs(ctx *gin.Context)
}

type controller struct {
	service Service
}

func (c controller) Verify(ctx *gin.Context) {
	var verifyRequest models2.VerifyRequest
	err := ctx.ShouldBindBodyWith(&verifyRequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	isRegistered := c.service.Verify(verifyRequest.MobileNumber)
	response := models2.VerifyResponse{IsRegistered: isRegistered}
	if response.IsRegistered {
		addTokenInHeader(ctx, verifyRequest.MobileNumber)
	}
	ctx.JSON(http.StatusOK, response)
}

func (c controller) Add(ctx *gin.Context) {
	err := c.service.Add()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c controller) AddSO(ctx *gin.Context) {
	var addSORequest models2.SODetails
	err := ctx.ShouldBindBodyWith(&addSORequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.service.AddSO(getMobileNumber(ctx), addSORequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c controller) GetSOs(ctx *gin.Context) {
	SOs, err := c.service.GetSOs(getMobileNumber(ctx))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, SOs)
}

func addTokenInHeader(ctx *gin.Context, mobileNumber string) {
	tokenBody := token.UserInfo{MobileNumber: mobileNumber, Role: constants.AdminRole}
	generatedToken, generatedTime, err := token.GenerateToken(tokenBody)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Header("Token", generatedToken)
	ctx.Header("Generated-At", generatedTime.String())
}

func getMobileNumber(ctx *gin.Context) string {
	mobileNumber, _ := ctx.Get(constants.MobileNumber)
	return mobileNumber.(string)
}

func NewController(service Service) Controller {
	return controller{service}
}
