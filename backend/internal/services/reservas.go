package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"path/filepath"
	"strings"
	"time"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

// Reglas de la seña de la reserva (CU-08).
const (
	// PorcentajeSena es la fracción del precio del vehículo que el cliente
	// transfiere para reservar (5 %). Se calcula siempre en el backend.
	PorcentajeSena = 0.05
	// PlazoComprobante es el tiempo que tiene el cliente para subir el
	// comprobante antes de que la reserva se anule automáticamente.
	PlazoComprobante = 2 * time.Hour
	// MaximoPesoComprobanteBytes limita el tamaño de la imagen del comprobante.
	MaximoPesoComprobanteBytes = 5 * 1024 * 1024
)

// formatosComprobante mapea las extensiones admitidas del comprobante al MIME
// correspondiente.
var formatosComprobante = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

// Errores de negocio de las reservas.
var (
	// ErrReservaNoEncontrada indica que la reserva no existe o no pertenece al
	// cliente.
	ErrReservaNoEncontrada = errors.New("reserva no encontrada")
	// ErrVehiculoYaNoDisponible indica que la unidad ya no está disponible para
	// reservar (fue reservada o vendida por otro cliente).
	ErrVehiculoYaNoDisponible = errors.New("la unidad ya no está disponible para reservar")
	// ErrReservaEstadoInvalido indica que la transición de estado solicitada no
	// es válida.
	ErrReservaEstadoInvalido = errors.New("no se puede cambiar el estado de la reserva")
	// ErrFiltroEstadoReservaInvalido indica que el filtro de estado no es válido.
	ErrFiltroEstadoReservaInvalido = errors.New("filtro de estado inválido")
	// ErrComprobanteInvalido indica que el archivo subido no es una imagen
	// admitida o supera el tamaño máximo.
	ErrComprobanteInvalido = errors.New("comprobante inválido")
	// ErrComprobanteFueraDePlazo indica que el plazo de 2 horas venció sin
	// comprobante y la reserva fue anulada automáticamente.
	ErrComprobanteFueraDePlazo = errors.New("el plazo para enviar el comprobante venció")
	// ErrComprobanteNoEncontrado indica que la reserva aún no tiene comprobante.
	ErrComprobanteNoEncontrado = errors.New("la reserva todavía no tiene un comprobante")
	// ErrReservaProhibida indica que el usuario no puede ver el comprobante de
	// una reserva ajena.
	ErrReservaProhibida = errors.New("no tenés permisos sobre esta reserva")
	// ErrMotivoRequerido indica que el vendedor intentó cancelar una reserva
	// sin el motivo obligatorio.
	ErrMotivoRequerido = errors.New("tenés que indicar el motivo de la cancelación")
)

// DatosTransferencia son los datos bancarios y el monto que el cliente usa
// para transferir la seña.
type DatosTransferencia struct {
	CBU   string  `json:"cbu"`
	Alias string  `json:"alias"`
	Monto float64 `json:"monto"`
}

// ReservaService define el contrato de la lógica de negocio de reservas.
type ReservaService interface {
	// Crear reserva un vehículo disponible: crea la reserva activa con plazo
	// de 2 horas para el comprobante y bloquea la unidad.
	Crear(ctx context.Context, clienteID uint, vehiculoID uint) (*models.Reserva, error)
	// ListarMisReservas lista las reservas de un cliente.
	ListarMisReservas(ctx context.Context, clienteID uint) ([]models.Reserva, error)
	// Cancelar cancela una reserva propia en estado activa y libera la unidad.
	Cancelar(ctx context.Context, reservaID uint, clienteID uint) (*models.Reserva, error)
	// Listar lista las reservas para el vendedor, con filtro de estado opcional.
	Listar(ctx context.Context, estado string) ([]models.Reserva, error)
	// ConfirmarVenta confirma la venta de una reserva activa.
	ConfirmarVenta(ctx context.Context, reservaID uint) (*models.Reserva, error)
	// CancelarComoVendedor cancela una reserva activa y libera la unidad,
	// registrando el motivo obligatorio que verá el cliente.
	CancelarComoVendedor(ctx context.Context, reservaID uint, motivo string) (*models.Reserva, error)
	// ObtenerDatosTransferencia devuelve CBU/alias de la concesionaria y el
	// monto de la seña (5 % del precio) del vehículo indicado.
	ObtenerDatosTransferencia(ctx context.Context, vehiculoID uint) (*DatosTransferencia, error)
	// SubirComprobante guarda la imagen del comprobante de una reserva propia
	// activa dentro del plazo.
	SubirComprobante(ctx context.Context, reservaID uint, clienteID uint, nombreArchivo string, datos []byte) (*models.Reserva, error)
	// ObtenerComprobante devuelve el comprobante cargado si el solicitante es
	// el dueño o personal (vendedor/administrador).
	ObtenerComprobante(ctx context.Context, reservaID uint, solicitanteID uint, esPersonal bool) (*models.ComprobanteReserva, error)
	// ExpirarVencidas anula las reservas activas vencidas sin comprobante y
	// libera sus unidades.
	ExpirarVencidas(ctx context.Context) error
}

