package main

import (
	"log"
	"os"

	"github.com/gbadegesintestimony/jwt-authentication/config"
	"github.com/gbadegesintestimony/jwt-authentication/database"
	"github.com/gbadegesintestimony/jwt-authentication/routes"
	"github.com/gbadegesintestimony/jwt-authentication/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load .env file
	config.LoadEnv()

	// Initialize database
	database.Connect()

	// Create a new gin engine
	r := gin.Default()

	// Setup routes
	routes.Setup(r)

	utils.Init()

	// Get server port from environment variable
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
