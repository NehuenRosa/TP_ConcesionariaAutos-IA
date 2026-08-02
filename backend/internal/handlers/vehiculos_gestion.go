package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// VehiculoEntrada es el cuerpo de alta y modificación de un vehículo.
type VehiculoEntrada struct {
	Marca       string   `json:"marca"`
	Modelo      string   `json:"modelo"`
	Anio        int      `json:"anio"`
	Kilometraje int      `json:"kilometraje"`
	Combustible string   `json:"combustible"`
	Transmision string   `json:"transmision"`
	Precio      float64  `json:"precio"`
	Condicion   string   `json:"condicion"`
	Estado      string   `json:"estado"`
	Imagenes    []string `json:"imagenes"`
}

// VehiculoGestionResumen es la ficha básica de un vehículo en el listado de
// administración, que además del catálogo incluye estado y kilometraje.
type VehiculoGestionResumen struct {
	ID          uint    `json:"id"`
	Marca       string  `json:"marca"`
	Modelo      string  `json:"modelo"`
	Anio        int     `json:"anio"`
	Kilometraje int     `json:"kilometraje"`
	Precio      float64 `json:"precio"`
	Condicion   string  `json:"condicion"`
	Estado      string  `json:"estado"`
	Imagen      string  `json:"imagen"`
}

// RespuestaPaginadaGestion envuelve un listado de gestión paginado.
type RespuestaPaginadaGestion struct {
	Datos  []VehiculoGestionResumen `json:"datos"`
	Pagina int                      `json:"pagina"`
	Tamano int                      `json:"tamano"`
	Total  int64                    `json:"total"`
}

// VehiculoGestionHandler agrupa los handlers del ABM administrativo de vehículos.
type VehiculoGestionHandler struct {
	servicio services.VehiculoService
}

// NuevoVehiculoGestionHandler crea un handler de gestión de vehículos.
func NuevoVehiculoGestionHandler(servicio services.VehiculoService) *VehiculoGestionHandler {
	return &VehiculoGestionHandler{servicio: servicio}
}

// Listar responde el listado administrativo de vehículos, con filtro opcional
// por estado y paginación.
func (h *VehiculoGestionHandler) Listar(c *gin.Context) {
	pagina, tamano := parsearPaginacion(c)
	estado := c.Query("estado")

	vehiculos, total, err := h.servicio.ListarParaGestion(c.Request.Context(), estado, pagina, tamano)
	if err != nil {
		if errors.Is(err, services.ErrPaginacionInvalida) || errors.Is(err, services.ErrFiltroEstadoInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el stock de vehículos"})
		return
	}

	resumenes := make([]VehiculoGestionResumen, 0, len(vehiculos))
	for _, vehiculo := range vehiculos {
		resumenes = append(resumenes, aResumenGestion(vehiculo))
	}

	c.JSON(http.StatusOK, RespuestaPaginadaGestion{
		Datos:  resumenes,
		Pagina: pagina,
		Tamano: tamano,
		Total:  total,
	})
}

// ObtenerDetalle responde la ficha técnica completa de cualquier vehículo.
func (h *VehiculoGestionHandler) ObtenerDetalle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de vehículo inválido"})
		return
	}

	vehiculo, err := h.servicio.ObtenerParaGestion(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrVehiculoNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el detalle del vehículo"})
		return
	}

	c.JSON(http.StatusOK, vehiculo)
}

// Crear da de alta un vehículo con su ficha técnica e imágenes.
func (h *VehiculoGestionHandler) Crear(c *gin.Context) {
	var entrada VehiculoEntrada
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de solicitud inválido"})
		return
	}

	vehiculo, err := h.servicio.Crear(c.Request.Context(), aEntrada(entrada))
	if err != nil {
		if errors.Is(err, services.ErrDatosVehiculoInvalidos) || errors.Is(err, services.ErrEstadoInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el vehículo"})
		return
	}

	c.JSON(http.StatusCreated, vehiculo)
}

// Actualizar modifica la ficha técnica, el estado y las imágenes de un vehículo.
func (h *VehiculoGestionHandler) Actualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de vehículo inválido"})
		return
	}

	var entrada VehiculoEntrada
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de solicitud inválido"})
		return
	}

	vehiculo, err := h.servicio.Actualizar(c.Request.Context(), uint(id), aEntrada(entrada))
	if err != nil {
		if errors.Is(err, services.ErrDatosVehiculoInvalidos) || errors.Is(err, services.ErrEstadoInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrVehiculoNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el vehículo"})
		return
	}

	c.JSON(http.StatusOK, vehiculo)
}

// DarDeBaja cambia el estado del vehículo a dado_de_baja sin eliminarlo.
func (h *VehiculoGestionHandler) DarDeBaja(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de vehículo inválido"})
		return
	}

	vehiculo, err := h.servicio.DarDeBaja(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrVehiculoNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo dar de baja el vehículo"})
		return
	}

	c.JSON(http.StatusOK, vehiculo)
}

// aEntrada convierte el DTO del handler en la entrada del service.
func aEntrada(entrada VehiculoEntrada) services.EntradaVehiculo {
	return services.EntradaVehiculo{
		Marca:       entrada.Marca,
		Modelo:      entrada.Modelo,
		Anio:        entrada.Anio,
		Kilometraje: entrada.Kilometraje,
		Combustible: entrada.Combustible,
		Transmision: entrada.Transmision,
		Precio:      entrada.Precio,
		Condicion:   entrada.Condicion,
		Estado:      entrada.Estado,
		Imagenes:    entrada.Imagenes,
	}
}

// aResumenGestion convierte un modelo en la ficha básica del listado de gestión.
func aResumenGestion(vehiculo models.Vehiculo) VehiculoGestionResumen {
	imagen := ""
	if len(vehiculo.Imagenes) > 0 {
		imagen = vehiculo.Imagenes[0].URL
	}

	return VehiculoGestionResumen{
		ID:          vehiculo.ID,
		Marca:       vehiculo.Marca,
		Modelo:      vehiculo.Modelo,
		Anio:        vehiculo.Anio,
		Kilometraje: vehiculo.Kilometraje,
		Precio:      vehiculo.Precio,
		Condicion:   vehiculo.Condicion,
		Estado:      vehiculo.Estado,
		Imagen:      imagen,
	}
}
