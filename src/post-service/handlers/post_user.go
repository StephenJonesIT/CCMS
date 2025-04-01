package handlers

import (
	"net/http"
	"time"
	page "github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/post-service/business"
	"github.com/StephenJonesIT/CCMS/src/post-service/common"
	"github.com/StephenJonesIT/CCMS/src/post-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type PostHandler struct {
	postService business.PostService
}

func NewPostHandler(postService business.PostService) *PostHandler{
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost godoc
// @Summary Create a new post with media
// @Description Create a new post with multiple media files (images/videos)
// @Tags user posts
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param title formData string false "Post title"
// @Param content formData string false "Post content"
// @Param categories formData []string false "Post categories" collectionFormat(multi)
// @Param tags formData []string false "Post tags" collectionFormat(multi)
// @Param status formData string true "Post status" Enums(draft, pending)
// @Param media formData []file true "Media files (images/videos)"
// @Security ApiKeyAuth
// @Success 201 {object} common.Response "Post created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request or file upload failed"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /user/posts [post]
func(h *PostHandler) CreatePost(ctx *gin.Context) {
	log.WithFields(log.Fields{
		"handler": "PostHandler",
		"endpoint": "CreatePost",
		"method": ctx.Request.Method,
		"path": ctx.FullPath(),
		"client": ctx.ClientIP(),
	}).Info("Request received")

	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusBadRequest, 
		}).Error("File too large")
		ctx.JSON(http.StatusBadRequest,common.NewErrorResponse("Failed to parse multipart form"))
		return
	}

	mediaList, err := common.HandleFileUpload(ctx, "media", "./uploads")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusBadRequest,
		}).Error("Failed to handle file upload")
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Failed to handle file upload"))
		return
	}

	post := &models.Post{
		Title:        &[]string{ctx.PostForm("title")}[0],
		Content:      &[]string{ctx.PostForm("content")}[0],
		AuthorID:     ctx.MustGet("userID").(uuid.UUID),
		Categories:   ctx.PostFormArray("categories"),
		Tags:         ctx.PostFormArray("tags"),
		Media:        mediaList,
		Status:       ctx.PostForm("status"),
		LikeCount:    0,
		CommentCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	} 

	result, err := h.postService.CreatePost(post) 
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to create post")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to create post"))
		return
	}


	log.WithFields(log.Fields{
		"status_code": http.StatusCreated,
		"parameter": ctx.Request.URL.Query(),
	}).Info("Post created successfully")
	ctx.JSON(http.StatusCreated, common.NewDetailResponse("Post created successfully", result))
}

// GetPostByID godoc
// @Summary Get a post by ID
// @Description Get a post by its ID
// @Tags user posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "Post ID"
// @Success 200 {object} common.Response "Post retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Post ID is required"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Security Bearer
// @Router /user/posts/{id} [get]
func (h *PostHandler) GetPostByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		log.WithFields(log.Fields{
			"handler": "PostHandler",
			"endpoint": "GetPostByID",
			"method": ctx.Request.Method,
			"path": ctx.FullPath(),
			"client": ctx.ClientIP(),
		}).Error("Post ID is required")
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Post ID is required"))
		return
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to get post by ID")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to get post by ID"))
		return
	}


	log.WithFields(log.Fields{
		"status_code": http.StatusOK,
	}).Info("Post retrieved successfully")
	ctx.JSON(http.StatusOK, common.NewDetailResponse("Post retrieved successfully", post))
}

// UpdatePost godoc
// @Summary Update a post by ID
// @Description Update an existing post with new data including optional media files
// @Tags user posts
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "Post ID"
// @Param title formData string false "Post title"
// @Param content formData string false "Post content"
// @Param categories formData []string false "Post categories" collectionFormat(multi)
// @Param tags formData []string false "Post tags" collectionFormat(multi)
// @Param status formData string true "Post status" Enums(draft, pending)
// @Param media formData []file false "Media files" collectionFormat(multi)
// @Success 200 {object} common.Response "Post updated successfully"
// @Failure 400 {object} common.ErrorResponse "Bad request"
// @Failure 404 {object} common.ErrorResponse "Post not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Security Bearer
// @Router /user/posts/{id} [put]
func (h *PostHandler) UpdatePost(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		log.WithFields(log.Fields{
			"handler": "PostHandler",
			"endpoint": "UpdatePost",
			"method": ctx.Request.Method,
			"path": ctx.FullPath(),
			"client": ctx.ClientIP(),
		}).Error("Post ID is required")
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Post ID is required"))
		return
	}

	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusBadRequest,
		}).Error("File too large")
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Failed to parse multipart form"))
		return
	}

	mediaList, err := common.HandleFileUpload(ctx, "media", "./uploads")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusBadRequest,
		}).Error("Failed to handle file upload")
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Failed to handle file upload"))
		return
	}

	post, err := h.postService.GetPostByID(id) 
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to get post by ID")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to get post by ID"))
		return
	}

	if post == nil {
		log.WithFields(log.Fields{
			"error": "Post not found",
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusNotFound,
		}).Error("Post not found")
		ctx.JSON(http.StatusNotFound, common.NewErrorResponse("Post not found"))
		return
	}

	if err := common.HandlerFileDeleted(ctx, &post.Media); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to delete file")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to delete file"))
		return
	}

	// Update the post fields with the new values from the request
	post.Title = &[]string{ctx.PostForm("title")}[0]
	post.Content = &[]string{ctx.PostForm("content")}[0]
	post.Categories = ctx.PostFormArray("categories")
	post.Tags = ctx.PostFormArray("tags")
	post.Status = ctx.PostForm("status")
	post.Media = mediaList


	err = h.postService.UpdatePost(id, post)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to update post")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to update post"))
		return
	}

	log.WithFields(log.Fields{
		"status_code": http.StatusOK,
	}).Info("Post updated successfully")
	ctx.JSON(http.StatusOK, common.NewDetailResponse("Post updated successfully", nil))
}

// DeletePost godoc
// @Summary Delete a post by ID
// @Description Delete a post by its ID
// @Tags user posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "Post ID"
// @Success 200 {object} common.Response "Post deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Post ID is required"
// @Failure 404 {object} common.ErrorResponse "Post not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Security Bearer
// @Router /user/posts/{id} [delete]
func (h *PostHandler) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		log.WithFields(log.Fields{
			"handler": "PostHandler",
			"endpoint": "DeletePost",
			"method": ctx.Request.Method,
			"path": ctx.FullPath(),
			"client": ctx.ClientIP(),
		}).Error("Post ID is required")
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Post ID is required"))
		return
	}

	post, err := h.postService.GetPostByID(id) 
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to get post by ID")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to get post by ID"))
		return
	}

	if post == nil {
		log.WithFields(log.Fields{
			"error": "Post not found",
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusNotFound,
		}).Error("Post not found")
		ctx.JSON(http.StatusNotFound, common.NewErrorResponse("Post not found"))
		return
	}

	if err := common.HandlerFileDeleted(ctx, &post.Media); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to delete file")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to delete file"))
		return
	}

	err = h.postService.DeletePost(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"parameter": ctx.Request.URL.Query(),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to delete post")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to delete post"))
		return
	}

	log.WithFields(log.Fields{
		"status_code": http.StatusOK,
	}).Info("Post deleted successfully")
	ctx.JSON(http.StatusOK, common.NewResponse("Post deleted successfully"))
}

func (h *PostHandler) GetPostsByUser(ctx *gin.Context) {
	var page page.Paging

	if err := ctx.ShouldBindJSON(&page); err != nil {
		log.WithFields(log.Fields{
			
		})
	}
}