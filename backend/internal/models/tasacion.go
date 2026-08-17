package models

import "time"

// Estado del flujo de una tasación.
const (
	// EstadoTasacionPendiente indica que la tasación se generó pero todavía
	// falta confirmar el día y la franja de la visita.
	EstadoTasacionPendiente = "pendiente"
	// EstadoTasacionConfirmada indica que el cliente confirmó día y franja y
	// ya tiene un código de presentación.
	EstadoTasacionConfirmada = "confirmada"
)

// Tasacion es la entidad de GORM que representa una tasación hecha con la IA.
// Guarda el vehículo identificado, los valores de referencia, el día y la
// franja que el cliente elige para acercarse y el código único de presentación.
type Tasacion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SesionID    string    `gorm:"size:64;not null;uniqueIndex" json:"sesionId"`
	Codigo      *string   `gorm:"size:10;uniqueIndex" json:"codigo,omitempty"`
	Marca       string    `gorm:"size:100;not null" json:"marca"`
	Modelo      string    `gorm:"size:100;not null" json:"modelo"`
	Version     string    `gorm:"size:150" json:"version"`
	Anio        int       `json:"anio"`
	Estado      string    `gorm:"size:30" json:"estado"`
	Kilometraje int       `json:"kilometraje"`
	PrecioUSD   float64   `json:"precioUsd"`
	PrecioARS   float64   `json:"precioArs"`
	Dia         string    `gorm:"size:60" json:"dia,omitempty"`
	Franja      string    `gorm:"size:30" json:"franja,omitempty"`
	EstadoFlujo string    `gorm:"size:20;not null;default:pendiente;index" json:"estadoFlujo"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// TableName define el nombre de la tabla en español.
func (Tasacion) TableName() string {
	return "tasaciones"
}