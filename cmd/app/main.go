package main

import (
	"log"
	"os"
	"payment-sandbox/internal/middleware"

	"payment-sandbox/internal/config"
	"payment-sandbox/internal/modules"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "payment-sandbox/docs"
)

// @title Payment Sandbox API
// @version 1.0
// @description Mini Payment Gateway Simulation
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("failed load .env:", err)
	}

	// sql connection
	db := config.InitDB()

	// mongoDB connection
	config.InitMongo()

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api/v1")
	modules.Init(api, db)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("server running on :", port)
	log.Println("swagger: http://localhost:" + port + "/swagger/index.html")

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
