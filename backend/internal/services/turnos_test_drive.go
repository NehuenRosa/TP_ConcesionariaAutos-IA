package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

// Errores de negocio de los turnos de test drive.
var (
	// ErrTurnoNoEncontrado indica que el turno no existe o no pertenece al cliente.
	ErrTurnoNoEncontrado = errors.New("turno de test drive no encontrado")
	// ErrDatosTurnoInvalidos indica que la fecha o la franja horaria no son válidas.
	ErrDatosTurnoInvalidos = errors.New("fecha o franja horaria inválida")
	// ErrTurnoEnPasado indica que la franja elegida para hoy ya comenzó.
	ErrTurnoEnPasado = errors.New("esa hora ya pasó para hoy, elegí otro horario")
	// ErrTurnoSuperpuesto indica que ya hay un turno activo para la misma
	// unidad en la misma fecha y franja.
	ErrTurnoSuperpuesto = errors.New("el turno ya está ocupado para esa unidad, fecha y franja")
	// ErrTurnoDuplicadoVehiculo indica que el cliente ya tiene un turno activo
	// para el mismo vehículo.
	ErrTurnoDuplicadoVehiculo = errors.New("ya tenés un test drive pendiente para este vehículo, esperá a que lo completen o cancelelo")
	// ErrTurnoEstadoInvalido indica que la transición de estado solicitada no es válida.
	ErrTurnoEstadoInvalido = errors.New("no se puede cambiar el estado del turno")
	// ErrFiltroEstadoTurnoInvalido indica que el filtro de estado no es válido.
	ErrFiltroEstadoTurnoInvalido = errors.New("filtro de estado inválido")
)

// TurnoTestDriveService define el contrato de la lógica de negocio de turnos.
type TurnoTestDriveService interface {
	// Solicitar crea un turno solicitado para un vehículo disponible.
	Solicitar(ctx context.Context, clienteID uint, vehiculoID uint, fecha string, franja string) (*models.TurnoTestDrive, error)
	// ListarMisTurnos lista los turnos de un cliente.
	ListarMisTurnos(ctx context.Context, clienteID uint) ([]models.TurnoTestDrive, error)
	// Cancelar cancela un turno propio en estado solicitado o confirmado.
	Cancelar(ctx context.Context, turnoID uint, clienteID uint) (*models.TurnoTestDrive, error)
	// Listar lista los turnos para el vendedor, con filtro de estado opcional.
	Listar(ctx context.Context, estado string) ([]models.TurnoTestDrive, error)
	// Confirmar confirma un turno solicitado.
	Confirmar(ctx context.Context, turnoID uint) (*models.TurnoTestDrive, error)
	// CancelarComoVendedor cancela un turno en estado solicitado o confirmado.
	CancelarComoVendedor(ctx context.Context, turnoID uint) (*models.TurnoTestDrive, error)
	// Completar marca como completado un turno confirmado.
	Completar(ctx context.Context, turnoID uint) (*models.TurnoTestDrive, error)
	// Franjas devuelve el catálogo de franjas horarias predefinidas.
	Franjas() []models.FranjaHoraria
	// FranjasConDisponibilidad devuelve el catálogo de franjas marcando las
	// ocupadas para una unidad y fecha. Si faltan vehiculoID o fecha, devuelve
	// el catálogo completo sin marcar ocupadas.
	FranjasConDisponibilidad(ctx context.Context, vehiculoID uint, fecha string) ([]models.FranjaHoraria, error)
}

// turnoTestDriveService implementa TurnoTestDriveService.
type turnoTestDriveService struct {
	repositorio  repositories.TurnoTestDriveRepository
	vehiculos    repositories.VehiculoRepository
}

// NuevoTurnoTestDriveService crea un servicio de turnos de test drive.
func NuevoTurnoTestDriveService(
	repositorio repositories.TurnoTestDriveRepository,
	vehiculos repositories.VehiculoRepository,
) TurnoTestDriveService {
	return &turnoTestDriveService{
		repositorio: repositorio,
		vehiculos:   vehiculos,
	}
}

