/*
 * @File: main.go
 * @Description: Creates HTTP server & API groups of the User Service
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package main

import (
	"time"
	"user-service/common"
	"user-service/config"
	_ "user-service/docs"
	"user-service/internal/business"
	"user-service/internal/handlers"
	"user-service/internal/middleware"
	"user-service/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Main struct {
	router *gin.Engine
}

func (m *Main) initServe(r *repository.UserRepoImpl,h *handlers.UserHandler) {
	m.router = gin.Default()
	m.router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // For development only
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    
    // Swagger route
    m.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, 
        ginSwagger.URL("/swagger/doc.json"), // Explicit URL
        ginSwagger.DefaultModelsExpandDepth(-1),
    ))


	v1 := m.router.Group("/api/v1")
	{
			v1.POST("/login", h.Login)
			v1.POST("/register", h.Register)

			authGroup := v1.Group("")
			authGroup.Use(middleware.AuthMiddleware(r))
			userGroup := authGroup.Group("/users")
			{
				userGroup.GET("",middleware.RBACMiddleware("Quản lý người dùng"), h.ListUser)
			}

			profileGroup := authGroup.Group("/profiles")
			{
				profileGroup.GET("",h.ListProfile)
			}
	}
	

	m.router.Run(config.Config.Port)
}
 
// @title UserManagement Service API Document
// @version 1.0
// @description List APIs of UserManagement Service
// @termsOfService http://swagger.io/terms/

// @host 127.0.0.1:9000
// @BasePath /api/v1
// @schemes http https
func main(){
	m := Main{}
	err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
	config.ConfigDatabase()

	if err := config.LoadConfig("config/config.json"); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	repoUser := repository.NewUserRepository(config.DB)
	serviceUser := business.NewUserService(repoUser)
	handler := handlers.NewUserHandler(serviceUser)
	
	common.InitRedis("localhost:6379", "", 0)
	m.initServe(repoUser, handler)
}