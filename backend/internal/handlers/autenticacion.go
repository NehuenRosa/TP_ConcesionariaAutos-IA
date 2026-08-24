package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

// EntradaGoogle es el cuerpo de la solicitud de inicio de sesión con Google:
// el credential (ID token) emitido por Google Identity Services.
type EntradaGoogle struct {
	Credencial string `json:"credencial"`
}

// RespuestaLogin envuelve el token y el perfil del usuario autenticado.
type RespuestaLogin struct {
	Token   string         `json:"token"`
	Usuario *models.Usuario `json:"usuario"`
}

// AutenticacionHandler agrupa los handlers de autenticación y usuarios.
type AutenticacionHandler struct {
	servicio         services.AutenticacionService
	googleHabilitado bool
	clienteIDGoogle  string
}

// NuevoAutenticacionHandler crea un handler de autenticación. El client ID de
// Google se expone en /auth/proveedores para que el frontend inicialice
// Google Identity Services.
func NuevoAutenticacionHandler(servicio services.AutenticacionService, googleHabilitado bool, clienteIDGoogle string) *AutenticacionHandler {
	return &AutenticacionHandler{servicio: servicio, googleHabilitado: googleHabilitado, clienteIDGoogle: clienteIDGoogle}
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

// IniciarSesionConGoogle valida el credential de Google y responde con el
// token JWT propio del sistema.
func (h *AutenticacionHandler) IniciarSesionConGoogle(c *gin.Context) {
	var entrada EntradaGoogle
	if err := c.ShouldBindJSON(&entrada); err != nil || strings.TrimSpace(entrada.Credencial) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falta la credencial de Google"})
		return
	}

	usuario, tokenString, err := h.servicio.IniciarSesionConGoogle(c.Request.Context(), entrada.Credencial)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrGoogleNoDisponible):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "El inicio de sesión con Google no está disponible"})
		case errors.Is(err, services.ErrCredencialGoogleInvalida):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "La credencial de Google no es válida o está vencida"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar sesión con Google"})
		}
		return
	}

	c.JSON(http.StatusOK, RespuestaLogin{Token: tokenString, Usuario: usuario})
}

// Proveedores informa qué métodos de inicio de sesión están habilitados.
func (h *AutenticacionHandler) Proveedores(c *gin.Context) {
	respuesta := gin.H{"google": false}
	if h.googleHabilitado && h.servicio.GoogleHabilitado() {
		respuesta["google"] = true
		respuesta["client_id"] = h.clienteIDGoogle
	}
	c.JSON(http.StatusOK, respuesta)
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
