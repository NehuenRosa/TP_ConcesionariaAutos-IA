package handlers

import (
	"errors"
	"fmt"
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
	Tipo      string  `json:"tipo"`
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

// Listar responde el catálogo paginado de vehículos disponibles, con búsqueda,
// filtros combinables y ordenamiento opcionales.
func (h *VehiculoHandler) Listar(c *gin.Context) {
	pagina, tamano := parsearPaginacion(c)
	filtros, err := parsearFiltros(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehiculos, total, err := h.servicio.ListarDisponibles(c.Request.Context(), filtros, pagina, tamano)
	if err != nil {
		if errors.Is(err, services.ErrPaginacionInvalida) || errors.Is(err, services.ErrFiltroInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// parsearFiltros lee los query params opcionales de búsqueda, filtros y
// ordenamiento del catálogo. Devuelve un error con mensaje en español si algún
// valor numérico no es válido.
func parsearFiltros(c *gin.Context) (services.FiltrosBusqueda, error) {
	filtros := services.FiltrosBusqueda{
		Busqueda:       c.Query("busqueda"),
		Marca:          c.Query("marca"),
		Modelo:         c.Query("modelo"),
		Tipo:           c.Query("tipo"),
		Combustible:    c.Query("combustible"),
		Condicion:      c.Query("condicion"),
		OrdenPor:       c.Query("orden_por"),
		OrdenDireccion: c.Query("orden_direccion"),
	}

	anioMin, err := parsearEnteroOpcional(c, "anio_min")
	if err != nil {
		return filtros, err
	}
	filtros.AnioMin = anioMin

	anioMax, err := parsearEnteroOpcional(c, "anio_max")
	if err != nil {
		return filtros, err
	}
	filtros.AnioMax = anioMax

	precioMin, err := parsearDecimalOpcional(c, "precio_min")
	if err != nil {
		return filtros, err
	}
	filtros.PrecioMin = precioMin

	precioMax, err := parsearDecimalOpcional(c, "precio_max")
	if err != nil {
		return filtros, err
	}
	filtros.PrecioMax = precioMax

	return filtros, nil
}

// parsearEnteroOpcional parsea un query param entero opcional. Si el parámetro
// no está presente devuelve nil; si está y no es numérico devuelve error.
func parsearEnteroOpcional(c *gin.Context, nombre string) (*int, error) {
	valor := c.Query(nombre)
	if valor == "" {
		return nil, nil
	}
	numero, err := strconv.Atoi(valor)
	if err != nil {
		return nil, fmt.Errorf("parámetro %s inválido", nombre)
	}
	return &numero, nil
}

// parsearDecimalOpcional parsea un query param decimal opcional. Si el
// parámetro no está presente devuelve nil; si está y no es numérico devuelve
// error.
func parsearDecimalOpcional(c *gin.Context, nombre string) (*float64, error) {
	valor := c.Query(nombre)
	if valor == "" {
		return nil, nil
	}
	numero, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		return nil, fmt.Errorf("parámetro %s inválido", nombre)
	}
	return &numero, nil
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
		Tipo:      vehiculo.Tipo,
		Imagen:    imagen,
	}
}
