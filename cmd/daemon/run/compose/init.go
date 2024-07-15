package compose

import (
	// External
	"github.com/gin-gonic/gin"
	// Custom
	"stamus-ctl/internal/app"
	handlers "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/pkg"
)

// Init godoc
// @Summary Initialize configuration
// @Description Initializes configuration with provided parameters.
// @Tags init
// @Accept json
// @Produce json
// @Param arbitraries body InitRequest true "Initialization parameters"
// @Success 200 {object} SuccessResponse "Initialization successful"
// @Failure 400 {object} ErrorResponse "Bad request with explanation"
// @Router /compose/init [post]

func initHandler(c *gin.Context) {
	// Extract request body
	var req pkg.InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Folder == "" {
		req.Folder = "tmp"
	}
	if req.Values == nil {
		req.Values = make(map[string]string)
	}
	if req.Project == "" {
		req.Project = "selks"
	}
	if req.Version == "" {
		req.Version = "latest"
	}

	// Call handler
	parameters := handlers.InitHandlerInputs{
		IsDefault:        req.IsDefault,
		FolderPath:       app.LatestSelksPath,
		BackupFolderPath: app.DefaultSelksPath,
		OutputPath:       req.Folder,
		Arbitrary:        req.Values,
		Project:          req.Project,
	}
	if err := handlers.InitHandler(false, parameters); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})

}
