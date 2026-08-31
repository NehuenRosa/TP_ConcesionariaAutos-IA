package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"concesionaria/backend/internal/cifrado"
	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

// Errores de negocio de cotizaciones.
var (
	// ErrCotizacionNoEncontrada indica que la cotización no existe.
	ErrCotizacionNoEncontrada = errors.New("cotización no encontrada")
	// ErrCotizacionNoPertenece indica que la cotización pertenece a otro cliente.
	ErrCotizacionNoPertenece = errors.New("la cotización no pertenece al cliente")
	// ErrCotizacionCerradaMensajes indica que no se pueden enviar mensajes a una
	// cotización cerrada.
	ErrCotizacionCerradaMensajes = errors.New("no se pueden enviar mensajes a una cotización cerrada")
	// ErrCotizacionYaCerrada indica que la cotización ya estaba cerrada.
	ErrCotizacionYaCerrada = errors.New("la cotización ya está cerrada")
	// ErrCotizacionYaAtendida indica que otro vendedor ya tomó la cotización.
	ErrCotizacionYaAtendida = errors.New("la cotización ya está siendo atendida por otro vendedor")
	// ErrCotizacionNoTomada indica que hay que tomar la cotización antes de
	// responder como vendedor.
	ErrCotizacionNoTomada = errors.New("tenés que tomar la cotización antes de responder")
)

// Lados para el marcado de lectura de un hilo de cotización.
const (
	// LadoCliente marca la lectura desde la vista del cliente.
	LadoCliente = "cliente"
	// LadoPersonal marca la lectura desde la vista del vendedor asignado.
	LadoPersonal = "personal"
)

// EstadoConversacionCotizacion agrupa los mensajes nuevos de un hilo (fetch
// incremental) con los datos de cabecera necesarios para refrescar la vista sin
// recargar el historial completo.
type EstadoConversacionCotizacion struct {
	// Mensajes son los mensajes posteriores a desdeID, ya descifrados.
	Mensajes []models.MensajeCotizacion
	// Total es la cantidad total de mensajes del hilo.
	Total int64
	// Estado es el estado actual de la cotización.
	Estado string
	// Vendedor es el vendedor que tomó la cotización (nil si está sin asignar).
	// Grupo de datos de cabecera para que la vista se entere de una toma.
	Vendedor *models.Usuario
	// FechaToma indica cuándo tomó la cotización el vendedor.
	FechaToma *time.Time
}

// generadorCotizacion convierte una respuesta preguntada dentro de una
// cotización en contexto del vehículo cotizado.
type generadorCotizacion interface {
	// GenerarCotizacion responde el mensaje del cliente usando la ficha real del
	// vehículo cotizado y el historial descifrado de la conversación.
	GenerarCotizacion(ctx context.Context, vehiculo models.Vehiculo, historial []TurnoChat, mensaje string) (string, error)
}

