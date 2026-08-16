package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

// Errores de negocio de consultas.
var (
	// ErrConsultaNoEncontrada indica que la consulta no existe.
	ErrConsultaNoEncontrada = errors.New("consulta no encontrada")
	// ErrVehiculoNoDisponible indica que el vehículo no está disponible para consulta.
	ErrVehiculoNoDisponible = errors.New("vehículo no disponible")
	// ErrMensajeVacio indica que el mensaje está vacío.
	ErrMensajeVacio = errors.New("el mensaje no puede estar vacío")
	// ErrConsultaYaTomada indica que la consulta ya fue tomada por otro vendedor.
	ErrConsultaYaTomada = errors.New("la consulta ya fue tomada por otro vendedor")
	// ErrConsultaNoPendiente indica que la consulta no está en estado pendiente.
	ErrConsultaNoPendiente = errors.New("la consulta no está pendiente")
	// ErrConsultaNoCerrada indica que la consulta no está cerrada.
	ErrConsultaNoCerrada = errors.New("la consulta no está cerrada")
	// ErrConsultaYaCerrada indica que la consulta ya está cerrada.
	ErrConsultaYaCerrada = errors.New("la consulta ya está cerrada")
	// ErrNoEsVendedorAsignado indica que el usuario no es el vendedor asignado.
	ErrNoEsVendedorAsignado = errors.New("no es el vendedor asignado a esta consulta")
	// ErrNoEsParticipante indica que el usuario no participa en la consulta.
	ErrNoEsParticipante = errors.New("no es participante de esta consulta")
	// ErrConsultaCerradaMensajes indica que no se pueden enviar mensajes a una consulta cerrada.
	ErrConsultaCerradaMensajes = errors.New("no se pueden enviar mensajes a una consulta cerrada")
)

// ConsultaService define el contrato de la lógica de negocio de consultas.
type ConsultaService interface {
	// Crear crea una nueva consulta con el primer mensaje.
	Crear(ctx context.Context, clienteID uint, vehiculoID uint, mensaje string) (*models.Consulta, error)
	// ObtenerPorID obtiene una consulta por su ID.
	ObtenerPorID(ctx context.Context, id uint) (*models.Consulta, error)
	// ListarPorCliente lista las consultas de un cliente.
	ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Consulta, error)
	// ListarPorVendedor lista las consultas asignadas a un vendedor.
	ListarPorVendedor(ctx context.Context, vendedorID uint) ([]models.Consulta, error)
	// ListarPorUsuario devuelve los IDs de consultas donde el usuario participa.
	ListarPorUsuario(ctx context.Context, usuarioID uint) ([]uint, error)
	// Tomar asigna un vendedor a una consulta pendiente.
	Tomar(ctx context.Context, consultaID uint, vendedorID uint) (*models.Consulta, error)
	// Cerrar cambia el estado de una consulta a cerrada.
	Cerrar(ctx context.Context, consultaID uint, vendedorID uint) (*models.Consulta, error)
	// Eliminar elimina una consulta cerrada.
	Eliminar(ctx context.Context, consultaID uint, vendedorID uint) error
	// EsParticipante verifica si un usuario es participante de la consulta.
	EsParticipante(ctx context.Context, consultaID uint, usuarioID uint) (bool, error)
	// EsVendedorAsignado verifica si un vendedor es el asignado a la consulta.
	EsVendedorAsignado(ctx context.Context, consultaID uint, vendedorID uint) (bool, error)
}

// consultaService implementa ConsultaService.
type consultaService struct {
	repositorio   repositories.ConsultaRepository
	vehiculos     repositories.VehiculoRepository
}

// NuevoConsultaService crea un servicio de consultas.
func NuevoConsultaService(
	repositorio repositories.ConsultaRepository,
	vehiculos repositories.VehiculoRepository,
) ConsultaService {
	return &consultaService{
		repositorio: repositorio,
		vehiculos:   vehiculos,
	}
}

// Crear crea una nueva consulta con el primer mensaje del cliente.
func (s *consultaService) Crear(ctx context.Context, clienteID uint, vehiculoID uint, mensaje string) (*models.Consulta, error) {
	if strings.TrimSpace(mensaje) == "" {
		return nil, ErrMensajeVacio
	}

	// Verificar que el vehículo existe y está disponible
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

	consulta := &models.Consulta{
		VehiculoID: vehiculoID,
		ClienteID:  clienteID,
		Estado:     models.EstadoPendiente,
		Mensajes: []models.Mensaje{
			{
				RemitenteID: clienteID,
				Contenido:   strings.TrimSpace(mensaje),
			},
		},
	}

	return s.repositorio.Crear(ctx, consulta)
}

