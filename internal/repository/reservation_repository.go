package repository

import (
	"errors"

	"github.com/Rahima-Akter/spotsync-golang/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ReservationRepository handles all database operations for reservations
type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

// CreateWithCapacityCheck creates a reservation ONLY if capacity is available
// This method uses a DATABASE TRANSACTION with ROW-LEVEL LOCKING to prevent race conditions
//
// THE CONCURRENCY PROBLEM:
// If two drivers try to book the last spot at the same time, both might
// read "19 active reservations" and both would succeed -> 21 cars in 20 spots!
//
// THE SOLUTION:
// 1. Start a transaction
// 2. Lock the parking zone row with FOR UPDATE (other transactions must wait)
// 3. Count active reservations (now safe because the row is locked)
// 4. Check if capacity is available
// 5. If yes, create the reservation
// 6. Commit the transaction (releases the lock)

func (r *ReservationRepository) CreateWithCapacityCheck(reservation *models.Reservation) error {
	// Start a database transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Lock the parking zone row with FOR UPDATE
		// This prevents other concurrent transactions from modifying this zone
		// until our transaction is complete
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, reservation.ZoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("parking zone not found")
			}
			return err
		}

		// Count current active reservations for this zone
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", reservation.ZoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// Check if there's available capacity
		if activeCount >= int64(zone.TotalCapacity) {
			return errors.New("parking zone is at full capacity")
		}

		// check if this license plate already has an active reservation
		var existingCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("license_plate = ? AND status = ?", reservation.LicensePlate, "active").
			Count(&existingCount).Error; err != nil {
			return err
		}
		if existingCount > 0 {
			return errors.New("license plate already has an active reservation")
		}

		// Create the reservation (capacity is guaranteed available)
		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindByID
func (r *ReservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	result := r.db.First(&reservation, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &reservation, nil
}

// FindByUserID
func (r *ReservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	result := r.db.Where("user_id = ?", userID).
		Preload("Zone"). // This is like: prisma.reservation.findMany({ include: { zone: true } })
		Order("created_at DESC").
		Find(&reservations)
	if result.Error != nil {
		return nil, result.Error
	}
	return reservations, nil
}

// FindAll
func (r *ReservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	result := r.db.
		Preload("User"). // Include user details
		Preload("Zone"). // Include zone details
		Order("created_at DESC").
		Find(&reservations)
	if result.Error != nil {
		return nil, result.Error
	}
	return reservations, nil
}

// CancelReservation - only the owner can cancel their reservation
func (r *ReservationRepository) CancelReservation(id uint) error {
	result := r.db.Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled")

	if result.Error != nil {
		return result.Error
	}

	// Check if any row was actually updated
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
