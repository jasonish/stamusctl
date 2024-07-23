package config

import (
	// External
	"github.com/gin-gonic/gin"
	//Custom
	handlers "stamus-ctl/internal/handlers/config"
	"stamus-ctl/pkg"
)

// setHandler godoc
// @Summary Set configuration
// @Description Sets configuration with provided parameters.
// @Tags set
// @Accept json
// @Produce json
// @Param set body SetRequest true "Set parameters"
// @Success 200 {object} SuccessResponse "Configuration set successfully"
// @Failure 400 {object} ErrorResponse "Bad request with explanation"
// @Failure 500 {object} ErrorResponse "Internal server error with explanation"
// @Router /compose/config [post]

func setHandler(c *gin.Context) {
	// Extract request body
	var req pkg.SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Project == "" {
		req.Project = "tmp"
	}
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
		Config:   req.Project,
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

// getHandler godoc
// @Summary Get configuration
// @Description Retrieves configuration for a given project.
// @Tags get
// @Produce json
// @Param project query string true "Project name"
// @Success 200 {object} GetResponse "Configuration retrieved successfully"
// @Failure 404 {object} ErrorResponse "Bad request with explanation"
// @Failure 500 {object} ErrorResponse "Internal server error with explanation"
// @Router /compose/config [get]

type GetResponse map[string]interface{}

func getHandler(c *gin.Context) {
	// Extract request body
	var req pkg.GetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Project == "" {
		req.Project = "tmp"
	}
	if req.Values == nil {
		req.Values = []string{}
	}

	// Call handler
	if req.Content {
		filesMap, err := handlers.GetGroupedContent(req.Project, req.Values)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, filesMap)
		return
	}

	groupedValues, err := handlers.GetGroupedConfig(req.Project, req.Values, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, groupedValues)
}
