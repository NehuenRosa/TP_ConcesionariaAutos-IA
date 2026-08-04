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

	repositorioUsuarios := repositories.NuevoUsuarioRepository(base)
	servicioAutenticacion := services.NuevoAutenticacionService(repositorioUsuarios, configuracion.JWTSecreto, services.DuracionToken)
	handlerAutenticacion := handlers.NuevoAutenticacionHandler(servicioAutenticacion)

	repositorioConsultas := repositories.NuevoConsultaRepository(base)
	repositorioMensajes := repositories.NuevoMensajeRepository(base)
	servicioConsultas := services.NuevoConsultaService(repositorioConsultas, repositorioVehiculos)
	servicioMensajes := services.NuevoMensajeService(repositorioMensajes, repositorioConsultas, servicioConsultas)
	handlerConsultas := handlers.NuevoConsultaHandler(servicioConsultas, servicioMensajes)
	handlerMensajes := handlers.NuevoMensajeHandler(servicioMensajes)
	handlerNotificaciones := handlers.NuevoNotificacionHandler(servicioConsultas, servicioMensajes)

	api := enrutador.Group("/api")
	{
		api.GET("/health", handlers.Salud)

		autenticacion := api.Group("/auth")
		{
			autenticacion.POST("/registro", handlerAutenticacion.Registrar)
			autenticacion.POST("/login", handlerAutenticacion.IniciarSesion)
			autenticacion.GET("/perfil", middleware.AutenticacionJWT(configuracion.JWTSecreto), handlerAutenticacion.Perfil)
		}

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

		consultas := api.Group("/consultas")
		consultas.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		{
			consultas.POST("", handlerConsultas.Crear)
			consultas.GET("/mis-consultas", handlerConsultas.ListarMisConsultas)
			consultas.GET("/bandeja", handlerConsultas.ListarBandeja)
			consultas.PUT("/:id/tomar", handlerConsultas.Tomar)
			consultas.PUT("/:id/cerrar", handlerConsultas.Cerrar)
			consultas.DELETE("/:id", handlerConsultas.Eliminar)

			consultas.GET("/:id/mensajes", handlerMensajes.ObtenerMensajes)
			consultas.GET("/:id/mensajes/nuevos", handlerMensajes.ObtenerNuevos)
			consultas.POST("/:id/mensajes", handlerMensajes.Enviar)
			consultas.PUT("/:id/mensajes/leidos", handlerMensajes.MarcarLeidos)
		}

		notificaciones := api.Group("/notificaciones")
		notificaciones.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		{
			notificaciones.GET("/contador", handlerNotificaciones.Contador)
		}
	}

	return enrutador
}
