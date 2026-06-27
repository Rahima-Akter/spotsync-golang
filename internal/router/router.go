package router

import (
	"github.com/Rahima-Akter/spotsync-golang/internal/config"
	"github.com/Rahima-Akter/spotsync-golang/internal/handler"
	"github.com/Rahima-Akter/spotsync-golang/internal/middleware"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
// This is like:
// const router = express.Router();
// router.post('/auth/register', authController.register);
// app.use('/api/v1', router);
func SetupRoutes(e *echo.Echo, authHandler *handler.AuthHandler, zoneHandler *handler.ZoneHandler,
	reservationHandler *handler.ReservationHandler, cfg *config.Config) {

	api := e.Group("/api/v1")

	// Auth routes (public routes)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// PROTECTED ROUTES
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg)) // All routes in this group require JWT

	// Test handler for verifying auth works
	testHandler := handler.NewTestHandler()
	protected.GET("/profile", testHandler.GetProfile)
	protected.GET("/admin/dashboard", testHandler.AdminDashboard, middleware.RequireRole("admin"))
	protected.GET("/driver/dashboard", testHandler.DriverDashboard, middleware.RequireRole("driver"))

	// Admin-only routes
	adminOnly := protected.Group("")
	adminOnly.Use(middleware.RequireRole("admin"))
	adminOnly.GET("/admin/dashboard", testHandler.AdminDashboard)

	// Driver routes
	driverOnly := protected.Group("")
	driverOnly.Use(middleware.RequireRole("driver"))
	driverOnly.GET("/driver/dashboard", testHandler.DriverDashboard)

	// zones
	// public
	api.GET("/zones", zoneHandler.GetAll)
	api.GET("/zones/:id", zoneHandler.GetByID)

	// admin-only
	adminZones := protected.Group("/zones")
	adminZones.Use(middleware.RequireRole("admin"))
	adminZones.POST("", zoneHandler.Create)
	adminZones.PUT("/:id", zoneHandler.Update)
	adminZones.DELETE("/:id", zoneHandler.Delete)

	// reservations
	reservations := protected.Group("/reservations")
	reservations.POST("", reservationHandler.Reserve)
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservations.DELETE("/:id", reservationHandler.Cancel)

	// Admin-only
	reservations.GET("", reservationHandler.GetAll, middleware.RequireRole("admin"))
}
