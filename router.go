package main

import (
	"github.com/gin-gonic/gin"
	"trackcoro/police"
	"trackcoro/quarantine"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()
	addHealthCheckRoute(router)
	addRoutesForQuarantine(router)
	addRoutesForPolice(router)
	return router
}

func addHealthCheckRoute(router *gin.Engine) {
	router.GET("/api/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
}

func addRoutesForQuarantine(router *gin.Engine) {
	controller := quarantine.NewController()
	quarantineGroup := router.Group("/api/v1/quarantine")
	{
		quarantineGroup.GET("/save", controller.SaveProfileDetails)
	}
}

func addRoutesForPolice(router *gin.Engine) {
	controller := police.NewController()
	quarantineGroup := router.Group("/api/v1/police")
	{
		quarantineGroup.GET("/save", controller.SaveProfileDetails)
	}
}