package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"trackcoro/constants"
	"trackcoro/controller"
	"trackcoro/docs"
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
	addRoutesForAdmin(router)
	addRoutesForQuarantine(router)
	addRoutesForSO(router)
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

func addRoutesForAdmin(router *gin.Engine) {
	router.POST("/api/v1/admin/verify", controller.AdminController.Verify)
	router.POST("/api/v1/admin/add", controller.AdminController.Add)
	adminGroup := router.Group("/api/v1/admin")
	{
		adminGroup.Use(TokenAuthMiddleware(constants.AdminRole))
		adminGroup.POST("/addSO", controller.AdminController.AddSO)
		adminGroup.GET("/SOs", controller.AdminController.GetSOs)
		adminGroup.POST("/quarantines", controller.AdminController.GetQuarantines)
		adminGroup.POST("/deleteSO", controller.AdminController.DeleteSO)
		adminGroup.POST("/replaceSO", controller.AdminController.ReplaceSO)
	}
}

func addRoutesForSO(router *gin.Engine) {
	router.POST("/api/v1/so/verify", controller.SOController.Verify)
	quarantineGroup := router.Group("/api/v1/so")
	{
		quarantineGroup.Use(TokenAuthMiddleware(constants.SORole))
		quarantineGroup.POST("/addQuarantine", controller.SOController.AddQuarantine)
		quarantineGroup.GET("/quarantines", controller.SOController.GetQuarantines)
		quarantineGroup.POST("/deleteQuarantine", controller.SOController.DeleteQuarantine)
	}
}

func addRoutesForQuarantine(router *gin.Engine) {
	router.POST("/api/v1/quarantine/verify", controller.QuarantineController.Verify)
	quarantineGroup := router.Group("/api/v1/quarantine")
	{
		quarantineGroup.Use(TokenAuthMiddleware(constants.QuarantineRole))
		quarantineGroup.GET("", controller.QuarantineController.GetProfileDetails)
		quarantineGroup.POST("/saveDetails", controller.QuarantineController.SaveProfileDetails)
		quarantineGroup.GET("/daysStatus", controller.QuarantineController.GetDaysStatus)
		quarantineGroup.POST("/uploadPhoto", controller.QuarantineController.UploadPhoto)
		quarantineGroup.POST("/saveCurrentLocation", controller.QuarantineController.UpdateCurrentLocation)
		quarantineGroup.POST("/saveDeviceTokenId", controller.QuarantineController.UpdateDeviceTokenId)
	}
}

func TokenAuthMiddleware(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader(constants.Authorization)
		if authToken == constants.Empty {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		userInfo, err := token.ReadToken(authToken)
		if err != nil || userInfo.MobileNumber == constants.Empty || userInfo.Role != role {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(constants.MobileNumber, userInfo.MobileNumber)
		ctx.Next()
	}
}
