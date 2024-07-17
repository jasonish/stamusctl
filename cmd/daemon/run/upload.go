package run

import (
	// Custom

	"os"
	"path/filepath"
	"stamus-ctl/internal/logging"

	// External
	"github.com/gin-gonic/gin"
)

// UploadHandler godoc
// @Summary Upload file example
// @Schemes
// @Description Handles file uploads
// @Tags example
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Upload file"
// @Success 200 {object} map[string]string{ "message": "Uploaded file" }
// @Failure 400 {object} map[string]string{ "error": "Error message" }
// @Failure 500 {object} map[string]string{ "error": "Error message" }
// @Router /upload [post]

func uploadHandler(c *gin.Context) {
	logging.LoggerWithRequest(c.Request).Info("Upload")

	// Validate request
	if c.Request.ContentLength == 0 {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}
	if c.Query("path") == "" {
		c.JSON(400, gin.H{"error": "No path provided"})
		return
	}
	project := c.Query("project")
	if project == "" {
		project = "tmp"
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File upload error: " + err.Error()})
		return
	}

	// Create directory if it doesn't exist
	err = os.MkdirAll(filepath.Join(project, c.Query("path")), 0755)
	if err != nil {
		c.JSON(500, gin.H{"error": "Directory creation error: " + err.Error()})
		return
	}

	// Save file to path
	err = c.SaveUploadedFile(file, filepath.Join(project, c.Query("path"), file.Filename))
	if err != nil {
		c.JSON(500, gin.H{"error": "File save error: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Uploaded file",
	})
}
