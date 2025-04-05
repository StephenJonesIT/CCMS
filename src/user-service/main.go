/*
 * @File: main.go
 * @Description: Creates HTTP server & API groups of the User Service
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package main

import (
	"net"
	"os"
	"time"

	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/config"
	_"github.com/StephenJonesIT/CCMS/src/user-service/docs"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/business"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/handlers"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/middleware"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/repository"
	pb "github.com/StephenJonesIT/CCMS/src/user-service/proto"
	"github.com/StephenJonesIT/CCMS/src/user-service/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

type Main struct {
	router *gin.Engine
}

func (m *Main) initServe(
    r *repository.UserRepoImpl,
    h *handlers.UserHandler, 
    p *handlers.ProfileHandler,
    f *handlers.FollowHandler,
    ) error{
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
			v1.Static("/uploads", "./uploads")
			v1.POST("/login", h.Login)
			v1.POST("/register", h.Register)
			v1.POST("/change-password", h.ChangePassword)

			authGroup := v1.Group("")
			authGroup.Use(middleware.AuthMiddleware(r))
			userGroup := authGroup.Group("/users")
			{
				userGroup.GET("",middleware.RBACMiddleware("Quản lý người dùng"), h.ListUser)
			}

			profileGroup := authGroup.Group("/profiles")
			{
				profileGroup.GET("",middleware.RBACMiddleware("Quản lý người dùng"),p.ListProfile)
				profileGroup.POST("", p.CreateProfile)
				profileGroup.PUT(":id_profile", p.UpdateProfile)
			}

            followGroup := authGroup.Group("")
            {
                followGroup.POST("/follow", f.Follow)
                followGroup.POST("/unfollow", f.Unfollow)
                followGroup.GET("/followees/:id", f.GetFollowees)
                followGroup.GET("/followers/:id", f.GetFollowers)
            }
	}
	

	log.Infof("Starting HTTP server on %s", config.Config.Port)
    return m.router.Run(config.Config.Port)
}
 
// @title UserManagement Service API Document
// @version 1.0
// @description List APIs of UserManagement Service
// @termsOfService http://swagger.io/terms/

// @host 127.0.0.1:9000
// @BasePath /api/v1
// @schemes http https
func main() {
    // 1. Khởi tạo logging
    log.SetFormatter(&log.JSONFormatter{})
    
    // 2. Load config
    if err := godotenv.Load(".env"); err != nil {
        log.Warn("No .env file found, using environment variables")
    }
    
    if err := config.LoadConfig("config/config.json"); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // 3. Khởi tạo database
    config.ConfigDatabase()
    
    // 4. Khởi tạo Redis
    common.InitRedis("localhost:6379", "", 0)
    
    // 5. Khởi tạo các dependency
    profileRepo := repository.NewProfileRepository(config.DB)
    serviceProfile := business.NewProfileService(profileRepo)
    profileHandler := handlers.NewProfileHandler(serviceProfile)

    followRepo := repository.NewFollowRepository(config.DB)
    followService := business.NewFollowService(profileRepo, followRepo)
    followHandler := handlers.NewFollowHandler(followService)
    // Khởi tạo User repository và service
    repoUser := repository.NewUserRepository(config.DB)
    serviceUser := business.NewUserService(repoUser, profileRepo)
    handler := handlers.NewUserHandler(serviceUser)
    
    // 6. Tạo auth service TRƯỚC KHI sử dụng trong goroutine
    authService := service.NewAuthService(
        repoUser,
        os.Getenv("JWT_SECRET_KEY"),
        24*time.Hour,
    )
    
    // 7. Tạo channels cho các server
    httpErr := make(chan error)
    grpcErr := make(chan error)
    
    // 8. Chạy HTTP server trong goroutine
    go func() {
        m := Main{}
        log.Info("Starting HTTP server on ", config.Config.Port)
        httpErr <- m.initServe(repoUser, handler, profileHandler, followHandler)
    }()
    
    // 9. Chạy gRPC server trong goroutine
    go func(authService pb.AuthServiceServer) { // Truyền authService như parameter
        lis, err := net.Listen("tcp", ":50051")
        if err != nil {
            grpcErr <- err
            return
        }

        grpcServer := grpc.NewServer()
        pb.RegisterAuthServiceServer(grpcServer, authService)
        log.Printf("gRPC auth server listening at %v", lis.Addr())
        grpcErr <- grpcServer.Serve(lis)
    }(authService) // Truyền authService vào goroutine
    
    // 10. Chờ lỗi từ một trong hai server
    select {
    case err := <-httpErr:
        log.Fatalf("HTTP server failed: %v", err)
    case err := <-grpcErr:
        log.Fatalf("gRPC server failed: %v", err)
    }
}
