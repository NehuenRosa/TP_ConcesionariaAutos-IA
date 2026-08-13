package services

import (
	"context"
	"errors"
	"fmt"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

// Errores de negocio de las reservas.
var (
	// ErrReservaNoEncontrada indica que la reserva no existe o no pertenece al
	// cliente.
	ErrReservaNoEncontrada = errors.New("reserva no encontrada")
	// ErrVehiculoYaNoDisponible indica que la unidad ya no está disponible para
	// reservar (fue reservada o vendida por otro cliente).
	ErrVehiculoYaNoDisponible = errors.New("la unidad ya no está disponible para reservar")
	// ErrReservaEstadoInvalido indica que la transición de estado solicitada no
	// es válida.
	ErrReservaEstadoInvalido = errors.New("no se puede cambiar el estado de la reserva")
	// ErrFiltroEstadoReservaInvalido indica que el filtro de estado no es válido.
	ErrFiltroEstadoReservaInvalido = errors.New("filtro de estado inválido")
)

// ReservaService define el contrato de la lógica de negocio de reservas.
type ReservaService interface {
	// Crear reserva un vehículo disponible: crea la reserva activa y bloquea la
	// unidad.
	Crear(ctx context.Context, clienteID uint, vehiculoID uint) (*models.Reserva, error)
	// ListarMisReservas lista las reservas de un cliente.
	ListarMisReservas(ctx context.Context, clienteID uint) ([]models.Reserva, error)
	// Cancelar cancela una reserva propia en estado activa y libera la unidad.
	Cancelar(ctx context.Context, reservaID uint, clienteID uint) (*models.Reserva, error)
	// Listar lista las reservas para el vendedor, con filtro de estado opcional.
	Listar(ctx context.Context, estado string) ([]models.Reserva, error)
	// ConfirmarVenta confirma la venta de una reserva activa.
	ConfirmarVenta(ctx context.Context, reservaID uint) (*models.Reserva, error)
	// CancelarComoVendedor cancela una reserva activa y libera la unidad.
	CancelarComoVendedor(ctx context.Context, reservaID uint) (*models.Reserva, error)
}

// reservaService implementa ReservaService.
type reservaService struct {
	repositorio repositories.ReservaRepository
	vehiculos   repositories.VehiculoRepository
}

// NuevoReservaService crea un servicio de reservas.
func NuevoReservaService(
	repositorio repositories.ReservaRepository,
	vehiculos repositories.VehiculoRepository,
) ReservaService {
	return &reservaService{
		repositorio: repositorio,
		vehiculos:   vehiculos,
	}
}

// Crear valida que el vehículo exista y esté disponible y crea la reserva en
// estado activa, bloqueando la unidad.
func (s *reservaService) Crear(ctx context.Context, clienteID uint, vehiculoID uint) (*models.Reserva, error) {
	vehiculo, err := s.vehiculos.ObtenerPorID(ctx, vehiculoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehiculoNoDisponible
		}
		return nil, fmt.Errorf("obtener vehículo: %w", err)
	}

	switch vehiculo.Estado {
	case models.EstadoDisponible:
		// La unidad puede reservarse.
	case models.EstadoReservado, models.EstadoVendido:
		return nil, ErrVehiculoYaNoDisponible
	default:
		// Dado de baja: no es comercializable.
		return nil, ErrVehiculoNoDisponible
	}

	reserva := &models.Reserva{
		VehiculoID: vehiculoID,
		ClienteID:  clienteID,
		Estado:     models.EstadoReservaActiva,
	}
	reservaCreada, err := s.repositorio.CrearYReservar(ctx, reserva)
	if err != nil {
		if errors.Is(err, repositories.ErrVehiculoYaNoDisponible) {
			return nil, ErrVehiculoYaNoDisponible
		}
		return nil, fmt.Errorf("crear reserva: %w", err)
	}
	return reservaCreada, nil
}

// ListarMisReservas lista las reservas de un cliente.
func (s *reservaService) ListarMisReservas(ctx context.Context, clienteID uint) ([]models.Reserva, error) {
	return s.repositorio.ListarPorCliente(ctx, clienteID)
}

// Cancelar cancela una reserva propia en estado activa y libera la unidad. Las
// reservas ajenas se tratan como inexistentes para no revelar su existencia.
func (s *reservaService) Cancelar(ctx context.Context, reservaID uint, clienteID uint) (*models.Reserva, error) {
	reserva, err := s.repositorio.ObtenerPorID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservaNoEncontrada
		}
		return nil, fmt.Errorf("obtener reserva: %w", err)
	}
	if reserva.ClienteID != clienteID {
		return nil, ErrReservaNoEncontrada
	}
	if !reserva.EsActiva() {
		return nil, ErrReservaEstadoInvalido
	}

	reservaCancelada, err := s.repositorio.CancelarYLiberar(ctx, reserva)
	if err != nil {
		if errors.Is(err, repositories.ErrReservaYaNoActiva) {
			return nil, ErrReservaEstadoInvalido
		}
		return nil, fmt.Errorf("cancelar reserva: %w", err)
	}
	return reservaCancelada, nil
}

// Listar lista las reservas para el vendedor con filtro de estado opcional.
func (s *reservaService) Listar(ctx context.Context, estado string) ([]models.Reserva, error) {
	if estado != "" && !esEstadoReservaValido(estado) {
		return nil, ErrFiltroEstadoReservaInvalido
	}
	return s.repositorio.Listar(ctx, estado)
}

// ConfirmarVenta confirma la venta de una reserva en estado activa.
func (s *reservaService) ConfirmarVenta(ctx context.Context, reservaID uint) (*models.Reserva, error) {
	reserva, err := s.repositorio.ObtenerPorID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservaNoEncontrada
		}
		return nil, fmt.Errorf("obtener reserva: %w", err)
	}
	if !reserva.EsActiva() {
		return nil, ErrReservaEstadoInvalido
	}

	reservaVendida, err := s.repositorio.ConfirmarVentaYMarcarVendido(ctx, reserva)
	if err != nil {
		if errors.Is(err, repositories.ErrReservaYaNoActiva) {
			return nil, ErrReservaEstadoInvalido
		}
		return nil, fmt.Errorf("confirmar venta: %w", err)
	}
	return reservaVendida, nil
}

// CancelarComoVendedor cancela una reserva en estado activa y libera la unidad.
func (s *reservaService) CancelarComoVendedor(ctx context.Context, reservaID uint) (*models.Reserva, error) {
	reserva, err := s.repositorio.ObtenerPorID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservaNoEncontrada
		}
		return nil, fmt.Errorf("obtener reserva: %w", err)
	}
	if !reserva.EsActiva() {
		return nil, ErrReservaEstadoInvalido
	}

	reservaCancelada, err := s.repositorio.CancelarYLiberar(ctx, reserva)
	if err != nil {
		if errors.Is(err, repositories.ErrReservaYaNoActiva) {
			return nil, ErrReservaEstadoInvalido
		}
		return nil, fmt.Errorf("cancelar reserva: %w", err)
	}
	return reservaCancelada, nil
}

// esEstadoReservaValido indica si el estado es uno de los conocidos.
func esEstadoReservaValido(estado string) bool {
	switch estado {
	case models.EstadoReservaActiva, models.EstadoReservaVendida, models.EstadoReservaCancelada:
		return true
	default:
		return false
	}
}
