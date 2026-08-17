package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// CotizacionResumen es la ficha resumida de una cotización para listados.
type CotizacionResumen struct {
	ID            uint            `json:"id"`
	Vehiculo      VehiculoResumen `json:"vehiculo"`
	Estado        string          `json:"estado"`
	UltimoMensaje *MensajeResumen `json:"ultimoMensaje,omitempty"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
}

// MensajeCotizacionResumen es un mensaje dentro de una cotización.
type MensajeCotizacionResumen struct {
	ID        uint   `json:"id"`
	Remitente string `json:"remitente"`
	Contenido string `json:"contenido"`
	CreatedAt string `json:"createdAt"`
}

// CotizacionDetalle es la ficha completa de una cotización con sus mensajes.
type CotizacionDetalle struct {
	ID        uint                       `json:"id"`
	Vehiculo  VehiculoResumen            `json:"vehiculo"`
	Estado    string                     `json:"estado"`
	Mensajes  []MensajeCotizacionResumen `json:"mensajes"`
	CreatedAt string                     `json:"createdAt"`
	UpdatedAt string                     `json:"updatedAt"`
}

// CotizacionHandler agrupa los handlers del panel de cotizaciones del cliente.
type CotizacionHandler struct {
	servicio services.CotizacionService
}

// NuevoCotizacionHandler crea un handler de cotizaciones.
func NuevoCotizacionHandler(servicio services.CotizacionService) *CotizacionHandler {
	return &CotizacionHandler{servicio: servicio}
}

// Crear crea una cotización para un vehículo del cliente autenticado, con un
// mensaje inicial opcional y la primera respuesta del asistente.
func (h *CotizacionHandler) Crear(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	var entrada struct {
		VehiculoID uint   `json:"vehiculoId"`
		Mensaje    string `json:"mensaje,omitempty"`
	}
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida"})
		return
	}
	if entrada.VehiculoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El vehículo es obligatorio"})
		return
	}

	cotizacion, err := h.servicio.Crear(c.Request.Context(), clienteID, entrada.VehiculoID, entrada.Mensaje)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrVehiculoNoDisponible):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la cotización"})
		}
		return
	}

	c.JSON(http.StatusCreated, aCotizacionDetalle(cotizacion))
}

// ListarMisCotizaciones lista las cotizaciones del cliente autenticado con su
// último mensaje descifrado.
func (h *CotizacionHandler) ListarMisCotizaciones(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	cotizaciones, err := h.servicio.ListarPorCliente(c.Request.Context(), clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las cotizaciones"})
		return
	}

	resumenes := make([]CotizacionResumen, 0, len(cotizaciones))
	for _, cotizacion := range cotizaciones {
		resumenes = append(resumenes, aCotizacionResumen(&cotizacion))
	}

	c.JSON(http.StatusOK, resumenes)
}

// Obtener responde el detalle de una cotización del cliente autenticado con
// todos sus mensajes descifrados.
func (h *CotizacionHandler) Obtener(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	cotizacionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de cotización inválido"})
		return
	}

	cotizacion, err := h.servicio.ObtenerPorCliente(c.Request.Context(), clienteID, uint(cotizacionID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCotizacionNoEncontrada):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrCotizacionNoPertenece):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la cotización"})
		}
		return
	}

	c.JSON(http.StatusOK, aCotizacionDetalle(cotizacion))
}

// peticionMensajeCotizacion es el cuerpo de POST /api/cotizaciones/:id/mensajes.
type peticionMensajeCotizacion struct {
	Mensaje string `json:"mensaje"`
}

// EnviarMensaje guarda el mensaje del cliente, genera la respuesta del
// asistente con el contexto del vehículo y devuelve el texto de la respuesta.
func (h *CotizacionHandler) EnviarMensaje(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	cotizacionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de cotización inválido"})
		return
	}

	var peticion peticionMensajeCotizacion
	if err := c.ShouldBindJSON(&peticion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida"})
		return
	}

	respuesta, err := h.servicio.EnviarMensaje(c.Request.Context(), clienteID, uint(cotizacionID), peticion.Mensaje)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCotizacionNoEncontrada):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrCotizacionNoPertenece):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrCotizacionCerradaMensajes):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrMensajeVacio),
			errors.Is(err, services.ErrMensajeMuyLargo):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo procesar el mensaje"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"respuesta": respuesta, "enviado": true})
}

// Cerrar cierra una cotización abierta del cliente autenticado.
func (h *CotizacionHandler) Cerrar(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	cotizacionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de cotización inválido"})
		return
	}

	cotizacion, err := h.servicio.Cerrar(c.Request.Context(), clienteID, uint(cotizacionID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCotizacionNoEncontrada):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrCotizacionNoPertenece):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrCotizacionYaCerrada):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo cerrar la cotización"})
		}
		return
	}

	c.JSON(http.StatusOK, aCotizacionDetalle(cotizacion))
}

// aCotizacionResumen convierte un modelo en el resumen para el listado. El
// último mensaje llega ya descifrado del servicio.
func aCotizacionResumen(cotizacion *models.Cotizacion) CotizacionResumen {
	resumen := CotizacionResumen{
		ID:        cotizacion.ID,
		Vehiculo:  aResumen(cotizacion.Vehiculo),
		Estado:    cotizacion.Estado,
		CreatedAt: cotizacion.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: cotizacion.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if len(cotizacion.Mensajes) > 0 {
		ultimo := cotizacion.Mensajes[len(cotizacion.Mensajes)-1]
		resumen.UltimoMensaje = &MensajeResumen{
			Contenido: ultimo.Contenido,
			CreatedAt: ultimo.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return resumen
}

// aCotizacionDetalle convierte un modelo en el detalle con sus mensajes. Los
// mensajes llegan ya descifrados del servicio.
func aCotizacionDetalle(cotizacion *models.Cotizacion) CotizacionDetalle {
	mensajes := make([]MensajeCotizacionResumen, 0, len(cotizacion.Mensajes))
	for _, mensaje := range cotizacion.Mensajes {
		mensajes = append(mensajes, MensajeCotizacionResumen{
			ID:        mensaje.ID,
			Remitente: mensaje.Remitente,
			Contenido: mensaje.Contenido,
			CreatedAt: mensaje.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return CotizacionDetalle{
		ID:        cotizacion.ID,
		Vehiculo:  aResumen(cotizacion.Vehiculo),
		Estado:    cotizacion.Estado,
		Mensajes:  mensajes,
		CreatedAt: cotizacion.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: cotizacion.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
