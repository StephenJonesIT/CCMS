package main

import (
	"context"

	"github.com/StephenJonesIT/CCMS/src/post-service/business"
	"github.com/StephenJonesIT/CCMS/src/post-service/config"
	"github.com/StephenJonesIT/CCMS/src/post-service/handlers"
	"github.com/StephenJonesIT/CCMS/src/post-service/repositories"
	"github.com/StephenJonesIT/CCMS/src/post-service/routes"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	_ "github.com/StephenJonesIT/CCMS/src/post-service/docs"
)

// @title Posts Service API Document
// @version 1.0
// @description List APIs of Posts Service
// @termsOfService http://swagger.io/terms/

// @host 127.0.0.1:9100
// @BasePath /api/v1
// @schemes http https
func main(){

	log.SetFormatter(&log.JSONFormatter{})
	if err := godotenv.Load(".env"); err != nil {
        log.Warn("No .env file found, using environment variables")
    }

	if err := config.LoadConfig("config/config.json"); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

	db := config.ConnectDB()
	defer db.Disconnect(context.Background())
	postRepo := repositories.NewPostRepositoryImpl(db)
	postService := business.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	m := routes.NewMain()
	if err := m.InitServe(postHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
		log.Info("Server started successfully")
}