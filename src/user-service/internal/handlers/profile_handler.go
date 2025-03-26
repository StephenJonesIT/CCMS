package handlers

import (
	"net/http"
	"runtime/debug"
	"user-service/common"

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
func(h *UserHandler) ListProfile(ctx *gin.Context) {
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
        "total_records": paging.Total,
        "returned_records": len(result),
        "total_pages":   paging.Page,
        "status_code":   http.StatusOK,
    }).Info("Successfully processed request")
    ctx.JSON(http.StatusOK, common.NewDetailResponse(result, paging))
}