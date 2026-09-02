package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/googleapi"

	"concesionaria/backend/internal/cifrado"
	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

// Límites de entrada del chatbot.
const (
	LargoMaximoMensaje = 1000
	// MaximoTurnosHistorial limita cuántos turnos previos (pares usuario/IA)
	// se envían al LLM como contexto. Es exclusivamente la memoria del modelo:
	// el historial completo se conserva y se muestra en la UI; el frontend y
	// los servicios de conversación solo recortan lo que entra al prompt
	// (ver docs/roadmap.md "Escalabilidad de conversaciones").
	MaximoTurnosHistorial  = 10
	MaximoImagenesTasacion = 5
	MaximoPesoImagenBytes  = 5 * 1024 * 1024
	TimeoutChatbot         = 120 * time.Second
	TimeoutVision          = 120 * time.Second
	// MaxTokensSalida es la salida máxima por modelo.
	MaxTokensSalida = 600
	// Reintentos ante errores transitorios de Gemini (503 UNAVAILABLE por picos
	// de demanda, 429 por cuota). La API de Google recomienda reintentar porque
	// los picos son temporales; la espera crece en forma exponencial (1s, 2s).
	MaximosReintentosGoogleAI = 3
	EsperaBaseReintento      = 1 * time.Second
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
	// ErrTasacionSinSesion indica que falta el identificador de la sesión de tasación.
	ErrTasacionSinSesion = errors.New("falta el identificador de la sesión de tasación")
	// ErrTasacionNoEncontrada indica que no hay una tasación pendiente para confirmar.
	ErrTasacionNoEncontrada = errors.New("no hay una tasación pendiente para esa sesión")
)

// Mensajes de degradación cuando el modelo no está disponible.
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
	// respuestaReintentarVisita se usa cuando no se pudo extraer día y franja.
	respuestaReintentarVisita = "No pude entender el día y la franja horaria. " +
		"Decime, por ejemplo: \"el jueves a las 15:00\" o \"el viernes entre las 14 y las 15\"."
	// respuestaVisitaConfirmada indica al cliente cómo presentarse en la concesionaria.
	respuestaVisitaConfirmada = "¡Listo! Coordinamos tu visita el %s, entre las %s.\n\n" +
		"Cuando te acerques a la concesionaria, decile al personal que querés terminar " +
		"de tasar tu automóvil y que ya hiciste la tasación con la IA. Tu código es: %s.\n\n" +
		"Presentalo al llegar y un supervisor va a retomar tu tasación con los valores de la guía."
)

// TurnoChat es un mensaje previo de la conversación enviado por el frontend.
type TurnoChat struct {
	Rol       string `json:"rol"`
	Contenido string `json:"contenido"`
}

// ImagenTasacion es una imagen adjunta por el usuario para la tasación.
type ImagenTasacion struct {
	MIME  string
	Datos []byte
}

// TasacionRespuesta es el resultado de una tasación: la respuesta visible y el
// identificador de sesión para confirmar la visita a la concesionaria.
type TasacionRespuesta struct {
	// Respuesta es el texto que se muestra al usuario.
	Respuesta string
	// SesionID identifica la sesión de tasación para el paso de confirmación.
	SesionID string
}

// ConfirmarTasacionResultado es el resultado de confirmar la visita.
type ConfirmarTasacionResultado struct {
	// Respuesta es el texto que se muestra al usuario.
	Respuesta string
	// Confirmada indica si la visita quedó coordinada (día, franja y código) o
	// si la IA necesitaba que el usuario reintente con otro mensaje.
	Confirmada bool
}

// RespuestaChat es la respuesta del asistente para un mensaje del chat general.
// Cuando el usuario pidió cotizar un vehículo del stock y está autenticado, el
// backend crea la cotización y expone su id para que el frontend lo redirija al
// panel de cotizaciones.
type RespuestaChat struct {
	// Respuesta es el texto mostrado al usuario (sin marcadores internos).
	Respuesta string
	// CotizacionID es el id de la cotización creada si el usuario pidió cotizar
	// un vehículo y la sesión correspondía a un cliente autenticado. Nil si no
	// se creó ninguna cotización.
	CotizacionID *uint
	// VehiculosMencionados son los ids únicos de vehículos del stock que el
	// asistente señaló con [VEHICULO:<id>] en su respuesta, ya validados contra
	// el contexto servido, para que el frontend muestre enlaces a sus fichas.
	VehiculosMencionados []uint
}