// reservaService implementa ReservaService.
type reservaService struct {
	repositorio repositories.ReservaRepository
	vehiculos   repositories.VehiculoRepository
	// cbuConcesionaria y aliasConcesionaria se muestran al cliente para
	// transferir la seña; vacíos si no están configurados.
	cbuConcesionaria   string
	aliasConcesionaria string
}

// NuevoReservaService crea un servicio de reservas.
func NuevoReservaService(
	repositorio repositories.ReservaRepository,
	vehiculos repositories.VehiculoRepository,
	cbuConcesionaria string,
	aliasConcesionaria string,
) ReservaService {
	return &reservaService{
		repositorio:        repositorio,
		vehiculos:          vehiculos,
		cbuConcesionaria:   cbuConcesionaria,
		aliasConcesionaria: aliasConcesionaria,
	}
}

// Crear valida que el vehículo exista y esté disponible y crea la reserva en
// estado activa, bloqueando la unidad.
func (s *reservaService) Crear(ctx context.Context, clienteID uint, vehiculoID uint) (*models.Reserva, error) {
	vehiculo, err := s.vehiculos.ObtenerPorID(ctx, vehiculoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehiculoNoDisponible
		}
		return nil, fmt.Errorf("obtener vehículo: %w", err)
	}

	switch vehiculo.Estado {
	case models.EstadoDisponible:
		// La unidad puede reservarse.
	case models.EstadoReservado, models.EstadoVendido:
		return nil, ErrVehiculoYaNoDisponible
	default:
		// Dado de baja: no es comercializable.
		return nil, ErrVehiculoNoDisponible
	}

	reserva := &models.Reserva{
		VehiculoID:             vehiculoID,
		ClienteID:              clienteID,
		Estado:                 models.EstadoReservaActiva,
		VencimientoComprobante: time.Now().Add(PlazoComprobante),
	}
	reservaCreada, err := s.repositorio.CrearYReservar(ctx, reserva)
	if err != nil {
		if errors.Is(err, repositories.ErrVehiculoYaNoDisponible) {
			return nil, ErrVehiculoYaNoDisponible
		}
		return nil, fmt.Errorf("crear reserva: %w", err)
	}
	return reservaCreada, nil
}

// ListarMisReservas lista las reservas de un cliente.
func (s *reservaService) ListarMisReservas(ctx context.Context, clienteID uint) ([]models.Reserva, error) {
	return s.repositorio.ListarPorCliente(ctx, clienteID)
}

