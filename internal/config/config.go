package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server settings
	AppPort string
	AppEnv  string

	// Supabase Database settings
	DatabaseURL string

	// JWT settings
	JWTSecret      string
	JWTExpiryHours int
}

func Load() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Warning: .env file not found, using environment variables")
	}

	// Parse JWT expiry hours (convert string to int)
	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRY_HOURS value")
	}

	config := &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", "jwt_very_strong_secret"),
		JWTExpiryHours: jwtExpiryHours,
	}

	return config
}

// DSN = Data Source Name
func (c *Config) GetDSN() string {
	// if database url exists
	return c.DatabaseURL

}

// This is a helper function to avoid repetitive code
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
