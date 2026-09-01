package repositories

import (
	"context"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// MensajeRepository define el acceso a datos de mensajes sobre GORM.
type MensajeRepository interface {
	// Crear persiste un mensaje nuevo.
	Crear(ctx context.Context, mensaje *models.Mensaje) (*models.Mensaje, error)
	// ListarPorConsulta devuelve todos los mensajes de una consulta ordenados
	// por fecha ascendente.
	ListarPorConsulta(ctx context.Context, consultaID uint) ([]models.Mensaje, error)
	// ObtenerDesdeID devuelve los mensajes de una consulta con id mayor al
	// indicado, ordenados cronológicamente. Es el mecanismo para traer en el
	// polling solo lo nuevo (ver docs/roadmap.md "Escalabilidad de conversaciones").
	ObtenerDesdeID(ctx context.Context, consultaID uint, desdeID uint) ([]models.Mensaje, error)
	// MarcarComoLeidos marca como leídos los mensajes de una consulta cuyo
	// remitente sea el indicado.
	MarcarComoLeidos(ctx context.Context, consultaID uint, remitenteID uint) error
	// ContarNoLeidosPorConsultas devuelve un mapa con la cantidad de mensajes
	// no leídos por consulta, para los remitentes indicados.
	ContarNoLeidosPorConsultas(ctx context.Context, consultaIDs []uint, remitenteID uint) (map[uint]int, error)
}

// mensajeRepository implementa MensajeRepository sobre GORM.
type mensajeRepository struct {
	base *gorm.DB
}

// NuevoMensajeRepository crea un repositorio de mensajes.
func NuevoMensajeRepository(base *gorm.DB) MensajeRepository {
	return &mensajeRepository{base: base}
}

// Crear persiste el mensaje y devuelve el registro completo.
func (r *mensajeRepository) Crear(ctx context.Context, mensaje *models.Mensaje) (*models.Mensaje, error) {
	err := r.base.WithContext(ctx).Create(mensaje).Error
	if err != nil {
		return nil, fmt.Errorf("crear mensaje: %w", err)
	}
	return r.obtenerPorID(ctx, mensaje.ID)
}

// obtenerPorID devuelve un mensaje con su remitente.
func (r *mensajeRepository) obtenerPorID(ctx context.Context, id uint) (*models.Mensaje, error) {
	var mensaje models.Mensaje
	err := r.base.WithContext(ctx).
		Preload("Remitente").
		First(&mensaje, id).Error
	if err != nil {
		return nil, err
	}
	return &mensaje, nil
}

// ListarPorConsulta devuelve todos los mensajes de una consulta ordenados
// cronológicamente.
func (r *mensajeRepository) ListarPorConsulta(ctx context.Context, consultaID uint) ([]models.Mensaje, error) {
	var mensajes []models.Mensaje
	err := r.base.WithContext(ctx).
		Where("consulta_id = ?", consultaID).
		Preload("Remitente").
		Order("created_at ASC").
		Find(&mensajes).Error
	if err != nil {
		return nil, fmt.Errorf("listar mensajes por consulta: %w", err)
	}
	return mensajes, nil
}

// ObtenerDesdeID devuelve los mensajes de una consulta con id mayor al
// indicado. Se ordena por id porque es monotónico y el índice compuesto de la
// tabla ya cubre la consulta por hilo.
func (r *mensajeRepository) ObtenerDesdeID(ctx context.Context, consultaID uint, desdeID uint) ([]models.Mensaje, error) {
	var mensajes []models.Mensaje
	err := r.base.WithContext(ctx).
		Where("consulta_id = ? AND id > ?", consultaID, desdeID).
		Preload("Remitente").
		Order("id ASC").
		Find(&mensajes).Error
	if err != nil {
		return nil, fmt.Errorf("obtener mensajes nuevos: %w", err)
	}
	return mensajes, nil
}

// MarcarComoLeidos marca como leídos los mensajes de otro remitente en una
// consulta.
func (r *mensajeRepository) MarcarComoLeidos(ctx context.Context, consultaID uint, remitenteID uint) error {
	err := r.base.WithContext(ctx).
		Model(&models.Mensaje{}).
		Where("consulta_id = ? AND remitente_id != ? AND leido = ?", consultaID, remitenteID, false).
		Update("leido", true).Error
	if err != nil {
		return fmt.Errorf("marcar mensajes como leídos: %w", err)
	}
	return nil
}

// ContarNoLeidosPorConsultas cuenta los mensajes no leídos de un remitente
// para cada consulta del listado.
func (r *mensajeRepository) ContarNoLeidosPorConsultas(ctx context.Context, consultaIDs []uint, remitenteID uint) (map[uint]int, error) {
	if len(consultaIDs) == 0 {
		return make(map[uint]int), nil
	}

	type resultado struct {
		ConsultaID uint
		Cantidad   int
	}

	var resultados []resultado
	err := r.base.WithContext(ctx).
		Model(&models.Mensaje{}).
		Select("consulta_id, COUNT(*) as cantidad").
		Where("consulta_id IN ? AND remitente_id != ? AND leido = ?", consultaIDs, remitenteID, false).
		Group("consulta_id").
		Scan(&resultados).Error
	if err != nil {
		return nil, fmt.Errorf("contar mensajes no leídos: %w", err)
	}

	mapa := make(map[uint]int, len(resultados))
	for _, r := range resultados {
		mapa[r.ConsultaID] = r.Cantidad
	}
	return mapa, nil
}
