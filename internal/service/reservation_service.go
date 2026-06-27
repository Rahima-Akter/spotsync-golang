package service

import (
	"errors"

	"github.com/Rahima-Akter/spotsync-golang/internal/dto"
	"github.com/Rahima-Akter/spotsync-golang/internal/models"
	"github.com/Rahima-Akter/spotsync-golang/internal/repository"
	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
)

// ReservationService handles all business logic for reservations
type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	zoneRepo        *repository.ZoneRepository
}

func NewReservationService(
	reservationRepo *repository.ReservationRepository,
	zoneRepo *repository.ZoneRepository,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

// Reserve creates a new reservation
func (s *ReservationService) Reserve(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	// Verify the zone exists
	zone, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, utils.ErrNotFound
	}

	// Create the reservation model
	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	// Try to create with capacity check (uses transaction + row locking)
	err = s.reservationRepo.CreateWithCapacityCheck(reservation)
	if err != nil {
		// Check what kind of error occurred
		switch err.Error() {
		case "parking zone is at full capacity":
			return nil, utils.ErrZoneFull
		case "license plate already has an active reservation":
			return nil, utils.ErrDuplicateReservation
		case "parking zone not found":
			return nil, utils.ErrNotFound
		default:
			return nil, err
		}
	}

	// return the reservation response
	return s.toReservationResponse(reservation), nil
}

// GetMyReservations
func (s *ReservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	var responses []dto.MyReservationResponse
	for _, r := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	return responses, nil
}

// CancelReservation
func (s *ReservationService) CancelReservation(reservationID uint, userID uint) error {
	// Find the reservation
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		return err
	}
	if reservation == nil {
		return utils.ErrNotFound
	}

	// Check if the user owns this reservation
	if reservation.UserID != userID {
		return utils.ErrForbidden
	}

	// Check if reservation is already cancelled
	if reservation.Status == "cancelled" {
		return errors.New("reservation is already cancelled")
	}

	return s.reservationRepo.CancelReservation(reservationID)
}

// GetAllReservations (admin only)
func (s *ReservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Convert to admin response DTOs
	var responses []dto.AdminReservationResponse
	for _, r := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.UserInfo{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	return responses, nil
}

// ReservationResponse DTO
func (s *ReservationService) toReservationResponse(r *models.Reservation) *dto.ReservationResponse {
	return &dto.ReservationResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		ZoneID:       r.ZoneID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
