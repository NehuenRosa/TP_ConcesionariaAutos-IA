package models

import "time"

// Roles de usuario del sistema.
const (
	RolCliente       = "cliente"
	RolVendedor      = "vendedor"
	RolAdministrador = "administrador"
)

// Usuario es la entidad de GORM que representa una cuenta del sistema.
type Usuario struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"not null" json:"nombre"`
	Email     string    `gorm:"not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Rol       string    `gorm:"not null;index;default:cliente" json:"rol"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TableName define el nombre de la tabla en español.
func (Usuario) TableName() string {
	return "usuarios"
}
