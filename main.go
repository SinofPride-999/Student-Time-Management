package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/SinofPride-999/student-time-management/config"
	"github.com/SinofPride-999/student-time-management/routes"
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
	
	// Start server
	log.Fatal(r.Run(":8080"))
}
