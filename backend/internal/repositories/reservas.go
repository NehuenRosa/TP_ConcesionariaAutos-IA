package repositories

import (
	"context"
	"errors"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// Errores de datos de las reservas ante condiciones de carrera.
var (
	// ErrVehiculoYaNoDisponible indica que la unidad dejó de estar disponible
	// al momento de persistir la reserva (otro proceso la tomó primero).
	ErrVehiculoYaNoDisponible = errors.New("el vehículo ya no está disponible")
	// ErrReservaYaNoActiva indica que la reserva dejó de estar activa al
	// momento de persistir la transición (otro proceso la cambió primero).
	ErrReservaYaNoActiva = errors.New("la reserva ya no está activa")
)

// ReservaRepository define el acceso a datos de reservas, incluyendo las
// operaciones que cambian el estado del vehículo en la misma transacción.
type ReservaRepository interface {
	// CrearYReservar crea la reserva en estado activa y bloquea el vehículo
	// (disponible → reservado) de forma atómica.
	CrearYReservar(ctx context.Context, reserva *models.Reserva) (*models.Reserva, error)
	// ObtenerPorID devuelve una reserva con su vehículo, imágenes y cliente.
	ObtenerPorID(ctx context.Context, id uint) (*models.Reserva, error)
	// ListarPorCliente devuelve las reservas de un cliente, ordenadas por fecha
	// de creación.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Reserva, error)
	// Listar devuelve las reservas con el filtro de estado opcional, ordenadas
	// por fecha de creación.
	Listar(ctx context.Context, estado string) ([]models.Reserva, error)
	// ConfirmarVentaYMarcarVendido cambia la reserva activa a vendida y el
	// vehículo a vendido de forma atómica.
	ConfirmarVentaYMarcarVendido(ctx context.Context, reserva *models.Reserva) (*models.Reserva, error)
	// CancelarYLiberar cambia la reserva activa a cancelada y libera el
	// vehículo (vuelve a disponible) de forma atómica.
	CancelarYLiberar(ctx context.Context, reserva *models.Reserva) (*models.Reserva, error)
}

// reservaRepository implementa ReservaRepository sobre GORM.
type reservaRepository struct {
	base *gorm.DB
}

// NuevoReservaRepository crea un repositorio de reservas.
func NuevoReservaRepository(base *gorm.DB) ReservaRepository {
	return &reservaRepository{base: base}
}

// CrearYReservar crea la reserva y bloquea el vehículo en una única
// transacción. Si la unidad dejó de estar disponible, revierte y devuelve
// ErrVehiculoYaNoDisponible.
func (r *reservaRepository) CrearYReservar(ctx context.Context, reserva *models.Reserva) (*models.Reserva, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reserva).Error; err != nil {
			return err
		}

		resultado := tx.Model(&models.Vehiculo{}).
			Where("id = ? AND estado = ?", reserva.VehiculoID, models.EstadoDisponible).
			Update("estado", models.EstadoReservado)
		if resultado.Error != nil {
			return resultado.Error
		}
		if resultado.RowsAffected == 0 {
			return ErrVehiculoYaNoDisponible
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.ObtenerPorID(ctx, reserva.ID)
}

// ObtenerPorID devuelve una reserva con su vehículo, imágenes y cliente.
func (r *reservaRepository) ObtenerPorID(ctx context.Context, id uint) (*models.Reserva, error) {
	var reserva models.Reserva
	err := r.base.WithContext(ctx).
		Preload("Vehiculo").
		Preload("Vehiculo.Imagenes").
		Preload("Cliente").
		First(&reserva, id).Error
	if err != nil {
		return nil, err
	}
	return &reserva, nil
}

// ListarPorCliente devuelve las reservas de un cliente, ordenadas por fecha de
// creación.
func (r *reservaRepository) ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Reserva, error) {
	var reservas []models.Reserva
	err := r.base.WithContext(ctx).
		Where("cliente_id = ?", clienteID).
		Preload("Vehiculo").
		Preload("Vehiculo.Imagenes").
		Order("created_at DESC").
		Find(&reservas).Error
	if err != nil {
		return nil, fmt.Errorf("listar reservas por cliente: %w", err)
	}
	return reservas, nil
}

// Listar devuelve las reservas con el filtro de estado opcional, ordenadas por
// fecha de creación.
func (r *reservaRepository) Listar(ctx context.Context, estado string) ([]models.Reserva, error) {
	consulta := r.base.WithContext(ctx)
	if estado != "" {
		consulta = consulta.Where("estado = ?", estado)
	}

	var reservas []models.Reserva
	err := consulta.
		Preload("Vehiculo").
		Preload("Vehiculo.Imagenes").
		Preload("Cliente").
		Order("created_at DESC").
		Find(&reservas).Error
	if err != nil {
		return nil, fmt.Errorf("listar reservas: %w", err)
	}
	return reservas, nil
}

// ConfirmarVentaYMarcarVendido cambia la reserva activa a vendida y el
// vehículo a vendido en una única transacción. Si la reserva ya no está
// activa, devuelve ErrReservaYaNoActiva.
func (r *reservaRepository) ConfirmarVentaYMarcarVendido(ctx context.Context, reserva *models.Reserva) (*models.Reserva, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		resultado := tx.Model(&models.Reserva{}).
			Where("id = ? AND estado = ?", reserva.ID, models.EstadoReservaActiva).
			Update("estado", models.EstadoReservaVendida)
		if resultado.Error != nil {
			return resultado.Error
		}
		if resultado.RowsAffected == 0 {
			return ErrReservaYaNoActiva
		}

		if err := tx.Model(&models.Vehiculo{}).
			Where("id = ?", reserva.VehiculoID).
			Update("estado", models.EstadoVendido).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.ObtenerPorID(ctx, reserva.ID)
}

// CancelarYLiberar cambia la reserva activa a cancelada y libera el vehículo
// (vuelve a disponible) en una única transacción. Si la reserva ya no está
// activa, devuelve ErrReservaYaNoActiva.
func (r *reservaRepository) CancelarYLiberar(ctx context.Context, reserva *models.Reserva) (*models.Reserva, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		resultado := tx.Model(&models.Reserva{}).
			Where("id = ? AND estado = ?", reserva.ID, models.EstadoReservaActiva).
			Update("estado", models.EstadoReservaCancelada)
		if resultado.Error != nil {
			return resultado.Error
		}
		if resultado.RowsAffected == 0 {
			return ErrReservaYaNoActiva
		}

		if err := tx.Model(&models.Vehiculo{}).
			Where("id = ?", reserva.VehiculoID).
			Update("estado", models.EstadoDisponible).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.ObtenerPorID(ctx, reserva.ID)
}
