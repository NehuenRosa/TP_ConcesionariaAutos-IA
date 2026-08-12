package database

import (
	"fmt"
	"log"
	"strings"

	"concesionaria/backend/internal/config"
	"concesionaria/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Conectar abre la conexión a PostgreSQL usando GORM. Si se define BD_URL
// (cadena completa, p. ej. en Render) se usa directamente; en caso contrario
// se arma el DSN a partir de las variables BD_*.
func Conectar(configuracion config.Configuracion) (*gorm.DB, error) {
	dsn := configuracion.BDURL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			configuracion.BDHost,
			configuracion.BDPuerto,
			configuracion.BDUsuario,
			configuracion.BDPassword,
			configuracion.BDNombre,
			configuracion.BDSSL,
		)
	} else if !strings.Contains(dsn, "sslmode=") {
		// Render exige SSL en Postgres: lo agregamos si la URL no lo indica.
		separador := "?"
		if strings.Contains(dsn, "?") {
			separador = "&"
		}
		dsn += separador + "sslmode=require"
	}

	base, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("conexión a PostgreSQL: %w", err)
	}

	log.Println("Conexión a la base de datos establecida")
	return base, nil
}

// AutoMigrar crea o actualiza las tablas a partir de los modelos.
func AutoMigrar(base *gorm.DB) error {
	if err := base.AutoMigrate(&models.Vehiculo{}, &models.Imagen{}, &models.Usuario{}, &models.Consulta{}, &models.Mensaje{}, &models.TurnoTestDrive{}); err != nil {
		return fmt.Errorf("auto-migración: %w", err)
	}
	return nil
}
