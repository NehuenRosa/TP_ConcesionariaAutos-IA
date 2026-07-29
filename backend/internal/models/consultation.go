package models

import "time"

type ConsultationStatus string

const (
	ConsultPending    ConsultationStatus = "pendiente"
	ConsultInProgress ConsultationStatus = "en_conversacion"
	ConsultClosed     ConsultationStatus = "cerrada"
)

type Consultation struct {
	ID               uint               `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	ClientID         uint               `gorm:"not null;index" json:"client_id"`
	Client           *User              `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	VehicleID        uint               `gorm:"not null;index" json:"vehicle_id"`
	Vehicle          *Vehicle           `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Message          string             `gorm:"type:text;not null" json:"message"`
	Status           ConsultationStatus `gorm:"size:20;not null;default:pendiente" json:"status"`
	AssignedTo       *uint              `gorm:"index" json:"assigned_to,omitempty"`
	Seller           *User              `gorm:"foreignKey:AssignedTo" json:"seller,omitempty"`
	HasUnreadMessages  bool              `gorm:"default:false" json:"has_unread_messages"`
	HasUnreadForClient bool              `gorm:"default:false" json:"has_unread_for_client"`
	Responses          []ConsultationResponse `gorm:"foreignKey:ConsultationID" json:"responses,omitempty"`
}

type ConsultationResponse struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	ConsultationID uint      `gorm:"not null;index" json:"consultation_id"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	User           *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Message        string    `gorm:"type:text;not null" json:"message"`
}

type CreateConsultationRequest struct {
	VehicleID uint   `json:"vehicle_id" binding:"required"`
	Message   string `json:"message" binding:"required,min=10"`
}
