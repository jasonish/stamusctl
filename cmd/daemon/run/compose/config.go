package compose

import (
	handlers "stamus-ctl/internal/handlers/compose"

	"github.com/gin-gonic/gin"
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

type SetRequest struct {
	Reload  bool              `json:"reload"`
	Project string            `json:"project"`
	Values  map[string]string `json:"values"`
	Apply   bool              `json:"apply"`
}

func setHandler(c *gin.Context) {
	// Extract request body
	var req SetRequest
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

	// Call handler
	if err := handlers.SetHandler(req.Project, valuesAsStrings, req.Reload, req.Apply); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

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

type GetRequest struct {
	Project string   `json:"project"`
	Values  []string `json:"values"`
}

type GetResponse map[string]interface{}

func getHandler(c *gin.Context) {
	// Extract request body
	var req GetRequest
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
	groupedValues, err := handlers.GetGroupedConfig(req.Project, req.Values, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, groupedValues)
}
