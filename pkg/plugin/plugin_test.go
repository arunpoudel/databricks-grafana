package plugin

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseConnectionSettings(t *testing.T) {
	settings := json.RawMessage(`{
		"maxOpenConns": "8",
		"maxIdleConns": "3",
		"connMaxLifetime": "7200",
		"connMaxIdleTime": "1800",
		"retries": "6",
		"retryBackoff": "2",
		"maxRetryDuration": "45",
		"timeout": "15",
		"maxRows": "2500"
	}`)

	got := parseConnectionSettings(settings)

	if got.MaxOpenConns != 8 || got.MaxIdleConns != 3 {
		t.Errorf("connection counts = (%d, %d), want (8, 3)", got.MaxOpenConns, got.MaxIdleConns)
	}
	if got.ConnMaxLifetime != 2*time.Hour || got.ConnMaxIdleTime != 30*time.Minute {
		t.Errorf("connection durations = (%s, %s), want (2h, 30m)", got.ConnMaxLifetime, got.ConnMaxIdleTime)
	}
	if got.Retries != 6 || got.RetryBackoff != 2*time.Second || got.MaxRetryDuration != 45*time.Second {
		t.Errorf("retry settings = (%d, %s, %s), want (6, 2s, 45s)", got.Retries, got.RetryBackoff, got.MaxRetryDuration)
	}
	if got.Timeout != 15*time.Second || got.MaxRows != 2500 {
		t.Errorf("query settings = (%s, %d), want (15s, 2500)", got.Timeout, got.MaxRows)
	}
}

func TestParseConnectionSettingsDefaults(t *testing.T) {
	got := parseConnectionSettings(json.RawMessage(`{}`))

	if got.Retries != 4 || got.RetryBackoff != time.Second || got.MaxRetryDuration != 30*time.Second {
		t.Errorf("retry defaults = (%d, %s, %s), want (4, 1s, 30s)", got.Retries, got.RetryBackoff, got.MaxRetryDuration)
	}
	if got.Timeout != 0 || got.MaxRows != 10000 {
		t.Errorf("query defaults = (%s, %d), want (0s, 10000)", got.Timeout, got.MaxRows)
	}
}
