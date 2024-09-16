package config

import (
	// External
	"github.com/gin-gonic/gin"
	//Custom
	"stamus-ctl/internal/app"
	handlers "stamus-ctl/internal/handlers/config"
	"stamus-ctl/pkg"
)

// setHandler godoc
// @Summary Set configuration
// @Description Sets configuration with provided parameters.
// @Tags set
// @Accept json
// @Produce json
// @Param set body pkg.SetRequest true "Set parameters"
// @Success 200 {object} pkg.SuccessResponse "Configuration set successfully"
// @Failure 400 {object} pkg.ErrorResponse "Bad request with explanation"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error with explanation"
// @Router /config [post]
func setHandler(c *gin.Context) {
	// Extract request body
	var req pkg.SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Values == nil {
		req.Values = make(map[string]string)
	}
	valuesAsStrings := []string{}
	for k, v := range req.Values {
		valuesAsStrings = append(valuesAsStrings, k+"="+v)
	}
	fromFile := ""
	if req.FromFile != nil {
		for k, v := range req.FromFile {
			fromFile += k + "=" + v + " "
		}
	}

	// Call handler
	params := handlers.SetHandlerInputs{
		Reload:   req.Reload,
		Apply:    req.Apply,
		Args:     valuesAsStrings,
		Values:   req.ValuesPath,
		FromFile: fromFile,
	}
	if err := handlers.SetHandler(params); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})

}

// setCurrentHandler godoc
// @Summary Set current configuration
// @Description Sets configuration with provided parameters.
// @Tags set
// @Accept json
// @Produce json
// @Param set body pkg.Config true "Configuration name to use"
// @Success 200 {object} pkg.SuccessResponse "Configuration set successfully"
// @Failure 400 {object} pkg.ErrorResponse "Bad request with explanation"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error with explanation"
// @Router /config [post]
func setCurrentHandler(c *gin.Context) {
	// Extract request body
	var req pkg.Config
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Value == "" {
		req.Value = app.DefaultConfigName
	}

	// Call handler
	if err := handlers.SetCurrent(req.Value); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
}

type GetResponse map[string]interface{}

// getHandler godoc
// @Summary Get configuration
// @Description Retrieves configuration for a given project.
// @Tags get
// @Produce json
// @Param set body pkg.GetRequest true "Get parameters"
// @Success 200 {object} GetResponse "Configuration retrieved successfully"
// @Failure 404 {object} pkg.ErrorResponse "Bad request with explanation"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error with explanation"
// @Router /config [get]
func getHandler(c *gin.Context) {
	// Extract request body
	var req pkg.GetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Values == nil {
		req.Values = []string{}
	}

	// Call handler
	if req.Content {
		filesMap, err := handlers.GetGroupedContent(req.Values)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, filesMap)
		return
	}

	groupedValues, err := handlers.GetGroupedConfig(req.Values, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, groupedValues)
}
