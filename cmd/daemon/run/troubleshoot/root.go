package troubleshoot

import (
	"github.com/gin-gonic/gin"
)

func NewTroubleshoot(router *gin.RouterGroup) {
	troubleshoot := router.Group("/troubleshoot")
	{
		troubleshoot.POST("/containers", logsHandler)
		troubleshoot.POST("/kernel", kernelHandler)
		troubleshoot.POST("/reboot", rebootHandler)
	}
}
