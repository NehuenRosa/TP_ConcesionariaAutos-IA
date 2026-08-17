package repositories

import (
	"context"
	"errors"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// ErrTasacionNoEncontrada indica que no existe una tasación pendiente para la
// sesión dada.
var ErrTasacionNoEncontrada = errors.New("no hay una tasación pendiente para esa sesión")

// ErrCodigoRepetido indica que el código generado ya existe en la base.
var ErrCodigoRepetido = errors.New("el código de tasación ya fue usado")

// TasacionRepository define el acceso a datos de tasaciones.
type TasacionRepository interface {
	// Crear guarda una nueva tasación.
	Crear(ctx context.Context, tasacion *models.Tasacion) error
	// Actualizar persiste los cambios de una tasación.
	Actualizar(ctx context.Context, tasacion *models.Tasacion) error
	// ObtenerPendientePorSesion devuelve la tasación pendiente de una sesión.
	ObtenerPendientePorSesion(ctx context.Context, sesionID string) (*models.Tasacion, error)
	// CodigoExiste indica si un código de presentación ya fue asignado.
	CodigoExiste(ctx context.Context, codigo string) (bool, error)
}

// tasacionRepository implementa TasacionRepository sobre GORM.
type tasacionRepository struct {
	base *gorm.DB
}

// NuevoTasacionRepository crea un repositorio de tasaciones.
func NuevoTasacionRepository(base *gorm.DB) TasacionRepository {
	return &tasacionRepository{base: base}
}

// Crear guarda una nueva tasación.
func (r *tasacionRepository) Crear(ctx context.Context, tasacion *models.Tasacion) error {
	if err := r.base.WithContext(ctx).Create(tasacion).Error; err != nil {
		return fmt.Errorf("crear tasación: %w", err)
	}
	return nil
}

// Actualizar persiste los cambios de una tasación.
func (r *tasacionRepository) Actualizar(ctx context.Context, tasacion *models.Tasacion) error {
	if err := r.base.WithContext(ctx).Save(tasacion).Error; err != nil {
		return fmt.Errorf("actualizar tasación: %w", err)
	}
	return nil
}

// ObtenerPendientePorSesion devuelve la tasación pendiente de una sesión.
func (r *tasacionRepository) ObtenerPendientePorSesion(ctx context.Context, sesionID string) (*models.Tasacion, error) {
	var tasacion models.Tasacion
	err := r.base.WithContext(ctx).
		Where("sesion_id = ? AND estado_flujo = ?", sesionID, models.EstadoTasacionPendiente).
		First(&tasacion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTasacionNoEncontrada
		}
		return nil, fmt.Errorf("buscar tasación pendiente: %w", err)
	}
	return &tasacion, nil
}

// CodigoExiste indica si un código de presentación ya fue asignado.
func (r *tasacionRepository) CodigoExiste(ctx context.Context, codigo string) (bool, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.Tasacion{}).
		Where("codigo = ?", codigo).
		Count(&total).Error
	if err != nil {
		return false, fmt.Errorf("consultar código de tasación: %w", err)
	}
	return total > 0, nil
}