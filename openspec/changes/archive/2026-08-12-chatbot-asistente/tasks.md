## 1. Backend: configuración y dependencias

- [x] 1.1 Agregar `github.com/tmc/langchaingo` con `go get github.com/tmc/langchaingo` y `go mod tidy`
- [x] 1.2 Ampliar `internal/config/config.go`: campos `OllamaURL`, `ModeloChatbot` y `ModeloVision` con defaults (`http://localhost:11434`, `llama3`, `llama3.2-vision`) leídos de `OLLAMA_URL`, `MODELO_CHATBOT` y `MODELO_VISION`
- [x] 1.3 Agregar `OLLAMA_URL`, `MODELO_CHATBOT` y `MODELO_VISION` a `.env.example` con comentarios en español
- [x] 1.4 Verificar compilación: `go build ./...` y `go vet ./...`

## 2. Backend: servicio del chatbot

- [x] 2.1 Crear `internal/services/chatbot.go`: struct `ChatbotService` con el repositorio de vehículos, `OllamaURL`, `ModeloChatbot` y `ModeloVision`
- [x] 2.2 Implementar el armado de contexto: consultar stock disponible con `Listar` (estado `disponible`, página 1, tamaño 100) y serializar los vehículos a texto plano (marca, modelo, año, precio, tipo, combustible, transmisión, condición, kilometraje, id)
- [x] 2.3 Definir el `PromptSistema` en español: rol, reglas (solo stock real, no inventar, indicar si no hay stock, orientar a consulta/test drive, y detectar pedidos de comparación para ofrecer comparar preguntando los aspectos)
- [x] 2.4 Implementar `Responder(ctx, mensaje, historial) (string, error)`: valida largo máximo (1000) y cantidad de turnos (máx. 20), mapea historial a `llms.HumanMessage`/`llms.AIMessage`, llama al LLM (ollama) con timeout de 120 s
- [x] 2.5 Implementar `Tasacion(ctx, descripcion, imagenes) (string, error)`: valida cantidad (máx. 5) y peso (máx. 5 MB) de imágenes; identifica el vehículo con el modelo de visión, consulta el valor real y compone la respuesta (ver 2.6, 2.8)
- [x] 2.6 Definir `promptIdentificacion` en español: el modelo de visión identifica marca/modelo/año/estado/kilometraje y devuelve un JSON estricto; implementar `identificarVehiculo` y `parsearIdentificacion`
- [x] 2.7 Implementar degradación: si la llamada a Ollama falla o la respuesta es vacía, devolver mensaje de fallback en español (chat y tasación); si no se identifica el vehículo, responder `respuestaNoIdentificado`; loguear el error con `slog`
- [x] 2.8 Crear `internal/services/precios.go`: `ServicioPrecios` con `Buscar(ctx, marca, modelo, anio)` consultando la API de ArgAutos (`/search` → precio USD y `versión`, `/versions/{id}/valuations?currency=ars` → precio ARS), selección por año más cercano y caché en memoria de 24 h; errores tipados `ErrPrecioNoEncontrado` y `ErrPrecioNoDisponible`
- [x] 2.9 Implementar `componerTasacionConReferencia` y `componerTasacionSinReferencia` (la respuesta se arma en código con los valores reales formateados en US$/ARS; el LLM nunca genera montos) y `NuevoServicioPrecios` en el router

## 3. Backend: handler y rutas

- [x] 3.1 Crear `internal/handlers/chatbot.go`: `NuevoChatbotHandler` con `Responder` (parsea `{mensaje, historial}`, valida mensaje vacío y largo máximo → 400, delega en el service, responde `{"respuesta": "..."}`)
- [x] 3.2 En el mismo handler, implementar `Tasacion`: `c.FormFile("fotos")` (múltiples), valida extensiones (JPG/PNG/WebP), límites de cantidad/peso y `descripcion` opcional; delega en el service y responde `{"respuesta": "..."}`
- [x] 3.3 Registrar en `internal/router/router.go`: `POST /api/chatbot/mensajes` y `POST /api/chatbot/tasacion` públicos (sin middleware de autenticación); inyectar `ServicioPrecios` en `NuevoChatbotService`
- [x] 3.4 Verificar con `go build ./...`, `go vet ./...` y prueba manual con `curl` (mensaje válido, con historial, vacío, tasación con y sin fotos, y con Ollama apagado)
- [x] 3.5 Agregar `ArgAutosURL` a `internal/config/config.go` (env `ARGAUTOS_URL`, default `https://argautos.com/api/v1`) y documentarlo en `.env.example`

## 4. Frontend: servicio y tipos

- [x] 4.1 Crear `src/types/chatbot.ts`: tipos `TurnoChat`, `MensajeChatbot`, `RespuestaChatbot` y `PeticionChatbot`
- [x] 4.2 Agregar `enviarMensajeChatbot(mensaje, historial)` en `src/services/api.ts` apuntando a `POST /chatbot/mensajes`
- [x] 4.3 Agregar `enviarTasacion(fotos, descripcion)` en `src/services/api.ts` enviando multipart/form-data a `POST /chatbot/tasacion` (sin `Content-Type` manual para que el navegador setee el boundary)
- [x] 4.4 Verificar `npm run build` del frontend
- [x] 4.5 Agregar timeout de 140 s vía `AbortController` (`TIEMPO_MAXIMO_MILISEGUNDOS`) a `peticion` y `peticionMultipart` en `src/services/api.ts` para no quedar colgado en "pensando"

## 5. Frontend: widget de chat

- [x] 5.1 Crear `src/components/Chatbot.tsx`: botón flotante que abre una ventana de chat, historial de conversación en estado local (incluye respuestas del asistente para reenviar como contexto), campo de entrada y manejo de estados "enviando"/"escribiendo…"
- [x] 5.2 Enviar el historial (últimos 20 turnos) en cada request de mensaje para mantener el contexto multi-turno
- [x] 5.3 Agregar botón de adjuntar fotos (hasta 5, JPG/PNG/WebP) y campo de descripción para la tasación; mostrar el resultado de la tasación como mensaje del asistente
- [x] 5.4 Manejar errores de red/API mostrando mensaje en español y permitiendo reintentar
- [x] 5.5 Montar `<Chatbot />` en `src/layouts/LayoutBase.tsx` para que sea visible en todas las páginas públicas
- [x] 5.6 Estilos con TailwindCSS (responsive, botón fijo, burbujas de usuario/asistente, vista previa de fotos adjuntas)
- [x] 5.7 Verificar con `npm run build` y prueba manual: chat simple, comparación pidiendo aspectos y tasación con fotos
- [x] 5.8 Aceptar fotos por extensión (`EXTENSIONES_PERMITIDAS`: jpg/jpeg/png/webp) aunque el tipo MIME venga vacío; mostrar aviso "Analizando tus fotos…" y tratar `AbortError` con mensaje de error en español

## 6. Infraestructura y documentación

- [x] 6.1 Actualizar `docker-compose.yml`: agregar servicio `ollama` (imagen `ollama/ollama`) con puerto `11434`, volumen de modelos, GPU y límites de VRAM; pasar `OLLAMA_URL=http://ollama:11434`, `MODELO_CHATBOT` y `MODELO_VISION` al backend
- [ ] 6.2 Probar `docker compose up --build` y `GET /api/health`
- [x] 6.3 Documentar en `.env.example` los modelos a descargar (`ollama pull llama3` y `ollama pull minicpm-v`)
- [x] 6.4 Actualizar `docs/roadmap.md`: marcar CU-10 como Resuelto
