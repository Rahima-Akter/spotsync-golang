package router

import (
	"github.com/Rahima-Akter/spotsync-golang/internal/config"
	"github.com/Rahima-Akter/spotsync-golang/internal/handler"
	"github.com/Rahima-Akter/spotsync-golang/internal/middleware"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	zoneHandler *handler.ZoneHandler,
	reservationHandler *handler.ReservationHandler,
	cfg *config.Config,
) {

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "ok",
			"message": "SpotSync API is running",
		})
	})

	api := e.Group("/api/v1")

	// AUTH ROUTES (Public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// ZONE ROUTES (Public Read)
	api.GET("/zones", zoneHandler.GetAll)
	api.GET("/zones/:id", zoneHandler.GetByID)

	// PROTECTED ROUTES (Require JWT)
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg))

	// ZONE ROUTES (Admin Only)
	adminZones := protected.Group("/zones")
	adminZones.Use(middleware.RequireRole("admin"))
	adminZones.POST("", zoneHandler.Create)
	adminZones.PUT("/:id", zoneHandler.Update)
	adminZones.DELETE("/:id", zoneHandler.Delete)

	// RESERVATION ROUTES (Authenticated Users)
	reservations := protected.Group("/reservations")
	reservations.POST("", reservationHandler.Reserve)
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservations.DELETE("/:id", reservationHandler.Cancel)

	// Admin-only: View all reservations
	reservations.GET("", reservationHandler.GetAll, middleware.RequireRole("admin"))
}
