# Design: cu10-enlaces-vehiculos-chatbot

## Contexto

El chat general ya resuelve un caso similar: cuando el usuario pide cotizar, el
modelo agrega `[COTIZACION:<id>]` y `limpiarMarcadorCotizacion` lo extrae y
limpia. Reutilizamos exactamente ese patrón para los enlaces, en vez de pedir
Markdown o URLs (el normalizador conversacional borra enlaces y formato a
propósito).

## Decisiones

- **D1 — Marcador textual, no Markdown.** El modelo appendea
  `[VEHICULO:<id>]` al final de la respuesta. Es robusto con ambos proveedores
  (Gemini y Ollama) y sobrevive a `normalizarRespuestaConversacional`, que hoy
  elimina `[texto](url)` y negritas.
- **D2 — Validación contra el contexto servido.** `construirContextoStock`
  ya lista hasta 100 vehículos disponibles con su id. `Responder` colecciona
  esos ids en un set antes de generar; tras la respuesta descarta marcadores
  con ids que no estén en el set. Así el modelo no puede inyectar enlaces a
  unidades inexistentes o no disponibles.
- **D3 — Ids únicos, ordenados por aparición.** Se deduplican conservando el
  orden de aparición y se limita a un máximo razonable (5) para no saturar la
  UI si el modelo se descontrola.
- **D4 — DTO plano.** `RespuestaChat.VehiculosMencionados []uint`; el handler
  serializa `vehiculosMencionados`. `omitempty`: sin menciones el JSON queda
  igual que hoy.
- **D5 — Chips en el widget.** Bajo cada mensaje del asistente con vehículos,
  chips clickeables ("Ver ficha") con `useNavigate` a `/catalogo/:id`. No se
  navega automáticamente: el usuario elige. El historial enviado de vuelta al
  backend usa el texto limpio (sin marcadores), evitando contaminar turnos.

## Riesgos

- El modelo puede olvidar el marcador → degradación suave: simplemente no hay
  chip (estado actual). No rompe nada.
- El modelo puede poner el marcador en el medio del texto → la extracción por
  regex lo quita de donde esté; el texto visible queda prolijo.

## Migración

Sin base de datos. Despliegue atómico backend+frontend; versiones viejas del
widget siguen funcionando (ignoran el campo nuevo).
