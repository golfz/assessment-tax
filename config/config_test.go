package config

import (
	"strconv"
	"testing"
)

func TestNewWith_Default(t *testing.T) {
	// Arrange
	cfgGetter := func(key string) string {
		return ""
	}

	// Act
	cfg := NewWith(cfgGetter)

	// Assert
	if cfg.Port != defaultPort {
		t.Errorf("want %d, got %d", defaultPort, cfg.Port)
	}
	if cfg.DatabaseURL != defaultDatabaseURL {
		t.Errorf("want %s, got %s", defaultDatabaseURL, cfg.DatabaseURL)
	}
}

func TestNewWith_Custom(t *testing.T) {
	// Arrange
	wantPort := 1234
	wantDatabaseURL := "database-url"
	cfgGetter := func(key string) string {
		if key == kPort {
			return strconv.Itoa(wantPort)
		}
		if key == kDatabaseURL {
			return wantDatabaseURL
		}
		return ""
	}

	// Act
	cfg := NewWith(cfgGetter)

	// Assert
	if cfg.Port != wantPort {
		t.Errorf("want %d, got %d", wantPort, cfg.Port)
	}
	if cfg.DatabaseURL != wantDatabaseURL {
		t.Errorf("want %s, got %s", wantDatabaseURL, cfg.DatabaseURL)
	}
}