// Cancelar cancela una reserva propia en estado activa y libera la unidad. Las
// reservas ajenas se tratan como inexistentes para no revelar su existencia.
func (s *reservaService) Cancelar(ctx context.Context, reservaID uint, clienteID uint) (*models.Reserva, error) {
	reserva, err := s.repositorio.ObtenerPorID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservaNoEncontrada
		}
		return nil, fmt.Errorf("obtener reserva: %w", err)
	}
	if reserva.ClienteID != clienteID {
		return nil, ErrReservaNoEncontrada
	}
	if !reserva.EsActiva() {
		return nil, ErrReservaEstadoInvalido
	}
	// Si la reserva ya venció sin comprobante, se aplica primero su anulación
	// automática y la cancelación manual queda rechazada como estado inválido.
	s.expirarSiCorresponde(ctx, reserva)

	reservaCancelada, err := s.repositorio.CancelarYLiberar(ctx, reserva)
	if err != nil {
		if errors.Is(err, repositories.ErrReservaYaNoActiva) {
			return nil, ErrReservaEstadoInvalido
		}
		return nil, fmt.Errorf("cancelar reserva: %w", err)
	}
	return reservaCancelada, nil
}

// Listar lista las reservas para el vendedor con filtro de estado opcional.
func (s *reservaService) Listar(ctx context.Context, estado string) ([]models.Reserva, error) {
	if estado != "" && !esEstadoReservaValido(estado) {
		return nil, ErrFiltroEstadoReservaInvalido
	}
	return s.repositorio.Listar(ctx, estado)
}

// ConfirmarVenta confirma la venta de una reserva en estado activa. Si la
// reserva venció el plazo sin comprobante, se aplica antes la expiración
// automática y la confirmación queda rechazada.
func (s *reservaService) ConfirmarVenta(ctx context.Context, reservaID uint) (*models.Reserva, error) {
	reserva, err := s.repositorio.ObtenerPorID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservaNoEncontrada
		}
		return nil, fmt.Errorf("obtener reserva: %w", err)
	}
	if !reserva.EsActiva() {
		return nil, ErrReservaEstadoInvalido
	}
	s.expirarSiCorresponde(ctx, reserva)

	reservaVendida, err := s.repositorio.ConfirmarVentaYMarcarVendido(ctx, reserva)
	if err != nil {
		if errors.Is(err, repositories.ErrReservaYaNoActiva) {
			return nil, ErrReservaEstadoInvalido
		}
		return nil, fmt.Errorf("confirmar venta: %w", err)
	}
	return reservaVendida, nil
}

// CancelarComoVendedor cancela una reserva en estado activa y libera la
// unidad, guardando el motivo obligatorio que el cliente verá en sus reservas.
func (s *reservaService) CancelarComoVendedor(ctx context.Context, reservaID uint, motivo string) (*models.Reserva, error) {
	motivo = strings.TrimSpace(motivo)
	if motivo == "" {
		return nil, ErrMotivoRequerido
	}
	reserva, err := s.repositorio.ObtenerPorID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservaNoEncontrada
		}
		return nil, fmt.Errorf("obtener reserva: %w", err)
	}
	if !reserva.EsActiva() {
		return nil, ErrReservaEstadoInvalido
	}
	s.expirarSiCorresponde(ctx, reserva)

	reserva.MotivoCancelacion = motivo
	reservaCancelada, err := s.repositorio.CancelarYLiberar(ctx, reserva)
	if err != nil {
		if errors.Is(err, repositories.ErrReservaYaNoActiva) {
			return nil, ErrReservaEstadoInvalido
		}
		return nil, fmt.Errorf("cancelar reserva: %w", err)
	}
	return reservaCancelada, nil
}

// esEstadoReservaValido indica si el estado es uno de los conocidos.
func esEstadoReservaValido(estado string) bool {
	switch estado {
	case models.EstadoReservaActiva, models.EstadoReservaVendida, models.EstadoReservaCancelada:
		return true
	default:
		return false
	}
}

// CalcularMontoSenia calcula el monto de la seña: el PorcentajeSena del precio,
// redondeado a dos decimales. El monto siempre se compone en el backend.
func CalcularMontoSenia(precio float64) float64 {
	return math.Round(precio*PorcentajeSena*100) / 100
}

// ObtenerDatosTransferencia devuelve los datos bancarios de la concesionaria y
// el monto de la seña del vehículo indicado, solo si está disponible.
func (s *reservaService) ObtenerDatosTransferencia(ctx context.Context, vehiculoID uint) (*DatosTransferencia, error) {
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
	return &DatosTransferencia{
		CBU:   s.cbuConcesionaria,
		Alias: s.aliasConcesionaria,
		Monto: CalcularMontoSenia(vehiculo.Precio),
	}, nil
}