// CotizacionService define el contrato de la lógica de negocio de cotizaciones.
type CotizacionService interface {
	// Crear crea la cotización de un vehículo con el primer mensaje del cliente
	// y la primera respuesta del asistente.
	Crear(ctx context.Context, clienteID uint, vehiculoID uint, mensaje string) (*models.Cotizacion, error)
	// ListarPorCliente lista las cotizaciones del cliente con su último mensaje
	// descifrado.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Cotizacion, error)
	// ObtenerPorCliente devuelve la cotización del cliente con los mensajes
	// descifrados.
	ObtenerPorCliente(ctx context.Context, clienteID uint, cotizacionID uint) (*models.Cotizacion, error)
	// EnviarMensaje agrega el mensaje del cliente, genera la respuesta del
	// asistente y devuelve el texto de la respuesta.
	EnviarMensaje(ctx context.Context, clienteID uint, cotizacionID uint, mensaje string) (string, error)
	// Cerrar cambia el estado de la cotización a cerrada.
	Cerrar(ctx context.Context, clienteID uint, cotizacionID uint) (*models.Cotizacion, error)
	// ListarBandeja devuelve todas las cotizaciones con su último mensaje
	// descifrado, para la bandeja del personal.
	ListarBandeja(ctx context.Context) ([]models.Cotizacion, error)
	// ObtenerPersonal devuelve cualquier cotización con sus mensajes
	// descifrados, para la vista de atención del vendedor.
	ObtenerPersonal(ctx context.Context, cotizacionID uint) (*models.Cotizacion, error)
	// ObtenerMensajesDesde devuelve los mensajes del hilo posteriores a desdeID
	// (ya descifrados) junto con el total y la cabecera actual (estado,
	// vendedor asignado y fecha de toma). Reemplaza el recorte completo en el
	// polling del chat por un fetch incremental. El rol cliente exige que la
	// cotización sea propia; cualquier rol de personal puede pedir el delta.
	ObtenerMensajesDesde(ctx context.Context, usuarioID uint, rol string, cotizacionID uint, desdeID uint) (*EstadoConversacionCotizacion, error)
	// Tomar asigna la cotización al vendedor autenticado: a partir de ese
	// momento la IA deja de responder. Es idempotente para el mismo vendedor y
	// falla si otro vendedor ya la tomó o si está cerrada.
	Tomar(ctx context.Context, vendedorID uint, cotizacionID uint) (*models.Cotizacion, error)
	// ResponderComoVendedor guarda el mensaje del vendedor cifrado sin pasar
	// por la IA. Requiere haber tomado la cotización y que esté abierta.
	ResponderComoVendedor(ctx context.Context, vendedorID uint, cotizacionID uint, mensaje string) (*models.Cotizacion, error)
	// CerrarPersonal cierra una cotización desde la bandeja del personal.
	CerrarPersonal(ctx context.Context, cotizacionID uint) (*models.Cotizacion, error)
	// ContarNoLeidos devuelve los mensajes de cotizaciones sin leer del
	// usuario según su rol: respuestas del vendedor para un cliente; mensajes
	// de cliente para el personal.
	ContarNoLeidos(ctx context.Context, usuarioID uint, rol string) (int64, error)
	// MarcarLeidas marca como leídos los mensajes del hilo para el lado
	// indicado (LadoCliente o LadoPersonal). El lado personal solo marca cuando
	// el usuario es el vendedor asignado; si no, no hace nada.
	MarcarLeidas(ctx context.Context, usuarioID uint, cotizacionID uint, lado string) error
}

// cotizacionService implementa CotizacionService. El contenido de los mensajes
// se cifra con el cifrador antes de persistir y se descifra al leerlo.
type cotizacionService struct {
	repositorioCotizaciones repositories.CotizacionRepository
	repositorioVehiculos    repositories.VehiculoRepository
	cifrador                cifrado.Cifrador
	generador               generadorCotizacion
}

// NuevoCotizacionService crea un servicio de cotizaciones.
func NuevoCotizacionService(
	repositorioCotizaciones repositories.CotizacionRepository,
	repositorioVehiculos repositories.VehiculoRepository,
	cifrador cifrado.Cifrador,
	generador generadorCotizacion,
) CotizacionService {
	return &cotizacionService{
		repositorioCotizaciones: repositorioCotizaciones,
		repositorioVehiculos:    repositorioVehiculos,
		cifrador:                cifrador,
		generador:               generador,
	}
}

// Crear valida que el vehículo esté disponible, genera la primera respuesta del
// asistente y guarda la cotización con sus dos mensajes iniciales cifrados.
func (s *cotizacionService) Crear(ctx context.Context, clienteID uint, vehiculoID uint, mensaje string) (*models.Cotizacion, error) {
	vehiculo, err := s.repositorioVehiculos.ObtenerPorID(ctx, vehiculoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehiculoNoDisponible
		}
		return nil, fmt.Errorf("obtener vehículo: %w", err)
	}
	if vehiculo.Estado != models.EstadoDisponible {
		return nil, ErrVehiculoNoDisponible
	}

	mensajeCliente := strings.TrimSpace(mensaje)
	if mensajeCliente == "" {
		mensajeCliente = "Hola, quiero cotizar este vehículo."
	}

	respuestaIA, err := s.generador.GenerarCotizacion(ctx, *vehiculo, nil, mensajeCliente)
	if err != nil {
		return nil, fmt.Errorf("generar cotización: %w", err)
	}

	cotizacion, err := construirCotizacion(s.cifrador, vehiculoID, clienteID, mensajeCliente, respuestaIA)
	if err != nil {
		return nil, err
	}

	creada, err := s.repositorioCotizaciones.Crear(ctx, cotizacion)
	if err != nil {
		return nil, err
	}
	if err := s.descifrarMensajes(creada); err != nil {
		return nil, err
	}
	return creada, nil
}