// Solicitar valida el vehículo, la fecha y la franja, verifica que el cliente
// no tenga otro turno activo para el mismo vehículo ni que la franja esté
// ocupada, y crea el turno en estado solicitado.
func (s *turnoTestDriveService) Solicitar(ctx context.Context, clienteID uint, vehiculoID uint, fecha string, franja string) (*models.TurnoTestDrive, error) {
	if !models.FranjaValida(franja) {
		return nil, ErrDatosTurnoInvalidos
	}
	if !esFechaValida(fecha) {
		return nil, ErrDatosTurnoInvalidos
	}
	if !esHoraDisponible(fecha, franja) {
		return nil, ErrTurnoEnPasado
	}

	vehiculo, err := s.vehiculos.ObtenerPorID(ctx, vehiculoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehiculoNoDisponible
		}
		return nil, fmt.Errorf("obtener vehículo: %w", err)
	}
	if vehiculo.Estado != models.EstadoDisponible {
		return nil, ErrVehiculoNoDisponible
	}

	turno := &models.TurnoTestDrive{
		VehiculoID: vehiculoID,
		ClienteID:  clienteID,
		Fecha:      fecha,
		Franja:     franja,
		Estado:     models.EstadoTurnoSolicitado,
	}
	creado, err := s.repositorio.CrearSiDisponible(ctx, turno)
	if err != nil {
		if errors.Is(err, repositories.ErrFranjaOcupada) {
			return nil, ErrTurnoSuperpuesto
		}
		if errors.Is(err, repositories.ErrClienteYaTieneTurno) {
			return nil, ErrTurnoDuplicadoVehiculo
		}
		return nil, err
	}
	return creado, nil
}

// ListarMisTurnos lista los turnos de un cliente.
func (s *turnoTestDriveService) ListarMisTurnos(ctx context.Context, clienteID uint) ([]models.TurnoTestDrive, error) {
	return s.repositorio.ListarPorCliente(ctx, clienteID)
}

// Cancelar cancela un turno propio en estado solicitado o confirmado.
// Los turnos ajenos se tratan como inexistentes para no revelar su existencia.
func (s *turnoTestDriveService) Cancelar(ctx context.Context, turnoID uint, clienteID uint) (*models.TurnoTestDrive, error) {
	turno, err := s.repositorio.ObtenerPorID(ctx, turnoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTurnoNoEncontrado
		}
		return nil, fmt.Errorf("obtener turno: %w", err)
	}
	if turno.ClienteID != clienteID {
		return nil, ErrTurnoNoEncontrado
	}
	if !turno.EsActivo() {
		return nil, ErrTurnoEstadoInvalido
	}

	turno.Estado = models.EstadoTurnoCancelado
	return s.repositorio.Actualizar(ctx, turno)
}

// Listar lista los turnos para el vendedor con filtro de estado opcional.
func (s *turnoTestDriveService) Listar(ctx context.Context, estado string) ([]models.TurnoTestDrive, error) {
	if estado != "" && !esEstadoTurnoValido(estado) {
		return nil, ErrFiltroEstadoTurnoInvalido
	}
	return s.repositorio.Listar(ctx, estado)
}

// Confirmar confirma un turno en estado solicitado.
func (s *turnoTestDriveService) Confirmar(ctx context.Context, turnoID uint) (*models.TurnoTestDrive, error) {
	turno, err := s.repositorio.ObtenerPorID(ctx, turnoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTurnoNoEncontrado
		}
		return nil, fmt.Errorf("obtener turno: %w", err)
	}
	if turno.Estado != models.EstadoTurnoSolicitado {
		return nil, ErrTurnoEstadoInvalido
	}

	turno.Estado = models.EstadoTurnoConfirmado
	return s.repositorio.Actualizar(ctx, turno)
}

