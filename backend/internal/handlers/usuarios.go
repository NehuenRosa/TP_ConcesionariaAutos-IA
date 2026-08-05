package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// EntradaUsuarioAdmin es el cuerpo de creación/edición de un usuario desde el
// panel de administración. En la edición, password puede venir vacío para
// conservar la contraseña actual.
type EntradaUsuarioAdmin struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Rol      string `json:"rol"`
}

// UsuariosHandler agrupa los handlers de gestión de usuarios (administrador).
type UsuariosHandler struct {
	servicio services.UsuariosService
}

// NuevoUsuariosHandler crea un handler de gestión de usuarios.
func NuevoUsuariosHandler(servicio services.UsuariosService) *UsuariosHandler {
	return &UsuariosHandler{servicio: servicio}
}

// Listar responde todos los usuarios del sistema.
func (h *UsuariosHandler) Listar(c *gin.Context) {
	usuarios, err := h.servicio.Listar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron listar los usuarios"})
		return
	}
	c.JSON(http.StatusOK, usuarios)
}

// Crear da de alta un usuario con el rol indicado.
func (h *UsuariosHandler) Crear(c *gin.Context) {
	var entrada EntradaUsuarioAdmin
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de solicitud inválido"})
		return
	}

	usuario, err := h.servicio.Crear(c.Request.Context(), services.EntradaUsuarioAdmin{
		Nombre:   entrada.Nombre,
		Email:    entrada.Email,
		Password: entrada.Password,
		Rol:      entrada.Rol,
	})
	if err != nil {
		h.responderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, usuario)
}

// Actualizar modifica los datos de un usuario.
func (h *UsuariosHandler) Actualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de usuario inválido"})
		return
	}

	var entrada EntradaUsuarioAdmin
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de solicitud inválido"})
		return
	}

	idSolicitante, err := idUsuarioActual(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	usuario, err := h.servicio.Actualizar(c.Request.Context(), uint(id), services.EntradaUsuarioAdmin{
		Nombre:   entrada.Nombre,
		Email:    entrada.Email,
		Password: entrada.Password,
		Rol:      entrada.Rol,
	}, idSolicitante)
	if err != nil {
		h.responderError(c, err)
		return
	}

	c.JSON(http.StatusOK, usuario)
}

// Eliminar da de baja a un usuario.
func (h *UsuariosHandler) Eliminar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador de usuario inválido"})
		return
	}

	idSolicitante, err := idUsuarioActual(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	if err := h.servicio.Eliminar(c.Request.Context(), uint(id), idSolicitante); err != nil {
		h.responderError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// responderError traduce los errores de negocio a códigos HTTP.
func (h *UsuariosHandler) responderError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrUsuarioNoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrEmailEnUso):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrRolInvalido),
		errors.Is(err, services.ErrDatosUsuarioInvalidos):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrNoPuedeModificarPropioUsuario):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo realizar la operación"})
	}
}

// idUsuarioActual devuelve el identificador del usuario autenticado.
func idUsuarioActual(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.GetString("usuario_id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
