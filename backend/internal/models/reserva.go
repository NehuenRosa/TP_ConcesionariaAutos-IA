package models

import "time"

// Estado de la reserva de un vehículo.
const (
	EstadoReservaActiva    = "activa"
	EstadoReservaVendida   = "vendida"
	EstadoReservaCancelada = "cancelada"
)

// Reserva es la entidad de GORM que representa la reserva de una unidad por un
// cliente. Mientras está activa, el vehículo queda en estado reservado.
type Reserva struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	VehiculoID uint     `gorm:"not null;index" json:"vehiculoId"`
	Vehiculo   Vehiculo `gorm:"foreignKey:VehiculoID" json:"-"`
	ClienteID  uint     `gorm:"not null;index" json:"clienteId"`
	Cliente    Usuario  `gorm:"foreignKey:ClienteID" json:"-"`
	Estado     string   `gorm:"not null;index;default:activa" json:"estado"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

// TableName define el nombre de la tabla en español.
func (Reserva) TableName() string {
	return "reservas"
}

// EsActiva indica si la reserva está bloqueando la unidad.
func (r Reserva) EsActiva() bool {
	return r.Estado == EstadoReservaActiva
}
