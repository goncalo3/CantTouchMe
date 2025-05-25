package config

// Config holds all configuration for our application
type DbConfig struct {
	Port   int
	DbUser string
	DbPwd  string
	DbName string
	DbHost string
}

// LoadConfig loads the configuration from environment variables
func LoadDbConfig() *DbConfig {
	config := &DbConfig{
		Port:   getEnvAsInt("MYSQL_PORT", 3306),
		DbUser: getEnv("MYSQL_USER", "user"),
		DbPwd:  getEnv("MYSQL_PASSWORD", "password"),
		DbName: getEnv("MYSQL_DATABASE", "canttouchme"),
		DbHost: getEnv("MYSQL_HOST", "127.0.0.1"),
	}

	return config
}
