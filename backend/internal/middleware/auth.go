package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AutenticacionJWT valida el token JWT en el encabezado Authorization.
// TODO: implementar la validación real junto con el caso de uso CU-01.
// Actualmente es un stub: permite el paso sin autenticar para no bloquear
// el desarrollo del esqueleto.
func AutenticacionJWT(_ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// ExigirRol verifica el rol del usuario autenticado.
// TODO: implementar junto con el caso de uso CU-01.
func ExigirRol(_ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func responderNoAutorizado(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "No autorizado",
	})
}
