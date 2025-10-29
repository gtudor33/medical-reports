package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	config "github.com/tudormiron/medical-reports/internal/configs"
	"github.com/tudormiron/medical-reports/internal/repository/postgres"
	"github.com/tudormiron/medical-reports/internal/services"
	"github.com/tudormiron/medical-reports/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	// Initialize repositories
	reportRepo := postgres.NewReportRepository(db)
	referenceRepo := postgres.NewReferenceRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// Initialize services
	reportService := services.NewReportService(reportRepo)
	referenceService := services.NewReferenceService(referenceRepo)

	// JWT secret (should be in config/env var in production)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "medical-reports-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	authService := services.NewAuthService(userRepo, jwtSecret)

	// Initialize server
	srv := server.NewServer(cfg, reportService, referenceService, authService)

	// Start server
	log.Printf("Medical Reports API starting...")
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
