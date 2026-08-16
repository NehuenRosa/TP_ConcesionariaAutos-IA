package repositories

import (
	"context"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// ConsultaRepository define el acceso a datos de consultas sobre GORM.
type ConsultaRepository interface {
	// Crear persiste una consulta nueva.
	Crear(ctx context.Context, consulta *models.Consulta) (*models.Consulta, error)
	// ObtenerPorID devuelve una consulta con todas sus relaciones.
	ObtenerPorID(ctx context.Context, id uint) (*models.Consulta, error)
	// ListarPorCliente devuelve las consultas de un cliente, ordenadas por
	// fecha de actualización descendente.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Consulta, error)
	// ListarPorVendedor devuelve las consultas asignadas a un vendedor,
	// ordenadas por fecha de actualización descendente.
	ListarPorVendedor(ctx context.Context, vendedorID uint) ([]models.Consulta, error)
	// ListarPorUsuario devuelve las consultas donde el usuario participa
	// como cliente o vendedor.
	ListarPorUsuario(ctx context.Context, usuarioID uint) ([]uint, error)
	// TomarSiPendiente asigna el vendedor a una consulta solo si sigue
	// pendiente. Devuelve false si otro vendedor la tomó antes (UPDATE atómico
	// para evitar carreras entre dos vendedores).
	TomarSiPendiente(ctx context.Context, consultaID uint, vendedorID uint) (bool, error)
	// Actualizar actualiza el estado y el vendedor de una consulta.
	Actualizar(ctx context.Context, consulta *models.Consulta) (*models.Consulta, error)
	// Eliminar elimina una consulta y sus mensajes.
	Eliminar(ctx context.Context, id uint) error
}

// consultaRepository implementa ConsultaRepository sobre GORM.
type consultaRepository struct {
	base *gorm.DB
}

// NuevoConsultaRepository crea un repositorio de consultas.
func NuevoConsultaRepository(base *gorm.DB) ConsultaRepository {
	return &consultaRepository{base: base}
}

// Crear persiste la consulta con su primer mensaje.
func (r *consultaRepository) Crear(ctx context.Context, consulta *models.Consulta) (*models.Consulta, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Guardar el mensaje antes de crear la consulta para tener el ConsultaID
		mensaje := consulta.Mensajes[0]
		consulta.Mensajes = nil

		if err := tx.Create(consulta).Error; err != nil {
			return err
		}

		mensaje.ConsultaID = consulta.ID
		if err := tx.Create(&mensaje).Error; err != nil {
			return err
		}

		consulta.Mensajes = []models.Mensaje{mensaje}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("crear consulta: %w", err)
	}
	return r.ObtenerPorID(ctx, consulta.ID)
}

// ObtenerPorID devuelve una consulta con sus relaciones.
func (r *consultaRepository) ObtenerPorID(ctx context.Context, id uint) (*models.Consulta, error) {
	var consulta models.Consulta
	err := r.base.WithContext(ctx).
		Preload("Vehiculo").
		Preload("Cliente").
		Preload("Vendedor").
		Preload("Mensajes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Mensajes.Remitente").
		First(&consulta, id).Error
	if err != nil {
		return nil, err
	}
	return &consulta, nil
}

// ListarPorCliente devuelve las consultas de un cliente con información
// resumida del vehículo, vendedor y último mensaje.
func (r *consultaRepository) ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Consulta, error) {
	var consultas []models.Consulta
	err := r.base.WithContext(ctx).
		Where("cliente_id = ?", clienteID).
		Preload("Vehiculo").
		Preload("Vendedor").
		Preload("Mensajes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Preload("Mensajes.Remitente").
		Order("updated_at DESC").
		Find(&consultas).Error
	if err != nil {
		return nil, fmt.Errorf("listar consultas por cliente: %w", err)
	}
	return consultas, nil
}

// ListarPorVendedor devuelve las consultas asignadas a un vendedor Y las
// consultas pendientes (sin vendedor), con información resumida del vehículo,
// cliente y último mensaje.
func (r *consultaRepository) ListarPorVendedor(ctx context.Context, vendedorID uint) ([]models.Consulta, error) {
	var consultas []models.Consulta
	err := r.base.WithContext(ctx).
		Where("vendedor_id = ? OR (vendedor_id IS NULL AND estado = ?)", vendedorID, models.EstadoPendiente).
		Preload("Vehiculo").
		Preload("Cliente").
		Preload("Mensajes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Preload("Mensajes.Remitente").
		Order("updated_at DESC").
		Find(&consultas).Error
	if err != nil {
		return nil, fmt.Errorf("listar consultas por vendedor: %w", err)
	}
	return consultas, nil
}

// ListarPendientes devuelve las consultas sin vendedor asignado.
func (r *consultaRepository) ListarPendientes(ctx context.Context) ([]models.Consulta, error) {
	var consultas []models.Consulta
	err := r.base.WithContext(ctx).
		Where("vendedor_id IS NULL AND estado = ?", models.EstadoPendiente).
		Preload("Vehiculo").
		Preload("Cliente").
		Preload("Mensajes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Preload("Mensajes.Remitente").
		Order("created_at DESC").
		Find(&consultas).Error
	if err != nil {
		return nil, fmt.Errorf("listar consultas pendientes: %w", err)
	}
	return consultas, nil
}

// TomarSiPendiente asigna el vendedor a la consulta solo si sigue pendiente.
func (r *consultaRepository) TomarSiPendiente(ctx context.Context, consultaID uint, vendedorID uint) (bool, error) {
	resultado := r.base.WithContext(ctx).
		Model(&models.Consulta{}).
		Where("id = ? AND estado = ? AND vendedor_id IS NULL", consultaID, models.EstadoPendiente).
		Updates(map[string]any{
			"vendedor_id": vendedorID,
			"estado":      models.EstadoEnConversacion,
		})
	if resultado.Error != nil {
		return false, fmt.Errorf("tomar consulta: %w", resultado.Error)
	}
	return resultado.RowsAffected > 0, nil
}

// Actualizar actualiza el estado y opcionalmente el vendedor de una consulta.
func (r *consultaRepository) Actualizar(ctx context.Context, consulta *models.Consulta) (*models.Consulta, error) {
	err := r.base.WithContext(ctx).Save(consulta).Error
	if err != nil {
		return nil, fmt.Errorf("actualizar consulta: %w", err)
	}
	return r.ObtenerPorID(ctx, consulta.ID)
}

// Eliminar elimina una consulta y sus mensajes en una transacción.
func (r *consultaRepository) Eliminar(ctx context.Context, id uint) error {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("consulta_id = ?", id).Delete(&models.Mensaje{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Consulta{}, id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("eliminar consulta: %w", err)
	}
	return nil
}

// ListarPorUsuario devuelve los IDs de consultas donde el usuario participa.
func (r *consultaRepository) ListarPorUsuario(ctx context.Context, usuarioID uint) ([]uint, error) {
	var ids []uint
	err := r.base.WithContext(ctx).
		Model(&models.Consulta{}).
		Where("cliente_id = ? OR vendedor_id = ?", usuarioID, usuarioID).
		Pluck("id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("listar consultas por usuario: %w", err)
	}
	return ids, nil
}
