//go:build unit

package config

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewWith_Default(t *testing.T) {
	// Arrange
	cfgGetter := func(_ string) string {
		return ""
	}

	// Act
	got := NewWith(cfgGetter)

	// Assert
	assert.Equal(t, defaultPort, got.Port)
	assert.Equal(t, defaultDatabaseURL, got.DatabaseURL)
	assert.Equal(t, defaultAdminUsername, got.AdminUsername)
	assert.Equal(t, defaultAdminPassword, got.AdminPassword)
}

func TestNewWith_Custom(t *testing.T) {
	// Arrange
	want := Config{
		Port:          1234,
		DatabaseURL:   "database-url",
		AdminUsername: "admin",
		AdminPassword: "password",
	}
	cfgGetter := func(key string) string {
		if key == kPort {
			return strconv.Itoa(want.Port)
		}
		if key == kDatabaseURL {
			return want.DatabaseURL
		}
		if key == kAdminUsername {
			return want.AdminUsername
		}
		if key == kAdminPassword {
			return want.AdminPassword
		}
		return ""
	}

	// Act
	got := NewWith(cfgGetter)

	// Assert
	assert.Equal(t, want.Port, got.Port)
	assert.Equal(t, want.DatabaseURL, got.DatabaseURL)
	assert.Equal(t, want.AdminUsername, got.AdminUsername)
	assert.Equal(t, want.AdminPassword, got.AdminPassword)
}
