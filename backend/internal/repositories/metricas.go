package repositories

import (
	"context"
	"time"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// ConteoPorEstado es el resultado de agrupar una entidad por su estado.
type ConteoPorEstado struct {
	Estado   string `json:"estado"`
	Cantidad int64  `json:"cantidad"`
}

// ConsultaPorDia es la cantidad de consultas creadas en un día.
type ConsultaPorDia struct {
	Fecha    string `json:"fecha"`
	Cantidad int64  `json:"cantidad"`
}

// MetricasRepository define las agregaciones de solo lectura que alimentan el
// panel de administración.
type MetricasRepository interface {
	// ContarVehiculosPorEstado devuelve la cantidad de vehículos agrupada por
	// estado.
	ContarVehiculosPorEstado(ctx context.Context) ([]ConteoPorEstado, error)
	// ContarConsultasPorDia devuelve la cantidad de consultas agrupada por día
	// (en la zona horaria del servidor) para las creadas desde la fecha
	// indicada. desplazamientoSegundos ajusta la fecha en la base a la zona
	// local del servidor.
	ContarConsultasPorDia(ctx context.Context, desde time.Time, desplazamientoSegundos int) ([]ConsultaPorDia, error)
	// ContarReservasActivas devuelve la cantidad de reservas en estado activa.
	ContarReservasActivas(ctx context.Context) (int64, error)
	// ContarReservasVendidas devuelve la cantidad de reservas en estado vendida.
	ContarReservasVendidas(ctx context.Context) (int64, error)
	// ContarTestDrivesAgendados devuelve la cantidad de turnos en estado
	// solicitado o confirmado.
	ContarTestDrivesAgendados(ctx context.Context) (int64, error)
	// ContarTestDrivesCompletados devuelve la cantidad de turnos en estado
	// completado.
	ContarTestDrivesCompletados(ctx context.Context) (int64, error)
	// ContarConsultasAbiertas devuelve la cantidad de consultas en estado
	// pendiente o en conversación.
	ContarConsultasAbiertas(ctx context.Context) (int64, error)
	// ContarUsuarios devuelve la cantidad total de cuentas de usuario.
	ContarUsuarios(ctx context.Context) (int64, error)
}

// metricasRepository implementa MetricasRepository sobre GORM.
type metricasRepository struct {
	base *gorm.DB
}

// NuevoMetricasRepository crea un repositorio de métricas.
func NuevoMetricasRepository(base *gorm.DB) MetricasRepository {
	return &metricasRepository{base: base}
}

// ContarVehiculosPorEstado agrupa los vehículos por estado.
func (r *metricasRepository) ContarVehiculosPorEstado(ctx context.Context) ([]ConteoPorEstado, error) {
	var resultado []ConteoPorEstado
	err := r.base.WithContext(ctx).
		Model(&models.Vehiculo{}).
		Select("estado, count(*) AS cantidad").
		Group("estado").
		Scan(&resultado).Error
	if err != nil {
		return nil, err
	}
	return resultado, nil
}

// ContarConsultasPorDia agrupa las consultas creadas desde la fecha indicada
// por día, aplicando el desplazamiento de la zona horaria del servidor para
// que el día coincida con la fecha local (la sesión de la base suele estar en
// UTC).
func (r *metricasRepository) ContarConsultasPorDia(ctx context.Context, desde time.Time, desplazamientoSegundos int) ([]ConsultaPorDia, error) {
	var resultado []ConsultaPorDia
	err := r.base.WithContext(ctx).
		Model(&models.Consulta{}).
		Select("date(created_at + make_interval(secs => ?))::text AS fecha, count(*) AS cantidad", desplazamientoSegundos).
		Where("created_at >= ?", desde).
		Group("fecha").
		Order("fecha").
		Scan(&resultado).Error
	if err != nil {
		return nil, err
	}
	return resultado, nil
}

// ContarReservasActivas cuenta las reservas en estado activa.
func (r *metricasRepository) ContarReservasActivas(ctx context.Context) (int64, error) {
	return r.contarReservas(ctx, models.EstadoReservaActiva)
}

// ContarReservasVendidas cuenta las reservas en estado vendida.
func (r *metricasRepository) ContarReservasVendidas(ctx context.Context) (int64, error) {
	return r.contarReservas(ctx, models.EstadoReservaVendida)
}

// contarReservas cuenta las reservas en un estado dado.
func (r *metricasRepository) contarReservas(ctx context.Context, estado string) (int64, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.Reserva{}).
		Where("estado = ?", estado).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// ContarTestDrivesAgendados cuenta los turnos en estado solicitado o confirmado.
func (r *metricasRepository) ContarTestDrivesAgendados(ctx context.Context) (int64, error) {
	return r.contarTestDrives(ctx, []string{models.EstadoTurnoSolicitado, models.EstadoTurnoConfirmado})
}

// ContarTestDrivesCompletados cuenta los turnos en estado completado.
func (r *metricasRepository) ContarTestDrivesCompletados(ctx context.Context) (int64, error) {
	return r.contarTestDrives(ctx, []string{models.EstadoTurnoCompletado})
}

// contarTestDrives cuenta los turnos en los estados indicados.
func (r *metricasRepository) contarTestDrives(ctx context.Context, estados []string) (int64, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.TurnoTestDrive{}).
		Where("estado IN ?", estados).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// ContarConsultasAbiertas cuenta las consultas en estado pendiente o en
// conversación.
func (r *metricasRepository) ContarConsultasAbiertas(ctx context.Context) (int64, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.Consulta{}).
		Where("estado IN ?", []string{models.EstadoPendiente, models.EstadoEnConversacion}).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// ContarUsuarios cuenta el total de cuentas de usuario.
func (r *metricasRepository) ContarUsuarios(ctx context.Context) (int64, error) {
	var total int64
	err := r.base.WithContext(ctx).
		Model(&models.Usuario{}).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
