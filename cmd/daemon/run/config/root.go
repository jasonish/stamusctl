package config

import "github.com/gin-gonic/gin"

func NewConfig(router *gin.RouterGroup) {
	router.POST("/config", setHandler)
	router.POST("/config/current", setCurrentHandler)
	router.GET("/config", getHandler)
}
