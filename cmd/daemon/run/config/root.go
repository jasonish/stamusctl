package compose

import "github.com/gin-gonic/gin"

func NewConfig(router *gin.RouterGroup) {
	compose := router.Group("/compose")
	{
		compose.POST("/config", setHandler)
		compose.GET("/config", getHandler)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
