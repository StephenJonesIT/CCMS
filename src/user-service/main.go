package main

import (
	"log"
	"net/http"
	"user-service/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Main struct {
	router *gin.Engine
}

func(m *Main) initServe(){
	m.router = gin.Default()
	m.router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message" : "Hello World",
		})
	})

	m.router.Run()
}

func main(){
	m := Main{}
	err := godotenv.Load(".env")
	
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
	config.ConfigDatabase()
	m.initServe()
}