// ChatbotService define el contrato del asistente conversacional.
type ChatbotService interface {
	// Responder genera la respuesta del asistente para el mensaje del usuario,
	// usando el historial previo como contexto. clienteID es el id del cliente
	// autenticado (0 si la sesión es anónima) y se usa para crear las
	// cotizaciones que el usuario pide por chat.
	Responder(ctx context.Context, clienteID uint, mensaje string, historial []TurnoChat) (RespuestaChat, error)
	// Tasacion estima el valor de permuta del vehículo del usuario a partir de
	// las fotos y una descripción opcional, guarda la tasación pendiente y
	// devuelve la respuesta junto con la sesión para confirmar la visita.
	Tasacion(ctx context.Context, sesionID string, descripcion string, imagenes []ImagenTasacion) (TasacionRespuesta, error)
	// ConfirmarTasacion procesa el día y la franja horaria que elige el cliente,
	// genera un código único de presentación y confirma la tasación pendiente.
	ConfirmarTasacion(ctx context.Context, sesionID string, mensaje string) (ConfirmarTasacionResultado, error)
	// GenerarCotizacion genera la respuesta del asistente dentro del chat de
	// una cotización, acotado a la ficha real del vehículo cotizado.
	GenerarCotizacion(ctx context.Context, vehiculo models.Vehiculo, historial []TurnoChat, mensaje string) (string, error)
}

// Proveedor de LLM soportado.
const (
	// ProveedorGoogleAI usa Gemini en la nube (no consume recursos locales
	// y ofrece contexto de 1M de tokens).
	ProveedorGoogleAI = "googleai"
)

// chatbotService implementa ChatbotService con LangChain y un LLM provisto por
// Google AI (Gemini, en la nube).
type chatbotService struct {
	repositorioVehiculos    repositories.VehiculoRepository
	repositorioTasaciones   repositories.TasacionRepository
	repositorioCotizaciones repositories.CotizacionRepository
	cifrador                cifrado.Cifrador
	proveedor               string
	googleAIKey             string
	modeloChatbot           string
	modeloVision            string
	precios                 ServicioPrecios
}

// NuevoChatbotService crea el servicio del asistente conversacional.
func NuevoChatbotService(repositorioVehiculos repositories.VehiculoRepository, repositorioTasaciones repositories.TasacionRepository, repositorioCotizaciones repositories.CotizacionRepository, cifrador cifrado.Cifrador, proveedor string, googleAIKey string, modeloChatbot string, modeloVision string, precios ServicioPrecios) ChatbotService {
	if proveedor == "" {
		proveedor = ProveedorGoogleAI
	}
	return &chatbotService{
		repositorioVehiculos:    repositorioVehiculos,
		repositorioTasaciones:   repositorioTasaciones,
		repositorioCotizaciones: repositorioCotizaciones,
		cifrador:                cifrador,
		proveedor:               proveedor,
		googleAIKey:             googleAIKey,
		modeloChatbot:           modeloChatbot,
		modeloVision:            modeloVision,
		precios:                 precios,
	}
}

// Responder valida el mensaje, arma el contexto del stock y la conversación, y
// delega la generación de la respuesta en el modelo local. Si el modelo no
// está disponible, devuelve un mensaje de fallback en español. Cuando la IA
// marca una cotización y la sesión es de un cliente autenticado, crea el
// registro con los mensajes cifrados y lo devuelve en CotizacionID.
func (s *chatbotService) Responder(ctx context.Context, clienteID uint, mensaje string, historial []TurnoChat) (RespuestaChat, error) {
	if strings.TrimSpace(mensaje) == "" {
		return RespuestaChat{}, ErrMensajeVacio
	}
	if len([]rune(mensaje)) > LargoMaximoMensaje {
		return RespuestaChat{}, ErrMensajeMuyLargo
	}

	contexto, idsServidos, err := s.construirContextoStock(ctx)
	if err != nil {
		return RespuestaChat{}, fmt.Errorf("obtener stock para el chatbot: %w", err)
	}

	mensajes := s.construirMensajes(contexto, historial, mensaje)
	respuesta, err := s.generar(ctx, s.modeloChatbot, mensajes, TimeoutChatbot)
	if err != nil {
		slog.Error("Error al generar respuesta del chatbot", "error", err.Error())
		return RespuestaChat{Respuesta: respuestaFallbackChat}, nil
	}

	respuestaEditable, vehiculoID := limpiarMarcadorCotizacion(respuesta)
	respuestaEditable = normalizarRespuestaConversacional(respuestaEditable)
	mencionados, respuestaEditable := extraerMarcadoresVehiculo(respuestaEditable)
	mencionados = filtrarIdsServidos(mencionados, idsServidos)

	var cotizacionID *uint
	if vehiculoID > 0 && clienteID > 0 {
		cotizacion, err := s.crearCotizacionDesdeChat(ctx, vehiculoID, clienteID, mensaje, respuestaEditable)
		if err != nil {
			slog.Error("No se pudo crear la cotización desde el chat", "error", err.Error())
		} else {
			cotizacionID = &cotizacion.ID
		}
	}

	return RespuestaChat{
		Respuesta:            respuestaEditable,
		CotizacionID:         cotizacionID,
		VehiculosMencionados: mencionados,
	}, nil
}

