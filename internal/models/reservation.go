package models

import (
	"time"
)

// This links a user to a parking zone they've booked
type Reservation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	ZoneID       uint      `gorm:"not null" json:"zone_id"`
	LicensePlate string    `gorm:"not null;size:15" json:"license_plate"` // Max 15 characters
	Status       string    `gorm:"default:active;not null" json:"status"` // "active", "completed", or "cancelled"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Foreign key relationships
	User User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone ParkingZone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
}

func (Reservation) TableName() string {
	return "reservations"
}
