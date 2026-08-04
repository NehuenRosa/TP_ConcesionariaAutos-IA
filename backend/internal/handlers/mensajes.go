package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// MensajeResumenCompleto es la información completa de un mensaje.
type MensajeResumenCompleto struct {
	ID          uint            `json:"id"`
	ConsultaID  uint            `json:"consultaId"`
	Remitente   UsuarioResumen  `json:"remitente"`
	Contenido   string          `json:"contenido"`
	Leido       bool            `json:"leido"`
	CreatedAt   string          `json:"createdAt"`
}

// MensajeHandler agrupa los handlers de mensajes.
type MensajeHandler struct {
	servicio services.MensajeService
}

// NuevoMensajeHandler crea un handler de mensajes.
func NuevoMensajeHandler(servicio services.MensajeService) *MensajeHandler {
	return &MensajeHandler{servicio: servicio}
}

// Enviar envía un mensaje en una consulta.
func (h *MensajeHandler) Enviar(c *gin.Context) {
	usuarioID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de consulta inválido"})
		return
	}

	var entrada struct {
		Contenido string `json:"contenido"`
	}
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	mensaje, err := h.servicio.Enviar(c.Request.Context(), uint(consultaID), usuarioID, entrada.Contenido)
	if err != nil {
		if errors.Is(err, services.ErrMensajeVacio) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrConsultaNoEncontrada) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrConsultaCerradaMensajes) || errors.Is(err, services.ErrNoEsParticipante) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar el mensaje"})
		return
	}

	c.JSON(http.StatusCreated, aMensajeResumenCompleto(mensaje))
}

// ObtenerMensajes obtiene todos los mensajes de una consulta.
func (h *MensajeHandler) ObtenerMensajes(c *gin.Context) {
	usuarioID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de consulta inválido"})
		return
	}

	mensajes, err := h.servicio.ObtenerPorConsulta(c.Request.Context(), uint(consultaID), usuarioID)
	if err != nil {
		if errors.Is(err, services.ErrNoEsParticipante) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los mensajes"})
		return
	}

	resumenes := make([]MensajeResumenCompleto, 0, len(mensajes))
	for _, mensaje := range mensajes {
		resumenes = append(resumenes, aMensajeResumenCompleto(&mensaje))
	}

	c.JSON(http.StatusOK, resumenes)
}

// ObtenerNuevos obtiene los mensajes nuevos desde un timestamp.
func (h *MensajeHandler) ObtenerNuevos(c *gin.Context) {
	usuarioID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de consulta inválido"})
		return
	}

	desdeStr := c.Query("desde")
	var desde time.Time
	if desdeStr != "" {
		desde, err = time.Parse(time.RFC3339, desdeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de timestamp inválido"})
			return
		}
	} else {
		desde = time.Time{}
	}

	mensajes, err := h.servicio.ObtenerNuevos(c.Request.Context(), uint(consultaID), usuarioID, desde)
	if err != nil {
		if errors.Is(err, services.ErrNoEsParticipante) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los mensajes nuevos"})
		return
	}

	resumenes := make([]MensajeResumenCompleto, 0, len(mensajes))
	for _, mensaje := range mensajes {
		resumenes = append(resumenes, aMensajeResumenCompleto(&mensaje))
	}

	c.JSON(http.StatusOK, resumenes)
}

// MarcarLeidos marca como leídos los mensajes de otros en una consulta.
func (h *MensajeHandler) MarcarLeidos(c *gin.Context) {
	usuarioID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de consulta inválido"})
		return
	}

	err = h.servicio.MarcarComoLeidos(c.Request.Context(), uint(consultaID), usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron marcar los mensajes como leídos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Mensajes marcados como leídos"})
}

// aMensajeResumenCompleto convierte un modelo en el resumen completo.
func aMensajeResumenCompleto(mensaje *models.Mensaje) MensajeResumenCompleto {
	return MensajeResumenCompleto{
		ID:         mensaje.ID,
		ConsultaID: mensaje.ConsultaID,
		Remitente: UsuarioResumen{
			ID:     mensaje.Remitente.ID,
			Nombre: mensaje.Remitente.Nombre,
		},
		Contenido: mensaje.Contenido,
		Leido:     mensaje.Leido,
		CreatedAt: mensaje.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
