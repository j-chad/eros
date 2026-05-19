package config

import (
	"backend/internal/testutil"
	"reflect"
	"testing"
	"time"
)

func TestMergeFromJSON_String(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{"auth": {"admin_api_key": "secret123"}}`)))
	testutil.Equal(t, cfg.Auth.AdminAPIKey, "secret123")
}

func TestMergeFromJSON_Int(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{"server": {"port": 9090}}`)))
	testutil.Equal(t, cfg.Server.Port, 9090)
}

func TestMergeFromJSON_Bool(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{"database": {"wal": true}}`)))
	testutil.Equal(t, cfg.Database.WAL, true)
}

func TestMergeFromJSON_Duration(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{"server": {"read_timeout": "30s"}}`)))
	testutil.Equal(t, cfg.Server.ReadTimeout, 30*time.Second)
}

func TestMergeFromJSON_StringSlice(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{"server": {"cors": {"allowed_origins": ["http://a.com", "http://b.com"]}}}`)))
	testutil.Equal(t, len(cfg.Server.CORS.AllowedOrigins), 2)
	testutil.Equal(t, cfg.Server.CORS.AllowedOrigins[0], "http://a.com")
}

func TestMergeFromJSON_NestedStruct(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{"logging": {"collector": {"enabled": true, "max_spans": 100}}}`)))
	testutil.Equal(t, cfg.Logging.Collector.Enabled, true)
	testutil.Equal(t, cfg.Logging.Collector.MaxSpans, 100)
}

func TestMergeFromJSON_InvalidJSON(t *testing.T) {
	cfg := &Config{}
	testutil.NotNilErr(t, mergeFromJSON(cfg, []byte(`not json`)))
}

func TestMergeFromJSON_SkipsNullValues(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Port: 8080}}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{"server": {"port": null}}`)))
	testutil.Equal(t, cfg.Server.Port, 8080)
}

func TestMergeFromJSON_PreservesExisting(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Host: "original", Port: 8080}}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{"server": {"port": 9090}}`)))
	testutil.Equal(t, cfg.Server.Host, "original")
	testutil.Equal(t, cfg.Server.Port, 9090)
}

func TestSetFieldFromString_String(t *testing.T) {
	var s string
	v := reflect.ValueOf(&s).Elem()
	testutil.NilErr(t, setFieldFromString(v, "hello"))
	testutil.Equal(t, s, "hello")
}

func TestSetFieldFromString_Int(t *testing.T) {
	var n int
	v := reflect.ValueOf(&n).Elem()
	testutil.NilErr(t, setFieldFromString(v, "42"))
	testutil.Equal(t, n, 42)
}

func TestSetFieldFromString_Bool(t *testing.T) {
	var b bool
	v := reflect.ValueOf(&b).Elem()
	testutil.NilErr(t, setFieldFromString(v, "true"))
	testutil.Equal(t, b, true)
}

func TestSetFieldFromString_Float(t *testing.T) {
	var f float64
	v := reflect.ValueOf(&f).Elem()
	testutil.NilErr(t, setFieldFromString(v, "3.14"))
	testutil.Equal(t, f, 3.14)
}

func TestSetFieldFromString_Duration(t *testing.T) {
	var d time.Duration
	v := reflect.ValueOf(&d).Elem()
	testutil.NilErr(t, setFieldFromString(v, "5m"))
	testutil.Equal(t, d, 5*time.Minute)
}

func TestSetFieldFromString_CommaSeparatedSlice(t *testing.T) {
	var s []string
	v := reflect.ValueOf(&s).Elem()
	testutil.NilErr(t, setFieldFromString(v, "a, b, c"))
	testutil.Equal(t, len(s), 3)
	testutil.Equal(t, s[0], "a")
	testutil.Equal(t, s[1], "b")
	testutil.Equal(t, s[2], "c")
}

func TestSetFieldFromString_EmptyCommaParts(t *testing.T) {
	var s []string
	v := reflect.ValueOf(&s).Elem()
	testutil.NilErr(t, setFieldFromString(v, "a,,b, ,c"))
	testutil.Equal(t, len(s), 3)
}

func TestSetFieldFromString_InvalidInt_Silently(t *testing.T) {
	var n int
	v := reflect.ValueOf(&n).Elem()
	// setFieldFromString silently ignores parse errors for int
	testutil.NilErr(t, setFieldFromString(v, "not-a-number"))
	testutil.Equal(t, n, 0)
}

// --- Map merge (SchedulerConfig.Tasks) ---

