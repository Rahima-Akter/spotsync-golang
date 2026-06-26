package handler

import (
	"net/http"

	"github.com/Rahima-Akter/spotsync-golang/internal/dto"
	"github.com/Rahima-Akter/spotsync-golang/internal/service"
	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// AuthHandler handles HTTP requests for authentication
// This is like an Express.js controller:
// router.post('/register', authController.register)
type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// Register handles POST /api/v1/auth/register
// This is like:
//
//	app.post('/api/v1/auth/register', async (req, res) => {
//	  const { name, email, password, role } = req.body;
//	})
func (h *AuthHandler) Register(c echo.Context) error {
	// Parse the request body into our DTO
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate the DTO fields
	if err := h.validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", validationErrors.Error())
	}

	// Call the service layer to register the user
	userResponse, err := h.authService.Register(&req)
	if err != nil {
		// Check error type
		switch err {
		case utils.ErrDuplicateEmail:
			return utils.ErrorResponse(c, http.StatusBadRequest, "Email already registered", err.Error())
		default:
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user", err.Error())
		}
	}

	// Return success response
	return utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", userResponse)
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	// Parse the request body
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate the fields
	if err := h.validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", validationErrors.Error())
	}

	// Call the service layer to authenticate
	loginResponse, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case utils.ErrInvalidCredentials:
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
		default:
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login", err.Error())
		}
	}

	//  Return the JWT token and user data
	return utils.SuccessResponse(c, http.StatusOK, "Login successful", loginResponse)
}
