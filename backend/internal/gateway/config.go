package gateway

import (
	"time"
)

// Config represents the configuration for the API Gateway
type Config struct {
	// Server configuration
	Server struct {
		Port         int           `json:"port" yaml:"port"`
		ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
		WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
		IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	} `json:"server" yaml:"server"`

	// Rate limiting configuration
	RateLimit struct {
		Enabled   bool          `json:"enabled" yaml:"enabled"`
		Requests  int           `json:"requests" yaml:"requests"`
		Period    time.Duration `json:"period" yaml:"period"`
		BurstSize int           `json:"burst_size" yaml:"burst_size"`
	} `json:"rate_limit" yaml:"rate_limit"`

	// CORS configuration
	CORS struct {
		Enabled          bool     `json:"enabled" yaml:"enabled"`
		AllowedOrigins   []string `json:"allowed_origins" yaml:"allowed_origins"`
		AllowedMethods   []string `json:"allowed_methods" yaml:"allowed_methods"`
		AllowedHeaders   []string `json:"allowed_headers" yaml:"allowed_headers"`
		ExposedHeaders   []string `json:"exposed_headers" yaml:"exposed_headers"`
		AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials"`
		MaxAge           int      `json:"max_age" yaml:"max_age"`
	} `json:"cors" yaml:"cors"`

	// Services configuration
	Services struct {
		Auth struct {
			URL     string        `json:"url" yaml:"url"`
			Timeout time.Duration `json:"timeout" yaml:"timeout"`
		} `json:"auth" yaml:"auth"`
		RBAC struct {
			URL     string        `json:"url" yaml:"url"`
			Timeout time.Duration `json:"timeout" yaml:"timeout"`
		} `json:"rbac" yaml:"rbac"`
		Taskake struct {
			URL     string        `json:"url" yaml:"url"`
			Timeout time.Duration `json:"timeout" yaml:"timeout"`
		} `json:"taskake" yaml:"taskake"`
		Qultrix struct {
			URL     string        `json:"url" yaml:"url"`
			Timeout time.Duration `json:"timeout" yaml:"timeout"`
		} `json:"qultrix" yaml:"qultrix"`
		AdminHub struct {
			URL     string        `json:"url" yaml:"url"`
			Timeout time.Duration `json:"timeout" yaml:"timeout"`
		} `json:"adminhub" yaml:"adminhub"`
		CustomerConnect struct {
			URL     string        `json:"url" yaml:"url"`
			Timeout time.Duration `json:"timeout" yaml:"timeout"`
		} `json:"customerconnect" yaml:"customerconnect"`
		Invantray struct {
			URL     string        `json:"url" yaml:"url"`
			Timeout time.Duration `json:"timeout" yaml:"timeout"`
		} `json:"invantray" yaml:"invantray"`
	} `json:"services" yaml:"services"`

	// Logging configuration
	Logging struct {
		Level  string `json:"level" yaml:"level"`
		Format string `json:"format" yaml:"format"`
	} `json:"logging" yaml:"logging"`

	// Metrics configuration
	Metrics struct {
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Path    string `json:"path" yaml:"path"`
	} `json:"metrics" yaml:"metrics"`

	// Health check configuration
	HealthCheck struct {
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Path    string `json:"path" yaml:"path"`
	} `json:"health_check" yaml:"health_check"`
}

// DefaultConfig returns a default configuration for the API Gateway
func DefaultConfig() *Config {
	config := &Config{}

	// Server configuration
	config.Server.Port = 8000
	config.Server.ReadTimeout = 15 * time.Second
	config.Server.WriteTimeout = 15 * time.Second
	config.Server.IdleTimeout = 60 * time.Second

	// Rate limiting configuration
	config.RateLimit.Enabled = true
	config.RateLimit.Requests = 100
	config.RateLimit.Period = time.Minute
	config.RateLimit.BurstSize = 20

	// CORS configuration
	config.CORS.Enabled = true
	config.CORS.AllowedOrigins = []string{"*"}
	config.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.CORS.AllowedHeaders = []string{"Content-Type", "Authorization", "X-Requested-With"}
	config.CORS.ExposedHeaders = []string{"Content-Length"}
	config.CORS.AllowCredentials = true
	config.CORS.MaxAge = 86400

	// Services configuration
	config.Services.Auth.URL = "http://localhost:8001"
	config.Services.Auth.Timeout = 5 * time.Second
	config.Services.RBAC.URL = "http://localhost:8002"
	config.Services.RBAC.Timeout = 5 * time.Second
	config.Services.Taskake.URL = "http://localhost:8003"
	config.Services.Taskake.Timeout = 5 * time.Second
	config.Services.Qultrix.URL = "http://localhost:8004"
	config.Services.Qultrix.Timeout = 5 * time.Second
	config.Services.AdminHub.URL = "http://localhost:8005"
	config.Services.AdminHub.Timeout = 5 * time.Second
	config.Services.CustomerConnect.URL = "http://localhost:8006"
	config.Services.CustomerConnect.Timeout = 5 * time.Second
	config.Services.Invantray.URL = "http://localhost:8007"
	config.Services.Invantray.Timeout = 5 * time.Second

	// Logging configuration
	config.Logging.Level = "info"
	config.Logging.Format = "json"

	// Metrics configuration
	config.Metrics.Enabled = true
	config.Metrics.Path = "/metrics"

	// Health check configuration
	config.HealthCheck.Enabled = true
	config.HealthCheck.Path = "/health"

	return config
}
