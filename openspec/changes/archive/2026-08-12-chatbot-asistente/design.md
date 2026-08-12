## Context

El backend ya tiene la capa handler → service → repository para vehículos, con
`VehiculoRepository.Listar` que filtra por estado y página el stock del
catálogo público. El stack exige LangChain en Go (langchaingo) y la API corre
con Gin. La motivación del cambio está en proposal.md (Why).

## Goals / Non-Goals

**Goals:**
- Exponer un endpoint público de chat que responda con base en el stock real.
- Mantener la arquitectura por capas (handler → service → repository).
- Degradar con elegancia si el modelo local no está disponible.
- Mínima superficie nueva: una dependencia (langchaingo), sin nuevas entidades.

**Non-Goals:**
- Persistir conversaciones o historial en base de datos (el historial vive en la
  sesión del widget).
- Multi-turno con memoria en el backend (endpoint sin estado).
- Autenticación/roles para el chat (es público por diseño).
- Rate limiting o control anti-abuso (se deja como mejora futura).
- WebSockets / streaming de tokens: la respuesta se devuelve completa por HTTP.

## Decisions

### 1. Proveedor LLM: Ollama local vía langchaingo

Se usa `github.com/tmc/langchaingo` con el paquete `ollama`, apuntando a
`OLLAMA_URL` (por defecto `http://localhost:11434`) y al modelo `MODELO_CHATBOT`
(por defecto `llama3`).

- **Por qué:** el stack lo define explícitamente (langchaingo) y ya hay Ollama
  instalado en el entorno de desarrollo. Evita claves externas y costos.
- **Alternativa descartada:** usar la API de OpenAI vía langchaingo (más
  capacidad, pero rompe el requisito de modelo local y agrega dependencias de
  red). Usar `llm.Register` directo sin langchaingo (evitaría la dependencia
  pero viola el stack).

### 2. Contexto del stock: reutilizar `VehiculoRepository.Listar`

El servicio de chatbot pide al repositorio los vehículos `disponibles`
(`Listar(ctx, models.EstadoDisponible, FiltrosBusqueda{}, 1, 100)`) y los
serializa a un bloque de texto plano (marca, modelo, año, precio, tipo,
combustible, transmisión, condición, kilometraje, id) que se inyecta en el
prompt del sistema.

- **Por qué:** reutiliza la capa existente sin duplicar consultas SQL ni crear
  un método nuevo específico.
- **Alternativa descartada:** crear un método `ListarParaChatbot` en el
  repository (más cerrado pero innecesario: `Listar` ya cubre el caso con
  estado + tamaño grande).

### 3. Prompt del sistema en español

Se define un `PromptSistema` constante en el service con el rol ("asistente de
la concesionaria"), las reglas de negocio (solo stock disponible, no inventar
vehículos ni precios, si no hay stock decirlo, orientar a consulta/test drive)
y el contexto de vehículos inyectado. El mensaje del usuario va como input
separado (no concatenado al prompt del sistema).

- **Por qué:** separar sistema/usuario reduce inyección de prompt y facilita
  mantener la regla de "no inventar".

### 4. Degradación controlada

Si la llamada a Ollama falla (servicio caído, timeout, error HTTP) o devuelve
respuesta vacía, el servicio retorna un mensaje fijo en español que orienta al
usuario a usar el catálogo, consultar o pedir test drive. El endpoint responde
`200` (spec) y el handler loguea el error interno con `slog`. Aplica tanto al
chat (`mensajes`) como a la tasación (`tasacion`).

- **Por qué:** un fallo del LLM no debe tumbar el endpoint ni dejar al usuario
  sin respuesta.
- **Riesgo:** enmascarar un fallo real. → se loguea siempre el error original.

### 5. Backend sin estado; el historial lo envía el frontend

`POST /api/chatbot/mensajes` recibe `{"mensaje": "...", "historial": [...]}`. El
`historial` (máx. 20 turnos) se mapea a mensajes `llms.HumanMessage`/
`llms.AIMessage` y se envía junto al prompt del sistema y el mensaje actual. No
se persiste nada en el backend.

