package handler

import (
	"net/http"
	"strconv"

	"github.com/Rahima-Akter/spotsync-golang/internal/dto"
	"github.com/Rahima-Akter/spotsync-golang/internal/service"
	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ZoneHandler handles HTTP requests for parking zones
type ZoneHandler struct {
	zoneService *service.ZoneService
	validator   *validator.Validate
}

func NewZoneHandler(zoneService *service.ZoneService) *ZoneHandler {
	return &ZoneHandler{
		zoneService: zoneService,
		validator:   validator.New(),
	}
}

// Create -> POST /api/v1/zones (Admin only)
func (h *ZoneHandler) Create(c echo.Context) error {
	// Parse request body
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate the request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", validationErrors.Error())
	}

	// Call service to create zone
	zone, err := h.zoneService.Create(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create parking zone", err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Parking zone created successfully", zone)
}

// GetAll -> GET /api/v1/zones (Public)
func (h *ZoneHandler) GetAll(c echo.Context) error {
	zones, err := h.zoneService.GetAll()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve parking zones", err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

// GetByID -> GET /api/v1/zones/:id (Public)
func (h *ZoneHandler) GetByID(c echo.Context) error {
	// Extract ID from URL parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	zone, err := h.zoneService.GetByID(uint(id))
	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return utils.ErrorResponse(c, http.StatusNotFound, "Parking zone not found", nil)
		default:
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve parking zone", err.Error())
		}
	}

	return utils.SuccessResponse(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}

// Update -> PUT /api/v1/zones/:id (Admin only)
func (h *ZoneHandler) Update(c echo.Context) error {
	// Extract ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	// Parse request body
	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate the request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", validationErrors.Error())
	}

	// Call service to update zone
	zone, err := h.zoneService.Update(uint(id), &req)
	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return utils.ErrorResponse(c, http.StatusNotFound, "Parking zone not found", nil)
		default:
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update parking zone", err.Error())
		}
	}

	return utils.SuccessResponse(c, http.StatusOK, "Parking zone updated successfully", zone)
}

// Delete -> DELETE /api/v1/zones/:id (Admin only)
func (h *ZoneHandler) Delete(c echo.Context) error {
	// Extract ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	// Call service to delete zone
	if err := h.zoneService.Delete(uint(id)); err != nil {
		switch err {
		case utils.ErrNotFound:
			return utils.ErrorResponse(c, http.StatusNotFound, "Parking zone not found", nil)
		default:
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete parking zone", err.Error())
		}
	}

	return utils.SuccessResponse(c, http.StatusOK, "Parking zone deleted successfully", nil)
}
