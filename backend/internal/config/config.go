package config

import (
	"log/slog"
	"time"
)

// Config holds all application configuration
type Config struct {
	Logging     LoggingConfig     `json:"logging"`
	Server      ServerConfig      `json:"server"`
	Database    DatabaseConfig    `json:"database"`
	Auth        AuthConfig        `json:"auth"`
	FileStorage FileStorageConfig `json:"file_storage"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level     slog.Level       `json:"level" env:"LOG_LEVEL"`
	JSON      bool             `json:"json" env:"LOG_JSON"`
	AddSource bool             `json:"add_source" env:"LOG_ADD_SOURCE"`
	Collector *CollectorConfig `json:"collector"` // nil if disabled
}

// CollectorConfig contains settings for the log collector
type CollectorConfig struct {
	Enabled     bool   `json:"enabled" env:"COLLECTOR_ENABLED"`
	MaxSpans    int    `json:"max_spans" env:"LOG_COLLECTOR_MAX_SPANS"`
	TraceHeader string `json:"trace_header" env:"LOG_COLLECTOR_TRACE_HEADER"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host         string        `json:"host" env:"SERVER_HOST" required:"true"`
	Port         int           `json:"port" env:"SERVER_PORT" required:"true"`
	ReadTimeout  time.Duration `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `json:"idle_timeout" env:"SERVER_IDLE_TIMEOUT"`
	CorsOrigins  []string      `json:"cors_origins" env:"CORS_ALLOWED_ORIGINS"`
}

// DatabaseConfig contains database settings
type DatabaseConfig struct {
	Path        string        `json:"path" env:"DATABASE_PATH"`
	BusyTimeout time.Duration `json:"busy_timeout" env:"DATABASE_BUSY_TIMEOUT"`
	WAL         bool          `json:"wal" env:"DATABASE_WAL"`
}

// AuthConfig contains auth settings
type AuthConfig struct {
	AdminAPIKey string `json:"admin_api_key" env:"ADMIN_API_KEY" required:"true"`
}

type FileStorageType string

const (
	FileStorageLocal FileStorageType = "local"
	FileStorageS3    FileStorageType = "s3"
)

type FileStorageConfig struct {
	Type  FileStorageType        `json:"type" env:"FILE_STORAGE_TYPE"`
	Local LocalFileStorageConfig `json:"local"`
	S3    S3FileStorageConfig    `json:"s3"`
}

type LocalFileStorageConfig struct {
	BasePath string `json:"base_path" env:"FILE_STORAGE_LOCAL_BASE_PATH"`
}

type S3FileStorageConfig struct {
	Region    string `json:"region" env:"FILE_STORAGE_S3_REGION"`
	Bucket    string `json:"bucket" env:"FILE_STORAGE_S3_BUCKET"`
	Endpoint  string `json:"endpoint" env:"FILE_STORAGE_S3_ENDPOINT"`
	AccessKey string `json:"access_key" env:"FILE_STORAGE_S3_ACCESS_KEY"`
	SecretKey string `json:"secret_key" env:"FILE_STORAGE_S3_SECRET_KEY"`
}