- **Por qué:** la comparación en vivo exige varios turnos; mantener el estado en
  el frontend evita sesiones y tablas nuevas, y mantiene el backend stateless.
- **Alternativa descartada:** sesiones en el backend con token (más estado que
  mantener, sin beneficio para un asistente público informativo).

### 6. Comparación en vivo: prompt + contexto de stock + historial

El asistente compara gracias a tres ingredientes que ya se inyectan en cada
request: (a) el contexto del stock disponible (fichas reales), (b) el historial,
y (c) un prompt del sistema que instruye: cuando el usuario describa su auto y el
auto que quiere comprar de nuestro catálogo, ofrecer la comparación, preguntar
qué aspectos quiere comparar y responder solo con los datos provistos (ficha real
del stock para nuestro auto, conocimiento general para el suyo).

- **Por qué:** no requiere tool/function-calling (soporte limitado en langchaingo
  + Ollama) ni parsing de JSON estructurado; es robusto y simple.
- **Alternativa descartada:** extracción estructurada de {autoUsuario,
  autoCatalogo, aspectos} con salida JSON estricta y matcheo por marca/modelo en
  código. Más frágil ante variaciones del lenguaje natural.

### 7. Tasación por fotos: identificación con visión + valor real de la CCA

`POST /api/chatbot/tasacion` recibe multipart con hasta 5 imágenes
(JPG/PNG/WebP, máx. 5 MB cada una) y una `descripcion` opcional. El flujo tiene
tres pasos:

1. **Identificar** (`identificarVehiculo`): con `GenerateContent` +
   `llms.BinaryPart` se envía el `promptIdentificacion` al modelo de visión
   (`MODELO_VISION`), que devuelve un JSON estricto
   `{marca, modelo, anio, estado, kilometraje}` (se parsea con
   `parsearIdentificacion`). Si no se identifican marca y modelo → se responde
   con `respuestaNoIdentificado` (no se inventa nada).
2. **Consultar el valor real** (`ServicioPrecios.Buscar`): se busca la versión
   más cercana al año y se obtiene el precio en US$ y ARS de la Guía de la CCA.
3. **Componer** (`componerTasacionConReferencia` /
   `componerTasacionSinReferencia`): la respuesta se arma en código con los
   valores reales formateados (US$ y ARS con separadores de miles), la versión,
   el año de referencia y la fuente. El LLM **no genera montos**: si no hay
   referencia, la respuesta dice que no encontró valor oficial y orienta a
   acercarse a la concesionaria.

- **Por qué:** un LLM alucina precios; componer en código con un valor oficial
  (CCA) elimina la invención de montos y mantiene la transparencia de la fuente.
- **Alternativa descartada:** pedir al LLM que estime un rango de valor
  (probado y descartado: los valores eran inventados e inconsistentes). API
  externa de tasación pagada (costo y claves, fuera del stack).

### 8. Widget de chat en el layout público

`Chatbot.tsx` se monta en `LayoutBase` (visible en todas las páginas públicas).
Botón flotante → ventana con historial en `useState`, indicador "escribiendo…",
envío vía `api.enviarMensajeChatbot` y manejo de errores con `ErrorApi`. Tipos
nuevos en `src/types/chatbot.ts`.

- **Por qué:** reutiliza `LayoutBase` y el cliente HTTP centralizado
  (`services/api.ts`), sin dependencias nuevas.

### 9. Configuración

Se agregan `OLLAMA_URL`, `MODELO_CHATBOT` y `MODELO_VISION` a `Configuracion`,
`.env.example`, `docker-compose.yml` (con `OLLAMA_URL=http://ollama:11434` dentro
de la red y modelos configurados) y al Dockerfile del backend si aplica. Se usa
un timeout de contexto (p. ej. 30 s para chat y 45 s para visión) para las
llamadas al LLM.

- **Por qué:** el stack dice respetar Docker Compose; el contenedor del
  backend debe apuntar al servicio `ollama` y no a `localhost`.

### 10. Servicio de precios reales: ArgAutos + caché

