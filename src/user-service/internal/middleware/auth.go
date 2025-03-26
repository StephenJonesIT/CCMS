package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-service/common"
	"user-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(userRepo repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract token từ header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, common.NewErrorResponse("Authorization header missing"))
            return
        }

        // 2. Validate Bearer token format
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, common.NewErrorResponse("Invalid token format"))
            return
        }

        // 3. Verify JWT
        userID, err := common.VerifyToken(tokenParts[1])
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, common.NewErrorResponse("Invalid or expired token"))
            return
        }

        // 4. Verify user exists (optional but recommended)
        exists, err := userRepo.GetUserByID(userID)
        if err != nil || exists == nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, common.NewErrorResponse("User not found"))
            return
        }

        // 5. Set context for RBAC
        c.Set("userID", userID)
        c.Set("userRepo", userRepo) // Truyền userRepo qua context để tái sử dụng
    }
}

func RBACMiddleware(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Lấy userID từ context (đã được AuthMiddleware set)
        userID, exists := c.Get("userID")
        if !exists {
            c.AbortWithStatusJSON(http.StatusInternalServerError, common.NewErrorResponse("User context missing"))
            return
        }

        // 2. Lấy userRepo từ context (tránh truyền tham số trùng)
        userRepo, ok := c.MustGet("userRepo").(repository.UserRepository)
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError, common.NewErrorResponse("System error"))
            return
        }

        // 3. Kiểm tra quyền với cache (ví dụ dùng Redis)
        cacheKey := fmt.Sprintf("permission:%s:%s", userID, permission)
        if cached, err := common.GetCache(cacheKey); err == nil {
            if cached == "true" {
                c.Next()
                return
            }
            c.AbortWithStatusJSON(http.StatusForbidden, common.NewErrorResponse("Insufficient permissions"))
            return
        }

        // 4. Query database nếu không có cache
        hasPerm := userRepo.HasPermission(userID.(uuid.UUID), permission)
        if !hasPerm {
            c.AbortWithStatusJSON(http.StatusForbidden, common.NewErrorResponse("Insufficient permissions"))
            return
        }

        // 5. Cache kết quả (ví dụ: 5 phút)
        common.SetCache(cacheKey, strconv.FormatBool(hasPerm), 5*time.Minute)

        if !hasPerm {
            c.AbortWithStatusJSON(http.StatusForbidden, common.NewErrorResponse("Insufficient permissions"))
            return
        }

        c.Next()
    }
}