// SubirComprobante valida y guarda la imagen del comprobante de una reserva
// propia activa dentro del plazo. Reenviar reemplaza la imagen anterior.
func (s *reservaService) SubirComprobante(ctx context.Context, reservaID uint, clienteID uint, nombreArchivo string, datos []byte) (*models.Reserva, error) {
	mime, err := validarComprobante(nombreArchivo, datos)
	if err != nil {
		return nil, err
	}

	reserva, err := s.repositorio.ObtenerPorID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservaNoEncontrada
		}
		return nil, fmt.Errorf("obtener reserva: %w", err)
	}
	if reserva.ClienteID != clienteID {
		return nil, ErrReservaNoEncontrada
	}
	if !reserva.EsActiva() {
		return nil, ErrReservaEstadoInvalido
	}
	if reserva.ComprobanteVencido(time.Now()) {
		expirada, err := s.repositorio.ExpirarSiVencida(ctx, reserva)
		if err != nil {
			return nil, fmt.Errorf("expirar reserva: %w", err)
		}
		if expirada {
			return nil, ErrComprobanteFueraDePlazo
		}
	}

	ahora := time.Now()
	reserva.ComprobanteEnviadoAt = &ahora
	comprobante := &models.ComprobanteReserva{MIME: mime, Datos: datos}
	if err := s.repositorio.GuardarComprobante(ctx, reserva, comprobante); err != nil {
		return nil, fmt.Errorf("guardar comprobante: %w", err)
	}
	return reserva, nil
}

// ObtenerComprobante devuelve el comprobante cargado si el solicitante es el
// dueño de la reserva o personal (vendedor/administrador).
func (s *reservaService) ObtenerComprobante(ctx context.Context, reservaID uint, solicitanteID uint, esPersonal bool) (*models.ComprobanteReserva, error) {
	reserva, err := s.repositorio.ObtenerPorID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservaNoEncontrada
		}
		return nil, fmt.Errorf("obtener reserva: %w", err)
	}
	if !esPersonal && reserva.ClienteID != solicitanteID {
		return nil, ErrReservaProhibida
	}
	comprobante, err := s.repositorio.ObtenerComprobantePorReservaID(ctx, reservaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrComprobanteNoEncontrado
		}
		return nil, fmt.Errorf("obtener comprobante: %w", err)
	}
	return comprobante, nil
}

// ExpirarVencidas anula las reservas activas vencidas sin comprobante y libera
// sus unidades. La invoca periódicamente el job del main.
func (s *reservaService) ExpirarVencidas(ctx context.Context) error {
	if _, err := s.repositorio.ExpirarVencidas(ctx); err != nil {
		return fmt.Errorf("expirar reservas vencidas: %w", err)
	}
	return nil
}

// expirarSiCorresponde aplica la anulación automática cuando una reserva
// activa venció el plazo sin comprobante (chequeo perezoso).
func (s *reservaService) expirarSiCorresponde(ctx context.Context, reserva *models.Reserva) {
	if !reserva.ComprobanteVencido(time.Now()) {
		return
	}
	if _, err := s.repositorio.ExpirarSiVencida(ctx, reserva); err != nil {
		// El fallo del chequeo perezoso no interrumpe la operación: la
		// transición final revalida el estado sobre la base igualmente.
		slog.Warn("expiración perezosa de reserva falló", "reservaId", reserva.ID, "error", err)
	}
}

// validarComprobante verifica extensión y tamaño de la imagen y devuelve su
// MIME. Mismo criterio que las fotos de la tasación del chatbot.
func validarComprobante(nombreArchivo string, datos []byte) (string, error) {
	if len(datos) == 0 || len(datos) > MaximoPesoComprobanteBytes {
		return "", ErrComprobanteInvalido
	}
	extension := strings.ToLower(filepath.Ext(nombreArchivo))
	mime, ok := formatosComprobante[extension]
	if !ok {
		return "", ErrComprobanteInvalido
	}
	return mime, nil
}