`internal/services/precios.go` define `ServicioPrecios` (interfaz `Buscar(ctx,
marca, modelo, anio)`) implementado por `servicioPrecios`, que consulta la API
pública de ArgAutos (`ARGAUTOS_URL`, por defecto `https://argautos.com/api/v1`):

- `GET /search?q=<marca> <modelo>` → versión candidata (`version_id`,
  `version`, `price` USD como texto, `price_year`); se elige la versión cuyo
  año se acerca más al año buscado (`cercaniaAnio`, año 0 = 0km).
- `GET /versions/{id}/valuations?currency=ars` → precios ARS por año; se elige
  el año más cercano.
- Caché en memoria con mutex, TTL 24 h, clave `marca|modelo|anio` (minúscula),
  para no superar el límite de pedidos anónimos de la API.
- Errores tipados: `ErrPrecioNoEncontrado` (sin referencia) y
  `ErrPrecioNoDisponible` (fuente caída). Timeout HTTP de 15 s.

- **Por qué:** ArgAutos agrega la Guía Oficial de Precios de la CCA (fuente de
  referencia del mercado argentino) y es gratuita y sin clave para uso
  anónimo; el caché evita abuso. Es una consulta HTTP en runtime, no una
  dependencia de código nueva.
- **Riesgo:** la API puede no tener la versión exacta o el año pedido → se elige
  la más cercana y el valor se aclara como referencia de la guía, no como
  tasación cerrada.

## Risks / Trade-offs

- [`llama3` no es modelo agente-capable; respuestas largas o razonamiento débil] →
  El modelo es configurable (`MODELO_CHATBOT`); documentar en `.env.example` la
  opción de usar un modelo mejor (p. ej. `gemma4:26b`) cuando haga falta.
- [El modelo de visión debe estar descargado en Ollama para la tasación] →
  `MODELO_VISION` configurable; el fallback orienta al usuario; se documenta el
  comando `ollama pull` necesario en `.env.example`/README.
- [La tasación con visión puede identificar mal el vehículo] → Si la
  identificación es incierta (sin marca/modelo), se responde con honestidad
  ("prefiero no inventar un valor") y se pide completar la descripción; nunca se
  genera un monto sin referencia. Se acepta `descripcion` (año, kilometraje,
  estado) para afinar la identificación.
- [El LLM podría inventar el monto si se le pidiera tasar] → El LLM solo
  identifica el vehículo (salida JSON); el monto lo compone el código con el
  valor oficial de la CCA. Mitigación estructural, no por prompt.
- [La fuente de precios no responde o no tiene el vehículo] → Caché de 24 h,
  errores tipados, respuesta honesta (`componerTasacionSinReferencia`); se
  loguea el motivo con `slog`.
- [El historial enviado por el frontend puede agrandar el prompt] → Se limita a
  los últimos 20 turnos y el contexto de stock se mantiene compacto (ficha en
  una línea por vehículo).
- [Ollama caído o lento en dev → latencia o fallback] → Timeout de 30-45 s y
  fallback en español; el widget muestra "escribiendo…".
- [Alucinación: el LLM inventa datos fuera del contexto] → El prompt prohíbe
  inventar y solo se inyecta el stock real; el fallback cubre el caso de LLM
  inestable.
- [Inyección de prompt desde el mensaje del usuario] → El mensaje se envía como
  input separado y el prompt del sistema fija la regla de responder solo con el
  contexto provisto. Mitigación parcial; suficiente para un asistente público
  informativo.
- [Nuevo request a la DB por cada mensaje] → Aceptable para el volumen de una
  concesionaria; se puede cachear el contexto si fuera necesario (open question).

## Migration Plan

- No hay migraciones de base de datos (no se agregan entidades).
- Orden de despliegue: 1) agregar config (incluye `ARGAUTOS_URL`) y dependencia
  langchaingo, 2) service + handler + rutas, 3) frontend widget (chat + fotos),
  4) `.env.example`/docker-compose + `ollama pull` de los modelos.
- Rollback: revertir el commit del change; el resto del sistema no depende del
  chatbot (el widget se quita del layout).

## Open Questions

- ¿Cachear el contexto de vehículos si el stock cambia poco? Deferrable.
