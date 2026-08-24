package models

import "time"

// Roles de usuario del sistema.
const (
	RolCliente       = "cliente"
	RolVendedor      = "vendedor"
	RolAdministrador = "administrador"
)

// Proveedores de identidad con los que puede crearse una cuenta.
const (
	ProveedorLocal  = "local"  // registro clásico con email y contraseña
	ProveedorGoogle = "google" // cuenta federada con Google Identity Services
)

// Usuario es la entidad de GORM que representa una cuenta del sistema.
type Usuario struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"not null" json:"nombre"`
	Email     string    `gorm:"not null;uniqueIndex" json:"email"`
	// Password guarda el hash bcrypt; cadena vacía para cuentas creadas solo
	// con Google (el login con contraseña responde 401 igualmente).
	Password string `gorm:"not null" json:"-"`
	Rol      string `gorm:"not null;index;default:cliente" json:"rol"`
	// Proveedor indica el origen de la identidad: local o google.
	Proveedor string `gorm:"not null;default:local;index" json:"proveedor"`
	// GoogleSub es el identificador estable de la cuenta de Google (claim sub).
	// Puntero nulo para cuentas locales: PostgreSQL admite múltiples NULL en un
	// índice único, así que no chocan entre sí.
	GoogleSub *string `gorm:"uniqueIndex" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TableName define el nombre de la tabla en español.
func (Usuario) TableName() string {
	return "usuarios"
}
