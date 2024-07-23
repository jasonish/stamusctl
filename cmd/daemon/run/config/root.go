package config

import "github.com/gin-gonic/gin"

func NewConfig(router *gin.RouterGroup) {
	router.POST("/config", setHandler)
	router.GET("/config", getHandler)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
