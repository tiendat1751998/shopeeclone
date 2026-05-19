package config

import "os"

type Config struct {
	AppName  string
	HTTPPort int
	LogLevel string
}

func Load() *Config {
	return &Config{
		AppName:  "rec-vector",
		HTTPPort: 8085,
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
