package compose

import (
	// External

	handlers "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/pkg"

	"github.com/gin-gonic/gin"
	// Custom
)

// UpHandler godoc
// @Summary Similar to docker compose up
// @Description Given a configuration, it will start the services defined in the configuration.
// @Tags up
// @Accept json
// @Produce json
// @Param arbitraries body pkg.UpRequest true "Up parameters"
// @Success 200 {object} SuccessResponse "Up successful"
// @Failure 400 {object} ErrorResponse "Bad request with explanation"
// @Router /compose/init [post]
func upHandler(c *gin.Context) {
	// Extract request body
	var req pkg.UpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Project == "" {
		req.Project = "tmp"
	}

	// Call handler
	err := handlers.HandleUp(req.Project)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})

}

// PsHandler godoc
// @Summary Similar to docker compose ps
// @Description Will return data about the containers running in the system.
// @Tags ps
// @Produce json
// @Success 200 {object} pkg.PsResponse "List of containers with their status"
// @Failure 400 {object} ErrorResponse "Bad request with explanation"
// @Router /compose/init [post]
func psHandler(c *gin.Context) {
	// Call handler
	containers, err := handlers.HandlePs()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, containers)
}
