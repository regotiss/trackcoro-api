package quarantine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"trackcoro/constants"
	models2 "trackcoro/models"
	"trackcoro/quarantine/models"
	"trackcoro/utils"
)

type Controller interface {
	Verify(ctx *gin.Context)
	SaveProfileDetails(ctx *gin.Context)
	GetProfileDetails(ctx *gin.Context)
	GetRemainingDays(ctx *gin.Context)
	UploadPhoto(ctx *gin.Context)
	DownloadPhoto(ctx *gin.Context)
	UpdateCurrentLocation(ctx *gin.Context)
	UpdateDeviceTokenId(ctx *gin.Context)
	NotifySO(ctx *gin.Context)
}

type controller struct {
	service Service
}

func (c controller) Verify(ctx *gin.Context) {
	var verifyRequest models2.VerifyRequest
	bindError := ctx.ShouldBind(&verifyRequest)
	if bindError != nil {
		logrus.Error("Request bind body failed", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	isRegistered, isSignupCompleted := c.service.Verify(verifyRequest.MobileNumber)

	if isRegistered {
		utils.AddTokenInHeader(ctx, verifyRequest.MobileNumber, constants.QuarantineRole)
	}
	ctx.JSON(http.StatusOK, models.QVerifyResponse{
		VerifyResponse: models2.VerifyResponse{IsRegistered: isRegistered},
		IsSignUpCompleted: isSignupCompleted,
	})
}

func (c controller) SaveProfileDetails(ctx *gin.Context) {
	var saveDetailsRequest models2.QuarantineDetails
	bindError := ctx.ShouldBindBodyWith(&saveDetailsRequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", &constants.BadRequestError)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	saveDetailsRequest.MobileNumber = utils.GetMobileNumber(ctx)

	err := c.service.SaveDetails(saveDetailsRequest)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) UploadPhoto(ctx *gin.Context) {
	file, header, bindError := ctx.Request.FormFile("photo")
	if bindError != nil {
		logrus.Error("Couldn't find form photo", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}
	err := c.service.UploadPhoto(utils.GetMobileNumber(ctx), file, header)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) GetProfileDetails(ctx *gin.Context) {
	profileDetails, err := c.service.GetDetails(utils.GetMobileNumber(ctx))

	utils.HandleResponse(ctx, err, profileDetails, getStatusCode)
}

func (c controller) UpdateCurrentLocation(ctx *gin.Context) {
	var request models2.Coordinates
	bindError := ctx.ShouldBindBodyWith(&request, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	err := c.service.UpdateCurrentLocation(utils.GetMobileNumber(ctx), request.Latitude, request.Longitude)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) UpdateDeviceTokenId(ctx *gin.Context) {
	var request models2.DeviceTokeIdRequest
	bindError := ctx.ShouldBindBodyWith(&request, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	err := c.service.UpdateDeviceTokenId(utils.GetMobileNumber(ctx), request.DeviceTokeId)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) NotifySO(ctx *gin.Context) {
	var notificationRequest models2.NotificationRequest
	bindError := ctx.ShouldBindBodyWith(&notificationRequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	err := c.service.NotifySO(notificationRequest, utils.GetMobileNumber(ctx))

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) GetRemainingDays(ctx *gin.Context) {
	daysStatusResponse, err := c.service.GetDaysStatus(utils.GetMobileNumber(ctx))

	utils.HandleResponse(ctx, err, daysStatusResponse, getStatusCode)
}

func (c controller) DownloadPhoto(ctx *gin.Context) {
	content, err := c.service.DownloadPhoto(utils.GetMobileNumber(ctx))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	_, writeErr := ctx.Writer.Write(content)
	if writeErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &constants.InternalError)
		return
	}
	ctx.Header("Content-Type", http.DetectContentType(content))
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v.jpg;", utils.GetMobileNumber(ctx)))
	logrus.Info("length: ", len(content))
	ctx.Header("Content-Length", strconv.Itoa(len(content)))
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func getStatusCode(err *models2.Error) int {
	if err == nil {
		return http.StatusOK
	}
	if err.Code == constants.DOBIncorrectFormatCode ||
		err.Code == constants.QuarantineDateIncorrectFormatCode ||
		err.Code == constants.TravelDateIncorrectFormatCode {
		return http.StatusBadRequest
	}
	if err.Code == constants.QuarantineNotExistsCode {
		return http.StatusForbidden
	}

	return http.StatusInternalServerError
}

func NewController(service Service) Controller {
	return controller{service}
}
