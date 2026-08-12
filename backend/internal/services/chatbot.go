package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// Límites de entrada del chatbot.
const (
	LargoMaximoMensaje     = 1000
	MaximoTurnosHistorial  = 20
	MaximoImagenesTasacion = 5
	MaximoPesoImagenBytes  = 5 * 1024 * 1024
	TimeoutChatbot         = 120 * time.Second
	TimeoutVision          = 120 * time.Second
	// Contexto (num_ctx) y salida máxima por modelo, ajustados para 8 GB de VRAM.
	NumCtxChatbot = 4096
	NumCtxVision  = 2048
	MaxTokensSalida = 600
)

// Errores de negocio del chatbot.
var (
	// ErrMensajeMuyLargo indica que el mensaje supera el largo máximo.
	ErrMensajeMuyLargo = errors.New("el mensaje es demasiado largo")
	// ErrSinImagenes indica que la tasación requiere al menos una foto o descripción.
	ErrSinImagenes = errors.New("se requiere al menos una foto o una descripción para tasar")
	// ErrDemasiadasImagenes indica que se superó el máximo de fotos por tasación.
	ErrDemasiadasImagenes = errors.New("se permiten como máximo 5 fotos por tasación")
	// ErrImagenMuyPesada indica que una foto supera el peso máximo.
	ErrImagenMuyPesada = errors.New("una de las fotos supera el peso máximo permitido")
)

// Mensajes de degradación cuando el modelo local no está disponible.
const (
	respuestaFallbackChat = "Disculpá, ahora mismo no puedo responder. " +
		"Podés consultar el catálogo, crear una consulta sobre un vehículo o pedir " +
		"un test drive desde el sitio, y un vendedor te va a ayudar."
	respuestaFallbackTasacion = "Disculpá, ahora mismo no puedo analizar las fotos. " +
		"Podés acercarte a la concesionaria o enviar una consulta desde el sitio, y " +
		"un vendedor se va a comunicar con vos para la tasación."
	respuestaNoIdentificado = "No pude identificar tu vehículo con certeza a partir de las fotos " +
		"y la descripción, así que prefiero no inventar un valor. " +
		"Escribí la marca, el modelo y el año en la descripción, o acercate a la " +
		"concesionaria donde un vendedor te va a hacer la tasación definitiva."
)

// TurnoChat es un mensaje previo de la conversación enviado por el frontend.
type TurnoChat struct {
	Rol      string `json:"rol"`
	Contenido string `json:"contenido"`
}

// ImagenTasacion es una imagen adjunta por el usuario para la tasación.
type ImagenTasacion struct {
	MIME  string
	Datos []byte
}

// ChatbotService define el contrato del asistente conversacional.
type ChatbotService interface {
	// Responder genera la respuesta del asistente para el mensaje del usuario,
	// usando el historial previo como contexto.
	Responder(ctx context.Context, mensaje string, historial []TurnoChat) (string, error)
	// Tasacion estima el valor de permuta del vehículo del usuario a partir de
	// las fotos y una descripción opcional.
	Tasacion(ctx context.Context, descripcion string, imagenes []ImagenTasacion) (string, error)
}

// chatbotService implementa ChatbotService con LangChain y Ollama.
type chatbotService struct {
	repositorioVehiculos repositories.VehiculoRepository
	ollamaURL            string
	modeloChatbot        string
	modeloVision         string
	precios              ServicioPrecios
}

// NuevoChatbotService crea el servicio del asistente conversacional.
func NuevoChatbotService(repositorioVehiculos repositories.VehiculoRepository, ollamaURL string, modeloChatbot string, modeloVision string, precios ServicioPrecios) ChatbotService {
	return &chatbotService{
		repositorioVehiculos: repositorioVehiculos,
		ollamaURL:            ollamaURL,
		modeloChatbot:        modeloChatbot,
		modeloVision:         modeloVision,
		precios:              precios,
	}
}

// Responder valida el mensaje, arma el contexto del stock y la conversación, y
// delega la generación de la respuesta en el modelo local. Si el modelo no
// está disponible, devuelve un mensaje de fallback en español.
func (s *chatbotService) Responder(ctx context.Context, mensaje string, historial []TurnoChat) (string, error) {
	if strings.TrimSpace(mensaje) == "" {
		return "", ErrMensajeVacio
	}
	if len([]rune(mensaje)) > LargoMaximoMensaje {
		return "", ErrMensajeMuyLargo
	}

	contexto, err := s.construirContextoStock(ctx)
	if err != nil {
		return "", fmt.Errorf("obtener stock para el chatbot: %w", err)
	}

	mensajes := s.construirMensajes(contexto, historial, mensaje)
	respuesta, err := s.generar(ctx, s.modeloChatbot, mensajes, NumCtxChatbot, TimeoutChatbot)
	if err != nil {
		slog.Error("Error al generar respuesta del chatbot", "error", err.Error())
		return respuestaFallbackChat, nil
	}
	return respuesta, nil
}

