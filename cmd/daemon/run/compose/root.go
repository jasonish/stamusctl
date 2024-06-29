package compose

import "github.com/gin-gonic/gin"

func NewCompose(router *gin.RouterGroup) {
	router.GET("/pingcompose", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pongcompose",
		})
	})
}
