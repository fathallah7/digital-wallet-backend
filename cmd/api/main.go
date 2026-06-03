package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/fathallah7/wallet-service/internal/config"
	"github.com/fathallah7/wallet-service/internal/database"
	"github.com/fathallah7/wallet-service/internal/handler"
	"github.com/fathallah7/wallet-service/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	cfg := config.Load()

	db, err := database.New(cfg.DBUrl)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()
	log.Println("Successfully connected to the database")

	h := handler.New(db, []byte(cfg.JWTSecret))
	r := router.Setup(h, []byte(cfg.JWTSecret))

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
