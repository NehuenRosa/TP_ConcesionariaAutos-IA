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
	// ObtenerCabecera devuelve la cotización con vehículo, cliente y vendedor,
	// sin cargar los mensajes. Sirve para validar la propiedad de un hilo y
	// conocer su estado, vendedor y fecha de toma sin traer el historial
	// completo.
	ObtenerCabecera(ctx context.Context, id uint) (*models.Cotizacion, error)
	// ObtenerMensajesDesde devuelve los mensajes de una cotización con id mayor
	// al indicado, ordenados cronológicamente. Es el mecanismo para traer en el
	// polling solo lo nuevo (ver docs/roadmap.md "Escalabilidad de conversaciones").
	ObtenerMensajesDesde(ctx context.Context, cotizacionID uint, desdeID uint) ([]models.MensajeCotizacion, error)
	// ContarMensajes cuenta el total de mensajes de una cotización.
	ContarMensajes(ctx context.Context, cotizacionID uint) (int64, error)
	// ListarPorCliente devuelve las cotizaciones de un cliente, ordenadas por
	// fecha de actualización descendente, con su vehículo y último mensaje.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Cotizacion, error)
	// ListarBandeja devuelve todas las cotizaciones para el personal, ordenadas
	// por fecha de actualización descendente, con vehículo, cliente, vendedor y
	// último mensaje.
	ListarBandeja(ctx context.Context) ([]models.Cotizacion, error)
	// Actualizar guarda los cambios de la cotización (por ejemplo, su estado).
	Actualizar(ctx context.Context, cotizacion *models.Cotizacion) error
	// ContarNoLeidosDeCliente cuenta los mensajes de la IA y del vendedor sin
	// leer en las cotizaciones del cliente indicado.
	ContarNoLeidosDeCliente(ctx context.Context, clienteID uint) (int64, error)
	// ContarNoLeidosParaPersonal cuenta los mensajes de cliente sin leer en
	// cotizaciones abiertas sin asignar o asignadas al vendedor indicado.
	ContarNoLeidosParaPersonal(ctx context.Context, vendedorID uint) (int64, error)
	// MarcarLeidasParaCliente marca como leídos (lado cliente) los mensajes de
	// ia/vendedor de una cotización.
	MarcarLeidasParaCliente(ctx context.Context, cotizacionID uint) error
	// MarcarLeidasParaPersonal marca como leídos (lado personal) los mensajes
	// de cliente de una cotización.
	MarcarLeidasParaPersonal(ctx context.Context, cotizacionID uint) error
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

