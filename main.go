package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ritikchawla/url-shortner/api"
	"github.com/ritikchawla/url-shortner/config"
	"github.com/ritikchawla/url-shortner/db"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize PostgreSQL
	if err := db.InitPostgres(cfg); err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	defer db.ClosePostgres()

	// Initialize Redis
	if err := db.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer db.CloseRedis()

	// Set up Gin router
	router := gin.Default()

	// API routes
	router.POST("/api/shorten", api.CreateShortURL)
	router.GET("/api/stats/:shortCode", api.GetURLStats)

	// Redirect route (should be last to avoid conflicts)
	router.GET("/:shortCode", api.RedirectURL)

	// Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
