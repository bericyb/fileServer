package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func main() {
	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		c.SaveUploadedFile(file, filepath.Join("uploads", file.Filename))
		c.String(http.StatusOK, "File uploaded successfully: %!(NOVERB)s", file.Filename)
	})

	router.GET("/", func(c *gin.Context) {
		filename := c.Query("filename")
		targetPath := filepath.Join("uploads", filename)
		c.File(targetPath)
	})

	router.Run("0.0.0.0:8080")
}
