package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Errores de negocio de la gestión de usuarios.
var (
	// ErrDatosUsuarioInvalidos indica que los datos del usuario no son válidos.
	ErrDatosUsuarioInvalidos = errors.New("datos de usuario inválidos")
	// ErrRolInvalido indica que el rol no es uno de los permitidos.
	ErrRolInvalido = errors.New("rol inválido")
	// ErrNoPuedeModificarPropioUsuario indica que el administrador no puede
	// modificar ni eliminar su propia cuenta.
	ErrNoPuedeModificarPropioUsuario = errors.New("no se puede modificar tu propia cuenta")
)

// EntradaUsuarioAdmin son los datos para crear o actualizar un usuario desde
// el panel de administración. Si Password queda vacío en una actualización,
// se conserva la contraseña actual.
type EntradaUsuarioAdmin struct {
	Nombre   string
	Email    string
	Password string
	Rol      string
}

// UsuariosService define el contrato de la lógica de negocio de gestión de
// usuarios por parte del administrador.
type UsuariosService interface {
	// Listar devuelve todos los usuarios del sistema.
	Listar(ctx context.Context) ([]models.Usuario, error)
	// Crear da de alta un usuario con el rol indicado.
	Crear(ctx context.Context, entrada EntradaUsuarioAdmin) (*models.Usuario, error)
	// Actualizar modifica los datos de un usuario existente.
	Actualizar(ctx context.Context, id uint, entrada EntradaUsuarioAdmin, idSolicitante uint) (*models.Usuario, error)
	// Eliminar da de baja a un usuario por su identificador.
	Eliminar(ctx context.Context, id uint, idSolicitante uint) error
}

// usuariosService implementa UsuariosService.
type usuariosService struct {
	repositorio repositories.UsuarioRepository
}

// NuevoUsuariosService crea un servicio de gestión de usuarios.
func NuevoUsuariosService(repositorio repositories.UsuarioRepository) UsuariosService {
	return &usuariosService{repositorio: repositorio}
}

// Listar devuelve todos los usuarios del sistema.
func (s *usuariosService) Listar(ctx context.Context) ([]models.Usuario, error) {
	return s.repositorio.Listar(ctx)
}

// Crear valida los datos, hashea la contraseña y crea el usuario con el rol
// solicitado por el administrador.
func (s *usuariosService) Crear(ctx context.Context, entrada EntradaUsuarioAdmin) (*models.Usuario, error) {
	nombre := strings.TrimSpace(entrada.Nombre)
	email := strings.ToLower(strings.TrimSpace(entrada.Email))

	if err := validarRegistro(nombre, email, entrada.Password); err != nil {
		return nil, ErrDatosUsuarioInvalidos
	}
	if !rolValido(entrada.Rol) {
		return nil, ErrRolInvalido
	}

	if _, err := s.repositorio.ObtenerPorEmail(ctx, email); err == nil {
		return nil, ErrEmailEnUso
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("verificar email existente: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(entrada.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generar hash de contraseña: %w", err)
	}

	usuario := &models.Usuario{
		Nombre:   nombre,
		Email:    email,
		Password: string(hash),
		Rol:      entrada.Rol,
	}
	return s.repositorio.Crear(ctx, usuario)
}

// Actualizar modifica los datos de un usuario. Si la contraseña llega vacía se
// conserva la actual, y el administrador no puede cambiar el rol de su propia
// cuenta para evitar quedarse sin permisos.
func (s *usuariosService) Actualizar(ctx context.Context, id uint, entrada EntradaUsuarioAdmin, idSolicitante uint) (*models.Usuario, error) {
	usuario, err := s.repositorio.ObtenerPorID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUsuarioNoEncontrado
		}
		return nil, fmt.Errorf("obtener usuario por ID: %w", err)
	}

	nombre := strings.TrimSpace(entrada.Nombre)
	email := strings.ToLower(strings.TrimSpace(entrada.Email))
	if nombre == "" {
		return nil, ErrDatosUsuarioInvalidos
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrDatosUsuarioInvalidos
	}
	if !rolValido(entrada.Rol) {
		return nil, ErrRolInvalido
	}
	if id == idSolicitante && entrada.Rol != usuario.Rol {
		return nil, ErrNoPuedeModificarPropioUsuario
	}
	if entrada.Password != "" && len(entrada.Password) < 8 {
		return nil, ErrDatosUsuarioInvalidos
	}

	otro, err := s.repositorio.ObtenerPorEmail(ctx, email)
	if err == nil && otro.ID != id {
		return nil, ErrEmailEnUso
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("verificar email existente: %w", err)
	}

	usuario.Nombre = nombre
	usuario.Email = email
	usuario.Rol = entrada.Rol
	if entrada.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(entrada.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("generar hash de contraseña: %w", err)
		}
		usuario.Password = string(hash)
	}

	if err := s.repositorio.Actualizar(ctx, usuario); err != nil {
		return nil, err
	}
	return usuario, nil
}

// Eliminar da de baja a un usuario. No se permite eliminar la propia cuenta.
func (s *usuariosService) Eliminar(ctx context.Context, id uint, idSolicitante uint) error {
	if id == idSolicitante {
		return ErrNoPuedeModificarPropioUsuario
	}

	usuario, err := s.repositorio.ObtenerPorID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUsuarioNoEncontrado
		}
		return fmt.Errorf("obtener usuario por ID: %w", err)
	}

	if err := s.repositorio.Eliminar(ctx, usuario.ID); err != nil {
		return fmt.Errorf("no se pudo eliminar el usuario (puede tener consultas o turnos asociados): %w", err)
	}
	return nil
}

// rolValido indica si el rol pertenece al catálogo de roles del sistema.
func rolValido(rol string) bool {
	switch rol {
	case models.RolCliente, models.RolVendedor, models.RolAdministrador:
		return true
	default:
		return false
	}
}
