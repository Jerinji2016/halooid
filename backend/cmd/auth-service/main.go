package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/auth"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "AUTH-SERVICE: ", log.LstdFlags)

	// Load configuration
	// In a real application, this would be loaded from environment variables or a config file
	config := auth.Config{
		AccessTokenSecret:  "your-access-token-secret",
		RefreshTokenSecret: "your-refresh-token-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "halooid-auth-service",
	}

	// Connect to PostgreSQL
	dbConnStr := "host=localhost port=5432 user=halooid password=halooid_password dbname=halooid sslmode=disable"
	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer redisClient.Close()

	// Ping Redis to check connection
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)

	// Initialize services
	authService := auth.NewService(userRepo, redisClient, config)

	// Initialize handlers
	authHandlers := auth.NewHandlers(authService)

	// Initialize router
	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/api/auth/register", authHandlers.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandlers.Login).Methods("POST")
	router.HandleFunc("/api/auth/refresh", authHandlers.RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/logout", authHandlers.Logout).Methods("POST")
	router.HandleFunc("/api/auth/me", authHandlers.Me).Methods("GET")

	// Create server
	server := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Starting server on port 8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Shutdown server
	logger.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited properly")
}
