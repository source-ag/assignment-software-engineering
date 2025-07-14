package main

import (
	"log"
	"os"

	"meteo-api/internal/database"
	"meteo-api/internal/handlers"
	"meteo-api/internal/middleware"
	"meteo-api/internal/repository"
	"meteo-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	dbConfig := database.NewConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize layers
	repo := repository.NewMeteoRepository(db)
	svc := service.NewMeteoService(repo)
	auth := middleware.NewAuthMiddleware()
	handler := handlers.NewMeteoHandler(svc, auth)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register routes
	handler.RegisterRoutes(router)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}