package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
)

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

	return s.repositorioCotizaciones.Crear(ctx, cotizacion)
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
// turnos que entiende el chat del asistente.
func aTurnosChat(mensajes []models.MensajeCotizacion) []TurnoChat {
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
