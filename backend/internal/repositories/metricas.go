package repositories

import (
	"context"
	"fmt"
	"time"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// ConteoPorEstado es el resultado de agrupar una entidad por su estado.
type ConteoPorEstado struct {
	Estado   string `json:"estado"`
	Cantidad int64  `json:"cantidad"`
}

// CantidadPorDia es la cantidad de eventos (consultas, ventas o turnos) que
// ocurrieron en un día.
type CantidadPorDia struct {
	Fecha    string `json:"fecha"`
	Cantidad int64  `json:"cantidad"`
}

// ConteoPorMarca es la cantidad de vehículos vendidos agrupada por marca.
type ConteoPorMarca struct {
	Marca    string `json:"marca"`
	Cantidad int64  `json:"cantidad"`
}

// VehiculoConAntiguedad es un vehículo junto con su fecha de alta, para
// calcular cuántos días lleva publicado.
type VehiculoConAntiguedad struct {
	ID        uint      `json:"id"`
	Marca     string    `json:"marca"`
	Modelo    string    `json:"modelo"`
	Anio      int       `json:"anio"`
	CreatedAt time.Time `json:"-"`
}

// VehiculoEnStock es un vehículo publicado con la cantidad de días que lleva
// en el lote.
type VehiculoEnStock struct {
	ID          uint   `json:"id"`
	Marca       string `json:"marca"`
	Modelo      string `json:"modelo"`
	Anio        int    `json:"anio"`
	DiasEnStock int    `json:"diasEnStock"`
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
	ContarConsultasPorDia(ctx context.Context, desde time.Time, desplazamientoSegundos int) ([]CantidadPorDia, error)
	// ContarVentasPorDia devuelve la cantidad de reservas vendidas agrupada por
	// día para las confirmadas desde la fecha indicada.
	ContarVentasPorDia(ctx context.Context, desde time.Time, desplazamientoSegundos int) ([]CantidadPorDia, error)
	// SumarIngresosPorPeriodo devuelve la suma del precio de los vehículos
	// cuyas ventas fueron confirmadas desde la fecha indicada.
	SumarIngresosPorPeriodo(ctx context.Context, desde time.Time) (float64, error)
	// ContarVentasPorMarca devuelve la cantidad de vehículos vendidos agrupada
	// por marca.
	ContarVentasPorMarca(ctx context.Context) ([]ConteoPorMarca, error)
	// ContarTestDrivesPorDia devuelve la cantidad de turnos agendados
	// (solicitados o confirmados) agrupada por día desde la fecha indicada.
	ContarTestDrivesPorDia(ctx context.Context, desde time.Time, desplazamientoSegundos int) ([]CantidadPorDia, error)
	// TraerVehiculosConAntiguedad devuelve los primeros limite vehículos
	// disponibles con mayor tiempo publicado (por orden de alta).
	TraerVehiculosConAntiguedad(ctx context.Context, limite int) ([]VehiculoConAntiguedad, error)
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
func (r *metricasRepository) ContarConsultasPorDia(ctx context.Context, desde time.Time, desplazamientoSegundos int) ([]CantidadPorDia, error) {
	var resultado []CantidadPorDia
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

// ContarVentasPorDia agrupa las reservas vendidas confirmadas desde la fecha
// indicada por día (zona local del servidor).
func (r *metricasRepository) ContarVentasPorDia(ctx context.Context, desde time.Time, desplazamientoSegundos int) ([]CantidadPorDia, error) {
	var resultado []CantidadPorDia
	err := r.base.WithContext(ctx).
		Model(&models.Reserva{}).
		Select("date(updated_at + make_interval(secs => ?))::text AS fecha, count(*) AS cantidad", desplazamientoSegundos).
		Where("estado = ? AND updated_at >= ?", models.EstadoReservaVendida, desde).
		Group("fecha").
		Order("fecha").
		Scan(&resultado).Error
	if err != nil {
		return nil, err
	}
	return resultado, nil
}

// SumarIngresosPorPeriodo suma el precio de los vehículos vendidos desde la
// fecha indicada.
func (r *metricasRepository) SumarIngresosPorPeriodo(ctx context.Context, desde time.Time) (float64, error) {
	var total float64
	err := r.base.WithContext(ctx).
		Model(&models.Reserva{}).
		Select("COALESCE(SUM(v.precio), 0)").
		Joins("JOIN vehiculos v ON v.id = reservas.vehiculo_id").
		Where("reservas.estado = ? AND reservas.updated_at >= ?", models.EstadoReservaVendida, desde).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// ContarVentasPorMarca agrupa las reservas vendidas por la marca del vehículo.
func (r *metricasRepository) ContarVentasPorMarca(ctx context.Context) ([]ConteoPorMarca, error) {
	var resultado []ConteoPorMarca
	err := r.base.WithContext(ctx).
		Model(&models.Reserva{}).
		Select("v.marca AS marca, count(*) AS cantidad").
		Joins("JOIN vehiculos v ON v.id = reservas.vehiculo_id").
		Where("reservas.estado = ?", models.EstadoReservaVendida).
		Group("v.marca").
		Order("cantidad DESC, v.marca").
		Scan(&resultado).Error
	if err != nil {
		return nil, err
	}
	return resultado, nil
}

// ContarTestDrivesPorDia agrupa los turnos agendados (solicitados o
// confirmados) creados desde la fecha indicada por día.
func (r *metricasRepository) ContarTestDrivesPorDia(ctx context.Context, desde time.Time, desplazamientoSegundos int) ([]CantidadPorDia, error) {
	var resultado []CantidadPorDia
	// Se agrupa por la expresión completa (no por el alias "fecha") porque la
	// tabla ya tiene una columna llamada fecha y Postgres la priorizaría.
	expresion := fmt.Sprintf("date(created_at + make_interval(secs => %d))::text", desplazamientoSegundos)
	err := r.base.WithContext(ctx).
		Model(&models.TurnoTestDrive{}).
		Select(expresion+" AS fecha, count(*) AS cantidad").
		Where("estado IN ? AND created_at >= ?", []string{models.EstadoTurnoSolicitado, models.EstadoTurnoConfirmado}, desde).
		Group(expresion).
		Order(expresion).
		Scan(&resultado).Error
	if err != nil {
		return nil, err
	}
	return resultado, nil
}

// TraerVehiculosConAntiguedad devuelve los primeros limite vehículos
// disponibles con mayor tiempo publicado.
func (r *metricasRepository) TraerVehiculosConAntiguedad(ctx context.Context, limite int) ([]VehiculoConAntiguedad, error) {
	var resultado []VehiculoConAntiguedad
	err := r.base.WithContext(ctx).
		Model(&models.Vehiculo{}).
		Select("id, marca, modelo, anio, created_at").
		Where("estado = ?", models.EstadoDisponible).
		Order("created_at").
		Limit(limite).
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
