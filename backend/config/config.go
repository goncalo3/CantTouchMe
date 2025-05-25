package config

var cfg *Config

// Config holds all configuration for our application
type Config struct {
	Port                    int
	Environment             string
	JWTSecret               string
	JWTExpiration           int // JWT expiration time in seconds
	ChallengeCleanupMinutes int // Challenge cleanup interval in minutes
}

// LoadConfig loads the configuration from environment variables
// Reurns - a pointer to Config struct
func LoadConfig() *Config {
	if cfg != nil {
		return cfg
	}

	cfg = &Config{
		Port:                    getEnvAsInt("API_PORT", 3000),
		Environment:             getEnv("ENVIRONMENT", "development"),
		JWTSecret:               getEnv("JWT_SECRET", "DEFAULT_JWT_DO_NOT_USE_IN_PRODUCTION"),
		JWTExpiration:           getEnvAsInt("JWT_EXPIRATION_SECONDS", 3600),  // Default to 1 hour
		ChallengeCleanupMinutes: getEnvAsInt("CHALLENGE_CLEANUP_MINUTES", 15), // Default to 15 minutes
	}

	return cfg
}

// GetConfig returns the current configuration without reloading from environment
func GetConfig() *Config {
	if cfg == nil {
		return LoadConfig()
	}
	return cfg
}
