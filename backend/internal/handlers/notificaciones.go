package handlers

import (
	"log/slog"
	"net/http"

	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// NotificacionHandler maneja las notificaciones del sistema.
type NotificacionHandler struct {
	servicioConsultas    services.ConsultaService
	servicioMensajes     services.MensajeService
	servicioCotizaciones services.CotizacionService
}

// NuevoNotificacionHandler crea un handler de notificaciones.
func NuevoNotificacionHandler(
	servicioConsultas services.ConsultaService,
	servicioMensajes services.MensajeService,
	servicioCotizaciones services.CotizacionService,
) *NotificacionHandler {
	return &NotificacionHandler{
		servicioConsultas:    servicioConsultas,
		servicioMensajes:     servicioMensajes,
		servicioCotizaciones: servicioCotizaciones,
	}
}

// ContadorDevuelve la cantidad total de mensajes no leídos del usuario
// autenticado, con el desglose por canal: consultas y cotizaciones.
func (h *NotificacionHandler) Contador(c *gin.Context) {
	usuarioID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	// Mensajes no leídos en consultas (comportamiento existente).
	noLeidosConsultas := 0
	ids, err := h.servicioConsultas.ListarPorUsuario(c.Request.Context(), usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el contador"})
		return
	}

	if len(ids) > 0 {
		noLeidos, err := h.servicioMensajes.ContarNoLeidosPorConsultas(c.Request.Context(), ids, usuarioID)
		if err != nil {
			slog.Warn("no se pudieron contar los no leídos de consultas", "error", err)
		} else {
			for _, cantidad := range noLeidos {
				noLeidosConsultas += cantidad
			}
		}
	}

	// Mensajes no leídos en cotizaciones según el rol. Ante error se degrada a
	// cero sin fallar el request.
	rol, _ := c.Get("rol")
	rolActual, _ := rol.(string)
	noLeidosCotizaciones, err := h.servicioCotizaciones.ContarNoLeidos(c.Request.Context(), usuarioID, rolActual)
	if err != nil {
		slog.Warn("no se pudieron contar los no leídos de cotizaciones", "error", err)
		noLeidosCotizaciones = 0
	}

	total := noLeidosConsultas + int(noLeidosCotizaciones)

	c.JSON(http.StatusOK, gin.H{
		"contador":      total,
		"consultas":     noLeidosConsultas,
		"cotizaciones":  noLeidosCotizaciones,
	})
}
