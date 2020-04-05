package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"trackcoro/constants"
	"trackcoro/controller"
	"trackcoro/docs"
	"trackcoro/police"
	"trackcoro/token"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()
	addSwagger(router)
	addRoutes(router)
	return router
}

func addSwagger(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Track-Coro API"
	docs.SwaggerInfo.Description = "This is track corona api server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
func addRoutes(router *gin.Engine) {
	addHealthCheckRoute(router)
	addRoutesForQuarantine(router)
	addRoutesForPolice(router)
}

// HealthCheck godoc
// @Success 200 {string} string	"ok"
// @Summary Check status
// @Router /api/healthz [get]
func addHealthCheckRoute(router *gin.Engine) {
	healthCheck := func(c *gin.Context) {
		c.String(200, "ok")
	}
	router.GET("/api/healthz", healthCheck)
}

func addRoutesForQuarantine(router *gin.Engine) {
	router.POST("/api/v1/quarantine/verify", controller.QuarantineController.Verify)
	quarantineGroup := router.Group("/api/v1/quarantine")
	{
		quarantineGroup.Use(TokenAuthMiddleware())
		quarantineGroup.POST("/saveDetails", controller.QuarantineController.SaveProfileDetails)
		quarantineGroup.GET("/daysStatus", controller.QuarantineController.GetDaysStatus)
		quarantineGroup.GET("", controller.QuarantineController.GetProfileDetails)
	}
}

func addRoutesForPolice(router *gin.Engine) {
	policeController := police.NewController()
	quarantineGroup := router.Group("/api/v1/police")
	{
		quarantineGroup.GET("/save", policeController.SaveProfileDetails)
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader(constants.Authorization)
		if authToken == constants.Empty {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		userInfo, err := token.ReadToken(authToken)
		if err != nil || userInfo.MobileNumber == constants.Empty {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(constants.MobileNumber, userInfo.MobileNumber)
		ctx.Next()
	}
}
