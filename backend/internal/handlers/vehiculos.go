package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// VehiculoResumen es la ficha básica de un vehículo en el listado del catálogo.
type VehiculoResumen struct {
	ID        uint    `json:"id"`
	Marca     string  `json:"marca"`
	Modelo    string  `json:"modelo"`
	Anio      int     `json:"anio"`
	Precio    float64 `json:"precio"`
	Condicion string  `json:"condicion"`
	Imagen    string  `json:"imagen"`
}

// RespuestaPaginada envuelve un listado paginado con sus metadatos.
type RespuestaPaginada struct {
	Datos  []VehiculoResumen `json:"datos"`
	Pagina int               `json:"pagina"`
	Tamano int               `json:"tamano"`
	Total  int64             `json:"total"`
}

// VehiculoHandler agrupa los handlers del catálogo público de vehículos.
type VehiculoHandler struct {
	servicio services.VehiculoService
}

// NuevoVehiculoHandler crea un handler de vehículos.
func NuevoVehiculoHandler(servicio services.VehiculoService) *VehiculoHandler {
	return &VehiculoHandler{servicio: servicio}
}

// Listar responde el catálogo paginado de vehículos disponibles.
func (h *VehiculoHandler) Listar(c *gin.Context) {
	pagina, tamano := parsearPaginacion(c)

	vehiculos, total, err := h.servicio.ListarDisponibles(c.Request.Context(), pagina, tamano)
	if err != nil {
		if errors.Is(err, services.ErrPaginacionInvalida) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Paginación inválida"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el catálogo"})
		return
	}

	resumenes := make([]VehiculoResumen, 0, len(vehiculos))
	for _, vehiculo := range vehiculos {
		resumenes = append(resumenes, aResumen(vehiculo))
	}

	c.JSON(http.StatusOK, RespuestaPaginada{
		Datos:  resumenes,
		Pagina: pagina,
		Tamano: tamano,
		Total:  total,
	})
}

// ObtenerDetalle responde la ficha técnica completa de un vehículo disponible.
func (h *VehiculoHandler) ObtenerDetalle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de vehículo inválido"})
		return
	}

	vehiculo, err := h.servicio.ObtenerPorID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrVehiculoNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado o no disponible"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el detalle del vehículo"})
		return
	}

	c.JSON(http.StatusOK, vehiculo)
}

// parsearPaginacion lee los parámetros de consulta página y tamaño con
// valores por defecto (1 y 12) y validando el rango permitido.
func parsearPaginacion(c *gin.Context) (int, int) {
	pagina, err := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	if err != nil || pagina < 1 {
		pagina = 1
	}

	tamano, err := strconv.Atoi(c.DefaultQuery("tamano", "12"))
	if err != nil || tamano < 1 {
		tamano = 12
	}
	if tamano > 100 {
		tamano = 100
	}

	return pagina, tamano
}

// aResumen convierte un modelo en la ficha básica del listado. La imagen es la
// primera de la galería, si existe.
func aResumen(vehiculo models.Vehiculo) VehiculoResumen {
	imagen := ""
	if len(vehiculo.Imagenes) > 0 {
		imagen = vehiculo.Imagenes[0].URL
	}

	return VehiculoResumen{
		ID:        vehiculo.ID,
		Marca:     vehiculo.Marca,
		Modelo:    vehiculo.Modelo,
		Anio:      vehiculo.Anio,
		Precio:    vehiculo.Precio,
		Condicion: vehiculo.Condicion,
		Imagen:    imagen,
	}
}
