package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for our application
type Config struct {
	Port        int
	Environment string
	LogLevel    string
	DatabaseURL string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Port:        getEnvAsInt("API_PORT", 3000),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
	
	return config
}

// Helper function to get an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}