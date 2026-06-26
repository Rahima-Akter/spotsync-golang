package handler

import (
	"net/http"

	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
	"github.com/labstack/echo/v4"
)

// TestHandler provides test endpoints to verify authentication
type TestHandler struct{}

// NewTestHandler creates a new test handler
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

// GetProfile is a protected route that returns the current user's info
// This proves the JWT middleware worked and extracted user data correctly
func (h *TestHandler) GetProfile(c echo.Context) error {
	// Get user info from context (set by JWTAuth middleware)
	userID := c.Get("user_id").(uint)
	email := c.Get("email").(string)
	role := c.Get("role").(string)

	// Return the user's profile from the token
	return utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", map[string]interface{}{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}

// AdminDashboard is a route that only admins can access
func (h *TestHandler) AdminDashboard(c echo.Context) error {
	return utils.SuccessResponse(c, http.StatusOK, "Welcome to admin dashboard", map[string]string{
		"message": "Only admins can see this!",
		"secret":  "The parking zone pricing formula is...",
	})
}

// DriverDashboard is a route that only drivers can access
func (h *TestHandler) DriverDashboard(c echo.Context) error {
	return utils.SuccessResponse(c, http.StatusOK, "Welcome to driver dashboard", map[string]string{
		"message": "Here are your parking options!",
	})
}