// ListarPorCliente devuelve las cotizaciones con su último mensaje descifrado.
func (s *cotizacionService) ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Cotizacion, error) {
	cotizaciones, err := s.repositorioCotizaciones.ListarPorCliente(ctx, clienteID)
	if err != nil {
		return nil, fmt.Errorf("listar cotizaciones: %w", err)
	}
	for i := range cotizaciones {
		if err := s.descifrarMensajes(&cotizaciones[i]); err != nil {
			return nil, err
		}
	}
	return cotizaciones, nil
}

// ObtenerPorCliente verifica que la cotización sea del cliente y devuelve sus
// mensajes descifrados.
func (s *cotizacionService) ObtenerPorCliente(ctx context.Context, clienteID uint, cotizacionID uint) (*models.Cotizacion, error) {
	cotizacion, err := s.repositorioCotizaciones.ObtenerPorID(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCotizacionNoEncontrada
		}
		return nil, fmt.Errorf("obtener cotización: %w", err)
	}

	if cotizacion.ClienteID != clienteID {
		return nil, ErrCotizacionNoPertenece
	}

	if err := s.descifrarMensajes(cotizacion); err != nil {
		return nil, err
	}
	return cotizacion, nil
}

// EnviarMensaje guarda el mensaje del cliente cifrado, enciende la respuesta
// del asistente con el contexto del vehículo y la guarda también cifrada.
func (s *cotizacionService) EnviarMensaje(ctx context.Context, clienteID uint, cotizacionID uint, mensaje string) (string, error) {
	mensajeCliente := strings.TrimSpace(mensaje)
	if mensajeCliente == "" {
		return "", ErrMensajeVacio
	}
	if len([]rune(mensajeCliente)) > LargoMaximoMensaje {
		return "", ErrMensajeMuyLargo
	}

	cotizacion, err := s.repositorioCotizaciones.ObtenerPorID(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrCotizacionNoEncontrada
		}
		return "", fmt.Errorf("obtener cotización: %w", err)
	}

	if cotizacion.ClienteID != clienteID {
		return "", ErrCotizacionNoPertenece
	}
	if cotizacion.Estado != models.EstadoCotizacionAbierta {
		return "", ErrCotizacionCerradaMensajes
	}

	if err := s.descifrarMensajes(cotizacion); err != nil {
		return "", err
	}
	historial := aTurnosChat(cotizacion.Mensajes)

	mensajeCifrado, err := s.cifrador.Cifrar(mensajeCliente)
	if err != nil {
		return "", fmt.Errorf("cifrar mensaje de cotización: %w", err)
	}
	if err := s.repositorioCotizaciones.AgregarMensaje(ctx, &models.MensajeCotizacion{
		CotizacionID: cotizacionID,
		Remitente:    models.RemitenteCliente,
		Contenido:    mensajeCifrado,
	}); err != nil {
		return "", err
	}

	// Si la cotización está atendida por un vendedor, la IA queda silenciada:
	// el mensaje se guarda y espera la respuesta del personal.
	if cotizacion.VendedorID != nil {
		return "", nil
	}

	respuestaIA, err := s.generador.GenerarCotizacion(ctx, cotizacion.Vehiculo, historial, mensajeCliente)
	if err != nil {
		return "", fmt.Errorf("generar respuesta de cotización: %w", err)
	}

	respuestaCifrada, err := s.cifrador.Cifrar(respuestaIA)
	if err != nil {
		return "", fmt.Errorf("cifrar respuesta de cotización: %w", err)
	}
	if err := s.repositorioCotizaciones.AgregarMensaje(ctx, &models.MensajeCotizacion{
		CotizacionID: cotizacionID,
		Remitente:    models.RemitenteIA,
		Contenido:    respuestaCifrada,
	}); err != nil {
		return "", err
	}

	return respuestaIA, nil
}

