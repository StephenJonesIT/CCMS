package handlers

import (
	"net/http"

	"github.com/StephenJonesIT/CCMS/src/post-service/business"
	"github.com/StephenJonesIT/CCMS/src/post-service/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeHandler struct {
	likeService business.LikeService
}

func NewLikeHandler(likeService business.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// ToggleLike godoc
// @Summary Toggle like on a post or comment
// @Description Toggle like status for a post or comment. Returns current like status after toggle.
// @Tags likes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Target ID (hex-encoded MongoDB ObjectID)" example(507f1f77bcf86cd799439011)
// @Param type query string true "Target type (post or comment)" example(post)
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Success 200 {object} object "Like status after toggle" { "liked": boolean }
// @Failure 400 {object} common.ErrorResponse "Invalid target ID or type"
// @Failure 401 "Unauthorized - Missing or invalid JWT token"
// @Failure 500 {object} common.ErrorResponse "Failed to toggle like"
// @Router /likes/{id} [post]
func (h *LikeHandler) ToggleLike(ctx *gin.Context) {
	targetID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "CommentHandler",
			"endpoint": "ToggleLike",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid post ID or comment ID")
		// Trả về lỗi 400 nếu postID không hợp lệ
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid post ID or comment ID"))
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)

	targetType := ctx.Query("type")
	if targetType != "post" && targetType != "comment" {
		logrus.WithFields(logrus.Fields{
			"handler":  "LikeHandler",
			"endpoint": "ToggleLike",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Invalid target type")
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("invalid target type"))
		return
	}

	isLiked, err := h.likeService.ToggleLike(ctx.Request.Context(), userID, targetID, targetType)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "LikeHandler",
			"endpoint": "ToggleLike",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
		}).Error("Faild to toggle like")
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
	ctx.JSON(http.StatusOK, gin.H{"liked": isLiked})
}
