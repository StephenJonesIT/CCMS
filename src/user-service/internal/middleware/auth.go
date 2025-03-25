package middleware

import (
	"net/http"
	"user-service/common"
	"user-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(userRepo repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        userID, err := common.VerifyToken(token) // Implement JWT verification
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        c.Set("userID", userID)
        c.Next()
    }
}

func RBACMiddleware(userRepo repository.UserRepository, permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        if !userRepo.HasPermission(userID.(uuid.UUID), permission) {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            return
        }

        c.Next()
    }
}