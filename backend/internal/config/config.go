package config

import (
	"time"
)

// Config holds all application configuration
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Auth        AuthConfig
	FileStorage FileStorageConfig
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	CORS_Origins []string
}

// DatabaseConfig contains database settings
type DatabaseConfig struct {
	Path        string
	BusyTimeout time.Duration
	WAL         bool
}

// AuthConfig contains auth settings
type AuthConfig struct {
	AdminAPIKey string
}

type FileStorageType string

const (
	FileStorageLocal FileStorageType = "local"
	FileStorageS3    FileStorageType = "s3"
)

type FileStorageConfig struct {
	Type  FileStorageType
	Local LocalFileStorageConfig
	S3    S3FileStorageConfig
}

type LocalFileStorageConfig struct {
	BasePath string
}

type S3FileStorageConfig struct {
	Region    string
	Bucket    string
	Endpoint  string
	AccessKey string
	SecretKey string
}
