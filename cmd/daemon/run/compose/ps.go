package compose

import (
	// Internal
	handlers "stamus-ctl/internal/handlers/compose"
	// External
	"github.com/gin-gonic/gin"
)

// PsHandler godoc
// @Summary Similar to docker compose ps
// @Description Will return data about the containers running in the system.
// @Tags ps
// @Produce json
// @Success 200 {object} pkg.PsResponse "List of containers with their status"
// @Failure 400 {object} pkg.ErrorResponse "Bad request with explanation"
// @Router /compose/ps [post]
func psHandler(c *gin.Context) {
	// Call handler
	containers, err := handlers.HandlePs()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, containers)
}
