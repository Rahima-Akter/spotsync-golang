package middleware

import (
	"net/http"
	"strings"

	"github.com/Rahima-Akter/spotsync-golang/internal/config"
	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// 1. Extract the token from "Authorization: Bearer <token>" header
// 2. Validate and parse the token
// 3. Store user info (id, email, role) in Echo context for later use
func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header
			authHeader := c.Request().Header.Get("Authorization")

			//  Check if the header is missing
			if authHeader == "" {
				return utils.ErrorResponse(
					c,
					http.StatusUnauthorized,
					"Missing authorization token",
					"Please provide a valid JWT token in the Authorization header",
				)
			}

			// Check if it's a Bearer token, Like -> "Bearer eyJhbGciOiJIUzI1NiIs..."
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return utils.ErrorResponse(
					c,
					http.StatusUnauthorized,
					"Invalid authorization format",
					"Authorization header must be in format: Bearer <token>",
				)
			}

			// Extract the token string
			tokenString := parts[1]

			// Parse and validate the JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Verify the signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				// Return the secret key for verification
				return []byte(cfg.JWTSecret), nil
			})

			// Handle invalid tokens
			if err != nil || !token.Valid {
				return utils.ErrorResponse(
					c,
					http.StatusUnauthorized,
					"Invalid or expired token",
					"Please login again to get a new token",
				)
			}

			// Extract claims (user data from token payload)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return utils.ErrorResponse(
					c,
					http.StatusUnauthorized,
					"Invalid token claims",
					nil,
				)
			}

			// Extract specific user information from claims
			userID, ok := claims["user_id"].(float64)
			if !ok {
				return utils.ErrorResponse(
					c,
					http.StatusUnauthorized,
					"Invalid user ID in token",
					nil,
				)
			}

			email, ok := claims["email"].(string)
			if !ok {
				return utils.ErrorResponse(
					c,
					http.StatusUnauthorized,
					"Invalid email in token",
					nil,
				)
			}

			role, ok := claims["role"].(string)
			if !ok {
				return utils.ErrorResponse(
					c,
					http.StatusUnauthorized,
					"Invalid role in token",
					nil,
				)
			}

			// Store user info in Echo context
			c.Set("user_id", uint(userID))
			c.Set("email", email)
			c.Set("role", role)

			// Call the next handler
			return next(c)
		}
	}
}
