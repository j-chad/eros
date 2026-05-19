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

func mergeFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %s: %w", path, err)
	}

	if err := mergeFromJSON(cfg, data); err != nil {
		return fmt.Errorf("merge config from %s: %w", path, err)
	}
	return nil
}

func mergeFromJSON(dst *Config, data []byte) error {
	// Decode into generic map so strings stay as strings
	// (avoids encoding/json trying to put "10s" into time.Duration)
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshal config JSON: %w", err)
	}
	return mergeMap(reflect.ValueOf(dst).Elem(), raw)
}

// mergeMap recursively applies values from a map[string]any onto a struct,
// skipping zero/nil values. It handles time.Duration from string values.
func mergeMap(dst reflect.Value, src map[string]any) error {
	dstType := dst.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dst.Field(i)
		ft := dstType.Field(i)

		jsonTag := ft.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		key := strings.Split(jsonTag, ",")[0]

		val, ok := src[key]
		if !ok || val == nil {
			continue
		}

		if err := setField(field, val); err != nil {
			return fmt.Errorf("field %s: %w", ft.Name, err)
		}
	}
	return nil
}

func setField(field reflect.Value, val any) error {
	// Nested struct - recurse with sub-map
	if field.Kind() == reflect.Struct && field.Type() != reflect.TypeOf(time.Time{}) {
		sub, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object, got %T", val)
		}
		return mergeMap(field, sub)
	}

	// time.Duration from string (e.g. "10s", "5m")
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		s, ok := val.(string)
		if !ok {
			return fmt.Errorf("expected duration string, got %T", val)
		}
		d, err := time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("invalid duration %q: %w", s, err)
		}
		field.Set(reflect.ValueOf(d))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		s, ok := val.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", val)
		}
		field.SetString(s)

	case reflect.Int, reflect.Int64:
		// JSON numbers decode as float64
		f, ok := val.(float64)
		if !ok {
			return fmt.Errorf("expected number, got %T", val)
		}
		field.SetInt(int64(f))

	case reflect.Bool:
		b, ok := val.(bool)
		if !ok {
			return fmt.Errorf("expected bool, got %T", val)
		}
		field.SetBool(b)

	case reflect.Float64:
		f, ok := val.(float64)
		if !ok {
			return fmt.Errorf("expected number, got %T", val)
		}
		field.SetFloat(f)

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			arr, ok := val.([]any)
			if !ok {
				return fmt.Errorf("expected array, got %T", val)
			}
			strs := make([]string, 0, len(arr))
			for _, v := range arr {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("expected string in array, got %T", v)
				}
				strs = append(strs, s)
			}
			field.Set(reflect.ValueOf(strs))
		}
	case reflect.Map:
		if field.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("expected map[string], got %T", val)
		}
		m, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object, got %T", val)
		}

		mapVal := field
		if mapVal.IsNil() {
			mapVal = reflect.MakeMap(field.Type())
		}

		for k, v := range m {
			key := reflect.ValueOf(k)
			value := reflect.New(field.Type().Elem()).Elem()
			if err := setField(value, v); err != nil {
				return fmt.Errorf("map value for key %s: %w", k, err)
			}
			mapVal.SetMapIndex(key, value)
		}
		field.Set(mapVal)
	default:
		return fmt.Errorf("unsupported field kind %s", field.Kind())
	}

	return nil
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
