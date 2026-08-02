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
	// ListarParaGestion devuelve los vehículos paginados para la administración.
	// Si estado está vacío incluye todos los estados; si no, filtra por el estado.
	ListarParaGestion(ctx context.Context, estado string, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	// ObtenerPorID devuelve un vehículo con su galería de imágenes.
	ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error)
	// Crear persiste un vehículo nuevo con sus imágenes.
	Crear(ctx context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error)
	// Actualizar actualiza la ficha y el estado de un vehículo, reemplazando su
	// galería de imágenes por la lista recibida.
	Actualizar(ctx context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error)
	// DarDeBaja cambia el estado del vehículo a dado_de_baja.
	DarDeBaja(ctx context.Context, id uint) error
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

// ListarParaGestion cuenta y devuelve la página solicitada de vehículos para la
// administración. Con estado vacío incluye todos los estados del stock.
func (r *vehiculoRepository) ListarParaGestion(ctx context.Context, estado string, pagina int, tamano int) ([]models.Vehiculo, int64, error) {
	consulta := r.base.WithContext(ctx).Model(&models.Vehiculo{})
	if estado != "" {
		consulta = consulta.Where("estado = ?", estado)
	}

	var total int64
	if err := consulta.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("contar vehículos para gestión: %w", err)
	}

	var vehiculos []models.Vehiculo
	err := consulta.
		Preload("Imagenes").
		Offset((pagina - 1) * tamano).
		Limit(tamano).
		Find(&vehiculos).Error
	if err != nil {
		return nil, 0, fmt.Errorf("listar vehículos para gestión: %w", err)
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

// Crear persiste el vehículo y su galería de imágenes, y devuelve el registro
// completo con los IDs asignados.
func (r *vehiculoRepository) Crear(ctx context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(vehiculo).Error; err != nil {
			return err
		}
		for i := range vehiculo.Imagenes {
			vehiculo.Imagenes[i].VehiculoID = vehiculo.ID
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("crear vehículo: %w", err)
	}
	return r.ObtenerPorID(ctx, vehiculo.ID)
}

// Actualizar reemplaza las imágenes existentes por la lista recibida, guarda la
// ficha técnica y el estado, y devuelve el registro actualizado.
func (r *vehiculoRepository) Actualizar(ctx context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("vehiculo_id = ?", vehiculo.ID).Delete(&models.Imagen{}).Error; err != nil {
			return err
		}

		imagenes := vehiculo.Imagenes
		vehiculo.Imagenes = nil
		if err := tx.Save(vehiculo).Error; err != nil {
			return err
		}

		for i := range imagenes {
			imagenes[i].ID = 0
			imagenes[i].VehiculoID = vehiculo.ID
		}
		if len(imagenes) > 0 {
			if err := tx.Create(&imagenes).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("actualizar vehículo: %w", err)
	}
	return r.ObtenerPorID(ctx, vehiculo.ID)
}

// DarDeBaja actualiza el estado del vehículo a dado_de_baja. Devuelve
// gorm.ErrRecordNotFound si el vehículo no existe.
func (r *vehiculoRepository) DarDeBaja(ctx context.Context, id uint) error {
	resultado := r.base.WithContext(ctx).
		Model(&models.Vehiculo{}).
		Where("id = ?", id).
		Update("estado", models.EstadoDadoDeBaja)
	if resultado.Error != nil {
		return fmt.Errorf("dar de baja vehículo: %w", resultado.Error)
	}
	if resultado.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
