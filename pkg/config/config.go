package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        int
	Environment string
	LogLevel    string
}

func Load() *Config {
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		Port:        port,
		Environment: env,
		LogLevel:    logLevel,
	}
}
