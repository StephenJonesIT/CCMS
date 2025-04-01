package middlewares

import (
	"net/http"

	"github.com/StephenJonesIT/CCMS/src/post-service/client"
	"github.com/StephenJonesIT/CCMS/src/post-service/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func AuthMiddleware(AuthClient *client.AuthClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx)
		if token == "" {
			log.WithFields(log.Fields{
				"handler": "AuthMiddleware",
				"endpoint": "AuthMiddleware",
				"method": ctx.Request.Method,
				"path": ctx.FullPath(),
				"client": ctx.ClientIP(),
			}).Error("Authorization header missing")
			ctx.JSON(401, common.NewErrorResponse("Authorization header missing"))
			ctx.Abort()
			return
		}

		claims, err := AuthClient.VerifyToken(token)
		if err != nil {
			log.WithFields(log.Fields{
				"handler": "AuthMiddleware",
				"endpoint": "AuthMiddleware",
				"method": ctx.Request.Method,
				"path": ctx.FullPath(),
				"client": ctx.ClientIP(),
			}).Error("Invalid or expired token")
			ctx.JSON(401, common.NewErrorResponse("Invalid or expired token"))
			ctx.Abort()
			return
		}

		if claims.UserId == "" {
			log.WithFields(log.Fields{
				"handler": "AuthMiddleware",
				"endpoint": "AuthMiddleware",
				"method": ctx.Request.Method,
				"path": ctx.FullPath(),
				"client": ctx.ClientIP(),
			}).Error("Invalid token claims")
			ctx.JSON(401, common.NewErrorResponse("Invalid token claims"))
			ctx.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserId)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.NewErrorResponse("invalid user ID format"))
            return
        }
		ctx.Set("userID", userID)
		ctx.Set("roleName", claims.RoleName)
		ctx.Next()
	}
}

func PermissionMiddleware(authClient *client.AuthClient, resource string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx)
		if token == "" {
			log.WithFields(log.Fields{
				"handler": "PermissionMiddleware",
				"endpoint": "PermissionMiddleware",
				"method": ctx.Request.Method,
				"path": ctx.FullPath(),
				"client": ctx.ClientIP(),
			}).Error("Authorization header missing")
			ctx.JSON(401, common.NewErrorResponse("Authorization header missing"))
			ctx.Abort()
			return
		}

		_, err := authClient.CheckPermission(token, resource)
		if err != nil {
			log.WithFields(log.Fields{
				"handler": "PermissionMiddleware",
				"endpoint": "PermissionMiddleware",
				"method": ctx.Request.Method,
				"path": ctx.FullPath(),
				"client": ctx.ClientIP(),
			}).Error("Permission denied")
			ctx.JSON(403, common.NewErrorResponse("Permission denied"))
			ctx.Abort()
			return
		}
		// Nếu có quyền, tiếp tục xử lý request
		ctx.Next()
	}
}

func extractToken(c *gin.Context) string {
	// Lấy từ header
	token := c.GetHeader("Authorization")
	if token != "" {
		// Loại bỏ "Bearer " nếu có
		if len(token) > 7 && token[:7] == "Bearer " {
			return token[7:]
		}
		return token
	}

	// Hoặc lấy từ query parameter
	token = c.Query("token")
	if token != "" {
		return token
	}

	// Hoặc lấy từ cookie
	token, _ = c.Cookie("token")
	return token
}