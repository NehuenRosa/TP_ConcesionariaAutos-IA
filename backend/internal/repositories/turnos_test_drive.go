package repositories

import (
	"context"
	"errors"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// errSuperposicionExistente es el centinela interno que indica que la regla de
// no superposición se disparó dentro de la transacción.
var errSuperposicionExistente = errors.New("ya existe un turno superpuesto")

// TurnoTestDriveRepository define el acceso a datos de turnos de test drive.
type TurnoTestDriveRepository interface {
	// CrearSiSinSuperposicion crea el turno en una transacción que bloquea la
	// fila del vehículo, garantizando la regla de no superposición ante
	// pedidos concurrentes. Devuelve el turno creado y true, o nil y false si
	// ya había un turno activo para la unidad, fecha y franja.
	CrearSiSinSuperposicion(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, bool, error)
	// ObtenerPorID devuelve un turno con su vehículo, imágenes y cliente.
	ObtenerPorID(ctx context.Context, id uint) (*models.TurnoTestDrive, error)
	// ListarPorCliente devuelve los turnos de un cliente, ordenados por fecha.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.TurnoTestDrive, error)
	// Listar devuelve los turnos con el filtro de estado opcional, ordenados
	// por fecha y franja.
	Listar(ctx context.Context, estado string) ([]models.TurnoTestDrive, error)
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

// CrearSiSinSuperposicion crea el turno en una transacción con bloqueo de la
// fila del vehículo (FOR UPDATE): los pedidos concurrentes para la misma
// unidad, fecha y franja se serializan y la regla de no superposición no se
// puede violar.
func (r *turnoTestDriveRepository) CrearSiSinSuperposicion(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, bool, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var vehiculo models.Vehiculo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&vehiculo, turno.VehiculoID).Error; err != nil {
			return err
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
			return errSuperposicionExistente
		}

		return tx.Create(turno).Error
	})
	if err != nil {
		if errors.Is(err, errSuperposicionExistente) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("crear turno de test drive: %w", err)
	}

	creado, err := r.ObtenerPorID(ctx, turno.ID)
	if err != nil {
		return nil, false, fmt.Errorf("recuperar turno de test drive: %w", err)
	}
	return creado, true, nil
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

// Actualizar persiste los cambios de estado de un turno y lo devuelve con sus
// relaciones.
func (r *turnoTestDriveRepository) Actualizar(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error) {
	if err := r.base.WithContext(ctx).Save(turno).Error; err != nil {
		return nil, fmt.Errorf("actualizar turno de test drive: %w", err)
	}
	return r.ObtenerPorID(ctx, turno.ID)
}
