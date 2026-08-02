package router

import (
	"concesionaria/backend/internal/config"
	"concesionaria/backend/internal/handlers"
	"concesionaria/backend/internal/middleware"
	"concesionaria/backend/internal/repositories"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Nuevo arma el enrutador de la API con sus rutas y middlewares.
func Nuevo(base *gorm.DB, configuracion config.Configuracion) *gin.Engine {
	enrutador := gin.Default()
	enrutador.Use(middleware.CORS(configuracion.OrigenesCORS))

	repositorioVehiculos := repositories.NuevoVehiculoRepository(base)
	servicioVehiculos := services.NuevoVehiculoService(repositorioVehiculos)
	handlerVehiculos := handlers.NuevoVehiculoHandler(servicioVehiculos)
	handlerGestionVehiculos := handlers.NuevoVehiculoGestionHandler(servicioVehiculos)

	api := enrutador.Group("/api")
	{
		api.GET("/health", handlers.Salud)

		vehiculos := api.Group("/vehiculos")
		{
			vehiculos.GET("", handlerVehiculos.Listar)
			vehiculos.GET("/:id", handlerVehiculos.ObtenerDetalle)
		}

		gestionVehiculos := api.Group("/admin/vehiculos")
		gestionVehiculos.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		gestionVehiculos.Use(middleware.ExigirRol("administrador"))
		{
			gestionVehiculos.GET("", handlerGestionVehiculos.Listar)
			gestionVehiculos.GET("/:id", handlerGestionVehiculos.ObtenerDetalle)
			gestionVehiculos.POST("", handlerGestionVehiculos.Crear)
			gestionVehiculos.PUT("/:id", handlerGestionVehiculos.Actualizar)
			gestionVehiculos.DELETE("/:id", handlerGestionVehiculos.DarDeBaja)
		}
	}

	return enrutador
}