// Tasacion identifica el vehículo con el modelo de visión, consulta el valor
// de referencia oficial en la guía de precios y compone la respuesta con datos
// reales, sin inventar valores. Cuando hay referencia, guarda la tasación como
// pendiente para poder confirmar después el día y la franja de la visita.
func (s *chatbotService) Tasacion(ctx context.Context, sesionID string, descripcion string, imagenes []ImagenTasacion) (TasacionRespuesta, error) {
	if len(imagenes) > MaximoImagenesTasacion {
		return TasacionRespuesta{}, ErrDemasiadasImagenes
	}
	for _, imagen := range imagenes {
		if len(imagen.Datos) > MaximoPesoImagenBytes {
			return TasacionRespuesta{}, ErrImagenMuyPesada
		}
	}
	if len(imagenes) == 0 && strings.TrimSpace(descripcion) == "" {
		return TasacionRespuesta{}, ErrSinImagenes
	}

	identificacion, err := s.identificarVehiculo(ctx, descripcion, imagenes)
	if err != nil {
		slog.Warn("No se pudo identificar el vehículo para tasar", "error", err.Error())
		return TasacionRespuesta{Respuesta: respuestaNoIdentificado}, nil
	}

	referencia, err := s.precios.Buscar(ctx, identificacion.Marca, identificacion.Modelo, identificacion.Anio)
	if err != nil {
		slog.Warn("No se encontró valor de referencia para la tasación",
			"marca", identificacion.Marca, "modelo", identificacion.Modelo, "anio", identificacion.Anio, "error", err.Error())
		return TasacionRespuesta{Respuesta: componerTasacionSinReferencia(identificacion)}, nil
	}

	if strings.TrimSpace(sesionID) == "" {
		sesionID, err = nuevaSesionID()
		if err != nil {
			return TasacionRespuesta{}, fmt.Errorf("generar sesión de tasación: %w", err)
		}
	}

	tasacion := &models.Tasacion{
		SesionID:    sesionID,
		Marca:       referencia.Marca,
		Modelo:      referencia.Modelo,
		Version:     referencia.Version,
		Anio:        identificacion.Anio,
		Estado:      identificacion.Estado,
		PrecioUSD:   referencia.PrecioUSD,
		PrecioARS:   referencia.PrecioARS,
		EstadoFlujo: models.EstadoTasacionPendiente,
	}
	if err := s.repositorioTasaciones.Crear(ctx, tasacion); err != nil {
		slog.Error("No se pudo guardar la tasación pendiente", "error", err.Error())
	}

	return TasacionRespuesta{
		Respuesta: componerTasacionConReferencia(referencia, identificacion),
		SesionID:  sesionID,
	}, nil
}

