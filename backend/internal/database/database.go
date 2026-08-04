package database

import (
	"fmt"
	"log"

	"concesionaria/backend/internal/config"
	"concesionaria/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Conectar abre la conexión a PostgreSQL usando GORM.
func Conectar(configuracion config.Configuracion) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		configuracion.BDHost,
		configuracion.BDPuerto,
		configuracion.BDUsuario,
		configuracion.BDPassword,
		configuracion.BDNombre,
		configuracion.BDSSL,
	)

	base, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("conexión a PostgreSQL: %w", err)
	}

	log.Println("Conexión a la base de datos establecida")
	return base, nil
}

// AutoMigrar crea o actualiza las tablas a partir de los modelos.
func AutoMigrar(base *gorm.DB) error {
	if err := base.AutoMigrate(&models.Vehiculo{}, &models.Imagen{}, &models.Usuario{}, &models.Consulta{}, &models.Mensaje{}); err != nil {
		return fmt.Errorf("auto-migración: %w", err)
	}
	return nil
}
