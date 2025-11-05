package main

import (
	"log"
	"os"

	"github.com/gbadegesintestimony/jwt-authentication/database"
	"github.com/gbadegesintestimony/jwt-authentication/models"
	"github.com/gbadegesintestimony/jwt-authentication/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	database.Connect()

	// Auto migrate the database
	database.DB.AutoMigrate(&models.User{})

	// Create a new gin engine
	r := gin.Default()

	// Setup routes
	routes.Setup(r)

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

// func main() {
// 	// load .env early if needed
// 	if _, err := os.Stat(".env"); err == nil {
// 		// godotenv loaded inside database.Connect; optional here
// 	}

// 	database.Connect()

// 	// Migrate models
// 	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
// 		log.Fatal("AutoMigrate error:", err)
// 	}

// 	r := gin.Default()
// 	routes.Setup(r)

// 	r.Run(":8080")
// }
