package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"concesionaria/backend/internal/repositories"
)

// Errores de negocio de las métricas del panel de administración.
var (
	// ErrPeriodoInvalido indica que el parámetro de período no es válido.
	ErrPeriodoInvalido = errors.New("el período debe ser 7, 30 o 90 días")
)

// periodosMetricasValidos son los valores permitidos para el parámetro periodo.
var periodosMetricasValidos = map[int]bool{7: true, 30: true, 90: true}

// limiteVehiculosEnStock es el tope de vehículos del gráfico de antigüedad.
const limiteVehiculosEnStock = 8

// Metricas es el payload completo del panel de administración.
type Metricas struct {
	VehiculosPorEstado    []repositories.ConteoPorEstado `json:"vehiculosPorEstado"`
	ConsultasPorPeriodo   []repositories.CantidadPorDia  `json:"consultasPorPeriodo"`
	VentasPorPeriodo      []repositories.CantidadPorDia  `json:"ventasPorPeriodo"`
	IngresoPorPeriodo     float64                         `json:"ingresoPorPeriodo"`
	VentasPorMarca        []repositories.ConteoPorMarca   `json:"ventasPorMarca"`
	TestDrivesPorPeriodo  []repositories.CantidadPorDia   `json:"testDrivesPorPeriodo"`
	VehiculosEnStock      []repositories.VehiculoEnStock  `json:"vehiculosEnStock"`
	ReservasActivas       int64                           `json:"reservasActivas"`
	ReservasVendidas      int64                           `json:"reservasVendidas"`
	TestDrivesAgendados   int64                           `json:"testDrivesAgendados"`
	TestDrivesCompletados int64                           `json:"testDrivesCompletados"`
	ConsultasAbiertas     int64                           `json:"consultasAbiertas"`
	TotalUsuarios         int64                           `json:"totalUsuarios"`
}

// MetricasService define el contrato de las métricas del panel.
type MetricasService interface {
	// Obtener arma el payload de métricas con el período indicado.
	Obtener(ctx context.Context, periodo string) (*Metricas, error)
}

// metricasService implementa MetricasService.
type metricasService struct {
	repositorio repositories.MetricasRepository
}

// NuevoMetricasService crea un servicio de métricas.
func NuevoMetricasService(repositorio repositories.MetricasRepository) MetricasService {
	return &metricasService{repositorio: repositorio}
}

// Obtener valida el período, consulta las agregaciones y arma el payload,
// rellenando con ceros los días del período sin consultas.
func (s *metricasService) Obtener(ctx context.Context, periodo string) (*Metricas, error) {
	dias, err := interpretarPeriodo(periodo)
	if err != nil {
		return nil, err
	}

	vehiculos, err := s.repositorio.ContarVehiculosPorEstado(ctx)
	if err != nil {
		return nil, fmt.Errorf("contar vehículos por estado: %w", err)
	}

	inicio := inicioPeriodo(dias)
	_, desplazamientoSegundos := time.Now().Zone()
	consultas, err := s.repositorio.ContarConsultasPorDia(ctx, inicio, desplazamientoSegundos)
	if err != nil {
		return nil, fmt.Errorf("contar consultas por día: %w", err)
	}

	ventas, err := s.repositorio.ContarVentasPorDia(ctx, inicio, desplazamientoSegundos)
	if err != nil {
		return nil, fmt.Errorf("contar ventas por día: %w", err)
	}

	ingreso, err := s.repositorio.SumarIngresosPorPeriodo(ctx, inicio)
	if err != nil {
		return nil, fmt.Errorf("sumar ingresos por período: %w", err)
	}

	ventasPorMarca, err := s.repositorio.ContarVentasPorMarca(ctx)
	if err != nil {
		return nil, fmt.Errorf("contar ventas por marca: %w", err)
	}

	testDrives, err := s.repositorio.ContarTestDrivesPorDia(ctx, inicio, desplazamientoSegundos)
	if err != nil {
		return nil, fmt.Errorf("contar test drives por día: %w", err)
	}

	vehiculosAntiguos, err := s.repositorio.TraerVehiculosConAntiguedad(ctx, limiteVehiculosEnStock)
	if err != nil {
		return nil, fmt.Errorf("traer vehículos con antigüedad: %w", err)
	}

	reservasActivas, err := s.repositorio.ContarReservasActivas(ctx)
	if err != nil {
		return nil, fmt.Errorf("contar reservas activas: %w", err)
	}

	reservasVendidas, err := s.repositorio.ContarReservasVendidas(ctx)
	if err != nil {
		return nil, fmt.Errorf("contar reservas vendidas: %w", err)
	}

	testDrivesAgendados, err := s.repositorio.ContarTestDrivesAgendados(ctx)
	if err != nil {
		return nil, fmt.Errorf("contar test drives agendados: %w", err)
	}

	testDrivesCompletados, err := s.repositorio.ContarTestDrivesCompletados(ctx)
	if err != nil {
		return nil, fmt.Errorf("contar test drives completados: %w", err)
	}

	consultasAbiertas, err := s.repositorio.ContarConsultasAbiertas(ctx)
	if err != nil {
		return nil, fmt.Errorf("contar consultas abiertas: %w", err)
	}

	totalUsuarios, err := s.repositorio.ContarUsuarios(ctx)
	if err != nil {
		return nil, fmt.Errorf("contar usuarios: %w", err)
	}

	return &Metricas{
		VehiculosPorEstado:    vehiculos,
		ConsultasPorPeriodo:   completarSerie(consultas, dias),
		VentasPorPeriodo:      completarSerie(ventas, dias),
		IngresoPorPeriodo:     ingreso,
		VentasPorMarca:        ventasPorMarca,
		TestDrivesPorPeriodo:  completarSerie(testDrives, dias),
		VehiculosEnStock:      conDiasEnStock(vehiculosAntiguos),
		ReservasActivas:       reservasActivas,
		ReservasVendidas:      reservasVendidas,
		TestDrivesAgendados:   testDrivesAgendados,
		TestDrivesCompletados: testDrivesCompletados,
		ConsultasAbiertas:     consultasAbiertas,
		TotalUsuarios:         totalUsuarios,
	}, nil
}

