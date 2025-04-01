package common

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/StephenJonesIT/CCMS/src/post-service/models"
	"github.com/gin-gonic/gin"
)

func HandleFileUpload(ctx *gin.Context, formFieldName string, uploadDir string) ([]models.Media, error) {
	if formFieldName == "" {
		return nil, nil // No file to upload is not an error
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err

	}

	files := form.File[formFieldName]
	if len(files) == 0 {
		return nil, fmt.Errorf("no file found in form field: %s", formFieldName)
	}

	var mediaList []models.Media
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(files))

	for index, file := range files {
		wg.Add(1)
		go func(idx int, file *multipart.FileHeader) {
			defer wg.Done()
			ext := filepath.Ext(strings.ToLower(file.Filename))
			fileName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))

			mediaType := "image"
			if ext == ".mp4" || ext == ".avi" || ext == ".mov" {
				mediaType = "video"
			}

			newFilename := fmt.Sprintf("%s%d%s", fileName,time.Now().UnixNano(), ext)

			tempPath := uploadDir + "/" + mediaType + "s"

			// Create a directory for the upload if it doesn't exist
			if err := os.MkdirAll(tempPath, 0755); err != nil {
				errChan <- fmt.Errorf("failed to create upload directory: %w", err)
				return
			}

			savePath := filepath.Join(tempPath, newFilename)
			if err := ctx.SaveUploadedFile(file, savePath); err != nil {
				errChan <- fmt.Errorf("failed to save uploaded file: %w", err)
				return
			}
			host := ctx.Request.Host // "127.0.0.1:9100" hoặc "example.com"

			// Tạo URL đầy đủ
			fileURL := fmt.Sprintf("http://%s/uploads/%ss/%s", host, mediaType, newFilename)
			mu.Lock()

			mediaList = append(mediaList, models.Media{
				URL:  fileURL,
				Type: mediaType,
			})

			mu.Unlock()

		}(index, file)
	}

	wg.Wait()
	close(errChan)

	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return mediaList, fmt.Errorf("encountered %d errors during upload: %s",
			len(errors), strings.Join(errors, "; "))
	}

	return mediaList, nil
}

func HandlerFileDeleted(ctx *gin.Context, mediaList *[]models.Media) error {
	if len(*mediaList) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(*mediaList))

	for _, media := range *mediaList {
		wg.Add(1)
		go func(media models.Media) {
			defer wg.Done()
			mediaType := media.Type
			parseUrl, err := url.Parse(media.URL)

			if err != nil {
				fileUrl := fmt.Sprintf("./uploads/%ss/%s", mediaType, media.URL)
				if err := os.Remove(fileUrl); err != nil {
					errChan <- fmt.Errorf("failed to delete file: %w", err)
					return
				}
			}

			// Lấy tên file từ URL
			filepath := parseUrl.Path
			filename := filepath[strings.LastIndex(filepath, "/")+1:]
			// Xóa file khỏi server
			fileUrl := fmt.Sprintf("./uploads/%ss/%s", mediaType, filename)
			if err := os.Remove(fileUrl); err != nil {
				errChan <- fmt.Errorf("failed to delete file: %w", err)
				return
			}
			// Xóa file khỏi mediaList

			mu.Lock()
			errChan <- nil
			mu.Unlock()
		}(media)

	}

	wg.Wait()
	close(errChan)

	var errors []string
	for err := range errChan {
		if err != nil {
			errors = append(errors, err.Error())
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during file deletion: %s",
			len(errors), strings.Join(errors, "; "))
	}
	return nil
}
