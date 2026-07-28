package models

import "time"

type TestDriveStatus string

const (
	TDStatusPending   TestDriveStatus = "pendiente"
	TDStatusConfirmed TestDriveStatus = "confirmado"
	TDStatusCancelled TestDriveStatus = "cancelado"
	TDStatusCompleted TestDriveStatus = "completado"
)

type TestDrive struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	ClientID    uint            `gorm:"not null;index" json:"client_id"`
	Client      *User           `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	VehicleID   uint            `gorm:"not null;index" json:"vehicle_id"`
	Vehicle     *Vehicle        `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	ScheduledAt time.Time       `gorm:"not null;index" json:"scheduled_at"`
	Status      TestDriveStatus `gorm:"size:20;not null;default:pendiente" json:"status"`
	Notes       string          `gorm:"type:text" json:"notes,omitempty"`
}

type CreateTestDriveRequest struct {
	VehicleID   uint   `json:"vehicle_id" binding:"required"`
	ScheduledAt string `json:"scheduled_at" binding:"required"`
	Notes       string `json:"notes,omitempty"`
}
