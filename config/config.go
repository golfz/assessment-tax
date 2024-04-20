package config

import (
	"os"
	"strconv"
)

const (
	kPort        = "PORT"
	kDatabaseURL = "DATABASE_URL"

	defaultPort        = 8080
	defaultDatabaseURL = "postgresql://postgres:postgres@localhost:5432/ktaxes?sslmode=disable"
)

type ConfigGetter func(string) string

type Config struct {
	Port        int
	DatabaseURL string
}

func NewWith(cfgGetter ConfigGetter) *Config {
	return &Config{
		Port:        getInt(cfgGetter, kPort, defaultPort),
		DatabaseURL: getString(cfgGetter, kDatabaseURL, defaultDatabaseURL),
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
	v := os.Getenv(key)
	result, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return result
}