// Tasacion identifica el vehículo con el modelo de visión, consulta el valor
// de referencia oficial en la guía de precios y compone la respuesta con datos
// reales, sin inventar valores. Si algo falla, responde de forma honesta.
func (s *chatbotService) Tasacion(ctx context.Context, descripcion string, imagenes []ImagenTasacion) (string, error) {
	if len(imagenes) > MaximoImagenesTasacion {
		return "", ErrDemasiadasImagenes
	}
	for _, imagen := range imagenes {
		if len(imagen.Datos) > MaximoPesoImagenBytes {
			return "", ErrImagenMuyPesada
		}
	}
	if len(imagenes) == 0 && strings.TrimSpace(descripcion) == "" {
		return "", ErrSinImagenes
	}

	identificacion, err := s.identificarVehiculo(ctx, descripcion, imagenes)
	if err != nil {
		slog.Warn("No se pudo identificar el vehículo para tasar", "error", err.Error())
		return respuestaNoIdentificado, nil
	}

	referencia, err := s.precios.Buscar(ctx, identificacion.Marca, identificacion.Modelo, identificacion.Anio)
	if err != nil {
		slog.Warn("No se encontró valor de referencia para la tasación",
			"marca", identificacion.Marca, "modelo", identificacion.Modelo, "anio", identificacion.Anio, "error", err.Error())
		return componerTasacionSinReferencia(identificacion), nil
	}

	return componerTasacionConReferencia(referencia, identificacion), nil
}

// identificacionVehiculo es el resultado estructurado del modelo de visión.
type identificacionVehiculo struct {
	Marca       string `json:"marca"`
	Modelo      string `json:"modelo"`
	Anio        int    `json:"anio"`
	Estado      string `json:"estado"`
	Kilometraje int    `json:"kilometraje"`
}

// identificarVehiculo pide al modelo de visión que identifique el vehículo y
// devuelva un JSON estructurado a partir de las fotos y la descripción.
func (s *chatbotService) identificarVehiculo(ctx context.Context, descripcion string, imagenes []ImagenTasacion) (identificacionVehiculo, error) {
	texto := promptIdentificacion
	if strings.TrimSpace(descripcion) != "" {
		texto += "\n\nDescripción del cliente:\n" + strings.TrimSpace(descripcion)
	}

	partes := []llms.ContentPart{llms.TextPart(texto)}
	for _, imagen := range imagenes {
		partes = append(partes, llms.BinaryPart(imagen.MIME, imagen.Datos))
	}

	mensajes := []llms.MessageContent{
		{Role: llms.ChatMessageTypeHuman, Parts: partes},
	}

	respuesta, err := s.generar(ctx, s.modeloVision, mensajes, NumCtxVision, TimeoutVision)
	if err != nil {
		slog.Error("Error al identificar el vehículo del chatbot", "error", err.Error())
		return identificacionVehiculo{}, err
	}

	identificacion, err := parsearIdentificacion(respuesta)
	if err != nil {
		return identificacionVehiculo{}, err
	}
	if strings.TrimSpace(identificacion.Marca) == "" || strings.TrimSpace(identificacion.Modelo) == "" {
		return identificacionVehiculo{}, errors.New("el modelo no pudo identificar marca o modelo")
	}
	return identificacion, nil
}

// parsearIdentificacion extrae el JSON de la respuesta del modelo de visión.
func parsearIdentificacion(respuesta string) (identificacionVehiculo, error) {
	inicio := strings.Index(respuesta, "{")
	fin := strings.LastIndex(respuesta, "}")
	if inicio == -1 || fin <= inicio {
		return identificacionVehiculo{}, errors.New("respuesta sin JSON válido")
	}

	var identificacion identificacionVehiculo
	if err := json.Unmarshal([]byte(respuesta[inicio:fin+1]), &identificacion); err != nil {
		return identificacionVehiculo{}, err
	}
	return identificacion, nil
}

