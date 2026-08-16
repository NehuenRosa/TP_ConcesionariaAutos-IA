package main

import (
	"log/slog"
	"time"

	"concesionaria/backend/internal/config"
	"concesionaria/backend/internal/database"
	"concesionaria/backend/internal/router"

	"gorm.io/gorm"
)

// cantidadIntentosConexion y esperaEntreIntentos controlan el reintento de
// conexión inicial a PostgreSQL (útil cuando el contenedor de BD arranca más
// tarde que la API).
const (
	cantidadIntentosConexion = 10
	esperaEntreIntentos      = 2 * time.Second
)

func main() {
	configuracion := config.Cargar()

	base, err := conectarConReintentos(configuracion)
	if err != nil {
		slog.Error("No se pudo conectar a la base de datos", "error", err.Error())
		return
	}

	if err := database.AutoMigrar(base); err != nil {
		slog.Error("No se pudieron aplicar las migraciones", "error", err.Error())
		return
	}

	if err := database.SembrarVehiculos(base); err != nil {
		slog.Error("No se pudieron sembrar los vehículos por defecto", "error", err.Error())
		return
	}

	if err := database.SembrarUsuarios(base); err != nil {
		slog.Error("No se pudieron sembrar los usuarios por defecto", "error", err.Error())
		return
	}

	enrutador := router.Nuevo(base, configuracion)

	slog.Info("API escuchando", "direccion", configuracion.Host+":"+configuracion.Puerto)
	if err := enrutador.Run(configuracion.Host + ":" + configuracion.Puerto); err != nil {
		slog.Error("Error al iniciar el servidor", "error", err.Error())
	}
}

// conectarConReintentos intenta conectar a PostgreSQL varias veces antes de
// rendirse, para tolerar arranques desordenados con Docker Compose.
func conectarConReintentos(configuracion config.Configuracion) (*gorm.DB, error) {
	var base *gorm.DB
	var err error

	for intento := 1; intento <= cantidadIntentosConexion; intento++ {
		base, err = database.Conectar(configuracion)
		if err == nil {
			return base, nil
		}

		slog.Warn(
			"No se pudo conectar a la base de datos, reintentando",
			"intento", intento,
			"de", cantidadIntentosConexion,
			"error", err.Error(),
		)
		time.Sleep(esperaEntreIntentos)
	}

	return nil, err
}