// ObtenerPorID devuelve la cotización con su vehículo, cliente, vendedor y
// mensajes.
func (r *cotizacionRepository) ObtenerPorID(ctx context.Context, id uint) (*models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := r.base.WithContext(ctx).
		Preload("Vehiculo").
		Preload("Cliente").
		Preload("Vendedor").
		Preload("Mensajes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&cotizacion, id).Error
	if err != nil {
		return nil, err
	}
	return &cotizacion, nil
}

// ObtenerCabecera devuelve la cotización con vehículo, cliente y vendedor,
// sin cargar los mensajes. Sirve para validar la propiedad de un hilo y
// conocer su estado, vendedor y fecha de toma sin traer el historial completo.
func (r *cotizacionRepository) ObtenerCabecera(ctx context.Context, id uint) (*models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := r.base.WithContext(ctx).
		Preload("Vehiculo").
		Preload("Cliente").
		Preload("Vendedor").
		First(&cotizacion, id).Error
	if err != nil {
		return nil, err
	}
	return &cotizacion, nil
}

// ObtenerMensajesDesde devuelve los mensajes de una cotización con id mayor al
// indicado, ordenados cronológicamente.
func (r *cotizacionRepository) ObtenerMensajesDesde(ctx context.Context, cotizacionID uint, desdeID uint) ([]models.MensajeCotizacion, error) {
	var mensajes []models.MensajeCotizacion
	err := r.base.WithContext(ctx).
		Where("cotizacion_id = ? AND id > ?", cotizacionID, desdeID).
		Order("id ASC").
		Find(&mensajes).Error
	if err != nil {
		return nil, fmt.Errorf("obtener mensajes nuevos de cotización: %w", err)
	}
	return mensajes, nil
}

// ContarMensajes cuenta el total de mensajes de una cotización.
func (r *cotizacionRepository) ContarMensajes(ctx context.Context, cotizacionID uint) (int64, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.MensajeCotizacion{}).
		Where("cotizacion_id = ?", cotizacionID).
		Count(&total).Error
	if err != nil {
		return 0, fmt.Errorf("contar mensajes de cotización: %w", err)
	}
	return total, nil
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

// ListarBandeja devuelve todas las cotizaciones para el personal, ordenadas
// por fecha de actualización descendente, con vehículo, cliente, vendedor y
// último mensaje.
func (r *cotizacionRepository) ListarBandeja(ctx context.Context) ([]models.Cotizacion, error) {
	var cotizaciones []models.Cotizacion
	err := r.base.WithContext(ctx).
		Preload("Vehiculo").
		Preload("Cliente").
		Preload("Vendedor").
		Preload("Mensajes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Order("updated_at DESC").
		Find(&cotizaciones).Error
	if err != nil {
		return nil, fmt.Errorf("listar bandeja de cotizaciones: %w", err)
	}
	return cotizaciones, nil
}

// Actualizar guarda los cambios de la cotización: su estado y, si corresponde,
// el vendedor que la tomó con la fecha de toma.
func (r *cotizacionRepository) Actualizar(ctx context.Context, cotizacion *models.Cotizacion) error {
	if err := r.base.WithContext(ctx).
		Model(cotizacion).
		Updates(map[string]interface{}{
			"estado":      cotizacion.Estado,
			"vendedor_id": cotizacion.VendedorID,
			"fecha_toma":  cotizacion.FechaToma,
		}).Error; err != nil {
		return fmt.Errorf("actualizar cotización: %w", err)
	}
	return nil
}

// ContarNoLeidosDeCliente cuenta los mensajes de la IA y del vendedor sin
// leer en las cotizaciones del cliente: una respuesta del asistente que aún no
// se abrió (o que llegó mientras el cliente estaba en otra pestaña) también
// debe avisarse, igual que una respuesta de un vendedor.
func (r *cotizacionRepository) ContarNoLeidosDeCliente(ctx context.Context, clienteID uint) (int64, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.MensajeCotizacion{}).
		Joins("JOIN cotizaciones ON cotizaciones.id = cotizacion_mensajes.cotizacion_id").
		Where(
			"cotizaciones.cliente_id = ? AND cotizacion_mensajes.remitente <> ? AND cotizacion_mensajes.leido_por_cliente = ?",
			clienteID, models.RemitenteCliente, false,
		).
		Count(&total).Error
	if err != nil {
		return 0, fmt.Errorf("contar mensajes no leídos del cliente: %w", err)
	}
	return total, nil
}

// ContarNoLeidosParaPersonal cuenta los mensajes de cliente sin leer en
// cotizaciones abiertas sin asignar o asignadas al vendedor indicado. Las
// cerradas y las atendidas por otro vendedor no cuentan.
func (r *cotizacionRepository) ContarNoLeidosParaPersonal(ctx context.Context, vendedorID uint) (int64, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.MensajeCotizacion{}).
		Joins("JOIN cotizaciones ON cotizaciones.id = cotizacion_mensajes.cotizacion_id").
		Where(
			"cotizaciones.estado = ? AND (cotizaciones.vendedor_id IS NULL OR cotizaciones.vendedor_id = ?)"+
				" AND cotizacion_mensajes.remitente = ? AND cotizacion_mensajes.leido_por_vendedor = ?",
			models.EstadoCotizacionAbierta, vendedorID, models.RemitenteCliente, false,
		).
		Count(&total).Error
	if err != nil {
		return 0, fmt.Errorf("contar mensajes no leídos para el personal: %w", err)
	}
	return total, nil
}

// MarcarLeidasParaCliente marca como leídos (lado cliente) los mensajes de
// ia/vendedor de la cotización indicada.
func (r *cotizacionRepository) MarcarLeidasParaCliente(ctx context.Context, cotizacionID uint) error {
	err := r.base.WithContext(ctx).
		Model(&models.MensajeCotizacion{}).
		Where("cotizacion_id = ? AND remitente <> ? AND leido_por_cliente = ?", cotizacionID, models.RemitenteCliente, false).
		Update("leido_por_cliente", true).Error
	if err != nil {
		return fmt.Errorf("marcar leídos para el cliente: %w", err)
	}
	return nil
}

// MarcarLeidasParaPersonal marca como leídos (lado personal) los mensajes de
// cliente de la cotización indicada.
func (r *cotizacionRepository) MarcarLeidasParaPersonal(ctx context.Context, cotizacionID uint) error {
	err := r.base.WithContext(ctx).
		Model(&models.MensajeCotizacion{}).
		Where("cotizacion_id = ? AND remitente = ? AND leido_por_vendedor = ?", cotizacionID, models.RemitenteCliente, false).
		Update("leido_por_vendedor", true).Error
	if err != nil {
		return fmt.Errorf("marcar leídos para el personal: %w", err)
	}
	return nil
}
