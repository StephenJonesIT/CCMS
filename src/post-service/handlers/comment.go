package handlers

import (
	"net/http"
	"strings"

	"github.com/StephenJonesIT/CCMS/src/post-service/business"
	"github.com/StephenJonesIT/CCMS/src/post-service/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	commentService business.CommentService
}

func NewCommentHandler(commentService business.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Creates a new comment on a specified post. Requires authentication (JWT token in Authorization header).
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param postId path string true "Hex-encoded MongoDB ObjectID of the post" example(507f1f77bcf86cd799439011)
// @Param request body object true "Comment content" { "content": "string" }
// @Success 201 {object} common.Response{data=models.Comment} "Comment created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid post ID or request body"
// @Failure 401 "Unauthorized - Missing or invalid JWT token"
// @Failure 500 {object} common.ErrorResponse "Failed to create comment"
// @Security BearerAuth
// @Router /posts/{postId}/comments [post]
func (h *CommentHandler) CreateComment(ctx *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(ctx.Param("postId"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "CreateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid post ID")
		// Trả về lỗi 400 nếu postID không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("invalid post ID"))
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "CreateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid request body")
		// Trả về lỗi 400 nếu request body không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error()))
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)

	comment, err := h.commentService.CreateComment(ctx.Request.Context(), postID, userID, req.Content)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "CreateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Failed to create comment")
		// Trả về lỗi 500 nếu có lỗi xảy ra trong quá trình tạo comment
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("failed to create comment"))
		return
	}

	logrus.WithFields(logrus.Fields{
		"handler":  "CommentHandler",
		"endpoint": "CreateComment",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Comment created successfully")
	ctx.JSON(http.StatusCreated, common.NewResponse(comment))
}

// UpdateComment godoc
// @Summary Update comment by id
// @Description Update an existing comment's content. Requires authentication (JWT token in Authorization header).
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param id path string true "Comment ID (hex-encoded MongoDB ObjectID)" example(507f1f77bcf86cd799439011)
// @Param postId path string true "Post ID (hex-encoded MongoDB ObjectID)" 
// @Param request body object true "Comment content" { "content": "string" }
// @Success 200 {object} common.Response{data=string} "Comment updated successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid comment ID or request body"
// @Failure 404 "Comment not found"
// @Failure 500 {object} common.ErrorResponse "Failed to update comment"
// @Router /posts/{postId}/comments/{id} [put]
func (h *CommentHandler) UpdateComment(ctx *gin.Context) {
	// Lấy commentID từ URL
	commentID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "UpdateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid comment ID")
		// Trả về lỗi 400 nếu commentID không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("invalid comment ID"))
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "UpdateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid request body")
		// Trả về lỗi 400 nếu request body không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error()))
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)
	if err := h.commentService.UpdateComment(ctx.Request.Context(), commentID, userID, req.Content); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "UpdateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Failed to update comment")
		// Trả về lỗi 500 nếu có lỗi xảy ra trong quá trình cập nhật comment
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("failed to update comment"))
		return
	}

	logrus.WithFields(logrus.Fields{
		"handler":  "CommentHandler",
		"endpoint": "UpdateComment",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Comment updated successfully")
	// Trả về thông báo thành công
	ctx.JSON(http.StatusOK, common.NewResponse("comment updated successfully"))
}

// DeleteComment godoc
// @Summary Delete comment by id
// @Description Delete an existing comment. Requires authentication (JWT token in Authorization header).
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param postId path string true "Post ID (hex-encoded MongoDB ObjectID)" example(507f1f77bcf86cd799439011)
// @Param id path string true "Comment ID (hex-encoded MongoDB ObjectID)" example(507f1f77bcf86cd799439012)
// @Success 200 {object} common.Response{data=string} "Comment deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid comment ID or post ID"
// @Failure 404 "Comment not found"
// @Failure 500 {object} common.ErrorResponse "Failed to delete comment"
// @Router /posts/{postId}/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(ctx *gin.Context) {
	commentID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "DeleteComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid comment ID")
		// Trả về lỗi 400 nếu commentID không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("invalid comment ID"))
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)
	if err := h.commentService.DeleteComment(ctx.Request.Context(), commentID, userID); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "DeleteComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Failed to delete comment")
		ctx.JSON(http.StatusInternalServerError,common.NewErrorResponse(err.Error()))
		return
	}

	logrus.WithFields(logrus.Fields{
		"handler":  "CommentHandler",
		"endpoint": "DeleteComment",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Comment deleted successfully")
	ctx.JSON(http.StatusOK, common.NewResponse("comment deleted successfully"))
}

// AddReply godoc
// @Summary Add a reply to a comment
// @Description Add a reply to an existing comment. Requires authentication (JWT token in Authorization header).
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Comment ID (hex-encoded MongoDB ObjectID)" example(507f1f77bcf86cd799439011)
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param request body object true "Comment content" { "content": "string" }
// @Success 201 {object} common.Response{data=string} "Reply added successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid comment ID or request body"
// @Failure 401 "Unauthorized - Missing or invalid JWT token"
// @Failure 404 "Comment not found"
// @Failure 500 {object} common.ErrorResponse "Failed to add reply"
// @Router /posts/{postId}/comments/{id}/replies [post]
func (h *CommentHandler) AddReply(ctx *gin.Context) {
	commentID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "AddReply",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid comment ID")
		// Trả về lỗi 400 nếu commentID không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("invalid comment ID"))
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "AddReply",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid request body")
		// Trả về lỗi 400 nếu request body không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error()))
		return
	}

	if err := h.commentService.AddReply(ctx.Request.Context(), commentID, userID, req.Content); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "AddReply",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Failed to add reply")
		// Trả về lỗi 500 nếu có lỗi xảy ra trong quá trình thêm reply
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

	logrus.WithFields(logrus.Fields{
		"handler":  "CommentHandler",
		"endpoint": "AddReply",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Reply added successfully")
	// Trả về thông báo thành công
	ctx.JSON(http.StatusCreated, common.NewResponse("reply added successfully"))
}

// GetComments godoc
// @Summary Get comments for a post
// @Description Get all comments for a specific post. Admin can see pending comments.
// @Tags comments
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param postId path string true "Post ID (hex-encoded MongoDB ObjectID)" example(507f1f77bcf86cd799439011)
// @Success 200 {object} common.Response{data=[]models.Comment} "List of comments"
// @Failure 400 {object} common.ErrorResponse "Invalid post ID"
// @Failure 500 {object} common.ErrorResponse "Failed to get comments"
// @Router /posts/{postId}/comments [get]
func (h *CommentHandler) GetComments(ctx *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(ctx.Param("postId"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "GetComments",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid post ID")
		// Trả về lỗi 400 nếu postID không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("invalid post ID"))
		return
	}

	// Chỉ admin mới có thể xem các comment pending
	roleName := ctx.MustGet("roleName").(string)
	showPending := false
	if strings.Compare("admin", roleName) == 0 {
		showPending = true
	}

	comments, err := h.commentService.GetCommentsByPostID(ctx.Request.Context(), postID, showPending)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "GetComments",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Failed to get comments")
		// Trả về lỗi 500 nếu có lỗi xảy ra trong quá trình lấy comment
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("failed to get comments"))
		return
	}

	logrus.WithFields(logrus.Fields{
		"handler":  "CommentHandler",
		"endpoint": "GetComments",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Comments retrieved successfully")
	// Trả về danh sách comment	
	ctx.JSON(http.StatusOK, common.NewResponse(comments))
}

// ModerateComment godoc
// @Summary Moderate a comment (admin only)
// @Description Approve or reject a comment (admin role required)
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param id path string true "Comment ID (hex-encoded MongoDB ObjectID)" example(507f1f77bcf86cd799439011)
// @Param request body object true "Moderation status (approved/rejected)"  { "status": "string" }
// @Success 200 {object} common.Response{data=string} "Comment moderated successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid comment ID or status"
// @Failure 404 "Comment not found"
// @Failure 500 {object} common.ErrorResponse "Failed to moderate comment"
// @Router /posts/{postId}/comments/{id}/moderate [put]
func (h *CommentHandler) ModerateComment(ctx *gin.Context) {
	commentID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "ModerateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid post ID")
		// Trả về lỗi 400 nếu postID không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("invalid post ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=approved rejected"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "ModerateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid request body")
		ctx.JSON(http.StatusBadRequest, common.NewResponse(err.Error()))
		return
	}

	if err := h.commentService.ModerateComment(ctx.Request.Context(), commentID, req.Status); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "ModerateComment",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Faild to moderate comment")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

	logrus.WithFields(logrus.Fields{
		"handler":  "CommentHandler",
		"endpoint": "ModerateComment",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Comments moderated successfully")

	ctx.JSON(http.StatusOK, common.NewResponse("Comment moderated successfully"))
}