package config

import (
	"embed"
	"fmt"
)

//go:embed config.*.json
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

	// Layer 2: embedded private config (config.private.json)
	if err := mergeFromEmbed(cfg, privateConfigFile, false); err != nil {
		return nil, fmt.Errorf("load private config: %w", err)
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
