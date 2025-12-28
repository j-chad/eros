package config

import (
	"time"
)

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Auth:     loadAuthConfig(),
	}

	return cfg, nil
}

// loadServerConfig loads server configuration
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Host:         getEnv("SERVER_HOST", "127.0.0.1"),
		Port:         getEnvInt("SERVER_PORT", 8080),
		ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		CORS_Origins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", ",", []string{}),
	}
}

// loadDatabaseConfig loads database configuration
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Path:        getEnv("DATABASE_PATH", "./db.sqlite"),
		BusyTimeout: getEnvDuration("DATABASE_BUSY_TIMEOUT", 5*time.Second),
		WAL:         getEnvBool("DATABASE_WAL", true),
	}
}

// loadAuthConfig loads authentication configuration
func loadAuthConfig() AuthConfig {
	return AuthConfig{
		AdminAPIKey: requireEnv("ADMIN_API_KEY"),
	}
}
