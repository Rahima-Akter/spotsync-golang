package middleware

import (
	"net/http"

	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
	"github.com/labstack/echo/v4"
)

func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the user's role from context set by JWTAuth middleware
			role, ok := c.Get("role").(string)
			if !ok {
				return utils.ErrorResponse(
					c,
					http.StatusForbidden,
					"Access denied",
					"User role not found in context",
				)
			}

			// Check if the user's role is in the allowed list
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					return next(c)
				}
			}

			// Role not allowed - return 403 Forbidden
			return utils.ErrorResponse(
				c,
				http.StatusForbidden,
				"Permission denied",
				"This action requires one of these roles: "+joinRoles(allowedRoles),
			)
		}
	}
}

// joinRoles is a helper function to format role names nicely
func joinRoles(roles []string) string {
	result := ""
	for i, role := range roles {
		if i > 0 {
			result += ", "
		}
		result += role
	}
	return result
}
