package database

import (
	"errors"

	"concesionaria/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// usuarioPorDefecto describe un usuario sembrado al arrancar el sistema.
type usuarioPorDefecto struct {
	nombre   string
	email    string
	password string
	rol      string
}

// usuariosPorDefecto son las cuentas de desarrollo creadas en el primer
// arranque para poder operar el sistema sin una pantalla de gestión de
// usuarios.
var usuariosPorDefecto = []usuarioPorDefecto{
	{nombre: "Administrador", email: "administrador@concesionaria.local", password: "Admin123!", rol: models.RolAdministrador},
	{nombre: "Vendedor", email: "vendedor@concesionaria.local", password: "Vendedor123!", rol: models.RolVendedor},
}

// SembrarUsuarios crea los usuarios por defecto si no existen (idempotente).
func SembrarUsuarios(base *gorm.DB) error {
	for _, u := range usuariosPorDefecto {
		var usuario models.Usuario
		err := base.Where("email = ?", u.email).First(&usuario).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		nuevo := models.Usuario{
			Nombre:   u.nombre,
			Email:    u.email,
			Password: string(hash),
			Rol:      u.rol,
		}
		if err := base.Create(&nuevo).Error; err != nil {
			return err
		}
	}
	return nil
}
