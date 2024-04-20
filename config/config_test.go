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
	got := NewWith(cfgGetter)

	// Assert
	if got.Port != defaultPort {
		t.Errorf("want %d, got %d", defaultPort, got.Port)
	}
	if got.DatabaseURL != defaultDatabaseURL {
		t.Errorf("want %s, got %s", defaultDatabaseURL, got.DatabaseURL)
	}
}

func TestNewWith_Custom(t *testing.T) {
	// Arrange
	want := Config{
		Port:        1234,
		DatabaseURL: "database-url",
	}
	cfgGetter := func(key string) string {
		if key == kPort {
			return strconv.Itoa(want.Port)
		}
		if key == kDatabaseURL {
			return want.DatabaseURL
		}
		return ""
	}

	// Act
	got := NewWith(cfgGetter)

	// Assert
	if got.Port != want.Port {
		t.Errorf("want %d, got %d", want.Port, got.Port)
	}
	if got.DatabaseURL != want.DatabaseURL {
		t.Errorf("want %s, got %s", want.DatabaseURL, got.DatabaseURL)
	}
}
