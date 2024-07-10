package compose

import (
	handlers "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/internal/models"

	"github.com/gin-gonic/gin"
)

// Init godoc
// @Summary Update configuration
// @Description Update configuration with provided parameters.
// @Tags init
// @Accept json
// @Produce json
// @Param arbitraries body UpdateRequest true "Update parameters"
// @Success 200 {object} SuccessResponse "Initialization successful"
// @Failure 400 {object} ErrorResponse "Bad request with explanation"
// @Router /compose/update [post]

type UpdateRequest struct {
	models.RegistryInfo
	Version string            `json:"version"`
	Project string            `json:"project"`
	Values  map[string]string `json:"values"`
}

func updateHandler(c *gin.Context) {
	// Extract request body
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate parameters
	if err := req.ValidateRegistry(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.Project == "" {
		req.Project = "tmp"
	}
	if req.Version == "" {
		req.Version = "latest"
	}
	if req.Values == nil {
		req.Values = make(map[string]string)
	}
	valuesVal := []string{}
	for k, v := range req.Values {
		valuesVal = append(valuesVal, k+"="+v)
	}

	// Call handler
	params := handlers.UpdateHandlerParams{
		RegistryInfo: models.RegistryInfo{
			Registry: req.Registry,
			Username: req.Username,
			Password: req.Password,
		},
		Config:  req.Project,
		Args:    valuesVal,
		Version: req.Version,
	}
	err := handlers.UpdateHandler(params)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Configuration updated successfully"})

}
