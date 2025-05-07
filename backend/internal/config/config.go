package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
	Storage  StorageConfig  `yaml:"storage"`
	API      APIConfig      `yaml:"api"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port         int `yaml:"port"`
	Timeout      int `yaml:"timeout"`
	ReadTimeout  int `yaml:"read_timeout"`
	WriteTimeout int `yaml:"write_timeout"`
	IdleTimeout  int `yaml:"idle_timeout"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	URL             string `yaml:"url"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// AuthConfig represents the authentication configuration
type AuthConfig struct {
	JWTSecret          string `yaml:"jwt_secret"`
	TokenExpiry        int    `yaml:"token_expiry"`
	RefreshTokenExpiry int    `yaml:"refresh_token_expiry"`
}

// StorageConfig represents the file storage configuration
type StorageConfig struct {
	BasePath string `yaml:"base_path"`
}

// APIConfig represents the API configuration
type APIConfig struct {
	BaseURL string `yaml:"base_url"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	// Read the configuration file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the configuration
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}
