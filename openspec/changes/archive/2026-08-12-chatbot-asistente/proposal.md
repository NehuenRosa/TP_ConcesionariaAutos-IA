## Why

El cliente/visitante que navega el catálogo no tiene una vía rápida para hacer
preguntas en lenguaje natural sobre el stock real (precios, marcas, modelos,
disponibilidad). Implementar el asistente conversacional (CU-10) mejora la
experiencia de consulta y guía al usuario hacia acciones concretas (consultar o
pedir un test drive) sin depender de un vendedor en línea.

## What Changes

- Nuevo endpoint público `POST /api/chatbot/mensajes` que recibe el mensaje del
  usuario (con historial opcional de la conversación) y devuelve la respuesta
  del asistente en lenguaje natural.
- Nuevo servicio backend `chatbot` que orquesta LangChain (langchaingo):
  construye el contexto con el stock real (vehículos `disponibles`), lo inyecta
  en el prompt del LLM local (Ollama) y devuelve la respuesta.
- **Tasación por fotos con valores reales**: nuevo endpoint `POST
  /api/chatbot/tasacion` que recibe fotos del auto del usuario (multipart) y una
  `descripcion` opcional. El modelo con visión (Ollama) identifica
  marca/modelo/año/estado y el servicio consulta el **valor real de la Guía
  Oficial de Precios de la CCA** (vía API pública de ArgAutos) para componer la
  respuesta en US$ y ARS con la versión y el año de referencia. **Nunca se
  inventa un valor**: si no se identifica el vehículo o no hay referencia, se
  responde con honestidad (mensaje en español que orienta al usuario).
- Nuevo servicio backend `internal/services/precios.go` (`ServicioPrecios`):
  consulta la API de ArgAutos (`/search` para el valor USD y
  `/versions/{id}/valuations?currency=ars` para el valor ARS), elige la versión
  cuyo año se acerca más al del vehículo y guarda en caché (24 h, clave
  `marca|modelo|anio`) para no superar los límites de pedidos anónimos.
- **Robustez del widget de tasación** en el frontend: las fotos se aceptan por
  extensión (JPG/JPEG/PNG/WebP) aunque el navegador no informe el tipo MIME,
  el cliente HTTP centralizado usa `AbortController` con un timeout
  (`TIEMPO_MAXIMO_MILISEGUNDOS = 140000`) para no quedar colgado en "pensando",
  y se muestra el aviso "Analizando tus fotos…" mientras procesa.
- **Comparación en vivo**: el asistente detecta cuando el usuario describe su
  auto actual y el auto que quiere comprar del catálogo, y ofrece compararlos
  preguntando qué aspectos quiere comparar; la comparación usa la ficha real del
  vehículo del stock y el conocimiento del modelo para el auto del usuario.
- Nueva configuración de entorno para el LLM: `OLLAMA_URL`, `MODELO_CHATBOT` y
  `MODELO_VISION`; y `ARGAUTOS_URL` (por defecto `https://argautos.com/api/v1`)
  para la fuente de precios.
- Degradación controlada: si el modelo local no está disponible, el servicio
  responde con un mensaje en español que orienta al usuario (no falla el endpoint).
- Nuevo widget de chat flotante en el frontend (`Chatbot.tsx`) visible en las
  páginas públicas, con historial en la sesión y estados de "escribiendo".
- Actualización de `docs/roadmap.md` marcando CU-10 como resuelto.

## Capabilities

### New Capabilities

- `chatbot-asistente`: asiste al usuario con respuestas en lenguaje natural
  sobre el stock real de la concesionaria, lo orienta a consultar o solicitar
  un test drive y tasa el auto del usuario por fotos con valores reales de la
  Guía de la CCA (sin inventar precios).

### Modified Capabilities

- (ninguna)

## Impact

- **Backend**: nuevo `internal/services/chatbot.go` (lógica LangChain, chat,
  comparación y tasación por visión), nuevo `internal/services/precios.go`
  (`ServicioPrecios` con caché), nuevo `internal/handlers/chatbot.go`
  (endpoints `POST /api/chatbot/mensajes` y `POST /api/chatbot/tasacion`),
  configuración ampliada en `internal/config/config.go`, registro de rutas en
  `internal/router/router.go`. Reutiliza `VehiculoRepository.Listar` para el
  contexto del stock.
- **Dependencias**: se agrega `github.com/tmc/langchaingo` (única dependencia
  de código nueva, permitida por el stack). La fuente de precios es un servicio
  externo consultado por HTTP en tiempo de ejecución (ArgAutos → Guía CCA), no
  una dependencia de código.
- **Frontend**: nuevo componente `src/components/Chatbot.tsx` (chat + adjuntar
  fotos para tasación), tipos en `src/types/chatbot.ts` y funciones en
  `services/api.ts` (con timeout vía `AbortController`), sin nuevas dependencias.
- **Infra**: variables `OLLAMA_URL`, `MODELO_CHATBOT`, `MODELO_VISION` y
  `ARGAUTOS_URL` en `.env.example`, `docker-compose.yml` (servicio `ollama`) y
  documentación.
- **Sin cambios de base de datos** (no hay nuevas entidades).
