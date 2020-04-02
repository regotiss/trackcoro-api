package main

import "github.com/gin-gonic/gin"

func InitializeRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	return router
}