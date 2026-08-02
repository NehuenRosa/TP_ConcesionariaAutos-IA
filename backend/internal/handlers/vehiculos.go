package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListarVehiculos es un stub del catálogo público (CU-03).
// TODO: implementar cuando se desarrolle el caso de uso.
func ListarVehiculos(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Catálogo de vehículos en construcción",
	})
}

// ObtenerVehiculo es un stub del detalle del catálogo público (CU-03).
// TODO: implementar cuando se desarrolle el caso de uso.
func ObtenerVehiculo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Detalle de vehículo en construcción",
	})
}
