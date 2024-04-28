package config

import (
	"strconv"
)

const (
	kPort       = "PORT"
	defaultPort = 8080

	kDatabaseURL       = "DATABASE_URL"
	defaultDatabaseURL = "postgresql://postgres:postgres@localhost:5432/ktaxes?sslmode=disable"

	kAdminUsername       = "ADMIN_USERNAME"
	defaultAdminUsername = ""

	kAdminPassword       = "ADMIN_PASSWORD"
	defaultAdminPassword = ""
)

type ConfigGetter func(string) string

type Config struct {
	Port          int
	DatabaseURL   string
	AdminUsername string
	AdminPassword string
}

func NewWith(cfgGetter ConfigGetter) *Config {
	return &Config{
		Port:          getInt(cfgGetter, kPort, defaultPort),
		DatabaseURL:   getString(cfgGetter, kDatabaseURL, defaultDatabaseURL),
		AdminUsername: getString(cfgGetter, kAdminUsername, defaultAdminUsername),
		AdminPassword: getString(cfgGetter, kAdminPassword, defaultAdminPassword),
	}
}

func getString(fn ConfigGetter, key, defaultValue string) string {
	result := fn(key)
	if result == "" {
		return defaultValue
	}
	return result
}

func getInt(fn ConfigGetter, key string, defaultValue int) int {
	v := fn(key)
	result, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return result
}
