package models

import "time"

type ReservationStatus string

const (
	ReservationActive    ReservationStatus = "activa"
	ReservationConfirmed ReservationStatus = "confirmada"
	ReservationCancelled ReservationStatus = "cancelada"
)

type Reservation struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	ClientID  uint              `gorm:"not null;index" json:"client_id"`
	Client    User              `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	VehicleID uint              `gorm:"not null;index" json:"vehicle_id"`
	Vehicle   Vehicle           `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Status    ReservationStatus `gorm:"size:20;not null;default:activa" json:"status"`
	Notes     string            `gorm:"type:text" json:"notes,omitempty"`
}

type CreateReservationRequest struct {
	VehicleID uint   `json:"vehicle_id" binding:"required"`
	Notes     string `json:"notes,omitempty"`
}
