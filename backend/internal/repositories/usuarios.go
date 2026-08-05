package repositories

import (
	"context"
	"fmt"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// UsuarioRepository define el acceso a datos de usuarios sobre GORM.
type UsuarioRepository interface {
	// Crear persiste un usuario nuevo.
	Crear(ctx context.Context, usuario *models.Usuario) (*models.Usuario, error)
	// ObtenerPorEmail devuelve un usuario por su email o un error de GORM
	// (ErrRecordNotFound si no existe).
	ObtenerPorEmail(ctx context.Context, email string) (*models.Usuario, error)
	// ObtenerPorID devuelve un usuario por su identificador.
	ObtenerPorID(ctx context.Context, id uint) (*models.Usuario, error)
	// Listar devuelve todos los usuarios ordenados por ID.
	Listar(ctx context.Context) ([]models.Usuario, error)
	// Actualizar persiste los cambios de un usuario existente.
	Actualizar(ctx context.Context, usuario *models.Usuario) error
	// Eliminar borra un usuario por su identificador.
	Eliminar(ctx context.Context, id uint) error
}

// usuarioRepository implementa UsuarioRepository sobre GORM.
type usuarioRepository struct {
	base *gorm.DB
}

// NuevoUsuarioRepository crea un repositorio de usuarios.
func NuevoUsuarioRepository(base *gorm.DB) UsuarioRepository {
	return &usuarioRepository{base: base}
}

// Crear persiste el usuario y devuelve el registro con su ID asignado.
func (r *usuarioRepository) Crear(ctx context.Context, usuario *models.Usuario) (*models.Usuario, error) {
	if err := r.base.WithContext(ctx).Create(usuario).Error; err != nil {
		return nil, fmt.Errorf("crear usuario: %w", err)
	}
	return usuario, nil
}

// ObtenerPorEmail devuelve un usuario por su email.
func (r *usuarioRepository) ObtenerPorEmail(ctx context.Context, email string) (*models.Usuario, error) {
	var usuario models.Usuario
	if err := r.base.WithContext(ctx).
		Where("email = ?", email).
		First(&usuario).Error; err != nil {
		return nil, err
	}
	return &usuario, nil
}

// ObtenerPorID devuelve un usuario por su identificador.
func (r *usuarioRepository) ObtenerPorID(ctx context.Context, id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	if err := r.base.WithContext(ctx).
		First(&usuario, id).Error; err != nil {
		return nil, err
	}
	return &usuario, nil
}

// Listar devuelve todos los usuarios ordenados por ID.
func (r *usuarioRepository) Listar(ctx context.Context) ([]models.Usuario, error) {
	var usuarios []models.Usuario
	if err := r.base.WithContext(ctx).
		Order("id ASC").
		Find(&usuarios).Error; err != nil {
		return nil, fmt.Errorf("listar usuarios: %w", err)
	}
	return usuarios, nil
}

// Actualizar persiste los cambios del usuario.
func (r *usuarioRepository) Actualizar(ctx context.Context, usuario *models.Usuario) error {
	if err := r.base.WithContext(ctx).
		Save(usuario).Error; err != nil {
		return fmt.Errorf("actualizar usuario: %w", err)
	}
	return nil
}

// Eliminar borra un usuario por su identificador.
func (r *usuarioRepository) Eliminar(ctx context.Context, id uint) error {
	if err := r.base.WithContext(ctx).
		Delete(&models.Usuario{}, id).Error; err != nil {
		return fmt.Errorf("eliminar usuario: %w", err)
	}
	return nil
}
