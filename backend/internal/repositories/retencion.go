package repositories

import (
	"context"
	"fmt"
	"time"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// RetencionRepository define el acceso a datos para purgar conversaciones
// cerradas viejas (política de retención). Las conversaciones abiertas nunca
// se tocan: solo las cerradas que superan el plazo de conservación.
type RetencionRepository interface {
	// PurgarCotizacionesCerradas elimina los mensajes y luego las cotizaciones
	// cerradas cuya última actualización sea anterior al corte. Devuelve la
	// cantidad de mensajes y de cotizaciones eliminadas.
	PurgarCotizacionesCerradas(ctx context.Context, corte time.Time) (int64, int64, error)
	// PurgarConsultasCerradas elimina los mensajes y luego las consultas
	// cerradas cuya última actualización sea anterior al corte. Devuelve la
	// cantidad de mensajes y de consultas eliminadas.
	PurgarConsultasCerradas(ctx context.Context, corte time.Time) (int64, int64, error)
}

// retencionRepository implementa RetencionRepository sobre GORM.
type retencionRepository struct {
	base *gorm.DB
}

// NuevoRetencionRepository crea un repositorio de retención.
func NuevoRetencionRepository(base *gorm.DB) RetencionRepository {
	return &retencionRepository{base: base}
}

// PurgarCotizacionesCerradas elimina en una transacción los mensajes de las
// cotizaciones cerradas viejas y después las cotizaciones mismas (los mensajes
// son los que más filas aportan, por eso se borran primero).
func (r *retencionRepository) PurgarCotizacionesCerradas(ctx context.Context, corte time.Time) (int64, int64, error) {
	var mensajesEliminados int64
	var cotizacionesEliminadas int64

	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		subconsulta := tx.Model(&models.Cotizacion{}).
			Where("estado = ? AND updated_at < ?", models.EstadoCotizacionCerrada, corte).
			Select("id")

		resultadoMensajes := tx.Model(&models.MensajeCotizacion{}).
			Where("cotizacion_id IN (?)", subconsulta).
			Delete(&models.MensajeCotizacion{})
		if err := resultadoMensajes.Error; err != nil {
			return err
		}
		mensajesEliminados = resultadoMensajes.RowsAffected

		resultado := tx.
			Where("estado = ? AND updated_at < ?", models.EstadoCotizacionCerrada, corte).
			Delete(&models.Cotizacion{})
		if err := resultado.Error; err != nil {
			return err
		}
		cotizacionesEliminadas = resultado.RowsAffected
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("purgar cotizaciones cerradas: %w", err)
	}
	return mensajesEliminados, cotizacionesEliminadas, nil
}

// PurgarConsultasCerradas elimina en una transacción los mensajes de las
// consultas cerradas viejas y después las consultas mismas.
func (r *retencionRepository) PurgarConsultasCerradas(ctx context.Context, corte time.Time) (int64, int64, error) {
	var mensajesEliminados int64
	var consultasEliminadas int64

	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		subconsulta := tx.Model(&models.Consulta{}).
			Where("estado = ? AND updated_at < ?", models.EstadoCerrada, corte).
			Select("id")

		resultadoMensajes := tx.Model(&models.Mensaje{}).
			Where("consulta_id IN (?)", subconsulta).
			Delete(&models.Mensaje{})
		if err := resultadoMensajes.Error; err != nil {
			return err
		}
		mensajesEliminados = resultadoMensajes.RowsAffected

		resultado := tx.
			Where("estado = ? AND updated_at < ?", models.EstadoCerrada, corte).
			Delete(&models.Consulta{})
		if err := resultado.Error; err != nil {
			return err
		}
		consultasEliminadas = resultado.RowsAffected
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("purgar consultas cerradas: %w", err)
	}
	return mensajesEliminados, consultasEliminadas, nil
}