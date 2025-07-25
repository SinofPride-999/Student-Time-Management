package main

import (
	"log"
	"os"

	"github.com/SinofPride-999/student-time-management/config"
	"github.com/SinofPride-999/student-time-management/routes"
	// "github.com/joho/godotenv"
)

func main() {
	// Initialize database (now using SQLite)
	config.ConnectDB()
	config.AutoMigrate()

	// Setup router
	r := routes.SetupRouter(config.DB)

	// Use environment variable or default to port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running at http://localhost:" + port)
	log.Fatal(r.Run(":" + port))
}
