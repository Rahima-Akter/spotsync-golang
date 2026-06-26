package service

import (
	"github.com/Rahima-Akter/spotsync-golang/internal/dto"
	"github.com/Rahima-Akter/spotsync-golang/internal/models"
	"github.com/Rahima-Akter/spotsync-golang/internal/repository"
	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
)

// ZoneService handles all business logic for parking zones
type ZoneService struct {
	zoneRepo *repository.ZoneRepository
}

func NewZoneService(zoneRepo *repository.ZoneRepository) *ZoneService {
	return &ZoneService{zoneRepo: zoneRepo}
}

// creates a new parking zone (admin only)
func (s *ZoneService) Create(req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	// Create the zone model from the request
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	// Save to database
	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}

	return s.toZoneResponse(zone, int64(zone.TotalCapacity)), nil
}

// GetByID
func (s *ZoneService) GetByID(id uint) (*dto.ZoneResponse, error) {
	// Find the zone
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, utils.ErrNotFound
	}

	// Count active reservations to calculate availability
	activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, err
	}

	// Calculate available spots
	availableSpots := int64(zone.TotalCapacity) - activeCount
	if availableSpots < 0 {
		availableSpots = 0
	}

	return s.toZoneResponse(zone, availableSpots), nil
}

// GetAll
func (s *ZoneService) GetAll() ([]dto.ZoneResponse, error) {
	// Get all zones
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Build response for each zone with availability calculation
	var zoneResponses []dto.ZoneResponse
	for _, zone := range zones {
		// Count active reservations for this zone
		activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
		if err != nil {
			return nil, err
		}

		// Calculate available spots
		availableSpots := int64(zone.TotalCapacity) - activeCount
		if availableSpots < 0 {
			availableSpots = 0
		}

		zoneResponses = append(zoneResponses, *s.toZoneResponse(&zone, availableSpots))
	}

	return zoneResponses, nil
}

// Update a zone (admin only)
func (s *ZoneService) Update(id uint, req *dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	// Find the existing zone
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, utils.ErrNotFound
	}

	// Only update fields that were provided in the request
	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity > 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour > 0 {
		zone.PricePerHour = req.PricePerHour
	}

	// Save changes to database
	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, err
	}

	// Get active reservation count for availability
	activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, err
	}

	availableSpots := int64(zone.TotalCapacity) - activeCount
	if availableSpots < 0 {
		availableSpots = 0
	}

	return s.toZoneResponse(zone, availableSpots), nil
}

// Delete a parking zone (admin only)
func (s *ZoneService) Delete(id uint) error {
	// Check if zone exists
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return err
	}
	if zone == nil {
		return utils.ErrNotFound
	}

	// Delete the zone
	return s.zoneRepo.Delete(id)
}

// This is a helper method to avoid code duplication
func (s *ZoneService) toZoneResponse(zone *models.ParkingZone, availableSpots int64) *dto.ZoneResponse {
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: int(availableSpots),
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}
}
