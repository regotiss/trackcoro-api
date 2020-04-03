package police

import "github.com/gin-gonic/gin"

type Controller interface {
	SaveProfileDetails(ctx *gin.Context)
}

type controller struct {

}

func (c controller) SaveProfileDetails(ctx *gin.Context) {
	panic("implement me")
}

func NewController() Controller {
	return controller{}
}