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

// ReservationHandler handles HTTP requests for reservations
type ReservationHandler struct {
	reservationService *service.ReservationService
	validator          *validator.Validate
}

func NewReservationHandler(reservationService *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		validator:          validator.New(),
	}
}

// Reserve -> POST /api/v1/reservations
func (h *ReservationHandler) Reserve(c echo.Context) error {
	// Get the user ID from JWT context
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", nil)
	}

	// Parse request body
	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate the request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", validationErrors.Error())
	}

	// Call service to create reservation
	reservation, err := h.reservationService.Reserve(userID, &req)
	if err != nil {
		switch err {
		case utils.ErrNotFound:
			return utils.ErrorResponse(c, http.StatusNotFound, "Parking zone not found", nil)
		case utils.ErrZoneFull:
			return utils.ErrorResponse(c, http.StatusConflict, "Parking zone is at full capacity", "No available spots in this zone")
		case utils.ErrDuplicateReservation:
			return utils.ErrorResponse(c, http.StatusConflict, "License plate already has an active reservation", nil)
		default:
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create reservation", err.Error())
		}
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Reservation confirmed successfully", reservation)
}

// GetMyReservations -> GET /api/v1/reservations/my-reservations
func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", nil)
	}

	// Call service
	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reservations", err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "My reservations retrieved successfully", reservations)
}

// Cancel -> DELETE /api/v1/reservations/:id (own reservation)
func (h *ReservationHandler) Cancel(c echo.Context) error {
	// Get user ID from JWT context
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", nil)
	}

	// Extract reservation ID from URL
	reservationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid reservation ID", nil)
	}

	// Call service to cancel
	if err := h.reservationService.CancelReservation(uint(reservationID), userID); err != nil {
		switch err {
		case utils.ErrNotFound:
			return utils.ErrorResponse(c, http.StatusNotFound, "Reservation not found", nil)
		case utils.ErrForbidden:
			return utils.ErrorResponse(c, http.StatusForbidden, "You can only cancel your own reservations", nil)
		default:
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel reservation", err.Error())
		}
	}

	return utils.SuccessResponse(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

// GetAll -> GET /api/v1/reservations (Admin only)
func (h *ReservationHandler) GetAll(c echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reservations", err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "All reservations retrieved successfully", reservations)
}
