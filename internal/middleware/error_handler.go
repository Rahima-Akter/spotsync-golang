package middleware

import (
	"net/http"

	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
	"github.com/labstack/echo/v4"
)

// CustomErrorHandler provides centralized error handling for the entire app
// This catches any unhandled errors and formats them consistently
func CustomErrorHandler(err error, c echo.Context) {
	// Check if it's an Echo HTTP error
	if he, ok := err.(*echo.HTTPError); ok {
		// Extract the error message
		message := "An error occurred"
		if he.Message != nil {
			message = he.Message.(string)
		}

		// Return formatted error response
		c.JSON(he.Code, utils.Response{
			Success: false,
			Message: message,
		})
		return
	}

	// For any other unexpected errors, return 500 Internal Server Error
	// In production, don't expose the actual error message to clients
	c.JSON(http.StatusInternalServerError, utils.Response{
		Success: false,
		Message: "Internal server error",
	})
}
