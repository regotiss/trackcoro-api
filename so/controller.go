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
	bindError := ctx.ShouldBindBodyWith(&verifyRequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	isRegistered := c.service.Verify(verifyRequest.MobileNumber)

	if isRegistered {
		utils.AddTokenInHeader(ctx, verifyRequest.MobileNumber, constants.SORole)
	}
	ctx.JSON(http.StatusOK, models2.VerifyResponse{IsRegistered: isRegistered})
}

func (c controller) AddQuarantine(ctx *gin.Context) {
	var addQuarantineRequest models2.VerifyRequest
	bindError := ctx.ShouldBindBodyWith(&addQuarantineRequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	err := c.service.AddQuarantine(utils.GetMobileNumber(ctx), addQuarantineRequest.MobileNumber)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) GetQuarantines(ctx *gin.Context) {
	quarantines, err := c.service.GetQuarantines(utils.GetMobileNumber(ctx))

	utils.HandleResponse(ctx, err, quarantines, getStatusCode)
}

func (c controller) DeleteQuarantine(ctx *gin.Context) {
	var removeQuarantineRequest models.RemoveQuarantineRequest
	bindError := ctx.ShouldBindBodyWith(&removeQuarantineRequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	err := c.service.DeleteQuarantine(utils.GetMobileNumber(ctx), removeQuarantineRequest.MobileNumber)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func getStatusCode(err *models2.Error) int {
	if err == nil {
		return http.StatusOK
	}
	if err.Code == constants.SONotExistsCode ||
		err.Code == constants.QuarantineNotExistsCode ||
		err.Code == constants.QuarantineAlreadyExistsCode {
		return http.StatusBadRequest
	}
	if err.Code == constants.QuarantineNotRegisteredBySOError.Code {
		return http.StatusForbidden
	}

	return http.StatusInternalServerError
}

func NewController(service Service) Controller {
	return controller{service}
}
