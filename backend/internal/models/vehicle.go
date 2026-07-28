package models

import (
	"time"

	"github.com/lib/pq"
)

type VehicleStatus string

const (
	VehicleAvailable VehicleStatus = "disponible"
	VehicleReserved  VehicleStatus = "reservado"
	VehicleSold      VehicleStatus = "vendido"
)

type FuelType string

const (
	FuelGasoline  FuelType = "nafta"
	FuelDiesel    FuelType = "diesel"
	FuelElectric  FuelType = "electrico"
	FuelHybrid    FuelType = "hibrido"
)

type TransmissionType string

const (
	TransmissionManual    TransmissionType = "manual"
	TransmissionAutomatic TransmissionType = "automatico"
)

type VehicleCondition string

const (
	ConditionNew  VehicleCondition = "nuevo"
	ConditionUsed VehicleCondition = "usado"
)

type Vehicle struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Brand        string            `gorm:"size:50;not null;index" json:"brand"`
	Model        string            `gorm:"size:50;not null;index" json:"model"`
	Year         int               `gorm:"not null;index" json:"year"`
	Price        float64           `gorm:"not null;index" json:"price"`
	Mileage      int               `json:"mileage"`
	Fuel         FuelType          `gorm:"size:20;not null" json:"fuel"`
	Transmission TransmissionType  `gorm:"size:20;not null" json:"transmission"`
	Condition    VehicleCondition  `gorm:"size:10;not null" json:"condition"`
	Color        string            `gorm:"size:30" json:"color,omitempty"`
	Description  string            `gorm:"type:text" json:"description,omitempty"`
	Images       pq.StringArray    `gorm:"type:text[]" json:"images"`
	Status       VehicleStatus     `gorm:"size:20;not null;default:disponible;index" json:"status"`
	VehicleType  string            `gorm:"size:30;index" json:"vehicle_type"`
}

type VehicleFilter struct {
	Search      string   `form:"search"`
	Brand       string   `form:"brand"`
	Model       string   `form:"model"`
	YearFrom    *int     `form:"year_from"`
	YearTo      *int     `form:"year_to"`
	PriceFrom   *float64 `form:"price_from"`
	PriceTo     *float64 `form:"price_to"`
	Fuel        string   `form:"fuel"`
	Condition   string   `form:"condition"`
	VehicleType string   `form:"vehicle_type"`
	SortBy      string   `form:"sort_by"`
	SortOrder   string   `form:"sort_order"`
	Page        int      `form:"page"`
	PageSize    int      `form:"page_size"`
}
