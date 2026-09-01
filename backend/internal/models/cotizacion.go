package models

import "time"

// Estado de una cotización conversada con la IA.
const (
	// EstadoCotizacionAbierta indica que la conversación de cotización sigue
	// activa y se pueden seguir enviando mensajes.
	EstadoCotizacionAbierta = "abierta"
	// EstadoCotizacionCerrada indica que la cotización quedó cerrada.
	EstadoCotizacionCerrada = "cerrada"
)

// Remitente de un mensaje dentro de una cotización.
const (
	// RemitenteCliente es el mensaje enviado por el cliente.
	RemitenteCliente = "cliente"
	// RemitenteIA es el mensaje generado por el asistente.
	RemitenteIA = "ia"
	// RemitenteVendedor es el mensaje enviado por un vendedor desde la bandeja
	// de cotizaciones.
	RemitenteVendedor = "vendedor"
)

// Cotizacion representa una conversación entre un cliente y el asistente IA
// sobre precios, financiación y formas de pago de un vehículo. El contenido de
// los mensajes se guarda cifrado en la base porque son datos sensibles. Un
// vendedor puede tomar la conversación para atenderla personalmente: desde ese
// momento la IA deja de responder y los mensajes del personal quedan firmados
// como "vendedor".
type Cotizacion struct {
	ID         uint                `gorm:"primaryKey" json:"id"`
	VehiculoID uint                `gorm:"not null;index" json:"vehiculoId"`
	Vehiculo   Vehiculo            `gorm:"foreignKey:VehiculoID" json:"vehiculo"`
	ClienteID  uint                `gorm:"not null;index" json:"clienteId"`
	Cliente    Usuario             `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
	VendedorID *uint               `gorm:"index" json:"vendedorId,omitempty"`
	Vendedor   *Usuario            `gorm:"foreignKey:VendedorID" json:"vendedor,omitempty"`
	FechaToma  *time.Time          `json:"fechaToma,omitempty"`
	Estado     string              `gorm:"type:varchar(20);not null;default:abierta;index" json:"estado"`
	Mensajes   []MensajeCotizacion `gorm:"foreignKey:CotizacionID" json:"mensajes,omitempty"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
}

// TableName define el nombre de la tabla en español.
func (Cotizacion) TableName() string {
	return "cotizaciones"
}

// MensajeCotizacion es un mensaje dentro de una cotización. El campo Contenido
// guarda el texto cifrado en la base: no se persiste en claro. Las marcas
// LeidoPorCliente y LeidoPorVendedor alimentan las notificaciones de mensajes
// sin leer de cada lado.
type MensajeCotizacion struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// El índice compuesto (cotizacion_id, created_at) acelera la carga del
	// historial de un hilo ordenado cronológicamente y el fetch incremental
	// por desdeID (CU-12, ver docs/roadmap.md "Escalabilidad de conversaciones").
	CotizacionID     uint      `gorm:"not null;index:idx_cotizacion_mensajes_hilo,priority:1" json:"cotizacionId"`
	Remitente        string    `gorm:"size:20;not null" json:"remitente"`
	Contenido        string    `gorm:"type:text;not null" json:"-"` // cifrado en reposo
	LeidoPorCliente  bool      `gorm:"default:false" json:"leidoPorCliente"`
	LeidoPorVendedor bool      `gorm:"default:false" json:"leidoPorVendedor"`
	CreatedAt        time.Time `gorm:"index:idx_cotizacion_mensajes_hilo,priority:2" json:"createdAt"`
}

// TableName define el nombre de la tabla en español.
func (MensajeCotizacion) TableName() string {
	return "cotizacion_mensajes"
}