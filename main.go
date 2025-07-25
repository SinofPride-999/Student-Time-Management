package main

import (
	"log"
	"os"

	"github.com/SinofPride-999/student-time-management/config"
	"github.com/SinofPride-999/student-time-management/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	// Initialize database
	config.ConnectDB()
	config.AutoMigrate()
	
	// Setup router
	r := routes.SetupRouter(config.DB)
	
	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}

	// Start server
	log.Fatal(r.Run(":" + port))
}
