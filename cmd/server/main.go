package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rahima-Akter/spotsync-golang/internal/config"
	"github.com/Rahima-Akter/spotsync-golang/internal/handler"
	customMiddleware "github.com/Rahima-Akter/spotsync-golang/internal/middleware"
	"github.com/Rahima-Akter/spotsync-golang/internal/repository"
	"github.com/Rahima-Akter/spotsync-golang/internal/router"
	"github.com/Rahima-Akter/spotsync-golang/internal/service"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {

	// Load configuration
	cfg := config.Load()

	log.Println("Starting SpotSync server...")
	log.Printf("Environment: %s", cfg.AppEnv)
	log.Printf("Port: %s", cfg.AppPort)

	// Connect database
	log.Println("🔌 Connecting to NeonDB database...")
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully!")

	// Get SQL DB for cleanup
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql db:", err)
	}
	defer sqlDB.Close()

	// Create Echo server
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	customMiddleware.SetupCORS(e)

	// Hide Echo banner in production
	if cfg.AppEnv == "production" {
		e.HideBanner = true
	}

	// Custom error handler
	e.HTTPErrorHandler = customMiddleware.CustomErrorHandler

	// Repositories
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	zoneService := service.NewZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo, zoneRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)

	// Setup Routes
	router.SetupRoutes(e, authHandler, zoneHandler, reservationHandler, cfg)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start Server
	go func() {
		log.Printf("SpotSync server is running on http://localhost:%s", cfg.AppPort)
		log.Printf("API Base URL: http://localhost:%s/api/v1", cfg.AppPort)
		log.Printf("Health Check: http://localhost:%s/health", cfg.AppPort)

		if err := e.Start(":" + cfg.AppPort); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for shutdown signal
	<-stop

	log.Println("Shutting down SpotSync server...")

	// Allow existing requests to finish (5 second timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully stop Echo
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	} else {
		log.Println("Echo server stopped")
	}

	// Close database connection
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed")
	}

	log.Println("Server stopped gracefully")
}
