package repository

import (
	"errors"

	"github.com/Rahima-Akter/spotsync-golang/internal/models"
	"gorm.io/gorm"
)

type ZoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) *ZoneRepository {
	return &ZoneRepository{db: db}
}

func (r *ZoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

// FindByID
func (r *ZoneRepository) FindByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	result := r.db.First(&zone, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &zone, nil
}

// FindAll
func (r *ZoneRepository) FindAll() ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	result := r.db.Order("created_at DESC").Find(&zones)
	if result.Error != nil {
		return nil, result.Error
	}
	return zones, nil
}

// Update Zone
func (r *ZoneRepository) Update(zone *models.ParkingZone) error {
	return r.db.Save(zone).Error
}

// Delete zone
func (r *ZoneRepository) Delete(id uint) error {
	return r.db.Delete(&models.ParkingZone{}, id).Error
}

// Count active zones
func (r *ZoneRepository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	result := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count)
	return count, result.Error
}