func TestMergeFromJSON_MapSingleEntry(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {
				"push_notifications": {
					"cron": "* * * * *",
					"timeout": "30s"
				}
			}
		}
	}`)))
	testutil.Equal(t, len(cfg.Scheduler.Tasks), 1)
	task := cfg.Scheduler.Tasks["push_notifications"]
	testutil.Equal(t, task.Cron, "* * * * *")
	testutil.Equal(t, task.Timeout, 30*time.Second)
	testutil.Equal(t, task.Disabled, false)
}

func TestMergeFromJSON_MapMultipleEntries(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {
				"push_notifications": {"cron": "* * * * *", "timeout": "30s"},
				"file_cleanup": {"cron": "0 3 * * *", "timeout": "5m"}
			}
		}
	}`)))
	testutil.Equal(t, len(cfg.Scheduler.Tasks), 2)
	testutil.Equal(t, cfg.Scheduler.Tasks["push_notifications"].Cron, "* * * * *")
	testutil.Equal(t, cfg.Scheduler.Tasks["file_cleanup"].Cron, "0 3 * * *")
	testutil.Equal(t, cfg.Scheduler.Tasks["file_cleanup"].Timeout, 5*time.Minute)
}

func TestMergeFromJSON_MapWithDisabledTask(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {
				"push_notifications": {"disabled": true, "cron": "* * * * *", "timeout": "30s"}
			}
		}
	}`)))
	testutil.Equal(t, cfg.Scheduler.Tasks["push_notifications"].Disabled, true)
}

func TestMergeFromJSON_MapMergePreservesExistingEntries(t *testing.T) {
	// First merge: default config with push_notifications
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {
				"push_notifications": {"cron": "* * * * *", "timeout": "30s"}
			}
		}
	}`)))

	// Second merge: production config adds file_cleanup
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {
				"file_cleanup": {"cron": "0 3 * * *", "timeout": "5m"}
			}
		}
	}`)))

	// Both entries should exist
	testutil.Equal(t, len(cfg.Scheduler.Tasks), 2)
	testutil.Equal(t, cfg.Scheduler.Tasks["push_notifications"].Cron, "* * * * *")
	testutil.Equal(t, cfg.Scheduler.Tasks["file_cleanup"].Cron, "0 3 * * *")
}

func TestMergeFromJSON_MapMergeOverridesExistingEntry(t *testing.T) {
	// First merge: default config
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {
				"push_notifications": {"cron": "* * * * *", "timeout": "30s"}
			}
		}
	}`)))

	// Second merge: production overrides the cron for push_notifications
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {
				"push_notifications": {"cron": "*/5 * * * *", "timeout": "1m"}
			}
		}
	}`)))

	testutil.Equal(t, len(cfg.Scheduler.Tasks), 1)
	testutil.Equal(t, cfg.Scheduler.Tasks["push_notifications"].Cron, "*/5 * * * *")
	testutil.Equal(t, cfg.Scheduler.Tasks["push_notifications"].Timeout, 1*time.Minute)
}

func TestMergeFromJSON_MapEmptyObject(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {}
		}
	}`)))
	testutil.Equal(t, len(cfg.Scheduler.Tasks), 0)
}

func TestMergeFromJSON_MapNullPreservesExisting(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": {
				"push_notifications": {"cron": "* * * * *", "timeout": "30s"}
			}
		}
	}`)))

	// Null tasks should not wipe out existing map
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {
			"tasks": null
		}
	}`)))
	testutil.Equal(t, len(cfg.Scheduler.Tasks), 1)
}

func TestMergeFromJSON_SchedulerDisabledViaJSON(t *testing.T) {
	cfg := &Config{}
	testutil.NilErr(t, mergeFromJSON(cfg, []byte(`{
		"scheduler": {"disabled": true}
	}`)))
	testutil.Equal(t, cfg.Scheduler.Disabled, true)
}

func TestValidateRequired_MissingField(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Host: "", Port: 0},
		Auth:   AuthConfig{AdminAPIKey: ""},
	}
	err := validateRequired(cfg)
	testutil.NotNilErr(t, err)
}

func TestValidateRequired_AllSet(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Host: "localhost", Port: 8080},
		Auth:   AuthConfig{AdminAPIKey: "key"},
	}
	testutil.NilErr(t, validateRequired(cfg))
}

func TestValidateRequired_PartiallySet(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Host: "localhost", Port: 8080},
		Auth:   AuthConfig{AdminAPIKey: ""},
	}
	err := validateRequired(cfg)
	testutil.NotNilErr(t, err)
}
