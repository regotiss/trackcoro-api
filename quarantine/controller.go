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
	isRegistered, err := c.service.Verify(verifyRequest.MobileNumber)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response := models.VerifyResponse{IsRegistered:isRegistered}
	ctx.JSON(http.StatusOK, response)
}

func (c controller) SaveProfileDetails(ctx *gin.Context) {
	var saveDetailsRequest models.SaveDetailsRequest
	err := ctx.ShouldBindBodyWith(&saveDetailsRequest, binding.JSON)
	if err != nil {
		logrus.Error("Request bind body failed", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.service.SaveDetails(saveDetailsRequest)
	if err != nil && (err.Error() == NotExists || err.Error() == TimeParseError){
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func NewController(service Service) Controller {
	return controller{service}
}