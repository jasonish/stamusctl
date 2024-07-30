package troubleshoot

import (
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// kernelHandler godoc
// @Summary Logs of the kernel
// @Description Will return the logs of the kernel
// @Tags logs
// @Produce json
// @Success 200 {object} pkg.SuccessResponse "Kernel logs"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error with explanation"
func kernelHandler(c *gin.Context) {
	// Execute the `dmesg` command
	cmd := exec.Command("dmesg")
	output, err := cmd.Output()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error executing dmesg: %s", err)})
	}
	c.JSON(200, gin.H{"message": string(output)})
}
