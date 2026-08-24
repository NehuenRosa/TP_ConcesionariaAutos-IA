package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"strings"
	"time"

	"concesionaria/backend/internal/googleid"
	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"
	"concesionaria/backend/internal/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// DuracionToken define el tiempo de vida de los tokens JWT (24 horas).
const DuracionToken = 24 * time.Hour

// Errores de negocio de la autenticación.
var (
	// ErrEmailEnUso indica que el email ya está registrado.
	ErrEmailEnUso = errors.New("el email ya está registrado")
	// ErrCredencialesInvalidas indica que el email o la contraseña son incorrectos.
	ErrCredencialesInvalidas = errors.New("email o contraseña incorrectos")
	// ErrDatosRegistroInvalidos indica que los datos del registro no son válidos.
	ErrDatosRegistroInvalidos = errors.New("datos de registro inválidos")
	// ErrUsuarioNoEncontrado indica que el usuario no existe.
	ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")
	// ErrGoogleNoDisponible indica que el acceso con Google no está habilitado
	// o el servicio de Google no está accesible.
	ErrGoogleNoDisponible = errors.New("el inicio de sesión con Google no está disponible")
	// ErrCredencialGoogleInvalida indica que el credential recibido no es un
	// ID token válido de Google para esta aplicación.
	ErrCredencialGoogleInvalida = errors.New("la credencial de Google no es válida")
)

// EntradaRegistro son los datos para crear una cuenta nueva.
type EntradaRegistro struct {
	Nombre   string
	Email    string
	Password string
}

// EntradaLogin son las credenciales para iniciar sesión.
type EntradaLogin struct {
	Email    string
	Password string
}

// VerificadorGoogle verifica credenciales (ID tokens) emitidas por Google.
type VerificadorGoogle interface {
	Verificar(ctx context.Context, idToken string) (*googleid.Identidad, error)
}

// AutenticacionService define el contrato de la lógica de negocio de
// autenticación y usuarios.
type AutenticacionService interface {
	Registrar(ctx context.Context, entrada EntradaRegistro) (*models.Usuario, error)
	IniciarSesion(ctx context.Context, entrada EntradaLogin) (*models.Usuario, string, error)
	IniciarSesionConGoogle(ctx context.Context, credencial string) (*models.Usuario, string, error)
	GoogleHabilitado() bool
	ObtenerPorID(ctx context.Context, id uint) (*models.Usuario, error)
}

// autenticacionService implementa AutenticacionService.
type autenticacionService struct {
	repositorio repositories.UsuarioRepository
	secreto     string
	duracion    time.Duration
	verificador VerificadorGoogle
}

// NuevoAutenticacionService crea un servicio de autenticación. El verificador
// de Google puede ser nil: en ese caso el acceso con Google queda
// deshabilitado.
func NuevoAutenticacionService(repositorio repositories.UsuarioRepository, secreto string, duracion time.Duration, verificador VerificadorGoogle) AutenticacionService {
	return &autenticacionService{repositorio: repositorio, secreto: secreto, duracion: duracion, verificador: verificador}
}

// Registrar valida los datos, hashea la contraseña y crea el usuario con rol
// cliente.
func (s *autenticacionService) Registrar(ctx context.Context, entrada EntradaRegistro) (*models.Usuario, error) {
	nombre := strings.TrimSpace(entrada.Nombre)
	email := strings.ToLower(strings.TrimSpace(entrada.Email))

	if err := validarRegistro(nombre, email, entrada.Password); err != nil {
		return nil, err
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
		Rol:      models.RolCliente,
	}
	return s.repositorio.Crear(ctx, usuario)
}

// IniciarSesion valida las credenciales y emite un token JWT de 24 horas.
func (s *autenticacionService) IniciarSesion(ctx context.Context, entrada EntradaLogin) (*models.Usuario, string, error) {
	email := strings.ToLower(strings.TrimSpace(entrada.Email))
	usuario, err := s.repositorio.ObtenerPorEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrCredencialesInvalidas
		}
		return nil, "", fmt.Errorf("obtener usuario por email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(entrada.Password)); err != nil {
		return nil, "", ErrCredencialesInvalidas
	}

	tokenString, err := token.Generar(usuario.ID, usuario.Rol, s.secreto, s.duracion)
	if err != nil {
		return nil, "", err
	}
	return usuario, tokenString, nil
}

// ObtenerPorID devuelve un usuario por su identificador.
func (s *autenticacionService) ObtenerPorID(ctx context.Context, id uint) (*models.Usuario, error) {
	usuario, err := s.repositorio.ObtenerPorID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUsuarioNoEncontrado
		}
		return nil, fmt.Errorf("obtener usuario por ID: %w", err)
	}
	return usuario, nil
}

// GoogleHabilitado informa si el acceso con Google está disponible.
func (s *autenticacionService) GoogleHabilitado() bool {
	return s.verificador != nil
}

// IniciarSesionConGoogle verifica el credential (ID token) de Google,
// vincula o crea el usuario y emite un JWT propio del sistema.
func (s *autenticacionService) IniciarSesionConGoogle(ctx context.Context, credencial string) (*models.Usuario, string, error) {
	if s.verificador == nil {
		return nil, "", ErrGoogleNoDisponible
	}

	identidad, err := s.verificador.Verificar(ctx, credencial)
	if err != nil {
		slog.Warn("credencial de Google rechazada", "error", err)
		return nil, "", ErrCredencialGoogleInvalida
	}

	email := strings.ToLower(strings.TrimSpace(identidad.Email))
	sub := identidad.Sub

	usuario, err := s.repositorio.ObtenerPorEmail(ctx, email)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// Alta automática como cliente; sin contraseña local.
		usuario = &models.Usuario{
			Nombre:    identidad.Nombre,
			Email:     email,
			Password:  "",
			Rol:       models.RolCliente,
			Proveedor: models.ProveedorGoogle,
			GoogleSub: &sub,
		}
		if _, err := s.repositorio.Crear(ctx, usuario); err != nil {
			return nil, "", fmt.Errorf("crear usuario desde Google: %w", err)
		}
	case err != nil:
		return nil, "", fmt.Errorf("obtener usuario por email: %w", err)
	default:
		// Vinculación automática de cuentas por email: se conserva rol,
		// nombre y contraseña del usuario existente.
		if usuario.Proveedor != models.ProveedorGoogle || usuario.GoogleSub == nil || *usuario.GoogleSub != sub {
			usuario.Proveedor = models.ProveedorGoogle
			usuario.GoogleSub = &sub
			if err := s.repositorio.Actualizar(ctx, usuario); err != nil {
				return nil, "", fmt.Errorf("vincular cuenta de Google: %w", err)
			}
		}
	}

	tokenString, err := token.Generar(usuario.ID, usuario.Rol, s.secreto, s.duracion)
	if err != nil {
		return nil, "", err
	}
	return usuario, tokenString, nil
}

// validarRegistro valida los datos del registro.
func validarRegistro(nombre string, email string, password string) error {
	if nombre == "" {
		return ErrDatosRegistroInvalidos
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrDatosRegistroInvalidos
	}
	if len(password) < 8 {
		return ErrDatosRegistroInvalidos
	}
	return nil
}
