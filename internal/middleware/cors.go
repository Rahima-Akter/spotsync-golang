package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupCORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// Allow all origins
		AllowOrigins: []string{"*"},

		// HTTP methods allowed
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},

		// Headers allowed in requests
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Accept",
		},

		// Expose these headers to the frontend
		ExposeHeaders: []string{
			"Content-Length",
		},

		// Allow credentials (cookies, authorization headers)
		AllowCredentials: true,

		// Cache preflight response for 1 hour (in seconds)
		MaxAge: 3600,
	}))
}
