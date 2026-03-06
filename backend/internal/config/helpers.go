package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func mergeFromEmbed(cfg *Config, configFile string, required bool) error {
	data, err := configFiles.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && !required {
			return nil
		}
		return err
	}

	err = mergeFromJSON(cfg, data)
	if err != nil {
		return fmt.Errorf("merge config from %s: %w", configFile, err)
	}
	return nil
}

func mergeFromJSON(dst *Config, data []byte) error {
	overlay := &Config{}
	if err := json.Unmarshal(data, overlay); err != nil {
		return fmt.Errorf("unmarshal config JSON: %w", err)
	}
	mergeStructs(reflect.ValueOf(dst).Elem(), reflect.ValueOf(overlay).Elem())
	return nil
}

// mergeStructs recursively copies non-zero fields from src into dst.
func mergeStructs(dst, src reflect.Value) {
	for i := 0; i < dst.NumField(); i++ {
		dstField := dst.Field(i)
		srcField := src.Field(i)

		if srcField.IsZero() {
			continue
		}

		if dstField.Kind() == reflect.Struct && dstField.Type() != reflect.TypeOf(time.Duration(0)) {
			mergeStructs(dstField, srcField)
		} else {
			dstField.Set(srcField)
		}
	}
}

// applyEnvVars walks the config struct tree and applies env var overrides
// based on the `env` struct tag.
func applyEnvVars(cfg *Config) error {
	return applyEnvToStruct(reflect.ValueOf(cfg).Elem())
}

func applyEnvToStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		ft := t.Field(i)

		// Recurse into nested structs (but not time.Duration which is int64)
		if field.Kind() == reflect.Struct && ft.Type != reflect.TypeOf(time.Duration(0)) {
			if err := applyEnvToStruct(field); err != nil {
				return err
			}
			continue
		}

		envKey := ft.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envVal, ok := os.LookupEnv(envKey)
		if !ok {
			continue
		}

		err := setFieldFromString(field, envVal)
		if err != nil {
			return fmt.Errorf("set field %s from env var %s: %w", ft.Name, envKey, err)
		}
	}

	return nil
}

func setFieldFromString(field reflect.Value, val string) error {
	// Handle time.Duration
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		if d, err := time.ParseDuration(val); err == nil {
			field.Set(reflect.ValueOf(d))
		}
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int64:
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			field.SetInt(n)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(val); err == nil {
			field.SetBool(b)
		}
	case reflect.Float64:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			field.SetFloat(f)
		}
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			parts := strings.Split(val, ",")
			cleaned := make([]string, 0, len(parts))
			for _, p := range parts {
				if s := strings.TrimSpace(p); s != "" {
					cleaned = append(cleaned, s)
				}
			}
			field.Set(reflect.ValueOf(cleaned))
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type())
	}

	return nil
}

// validateRequired walks the struct and returns an error if any field
// tagged `required:"true"` is still at its zero value.
func validateRequired(cfg *Config) error {
	return walkRequired(reflect.ValueOf(cfg).Elem(), reflect.TypeOf(*cfg), "")
}

func walkRequired(v reflect.Value, t reflect.Type, prefix string) error {
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		ft := t.Field(i)

		path := prefix + ft.Name
		if prefix != "" {
			path = prefix + "." + ft.Name
		}

		if field.Kind() == reflect.Struct && ft.Type != reflect.TypeOf(time.Duration(0)) {
			if err := walkRequired(field, ft.Type, path); err != nil {
				return err
			}
			continue
		}

		if ft.Tag.Get("required") == "true" && field.IsZero() {
			envKey := ft.Tag.Get("env")
			hint := ""
			if envKey != "" {
				hint = fmt.Sprintf(" (set via env %s or config file)", envKey)
			}
			return fmt.Errorf("config: %s is required but not set%s", path, hint)
		}
	}
	return nil
}