// ConfirmarTasacion extrae el día y la franja horaria que elige el cliente de
// su mensaje, genera un código único de presentación, confirma la tasación
// pendiente y responde indicándole cómo presentarse en la concesionaria.
func (s *chatbotService) ConfirmarTasacion(ctx context.Context, sesionID string, mensaje string) (ConfirmarTasacionResultado, error) {
	if strings.TrimSpace(sesionID) == "" {
		return ConfirmarTasacionResultado{}, ErrTasacionSinSesion
	}
	if strings.TrimSpace(mensaje) == "" {
		return ConfirmarTasacionResultado{}, ErrMensajeVacio
	}
	if len([]rune(mensaje)) > LargoMaximoMensaje {
		return ConfirmarTasacionResultado{}, ErrMensajeMuyLargo
	}

	tasacion, err := s.repositorioTasaciones.ObtenerPendientePorSesion(ctx, sesionID)
	if err != nil {
		if errors.Is(err, repositories.ErrTasacionNoEncontrada) {
			return ConfirmarTasacionResultado{}, ErrTasacionNoEncontrada
		}
		slog.Error("Error al buscar la tasación pendiente", "sesionID", sesionID, "error", err.Error())
		return ConfirmarTasacionResultado{}, err
	}

	visita, err := s.extraerVisita(ctx, mensaje)
	if err != nil {
		slog.Warn("No se pudo extraer día y franja de la confirmación", "error", err.Error())
		return ConfirmarTasacionResultado{Respuesta: respuestaReintentarVisita, Confirmada: false}, nil
	}
	if !models.FranjaValida(visita.Franja) {
		return ConfirmarTasacionResultado{Respuesta: respuestaReintentarVisita, Confirmada: false}, nil
	}

	codigo, err := s.generarCodigoUnico(ctx)
	if err != nil {
		slog.Error("No se pudo generar un código de tasación único", "error", err.Error())
		return ConfirmarTasacionResultado{Respuesta: respuestaReintentarVisita, Confirmada: false}, nil
	}

	tasacion.Dia = visita.Dia
	tasacion.Franja = visita.Franja
	tasacion.Codigo = &codigo
	tasacion.EstadoFlujo = models.EstadoTasacionConfirmada
	if err := s.repositorioTasaciones.Actualizar(ctx, tasacion); err != nil {
		slog.Error("No se pudo confirmar la tasación", "sesionID", sesionID, "error", err.Error())
		return ConfirmarTasacionResultado{}, err
	}

	return ConfirmarTasacionResultado{
		Respuesta:  fmt.Sprintf(respuestaVisitaConfirmada, visita.Dia, visita.Franja, codigo),
		Confirmada: true,
	}, nil
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

	respuesta, err := s.generar(ctx, s.modeloVision, mensajes, TimeoutVision)
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

// nuevaSesionID genera un identificador aleatorio para una sesión de tasación.
func nuevaSesionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// visitaTasacion es el día y la franja horaria que elige el cliente.
type visitaTasacion struct {
	Dia    string `json:"dia"`
	Franja string `json:"franja"`
}

// extraerVisita pide al modelo que interprete el día y la franja horaria del
// mensaje del cliente y los devuelva estructurados.
func (s *chatbotService) extraerVisita(ctx context.Context, mensaje string) (visitaTasacion, error) {
	mensajes := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextPart(promptExtraerVisita),
				llms.TextPart("Mensaje del cliente:\n" + strings.TrimSpace(mensaje)),
			},
		},
	}

	respuesta, err := s.generar(ctx, s.modeloChatbot, mensajes, TimeoutChatbot)
	if err != nil {
		slog.Error("Error al extraer día y franja del chatbot", "error", err.Error())
		return visitaTasacion{}, err
	}

	var visita visitaTasacion
	inicio := strings.Index(respuesta, "{")
	fin := strings.LastIndex(respuesta, "}")
	if inicio == -1 || fin <= inicio {
		return visitaTasacion{}, errors.New("respuesta sin JSON válido")
	}
	if err := json.Unmarshal([]byte(respuesta[inicio:fin+1]), &visita); err != nil {
		return visitaTasacion{}, err
	}
	visita.Dia = strings.TrimSpace(visita.Dia)
	visita.Franja = strings.TrimSpace(visita.Franja)
	return visita, nil
}

// generarCodigoUnico genera un código de presentación aleatorio que no se
// repita en la base de datos.
func (s *chatbotService) generarCodigoUnico(ctx context.Context) (string, error) {
	const caracteres = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	for intento := 0; intento < 10; intento++ {
		var b strings.Builder
		for i := 0; i < 6; i++ {
			indice, err := rand.Int(rand.Reader, big.NewInt(int64(len(caracteres))))
			if err != nil {
				return "", err
			}
			b.WriteByte(caracteres[indice.Int64()])
		}
		codigo := b.String()
		existe, err := s.repositorioTasaciones.CodigoExiste(ctx, codigo)
		if err != nil {
			return "", err
		}
		if !existe {
			return codigo, nil
		}
	}
	return "", errors.New("no se pudo generar un código único")
}

// construirContextoStock devuelve el texto con la ficha de los vehículos
// disponibles (usado como contexto del asistente) y un set con los ids que
// figuran en ese contexto, para validar los marcadores de vehículos.
func (s *chatbotService) construirContextoStock(ctx context.Context) (string, map[uint]struct{}, error) {
	vehiculos, _, err := s.repositorioVehiculos.Listar(ctx, models.EstadoDisponible, repositories.FiltrosBusqueda{}, 1, 100)
	if err != nil {
		return "", nil, err
	}

	if len(vehiculos) == 0 {
		return "No hay vehículos disponibles en el stock actualmente.", map[uint]struct{}{}, nil
	}

	idsServidos := make(map[uint]struct{}, len(vehiculos))
	var partes []string
	for _, vehiculo := range vehiculos {
		idsServidos[vehiculo.ID] = struct{}{}
		partes = append(partes, fmt.Sprintf(
			"- [id:%d] %s %s %d | tipo: %s | combustible: %s | transmisión: %s | condición: %s | kilometraje: %d km | precio: $%.0f",
			vehiculo.ID, vehiculo.Marca, vehiculo.Modelo, vehiculo.Anio,
			vehiculo.Tipo, vehiculo.Combustible, vehiculo.Transmision,
			vehiculo.Condicion, vehiculo.Kilometraje, vehiculo.Precio,
		))
	}
	return strings.Join(partes, "\n"), idsServidos, nil
}

