package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ReservaResumen es la ficha de una reserva para las respuestas de la API.
type ReservaResumen struct {
	ID       uint            `json:"id"`
	Vehiculo VehiculoResumen `json:"vehiculo"`
	Cliente  UsuarioResumen  `json:"cliente"`
	Estado   string          `json:"estado"`
	CreatedAt string         `json:"createdAt"`
}

// ReservaHandler agrupa los handlers de reservas.
type ReservaHandler struct {
	servicio services.ReservaService
}

// NuevoReservaHandler crea un handler de reservas.
func NuevoReservaHandler(servicio services.ReservaService) *ReservaHandler {
	return &ReservaHandler{servicio: servicio}
}

// Crear reserva un vehículo disponible desde el cliente autenticado.
func (h *ReservaHandler) Crear(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	var entrada struct {
		VehiculoID uint `json:"vehiculoId"`
	}
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if entrada.VehiculoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El vehículo es obligatorio"})
		return
	}

	reserva, err := h.servicio.Crear(c.Request.Context(), clienteID, entrada.VehiculoID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrVehiculoNoDisponible):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrVehiculoYaNoDisponible):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo reservar el vehículo"})
		}
		return
	}

	c.JSON(http.StatusCreated, aReservaResumen(reserva))
}

// ListarMisReservas lista las reservas del cliente autenticado.
func (h *ReservaHandler) ListarMisReservas(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	reservas, err := h.servicio.ListarMisReservas(c.Request.Context(), clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las reservas"})
		return
	}

	c.JSON(http.StatusOK, aReservasResumen(reservas))
}

// Cancelar cancela una reserva propia del cliente autenticado.
func (h *ReservaHandler) Cancelar(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	reservaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de reserva inválido"})
		return
	}

	reserva, err := h.servicio.Cancelar(c.Request.Context(), uint(reservaID), clienteID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReservaNoEncontrada):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrReservaEstadoInvalido):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo cancelar la reserva"})
		}
		return
	}

	c.JSON(http.StatusOK, aReservaResumen(reserva))
}

// Listar lista las reservas para el vendedor autenticado.
func (h *ReservaHandler) Listar(c *gin.Context) {
	estado := c.Query("estado")

	reservas, err := h.servicio.Listar(c.Request.Context(), estado)
	if err != nil {
		if errors.Is(err, services.ErrFiltroEstadoReservaInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las reservas"})
		return
	}

	c.JSON(http.StatusOK, aReservasResumen(reservas))
}

// ConfirmarVenta confirma la venta de una reserva activa.
func (h *ReservaHandler) ConfirmarVenta(c *gin.Context) {
	h.cambiarEstado(c, func(id uint) (*models.Reserva, error) {
		return h.servicio.ConfirmarVenta(c.Request.Context(), id)
	})
}

// CancelarComoVendedor cancela una reserva activa.
func (h *ReservaHandler) CancelarComoVendedor(c *gin.Context) {
	h.cambiarEstado(c, func(id uint) (*models.Reserva, error) {
		return h.servicio.CancelarComoVendedor(c.Request.Context(), id)
	})
}

// cambiarEstado comparte la lógica de transición de estado de las reservas.
func (h *ReservaHandler) cambiarEstado(c *gin.Context, accion func(id uint) (*models.Reserva, error)) {
	reservaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de reserva inválido"})
		return
	}

	reserva, err := accion(uint(reservaID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReservaNoEncontrada):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrReservaEstadoInvalido):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la reserva"})
		}
		return
	}

	c.JSON(http.StatusOK, aReservaResumen(reserva))
}

// aReservaResumen convierte un modelo en el resumen para las respuestas.
func aReservaResumen(reserva *models.Reserva) ReservaResumen {
	resumen := ReservaResumen{
		ID:        reserva.ID,
		Estado:    reserva.Estado,
		CreatedAt: reserva.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if reserva.Vehiculo.ID != 0 {
		imagen := ""
		if len(reserva.Vehiculo.Imagenes) > 0 {
			imagen = reserva.Vehiculo.Imagenes[0].URL
		}
		resumen.Vehiculo = VehiculoResumen{
			ID:        reserva.Vehiculo.ID,
			Marca:     reserva.Vehiculo.Marca,
			Modelo:    reserva.Vehiculo.Modelo,
			Anio:      reserva.Vehiculo.Anio,
			Precio:    reserva.Vehiculo.Precio,
			Condicion: reserva.Vehiculo.Condicion,
			Tipo:      reserva.Vehiculo.Tipo,
			Imagen:    imagen,
		}
	}

	resumen.Cliente = UsuarioResumen{
		ID:     reserva.Cliente.ID,
		Nombre: reserva.Cliente.Nombre,
	}

	return resumen
}

// aReservasResumen convierte un listado de modelos en resúmenes.
func aReservasResumen(reservas []models.Reserva) []ReservaResumen {
	resumenes := make([]ReservaResumen, 0, len(reservas))
	for i := range reservas {
		resumenes = append(resumenes, aReservaResumen(&reservas[i]))
	}
	return resumenes
}