// construirContextoStock devuelve el texto con la ficha de los vehículos
// disponibles, usado como contexto del asistente.
func (s *chatbotService) construirContextoStock(ctx context.Context) (string, error) {
	vehiculos, _, err := s.repositorioVehiculos.Listar(ctx, models.EstadoDisponible, repositories.FiltrosBusqueda{}, 1, 100)
	if err != nil {
		return "", err
	}

	if len(vehiculos) == 0 {
		return "No hay vehículos disponibles en el stock actualmente.", nil
	}

	var partes []string
	for _, vehiculo := range vehiculos {
		partes = append(partes, fmt.Sprintf(
			"- %s %s %d | tipo: %s | combustible: %s | transmisión: %s | condición: %s | kilometraje: %d km | precio: $%.0f | id: %d",
			vehiculo.Marca, vehiculo.Modelo, vehiculo.Anio,
			vehiculo.Tipo, vehiculo.Combustible, vehiculo.Transmision,
			vehiculo.Condicion, vehiculo.Kilometraje, vehiculo.Precio, vehiculo.ID,
		))
	}
	return strings.Join(partes, "\n"), nil
}

// construirMensajes arma la secuencia de mensajes para el chat: prompt del
// sistema con el contexto, historial previo y el mensaje actual del usuario.
func (s *chatbotService) construirMensajes(contexto string, historial []TurnoChat, mensaje string) []llms.MessageContent {
	sistema := strings.Replace(promptSistema, "{{contexto}}", contexto, 1)

	mensajes := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, sistema),
	}

	turnos := historial
	if len(turnos) > MaximoTurnosHistorial {
		turnos = turnos[len(turnos)-MaximoTurnosHistorial:]
	}
	for _, turno := range turnos {
		contenido := strings.TrimSpace(turno.Contenido)
		if contenido == "" {
			continue
		}
		switch turno.Rol {
		case "asistente", "assistant":
			mensajes = append(mensajes, llms.TextParts(llms.ChatMessageTypeAI, contenido))
		case "usuario", "user":
			mensajes = append(mensajes, llms.TextParts(llms.ChatMessageTypeHuman, contenido))
		}
	}

	mensajes = append(mensajes, llms.TextParts(llms.ChatMessageTypeHuman, mensaje))
	return mensajes
}

// componerTasacionConReferencia arma la respuesta con el valor oficial real.
func componerTasacionConReferencia(referencia ReferenciaPrecio, identificacion identificacionVehiculo) string {
	var detalle []string
	if identificacion.Anio > 0 {
		detalle = append(detalle, fmt.Sprintf("año %d", identificacion.Anio))
	}
	if strings.TrimSpace(identificacion.Estado) != "" {
		detalle = append(detalle, fmt.Sprintf("estado %s", identificacion.Estado))
	}
	if identificacion.Kilometraje > 0 {
		detalle = append(detalle, fmt.Sprintf("%d km", identificacion.Kilometraje))
	}

	cabecera := fmt.Sprintf("Identifiqué tu vehículo como un %s %s", referencia.Marca, referencia.Modelo)
	if len(detalle) > 0 {
		cabecera += " (" + strings.Join(detalle, ", ") + ")."
	} else {
		cabecera += "."
	}

	return fmt.Sprintf(`%s

Valor de referencia según la %s:

• %s (≈ %s, a la cotización oficial vigente)
• Versión de referencia: %s

Este es el valor promedio de mercado para esa versión y año-modelo. Tu caso puntual puede variar según el estado, el equipamiento y el kilometraje. La tasación definitiva se realiza en la concesionaria.`,
		cabecera, referencia.Fuente,
		formatearUSD(referencia.PrecioUSD), formatearARS(referencia.PrecioARS),
		referencia.Version)
}

// componerTasacionSinReferencia arma la respuesta honesta cuando el vehículo
// identificado no tiene valor de referencia oficial.
func componerTasacionSinReferencia(identificacion identificacionVehiculo) string {
	var detalle []string
	if identificacion.Anio > 0 {
		detalle = append(detalle, fmt.Sprintf("año %d", identificacion.Anio))
	}
	if strings.TrimSpace(identificacion.Estado) != "" {
		detalle = append(detalle, fmt.Sprintf("estado %s", identificacion.Estado))
	}

	texto := fmt.Sprintf("Identifiqué tu vehículo como un %s %s",
		capitalizar(identificacion.Marca), capitalizar(identificacion.Modelo))
	if len(detalle) > 0 {
		texto += " (" + strings.Join(detalle, ", ") + ")."
	} else {
		texto += "."
	}

	return texto + `

No encontré un valor de referencia oficial para ese vehículo en la guía de precios actualizada, así que prefiero no inventar un número. Podés acercarte a la concesionaria o crear una consulta desde el sitio para que un vendedor te haga la tasación definitiva.`
}

// formatearARS devuelve un monto en pesos argentinos con separadores de miles.
func formatearARS(valor float64) string {
	return "$" + formatearEntero(valor)
}

