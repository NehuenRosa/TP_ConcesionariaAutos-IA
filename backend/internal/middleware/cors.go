package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS configura los encabezados de acceso cruzado para el frontend.
func CORS(origenesPermitidos string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origen := c.GetHeader("Origin")
		permitido := origenesPermitidos == "*" || contiene(origenesPermitidos, origen)

		if permitido {
			if origenesPermitidos == "*" {
				// Con origen comodín no se puede enviar Allow-Credentials
				// (la spec lo prohíbe); la app usa tokens Bearer, no cookies.
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				c.Header("Access-Control-Allow-Origin", origen)
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func contiene(lista string, valor string) bool {
	for _, item := range strings.Split(lista, ",") {
		if strings.TrimSpace(item) == valor {
			return true
		}
	}
	return false
}
