package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ConsultaResumen es la ficha resumida de una consulta para listados.
type ConsultaResumen struct {
	ID           uint              `json:"id"`
	Vehiculo     VehiculoResumen   `json:"vehiculo"`
	Cliente      UsuarioResumen    `json:"cliente"`
	Vendedor     *UsuarioResumen   `json:"vendedor,omitempty"`
	Estado       string            `json:"estado"`
	UltimoMensaje *MensajeResumen   `json:"ultimoMensaje,omitempty"`
	MensajesNuevos int             `json:"mensajesNuevos"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
}

// UsuarioResumen es la ficha básica de un usuario para consultas.
type UsuarioResumen struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
}

// MensajeResumen es la información mínima de un mensaje para preview.
type MensajeResumen struct {
	Contenido string `json:"contenido"`
	CreatedAt string `json:"createdAt"`
}

// ConsultaHandler agrupa los handlers de consultas.
type ConsultaHandler struct {
	servicio       services.ConsultaService
	servicioMensajes services.MensajeService
}

// NuevoConsultaHandler crea un handler de consultas.
func NuevoConsultaHandler(
	servicio services.ConsultaService,
	servicioMensajes services.MensajeService,
) *ConsultaHandler {
	return &ConsultaHandler{
		servicio:         servicio,
		servicioMensajes: servicioMensajes,
	}
}

// Crear crea una nueva consulta desde el cliente.
func (h *ConsultaHandler) Crear(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	var entrada struct {
		VehiculoID uint   `json:"vehiculoId"`
		Mensaje    string `json:"mensaje"`
	}
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if entrada.VehiculoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El vehículo es obligatorio"})
		return
	}

	consulta, err := h.servicio.Crear(c.Request.Context(), clienteID, entrada.VehiculoID, entrada.Mensaje)
	if err != nil {
		if errors.Is(err, services.ErrMensajeVacio) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrVehiculoNoDisponible) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la consulta"})
		return
	}

	c.JSON(http.StatusCreated, aConsultaResumen(consulta))
}

// ListarMisConsultas lista las consultas del cliente autenticado.
func (h *ConsultaHandler) ListarMisConsultas(c *gin.Context) {
	clienteID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultas, err := h.servicio.ListarPorCliente(c.Request.Context(), clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las consultas"})
		return
	}

	// Calcular mensajes no leídos por consulta
	ids := make([]uint, len(consultas))
	for i, consulta := range consultas {
		ids[i] = consulta.ID
	}
	noLeidos, err := h.servicioMensajes.ContarNoLeidosPorConsultas(c.Request.Context(), ids, clienteID)
	if err != nil {
		noLeidos = make(map[uint]int)
	}

	resumenes := make([]ConsultaResumen, 0, len(consultas))
	for _, consulta := range consultas {
		resumen := aConsultaResumen(&consulta)
		resumen.MensajesNuevos = noLeidos[consulta.ID]
		resumenes = append(resumenes, resumen)
	}

	c.JSON(http.StatusOK, resumenes)
}

// ListarBandeja lista las consultas del vendedor autenticado.
func (h *ConsultaHandler) ListarBandeja(c *gin.Context) {
	vendedorID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultas, err := h.servicio.ListarPorVendedor(c.Request.Context(), vendedorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las consultas"})
		return
	}

	// Calcular mensajes no leídos por consulta
	ids := make([]uint, len(consultas))
	for i, consulta := range consultas {
		ids[i] = consulta.ID
	}
	noLeidos, err := h.servicioMensajes.ContarNoLeidosPorConsultas(c.Request.Context(), ids, vendedorID)
	if err != nil {
		noLeidos = make(map[uint]int)
	}

	resumenes := make([]ConsultaResumen, 0, len(consultas))
	for _, consulta := range consultas {
		resumen := aConsultaResumen(&consulta)
		resumen.MensajesNuevos = noLeidos[consulta.ID]
		resumenes = append(resumenes, resumen)
	}

	c.JSON(http.StatusOK, resumenes)
}

// Tomar asigna el vendedor a una consulta pendiente.
func (h *ConsultaHandler) Tomar(c *gin.Context) {
	vendedorID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de consulta inválido"})
		return
	}

	consulta, err := h.servicio.Tomar(c.Request.Context(), uint(consultaID), vendedorID)
	if err != nil {
		if errors.Is(err, services.ErrConsultaNoEncontrada) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrConsultaNoPendiente) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo tomar la consulta"})
		return
	}

	c.JSON(http.StatusOK, aConsultaResumen(consulta))
}

// Cerrar cierra una consulta asignada al vendedor.
func (h *ConsultaHandler) Cerrar(c *gin.Context) {
	vendedorID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de consulta inválido"})
		return
	}

	consulta, err := h.servicio.Cerrar(c.Request.Context(), uint(consultaID), vendedorID)
	if err != nil {
		if errors.Is(err, services.ErrConsultaNoEncontrada) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrConsultaYaCerrada) || errors.Is(err, services.ErrNoEsVendedorAsignado) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo cerrar la consulta"})
		return
	}

	c.JSON(http.StatusOK, aConsultaResumen(consulta))
}

// Eliminar elimina una consulta cerrada.
func (h *ConsultaHandler) Eliminar(c *gin.Context) {
	vendedorID, err := extraerUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	consultaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de consulta inválido"})
		return
	}

	err = h.servicio.Eliminar(c.Request.Context(), uint(consultaID), vendedorID)
	if err != nil {
		if errors.Is(err, services.ErrConsultaNoEncontrada) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrConsultaNoCerrada) || errors.Is(err, services.ErrNoEsVendedorAsignado) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la consulta"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// extraerUsuarioID obtiene el ID del usuario del contexto JWT.
func extraerUsuarioID(c *gin.Context) (uint, error) {
	usuarioIDStr, existe := c.Get("usuario_id")
	if !existe {
		return 0, errors.New("no autorizado")
	}

	idStr, ok := usuarioIDStr.(string)
	if !ok {
		return 0, errors.New("no autorizado")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("no autorizado")
	}

	return uint(id), nil
}

// aConsultaResumen convierte un modelo en el resumen para listados.
func aConsultaResumen(consulta *models.Consulta) ConsultaResumen {
	resumen := ConsultaResumen{
		ID:        consulta.ID,
		Estado:    consulta.Estado,
		CreatedAt: consulta.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: consulta.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Vehículo
	if consulta.Vehiculo.ID != 0 {
		imagen := ""
		if len(consulta.Vehiculo.Imagenes) > 0 {
			imagen = consulta.Vehiculo.Imagenes[0].URL
		}
		resumen.Vehiculo = VehiculoResumen{
			ID:        consulta.Vehiculo.ID,
			Marca:     consulta.Vehiculo.Marca,
			Modelo:    consulta.Vehiculo.Modelo,
			Anio:      consulta.Vehiculo.Anio,
			Precio:    consulta.Vehiculo.Precio,
			Condicion: consulta.Vehiculo.Condicion,
			Tipo:      consulta.Vehiculo.Tipo,
			Imagen:    imagen,
		}
	}

	// Cliente
	resumen.Cliente = UsuarioResumen{
		ID:     consulta.Cliente.ID,
		Nombre: consulta.Cliente.Nombre,
	}

	// Vendedor
	if consulta.Vendedor != nil {
		resumen.Vendedor = &UsuarioResumen{
			ID:     consulta.Vendedor.ID,
			Nombre: consulta.Vendedor.Nombre,
		}
	}

	// Último mensaje
	if len(consulta.Mensajes) > 0 {
		ultimo := consulta.Mensajes[len(consulta.Mensajes)-1]
		resumen.UltimoMensaje = &MensajeResumen{
			Contenido: ultimo.Contenido,
			CreatedAt: ultimo.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return resumen
}

// aConsultaResumenFromModel convierte un modelo de servicios en resumen.
func aConsultaResumenFromModel(consulta *models.Consulta) ConsultaResumen {
	return aConsultaResumen(consulta)
}