// formatearUSD devuelve un monto en dólares con separadores de miles.
func formatearUSD(valor float64) string {
	return "US$ " + formatearEntero(valor)
}

// formatearEntero formatea un monto redondeado con puntos de miles.
func formatearEntero(valor float64) string {
	entero := int64(valor + 0.5)
	if entero < 0 {
		return "-" + formatearEntero(-float64(entero))
	}

	digitos := fmt.Sprintf("%d", entero)
	var partes []string
	resto := len(digitos) % 3
	if resto > 0 {
		partes = append(partes, digitos[:resto])
	}
	for i := resto; i < len(digitos); i += 3 {
		partes = append(partes, digitos[i:i+3])
	}
	return strings.Join(partes, ".")
}

// capitalizar pone la primera letra en mayúscula.
func capitalizar(texto string) string {
	if texto == "" {
		return texto
	}
	return strings.ToUpper(texto[:1]) + texto[1:]
}

// generar conecta con Ollama y devuelve la respuesta de texto del modelo.
// numCtx acota la ventana de contexto del request y keepAlive mantiene el
// modelo cargado en memoria entre consultas.
func (s *chatbotService) generar(ctx context.Context, modelo string, mensajes []llms.MessageContent, numCtx int, timeout time.Duration) (string, error) {
	modeloLocal, err := ollama.New(
		ollama.WithServerURL(s.ollamaURL),
		ollama.WithModel(modelo),
		ollama.WithRunnerNumCtx(numCtx),
		ollama.WithKeepAlive("20m"),
	)
	if err != nil {
		return "", fmt.Errorf("crear cliente de Ollama: %w", err)
	}

	ctxConTimeout, cancelar := context.WithTimeout(ctx, timeout)
	defer cancelar()

	respuesta, err := modeloLocal.GenerateContent(ctxConTimeout, mensajes, llms.WithMaxTokens(MaxTokensSalida))
	if err != nil {
		return "", err
	}
	if len(respuesta.Choices) == 0 || strings.TrimSpace(respuesta.Choices[0].Content) == "" {
		return "", errors.New("el modelo no devolvió contenido")
	}
	return strings.TrimSpace(respuesta.Choices[0].Content), nil
}

// promptSistema es el prompt base del asistente con el contexto inyectado.
const promptSistema = `Sos el asistente virtual de una concesionaria de autos. Respondés en español, de forma breve y útil.

Reglas:
- Solo podés informar sobre los vehículos del stock disponible que figuran en el contexto provisto.
- No inventes vehículos, precios ni fichas técnicas que no estén en el contexto.
- Si te preguntan por un vehículo que no está en el stock, decí que no está disponible y ofrecé orientar al usuario.
- Si no hay vehículos disponibles, decilo y orientá al usuario a volver más tarde.
- Orientá siempre al usuario a acciones concretas del sitio: crear una consulta o cotización sobre un vehículo, o solicitar un test drive.
- Comparación en vivo: si el usuario menciona el auto que tiene y el auto que quiere comprar de nuestro catálogo, ofrecé compararlos y preguntale qué aspectos quiere comparar (precio, consumo, potencia, seguridad, equipamiento, etc.). Usá la ficha real del stock para el auto de la concesionaria y tu conocimiento general para el auto del usuario. Nunca inventes datos de nuestro stock.

Contexto del stock disponible:
{{contexto}}`

// promptIdentificacion pide al modelo de visión identificar el vehículo y
// devolver un JSON estructurado. Los valores de la tasación no los genera la
// IA: se consultan en la guía oficial de precios.
const promptIdentificacion = `Sos el asistente de tasación de una concesionaria de autos.
Analizá las fotos del vehículo del cliente y tené en cuenta la descripción si la hay.
Tu única tarea es IDENTIFICAR el vehículo. No inventes precios ni valores.

Respondé únicamente con un objeto JSON válido, sin texto adicional, con este formato exacto:
{"marca": "Fiat", "modelo": "Cronos", "anio": 2020, "estado": "bueno", "kilometraje": 60000}

Reglas:
- marca y modelo: nombres comerciales en español (ej. "Fiat", "Cronos", "Ford", "Ranger").
- anio: número entero con el año-modelo; si no se puede determinar, poné 0.
- estado: uno de "excelente", "bueno", "regular" o "malo".
- kilometraje: número entero aproximado; si no se puede determinar, poné 0.
- Si no podés identificar la marca o el modelo, respondé exactamente: {"marca": "", "modelo": "", "anio": 0, "estado": "", "kilometraje": 0}`
