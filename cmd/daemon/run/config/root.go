package config

import "github.com/gin-gonic/gin"

func NewConfig(router *gin.RouterGroup) {
	router.POST("/config", setHandler)
	router.POST("/config/current", setCurrentHandler)
	router.POST("/config/list", getConfigListHandler)
	router.GET("/config", getHandler)
}
