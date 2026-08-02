package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// EntradaRegistro es el cuerpo de la solicitud de registro.
type EntradaRegistro struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// EntradaLogin es el cuerpo de la solicitud de inicio de sesión.
type EntradaLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RespuestaLogin envuelve el token y el perfil del usuario autenticado.
type RespuestaLogin struct {
	Token   string         `json:"token"`
	Usuario *models.Usuario `json:"usuario"`
}

// AutenticacionHandler agrupa los handlers de autenticación y usuarios.
type AutenticacionHandler struct {
	servicio services.AutenticacionService
}

// NuevoAutenticacionHandler crea un handler de autenticación.
func NuevoAutenticacionHandler(servicio services.AutenticacionService) *AutenticacionHandler {
	return &AutenticacionHandler{servicio: servicio}
}

// Registrar crea una cuenta nueva con rol cliente.
func (h *AutenticacionHandler) Registrar(c *gin.Context) {
	var entrada EntradaRegistro
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de solicitud inválido"})
		return
	}

	usuario, err := h.servicio.Registrar(c.Request.Context(), services.EntradaRegistro{
		Nombre:   entrada.Nombre,
		Email:    entrada.Email,
		Password: entrada.Password,
	})
	if err != nil {
		if errors.Is(err, services.ErrDatosRegistroInvalidos) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrEmailEnUso) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la cuenta"})
		return
	}

	c.JSON(http.StatusCreated, usuario)
}

// IniciarSesion valida las credenciales y responde con el token JWT.
func (h *AutenticacionHandler) IniciarSesion(c *gin.Context) {
	var entrada EntradaLogin
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de solicitud inválido"})
		return
	}

	usuario, tokenString, err := h.servicio.IniciarSesion(c.Request.Context(), services.EntradaLogin{
		Email:    entrada.Email,
		Password: entrada.Password,
	})
	if err != nil {
		if errors.Is(err, services.ErrCredencialesInvalidas) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar sesión"})
		return
	}

	c.JSON(http.StatusOK, RespuestaLogin{Token: tokenString, Usuario: usuario})
}

// Perfil responde el usuario autenticado a partir del token del request.
func (h *AutenticacionHandler) Perfil(c *gin.Context) {
	id, err := strconv.ParseUint(c.GetString("usuario_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	usuario, err := h.servicio.ObtenerPorID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrUsuarioNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el perfil"})
		return
	}

	c.JSON(http.StatusOK, usuario)
}
