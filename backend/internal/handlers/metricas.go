package handlers

import (
	"errors"
	"net/http"

	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// MetricasHandler agrupa los handlers de métricas del panel de administración.
type MetricasHandler struct {
	servicio services.MetricasService
}

// NuevoMetricasHandler crea un handler de métricas.
func NuevoMetricasHandler(servicio services.MetricasService) *MetricasHandler {
	return &MetricasHandler{servicio: servicio}
}

// ObtenerMetricas devuelve el payload de métricas del panel de administración.
func (h *MetricasHandler) ObtenerMetricas(c *gin.Context) {
	metricas, err := h.servicio.Obtener(c.Request.Context(), c.Query("periodo"))
	if err != nil {
		if errors.Is(err, services.ErrPeriodoInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las métricas"})
		return
	}

	c.JSON(http.StatusOK, metricas)
}
