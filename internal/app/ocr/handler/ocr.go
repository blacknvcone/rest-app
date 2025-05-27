package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"rest-app/internal/app/ocr/port"
	"rest-app/pkg/helper"
	"strings"

	"github.com/gin-gonic/gin"
)

type handler struct {
	ocrService port.IOCRService
}

func New(ocrService port.IOCRService) port.IOCRHandler {
	return &handler{
		ocrService: ocrService,
	}
}

func (h *handler) ProcessReceipt(c *gin.Context) {
	const maxFileSize = 5 << 20 // 5 MB

	if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
		if err == http.ErrNotMultipart {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Request must be multipart/form-data",
			})
			return
		}
		if strings.Contains(err.Error(), "request body too large") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "File too large, max size is 5MB",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".pdf":  true,
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File extension %s not allowed", ext),
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open uploaded file",
		})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read file content",
		})
		return
	}

	// Process the file with your OCR service
	res, err := h.ocrService.ReceiptDataGenerator(c, fileBytes)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, &helper.Response{
		Success: true,
		Message: "Successfully processing image",
		Data:    res, // Include the result if available
	})
}
