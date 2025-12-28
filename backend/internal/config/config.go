package config

import (
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Admin    AdminConfig
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig contains database settings
type DatabaseConfig struct {
	Path        string
	BusyTimeout time.Duration
	WAL         bool
}

// AdminConfig contains admin interface settings
type AdminConfig struct {
	APIKey string
}
