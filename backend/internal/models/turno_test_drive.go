package models

import "time"

// Estado del turno de test drive.
const (
	EstadoTurnoSolicitado = "solicitado"
	EstadoTurnoConfirmado = "confirmado"
	EstadoTurnoCancelado  = "cancelado"
	EstadoTurnoCompletado = "completado"
)

// FranjaHoraria es una franja de media hora para los turnos de test drive. El
// identificador es la hora de inicio en formato "HH:MM" (ej. "10:00"). Ocupada
// indica que ya existe un turno activo para esa unidad, fecha y franja.
type FranjaHoraria struct {
	ID      string `json:"id"`
	Inicio  string `json:"inicio"`
	Fin     string `json:"fin"`
	Ocupada bool   `json:"ocupada,omitempty"`
}

// FranjasDisponibles devuelve el catálogo de franjas horarias de media hora en
// horario comercial: de 09:00 a 11:00 y de 14:00 a 17:00.
func FranjasDisponibles() []FranjaHoraria {
	horas := []string{
		"09:00", "09:30", "10:00", "10:30", "11:00",
		"14:00", "14:30", "15:00", "15:30", "16:00", "16:30", "17:00",
	}
	franjas := make([]FranjaHoraria, 0, len(horas))
	for _, inicio := range horas {
		inicioT, err := time.Parse("15:04", inicio)
		if err != nil {
			continue
		}
		franjas = append(franjas, FranjaHoraria{
			ID:     inicio,
			Inicio: inicio,
			Fin:    inicioT.Add(30 * time.Minute).Format("15:04"),
		})
	}
	return franjas
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
	ID                 uint      `gorm:"primaryKey" json:"id"`
	VehiculoID         uint      `gorm:"not null;index" json:"vehiculoId"`
	Vehiculo           Vehiculo  `gorm:"foreignKey:VehiculoID" json:"-"`
	ClienteID          uint      `gorm:"not null;index" json:"clienteId"`
	Cliente            Usuario   `gorm:"foreignKey:ClienteID" json:"-"`
	Fecha              string    `gorm:"type:date;not null;index" json:"fecha"`
	Franja             string    `gorm:"not null" json:"franja"`
	Estado             string    `gorm:"not null;index;default:solicitado" json:"estado"`
	BorradoPorCliente  bool      `gorm:"default:false;index" json:"borradoPorCliente"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
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
