package so

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"trackcoro/constants"
	models2 "trackcoro/models"
	"trackcoro/so/models"
	"trackcoro/utils"
)

type Controller interface {
	Verify(ctx *gin.Context)
	AddQuarantine(ctx *gin.Context)
	GetQuarantines(ctx *gin.Context)
	DeleteQuarantine(ctx *gin.Context)
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
		utils.AddTokenInHeader(ctx, verifyRequest.MobileNumber, constants.SORole)
	}
	ctx.JSON(http.StatusOK, response)
}

func (c controller) AddQuarantine(ctx *gin.Context) {
	var addQuarantineRequest models2.VerifyRequest
	err := ctx.ShouldBindBodyWith(&addQuarantineRequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.service.AddQuarantine(utils.GetMobileNumber(ctx), addQuarantineRequest.MobileNumber)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c controller) GetQuarantines(ctx *gin.Context) {
	quarantines, err := c.service.GetQuarantines(utils.GetMobileNumber(ctx))

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, quarantines)
}

func (c controller) DeleteQuarantine(ctx *gin.Context) {
	var removeQuarantineRequest models.RemoveQuarantineRequest
	err := ctx.ShouldBindBodyWith(&removeQuarantineRequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = c.service.DeleteQuarantine(utils.GetMobileNumber(ctx), removeQuarantineRequest.MobileNumber)

	ctx.Status(getStatusCode(err))
}

func getStatusCode(err error) int {
	if err != nil && (err.Error() == constants.SONotExistsError || err.Error() == constants.QuarantineNotExistsError) {
		return http.StatusBadRequest
	}
	if err != nil && err.Error() == constants.QuarantineNotRegisteredBySOError {
		return http.StatusUnauthorized
	}
	if err != nil {
		return http.StatusInternalServerError
	}
	return  http.StatusOK
}

func NewController(service Service) Controller {
	return controller{service}
}

