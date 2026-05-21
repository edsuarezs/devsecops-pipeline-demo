package config

import (
	"log/slog"
	"os"

)

// Config holds application settings loaded from environment variables.
// Follows 12-factor app methodology — no hardcoded values.
type Config struct {
	AppName     string
	AppVersion  string
	Environment string // development | staging | production
	Port        string
	Debug       bool
	LogLevelStr string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		AppName:     getEnv("APP_NAME", "devsecops-pipeline-demo"),
		AppVersion:  getEnv("APP_VERSION", "0.1.0"),
		Environment: getEnv("APP_ENVIRONMENT", "development"),
		Port:        getEnv("APP_PORT", "8080"),
		Debug:       getEnv("APP_DEBUG", "true") == "true",
		LogLevelStr: getEnv("APP_LOG_LEVEL", "INFO"),
	}
}

// LogLevel returns the slog.Level based on the configured log level string.
func (c *Config) LogLevel() slog.Level {
	switch c.LogLevelStr {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
