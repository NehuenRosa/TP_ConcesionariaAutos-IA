package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/token"

	"github.com/gin-gonic/gin"
)

// AutenticacionJWT valida el token JWT del encabezado Authorization y guarda
// el identificador y el rol del usuario en el contexto de la petición.
func AutenticacionJWT(secreto string) gin.HandlerFunc {
	return func(c *gin.Context) {
		encabezado := c.GetHeader("Authorization")
		if !strings.HasPrefix(encabezado, "Bearer ") {
			responderNoAutorizado(c)
			return
		}

		reclamos, err := token.Validar(strings.TrimPrefix(encabezado, "Bearer "), secreto)
		if err != nil {
			responderNoAutorizado(c)
			return
		}

		c.Set("usuario_id", strconv.FormatUint(uint64(reclamos.UsuarioID), 10))
		c.Set("rol", reclamos.Rol)
		c.Next()
	}
}

// ExigirRol verifica que el rol del usuario autenticado sea el requerido o
// administrador (que tiene acceso a todas las rutas protegidas).
func ExigirRol(rolRequerido string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rol, existe := c.Get("rol")
		if !existe {
			responderNoAutorizado(c)
			return
		}

		rolActual, ok := rol.(string)
		if !ok {
			responderNoAutorizado(c)
			return
		}

		if rolActual != rolRequerido && rolActual != models.RolAdministrador {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No tiene permisos para acceder a este recurso",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func responderNoAutorizado(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "No autorizado",
	})
	c.Abort()
}
