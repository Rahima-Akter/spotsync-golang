package models

import (
	"time"
)

// For example: "Terminal 1 EV Charging" with 20 spots
type ParkingZone struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Type          string    `gorm:"not null" json:"type"` // "general", "ev_charging", or "covered"
	TotalCapacity int       `gorm:"not null" json:"total_capacity"`
	PricePerHour  float64   `gorm:"not null" json:"price_per_hour"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ParkingZone) TableName() string {
	return "parking_zones"
}
