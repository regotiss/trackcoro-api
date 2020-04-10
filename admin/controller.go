package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"trackcoro/admin/models"
	"trackcoro/constants"
	models2 "trackcoro/models"
	"trackcoro/utils"
)

type Controller interface {
	Verify(ctx *gin.Context)
	Add(ctx *gin.Context)
	AddSO(ctx *gin.Context)
	GetSOs(ctx *gin.Context)
	GetQuarantines(ctx *gin.Context)
	DeleteSO(ctx *gin.Context)
	ReplaceSO(ctx *gin.Context)
	DeleteAllSOs(ctx *gin.Context)
}

type controller struct {
	service   Service
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
	response := models2.VerifyResponse{IsRegistered: isRegistered}
	if response.IsRegistered {
		utils.AddTokenInHeader(ctx, verifyRequest.MobileNumber, constants.AdminRole)
	}
	ctx.JSON(http.StatusOK, response)
}

func (c controller) Add(ctx *gin.Context) {
	err := c.service.Add()
	utils.HandleResponse(ctx, err, nil, handleError)
}

func (c controller) AddSO(ctx *gin.Context) {
	var addSORequest models2.SODetails
	bindError := ctx.ShouldBindBodyWith(&addSORequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed ", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}
	err := c.service.AddSO(utils.GetMobileNumber(ctx), addSORequest)

	utils.HandleResponse(ctx, err, nil, handleError)
}

func (c controller) GetSOs(ctx *gin.Context) {
	SOs, err := c.service.GetSOs(utils.GetMobileNumber(ctx))

	utils.HandleResponse(ctx, err, SOs, handleError)
}

func (c controller) GetQuarantines(ctx *gin.Context) {
	var quarantinesRequest models.GetQuarantinesRequest
	bindError := ctx.ShouldBindBodyWith(&quarantinesRequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed ", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}
	quarantines, err := c.service.GetQuarantines(utils.GetMobileNumber(ctx), quarantinesRequest.MobileNumber)

	utils.HandleResponse(ctx, err, quarantines, handleError)
}

func (c controller) DeleteSO(ctx *gin.Context) {
	var deleteSORequest models.GetQuarantinesRequest
	bindError := ctx.ShouldBindBodyWith(&deleteSORequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed ", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	err := c.service.DeleteSO(utils.GetMobileNumber(ctx), deleteSORequest.MobileNumber)

	utils.HandleResponse(ctx, err, nil, handleError)
}

func (c controller) ReplaceSO(ctx *gin.Context) {
	var replaceSORequest models.ReplaceSORequest
	bindError := ctx.ShouldBindBodyWith(&replaceSORequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed ", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	err := c.service.ReplaceSO(utils.GetMobileNumber(ctx), replaceSORequest.OldSOMobileNumber, replaceSORequest.NewSOMobileNumber)

	utils.HandleResponse(ctx, err, nil, handleError)
}

func (c controller) DeleteAllSOs(ctx *gin.Context) {
	err := c.service.DeleteAllSOs(utils.GetMobileNumber(ctx))

	utils.HandleResponse(ctx, err, nil, handleError)
}

func handleError(err *models2.Error) int {
	if err == nil {
		return http.StatusOK
	}
	if err.Code == constants.AdminNotExistsCode {
		return http.StatusForbidden
	}
	if err.Code == constants.SONotExistsCode ||
		err.Code == constants.SONotRegisteredByAdminCode ||
		err.Code == constants.SOAlreadyExistsCode {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func NewController(service Service) Controller {
	return controller{service}
}
