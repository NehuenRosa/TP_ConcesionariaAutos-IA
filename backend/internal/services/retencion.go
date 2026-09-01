package services

import (
	"context"
	"fmt"
	"time"

	"concesionaria/backend/internal/repositories"
)

// ResultadoRetencion resume lo purgado en una pasada del job de retención.
type ResultadoRetencion struct {
	// Cotizaciones y MensajesCotizaciones son las filas eliminadas de
	// cotizaciones cerradas viejas.
	Cotizaciones        int64
	MensajesCotizaciones int64
	// Consultas y MensajesConsultas son las filas eliminadas de consultas
	// cerradas viejas.
	Consultas         int64
	MensajesConsultas int64
}

// RetencionService define el contrato del job de conservación de datos.
type RetencionService interface {
	// Ejecutar purga las conversaciones cerradas que superan el plazo de
	// conservación (dias). Nunca toca conversaciones abiertas ni cerradas
	// recientes.
	Ejecutar(ctx context.Context, dias int) (*ResultadoRetencion, error)
}

// retencionService implementa RetencionService.
type retencionService struct {
	repositorio repositories.RetencionRepository
}

// NuevoRetencionService crea un servicio de retención.
func NuevoRetencionService(repositorio repositories.RetencionRepository) RetencionService {
	return &retencionService{repositorio: repositorio}
}

// Ejecutar purga las conversaciones cerradas anteriores al plazo indicado.
func (s *retencionService) Ejecutar(ctx context.Context, dias int) (*ResultadoRetencion, error) {
	if dias < 1 {
		return nil, fmt.Errorf("el plazo de retención debe ser mayor a 0 días")
	}

	corte := time.Now().AddDate(0, 0, -dias)

	mensajesCotizaciones, cotizaciones, err := s.repositorio.PurgarCotizacionesCerradas(ctx, corte)
	if err != nil {
		return nil, err
	}
	mensajesConsultas, consultas, err := s.repositorio.PurgarConsultasCerradas(ctx, corte)
	if err != nil {
		return nil, err
	}

	return &ResultadoRetencion{
		Cotizaciones:         cotizaciones,
		MensajesCotizaciones: mensajesCotizaciones,
		Consultas:            consultas,
		MensajesConsultas:    mensajesConsultas,
	}, nil
}