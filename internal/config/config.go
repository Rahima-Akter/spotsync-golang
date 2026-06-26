package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	// Server settings
	AppPort string
	AppEnv  string

	// Database settings
	DatabaseURL string

	// JWT settings
	JWTSecret      string
	JWTExpiryHours int
}

// Load reads configuration from .env file and environment variables
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Warning: .env file not found, using environment variables")
	}

	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRY_HOURS value")
	}

	config := &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-me"),
		JWTExpiryHours: jwtExpiryHours,
	}

	return config
}

// GetDSN returns the database connection URL directly
func (c *Config) GetDSN() string {
	return c.DatabaseURL
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
