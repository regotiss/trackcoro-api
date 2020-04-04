package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"trackcoro/controller"
	"trackcoro/docs"
	"trackcoro/police"
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
	quarantineGroup := router.Group("/api/v1/quarantine")
	{
		quarantineGroup.POST("/verify", controller.QuarantineController.Verify)
		quarantineGroup.POST("/saveDetails", controller.QuarantineController.SaveProfileDetails)
	}
}

func addRoutesForPolice(router *gin.Engine) {
	controller := police.NewController()
	quarantineGroup := router.Group("/api/v1/police")
	{
		quarantineGroup.GET("/save", controller.SaveProfileDetails)
	}
}
