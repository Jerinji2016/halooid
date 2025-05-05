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
	"github.com/Jerinji2016/halooid/backend/internal/rbac"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/Jerinji2016/halooid/backend/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "RBAC-SERVICE: ", log.LstdFlags)

	// Load configuration
	// In a real application, this would be loaded from environment variables or a config file
	authConfig := auth.Config{
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
	roleRepo := repository.NewPostgresRoleRepository(db)

	// Initialize services
	authService := auth.NewService(userRepo, redisClient, authConfig)
	rbacService := rbac.NewService(roleRepo, userRepo)

	// Initialize handlers
	authHandlers := auth.NewHandlers(authService)
	rbacHandlers := rbac.NewHandlers(rbacService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService)

	// Initialize router
	router := mux.NewRouter()

	// Register auth routes
	router.HandleFunc("/api/auth/register", authHandlers.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandlers.Login).Methods("POST")
	router.HandleFunc("/api/auth/refresh", authHandlers.RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/logout", authHandlers.Logout).Methods("POST")
	router.HandleFunc("/api/auth/me", authHandlers.Me).Methods("GET")

	// Register RBAC routes
	// Role routes
	roleRouter := router.PathPrefix("/api/rbac/roles").Subrouter()
	roleRouter.Use(authMiddleware.Authenticate)
	roleRouter.HandleFunc("", rbacHandlers.CreateRole).Methods("POST")
	roleRouter.HandleFunc("", rbacHandlers.ListRoles).Methods("GET")
	roleRouter.HandleFunc("/{id}", rbacHandlers.GetRole).Methods("GET")
	roleRouter.HandleFunc("/{id}", rbacHandlers.UpdateRole).Methods("PUT")
	roleRouter.HandleFunc("/{id}", rbacHandlers.DeleteRole).Methods("DELETE")

	// Permission routes
	permissionRouter := router.PathPrefix("/api/rbac/permissions").Subrouter()
	permissionRouter.Use(authMiddleware.Authenticate)
	permissionRouter.HandleFunc("", rbacHandlers.CreatePermission).Methods("POST")
	permissionRouter.HandleFunc("", rbacHandlers.ListPermissions).Methods("GET")
	permissionRouter.HandleFunc("/{id}", rbacHandlers.GetPermission).Methods("GET")
	permissionRouter.HandleFunc("/{id}", rbacHandlers.UpdatePermission).Methods("PUT")
	permissionRouter.HandleFunc("/{id}", rbacHandlers.DeletePermission).Methods("DELETE")

	// User role routes
	userRoleRouter := router.PathPrefix("/api/rbac/user-roles").Subrouter()
	userRoleRouter.Use(authMiddleware.Authenticate)
	userRoleRouter.HandleFunc("", rbacHandlers.AssignRoleToUser).Methods("POST")
	userRoleRouter.HandleFunc("/{user_id}/{role_id}/{org_id}", rbacHandlers.RemoveRoleFromUser).Methods("DELETE")
	userRoleRouter.HandleFunc("/{user_id}/{org_id}", rbacHandlers.GetUserRoles).Methods("GET")
	userRoleRouter.HandleFunc("/{user_id}/{org_id}/{permission}", rbacHandlers.CheckPermission).Methods("GET")

	// Protected routes example
	protectedRouter := router.PathPrefix("/api/protected").Subrouter()
	protectedRouter.Use(authMiddleware.Authenticate)
	
	// Route that requires a specific permission
	adminRouter := protectedRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(rbacMiddleware.RequirePermission("admin:access"))
	adminRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Admin access granted"}`))
	}).Methods("GET")
	
	// Route that requires any of the specified permissions
	userRouter := protectedRouter.PathPrefix("/user").Subrouter()
	userRouter.Use(rbacMiddleware.RequireAnyPermission("user:read", "user:write"))
	userRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "User access granted"}`))
	}).Methods("GET")
	
	// Route that requires all of the specified permissions
	superUserRouter := protectedRouter.PathPrefix("/super-user").Subrouter()
	superUserRouter.Use(rbacMiddleware.RequireAllPermissions("user:read", "user:write", "user:delete"))
	superUserRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Super user access granted"}`))
	}).Methods("GET")

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