// ObtenerPorID obtiene una consulta por su ID.
func (s *consultaService) ObtenerPorID(ctx context.Context, id uint) (*models.Consulta, error) {
	consulta, err := s.repositorio.ObtenerPorID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConsultaNoEncontrada
		}
		return nil, fmt.Errorf("obtener consulta: %w", err)
	}
	return consulta, nil
}

// ListarPorCliente lista las consultas de un cliente.
func (s *consultaService) ListarPorCliente(ctx context.Context, clienteID uint) ([]models.Consulta, error) {
	return s.repositorio.ListarPorCliente(ctx, clienteID)
}

// ListarPorVendedor lista las consultas asignadas a un vendedor.
func (s *consultaService) ListarPorVendedor(ctx context.Context, vendedorID uint) ([]models.Consulta, error) {
	return s.repositorio.ListarPorVendedor(ctx, vendedorID)
}

// ListarPorUsuario devuelve los IDs de consultas donde el usuario participa.
func (s *consultaService) ListarPorUsuario(ctx context.Context, usuarioID uint) ([]uint, error) {
	return s.repositorio.ListarPorUsuario(ctx, usuarioID)
}

// Tomar asigna un vendedor a una consulta pendiente. Usa un UPDATE atómico
// para que dos vendedores no puedan tomar la misma consulta a la vez.
func (s *consultaService) Tomar(ctx context.Context, consultaID uint, vendedorID uint) (*models.Consulta, error) {
	tomada, err := s.repositorio.TomarSiPendiente(ctx, consultaID, vendedorID)
	if err != nil {
		return nil, err
	}

	if !tomada {
		// Distinguir "no existe" de "ya no está pendiente".
		if _, err := s.ObtenerPorID(ctx, consultaID); err != nil {
			return nil, err
		}
		return nil, ErrConsultaNoPendiente
	}

	return s.repositorio.ObtenerPorID(ctx, consultaID)
}

// Cerrar cambia el estado de una consulta a cerrada.
func (s *consultaService) Cerrar(ctx context.Context, consultaID uint, vendedorID uint) (*models.Consulta, error) {
	consulta, err := s.repositorio.ObtenerPorID(ctx, consultaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConsultaNoEncontrada
		}
		return nil, fmt.Errorf("obtener consulta: %w", err)
	}

	if consulta.Estado != models.EstadoEnConversacion {
		return nil, ErrConsultaYaCerrada
	}

	if consulta.VendedorID == nil || *consulta.VendedorID != vendedorID {
		return nil, ErrNoEsVendedorAsignado
	}

	consulta.Estado = models.EstadoCerrada

	return s.repositorio.Actualizar(ctx, consulta)
}

// Eliminar elimina una consulta cerrada.
func (s *consultaService) Eliminar(ctx context.Context, consultaID uint, vendedorID uint) error {
	consulta, err := s.repositorio.ObtenerPorID(ctx, consultaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConsultaNoEncontrada
		}
		return fmt.Errorf("obtener consulta: %w", err)
	}

	if consulta.Estado != models.EstadoCerrada {
		return ErrConsultaNoCerrada
	}

	if consulta.VendedorID == nil || *consulta.VendedorID != vendedorID {
		return ErrNoEsVendedorAsignado
	}

	return s.repositorio.Eliminar(ctx, consultaID)
}

// EsParticipante verifica si un usuario es el cliente o el vendedor asignado.
func (s *consultaService) EsParticipante(ctx context.Context, consultaID uint, usuarioID uint) (bool, error) {
	consulta, err := s.repositorio.ObtenerPorID(ctx, consultaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("obtener consulta: %w", err)
	}

	if consulta.ClienteID == usuarioID {
		return true, nil
	}
	if consulta.VendedorID != nil && *consulta.VendedorID == usuarioID {
		return true, nil
	}

	return false, nil
}

// EsVendedorAsignado verifica si un vendedor es el asignado a la consulta.
func (s *consultaService) EsVendedorAsignado(ctx context.Context, consultaID uint, vendedorID uint) (bool, error) {
	consulta, err := s.repositorio.ObtenerPorID(ctx, consultaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("obtener consulta: %w", err)
	}

	return consulta.VendedorID != nil && *consulta.VendedorID == vendedorID, nil
}