// construirMensajes arma la secuencia de mensajes para el chat general: prompt
// del sistema con el contexto del stock, historial previo y el mensaje actual.
func (s *chatbotService) construirMensajes(contexto string, historial []TurnoChat, mensaje string) []llms.MessageContent {
	sistema := strings.Replace(promptSistema, "{{contexto}}", contexto, 1)
	return construirMensajesDesde(sistema, historial, mensaje)
}

// construirMensajesDesde arma la secuencia de mensajes para un prompt de
// sistema dado, con el historial previo y el mensaje actual del usuario.
func construirMensajesDesde(sistema string, historial []TurnoChat, mensaje string) []llms.MessageContent {
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

// regexMarcadorCotizacion reconoce el marcador interno con el que la IA señala
// que el usuario pidió cotizar un vehículo del stock, por ejemplo
// "[COTIZACION:7]".
var regexMarcadorCotizacion = regexp.MustCompile(`\[COTIZACION:(\d+)\]`)

// limpiarMarcadorCotizacion quita el marcador de la respuesta del modelo y
// devuelve el id del vehículo a cotizar (0 si no venía marcador).
func limpiarMarcadorCotizacion(respuesta string) (string, uint) {
	coincidencia := regexMarcadorCotizacion.FindStringSubmatch(respuesta)
	if len(coincidencia) < 2 {
		return strings.TrimSpace(respuesta), 0
	}

	var id uint64
	if _, err := fmt.Sscanf(coincidencia[1], "%d", &id); err != nil {
		return strings.TrimSpace(respuesta), 0
	}

	limpia := strings.TrimSpace(regexMarcadorCotizacion.ReplaceAllString(respuesta, ""))
	return limpia, uint(id)
}

// regexMarcadorVehiculo reconoce el marcador interno con el que la IA señala
// un vehículo puntual del stock que mencionó en su respuesta, por ejemplo
// "[VEHICULO:7]".
var regexMarcadorVehiculo = regexp.MustCompile(`\[VEHICULO:(\d+)\]`)

// regexEspaciosDobles colapsa espacios repetidos en una misma línea sin tocar
// los saltos de línea ni los párrafos.
var regexEspaciosDobles = regexp.MustCompile(`[^\S\n]{2,}`)

// regexEspacioAntesPuntuacion quita el espacio que queda antes de un signo de
// puntuación al eliminar un marcador intermedio.
var regexEspacioAntesPuntuacion = regexp.MustCompile(`[^\S\n]+([.,;:!?])`)

// MaximoVehiculosMencionados limita los enlaces por respuesta para no saturar
// la interfaz si el modelo repite el marcador muchas veces.
const MaximoVehiculosMencionados = 5

// extraerMarcadoresVehiculo quita los marcadores de vehículos del texto y
// devuelve los ids únicos en orden de aparición (hasta 5) más la respuesta
// limpia lista para mostrarse al usuario.
func extraerMarcadoresVehiculo(respuesta string) ([]uint, string) {
	coincidencias := regexMarcadorVehiculo.FindAllStringSubmatch(respuesta, -1)
	textoLimpio := regexEspacioAntesPuntuacion.ReplaceAllString(
		regexEspaciosDobles.ReplaceAllString(regexMarcadorVehiculo.ReplaceAllString(respuesta, ""), " "),
		"$1",
	)
	textoLimpio = strings.TrimSpace(textoLimpio)
	if len(coincidencias) == 0 {
		return nil, textoLimpio
	}

	vistos := make(map[uint]struct{}, len(coincidencias))
	ids := make([]uint, 0, len(coincidencias))
	for _, coincidencia := range coincidencias {
		id64, err := strconv.ParseUint(coincidencia[1], 10, 32)
		if err != nil || id64 == 0 {
			continue
		}
		id := uint(id64)
		if _, repetido := vistos[id]; repetido {
			continue
		}
		vistos[id] = struct{}{}
		if len(ids) < MaximoVehiculosMencionados {
			ids = append(ids, id)
		}
	}
	return ids, textoLimpio
}

// filtrarIdsServidos descarta los ids que no figuraban en el contexto del
// stock servido al modelo, para nunca enlazar un vehículo inexistente, dado
// de baja o reservado.
func filtrarIdsServidos(ids []uint, servidos map[uint]struct{}) []uint {
	if len(ids) == 0 || len(servidos) == 0 {
		return nil
	}
	filtrados := make([]uint, 0, len(ids))
	for _, id := range ids {
		if _, ok := servidos[id]; ok {
			filtrados = append(filtrados, id)
		}
	}
	return filtrados
}

// normalizarRespuestaConversacional quita el formato Markdown que algunos
// modelos agregan aunque se les pida responder como una conversación.
func normalizarRespuestaConversacional(respuesta string) string {
	respuesta = strings.TrimSpace(respuesta)
	respuesta = regexp.MustCompile("(?m)^```(?:markdown|md|texto|text)?\\s*$|^```\\s*$|^#{1,6}\\s+").ReplaceAllString(respuesta, "")
	respuesta = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`).ReplaceAllString(respuesta, "$1")
	respuesta = regexp.MustCompile(`(?m)^\s*[-*+]\s+`).ReplaceAllString(respuesta, "")
	respuesta = regexp.MustCompile(`\*\*([^*\n]+)\*\*|__([^_\n]+)__`).ReplaceAllString(respuesta, "$1$2")
	respuesta = regexp.MustCompile(`\*([^*\n]+)\*|_([^_\n]+)_`).ReplaceAllString(respuesta, "$1$2")
	return strings.TrimSpace(respuesta)
}

// crearCotizacionDesdeChat crea una cotización desde el chat general con el
// mensaje del cliente y la respuesta del asistente, ambos cifrados. Valida que
// el vehículo exista y siga disponible en el stock.
func (s *chatbotService) crearCotizacionDesdeChat(ctx context.Context, vehiculoID uint, clienteID uint, mensajeCliente string, respuestaIA string) (*models.Cotizacion, error) {
	vehiculo, err := s.repositorioVehiculos.ObtenerPorID(ctx, vehiculoID)
	if err != nil {
		return nil, fmt.Errorf("obtener vehículo a cotizar: %w", err)
	}
	if vehiculo.Estado != models.EstadoDisponible {
		return nil, ErrVehiculoNoDisponible
	}

	cotizacion, err := construirCotizacion(s.cifrador, vehiculoID, clienteID, mensajeCliente, respuestaIA)
	if err != nil {
		return nil, err
	}
	return s.repositorioCotizaciones.Crear(ctx, cotizacion)
}

// GenerarCotizacion genera la respuesta del asistente dentro de la
// conversación de una cotización, acotada a la ficha real del vehículo
// cotizado y el historial descifrado. Si el modelo está caído, devuelve el
// fallback que orienta al cliente a seguir en el sitio.
func (s *chatbotService) GenerarCotizacion(ctx context.Context, vehiculo models.Vehiculo, historial []TurnoChat, mensaje string) (string, error) {
	ficha := fichaVehiculo(vehiculo)
	sistema := strings.Replace(promptCotizacion, "{{ficha}}", ficha, 1)
	sistema = strings.Replace(sistema, "{{condiciones}}", textoCondicionesComerciales(vehiculo.Precio), 1)

	respuesta, err := s.generar(ctx, s.modeloChatbot, construirMensajesDesde(sistema, historial, mensaje), TimeoutChatbot)
	if err != nil {
		slog.Error("Error al generar respuesta de cotización", "error", err.Error())
		return respuestaFallbackChat, nil
	}
	return normalizarRespuestaConversacional(respuesta), nil
}

// fichaVehiculo devuelve el texto de la ficha real de un vehículo del stock.
func fichaVehiculo(vehiculo models.Vehiculo) string {
	return fmt.Sprintf(
		"%s %s %d | tipo: %s | combustible: %s | transmisión: %s | condición: %s | kilometraje: %d km | precio: $%.0f",
		vehiculo.Marca, vehiculo.Modelo, vehiculo.Anio,
		vehiculo.Tipo, vehiculo.Combustible, vehiculo.Transmision,
		vehiculo.Condicion, vehiculo.Kilometraje, vehiculo.Precio,
	)
}

// Porcentajes de las opciones de cobro de la tasación, siempre sobre el valor
// base de la guía. Se componen en código: la IA nunca genera montos.
const (
	recargoParteDePago = 0.03
	recargoCobro30Dias = 0.02
)

// opcionTasacion describe una forma de cobrar el auto del cliente.
type opcionTasacion struct {
	titulo  string
	nota    string
	recargo float64
}

// opcionesDeCobro son las tres alternativas que se ofrecen sobre el valor base.
var opcionesDeCobro = []opcionTasacion{
	{titulo: "Venta normal", nota: "", recargo: 0},
	{titulo: "Como parte de pago para comprar otro auto", nota: "+3% sobre el valor base", recargo: recargoParteDePago},
	{titulo: "Cobrar el dinero a 30 días", nota: "+2% sobre el valor base", recargo: recargoCobro30Dias},
}

// componerTasacionConReferencia arma la respuesta con el valor oficial real y
// las tres opciones de cobro calculadas en código.
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

	lineas := make([]string, 0, len(opcionesDeCobro))
	for i, opcion := range opcionesDeCobro {
		udis := formatearUSD(referencia.PrecioUSD * (1 + opcion.recargo))
		ars := formatearARS(referencia.PrecioARS * (1 + opcion.recargo))
		linea := fmt.Sprintf("%d. %s: %s (≈ %s)", i+1, opcion.titulo, udis, ars)
		if opcion.nota != "" {
			linea += " (" + opcion.nota + ")"
		}
		lineas = append(lineas, linea)
	}

	return fmt.Sprintf(`%s

Valor de referencia según la %s:

• %s (≈ %s, a la cotización oficial vigente)
• Versión de referencia: %s

Además, te ofrezco tres opciones para cobrar tu vehículo:

%s

Estos montos son de referencia y pueden variar según el estado, el equipamiento y el kilometraje de tu unidad. El precio real te lo va a confirmar uno de nuestros supervisores al ver el auto en persona.

Para coordinar tu visita, decime: ¿qué día y entre qué franja horaria te podés acercar a la concesionaria?`,
		cabecera, referencia.Fuente,
		formatearUSD(referencia.PrecioUSD), formatearARS(referencia.PrecioARS),
		referencia.Version,
		strings.Join(lineas, "\n"))
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

// generar delega la generación al proveedor de LLM configurado.
func (s *chatbotService) generar(ctx context.Context, modelo string, mensajes []llms.MessageContent, timeout time.Duration) (string, error) {
	return s.generarConGoogleAI(ctx, modelo, mensajes, timeout)
}

// generarConGoogleAI genera con Gemini en la nube. El mismo modelo soporta
// texto y visión, por eso el mismo camino sirve para el chat y la tasación.
// Replica los errores transitorios (503 por picos de demanda, 429 de cuota y
// otros 5xx) con espera exponencial, porque suelen ser pasajeros.
func (s *chatbotService) generarConGoogleAI(ctx context.Context, modelo string, mensajes []llms.MessageContent, timeout time.Duration) (string, error) {
	modeloNube, err := googleai.New(ctx, googleai.WithAPIKey(s.googleAIKey))
	if err != nil {
		return "", fmt.Errorf("crear cliente de Google AI: %w", err)
	}
	defer modeloNube.Close()

	ctxConTimeout, cancelar := context.WithTimeout(ctx, timeout)
	defer cancelar()

	var respuesta *llms.ContentResponse
	for intento := 0; intento <= MaximosReintentosGoogleAI; intento++ {
		if intento > 0 {
			espera := EsperaBaseReintento << (intento - 1)
			select {
			case <-ctxConTimeout.Done():
				return "", err
			case <-time.After(espera):
			}
		}

		respuesta, err = modeloNube.GenerateContent(ctxConTimeout, mensajes,
			llms.WithModel(modelo),
			llms.WithMaxTokens(MaxTokensSalida),
		)
		if err == nil {
			break
		}
		if !esErrorTransitorioGoogleAI(err) || intento == MaximosReintentosGoogleAI {
			return "", err
		}
		slog.Info("Gemini transitoriamente no disponible, reintentando",
			"modelo", modelo, "intento", intento+1, "error", err.Error())
	}
	if len(respuesta.Choices) == 0 || strings.TrimSpace(respuesta.Choices[0].Content) == "" {
		return "", errors.New("el modelo no devolvió contenido")
	}
	return strings.TrimSpace(respuesta.Choices[0].Content), nil
}

// esErrorTransitorioGoogleAI indica si el error de Gemini responde a una
// condición temporal (503 UNAVAILABLE por alta demanda, 429 por cuota o
// cualquier 5xx) que justifica reintentar. Reconoce el error tipado de la API
// de Google, el mapeado por langchaingo y el texto por las dudas.
func esErrorTransitorioGoogleAI(err error) bool {
	var errAPI *googleapi.Error
	if errors.As(err, &errAPI) {
		return errAPI.Code == http.StatusTooManyRequests || errAPI.Code >= http.StatusInternalServerError
	}
	var errLangchaingo *llms.Error
	if errors.As(err, &errLangchaingo) {
		return errLangchaingo.Code == llms.ErrCodeRateLimit || errLangchaingo.Code == llms.ErrCodeProviderUnavailable
	}
	texto := strings.ToLower(err.Error())
	return strings.Contains(texto, "503") ||
		strings.Contains(texto, "429") ||
		strings.Contains(texto, "service unavailable") ||
		strings.Contains(texto, "unavailable") ||
		strings.Contains(texto, "high demand") ||
		strings.Contains(texto, "quota") ||
		strings.Contains(texto, "rate limit")
}

// promptSistema es el prompt base del asistente con el contexto inyectado.
const promptSistema = `Sos el asistente virtual de una concesionaria de autos. Respondés en español, de forma breve, cálida y útil.

Estilo de respuesta:
- Escribí como una persona que atiende por chat, con frases naturales y párrafos cortos.
- Usá texto plano: no uses Markdown, títulos, negritas, cursivas, bloques de código ni listas con viñetas o numeradas.
- No empieces cada frase en una línea nueva; agrupá las ideas en uno o dos párrafos cuando corresponda.
- Mencioná los datos importantes dentro de la oración, sin formato de ficha técnica.

Reglas:
- Solo podés informar sobre los vehículos del stock disponible que figuran en el contexto provisto.
- No inventes vehículos, precios ni fichas técnicas que no estén en el contexto.
- Si te preguntan por un vehículo que no está en el stock, decí que no está disponible y ofrecé orientar al usuario.
- Si no hay vehículos disponibles, decilo y orientá al usuario a volver más tarde.
- Orientá siempre al usuario a acciones concretas del sitio: crear una consulta o cotización sobre un vehículo, o solicitar un test drive.
- Comparación en vivo: si el usuario menciona el auto que tiene y el auto que quiere comprar de nuestro catálogo, ofrecé compararlos y preguntale qué aspectos quiere comparar (precio, consumo, potencia, seguridad, equipamiento, etc.). Usá la ficha real del stock para el auto de la concesionaria y tu conocimiento general para el auto del usuario. Nunca inventes datos de nuestro stock.
- Cotización: usá el marcador [COTIZACION:<id>] SOLO si el usuario pide explícitamente cotizar, presupuestar o conocer cómo pagar/financiar un vehículo del stock. No lo uses cuando solo pregunta características, disponibilidad, compara modelos o quiere un test drive: en esos casos respondé sin marcador. Si el vehículo mencionado no está en el stock, no pongas el marcador. Cuando pongas el marcador, finalizá tu respuesta con la línea exacta [COTIZACION:<id>] usando el id que figura entre corchetes en el contexto (ej. [id:3]) y que sea lo último que escribas, sin texto después.
- Enlaces a fichas: cada vez que tu respuesta mencione uno o más vehículos puntuales del stock, agregá al final, después de todo el texto, el marcador [VEHICULO:<id>] por cada vehículo mencionado, usando el id que figura entre corchetes en el contexto (ej. [id:3] → [VEHICULO:3]). Solo vehículos que figuren en el contexto; nunca inventes ids ni uses el marcador si hablás en general sin nombrar un vehículo concreto.

Contexto del stock disponible:
{{contexto}}`

// promptCotizacion es el prompt del chat dentro de una cotización: la IA se
// centra en un único vehículo real del stock y en las condiciones comerciales
// oficiales de la concesionaria. Nunca inventa precios, tasas ni montos.
const promptCotizacion = `Sos el asistente de cotizaciones de una concesionaria de autos. Estás atendiendo una conversación sobre un vehículo específico del stock. Respondés en español, de forma breve, cálida y útil.

Estilo de respuesta:
- Escribí como una persona que atiende por chat, con frases naturales y párrafos cortos.
- Usá texto plano: no uses Markdown, títulos, negritas, cursivas, bloques de código ni listas con viñetas o numeradas.
- Presentá los datos comerciales dentro de la conversación, sin formato de ficha técnica.

Vehículo a cotizar (ficha real del stock):
{{ficha}}

Condiciones comerciales de la concesionaria (datos oficiales, única fuente de verdad para pagos y financiación):
{{condiciones}}

Reglas:
- Hablá solo de este vehículo y de los datos reales de su ficha.
- Para hablar de métodos de pago, planes, cuotas, tasas o descuentos usá EXCLUSIVAMENTE las condiciones comerciales provistas arriba. Nunca inventes tasas, plazos, cuotas, descuentos ni montos que no figuren ahí.
- Cuando el cliente pregunte cómo pagar o financiar, presentale los métodos de pago y proponé el plan de financiación que más le convenga según lo que pida (adelanto, cantidad de cuotas), explicando la cuota estimada calculada.
- Si el cliente no define su elección, orientalo a reservar el vehículo, solicitar un test drive o dejar una consulta para que un vendedor confirme las condiciones finales.`

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

// promptExtraerVisita pide al modelo extraer el día y la franja horaria de la
// visita del cliente en forma estructurada. La disponibilidad la valida el
// código contra el catálogo de franjas: el modelo solo interpreta el mensaje.
const promptExtraerVisita = `Sos el asistente de tasación de una concesionaria de autos.
El cliente respondió en qué día y franja horaria puede acercarse a la concesionaria.
Tu única tarea es interpretar ese mensaje y devolver la fecha y la franja horaria.

Respondé únicamente con un objeto JSON válido, sin texto adicional, con este formato exacto:
{"dia": "el jueves", "franja": "15:00"}

Reglas:
- dia: el día que menciona el cliente en su idioma natural (ej. "el jueves", "mañana", "el martes 10").
- franja: la hora de inicio de la franja en formato "HH:MM" de 24 h (ej. "09:00", "10:00", "11:00", "14:00", "15:00", "16:00", "17:00"). Solo esas horas; si menciona otra hora, redondeá a la más cercana.
- Si no menciona día o franja, poné el campo como string vacío: {"dia": "", "franja": ""}`
