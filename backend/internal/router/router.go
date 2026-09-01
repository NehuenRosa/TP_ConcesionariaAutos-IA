package router

import (
	"concesionaria/backend/internal/cifrado"
	"concesionaria/backend/internal/config"
	"concesionaria/backend/internal/googleid"
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
	// "Continuar con Google" (CU-11): habilitado solo si hay client ID.
	googleHabilitado := configuracion.GoogleClientID != ""
	var verificadorGoogle services.VerificadorGoogle
	if googleHabilitado {
		verificadorGoogle = googleid.NuevoVerificador(configuracion.GoogleClientID)
	}
	servicioAutenticacion := services.NuevoAutenticacionService(repositorioUsuarios, configuracion.JWTSecreto, services.DuracionToken, verificadorGoogle)
	handlerAutenticacion := handlers.NuevoAutenticacionHandler(servicioAutenticacion, googleHabilitado, configuracion.GoogleClientID)
	servicioUsuarios := services.NuevoUsuariosService(repositorioUsuarios)
	handlerUsuarios := handlers.NuevoUsuariosHandler(servicioUsuarios)

	repositorioConsultas := repositories.NuevoConsultaRepository(base)
	repositorioMensajes := repositories.NuevoMensajeRepository(base)
	servicioConsultas := services.NuevoConsultaService(repositorioConsultas, repositorioVehiculos)
	servicioMensajes := services.NuevoMensajeService(repositorioMensajes, repositorioConsultas, servicioConsultas)
	handlerConsultas := handlers.NuevoConsultaHandler(servicioConsultas, servicioMensajes)
	handlerMensajes := handlers.NuevoMensajeHandler(servicioMensajes)

	repositorioTurnos := repositories.NuevoTurnoTestDriveRepository(base)
	servicioTurnos := services.NuevoTurnoTestDriveService(repositorioTurnos, repositorioVehiculos)
	handlerTurnos := handlers.NuevoTurnoTestDriveHandler(servicioTurnos)

	repositorioReservas := repositories.NuevoReservaRepository(base)
	servicioReservas := services.NuevoReservaService(repositorioReservas, repositorioVehiculos, configuracion.CbuConcesionaria, configuracion.AliasConcesionaria)
	handlerReservas := handlers.NuevoReservaHandler(servicioReservas)

	repositorioMetricas := repositories.NuevoMetricasRepository(base)
	servicioMetricas := services.NuevoMetricasService(repositorioMetricas)
	handlerMetricas := handlers.NuevoMetricasHandler(servicioMetricas)

	claveEncriptacion := configuracion.ClaveEncriptacion
	if claveEncriptacion == "" {
		claveEncriptacion = configuracion.JWTSecreto
	}
	cifrador, err := cifrado.NuevoCifrador(claveEncriptacion)
	if err != nil {
		panic(err)
	}

	servicioPrecios := services.NuevoServicioPrecios(configuracion.ArgAutosURL)
	repositorioTasaciones := repositories.NuevoTasacionRepository(base)
	repositorioCotizaciones := repositories.NuevoCotizacionRepository(base)
	servicioChatbot := services.NuevoChatbotService(repositorioVehiculos, repositorioTasaciones, repositorioCotizaciones, cifrador, configuracion.ProveedorLLM, configuracion.GoogleAIKey, configuracion.OllamaURL, configuracion.ModeloChatbot, configuracion.ModeloVision, servicioPrecios)
	handlerChatbot := handlers.NuevoChatbotHandler(servicioChatbot)

	servicioCotizaciones := services.NuevoCotizacionService(repositorioCotizaciones, repositorioVehiculos, cifrador, servicioChatbot)
	handlerCotizaciones := handlers.NuevoCotizacionHandler(servicioCotizaciones)
	handlerNotificaciones := handlers.NuevoNotificacionHandler(servicioConsultas, servicioMensajes, servicioCotizaciones)

	api := enrutador.Group("/api")
	{
		api.GET("/health", handlers.Salud)

		chatbot := api.Group("/chatbot")
		{
			chatbot.POST("/mensajes", middleware.AutenticacionOpcional(configuracion.JWTSecreto), handlerChatbot.Responder)
			chatbot.POST("/tasacion", handlerChatbot.Tasacion)
			chatbot.POST("/tasacion/confirmar", handlerChatbot.ConfirmarTasacion)
		}

		autenticacion := api.Group("/auth")
		{
			autenticacion.POST("/registro", handlerAutenticacion.Registrar)
			autenticacion.POST("/login", handlerAutenticacion.IniciarSesion)
			autenticacion.POST("/google", handlerAutenticacion.IniciarSesionConGoogle)
			autenticacion.GET("/proveedores", handlerAutenticacion.Proveedores)
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

		gestionUsuarios := api.Group("/admin/usuarios")
		gestionUsuarios.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		gestionUsuarios.Use(middleware.ExigirRol("administrador"))
		{
			gestionUsuarios.GET("", handlerUsuarios.Listar)
			gestionUsuarios.POST("", handlerUsuarios.Crear)
			gestionUsuarios.PUT("/:id", handlerUsuarios.Actualizar)
			gestionUsuarios.DELETE("/:id", handlerUsuarios.Eliminar)
		}

		metricas := api.Group("/admin/metricas")
		metricas.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		metricas.Use(middleware.ExigirRol("administrador"))
		{
			metricas.GET("", handlerMetricas.ObtenerMetricas)
		}

		consultas := api.Group("/consultas")
		consultas.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		{
			consultas.POST("", handlerConsultas.Crear)
			consultas.GET("/mis-consultas", handlerConsultas.ListarMisConsultas)

			consultas.GET("/:id/mensajes", handlerMensajes.ObtenerMensajes)
			consultas.GET("/:id/mensajes/nuevos", handlerMensajes.ObtenerNuevos)
			consultas.POST("/:id/mensajes", handlerMensajes.Enviar)
			consultas.PUT("/:id/mensajes/leidos", handlerMensajes.MarcarLeidos)

			gestionConsultas := consultas.Group("")
			gestionConsultas.Use(middleware.ExigirRol("vendedor"))
			{
				gestionConsultas.GET("/bandeja", handlerConsultas.ListarBandeja)
				gestionConsultas.PUT("/:id/tomar", handlerConsultas.Tomar)
				gestionConsultas.PUT("/:id/cerrar", handlerConsultas.Cerrar)
				gestionConsultas.DELETE("/:id", handlerConsultas.Eliminar)
			}
		}

		notificaciones := api.Group("/notificaciones")
		notificaciones.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		{
			notificaciones.GET("/contador", handlerNotificaciones.Contador)
		}

		api.GET("/test-drives/franjas", handlerTurnos.Franjas)

		testDrives := api.Group("/test-drives")
		testDrives.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		{
			testDrives.POST("", handlerTurnos.Solicitar)
			testDrives.GET("/mis-turnos", handlerTurnos.ListarMisTurnos)
			testDrives.DELETE("/:id", handlerTurnos.Cancelar)
			testDrives.DELETE("/:id/eliminar", handlerTurnos.Eliminar)

			gestionTestDrives := testDrives.Group("")
			gestionTestDrives.Use(middleware.ExigirRol("vendedor"))
			{
				gestionTestDrives.GET("", handlerTurnos.Listar)
				gestionTestDrives.PUT("/:id/confirmar", handlerTurnos.Confirmar)
				gestionTestDrives.PUT("/:id/cancelar", handlerTurnos.CancelarComoVendedor)
				gestionTestDrives.PUT("/:id/completar", handlerTurnos.Completar)
			}
		}

		reservas := api.Group("/reservas")
		reservas.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		{
			reservas.POST("", handlerReservas.Crear)
			reservas.GET("/datos-transferencia", handlerReservas.DatosTransferencia)
			reservas.GET("/mis-reservas", handlerReservas.ListarMisReservas)
			reservas.DELETE("/:id", handlerReservas.Cancelar)
			// El comprobante lo ve el dueño o el personal: el control fino va
			// dentro del handler (dueño de la reserva o rol vendedor/admin).
			reservas.POST("/:id/comprobante", handlerReservas.SubirComprobante)
			reservas.GET("/:id/comprobante", handlerReservas.ObtenerComprobante)

			gestionReservas := reservas.Group("")
			gestionReservas.Use(middleware.ExigirRol("vendedor"))
			{
				gestionReservas.GET("", handlerReservas.Listar)
				gestionReservas.PUT("/:id/confirmar", handlerReservas.ConfirmarVenta)
				gestionReservas.PUT("/:id/cancelar", handlerReservas.CancelarComoVendedor)
			}
		}

		cotizaciones := api.Group("/cotizaciones")
		cotizaciones.Use(middleware.AutenticacionJWT(configuracion.JWTSecreto))
		{
			cotizaciones.POST("", handlerCotizaciones.Crear)
			cotizaciones.GET("/mis-cotizaciones", handlerCotizaciones.ListarMisCotizaciones)
			// Rutas de atención personal: solo vendedores.
			cotizaciones.GET("/bandeja", middleware.ExigirRol("vendedor"), handlerCotizaciones.ListarBandeja)
			cotizaciones.GET("/:id", handlerCotizaciones.Obtener)
			cotizaciones.GET("/:id/personal", middleware.ExigirRol("vendedor"), handlerCotizaciones.ObtenerPersonal)
			// Fetch incremental del chat: devuelven solo los mensajes con id
			// mayor a ?desdeId (el polling no recarga el historial completo).
			cotizaciones.GET("/:id/mensajes", handlerCotizaciones.ObtenerMensajesNuevos)
			cotizaciones.GET("/:id/mensajes/personal", middleware.ExigirRol("vendedor"), handlerCotizaciones.ObtenerMensajesNuevosPersonal)
			cotizaciones.POST("/:id/mensajes", handlerCotizaciones.EnviarMensaje)
			cotizaciones.PUT("/:id/tomar", middleware.ExigirRol("vendedor"), handlerCotizaciones.Tomar)
			cotizaciones.POST("/:id/mensajes-vendedor", middleware.ExigirRol("vendedor"), handlerCotizaciones.ResponderComoVendedor)
			cotizaciones.PUT("/:id/cerrar", handlerCotizaciones.Cerrar)
			cotizaciones.PUT("/:id/cerrar-personal", middleware.ExigirRol("vendedor"), handlerCotizaciones.CerrarPersonal)
		}
	}

	return enrutador
}
