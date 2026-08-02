package main

import (
	"log/slog"

	"concesionaria/backend/internal/config"
	"concesionaria/backend/internal/database"
	"concesionaria/backend/internal/router"
)

func main() {
	configuracion := config.Cargar()

	base, err := database.Conectar(configuracion)
	if err != nil {
		slog.Error("No se pudo conectar a la base de datos", "error", err.Error())
		return
	}

	if err := database.AutoMigrar(base); err != nil {
		slog.Error("No se pudieron aplicar las migraciones", "error", err.Error())
		return
	}

	enrutador := router.Nuevo(base, configuracion)

	slog.Info("API escuchando", "direccion", configuracion.Host+":"+configuracion.Puerto)
	if err := enrutador.Run(configuracion.Host + ":" + configuracion.Puerto); err != nil {
		slog.Error("Error al iniciar el servidor", "error", err.Error())
	}
}
