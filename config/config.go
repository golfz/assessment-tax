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
	cfgGetter   ConfigGetter
}

func New() *Config {
	return &Config{
		Port:        getInt(kPort, defaultPort),
		DatabaseURL: getString(kDatabaseURL, defaultDatabaseURL),
	}
}

func getString(key, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		return defaultValue
	}
	return result
}

func getInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	result, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return result
}
