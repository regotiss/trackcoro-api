package quarantine

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"trackcoro/constants"
	models2 "trackcoro/models"
	"trackcoro/quarantine/models"
	"trackcoro/utils"
)

type Controller interface {
	Verify(ctx *gin.Context)
	SaveProfileDetails(ctx *gin.Context)
	GetDaysStatus(ctx *gin.Context)
	GetProfileDetails(ctx *gin.Context)
	UploadPhoto(ctx *gin.Context)
	UpdateCurrentLocation(ctx *gin.Context)
	UpdateDeviceTokenId(ctx *gin.Context)
	NotifySO(ctx *gin.Context)
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
		utils.AddTokenInHeader(ctx, verifyRequest.MobileNumber, constants.QuarantineRole)
	}
	ctx.JSON(http.StatusOK, models2.VerifyResponse{IsRegistered: isRegistered})
}

func (c controller) SaveProfileDetails(ctx *gin.Context) {
	var saveDetailsRequest models2.QuarantineDetails
	bindError := ctx.ShouldBindBodyWith(&saveDetailsRequest, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", &constants.BadRequestError)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	saveDetailsRequest.MobileNumber = getMobileNumber(ctx)

	err := c.service.SaveDetails(saveDetailsRequest)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) GetDaysStatus(ctx *gin.Context) {
	daysStatusResponse, err := c.service.GetDaysStatus(getMobileNumber(ctx))

	utils.HandleResponse(ctx, err, daysStatusResponse, getStatusCode)
}

func (c controller) UploadPhoto(ctx *gin.Context) {
	file, header, bindError := ctx.Request.FormFile("photo")
	if bindError != nil {
		logrus.Error("Couldn't find form photo", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}
	contentType := header.Header.Get("Content-Type")

	err := c.service.UploadPhoto(getMobileNumber(ctx), file, header.Size, contentType)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) GetProfileDetails(ctx *gin.Context) {
	profileDetails, err := c.service.GetDetails(getMobileNumber(ctx))

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

	err := c.service.UpdateCurrentLocation(getMobileNumber(ctx), request.Latitude, request.Longitude)

	utils.HandleResponse(ctx, err, nil, getStatusCode)
}

func (c controller) UpdateDeviceTokenId(ctx *gin.Context) {
	var request models.DeviceTokeIdRequest
	bindError := ctx.ShouldBindBodyWith(&request, binding.JSON)
	if bindError != nil {
		logrus.Error("Request bind body failed", bindError)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &constants.BadRequestError)
		return
	}

	err := c.service.UpdateDeviceTokenId(getMobileNumber(ctx), request.DeviceTokeId)

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

	err := c.service.NotifySO(notificationRequest, getMobileNumber(ctx))

	utils.HandleResponse(ctx, err, nil, getStatusCode)
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

func getMobileNumber(ctx *gin.Context) string {
	mobileNumber, _ := ctx.Get(constants.MobileNumber)
	return mobileNumber.(string)
}

func NewController(service Service) Controller {
	return controller{service}
}
