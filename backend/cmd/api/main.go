package main

import (
	"context"
	"log/slog"
	"time"

	"concesionaria/backend/internal/config"
	"concesionaria/backend/internal/database"
	"concesionaria/backend/internal/repositories"
	"concesionaria/backend/internal/router"
	"concesionaria/backend/internal/services"

	"gorm.io/gorm"
)

// cantidadIntentosConexion y esperaEntreIntentos controlan el reintento de
// conexión inicial a PostgreSQL (útil cuando el contenedor de BD arranca más
// tarde que la API).
const (
	cantidadIntentosConexion = 10
	esperaEntreIntentos      = 2 * time.Second
	// intervaloExpiracionReservas define cada cuánto el job interno anula las
	// reservas activas que vencieron las 2 horas sin comprobante (CU-08).
	intervaloExpiracionReservas = 30 * time.Second
	// intervaloRetencionConversaciones define cada cuánto el job interno purga
	// las conversaciones cerradas que superan el plazo de conservación
	// (RETENCION_CONVERSACIONES_DIAS).
	intervaloRetencionConversaciones = 1 * time.Hour
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

	// Job de expiración de reservas (CU-08): anula las activas que vencieron
	// sin comprobante y libera sus unidades. Corre una vez al arrancar (barre
	// rezagos de caídas) y luego cada 30 segundos.
	servicioReservas := services.NuevoReservaService(
		repositories.NuevoReservaRepository(base),
		repositories.NuevoVehiculoRepository(base),
		configuracion.CbuConcesionaria,
		configuracion.AliasConcesionaria,
	)
	go ejecutarExpiracionReservas(servicioReservas)

	// Job de retención de conversaciones: purga las cotizaciones y consultas
	// cerradas que superan el plazo de conservación (privacidad y control del
	// crecimiento de la base). Corre una vez al arrancar y luego cada hora.
	servicioRetencion := services.NuevoRetencionService(repositories.NuevoRetencionRepository(base))
	go ejecutarRetencionConversaciones(servicioRetencion, configuracion.RetencionConversacionesDias)

	slog.Info("API escuchando", "direccion", configuracion.Host+":"+configuracion.Puerto)
	if err := enrutador.Run(configuracion.Host + ":" + configuracion.Puerto); err != nil {
		slog.Error("Error al iniciar el servidor", "error", err.Error())
	}
}

// ejecutarExpiracionReservas corre el barrido periódico de reservas vencidas
// sin comprobante.
func ejecutarExpiracionReservas(servicioReservas services.ReservaService) {
	if err := servicioReservas.ExpirarVencidas(context.Background()); err != nil {
		slog.Warn("expiración inicial de reservas falló", "error", err.Error())
	}

	ticker := time.NewTicker(intervaloExpiracionReservas)
	defer ticker.Stop()
	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := servicioReservas.ExpirarVencidas(ctx); err != nil {
			slog.Warn("expiración de reservas falló", "error", err.Error())
		}
		cancel()
	}
}

// ejecutarRetencionConversaciones corre el barrido periódico que purga las
// conversaciones cerradas que superan el plazo de conservación.
func ejecutarRetencionConversaciones(servicioRetencion services.RetencionService, dias int) {
	if err := ejecutarUnaRetencion(servicioRetencion, dias); err != nil {
		slog.Warn("retención inicial de conversaciones falló", "error", err.Error())
	}

	ticker := time.NewTicker(intervaloRetencionConversaciones)
	defer ticker.Stop()
	for range ticker.C {
		if err := ejecutarUnaRetencion(servicioRetencion, dias); err != nil {
			slog.Warn("retención de conversaciones falló", "error", err.Error())
		}
	}
}

// ejecutarUnaRetencion ejecuta una pasada de la purga y loguea el resultado.
func ejecutarUnaRetencion(servicioRetencion services.RetencionService, dias int) error {
	ctx, cancelar := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelar()

	resultado, err := servicioRetencion.Ejecutar(ctx, dias)
	if err != nil {
		return err
	}
	if resultado.Cotizaciones+resultado.Consultas > 0 {
		slog.Info("retención de conversaciones ejecutada",
			"dias", dias,
			"cotizaciones", resultado.Cotizaciones,
			"mensajesCotizaciones", resultado.MensajesCotizaciones,
			"consultas", resultado.Consultas,
			"mensajesConsultas", resultado.MensajesConsultas,
		)
	}
	return nil
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
