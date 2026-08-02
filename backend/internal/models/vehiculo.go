package models

import "time"

// Estado del vehículo en el concesionario.
const (
	EstadoDisponible  = "disponible"
	EstadoReservado   = "reservado"
	EstadoVendido     = "vendido"
	EstadoDadoDeBaja  = "dado_de_baja"
)

// Condicion del vehículo: nuevo o usado.
const (
	CondicionNuevo = "nuevo"
	CondicionUsado = "usado"
)

// Vehiculo es la entidad de GORM que representa una unidad del concesionario.
type Vehiculo struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Marca       string  `gorm:"not null" json:"marca"`
	Modelo      string  `gorm:"not null" json:"modelo"`
	Anio        int     `gorm:"not null" json:"anio"`
	Kilometraje int     `json:"kilometraje"`
	Combustible string  `json:"combustible"`
	Transmision string  `json:"transmision"`
	Precio      float64 `gorm:"not null" json:"precio"`
	Condicion   string  `gorm:"not null" json:"condicion"`
	Estado      string  `gorm:"not null;index;default:disponible" json:"estado"`
	Imagenes    []Imagen `json:"imagenes"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// TableName define el nombre de la tabla en español.
func (Vehiculo) TableName() string {
	return "vehiculos"
}

// Imagen es una URL asociada a un vehículo (galería del catálogo).
type Imagen struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	VehiculoID uint   `gorm:"not null;index" json:"-"`
	URL        string `gorm:"not null" json:"url"`
	CreatedAt  time.Time `json:"-"`
}

// TableName define el nombre de la tabla en español.
func (Imagen) TableName() string {
	return "imagenes"
}
