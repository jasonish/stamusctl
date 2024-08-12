package compose

import "github.com/gin-gonic/gin"

func NewCompose(router *gin.RouterGroup) {
	compose := router.Group("/compose")
	{
		compose.POST("/init", initHandler)
		compose.POST("/update", updateHandler)
		compose.POST("/up", upHandler)
		compose.POST("/down", downHandler)
		compose.POST("/ps", psHandler)
		compose.POST("/restart/config", restartConfigHandler)
		compose.POST("/restart/containers", restartContainersHandler)
	}
}
