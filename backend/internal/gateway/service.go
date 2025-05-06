package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Service represents the API Gateway service
type Service struct {
	config      *Config
	router      *chi.Mux
	authService auth.Service
	logger      *log.Logger
}

// NewService creates a new API Gateway service
func NewService(config *Config, authService auth.Service, logger *log.Logger) *Service {
	return &Service{
		config:      config,
		router:      chi.NewRouter(),
		authService: authService,
		logger:      logger,
	}
}

// Setup sets up the API Gateway service
func (s *Service) Setup() {
	// Set up middleware
	s.setupMiddleware()

	// Set up routes
	s.setupRoutes()
}

// setupMiddleware sets up middleware for the API Gateway
func (s *Service) setupMiddleware() {
	// Basic middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	if s.config.CORS.Enabled {
		s.router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   s.config.CORS.AllowedOrigins,
			AllowedMethods:   s.config.CORS.AllowedMethods,
			AllowedHeaders:   s.config.CORS.AllowedHeaders,
			ExposedHeaders:   s.config.CORS.ExposedHeaders,
			AllowCredentials: s.config.CORS.AllowCredentials,
			MaxAge:           s.config.CORS.MaxAge,
		}))
	}

	// Rate limiting middleware
	if s.config.RateLimit.Enabled {
		s.router.Use(httprate.LimitByIP(
			s.config.RateLimit.Requests,
			s.config.RateLimit.Period,
		))
	}
}

// setupRoutes sets up routes for the API Gateway
func (s *Service) setupRoutes() {
	// Health check endpoint
	if s.config.HealthCheck.Enabled {
		s.router.Get(s.config.HealthCheck.Path, s.handleHealthCheck)
	}

	// Metrics endpoint
	if s.config.Metrics.Enabled {
		s.router.Handle(s.config.Metrics.Path, promhttp.Handler())
	}

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		// Auth service routes
		r.Route("/auth", func(r chi.Router) {
			s.handleProxy(r, s.config.Services.Auth.URL)
		})

		// RBAC service routes
		r.Route("/rbac", func(r chi.Router) {
			s.handleProxy(r, s.config.Services.RBAC.URL)
		})

		// Taskake service routes
		r.Route("/taskake", func(r chi.Router) {
			s.handleProxy(r, s.config.Services.Taskake.URL)
		})

		// Qultrix service routes
		r.Route("/qultrix", func(r chi.Router) {
			s.handleProxy(r, s.config.Services.Qultrix.URL)
		})

		// AdminHub service routes
		r.Route("/adminhub", func(r chi.Router) {
			s.handleProxy(r, s.config.Services.AdminHub.URL)
		})

		// CustomerConnect service routes
		r.Route("/customerconnect", func(r chi.Router) {
			s.handleProxy(r, s.config.Services.CustomerConnect.URL)
		})

		// Invantray service routes
		r.Route("/invantray", func(r chi.Router) {
			s.handleProxy(r, s.config.Services.Invantray.URL)
		})
	})
}

// handleHealthCheck handles health check requests
func (s *Service) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// handleProxy handles proxying requests to a service
func (s *Service) handleProxy(r chi.Router, targetURL string) {
	// Parse target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		s.logger.Printf("Error parsing target URL %s: %v", targetURL, err)
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Set up proxy director
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Proxy", "Halooid-API-Gateway")
		req.Host = target.Host
	}

	// Set up proxy error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		s.logger.Printf("Proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"Bad Gateway"}`))
	}

	// Handle all methods
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		// Remove the prefix from the path
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) > 2 {
			r.URL.Path = "/" + strings.Join(parts[3:], "/")
		}

		// Proxy the request
		proxy.ServeHTTP(w, r)
	})
}

// Run runs the API Gateway service
func (s *Service) Run() error {
	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	// Log server start
	s.logger.Printf("Starting API Gateway on port %d", s.config.Server.Port)

	// Start server
	return server.ListenAndServe()
}

// Shutdown shuts down the API Gateway service
func (s *Service) Shutdown(ctx context.Context) error {
	// Log server shutdown
	s.logger.Println("Shutting down API Gateway")

	// Nothing to do here yet
	return nil
}
