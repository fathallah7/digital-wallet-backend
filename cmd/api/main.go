package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/fathallah7/wallet-service/internal/config"
	"github.com/fathallah7/wallet-service/internal/database"
	"github.com/fathallah7/wallet-service/internal/handler"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Load Config
	cfg := config.Load()

	// Connect to the database
	db, err := database.New(cfg.DBUrl)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()
	log.Println("Successfully connected to the database")

	// mux router
	mux := http.NewServeMux()
	// routes
	mux.HandleFunc("GET /health", handler.HealthHandler)

	// start server
	log.Printf("Server starting on port %s", cfg.Port)

	err = http.ListenAndServe(cfg.Port, mux)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}