// Cerrar cambia el estado de la cotización del cliente a cerrada.
func (s *cotizacionService) Cerrar(ctx context.Context, clienteID uint, cotizacionID uint) (*models.Cotizacion, error) {
	cotizacion, err := s.repositorioCotizaciones.ObtenerPorID(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCotizacionNoEncontrada
		}
		return nil, fmt.Errorf("obtener cotización: %w", err)
	}

	if cotizacion.ClienteID != clienteID {
		return nil, ErrCotizacionNoPertenece
	}
	if cotizacion.Estado != models.EstadoCotizacionAbierta {
		return nil, ErrCotizacionYaCerrada
	}

	cotizacion.Estado = models.EstadoCotizacionCerrada
	if err := s.repositorioCotizaciones.Actualizar(ctx, cotizacion); err != nil {
		return nil, fmt.Errorf("cerrar cotización: %w", err)
	}
	if err := s.descifrarMensajes(cotizacion); err != nil {
		return nil, err
	}
	return cotizacion, nil
}

// ListarBandeja devuelve todas las cotizaciones con su último mensaje
// descifrado, para la bandeja del personal.
func (s *cotizacionService) ListarBandeja(ctx context.Context) ([]models.Cotizacion, error) {
	cotizaciones, err := s.repositorioCotizaciones.ListarBandeja(ctx)
	if err != nil {
		return nil, fmt.Errorf("listar bandeja de cotizaciones: %w", err)
	}
	for i := range cotizaciones {
		if err := s.descifrarMensajes(&cotizaciones[i]); err != nil {
			return nil, err
		}
	}
	return cotizaciones, nil
}

// ObtenerPersonal devuelve cualquier cotización con sus mensajes descifrados,
// para la vista de atención del vendedor.
func (s *cotizacionService) ObtenerPersonal(ctx context.Context, cotizacionID uint) (*models.Cotizacion, error) {
	cotizacion, err := s.repositorioCotizaciones.ObtenerPorID(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCotizacionNoEncontrada
		}
		return nil, fmt.Errorf("obtener cotización: %w", err)
	}

	if err := s.descifrarMensajes(cotizacion); err != nil {
		return nil, err
	}
	return cotizacion, nil
}

// ObtenerMensajesDesde devuelve el delta de mensajes del hilo junto con su
// total y cabecera. Para un cliente valida que la cotización sea propia; para
// el personal no hay restricción de asignación (el router ya exige rol de
// vendedor). Los mensajes se devuelven ya descifrados.
func (s *cotizacionService) ObtenerMensajesDesde(ctx context.Context, usuarioID uint, rol string, cotizacionID uint, desdeID uint) (*EstadoConversacionCotizacion, error) {
	cotizacion, err := s.repositorioCotizaciones.ObtenerCabecera(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCotizacionNoEncontrada
		}
		return nil, fmt.Errorf("obtener cotización: %w", err)
	}
	if rol == models.RolCliente && cotizacion.ClienteID != usuarioID {
		return nil, ErrCotizacionNoPertenece
	}

	mensajes, err := s.repositorioCotizaciones.ObtenerMensajesDesde(ctx, cotizacionID, desdeID)
	if err != nil {
		return nil, fmt.Errorf("obtener mensajes nuevos de cotización: %w", err)
	}
	for i := range mensajes {
		contenido, err := s.cifrador.Descifrar(mensajes[i].Contenido)
		if err != nil {
			return nil, fmt.Errorf("descifrar mensaje de cotización: %w", err)
		}
		mensajes[i].Contenido = contenido
	}

	total, err := s.repositorioCotizaciones.ContarMensajes(ctx, cotizacionID)
	if err != nil {
		return nil, fmt.Errorf("contar mensajes de cotización: %w", err)
	}

	return &EstadoConversacionCotizacion{
		Mensajes:  mensajes,
		Total:     total,
		Estado:    cotizacion.Estado,
		Vendedor:  cotizacion.Vendedor,
		FechaToma: cotizacion.FechaToma,
	}, nil
}