// CancelarComoVendedor cancela un turno en estado solicitado o confirmado.
func (s *turnoTestDriveService) CancelarComoVendedor(ctx context.Context, turnoID uint) (*models.TurnoTestDrive, error) {
	turno, err := s.repositorio.ObtenerPorID(ctx, turnoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTurnoNoEncontrado
		}
		return nil, fmt.Errorf("obtener turno: %w", err)
	}
	if !turno.EsActivo() {
		return nil, ErrTurnoEstadoInvalido
	}

	turno.Estado = models.EstadoTurnoCancelado
	return s.repositorio.Actualizar(ctx, turno)
}

// Completar marca como completado un turno en estado confirmado.
func (s *turnoTestDriveService) Completar(ctx context.Context, turnoID uint) (*models.TurnoTestDrive, error) {
	turno, err := s.repositorio.ObtenerPorID(ctx, turnoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTurnoNoEncontrado
		}
		return nil, fmt.Errorf("obtener turno: %w", err)
	}
	if turno.Estado != models.EstadoTurnoConfirmado {
		return nil, ErrTurnoEstadoInvalido
	}

	turno.Estado = models.EstadoTurnoCompletado
	return s.repositorio.Actualizar(ctx, turno)
}

// Franjas devuelve el catálogo de franjas horarias predefinidas.
func (s *turnoTestDriveService) Franjas() []models.FranjaHoraria {
	return models.FranjasDisponibles()
}

// FranjasConDisponibilidad devuelve las franjas del catálogo con el flag
// Ocupada seteado según los turnos activos de la unidad y fecha pedidas.
func (s *turnoTestDriveService) FranjasConDisponibilidad(ctx context.Context, vehiculoID uint, fecha string) ([]models.FranjaHoraria, error) {
	franjas := models.FranjasDisponibles()
	if vehiculoID == 0 || fecha == "" {
		return franjas, nil
	}

	ocupadas, err := s.repositorio.Ocupadas(ctx, vehiculoID, fecha)
	if err != nil {
		return nil, err
	}
	ocupado := make(map[string]bool, len(ocupadas))
	for _, id := range ocupadas {
		ocupado[id] = true
	}
	for i := range franjas {
		franjas[i].Ocupada = ocupado[franjas[i].ID]
	}
	return franjas, nil
}

// esFechaValida valida el formato YYYY-MM-DD y que la fecha no sea anterior a
// hoy. Compara por componentes en la zona local: Truncate sobre la hora actual
// recorta al mediodía UTC y descalificaba la fecha de hoy en zonas al oeste.
func esFechaValida(fecha string) bool {
	parseada, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return false
	}
	ahora := time.Now()
	hoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, ahora.Location())
	pedida := time.Date(parseada.Year(), parseada.Month(), parseada.Day(), 0, 0, 0, 0, ahora.Location())
	return !pedida.Before(hoy)
}

// esHoraDisponible indica si la franja elegida todavía no comenzó cuando la
// fecha pedida es hoy. Para cualquier otra fecha siempre está disponible.
func esHoraDisponible(fecha string, franja string) bool {
	return esHoraDisponibleEn(fecha, franja, time.Now())
}

// esHoraDisponibleEn compara la franja contra un instante dado (inyectado para
// poder testear). Solo bloquea franjas del día actual que ya hayan comenzado.
func esHoraDisponibleEn(fecha string, franja string, ahora time.Time) bool {
	if fecha != ahora.Format("2006-01-02") {
		return true
	}
	for _, disponible := range models.FranjasDisponibles() {
		if disponible.ID == franja {
			inicio, err := time.Parse("15:04", disponible.Inicio)
			if err != nil {
				return false
			}
			comienzo := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), inicio.Hour(), inicio.Minute(), 0, 0, ahora.Location())
			return comienzo.After(ahora)
		}
	}
	return false
}

// esEstadoTurnoValido indica si el estado es uno de los conocidos.
func esEstadoTurnoValido(estado string) bool {
	switch estado {
	case models.EstadoTurnoSolicitado, models.EstadoTurnoConfirmado, models.EstadoTurnoCancelado, models.EstadoTurnoCompletado:
		return true
	default:
		return false
	}
}
