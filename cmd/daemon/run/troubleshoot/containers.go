package troubleshoot

import (
	handlers "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/pkg"

	"github.com/gin-gonic/gin"
)

// logsHandler godoc
// @Summary Logs of the containers
// @Description Will return the logs of the containers specified in the request, or all the containers if none are specified.
// @Tags logs
// @Accept json
// @Produce json
// @Param set body pkg.LogsRequest true "Logs parameters"
// @Success 200 {object} pkg.LogsResponse "Containers logs"
// @Failure 400 {object} pkg.ErrorResponse "Bad request with explanation"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error with explanation"
func logsHandler(c *gin.Context) {
	// Extract request body
	var req pkg.LogsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Since == "" {
		req.Since = "525600m" // 1 year
	}
	if req.Until == "" {
		req.Until = "0m"
	}
	if req.Tail == "" {
		req.Tail = "all"
	}

	// Call handler
	response, err := handlers.HandleLogs(pkg.LogsRequest{
		Containers: req.Containers,
		Since:      req.Since,
		Until:      req.Until,
		Tail:       req.Tail,
		Timestamps: req.Timestamps,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
}
