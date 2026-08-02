package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Salud responde el estado del servicio (health check).
func Salud(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"estado": "ok",
	})
}
