package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

// MensajeService define el contrato de la lógica de negocio de mensajes.
type MensajeService interface {
	// Enviar envía un mensaje en una consulta.
	Enviar(ctx context.Context, consultaID uint, remitenteID uint, contenido string) (*models.Mensaje, error)
	// ObtenerPorConsulta obtiene todos los mensajes de una consulta.
	ObtenerPorConsulta(ctx context.Context, consultaID uint, usuarioID uint) ([]models.Mensaje, error)
	// ObtenerDesdeID obtiene los mensajes de la consulta con id mayor al
	// indicado y los marca como leídos. Reemplaza el recorte por timestamp:
	// como los timestamps se serializan sin sub-segundos, dos mensajes del
	// mismo segundo podían perderse.
	ObtenerDesdeID(ctx context.Context, consultaID uint, usuarioID uint, desdeID uint) ([]models.Mensaje, error)
	// ContarNoLeidosPorConsultas cuenta los mensajes no leídos de un usuario
	// en varias consultas.
	ContarNoLeidosPorConsultas(ctx context.Context, consultaIDs []uint, usuarioID uint) (map[uint]int, error)
	// MarcarComoLeidos marca como leídos los mensajes de otros en una consulta.
	MarcarComoLeidos(ctx context.Context, consultaID uint, usuarioID uint) error
}

// mensajeService implementa MensajeService.
type mensajeService struct {
	repositorioMensajes  repositories.MensajeRepository
	repositorioConsultas repositories.ConsultaRepository
	consultaService     ConsultaService
}

// NuevoMensajeService crea un servicio de mensajes.
func NuevoMensajeService(
	repositorioMensajes repositories.MensajeRepository,
	repositorioConsultas repositories.ConsultaRepository,
	consultaService ConsultaService,
) MensajeService {
	return &mensajeService{
		repositorioMensajes:  repositorioMensajes,
		repositorioConsultas: repositorioConsultas,
		consultaService:     consultaService,
	}
}

// Enviar envía un mensaje en una consulta.
func (s *mensajeService) Enviar(ctx context.Context, consultaID uint, remitenteID uint, contenido string) (*models.Mensaje, error) {
	if strings.TrimSpace(contenido) == "" {
		return nil, ErrMensajeVacio
	}

	// Verificar que la consulta existe
	consulta, err := s.repositorioConsultas.ObtenerPorID(ctx, consultaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConsultaNoEncontrada
		}
		return nil, fmt.Errorf("obtener consulta: %w", err)
	}

	// Verificar que la consulta no esté cerrada
	if consulta.Estado == models.EstadoCerrada {
		return nil, ErrConsultaCerradaMensajes
	}

	// Verificar que el remitente es participante
	esParticipante, err := s.consultaService.EsParticipante(ctx, consultaID, remitenteID)
	if err != nil {
		return nil, fmt.Errorf("verificar participante: %w", err)
	}
	if !esParticipante {
		return nil, ErrNoEsParticipante
	}

	mensaje := &models.Mensaje{
		ConsultaID:  consultaID,
		RemitenteID: remitenteID,
		Contenido:   strings.TrimSpace(contenido),
	}

	return s.repositorioMensajes.Crear(ctx, mensaje)
}

// ObtenerPorConsulta obtiene todos los mensajes de una consulta.
func (s *mensajeService) ObtenerPorConsulta(ctx context.Context, consultaID uint, usuarioID uint) ([]models.Mensaje, error) {
	// Verificar que el usuario es participante
	esParticipante, err := s.consultaService.EsParticipante(ctx, consultaID, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("verificar participante: %w", err)
	}
	if !esParticipante {
		return nil, ErrNoEsParticipante
	}

	return s.repositorioMensajes.ListarPorConsulta(ctx, consultaID)
}

// ObtenerDesdeID obtiene los mensajes de la consulta posteriores a desdeID y
// marca como leídos los del otro participante.
func (s *mensajeService) ObtenerDesdeID(ctx context.Context, consultaID uint, usuarioID uint, desdeID uint) ([]models.Mensaje, error) {
	// Verificar que el usuario es participante
	esParticipante, err := s.consultaService.EsParticipante(ctx, consultaID, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("verificar participante: %w", err)
	}
	if !esParticipante {
		return nil, ErrNoEsParticipante
	}

	mensajes, err := s.repositorioMensajes.ObtenerDesdeID(ctx, consultaID, desdeID)
	if err != nil {
		return nil, fmt.Errorf("obtener mensajes nuevos: %w", err)
	}

	// Marcar como leídos los mensajes del otro participante
	if err := s.repositorioMensajes.MarcarComoLeidos(ctx, consultaID, usuarioID); err != nil {
		return nil, fmt.Errorf("marcar mensajes como leídos: %w", err)
	}

	return mensajes, nil
}

// ContarNoLeidosPorConsultas cuenta los mensajes no leídos de un usuario
// en varias consultas.
func (s *mensajeService) ContarNoLeidosPorConsultas(ctx context.Context, consultaIDs []uint, usuarioID uint) (map[uint]int, error) {
	return s.repositorioMensajes.ContarNoLeidosPorConsultas(ctx, consultaIDs, usuarioID)
}

// MarcarComoLeidos marca como leídos los mensajes de otros en una consulta.
func (s *mensajeService) MarcarComoLeidos(ctx context.Context, consultaID uint, usuarioID uint) error {
	esParticipante, err := s.consultaService.EsParticipante(ctx, consultaID, usuarioID)
	if err != nil {
		return fmt.Errorf("verificar participante: %w", err)
	}
	if !esParticipante {
		return ErrNoEsParticipante
	}

	return s.repositorioMensajes.MarcarComoLeidos(ctx, consultaID, usuarioID)
}
