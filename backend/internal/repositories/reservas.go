package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// horaCero excluye del barrido de expiración a las reservas históricas cuyo
// vencimiento quedó en cero o nulo tras la migración.
var horaCero = time.Time{}

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
	// GuardarComprobante persiste la imagen del comprobante y marca su envío
	// en la reserva dentro de una misma transacción.
	GuardarComprobante(ctx context.Context, reserva *models.Reserva, comprobante *models.ComprobanteReserva) error
	// ObtenerComprobantePorReservaID devuelve el comprobante cargado de una
	// reserva, o gorm.ErrRecordNotFound si nunca se envió.
	ObtenerComprobantePorReservaID(ctx context.Context, reservaID uint) (*models.ComprobanteReserva, error)
	// ExpirarSiVencida cancela atómicamente la reserva si sigue activa,
	// pendiente de comprobante y fuera de plazo; libera el vehículo. Devuelve
	// true cuando aplicó la expiración.
	ExpirarSiVencida(ctx context.Context, reserva *models.Reserva) (bool, error)
	// ExpirarVencidas cancela todas las reservas activas vencidas sin
	// comprobante y libera sus vehículos. Es idempotente.
	ExpirarVencidas(ctx context.Context) (int64, error)
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
			Updates(map[string]any{
				"estado":             models.EstadoReservaCancelada,
				"motivo_cancelacion": reserva.MotivoCancelacion,
			})
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

// GuardarComprobante persiste la imagen del comprobante (creándola o
// reemplazando la anterior si el cliente reenvió) y marca la hora de envío en
// la reserva, todo dentro de una misma transacción.
func (r *reservaRepository) GuardarComprobante(ctx context.Context, reserva *models.Reserva, comprobante *models.ComprobanteReserva) error {
	return r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existente models.ComprobanteReserva
		err := tx.Where("reserva_id = ?", reserva.ID).First(&existente).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			comprobante.ReservaID = reserva.ID
			if err := tx.Create(comprobante).Error; err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			existente.MIME = comprobante.MIME
			existente.Datos = comprobante.Datos
			if err := tx.Save(&existente).Error; err != nil {
				return err
			}
			comprobante.ID = existente.ID
		}
		return tx.Model(&models.Reserva{}).
			Where("id = ?", reserva.ID).
			Update("comprobante_enviado_at", reserva.ComprobanteEnviadoAt).Error
	})
}

// ObtenerComprobantePorReservaID devuelve el comprobante cargado de una
// reserva, o gorm.ErrRecordNotFound si nunca se envió.
func (r *reservaRepository) ObtenerComprobantePorReservaID(ctx context.Context, reservaID uint) (*models.ComprobanteReserva, error) {
	var comprobante models.ComprobanteReserva
	err := r.base.WithContext(ctx).
		Where("reserva_id = ?", reservaID).
		First(&comprobante).Error
	if err != nil {
		return nil, err
	}
	return &comprobante, nil
}

// ExpirarSiVencida cancela la reserva solo si sigue activa, pendiente de
// comprobante y fuera de plazo (condiciones revalidadas sobre la base para
// evitar carreras con un envío simultáneo) y libera el vehículo.
func (r *reservaRepository) ExpirarSiVencida(ctx context.Context, reserva *models.Reserva) (bool, error) {
	expirada := false
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		resultado := condicionDeExpiracion(tx.Where("id = ?", reserva.ID)).
			Update("estado", models.EstadoReservaCancelada)
		if resultado.Error != nil {
			return resultado.Error
		}
		if resultado.RowsAffected == 0 {
			return nil
		}
		expirada = true
		return liberarVehiculoReservado(tx, reserva.VehiculoID)
	})
	return expirada, err
}

// ExpirarVencidas cancela todas las reservas activas vencidas sin comprobante
// y libera sus unidades. Además repara cualquier unidad que siga reservada por
// una reserva ya expirada (autocurativo e idempotente).
func (r *reservaRepository) ExpirarVencidas(ctx context.Context) (int64, error) {
	var cantidad int64
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ahora := time.Now()
		resultado := condicionDeExpiracion(tx.Model(&models.Reserva{})).
			Update("estado", models.EstadoReservaCancelada)
		if resultado.Error != nil {
			return resultado.Error
		}
		cantidad = resultado.RowsAffected

		vencidas := tx.Model(&models.Reserva{}).Select("vehiculo_id").
			Where("estado = ? AND comprobante_enviado_at IS NULL AND vencimiento_comprobante > ? AND vencimiento_comprobante < ?",
				models.EstadoReservaCancelada, horaCero, ahora)
		return tx.Model(&models.Vehiculo{}).
			Where("estado = ? AND id IN (?)", models.EstadoReservado, vencidas).
			Update("estado", models.EstadoDisponible).Error
	})
	return cantidad, err
}

// condicionDeExpiracion aplica el filtro de reserva expirable: activa, sin
// comprobante y con plazo vencido. El vencimiento debe ser posterior a la
// hora cero para no tocar reservas históricas.
func condicionDeExpiracion(consulta *gorm.DB) *gorm.DB {
	return consulta.Model(&models.Reserva{}).
		Where("estado = ? AND comprobante_enviado_at IS NULL AND vencimiento_comprobante > ? AND vencimiento_comprobante < ?",
			models.EstadoReservaActiva, horaCero, time.Now())
}

// liberarVehiculoReservado devuelve la unidad a disponible solo si sigue
// reservada.
func liberarVehiculoReservado(tx *gorm.DB, vehiculoID uint) error {
	return tx.Model(&models.Vehiculo{}).
		Where("id = ? AND estado = ?", vehiculoID, models.EstadoReservado).
		Update("estado", models.EstadoDisponible).Error
}
