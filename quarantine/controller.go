package quarantine

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"trackcoro/quarantine/models"
)

type Controller interface {
	Verify(ctx *gin.Context)
	SaveProfileDetails(ctx *gin.Context)
	GetDaysStatus(ctx *gin.Context)
	GetProfileDetails(ctx *gin.Context)
}

type controller struct {
	service Service
}

func (c controller) Verify(ctx *gin.Context) {
	var verifyRequest models.VerifyRequest
	err := ctx.ShouldBindBodyWith(&verifyRequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	isRegistered := c.service.Verify(verifyRequest.MobileNumber)
	response := models.VerifyResponse{IsRegistered: isRegistered}
	ctx.JSON(http.StatusOK, response)
}

func (c controller) SaveProfileDetails(ctx *gin.Context) {
	var saveDetailsRequest models.ProfileDetails
	err := ctx.ShouldBindBodyWith(&saveDetailsRequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.service.SaveDetails(saveDetailsRequest)
	ctx.Status(getStatusCode(err))
}

func (c controller) GetDaysStatus(ctx *gin.Context) {
	var daysStatusRequest models.VerifyRequest
	err := ctx.ShouldBindBodyWith(&daysStatusRequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	daysStatusResponse, err := c.service.GetDaysStatus(daysStatusRequest.MobileNumber)
	status := getStatusCode(err)
	if status != http.StatusOK {
		ctx.AbortWithStatus(status)
		return
	}
	ctx.JSON(status, daysStatusResponse)
}

func (c controller) GetProfileDetails(ctx *gin.Context) {
	var getProfileDetailsRequest models.VerifyRequest
	err := ctx.ShouldBindBodyWith(&getProfileDetailsRequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	profileDetails, err := c.service.GetDetails(getProfileDetailsRequest.MobileNumber)
	status := getStatusCode(err)
	if status != http.StatusOK {
		ctx.AbortWithStatus(status)
		return
	}
	ctx.JSON(status, profileDetails)
}

func getStatusCode(err error) int {
	if err != nil && err.Error() == NotExists {
		return http.StatusUnauthorized
	}
	if err != nil && err.Error() == TimeParseError {
		return http.StatusBadRequest
	}
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func NewController(service Service) Controller {
	return controller{service}
}