package gateway

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads the configuration from a file
func LoadConfig(filePath string) (*Config, error) {
	// Load default configuration
	config := DefaultConfig()

	// Read configuration file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse configuration file
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Override configuration with environment variables
	overrideWithEnv(config)

	return config, nil
}

// overrideWithEnv overrides configuration with environment variables
func overrideWithEnv(config *Config) {
	// Server configuration
	if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil && port > 0 {
		config.Server.Port = port
	}

	// Services configuration
	if url := os.Getenv("AUTH_SERVICE_URL"); url != "" {
		config.Services.Auth.URL = url
	}
	if url := os.Getenv("RBAC_SERVICE_URL"); url != "" {
		config.Services.RBAC.URL = url
	}
	if url := os.Getenv("TASKAKE_SERVICE_URL"); url != "" {
		config.Services.Taskake.URL = url
	}
	if url := os.Getenv("QULTRIX_SERVICE_URL"); url != "" {
		config.Services.Qultrix.URL = url
	}
	if url := os.Getenv("ADMINHUB_SERVICE_URL"); url != "" {
		config.Services.AdminHub.URL = url
	}
	if url := os.Getenv("CUSTOMERCONNECT_SERVICE_URL"); url != "" {
		config.Services.CustomerConnect.URL = url
	}
	if url := os.Getenv("INVANTRAY_SERVICE_URL"); url != "" {
		config.Services.Invantray.URL = url
	}

	// Timeouts
	if timeout, err := time.ParseDuration(os.Getenv("SERVICE_TIMEOUT")); err == nil {
		config.Services.Auth.Timeout = timeout
		config.Services.RBAC.Timeout = timeout
		config.Services.Taskake.Timeout = timeout
		config.Services.Qultrix.Timeout = timeout
		config.Services.AdminHub.Timeout = timeout
		config.Services.CustomerConnect.Timeout = timeout
		config.Services.Invantray.Timeout = timeout
	}

	// Logging configuration
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Logging.Format = format
	}
}
