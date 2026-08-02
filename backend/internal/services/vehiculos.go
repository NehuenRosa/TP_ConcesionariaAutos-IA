package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

// Errores de negocio del catálogo y la gestión de vehículos.
var (
	// ErrVehiculoNoEncontrado indica que el vehículo no existe o no está disponible.
	ErrVehiculoNoEncontrado = errors.New("vehículo no encontrado o no disponible")
	// ErrPaginacionInvalida indica que la página o el tamaño solicitado no son válidos.
	ErrPaginacionInvalida = errors.New("paginación inválida")
	// ErrDatosVehiculoInvalidos indica que la ficha técnica enviada no es válida.
	ErrDatosVehiculoInvalidos = errors.New("ficha técnica del vehículo inválida")
	// ErrEstadoInvalido indica que el estado del vehículo no es uno de los conocidos.
	ErrEstadoInvalido = errors.New("estado de vehículo inválido")
	// ErrFiltroEstadoInvalido indica que el filtro de estado de la consulta no es válido.
	ErrFiltroEstadoInvalido = errors.New("filtro de estado inválido")
	// ErrFiltroInvalido indica que la búsqueda, un filtro o el ordenamiento
	// solicitado en el catálogo no son válidos.
	ErrFiltroInvalido = errors.New("filtro de búsqueda inválido")
)

// FiltrosBusqueda agrupa los criterios opcionales de búsqueda y filtrado del
// catálogo público. Es el mismo tipo que usa el repositorio.
type FiltrosBusqueda = repositories.FiltrosBusqueda

// EntradaVehiculo es el conjunto de datos para crear o actualizar un vehículo.
type EntradaVehiculo struct {
	Marca       string
	Modelo      string
	Anio        int
	Kilometraje int
	Combustible string
	Transmision string
	Tipo        string
	Precio      float64
	Condicion   string
	Estado      string
	Imagenes    []string
}

// VehiculoService define el contrato de la lógica de negocio de vehículos.
type VehiculoService interface {
	ListarDisponibles(ctx context.Context, filtros FiltrosBusqueda, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error)
	ListarParaGestion(ctx context.Context, estado string, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	ObtenerParaGestion(ctx context.Context, id uint) (*models.Vehiculo, error)
	Crear(ctx context.Context, entrada EntradaVehiculo) (*models.Vehiculo, error)
	Actualizar(ctx context.Context, id uint, entrada EntradaVehiculo) (*models.Vehiculo, error)
	DarDeBaja(ctx context.Context, id uint) (*models.Vehiculo, error)
}

// vehiculoService implementa VehiculoService.
type vehiculoService struct {
	repositorio repositories.VehiculoRepository
}

// NuevoVehiculoService crea un servicio de vehículos.
func NuevoVehiculoService(repositorio repositories.VehiculoRepository) VehiculoService {
	return &vehiculoService{repositorio: repositorio}
}

// ListarDisponibles valida la paginación y los filtros, y delega en el
// repositorio el listado de vehículos con estado "disponible".
func (s *vehiculoService) ListarDisponibles(ctx context.Context, filtros FiltrosBusqueda, pagina int, tamano int) ([]models.Vehiculo, int64, error) {
	if pagina < 1 || tamano < 1 {
		return nil, 0, ErrPaginacionInvalida
	}
	if err := validarFiltros(filtros); err != nil {
		return nil, 0, err
	}
	return s.repositorio.Listar(ctx, models.EstadoDisponible, filtros, pagina, tamano)
}

// ObtenerPorID devuelve un vehículo solo si existe y está disponible.
// Para los demás estados se retorna ErrVehiculoNoEncontrado, ocultando la
// existencia de unidades no comercializables.
func (s *vehiculoService) ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error) {
	vehiculo, err := s.repositorio.ObtenerPorID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehiculoNoEncontrado
		}
		return nil, fmt.Errorf("obtener vehículo por ID: %w", err)
	}
	if vehiculo.Estado != models.EstadoDisponible {
		return nil, ErrVehiculoNoEncontrado
	}
	return vehiculo, nil
}

// ListarParaGestion valida la paginación y el filtro de estado, y delega en el
// repositorio el listado administrativo de todos los estados.
func (s *vehiculoService) ListarParaGestion(ctx context.Context, estado string, pagina int, tamano int) ([]models.Vehiculo, int64, error) {
	if pagina < 1 || tamano < 1 {
		return nil, 0, ErrPaginacionInvalida
	}
	if estado != "" && !esEstadoValido(estado) {
		return nil, 0, ErrFiltroEstadoInvalido
	}
	return s.repositorio.ListarParaGestion(ctx, estado, pagina, tamano)
}

// ObtenerParaGestion devuelve un vehículo en cualquier estado o
// ErrVehiculoNoEncontrado si no existe.
func (s *vehiculoService) ObtenerParaGestion(ctx context.Context, id uint) (*models.Vehiculo, error) {
	vehiculo, err := s.repositorio.ObtenerPorID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehiculoNoEncontrado
		}
		return nil, fmt.Errorf("obtener vehículo para gestión por ID: %w", err)
	}
	return vehiculo, nil
}

