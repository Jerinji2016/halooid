package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/auth"
	"github.com/Jerinji2016/halooid/backend/internal/gateway"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "API-GATEWAY: ", log.LstdFlags)

	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/api-gateway.yaml"
	}

	config, err := gateway.LoadConfig(configPath)
	if err != nil {
		logger.Printf("Failed to load configuration from %s: %v", configPath, err)
		logger.Println("Using default configuration")
		config = gateway.DefaultConfig()
	}

	// Get database connection parameters from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "halooid")
	dbPassword := getEnv("DB_PASSWORD", "halooid_password")
	dbName := getEnv("DB_NAME", "halooid")

	// Connect to PostgreSQL
	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get Redis connection parameters from environment variables
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})
	defer redisClient.Close()

	// Ping Redis to check connection
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)

	// Get auth configuration from environment variables
	accessTokenSecret := getEnv("ACCESS_TOKEN_SECRET", "your-access-token-secret")
	refreshTokenSecret := getEnv("REFRESH_TOKEN_SECRET", "your-refresh-token-secret")
	accessTokenExpiryStr := getEnv("ACCESS_TOKEN_EXPIRY", "15m")
	refreshTokenExpiryStr := getEnv("REFRESH_TOKEN_EXPIRY", "168h")
	issuer := getEnv("ISSUER", "halooid-auth-service")

	// Parse token expiry durations
	accessTokenExpiry, err := time.ParseDuration(accessTokenExpiryStr)
	if err != nil {
		logger.Printf("Invalid access token expiry: %v, using default", err)
		accessTokenExpiry = 15 * time.Minute
	}

	refreshTokenExpiry, err := time.ParseDuration(refreshTokenExpiryStr)
	if err != nil {
		logger.Printf("Invalid refresh token expiry: %v, using default", err)
		refreshTokenExpiry = 7 * 24 * time.Hour
	}

	// Initialize services
	authService := auth.NewService(userRepo, redisClient, auth.Config{
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
		Issuer:             issuer,
	})

	// Initialize API Gateway
	apiGateway := gateway.NewService(config, authService, logger)
	apiGateway.Setup()

	// Start API Gateway in a goroutine
	go func() {
		if err := apiGateway.Run(); err != nil {
			logger.Fatalf("Failed to start API Gateway: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Shutdown API Gateway
	if err := apiGateway.Shutdown(ctx); err != nil {
		logger.Fatalf("API Gateway forced to shutdown: %v", err)
	}

	logger.Println("API Gateway exited properly")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
