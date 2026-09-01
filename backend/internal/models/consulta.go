package models

import "time"

// Estado de una consulta en el sistema.
const (
	EstadoPendiente        = "pendiente"
	EstadoEnConversacion   = "en_conversacion"
	EstadoCerrada          = "cerrada"
)

// Consulta representa una conversación entre un cliente y un vendedor
// sobre un vehículo específico.
type Consulta struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VehiculoID uint      `gorm:"not null;index" json:"vehiculoId"`
	Vehiculo   Vehiculo  `gorm:"foreignKey:VehiculoID" json:"vehiculo"`
	ClienteID  uint      `gorm:"not null;index" json:"clienteId"`
	Cliente    Usuario   `gorm:"foreignKey:ClienteID" json:"cliente"`
	VendedorID *uint     `gorm:"index" json:"vendedorId"`
	Vendedor   *Usuario  `gorm:"foreignKey:VendedorID" json:"vendedor,omitempty"`
	Estado     string    `gorm:"type:varchar(20);not null;default:pendiente;index" json:"estado"`
	Mensajes   []Mensaje `gorm:"foreignKey:ConsultaID" json:"mensajes,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// TableName define el nombre de la tabla en español.
func (Consulta) TableName() string {
	return "consultas"
}

// Mensaje es un mensaje dentro de una consulta/conversación.
type Mensaje struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// El índice compuesto (consulta_id, created_at) acelera la carga del
	// historial de un hilo ordenado cronológicamente y el fetch incremental
	// por desdeID.
	ConsultaID  uint      `gorm:"not null;index:idx_mensajes_consulta_hilo,priority:1" json:"consultaId"`
	RemitenteID uint      `gorm:"not null" json:"remitenteId"`
	Remitente   Usuario   `gorm:"foreignKey:RemitenteID" json:"remitente"`
	Contenido   string    `gorm:"type:text;not null" json:"contenido"`
	Leido       bool      `gorm:"default:false" json:"leido"`
	CreatedAt   time.Time `gorm:"index:idx_mensajes_consulta_hilo,priority:2" json:"createdAt"`
}

// TableName define el nombre de la tabla en español.
func (Mensaje) TableName() string {
	return "mensajes"
}
