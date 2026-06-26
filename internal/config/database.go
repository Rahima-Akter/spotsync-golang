package config

import (
	"fmt"
	"log"

	"github.com/Rahima-Akter/spotsync-golang/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB creates a connection to the neon PostgreSQL database
// This is like PrismaClient in Express.js
func ConnectDB(cfg *Config) (*gorm.DB, error) {
	// Get the connection string (DSN)
	dsn := cfg.GetDSN()

	// Open a connection to PostgreSQL using GORM
	// gorm.Open is like new PrismaClient() in Express.js
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Set logger level based on environment
		// In development, we want to see all SQL queries
		// In production, we only want to see errors
		Logger: logger.Default.LogMode(getLogLevel(cfg.AppEnv)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying SQL database connection for ping test
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Test the connection by pinging the database
	// This ensures neondb is reachable
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to NeonDB PostgreSQL database!")

	// Auto-migrate the schema
	// This is like running 'prisma migrate dev' in Express.js
	// It automatically creates/updates tables based on our models
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}

// autoMigrate creates/updates database tables based on our models
// This is like Prisma's migrate command
func autoMigrate(db *gorm.DB) error {
	log.Println("Running auto-migration...")

	// AutoMigrate will create tables, add missing columns, and add missing indexes
	err := db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully!")
	return nil
}

// getLogLevel returns the appropriate GORM log level based on environment
func getLogLevel(env string) logger.LogLevel {
	if env == "production" {
		// In production, only log errors to avoid noise
		return logger.Error
	}
	// In development, log all SQL queries (helpful for debugging)
	return logger.Info
}
