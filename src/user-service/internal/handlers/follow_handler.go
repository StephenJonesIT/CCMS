package handlers

import (
	"net/http"

	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/business"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FollowHandler struct {
	followService business.FollowService
}

func NewFollowHandler(followService business.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

// Follow godoc
// @Summary Follow a user
// @Description Create a follow relationship between two users
// @Tags Follow
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param follow body models.Follows true "Follow object"
// @Success 200 {object} common.Response "Follow created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid input"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Security Bearer Token
// @Router /follow [post]
func (h *FollowHandler) Follow(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"handler":  "FollowHandler",
		"endpoint": "Follow",
		"method":   c.Request.Method,
		"path":     c.FullPath(),
		"client":   c.ClientIP(),
	}).Info("Request received")

	var follow models.Follows
	if err := c.ShouldBindJSON(&follow); err != nil {
		logrus.Warn("Invalid parameters follow")
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid parameters follow"))
		return
	}

	if err := h.followService.Follow(c.Request.Context(), &follow); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

	logrus.Info("Follow created successfully")
	c.JSON(http.StatusOK, common.NewResponse("Follow created successfully"))
}

// Unfollow godoc
// @Summary Unfollow a user
// @Description Remove a follow relationship between two users
// @Tags Follow
// @Accept json
// @Produce json
// @Security Bearer Token
// @Param Authorization header string true "Bearer token"
// @Param follow body models.Follows true "Follow object"
// @Success 200 {object} common.Response "Unfollowed successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid input"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /unfollow [post]
func (h *FollowHandler) Unfollow(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"handler":  "FollowHandler",
		"endpoint": "UnFollow",
		"method":   c.Request.Method,
		"path":     c.FullPath(),
		"client":   c.ClientIP(),
	}).Info("Request received")
	var follow models.Follows
	if err := c.ShouldBindJSON(&follow); err != nil {
		logrus.Warn("Invalid parameters unfollow")
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid parameters unfollow"))
		return
	}

	if err := h.followService.Unfollow(c.Request.Context(), &follow); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

	logrus.Info("UnFollow successfully")
	c.JSON(http.StatusOK, common.NewResponse("UnFollow successfully"))
}

// GetFollowers godoc
// @Summary Get followers of a user
// @Description Retrieve a list of followers for a given user
// @Tags Follow
// @Produce json
// @Security Bearer Token
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User ID"
// @Success 200 {object} common.Response "List of followers"
// @Failure 400 {object} common.ErrorResponse "Invalid input"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /followers/{id} [get]
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"handler":  "FollowHandler",
		"endpoint": "UnFollow",
		"method":   c.Request.Method,
		"path":     c.FullPath(),
		"client":   c.ClientIP(),
	}).Info("Request received")

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		logrus.Warn("Invalid user ID")
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid user ID"))
		return
	}

	followers, err := h.followService.GetFollowers(c.Request.Context(), userID)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

	logrus.Info("Retrieve followers successfully")
	c.JSON(http.StatusOK, common.NewResponse(followers))
}

// GetFollowees godoc
// @Summary Get followees of a user
// @Description Retrieve a list of followees for a given user
// @Tags Follow
// @Produce json
// @Security Bearer Token
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User ID"
// @Success 200 {object} common.Response "List of followees"
// @Failure 400 {object} common.ErrorResponse "Invalid input"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /followees/{id} [get]
func (h *FollowHandler) GetFollowees(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		logrus.Warn("Invalid user ID")
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid user ID"))
		return
	}

	followees, err := h.followService.GetFollowees(c.Request.Context(), userID)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

	logrus.Info("Retrieve followees successfully")
	c.JSON(http.StatusOK, common.NewResponse(followees))
}
