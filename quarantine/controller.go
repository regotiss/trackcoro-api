package quarantine

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"trackcoro/constants"
	"trackcoro/quarantine/models"
	"trackcoro/token"
)

type Controller interface {
	Verify(ctx *gin.Context)
	SaveProfileDetails(ctx *gin.Context)
	GetDaysStatus(ctx *gin.Context)
	GetProfileDetails(ctx *gin.Context)
	UploadPhoto(ctx *gin.Context)
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
	if response.IsRegistered {
		addTokenInHeader(ctx, verifyRequest.MobileNumber)
	}
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
	saveDetailsRequest.MobileNumber = getMobileNumber(ctx)
	err = c.service.SaveDetails(saveDetailsRequest)
	ctx.Status(getStatusCode(err))
}

func (c controller) GetDaysStatus(ctx *gin.Context) {
	daysStatusResponse, err := c.service.GetDaysStatus(getMobileNumber(ctx))
	status := getStatusCode(err)
	if status != http.StatusOK {
		ctx.AbortWithStatus(status)
		return
	}
	ctx.JSON(status, daysStatusResponse)
}

func (c controller) UploadPhoto(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("photo")
	if err != nil {
		logrus.Error("Couldn't find form photo", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.service.UploadPhoto(getMobileNumber(ctx), file, header.Size)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c controller) GetProfileDetails(ctx *gin.Context) {
	profileDetails, err := c.service.GetDetails(getMobileNumber(ctx))
	status := getStatusCode(err)
	if status != http.StatusOK {
		ctx.AbortWithStatus(status)
		return
	}
	ctx.JSON(status, profileDetails)
}

func addTokenInHeader(ctx *gin.Context, mobileNumber string) {
	tokenBody := token.UserInfo{MobileNumber: mobileNumber, Role: constants.QuarantineRole}
	generatedToken, generatedTime, err := token.GenerateToken(tokenBody)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Header("Token", generatedToken)
	ctx.Header("Generated-At", generatedTime.String())
}

func getStatusCode(err error) int {
	if err != nil && err.Error() == constants.NotExists {
		return http.StatusUnauthorized
	}
	if err != nil && err.Error() == constants.TimeParseError {
		return http.StatusBadRequest
	}
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func getMobileNumber(ctx *gin.Context) string {
	mobileNumber, _ := ctx.Get(constants.MobileNumber)
	return mobileNumber.(string)
}

func NewController(service Service) Controller {
	return controller{service}
}
