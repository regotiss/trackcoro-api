package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"trackcoro/constants"
	models2 "trackcoro/models"
	"trackcoro/utils"
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
		utils.AddTokenInHeader(ctx, verifyRequest.MobileNumber, constants.AdminRole)
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
	err = c.service.AddSO(utils.GetMobileNumber(ctx), addSORequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c controller) GetSOs(ctx *gin.Context) {
	SOs, err := c.service.GetSOs(utils.GetMobileNumber(ctx))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, SOs)
}

func NewController(service Service) Controller {
	return controller{service}
}
