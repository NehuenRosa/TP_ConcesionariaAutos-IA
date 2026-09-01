package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// TurnoTestDriveResumen es la ficha de un turno de test drive para las
// respuestas de la API.
type TurnoTestDriveResumen struct {
	ID       uint             `json:"id"`
	Vehiculo VehiculoResumen  `json:"vehiculo"`
	Cliente  UsuarioResumen   `json:"cliente"`
	Fecha    string           `json:"fecha"`
	Franja   string           `json:"franja"`
	Estado   string           `json:"estado"`
}

// TurnoTestDriveHandler agrupa los handlers de turnos de test drive.
type TurnoTestDriveHandler struct {
	servicio services.TurnoTestDriveService
}

// NuevoTurnoTestDriveHandler crea un handler de turnos de test drive.
func NuevoTurnoTestDriveHandler(servicio services.TurnoTestDriveService) *TurnoTestDriveHandler {
	return &TurnoTestDriveHandler{servicio: servicio}
}

// Franjas responde el catálogo público de franjas horarias. Si llegan los
// query params vehiculoId y fecha, marca cuáles ya están ocupadas.
func (h *TurnoTestDriveHandler) Franjas(c *gin.Context) {
	vehiculoID, _ := strconv.ParseUint(c.Query("vehiculoId"), 10, 64)
	fecha := c.Query("fecha")

	franjas, err := h.servicio.FranjasConDisponibilidad(c.Request.Context(), uint(vehiculoID), fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las franjas"})
		return
	}
	c.JSON(http.StatusOK, franjas)
}

// Solicitar crea un turno solicitado desde el cliente autenticado.
func (h *TurnoTestDriveHandler) Solicitar(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	var entrada struct {
		VehiculoID uint   `json:"vehiculoId"`
		Fecha      string `json:"fecha"`
		Franja     string `json:"franja"`
	}
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if entrada.VehiculoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El vehículo es obligatorio"})
		return
	}

	turno, err := h.servicio.Solicitar(c.Request.Context(), clienteID, entrada.VehiculoID, entrada.Fecha, entrada.Franja)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDatosTurnoInvalidos):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTurnoEnPasado):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrVehiculoNoDisponible):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTurnoSuperpuesto):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTurnoDuplicadoVehiculo):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo solicitar el test drive"})
		}
		return
	}

	c.JSON(http.StatusCreated, aTurnoResumen(turno))
}

// ListarMisTurnos lista los turnos del cliente autenticado.
func (h *TurnoTestDriveHandler) ListarMisTurnos(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	turnos, err := h.servicio.ListarMisTurnos(c.Request.Context(), clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los turnos"})
		return
	}

	c.JSON(http.StatusOK, aTurnosResumen(turnos))
}

// Cancelar cancela un turno propio del cliente autenticado.
func (h *TurnoTestDriveHandler) Cancelar(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	turnoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de turno inválido"})
		return
	}

	turno, err := h.servicio.Cancelar(c.Request.Context(), uint(turnoID), clienteID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTurnoNoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTurnoEstadoInvalido):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo cancelar el turno"})
		}
		return
	}

	c.JSON(http.StatusOK, aTurnoResumen(turno))
}

// Eliminar borra lógicamente un turno propio para que el cliente no lo vea en
// su listado. Los turnos activos se cancelan antes de ocultarlos.
func (h *TurnoTestDriveHandler) Eliminar(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	turnoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de turno inválido"})
		return
	}

	turno, err := h.servicio.Eliminar(c.Request.Context(), uint(turnoID), clienteID)
	if err != nil {
		if errors.Is(err, services.ErrTurnoNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el turno"})
		return
	}

	c.JSON(http.StatusOK, aTurnoResumen(turno))
}
func (h *TurnoTestDriveHandler) Listar(c *gin.Context) {
	estado := c.Query("estado")

	turnos, err := h.servicio.Listar(c.Request.Context(), estado)
	if err != nil {
		if errors.Is(err, services.ErrFiltroEstadoTurnoInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los turnos"})
		return
	}

	c.JSON(http.StatusOK, aTurnosResumen(turnos))
}

// Confirmar confirma un turno solicitado.
func (h *TurnoTestDriveHandler) Confirmar(c *gin.Context) {
	h.cambiarEstado(c, func(id uint) (*models.TurnoTestDrive, error) {
		return h.servicio.Confirmar(c.Request.Context(), id)
	})
}

// CancelarComoVendedor cancela un turno solicitado o confirmado.
func (h *TurnoTestDriveHandler) CancelarComoVendedor(c *gin.Context) {
	h.cambiarEstado(c, func(id uint) (*models.TurnoTestDrive, error) {
		return h.servicio.CancelarComoVendedor(c.Request.Context(), id)
	})
}

// Completar marca como completado un turno confirmado.
func (h *TurnoTestDriveHandler) Completar(c *gin.Context) {
	h.cambiarEstado(c, func(id uint) (*models.TurnoTestDrive, error) {
		return h.servicio.Completar(c.Request.Context(), id)
	})
}

// cambiarEstado comparte la lógica de transición de estado de los turnos.
func (h *TurnoTestDriveHandler) cambiarEstado(c *gin.Context, accion func(id uint) (*models.TurnoTestDrive, error)) {
	turnoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de turno inválido"})
		return
	}

	turno, err := accion(uint(turnoID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTurnoNoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrTurnoEstadoInvalido):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el turno"})
		}
		return
	}

	c.JSON(http.StatusOK, aTurnoResumen(turno))
}

// aTurnoResumen convierte un modelo en el resumen para las respuestas.
func aTurnoResumen(turno *models.TurnoTestDrive) TurnoTestDriveResumen {
	resumen := TurnoTestDriveResumen{
		ID:     turno.ID,
		Fecha:  recortarFecha(turno.Fecha),
		Franja: turno.Franja,
		Estado: turno.Estado,
	}

	if turno.Vehiculo.ID != 0 {
		imagen := ""
		if len(turno.Vehiculo.Imagenes) > 0 {
			imagen = turno.Vehiculo.Imagenes[0].URL
		}
		resumen.Vehiculo = VehiculoResumen{
			ID:        turno.Vehiculo.ID,
			Marca:     turno.Vehiculo.Marca,
			Modelo:    turno.Vehiculo.Modelo,
			Anio:      turno.Vehiculo.Anio,
			Precio:    turno.Vehiculo.Precio,
			Condicion: turno.Vehiculo.Condicion,
			Tipo:      turno.Vehiculo.Tipo,
			Imagen:    imagen,
		}
	}

	resumen.Cliente = UsuarioResumen{
		ID:     turno.Cliente.ID,
		Nombre: turno.Cliente.Nombre,
	}

	return resumen
}

// aTurnosResumen convierte un listado de modelos en resúmenes.
func aTurnosResumen(turnos []models.TurnoTestDrive) []TurnoTestDriveResumen {
	resumenes := make([]TurnoTestDriveResumen, 0, len(turnos))
	for i := range turnos {
		resumenes = append(resumenes, aTurnoResumen(&turnos[i]))
	}
	return resumenes
}

// recortarFecha normaliza la fecha al formato YYYY-MM-DD. GORM escanea la
// columna date como time.Time y la convierte a RFC3339 al asignarla al campo
// string; el cliente espera el formato ISO de calendario.
func recortarFecha(fecha string) string {
	if len(fecha) >= 10 {
		return fecha[:10]
	}
	return fecha
}
