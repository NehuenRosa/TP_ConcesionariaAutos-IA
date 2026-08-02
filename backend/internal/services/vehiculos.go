package services

import (
	"context"
	"errors"
	"fmt"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

// Errores de negocio del catálogo de vehículos.
var (
	// ErrVehiculoNoEncontrado indica que el vehículo no existe o no está disponible.
	ErrVehiculoNoEncontrado = errors.New("vehículo no encontrado o no disponible")
	// ErrPaginacionInvalida indica que la página o el tamaño solicitado no son válidos.
	ErrPaginacionInvalida = errors.New("paginación inválida")
)

// VehiculoService define el contrato de la lógica de negocio de vehículos.
type VehiculoService interface {
	ListarDisponibles(ctx context.Context, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error)
}

// vehiculoService implementa VehiculoService.
type vehiculoService struct {
	repositorio repositories.VehiculoRepository
}

// NuevoVehiculoService crea un servicio de vehículos.
func NuevoVehiculoService(repositorio repositories.VehiculoRepository) VehiculoService {
	return &vehiculoService{repositorio: repositorio}
}

// ListarDisponibles valida la paginación y delega en el repositorio el listado
// de vehículos con estado "disponible".
func (s *vehiculoService) ListarDisponibles(ctx context.Context, pagina int, tamano int) ([]models.Vehiculo, int64, error) {
	if pagina < 1 || tamano < 1 {
		return nil, 0, ErrPaginacionInvalida
	}
	return s.repositorio.Listar(ctx, models.EstadoDisponible, pagina, tamano)
}

// ObtenerPorID devuelve un vehículo solo si existe y está disponible.
// Para los demás estados se retorna ErrVehiculoNoEncontrado, ocultando la
// existencia de unidades no comercializables.
func (s *vehiculoService) ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error) {
	vehiculo, err := s.repositorio.ObtenerPorID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehiculoNoEncontrado
		}
		return nil, fmt.Errorf("obtener vehículo por ID: %w", err)
	}
	if vehiculo.Estado != models.EstadoDisponible {
		return nil, ErrVehiculoNoEncontrado
	}
	return vehiculo, nil
}
