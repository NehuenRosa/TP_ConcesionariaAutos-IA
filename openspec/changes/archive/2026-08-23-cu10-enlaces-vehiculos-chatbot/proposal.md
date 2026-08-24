# Proposal: cu10-enlaces-vehiculos-chatbot

## Why

Hoy el asistente describe vehículos del stock en texto plano y el usuario tiene
que buscar la ficha a mano en el catálogo. Que cada mención quede acompañada de
un enlace cliqueable al detalle del vehículo acorta el camino hacia la
consulta, la cotización y la reserva, y aprovecha que el contexto ya incluye
los ids reales del stock.

## What Changes

- El prompt del sistema instruye al modelo para que, al mencionar un vehículo
  concreto del stock, agregue al final de su respuesta el marcador interno
  `[VEHICULO:<id>]` (mismo mecanismo ya probado con `[COTIZACION:<id>]`).
- `services.Responder` post-procesa la respuesta: extrae los ids únicos,
  valida que figuren en el contexto servido y los remueve del texto visible.
- `RespuestaChat` expone `VehiculosMencionados []uint` y el handler lo devuelve
  como `vehiculosMencionados` en el JSON.
- El widget (`Chatbot.tsx`) renderiza debajo del mensaje del asistente un chip
  "Ver ficha" por vehículo, que navega a `/catalogo/:id`.
- Sin marcadores no cambia nada: la respuesta se muestra igual que hoy.

## Capabilities

### New Capabilities

*(ninguna)*

### Modified Capabilities

- `chatbot-asistente`: nuevas reglas de marcadores y enlaces a fichas.

## Impact

- Backend: `services/chatbot.go` (prompt + extracción), `handlers/chatbot.go`
  (DTO), tests de servicio.
- Frontend: `types/chatbot.ts`, `components/Chatbot.tsx`.
- Sin cambios de base de datos ni de rutas.
