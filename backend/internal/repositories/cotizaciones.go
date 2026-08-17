package repositories

import (
	"context"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// CotizacionRepository define el acceso a datos de cotizaciones sobre GORM.
type CotizacionRepository interface {
	// Crear persiste la cotización junto con sus mensajes iniciales.
	Crear(ctx context.Context, cotizacion *models.Cotizacion) (*models.Cotizacion, error)
	// AgregarMensaje guarda un mensaje nuevo en la cotización.
	AgregarMensaje(ctx context.Context, mensaje *models.MensajeCotizacion) error
	// ObtenerPorID devuelve la cotización con su vehículo, cliente y mensajes.
	ObtenerPorID(ctx context.Context, id uint) (*models.Cotizacion, error)
	// ListarPorCliente devuelve las cotizaciones de un cliente, ordenadas por
	// fecha de actualización descendente, con su vehículo y último mensaje.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Cotizacion, error)
	// Actualizar guarda los cambios de la cotización (por ejemplo, su estado).
	Actualizar(ctx context.Context, cotizacion *models.Cotizacion) error
}

// cotizacionRepository implementa CotizacionRepository sobre GORM.
type cotizacionRepository struct {
	base *gorm.DB
}

// NuevoCotizacionRepository crea un repositorio de cotizaciones.
func NuevoCotizacionRepository(base *gorm.DB) CotizacionRepository {
	return &cotizacionRepository{base: base}
}

// Crear persiste la cotización y sus mensajes iniciales en una transacción.
func (r *cotizacionRepository) Crear(ctx context.Context, cotizacion *models.Cotizacion) (*models.Cotizacion, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		mensajes := cotizacion.Mensajes
		cotizacion.Mensajes = nil

		if err := tx.Create(cotizacion).Error; err != nil {
			return err
		}

		for i := range mensajes {
			mensajes[i].CotizacionID = cotizacion.ID
		}
		if len(mensajes) > 0 {
			if err := tx.Create(&mensajes).Error; err != nil {
				return err
			}
		}

		cotizacion.Mensajes = mensajes
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("crear cotización: %w", err)
	}
	return r.ObtenerPorID(ctx, cotizacion.ID)
}

// AgregarMensaje guarda un mensaje nuevo en la cotización.
func (r *cotizacionRepository) AgregarMensaje(ctx context.Context, mensaje *models.MensajeCotizacion) error {
	if err := r.base.WithContext(ctx).Create(mensaje).Error; err != nil {
		return fmt.Errorf("agregar mensaje de cotización: %w", err)
	}
	return nil
}

// ObtenerPorID devuelve la cotización con su vehículo, cliente y mensajes.
func (r *cotizacionRepository) ObtenerPorID(ctx context.Context, id uint) (*models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := r.base.WithContext(ctx).
		Preload("Vehiculo").
		Preload("Cliente").
		Preload("Mensajes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&cotizacion, id).Error
	if err != nil {
		return nil, err
	}
	return &cotizacion, nil
}

// ListarPorCliente devuelve las cotizaciones de un cliente con su vehículo y
// el último mensaje, ordenadas por fecha de actualización descendente.
func (r *cotizacionRepository) ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Cotizacion, error) {
	var cotizaciones []models.Cotizacion
	err := r.base.WithContext(ctx).
		Where("cliente_id = ?", clienteID).
		Preload("Vehiculo").
		Preload("Mensajes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Order("updated_at DESC").
		Find(&cotizaciones).Error
	if err != nil {
		return nil, fmt.Errorf("listar cotizaciones por cliente: %w", err)
	}
	return cotizaciones, nil
}

// Actualizar guarda los cambios de la cotización.
func (r *cotizacionRepository) Actualizar(ctx context.Context, cotizacion *models.Cotizacion) error {
	if err := r.base.WithContext(ctx).
		Model(cotizacion).
		Updates(map[string]interface{}{"estado": cotizacion.Estado}).Error; err != nil {
		return fmt.Errorf("actualizar cotización: %w", err)
	}
	return nil
}
