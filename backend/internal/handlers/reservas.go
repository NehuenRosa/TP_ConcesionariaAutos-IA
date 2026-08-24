package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ReservaResumen es la ficha de una reserva para las respuestas de la API.
type ReservaResumen struct {
	ID        uint            `json:"id"`
	Vehiculo  VehiculoResumen `json:"vehiculo"`
	Cliente   UsuarioResumen  `json:"cliente"`
	Estado    string          `json:"estado"`
	MontoSenia float64        `json:"montoSenia"`
	// VencimientoComprobante es el límite (RFC3339) para subir el comprobante;
	// vacío en reservas históricas.
	VencimientoComprobante string `json:"vencimientoComprobante,omitempty"`
	// ComprobanteEnviadoAt marca cuándo se subió el comprobante; vacío si
	// sigue pendiente.
	ComprobanteEnviadoAt string `json:"comprobanteEnviadoAt,omitempty"`
	// MotivoCancelacion es la explicación del vendedor cuando anuló la
	// reserva; vacío en cancelaciones propias del cliente o expiraciones.
	MotivoCancelacion string `json:"motivoCancelacion,omitempty"`
	CreatedAt string          `json:"createdAt"`
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

// CancelarComoVendedor cancela una reserva activa con el motivo obligatorio
// que verá el cliente.
func (h *ReservaHandler) CancelarComoVendedor(c *gin.Context) {
	var entrada struct {
		Motivo string `json:"motivo"`
	}
	// El cuerpo es opcional para tolerar clientes viejos; el servicio valida
	// que el motivo no quede vacío.
	if c.Request.Body != nil {
		_ = c.ShouldBindJSON(&entrada)
	}

	reservaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de reserva inválido"})
		return
	}

	reserva, err := h.servicio.CancelarComoVendedor(c.Request.Context(), uint(reservaID), entrada.Motivo)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMotivoRequerido):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenés que indicar el motivo de la cancelación"})
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
		resumen.MontoSenia = services.CalcularMontoSenia(reserva.Vehiculo.Precio)
	}

	if !reserva.VencimientoComprobante.IsZero() {
		resumen.VencimientoComprobante = reserva.VencimientoComprobante.Format(time.RFC3339)
	}
	if reserva.ComprobanteEnviadoAt != nil {
		resumen.ComprobanteEnviadoAt = reserva.ComprobanteEnviadoAt.Format(time.RFC3339)
	}
	resumen.MotivoCancelacion = strings.TrimSpace(reserva.MotivoCancelacion)

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

// DatosTransferencia devuelve CBU/alias de la concesionaria y el monto de la
// seña (5 % del precio) del vehículo indicado, para que el cliente transfiera
// antes o después de reservar.
func (h *ReservaHandler) DatosTransferencia(c *gin.Context) {
	if _, err := extraerUsuarioID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	vehiculoID, err := strconv.ParseUint(c.Query("vehiculoId"), 10, 64)
	if err != nil || vehiculoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El vehículo es obligatorio"})
		return
	}

	datos, err := h.servicio.ObtenerDatosTransferencia(c.Request.Context(), uint(vehiculoID))
	if err != nil {
		if errors.Is(err, services.ErrVehiculoNoDisponible) {
			c.JSON(http.StatusNotFound, gin.H{"error": "El vehículo no está disponible"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los datos de transferencia"})
		return
	}

	c.JSON(http.StatusOK, datos)
}

// SubirComprobante guarda la imagen del comprobante de transferencia de una
// reserva propia activa dentro del plazo de 2 horas.
func (h *ReservaHandler) SubirComprobante(c *gin.Context) {
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

	archivo, err := c.FormFile("comprobante")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La imagen del comprobante es obligatoria"})
		return
	}
	abierta, err := archivo.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el comprobante"})
		return
	}
	defer abierta.Close()
	datos, err := io.ReadAll(abierta)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el comprobante"})
		return
	}

	reserva, err := h.servicio.SubirComprobante(c.Request.Context(), uint(reservaID), clienteID, archivo.Filename, datos)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrComprobanteInvalido):
			c.JSON(http.StatusBadRequest, gin.H{"error": "El comprobante debe ser una imagen JPG, PNG o WebP de hasta 5 MB"})
		case errors.Is(err, services.ErrReservaNoEncontrada):
			c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
		case errors.Is(err, services.ErrReservaEstadoInvalido):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrComprobanteFueraDePlazo):
			c.JSON(http.StatusConflict, gin.H{"error": "El plazo de 2 horas para enviar el comprobante venció y la reserva fue anulada"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el comprobante"})
		}
		return
	}

	c.JSON(http.StatusOK, aReservaResumen(reserva))
}

// ObtenerComprobante sirve la imagen del comprobante al dueño de la reserva y
// a vendedores/administradores.
func (h *ReservaHandler) ObtenerComprobante(c *gin.Context) {
	solicitanteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	reservaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de reserva inválido"})
		return
	}

	rol := c.GetString("rol")
	esPersonal := rol == models.RolVendedor || rol == models.RolAdministrador

	comprobante, err := h.servicio.ObtenerComprobante(c.Request.Context(), uint(reservaID), solicitanteID, esPersonal)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReservaProhibida):
			c.JSON(http.StatusForbidden, gin.H{"error": "No tenés permisos para ver este comprobante"})
		case errors.Is(err, services.ErrReservaNoEncontrada), errors.Is(err, services.ErrComprobanteNoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": "La reserva todavía no tiene un comprobante"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el comprobante"})
		}
		return
	}

	c.Data(http.StatusOK, comprobante.MIME, comprobante.Datos)
}
