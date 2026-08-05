package repositories

import (
	"context"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// TurnoTestDriveRepository define el acceso a datos de turnos de test drive.
type TurnoTestDriveRepository interface {
	// Crear persiste un turno nuevo.
	Crear(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error)
	// ObtenerPorID devuelve un turno con su vehículo, imágenes y cliente.
	ObtenerPorID(ctx context.Context, id uint) (*models.TurnoTestDrive, error)
	// ListarPorCliente devuelve los turnos de un cliente, ordenados por fecha.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.TurnoTestDrive, error)
	// Listar devuelve los turnos con el filtro de estado opcional, ordenados
	// por fecha y franja.
	Listar(ctx context.Context, estado string) ([]models.TurnoTestDrive, error)
	// ExisteSuperposicion indica si hay un turno activo para la misma unidad,
	// fecha y franja.
	ExisteSuperposicion(ctx context.Context, vehiculoID uint, fecha string, franja string) (bool, error)
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

// Crear persiste un turno nuevo y lo devuelve con sus relaciones.
func (r *turnoTestDriveRepository) Crear(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error) {
	if err := r.base.WithContext(ctx).Create(turno).Error; err != nil {
		return nil, fmt.Errorf("crear turno de test drive: %w", err)
	}
	return r.ObtenerPorID(ctx, turno.ID)
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

// ExisteSuperposicion indica si hay un turno activo (solicitado o confirmado)
// para la misma unidad, fecha y franja.
func (r *turnoTestDriveRepository) ExisteSuperposicion(ctx context.Context, vehiculoID uint, fecha string, franja string) (bool, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.TurnoTestDrive{}).
		Where("vehiculo_id = ? AND fecha = ? AND franja = ? AND estado IN ?",
			vehiculoID, fecha, franja, []string{models.EstadoTurnoSolicitado, models.EstadoTurnoConfirmado}).
		Count(&total).Error
	if err != nil {
		return false, fmt.Errorf("verificar superposición de turnos: %w", err)
	}
	return total > 0, nil
}

// Actualizar persiste los cambios de estado de un turno y lo devuelve con sus
// relaciones.
func (r *turnoTestDriveRepository) Actualizar(ctx context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error) {
	if err := r.base.WithContext(ctx).Save(turno).Error; err != nil {
		return nil, fmt.Errorf("actualizar turno de test drive: %w", err)
	}
	return r.ObtenerPorID(ctx, turno.ID)
}
