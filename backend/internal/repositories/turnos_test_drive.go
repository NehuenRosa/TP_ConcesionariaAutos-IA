package repositories

import (
	"context"
	"errors"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrFranjaOcupada es el centinela que indica que ya existe un turno activo
// para la misma unidad, fecha y franja.
var ErrFranjaOcupada = errors.New("ya existe un turno superpuesto")

// ErrClienteYaTieneTurno es el centinela que indica que el cliente ya tiene un
// turno activo para el mismo vehículo y no puede pedir otro.
var ErrClienteYaTieneTurno = errors.New("el cliente ya tiene un turno para este vehículo")

// TurnoTestDriveRepository define el acceso a datos de turnos de test drive.
type TurnoTestDriveRepository interface {
	// CrearSiDisponible crea el turno en una transacción que bloquea la fila
	// del vehículo, garantizando ante pedidos concurrentes que no haya un
	// turno activo del mismo cliente para el vehículo (ErrClienteYaTieneTurno)
	// ni un turno activo para la unidad, fecha y franja (ErrFranjaOcupada).
	// Devuelve el turno creado o el error de la regla violada.
	CrearSiDisponible(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error)
	// ObtenerPorID devuelve un turno con su vehículo, imágenes y cliente.
	ObtenerPorID(ctx context.Context, id uint) (*models.TurnoTestDrive, error)
	// ListarPorCliente devuelve los turnos de un cliente, ordenados por fecha.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.TurnoTestDrive, error)
	// Listar devuelve los turnos con el filtro de estado opcional, ordenados
	// por fecha y franja.
	Listar(ctx context.Context, estado string) ([]models.TurnoTestDrive, error)
	// Ocupadas devuelve las franjas ya ocupadas (solicitado o confirmado) para
	// una unidad y fecha dadas.
	Ocupadas(ctx context.Context, vehiculoID uint, fecha string) ([]string, error)
	// Actualizar persiste los cambios de estado de un turno.
	Actualizar(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error)
}

// turnoTestDriveRepository implementa TurnoTestDriveRepository sobre GORM.
type turnoTestDriveRepository struct {
	base *gorm.DB
}

// NuevoTurnoTestDriveRepository crea un repositorio de turnos de test drive.
func NuevoTurnoTestDriveRepository(base *gorm.DB) TurnoTestDriveRepository {
	return &turnoTestDriveRepository{base: base}
}

// CrearSiDisponible crea el turno en una transacción con bloqueo de la fila
// del vehículo (FOR UPDATE): los pedidos concurrentes para el mismo vehículo
// se serializan y las dos reglas no se pueden violar (turno activo del mismo
// cliente para el vehículo, y superposición de unidad-fecha-franja).
func (r *turnoTestDriveRepository) CrearSiDisponible(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var vehiculo models.Vehiculo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&vehiculo, turno.VehiculoID).Error; err != nil {
			return err
		}

		var delCliente int64
		if err := tx.Model(&models.TurnoTestDrive{}).
			Where("vehiculo_id = ? AND cliente_id = ? AND estado IN ?",
				turno.VehiculoID, turno.ClienteID,
				[]string{models.EstadoTurnoSolicitado, models.EstadoTurnoConfirmado}).
			Count(&delCliente).Error; err != nil {
			return err
		}
		if delCliente > 0 {
			return ErrClienteYaTieneTurno
		}

		var total int64
		if err := tx.Model(&models.TurnoTestDrive{}).
			Where("vehiculo_id = ? AND fecha = ? AND franja = ? AND estado IN ?",
				turno.VehiculoID, turno.Fecha, turno.Franja,
				[]string{models.EstadoTurnoSolicitado, models.EstadoTurnoConfirmado}).
			Count(&total).Error; err != nil {
			return err
		}
		if total > 0 {
			return ErrFranjaOcupada
		}

		return tx.Create(turno).Error
	})
	if err != nil {
		return nil, err
	}

	creado, err := r.ObtenerPorID(ctx, turno.ID)
	if err != nil {
		return nil, fmt.Errorf("recuperar turno de test drive: %w", err)
	}
	return creado, nil
}

// ObtenerPorID devuelve un turno con su vehículo, imágenes y cliente.
func (r *turnoTestDriveRepository) ObtenerPorID(ctx context.Context, id uint) (*models.TurnoTestDrive, error) {
	var turno models.TurnoTestDrive
	err := r.base.WithContext(ctx).
		Preload("Vehiculo").
		Preload("Vehiculo.Imagenes").
		Preload("Cliente").
		First(&turno, id).Error
	if err != nil {
		return nil, err
	}
	return &turno, nil
}

// ListarPorCliente devuelve los turnos de un cliente, ordenados por fecha.
func (r *turnoTestDriveRepository) ListarPorCliente(ctx context.Context, clienteID uint) ([]models.TurnoTestDrive, error) {
	var turnos []models.TurnoTestDrive
	err := r.base.WithContext(ctx).
		Where("cliente_id = ?", clienteID).
		Preload("Vehiculo").
		Preload("Vehiculo.Imagenes").
		Order("fecha ASC, franja ASC").
		Find(&turnos).Error
	if err != nil {
		return nil, fmt.Errorf("listar turnos por cliente: %w", err)
	}
	return turnos, nil
}

// Listar devuelve los turnos con el filtro de estado opcional, ordenados por
// fecha y franja.
func (r *turnoTestDriveRepository) Listar(ctx context.Context, estado string) ([]models.TurnoTestDrive, error) {
	consulta := r.base.WithContext(ctx)
	if estado != "" {
		consulta = consulta.Where("estado = ?", estado)
	}

	var turnos []models.TurnoTestDrive
	err := consulta.
		Preload("Vehiculo").
		Preload("Vehiculo.Imagenes").
		Preload("Cliente").
		Order("fecha ASC, franja ASC").
		Find(&turnos).Error
	if err != nil {
		return nil, fmt.Errorf("listar turnos de test drive: %w", err)
	}
	return turnos, nil
}

// Ocupadas devuelve las franjas de una unidad y fecha que ya tienen un turno
// activo (solicitado o confirmado) y por lo tanto no se pueden elegir.
func (r *turnoTestDriveRepository) Ocupadas(ctx context.Context, vehiculoID uint, fecha string) ([]string, error) {
	var franjas []string
	err := r.base.WithContext(ctx).
		Model(&models.TurnoTestDrive{}).
		Where("vehiculo_id = ? AND fecha = ? AND estado IN ?",
			vehiculoID, fecha,
			[]string{models.EstadoTurnoSolicitado, models.EstadoTurnoConfirmado}).
		Pluck("franja", &franjas).Error
	if err != nil {
		return nil, fmt.Errorf("obtener franjas ocupadas: %w", err)
	}
	return franjas, nil
}

// Actualizar persiste los cambios de estado de un turno y lo devuelve con sus
// relaciones.
func (r *turnoTestDriveRepository) Actualizar(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error) {
	if err := r.base.WithContext(ctx).Save(turno).Error; err != nil {
		return nil, fmt.Errorf("actualizar turno de test drive: %w", err)
	}
	return r.ObtenerPorID(ctx, turno.ID)
}
