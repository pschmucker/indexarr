package config

import (
	"os"
)

type Config struct {
	ServerPort string
	DBPath     string
	TMDBAPIKey string
	TVDBAPIKey string
	ScanInterval int // hours
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBPath: getEnv("DB_PATH", "./indexarr.db"),
		TMDBAPIKey: getEnv("TMDB_API_KEY", ""),
		TVDBAPIKey: getEnv("TVDB_API_KEY", ""),
		ScanInterval: 24, // hours
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
