package routes

import (
	"log"
	"time"

	"github.com/StephenJonesIT/CCMS/src/post-service/client"
	"github.com/StephenJonesIT/CCMS/src/post-service/config"
	"github.com/StephenJonesIT/CCMS/src/post-service/handlers"
	"github.com/StephenJonesIT/CCMS/src/post-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Main struct {
	router *gin.Engine
}

func NewMain() *Main {
	return &Main{
		router: gin.Default(),
	}
}

func (m *Main) InitServe(h *handlers.PostHandler) error {
	authClient, err := client.NewAuthClient("0.0.0.0:50051", time.Second*5)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}
	defer authClient.Close()

	m.router.Use(gin.Logger())
	m.router.Use(gin.Recovery())


	// Set up CORS middleware
	m.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // For development only
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	m.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"), // Explicit URL
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	m.router.Static("/uploads", "./uploads")
	v1 := m.router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.Use(middlewares.AuthMiddleware(authClient))
			SetUpUserRoutes(user, h)
		}


		// admin := v1.Group("/admin")
		// {

		// }
	}

	m.router.Run(config.Config.Port)
	return nil
}