// Crear valida la ficha técnica y persiste un vehículo nuevo.
func (s *vehiculoService) Crear(ctx context.Context, entrada EntradaVehiculo) (*models.Vehiculo, error) {
	if err := validarEntrada(entrada); err != nil {
		return nil, err
	}
	estado := entrada.Estado
	if estado == "" {
		estado = models.EstadoDisponible
	}

	vehiculo := aModelo(entrada, estado)
	return s.repositorio.Crear(ctx, vehiculo)
}

// Actualizar valida la ficha técnica y actualiza un vehículo existente,
// reemplazando su galería de imágenes.
func (s *vehiculoService) Actualizar(ctx context.Context, id uint, entrada EntradaVehiculo) (*models.Vehiculo, error) {
	if err := validarEntrada(entrada); err != nil {
		return nil, err
	}
	if _, err := s.ObtenerParaGestion(ctx, id); err != nil {
		return nil, err
	}

	estado := entrada.Estado
	if estado == "" {
		estado = models.EstadoDisponible
	}

	vehiculo := aModelo(entrada, estado)
	vehiculo.ID = id
	return s.repositorio.Actualizar(ctx, vehiculo)
}

// DarDeBaja cambia el estado de un vehículo a dado_de_baja sin eliminarlo.
// Es idempotente: si ya está dado de baja, responde con el vehículo actual.
func (s *vehiculoService) DarDeBaja(ctx context.Context, id uint) (*models.Vehiculo, error) {
	vehiculo, err := s.ObtenerParaGestion(ctx, id)
	if err != nil {
		return nil, err
	}
	if vehiculo.Estado == models.EstadoDadoDeBaja {
		return vehiculo, nil
	}

	if err := s.repositorio.DarDeBaja(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehiculoNoEncontrado
		}
		return nil, fmt.Errorf("dar de baja vehículo: %w", err)
	}
	return s.repositorio.ObtenerPorID(ctx, id)
}

// validarFiltros valida los criterios de búsqueda y ordenamiento del catálogo.
func validarFiltros(filtros FiltrosBusqueda) error {
	if filtros.AnioMin != nil && filtros.AnioMax != nil && *filtros.AnioMin > *filtros.AnioMax {
		return ErrFiltroInvalido
	}
	if filtros.PrecioMin != nil && filtros.PrecioMax != nil && *filtros.PrecioMin > *filtros.PrecioMax {
		return ErrFiltroInvalido
	}
	if filtros.Condicion != "" && filtros.Condicion != models.CondicionNuevo && filtros.Condicion != models.CondicionUsado {
		return ErrFiltroInvalido
	}
	switch filtros.OrdenPor {
	case "", "precio", "anio":
	default:
		return ErrFiltroInvalido
	}
	switch filtros.OrdenDireccion {
	case "", "asc", "desc":
	default:
		return ErrFiltroInvalido
	}
	return nil
}

// validarEntrada valida los campos de la ficha técnica y el estado.
func validarEntrada(entrada EntradaVehiculo) error {
	if strings.TrimSpace(entrada.Marca) == "" || strings.TrimSpace(entrada.Modelo) == "" {
		return ErrDatosVehiculoInvalidos
	}
	if entrada.Anio < 1900 || entrada.Anio > time.Now().Year()+1 {
		return ErrDatosVehiculoInvalidos
	}
	if entrada.Precio <= 0 {
		return ErrDatosVehiculoInvalidos
	}
	if entrada.Condicion != models.CondicionNuevo && entrada.Condicion != models.CondicionUsado {
		return ErrDatosVehiculoInvalidos
	}
	if strings.TrimSpace(entrada.Tipo) == "" {
		return ErrDatosVehiculoInvalidos
	}
	if entrada.Estado != "" && !esEstadoValido(entrada.Estado) {
		return ErrEstadoInvalido
	}
	return nil
}

// esEstadoValido indica si el estado es uno de los conocidos del modelo.
func esEstadoValido(estado string) bool {
	switch estado {
	case models.EstadoDisponible, models.EstadoReservado, models.EstadoVendido, models.EstadoDadoDeBaja:
		return true
	default:
		return false
	}
}

// aModelo convierte una entrada validada en un modelo de GORM con sus imágenes.
func aModelo(entrada EntradaVehiculo, estado string) *models.Vehiculo {
	imagenes := make([]models.Imagen, 0, len(entrada.Imagenes))
	for _, url := range entrada.Imagenes {
		if strings.TrimSpace(url) == "" {
			continue
		}
		imagenes = append(imagenes, models.Imagen{URL: url})
	}

	return &models.Vehiculo{
		Marca:       entrada.Marca,
		Modelo:      entrada.Modelo,
		Anio:        entrada.Anio,
		Kilometraje: entrada.Kilometraje,
		Combustible: entrada.Combustible,
		Transmision: entrada.Transmision,
		Tipo:        entrada.Tipo,
		Precio:      entrada.Precio,
		Condicion:   entrada.Condicion,
		Estado:      estado,
		Imagenes:    imagenes,
	}
}
