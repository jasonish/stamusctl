package config

import "github.com/gin-gonic/gin"

func NewConfig(router *gin.RouterGroup) {
	router.POST("/config", setHandler)
	router.GET("/config", getHandler)
}