// Tomar asigna la cotización al vendedor autenticado. Es idempotente si la
// tomó el mismo vendedor; falla si otro vendedor la tiene o está cerrada.
func (s *cotizacionService) Tomar(ctx context.Context, vendedorID uint, cotizacionID uint) (*models.Cotizacion, error) {
	cotizacion, err := s.repositorioCotizaciones.ObtenerPorID(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCotizacionNoEncontrada
		}
		return nil, fmt.Errorf("obtener cotización: %w", err)
	}

	if cotizacion.Estado == models.EstadoCotizacionCerrada {
		return nil, ErrCotizacionYaCerrada
	}
	if cotizacion.VendedorID != nil {
		if *cotizacion.VendedorID != vendedorID {
			return nil, ErrCotizacionYaAtendida
		}
		if err := s.descifrarMensajes(cotizacion); err != nil {
			return nil, err
		}
		return cotizacion, nil
	}

	now := time.Now()
	cotizacion.VendedorID = &vendedorID
	cotizacion.FechaToma = &now
	if err := s.repositorioCotizaciones.Actualizar(ctx, cotizacion); err != nil {
		return nil, fmt.Errorf("tomar cotización: %w", err)
	}
	if err := s.descifrarMensajes(cotizacion); err != nil {
		return nil, err
	}
	return cotizacion, nil
}

// ResponderComoVendedor guarda el mensaje del vendedor cifrado sin pasar por
// la IA: al estar la cotización atendida, el asistente no genera respuestas.
func (s *cotizacionService) ResponderComoVendedor(ctx context.Context, vendedorID uint, cotizacionID uint, mensaje string) (*models.Cotizacion, error) {
	mensajeVendedor := strings.TrimSpace(mensaje)
	if mensajeVendedor == "" {
		return nil, ErrMensajeVacio
	}
	if len([]rune(mensajeVendedor)) > LargoMaximoMensaje {
		return nil, ErrMensajeMuyLargo
	}

	cotizacion, err := s.repositorioCotizaciones.ObtenerPorID(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCotizacionNoEncontrada
		}
		return nil, fmt.Errorf("obtener cotización: %w", err)
	}

	switch {
	case cotizacion.VendedorID == nil:
		return nil, ErrCotizacionNoTomada
	case *cotizacion.VendedorID != vendedorID:
		return nil, ErrCotizacionYaAtendida
	case cotizacion.Estado != models.EstadoCotizacionAbierta:
		return nil, ErrCotizacionYaCerrada
	}

	mensajeCifrado, err := s.cifrador.Cifrar(mensajeVendedor)
	if err != nil {
		return nil, fmt.Errorf("cifrar mensaje de cotización: %w", err)
	}
	if err := s.repositorioCotizaciones.AgregarMensaje(ctx, &models.MensajeCotizacion{
		CotizacionID: cotizacionID,
		Remitente:    models.RemitenteVendedor,
		Contenido:    mensajeCifrado,
	}); err != nil {
		return nil, fmt.Errorf("guardar mensaje de vendedor: %w", err)
	}

	return s.ObtenerPersonal(ctx, cotizacionID)
}

// CerrarPersonal cierra una cotización abierta desde la bandeja del personal.
// El administrador puede cerrarla aunque otro vendedor la tenga tomada.
func (s *cotizacionService) CerrarPersonal(ctx context.Context, cotizacionID uint) (*models.Cotizacion, error) {
	cotizacion, err := s.repositorioCotizaciones.ObtenerPorID(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCotizacionNoEncontrada
		}
		return nil, fmt.Errorf("obtener cotización: %w", err)
	}

	if cotizacion.Estado == models.EstadoCotizacionCerrada {
		return nil, ErrCotizacionYaCerrada
	}

	cotizacion.Estado = models.EstadoCotizacionCerrada
	if err := s.repositorioCotizaciones.Actualizar(ctx, cotizacion); err != nil {
		return nil, fmt.Errorf("cerrar cotización: %w", err)
	}
	if err := s.descifrarMensajes(cotizacion); err != nil {
		return nil, err
	}
	return cotizacion, nil
}

// descifrarMensajes reemplaza el contenido cifrado de cada mensaje por el texto
// original. El handler arma los DTO leyendo este campo descifrado.
func (s *cotizacionService) descifrarMensajes(cotizacion *models.Cotizacion) error {
	for i := range cotizacion.Mensajes {
		contenido, err := s.cifrador.Descifrar(cotizacion.Mensajes[i].Contenido)
		if err != nil {
			return fmt.Errorf("descifrar mensaje de cotización: %w", err)
		}
		cotizacion.Mensajes[i].Contenido = contenido
	}
	return nil
}

