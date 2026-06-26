package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rahima-Akter/spotsync-golang/internal/config"
	"github.com/Rahima-Akter/spotsync-golang/internal/handler"
	"github.com/Rahima-Akter/spotsync-golang/internal/repository"
	"github.com/Rahima-Akter/spotsync-golang/internal/router"
	"github.com/Rahima-Akter/spotsync-golang/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Load config
	cfg := config.Load()

	log.Println("Starting SpotSync server...")

	// Connect database
	log.Println("Connecting to database...")
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully!")

	// Get underlying SQL DB for cleanup (from OLD code)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql db:", err)
	}
	defer sqlDB.Close()

	// Create Echo server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	// Setup Routes
	router.SetupRoutes(e, authHandler, cfg)

	// Start server
	go func() {
		log.Printf("Server running on http://localhost:%s", cfg.AppPort)
		log.Printf("API Base URL: http://localhost:%s/api/v1", cfg.AppPort)

		if err := e.Start(":" + cfg.AppPort); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down SpotSync server...")

	// Echo shutdown
	if err := e.Close(); err != nil {
		log.Println("Error during server shutdown:", err)
	}

	log.Println("Server stopped cleanly")
}
