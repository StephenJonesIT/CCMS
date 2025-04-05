/*
 * @File: handlers.user_handler.go
 * @Description: Implements User API logic functions
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package handlers

import (
	"net/http"
	"runtime/debug"

	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/business"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	service business.UserService
}

func NewUserHandler(ser business.UserService) *UserHandler {
	return &UserHandler{
		service: ser,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and get access token
// @Tags authentication
// @Accept json
// @Produce json
// @Param loginRequest body common.LoginRequest true "Login credentials"
// @Success 200 {object} common.LoginResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
    
    log.WithFields(log.Fields{
        "handler":  "UserHandler",
        "endpoint": "Login",
        "method":   ctx.Request.Method,
        "clientIP": ctx.ClientIP(),
    }).Info("Login request received")

    var LoginRequest common.LoginRequest

    if err := ctx.ShouldBindJSON(&LoginRequest); err != nil {
        ctx.JSON(http.StatusBadRequest,  common.NewErrorResponse("invalid request"))
        return
    }

    user, err := h.service.Login(LoginRequest.Username, LoginRequest.Password)
    if err != nil {
        if err.Error() == "user not found" {
            log.WithFields(log.Fields{
                "username": LoginRequest.Username,
                "status":   http.StatusNotFound,
            }).Warn("Login failed - user not found")
            ctx.JSON(http.StatusNotFound, common.NewErrorResponse("user not found"))

        } else if err.Error() == "invalid password" {
            log.WithFields(log.Fields{
                "username": LoginRequest.Username,
                "status":   http.StatusBadRequest,
            }).Warn("Login failed - invalid password")
            ctx.JSON(http.StatusInternalServerError,common.NewErrorResponse("invalid password"))

        } else {
            log.WithFields(log.Fields{
                "error":       err.Error(),
                "stack_trace": string(debug.Stack()),
                "status":      http.StatusInternalServerError,
            }).Error("Login processing error")
            ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("login failed"))
        }
        return
    }

    token, err := common.GenerateToken(user.UserID, user.Role.RoleName)
    if err != nil {
        log.WithFields(log.Fields{
            "user_id":    user.UserID,
            "error":      err.Error(),
            "stack_trace": string(debug.Stack()),
            "status":     http.StatusInternalServerError,
        }).Error("Failed to generate JWT token")

        ctx.JSON(http.StatusInternalServerError,common.NewErrorResponse("failed to generate token"))
        return
    }

    log.WithFields(log.Fields{
        "user_id":  user.UserID,
        "username": user.UserName,
        "status":   http.StatusOK,
    }).Info("Login successful")
    ctx.JSON(http.StatusOK, common.NewLoginResponse("Login success", token))
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body models.UserRegister true "User registration data"
// @Success 201 {object} common.Response
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /register [post]
func(h *UserHandler) Register(ctx *gin.Context) {
    log.WithFields(log.Fields{
        "handler":  "UserHandler",
        "endpoint": "Register",
        "method":   ctx.Request.Method,
        "path":     ctx.FullPath(),
        "client":   ctx.ClientIP(),
    }).Info("Request received")
	var user models.UserRegister

	if err := ctx.ShouldBindJSON(&user); err != nil{
        log.WithFields(log.Fields{
            "error":       err.Error(),
            "parameters":  ctx.Request.URL.Query(),
            "status_code": http.StatusBadRequest,
        }).Warn("Invalid user parameters")

		ctx.JSON(http.StatusBadRequest, common.NewErrorResponse("invalid request"))
        return
	}

	if err := h.service.Register(&user); err != nil {
        log.WithFields(log.Fields{
            "error":       err.Error(),
            "stack_trace": string(debug.Stack()),
            "status_code": http.StatusInternalServerError,
        }).Error("Registration failed")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error()))
		return
	}

    log.WithFields(log.Fields{
        "user_id":  user.UserID,
        "username": user.UserName,
        "status":   http.StatusOK,
    }).Info("User registered successfully")
	ctx.JSON(http.StatusOK, common.NewResponse(user))
}

// @Summary Get list of users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /users [get]
func(h *UserHandler) ListUser(ctx *gin.Context) {
    log.WithFields(log.Fields{
        "handler":  "UserHandler",
        "endpoint": "ListUser",
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

    result, err := h.service.GetListUser(&paging)
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
        "total_records": paging.Total,
        "returned_records": len(result),
        "total_pages":   paging.Page,
        "status_code":   http.StatusOK,
    }).Info("Successfully processed request")
    ctx.JSON(http.StatusOK, common.NewDetailResponse(result, paging))
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for authenticated user (requires current password)
// @Tags authentication
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param changePasswordRequest body common.LoginRequest true "Password change details"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /change-password [post]
func(h *UserHandler) ChangePassword(ctx *gin.Context) {
    log.WithFields(log.Fields{
        "handler":  "UserHandler",
        "endpoint": "ChangePassword",
        "method":   ctx.Request.Method,
        "path":     ctx.FullPath(),
        "client":   ctx.ClientIP(),
    }).Info("Request received")

    var request common.LoginRequest

    if err := ctx.ShouldBindJSON(&request); err != nil {
        log.WithFields(log.Fields{
            "error":       err.Error(),
            "parameters":  ctx.Request.URL.Query(),
            "status_code": http.StatusBadRequest,
        }).Warn("Invalid request parameters")
    }

    
    if err := h.service.ChangePassword(request.Username, request.Password); err != nil {
        switch err.Error() { // Compare error message strings
        case "user not found":
            log.WithFields(log.Fields{
                "username": request.Username,
                "status":   http.StatusNotFound,
            }).Warn("ChangePassword failed - user not found")
            ctx.JSON(http.StatusNotFound, common.NewErrorResponse("user not found"))
        case "invalid password":
            log.WithFields(log.Fields{
                "username": request.Username,
                "status":   http.StatusUnauthorized,
            }).Warn("ChangePassword failed - invalid password")
            ctx.JSON(http.StatusUnauthorized, common.NewErrorResponse("current password is incorrect"))
        case "failed to hash password":
            log.WithFields(log.Fields{
                "username": request.Username,
                "status":   http.StatusInternalServerError,
                "error":    err.Error(),
            }).Error("ChangePassword failed - password hashing error")
            ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("password processing failed"))
        default:
            log.WithFields(log.Fields{
                "username": request.Username,
                "status":   http.StatusInternalServerError,
                "error":    err.Error(),
            }).Error("ChangePassword failed - unexpected error")
            ctx.JSON(http.StatusInternalServerError, common.NewErrorResponse("failed to change password"))
        }
        return
    }

    log.WithFields(log.Fields{
        "username": request.Username,
        "status":   http.StatusOK,
    }).Info("ChangePassword successful")
    ctx.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