// construirCotizacion arma una cotización con sus dos mensajes iniciales
// (cliente e IA) ya cifrados. Se comparte con el chatbot para crear el registro
// desde el chat general sin duplicar la lógica de cifrado.
func construirCotizacion(cifrador cifrado.Cifrador, vehiculoID uint, clienteID uint, mensajeCliente string, respuestaIA string) (*models.Cotizacion, error) {
	mensajeCifrado, err := cifrador.Cifrar(mensajeCliente)
	if err != nil {
		return nil, fmt.Errorf("cifrar mensaje de cotización: %w", err)
	}
	respuestaCifrada, err := cifrador.Cifrar(respuestaIA)
	if err != nil {
		return nil, fmt.Errorf("cifrar respuesta de cotización: %w", err)
	}

	return &models.Cotizacion{
		VehiculoID: vehiculoID,
		ClienteID:  clienteID,
		Estado:     models.EstadoCotizacionAbierta,
		Mensajes: []models.MensajeCotizacion{
			{Remitente: models.RemitenteCliente, Contenido: mensajeCifrado},
			{Remitente: models.RemitenteIA, Contenido: respuestaCifrada},
		},
	}, nil
}

// aTurnosChat convierte los mensajes (ya descifrados) de una cotización en el
// turnos que entiende el chat del asistente. Solo se pasan al LLM los últimos
// MaximoTurnosHistorial turnos: el historial completo se conserva en la base y
// se muestra en la UI, pero el modelo no necesita recordar la conversación
// entera (ver docs/roadmap.md "Escalabilidad de conversaciones").
func aTurnosChat(mensajes []models.MensajeCotizacion) []TurnoChat {
	mensajes = ultimosTurnos(mensajes, MaximoTurnosHistorial)
	turnos := make([]TurnoChat, 0, len(mensajes))
	for _, mensaje := range mensajes {
		rol := "asistente"
		if mensaje.Remitente == models.RemitenteCliente {
			rol = "usuario"
		}
		turnos = append(turnos, TurnoChat{Rol: rol, Contenido: mensaje.Contenido})
	}
	return turnos
}

// ultimosTurnos devuelve el segmento final de N turnos (la conversación viene
// ordenada de más vieja a más nueva).
func ultimosTurnos[T any](turnos []T, cantidad int) []T {
	if len(turnos) <= cantidad {
		return turnos
	}
	return turnos[len(turnos)-cantidad:]
}

// ContarNoLeidos devuelve los mensajes de cotizaciones sin leer del usuario
// según su rol: para un cliente cuentan las respuestas de la IA o del
// vendedor; para el personal, los mensajes de cliente en cotizaciones abiertas
// propias o sin asignar.
func (s *cotizacionService) ContarNoLeidos(ctx context.Context, usuarioID uint, rol string) (int64, error) {
	if rol == models.RolCliente {
		return s.repositorioCotizaciones.ContarNoLeidosDeCliente(ctx, usuarioID)
	}
	return s.repositorioCotizaciones.ContarNoLeidosParaPersonal(ctx, usuarioID)
}

// MarcarLeidas marca como leídos los mensajes del hilo para el lado indicado.
// El lado cliente exige que la cotización pertenezca al usuario; el lado
// personal solo marca cuando el usuario es el vendedor asignado (si no hay
// asignación o es de otro vendedor, no hace nada).
func (s *cotizacionService) MarcarLeidas(ctx context.Context, usuarioID uint, cotizacionID uint, lado string) error {
	cotizacion, err := s.repositorioCotizaciones.ObtenerPorID(ctx, cotizacionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCotizacionNoEncontrada
		}
		return fmt.Errorf("obtener cotización: %w", err)
	}

	switch lado {
	case LadoCliente:
		if cotizacion.ClienteID != usuarioID {
			return ErrCotizacionNoPertenece
		}
		return s.repositorioCotizaciones.MarcarLeidasParaCliente(ctx, cotizacionID)
	case LadoPersonal:
		if cotizacion.VendedorID == nil || *cotizacion.VendedorID != usuarioID {
			return nil
		}
		return s.repositorioCotizaciones.MarcarLeidasParaPersonal(ctx, cotizacionID)
	default:
		return fmt.Errorf("lado de lectura inválido: %s", lado)
	}
}
