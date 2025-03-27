package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// deleteFile safely removes a file at the specified path
// Returns error if:
// - file doesn't exist
// - no permission to delete
// - path is a directory
// - other OS-level errors
func DeleteFile(filePath string) error {
    // Check if file exists first
    fileInfo, err := os.Stat(filePath)
    if os.IsNotExist(err) {
        return fmt.Errorf("file does not exist: %s", filePath)
    }
    
    // Prevent accidental directory deletion
    if fileInfo.IsDir() {
        return fmt.Errorf("path is a directory, not a file: %s", filePath)
    }

    // Attempt to remove the file
    if err := os.Remove(filePath); err != nil {
        return fmt.Errorf("failed to delete file: %w", err) // %w wraps the original error
    }
    
    return nil
}

// handleFileUpload processes and validates the uploaded file
func HandleFileUpload(ctx *gin.Context, formFieldName string, uploadDir string) (string, error) {
    // Get file from form data
    file, err := ctx.FormFile(formFieldName)
    if err != nil {
        return "", fmt.Errorf("failed to get uploaded file: %w", err)
    }

    // Validate file size (5MB limit)
    const maxFileSize = 5 << 20 // 5MB in bytes
    if file.Size > maxFileSize {
        return "", fmt.Errorf("file size exceeds %dMB limit", maxFileSize>>20)
    }

    // Validate file extension
    fileExt := strings.ToLower(filepath.Ext(file.Filename))
    allowedExts := map[string]bool{
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
    }
    if !allowedExts[fileExt] {
        return "", errors.New("only .jpg, .jpeg or .png files are allowed")
    }

    // Generate unique filename
    nameWithoutExt := strings.TrimSuffix(
        filepath.Base(file.Filename),
        filepath.Ext(file.Filename),
    )
    timestamp := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS format
    filename := fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, fileExt)
    
    // Create upload directory if not exists
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create upload directory: %w", err)
    }

    // Save file
    savePath := filepath.Join(uploadDir, filename)
    if err := ctx.SaveUploadedFile(file, savePath); err != nil {
        return "", fmt.Errorf("failed to save uploaded file: %w", err)
    }

    return filename, nil
}