// conDiasEnStock convierte los vehículos con su fecha de alta al DTO del panel
// calculando la cantidad de días que llevan publicados.
func conDiasEnStock(vehiculos []repositories.VehiculoConAntiguedad) []repositories.VehiculoEnStock {
	ahora := time.Now()
	resultado := make([]repositories.VehiculoEnStock, 0, len(vehiculos))
	for _, vehiculo := range vehiculos {
		dias := int(ahora.Sub(vehiculo.CreatedAt).Hours() / 24)
		if dias < 0 {
			dias = 0
		}
		resultado = append(resultado, repositories.VehiculoEnStock{
			ID:          vehiculo.ID,
			Marca:       vehiculo.Marca,
			Modelo:      vehiculo.Modelo,
			Anio:        vehiculo.Anio,
			DiasEnStock: dias,
		})
	}
	return resultado
}

// interpretarPeriodo devuelve la cantidad de días del período solicitado,
// aplicando el valor por defecto de 30 días cuando no se envía.
func interpretarPeriodo(periodo string) (int, error) {
	if periodo == "" {
		return 30, nil
	}
	dias, err := strconv.Atoi(periodo)
	if err != nil || !periodosMetricasValidos[dias] {
		return 0, ErrPeriodoInvalido
	}
	return dias, nil
}

// inicioPeriodo devuelve el inicio (00:00, zona horaria local) del día más
// antiguo del período.
func inicioPeriodo(dias int) time.Time {
	inicio := time.Now().AddDate(0, 0, -(dias - 1))
	return time.Date(inicio.Year(), inicio.Month(), inicio.Day(), 0, 0, 0, 0, inicio.Location())
}

// completarSerie rellena con ceros los días del período que no tienen
// eventos para que la serie tenga exactamente un registro por día.
func completarSerie(consultas []repositories.CantidadPorDia, dias int) []repositories.CantidadPorDia {
	porDia := make(map[string]int64, len(consultas))
	for _, consulta := range consultas {
		porDia[consulta.Fecha] = consulta.Cantidad
	}

	inicio := inicioPeriodo(dias)
	serie := make([]repositories.CantidadPorDia, 0, dias)
	for dia := inicio; dia.Before(inicio.AddDate(0, 0, dias)); dia = dia.AddDate(0, 0, 1) {
		clave := dia.Format("2006-01-02")
		serie = append(serie, repositories.CantidadPorDia{Fecha: clave, Cantidad: porDia[clave]})
	}
	return serie
}
