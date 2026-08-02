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
	ID          uint   `gorm:"primaryKey"`
	Marca       string `gorm:"not null"`
	Modelo      string `gorm:"not null"`
	Anio        int    `gorm:"not null"`
	Kilometraje int
	Combustible string
	Transmision string
	Precio      float64 `gorm:"not null"`
	Condicion   string  `gorm:"not null"`
	Estado      string  `gorm:"not null;index;default:disponible"`
	Imagenes    []Imagen
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName define el nombre de la tabla en español.
func (Vehiculo) TableName() string {
	return "vehiculos"
}

// Imagen es una URL asociada a un vehículo (galería del catálogo).
type Imagen struct {
	ID         uint `gorm:"primaryKey"`
	VehiculoID uint `gorm:"not null;index"`
	URL        string `gorm:"not null"`
	CreatedAt  time.Time
}

// TableName define el nombre de la tabla en español.
func (Imagen) TableName() string {
	return "imagenes"
}
