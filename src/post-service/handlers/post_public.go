package handlers

import (
	"net/http"

	p "github.com/StephenJonesIT/CCMS/src/post-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *PostHandler) GetPostsByStatus(ctx *gin.Context) {

	var page common.Paging
	if err := ctx.ShouldBindQuery(&page); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "PostHandler",
			"endpoint": "GetPostsByStatus",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
	}).Error("Invalid query parameters")
		ctx.JSON(http.StatusBadRequest, common.NewDetailResponse("Invalid query parameters", nil))
		return
	}


	posts, err := h.postService.GetPostsByStatus(page, "published")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":  "PostHandler",
			"endpoint": "GetPostsByStatus",
			"method":   ctx.Request.Method,
			"path":     ctx.FullPath(),
			"client":   ctx.ClientIP(),
	}).Error("Failed to get posts by status")
		ctx.JSON(http.StatusInternalServerError, p.NewDetailResponse("Failed to get posts by status", nil))
		return
	}
	ctx.JSON(http.StatusOK, p.NewResponse(posts))
}