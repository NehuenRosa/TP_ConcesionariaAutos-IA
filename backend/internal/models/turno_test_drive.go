package models

import "time"

// Estado del turno de test drive.
const (
	EstadoTurnoSolicitado = "solicitado"
	EstadoTurnoConfirmado = "confirmado"
	EstadoTurnoCancelado  = "cancelado"
	EstadoTurnoCompletado = "completado"
)

// FranjaHoraria es una franja predefinida para los turnos de test drive.
type FranjaHoraria struct {
	ID     string `json:"id"`
	Inicio string `json:"inicio"`
	Fin    string `json:"fin"`
}

// FranjasDisponibles devuelve el catálogo de franjas horarias predefinidas.
func FranjasDisponibles() []FranjaHoraria {
	return []FranjaHoraria{
		{ID: "manana", Inicio: "09:00", Fin: "12:00"},
		{ID: "tarde", Inicio: "14:00", Fin: "18:00"},
	}
}

// FranjaValida indica si un identificador pertenece al catálogo de franjas.
func FranjaValida(franja string) bool {
	for _, disponible := range FranjasDisponibles() {
		if disponible.ID == franja {
			return true
		}
	}
	return false
}

// TurnoTestDrive es la entidad de GORM que representa un turno de prueba de
// manejo solicitado por un cliente para un vehículo.
type TurnoTestDrive struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VehiculoID uint      `gorm:"not null;index" json:"vehiculoId"`
	Vehiculo   Vehiculo  `gorm:"foreignKey:VehiculoID" json:"-"`
	ClienteID  uint      `gorm:"not null;index" json:"clienteId"`
	Cliente    Usuario   `gorm:"foreignKey:ClienteID" json:"-"`
	Fecha      string    `gorm:"type:date;not null;index" json:"fecha"`
	Franja     string    `gorm:"not null" json:"franja"`
	Estado     string    `gorm:"not null;index;default:solicitado" json:"estado"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

// TableName define el nombre de la tabla en español.
func (TurnoTestDrive) TableName() string {
	return "turnos_test_drive"
}

// EsActivo indica si el turno está en un estado que bloquea la franja
// (solicitado o confirmado).
func (t TurnoTestDrive) EsActivo() bool {
	return t.Estado == EstadoTurnoSolicitado || t.Estado == EstadoTurnoConfirmado
}
