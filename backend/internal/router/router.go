package router

import (
	"concesionaria/backend/internal/config"
	"concesionaria/backend/internal/handlers"
	"concesionaria/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Nuevo arma el enrutador de la API con sus rutas y middlewares.
func Nuevo(_ *gorm.DB, configuracion config.Configuracion) *gin.Engine {
	enrutador := gin.Default()
	enrutador.Use(middleware.CORS(configuracion.OrigenesCORS))

	api := enrutador.Group("/api")
	{
		api.GET("/health", handlers.Salud)

		vehiculos := api.Group("/vehiculos")
		{
			vehiculos.GET("", handlers.ListarVehiculos)
			vehiculos.GET("/:id", handlers.ObtenerVehiculo)
		}
	}

	return enrutador
}
