package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

var Validate = validator.New()

func SaveUploadedFile(c *gin.Context, formFieldName string, uploadDir string) (string, error) {
	file, err := c.FormFile(formFieldName)
	if err != nil {
		Logger.WithError(err).Error("File upload failed")
		return "", err
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		Logger.Warnf("Invalid file type: %s", file.Header.Get("Content-Type"))
		return "", errors.New("invalid file type")
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		Logger.Warnf("Invalid file extension: %s", ext)
		return "", errors.New("invalid file extension")
	}

	if file.Size > 5*1024*1024 {
		Logger.Warnf("File size too large: %d bytes", file.Size)
		return "", errors.New("file size exceeds limit (5MB)")
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		Logger.WithError(err).Error("Failed to create upload directory")
		return "", err
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Logger.WithError(err).Error("Failed to save uploaded file")
		return "", err
	}

	Logger.Infof("File uploaded successfully: %s (Size: %d bytes)", filePath, file.Size)
	return filePath, nil
}
