package handlers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"concesionaria/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// extensionesImagenValidas son las extensiones de archivo aceptadas para la tasación.
var extensionesImagenValidas = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

// ChatbotHandler agrupa los handlers públicos del asistente conversacional.
type ChatbotHandler struct {
	servicio services.ChatbotService
}

// NuevoChatbotHandler crea un handler del chatbot.
func NuevoChatbotHandler(servicio services.ChatbotService) *ChatbotHandler {
	return &ChatbotHandler{servicio: servicio}
}

// peticionChat es el cuerpo de POST /api/chatbot/mensajes.
type peticionChat struct {
	Mensaje   string               `json:"mensaje"`
	Historial []services.TurnoChat `json:"historial"`
}

// Responder procesa el mensaje del usuario y responde con la respuesta del
// asistente en lenguaje natural. Si el usuario es un cliente autenticado y la
// IA detectó que quiere cotizar un vehículo, responde además con el id de la
// cotización creada para redirigir al panel.
func (h *ChatbotHandler) Responder(c *gin.Context) {
	var peticion peticionChat
	if err := c.ShouldBindJSON(&peticion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida"})
		return
	}

	clienteID, _ := leerIDOpcional(c)
	respuesta, err := h.servicio.Responder(c.Request.Context(), clienteID, peticion.Mensaje, peticion.Historial)
	if err != nil {
		if errors.Is(err, services.ErrMensajeVacio) || errors.Is(err, services.ErrMensajeMuyLargo) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo procesar el mensaje"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"respuesta": respuesta.Respuesta, "cotizacionId": respuesta.CotizacionID})
}

// Tasacion procesa las fotos del auto del usuario y responde con la
// estimación del valor de permuta en lenguaje natural. Devuelve además el
// identificador de la sesión para confirmar la visita.
func (h *ChatbotHandler) Tasacion(c *gin.Context) {
	formulario, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida: se esperaba multipart/form-data"})
		return
	}

	imagenes, err := leerImagenes(formulario.File["fotos"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	descripcion := ""
	if valores := formulario.Value["descripcion"]; len(valores) > 0 {
		descripcion = strings.TrimSpace(valores[0])
	}
	sesionID := ""
	if valores := formulario.Value["sesion_id"]; len(valores) > 0 {
		sesionID = strings.TrimSpace(valores[0])
	}

	resultado, err := h.servicio.Tasacion(c.Request.Context(), sesionID, descripcion, imagenes)
	if err != nil {
		if errors.Is(err, services.ErrSinImagenes) ||
			errors.Is(err, services.ErrDemasiadasImagenes) ||
			errors.Is(err, services.ErrImagenMuyPesada) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo realizar la tasación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"respuesta": resultado.Respuesta, "sesionId": resultado.SesionID})
}

// peticionConfirmarTasacion es el cuerpo de POST /api/chatbot/tasacion/confirmar.
type peticionConfirmarTasacion struct {
	SesionID string `json:"sesionId"`
	Mensaje  string `json:"mensaje"`
}

// ConfirmarTasacion coordina la visita del cliente a la concesionaria: extrae
// el día y la franja del mensaje, genera un código único y confirma la
// tasación pendiente de la sesión.
func (h *ChatbotHandler) ConfirmarTasacion(c *gin.Context) {
	var peticion peticionConfirmarTasacion
	if err := c.ShouldBindJSON(&peticion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida"})
		return
	}

	resultado, err := h.servicio.ConfirmarTasacion(c.Request.Context(), peticion.SesionID, peticion.Mensaje)
	if err != nil {
		if errors.Is(err, services.ErrTasacionSinSesion) ||
			errors.Is(err, services.ErrTasacionNoEncontrada) ||
			errors.Is(err, services.ErrMensajeVacio) ||
			errors.Is(err, services.ErrMensajeMuyLargo) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo confirmar la tasación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"respuesta": resultado.Respuesta, "confirmada": resultado.Confirmada})
}

// leerIDOpcional devuelve el id del usuario autenticado, o 0 si la sesión es
// anónima (endpoint público con autenticación opcional).
func leerIDOpcional(c *gin.Context) (uint, error) {
	id, err := extraerUsuarioID(c)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

// leerImagenes valida la extensión y el peso de cada archivo adjunto y los
// convierte en imágenes para la tasación.
func leerImagenes(archivos []*multipart.FileHeader) ([]services.ImagenTasacion, error) {
	if len(archivos) > services.MaximoImagenesTasacion {
		return nil, services.ErrDemasiadasImagenes
	}

	imagenes := make([]services.ImagenTasacion, 0, len(archivos))
	for _, archivo := range archivos {
		if archivo == nil {
			continue
		}
		if archivo.Size > services.MaximoPesoImagenBytes {
			return nil, services.ErrImagenMuyPesada
		}

		mime, ok := extensionesImagenValidas[strings.ToLower(filepath.Ext(archivo.Filename))]
		if !ok {
			return nil, errors.New("formato de imagen no soportado: solo se aceptan JPG, PNG o WebP")
		}

		contenido, err := archivo.Open()
		if err != nil {
			return nil, errors.New("no se pudo leer la imagen adjunta")
		}
		defer contenido.Close()

		datos, err := io.ReadAll(contenido)
		if err != nil {
			return nil, errors.New("no se pudo leer la imagen adjunta")
		}

		imagenes = append(imagenes, services.ImagenTasacion{MIME: mime, Datos: datos})
	}
	return imagenes, nil
}
