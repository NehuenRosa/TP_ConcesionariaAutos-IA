package models

import "time"

// Estado de la reserva de un vehículo.
const (
	EstadoReservaActiva    = "activa"
	EstadoReservaVendida   = "vendida"
	EstadoReservaCancelada = "cancelada"
)

// Reserva es la entidad de GORM que representa la reserva de una unidad por un
// cliente. Mientras está activa, el vehículo queda en estado reservado y el
// cliente tiene PlazoComprobante para enviar la seña; sin comprobante al
// vencer, la reserva se anula automáticamente.
type Reserva struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	VehiculoID uint     `gorm:"not null;index" json:"vehiculoId"`
	Vehiculo   Vehiculo `gorm:"foreignKey:VehiculoID" json:"-"`
	ClienteID  uint     `gorm:"not null;index" json:"clienteId"`
	Cliente    Usuario  `gorm:"foreignKey:ClienteID" json:"-"`
	Estado     string   `gorm:"not null;index;default:activa" json:"estado"`
	// VencimientoComprobante es el límite para subir el comprobante de la
	// seña. Las reservas históricas lo tienen en cero y nunca expiran.
	VencimientoComprobante time.Time `json:"-"`
	// ComprobanteEnviadoAt marca cuándo el cliente subió el comprobante;
	// nulo significa pendiente de envío (sujeto a expiración).
	ComprobanteEnviadoAt *time.Time `json:"-"`
	// MotivoCancelacion guarda la explicación del vendedor cuando anula una
	// reserva (ej. comprobante ilegible). La cancelación propia del cliente y
	// la expiración automática lo dejan vacío.
	MotivoCancelacion string `gorm:"type:text" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TableName define el nombre de la tabla en español.
func (Reserva) TableName() string {
	return "reservas"
}

// EsActiva indica si la reserva está bloqueando la unidad.
func (r Reserva) EsActiva() bool {
	return r.Estado == EstadoReservaActiva
}

// PendienteDeComprobante indica si falta subir el comprobante de la seña.
func (r Reserva) PendienteDeComprobante() bool {
	return r.ComprobanteEnviadoAt == nil
}

// ComprobanteVencido indica si una reserva activa superó el plazo de las 2
// horas sin recibir el comprobante.
func (r Reserva) ComprobanteVencido(ahora time.Time) bool {
	return r.EsActiva() && r.PendienteDeComprobante() &&
		!r.VencimientoComprobante.IsZero() && ahora.After(r.VencimientoComprobante)
}

// ComprobanteReserva guarda la imagen del comprobante de transferencia de la
// seña (1:1 con la reserva). Se persiste en la base como bytea.
type ComprobanteReserva struct {
	ID        uint   `gorm:"primaryKey"`
	ReservaID uint   `gorm:"not null;uniqueIndex"`
	MIME      string `gorm:"not null"`
	Datos     []byte `gorm:"not null"`
	CreatedAt time.Time
}

// TableName define el nombre de la tabla en español.
func (ComprobanteReserva) TableName() string {
	return "comprobantes_reserva"
}
