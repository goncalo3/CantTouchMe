package config

import (
	"os"
	"strconv"
)

// function to get an environment variable with a default value
// Parameters:
// - key: the name of the environment variable to retrieve
// - defaultValue: the default value to return if the environment variable is not set
// Returns: the value of the environment variable if set, otherwise the default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// function to get an environment variable as an integer
// Parameters:
// - key: the name of the environment variable to retrieve
// - defaultValue: the default integer value to return if the environment variable is not set or cannot be converted to an integer
// Returns: the integer value of the environment variable if set and valid, otherwise the default integer value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
