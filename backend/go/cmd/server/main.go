package main

import (
	"fmt"
	"log"
	"net/http"

	"indexarr/internal/api"
	"indexarr/internal/config"
	"indexarr/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := repository.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Seed database with mock data
	if err := repository.SeedMockData(db); err != nil {
		log.Fatalf("Failed to seed mock data: %v", err)
	}

	// Setup API router
	router := api.SetupRoutes(db)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🎬 Indexarr server running on http://localhost:%s", cfg.ServerPort)
	log.Printf("📁 Database: %s", cfg.DBPath)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
