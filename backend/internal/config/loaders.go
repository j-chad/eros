package config

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

//go:embed config.default.json config.develop.json config.production.json
var configFiles embed.FS

const defaultConfigFile = "config.default.json"
const privateConfigFile = "config.private.json"

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	// Layer 0: embedded default config
	if err := mergeFromEmbed(cfg, defaultConfigFile, true); err != nil {
		return nil, fmt.Errorf("load default config: %w", err)
	}

	// Layer 1: embedded environment-specific config (e.g. config.production.json)
	env := getEnv("APP_ENV", "production")
	envConfigFile := fmt.Sprintf("config.%s.json", env)
	if err := mergeFromEmbed(cfg, envConfigFile, false); err != nil {
		return nil, fmt.Errorf("load environment config: %w", err)
	}

	// Layer 2: private config loaded from disk (never embedded in binary)
	privatePath := getEnv("CONFIG_PRIVATE_PATH", privateConfigFile)
	if err := mergeFromFile(cfg, privatePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("load private config: %w", err)
		}

		slog.Warn("private config file not found, skipping", "path", privatePath)
	}

	// Layer 3: environment variables
	if err := applyEnvVars(cfg); err != nil {
		return nil, fmt.Errorf("apply environment variables: %w", err)
	}

	// Validate required fields
	if err := validateRequired(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}
