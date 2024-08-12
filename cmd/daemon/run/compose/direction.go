package compose

import (
	// Internal
	handlers "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/pkg"

	// External
	"github.com/gin-gonic/gin"
)

// UpHandler godoc
// @Summary Similar to docker compose up
// @Description Given a configuration, it will start the services defined in the configuration.
// @Tags compose
// @Accept json
// @Produce json
// @Param arbitraries body pkg.Config true "Configuration to start"
// @Success 200 {object} pkg.SuccessResponse "Up successful"
// @Failure 400 {object} pkg.ErrorResponse "Bad request with explanation"
// @Router /compose/up [post]
func upHandler(c *gin.Context) {
	// Extract request body
	var req pkg.Config
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Value == "" {
		req.Value = "tmp"
	}

	// Call handler
	err := handlers.HandleUp(req.Value)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
}

// UpHandler godoc
// @Summary Similar to docker compose down
// @Description Given a configuration, it will stop the services defined in the configuration.
// @Tags compose
// @Accept json
// @Produce json
// @Param arbitraries body pkg.Config true "Configuration to stop"
// @Success 200 {object} pkg.SuccessResponse "Down successful"
// @Failure 400 {object} pkg.ErrorResponse "Bad request with explanation"
// @Router /compose/down [post]
func downHandler(c *gin.Context) {
	// Extract request body
	var req pkg.Config
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Value == "" {
		req.Value = "tmp"
	}

	// Call handler
	err := handlers.HandleDown(req.Value, false, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
}
