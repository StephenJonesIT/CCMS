package routes

import (
	"github.com/StephenJonesIT/CCMS/src/post-service/handlers"
	"github.com/gin-gonic/gin"
)

func SetUpUserRoutes(router *gin.RouterGroup, postHandler *handlers.PostHandler) {
	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts/:id", postHandler.GetPostByID)
	router.PUT("/posts/:id", postHandler.UpdatePost)
	router.DELETE("/posts/:id", postHandler.DeletePost)
	// router.GET("/posts/user/:userId", postHandler.GetPostsByUser)
	// router.GET("/posts/status", postHandler.GetPostsByStatus)
}