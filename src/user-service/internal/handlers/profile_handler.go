package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @Summary Get list of profiles
// @Description Get paginated list of profiles
// @Tags profiles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /profiles [get]
func (h *UserHandler) ListProfile(ctx *gin.Context) {
	log.WithFields(log.Fields{
		"handler":  "UserHandler",
		"endpoint": "ListProfile",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Request received")

	var paging common.Paging
	if err := ctx.ShouldBind(&paging); err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"parameters":  ctx.Request.URL.Query(),
			"status_code": http.StatusBadRequest,
		}).Warn("Invalid paging parameters")

		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid paging parameter"))
		return
	}

	result, err := h.service.GetListProfile(&paging)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"stack_trace": string(debug.Stack()),
			"status_code": http.StatusInternalServerError,
		}).Error("Failed to retrieve user list")

		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

	// Log successful response (without sensitive data)
	log.WithFields(log.Fields{
		"total_records":    paging.Total,
		"returned_records": len(result),
		"total_pages":      paging.Page,
		"status_code":      http.StatusOK,
	}).Info("Successfully processed request")
	ctx.JSON(http.StatusOK, common.NewDetailResponse(result, paging))
}

// CreateProfile godoc
// @Summary Create a new user profile
// @Description Create user profile with image upload and profile data
// @Tags profiles
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Profile image file (max 5MB)"
// @Param profile formData string true "Profile data in JSON format"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} common.Response "Successfully created profile"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 413 {object} common.ErrorResponse "File too large"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /profiles [post]
func (h *UserHandler) CreateProfile(ctx *gin.Context) {
	log.WithFields(log.Fields{
		"handler":  "UserHandler",
		"endpoint": "CreateProfile",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Request received")

	filename, err := common.HandleFileUpload(ctx, "image", "./uploads")
    if err != nil {
        switch {
        case errors.Is(err, http.ErrMissingFile):
            // File not provided is acceptable for update
        case strings.Contains(err.Error(), "file size exceeds"):
            log.Warn(err.Error())
            ctx.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error()))
            return
        case strings.Contains(err.Error(), "only .jpg, .jpeg or .png"):
            log.Warn(err.Error())
            ctx.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error()))
            return
        default:
            log.WithError(err).Error("File upload error")
            ctx.JSON(http.StatusInternalServerError, 
                common.NewErrorResponse("Failed to process uploaded file"))
            return
        }
    }

	profileJSON := ctx.PostForm("profile")
	if profileJSON == "" {
        log.Error("Missing profile JSON data")
        ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Missing profile data"))
        return
    }

	var profile models.Profile
    if err := json.Unmarshal([]byte(profileJSON), &profile); err != nil {
        log.WithError(err).Error("Invalid profile JSON format")
        ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid profile data format"))
        return
    }

	profile.Picture_URL = "/uploads/" + filename
	if err := h.service.CreateProfile(&profile); err != nil {
		log.WithError(err).Error("Failed to save profile to database")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to create profile"))
        return
	}

	log.WithFields(log.Fields{
		"profile_id": profile.ProfileID,      // Nếu có ID
		"file_name":  filename,
		"client_ip":  ctx.ClientIP(),
	}).Info("Profile created successfully")
	ctx.JSON(http.StatusOK, common.NewCreateOrUpdate("Profile created successfully", profile))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile with image upload and profile data
// @Tags profiles
// @Accept multipart/form-data
// @Produce json
// @Param id_profile path int true "Profile ID"
// @Param image formData file false "New profile image file (max 5MB)"
// @Param profile formData string true "Profile data in JSON format"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} common.Response "Successfully updated profile"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 413 {object} common.ErrorResponse "File too large"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /profiles/{id_profile} [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	log.WithFields(log.Fields{
		"handler":  "UserHandler",
		"endpoint": "UpdateProfile",
		"method":   ctx.Request.Method,
		"path":     ctx.FullPath(),
		"client":   ctx.ClientIP(),
	}).Info("Request received")

	idProfile, err := strconv.Atoi(ctx.Param("id_profile"))
	if err != nil {
		log.WithError(err).Warn("Invalid idProfile parameter")
		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid idProfile parameter"))
		return
	}

	filename, err := common.HandleFileUpload(ctx, "image", "./uploads")
    if err != nil {
        switch {
        case errors.Is(err, http.ErrMissingFile):
            // File not provided is acceptable for update
        case strings.Contains(err.Error(), "file size exceeds"):
            log.Warn(err.Error())
            ctx.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error()))
            return
        case strings.Contains(err.Error(), "only .jpg, .jpeg or .png"):
            log.Warn(err.Error())
            ctx.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error()))
            return
        default:
            log.WithError(err).Error("File upload error")
            ctx.JSON(http.StatusInternalServerError, 
                common.NewErrorResponse("Failed to process uploaded file"))
            return
        }
    }

	profileJSON := ctx.PostForm("profile")
	if profileJSON == "" {
        log.Error("Missing profile JSON data")
        ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Missing profile data"))
        return
    }

	var profile models.Profile
    if err := json.Unmarshal([]byte(profileJSON), &profile); err != nil {
        log.WithError(err).Error("Invalid profile JSON format")
        ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid profile data format"))
        return
    }

	// If file was uploaded, update profile picture URL
    if filename != "" {
        // Delete old picture if exists
        if profile.Picture_URL != "" {
            if err := common.DeleteFile(filepath.Base(profile.Picture_URL)); err != nil {
                log.WithError(err).Warn("Failed to delete old profile picture")
            }
        }
        profile.Picture_URL = "/uploads/"+ filename
    }

	profile.ProfileID = int64(idProfile)
	if err := h.service.UpdateProfile(&profile); err != nil {
		log.WithError(err).Error("Failed to save profile to database")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("Failed to update profile"))
        return
	}

	log.WithFields(log.Fields{
		"profile_id": profile.ProfileID,      // Nếu có ID
		"file_name":  filename,
		"client_ip":  ctx.ClientIP(),
	}).Info("Profile updated successfully")
	ctx.JSON(http.StatusOK, common.NewCreateOrUpdate("Profile updated successfully", profile))
}
