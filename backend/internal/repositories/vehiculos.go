package repositories

import (
	"context"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// VehiculoRepository define el acceso a datos de vehículos sobre GORM.
type VehiculoRepository interface {
	// Listar devuelve los vehículos que cumplen el estado indicado, paginados,
	// junto con el total de registros que cumplen la condición.
	Listar(ctx context.Context, estado string, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	// ObtenerPorID devuelve un vehículo con su galería de imágenes.
	ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error)
}

// vehiculoRepository implementa VehiculoRepository sobre GORM.
type vehiculoRepository struct {
	base *gorm.DB
}

// NuevoVehiculoRepository crea un repositorio de vehículos.
func NuevoVehiculoRepository(base *gorm.DB) VehiculoRepository {
	return &vehiculoRepository{base: base}
}

// Listar cuenta y devuelve la página solicitada de vehículos.
func (r *vehiculoRepository) Listar(ctx context.Context, estado string, pagina int, tamano int) ([]models.Vehiculo, int64, error) {
	var total int64
	if err := r.base.WithContext(ctx).
		Model(&models.Vehiculo{}).
		Where("estado = ?", estado).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("contar vehículos: %w", err)
	}

	var vehiculos []models.Vehiculo
	err := r.base.WithContext(ctx).
		Model(&models.Vehiculo{}).
		Where("estado = ?", estado).
		Preload("Imagenes").
		Offset((pagina - 1) * tamano).
		Limit(tamano).
		Find(&vehiculos).Error
	if err != nil {
		return nil, 0, fmt.Errorf("listar vehículos: %w", err)
	}

	return vehiculos, total, nil
}

// ObtenerPorID devuelve un vehículo con sus imágenes o un error de GORM
// (ErrRecordNotFound si no existe).
func (r *vehiculoRepository) ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error) {
	var vehiculo models.Vehiculo
	err := r.base.WithContext(ctx).
		Preload("Imagenes").
		First(&vehiculo, id).Error
	if err != nil {
		return nil, err
	}
	return &vehiculo, nil
}
