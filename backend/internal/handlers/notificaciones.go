package handlers

import (
	"net/http"

	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// NotificacionHandler maneja las notificaciones del sistema.
type NotificacionHandler struct {
	servicioConsultas services.ConsultaService
	servicioMensajes  services.MensajeService
}

// NuevoNotificacionHandler crea un handler de notificaciones.
func NuevoNotificacionHandler(
	servicioConsultas services.ConsultaService,
	servicioMensajes services.MensajeService,
) *NotificacionHandler {
	return &NotificacionHandler{
		servicioConsultas: servicioConsultas,
		servicioMensajes:  servicioMensajes,
	}
}

// ContadorDevuelve la cantidad total de mensajes no leídos del usuario autenticado.
func (h *NotificacionHandler) Contador(c *gin.Context) {
	usuarioID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	// Obtener IDs de consultas donde el usuario participa
	ids, err := h.servicioConsultas.ListarPorUsuario(c.Request.Context(), usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el contador"})
		return
	}

	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{"contador": 0})
		return
	}

	noLeidos, err := h.servicioMensajes.ContarNoLeidosPorConsultas(c.Request.Context(), ids, usuarioID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"contador": 0})
		return
	}

	total := 0
	for _, cantidad := range noLeidos {
		total += cantidad
	}

	c.JSON(http.StatusOK, gin.H{"contador": total})
}
