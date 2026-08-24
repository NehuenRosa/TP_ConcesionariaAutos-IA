# Tasks: cu10-enlaces-vehiculos-chatbot

## 1. Backend — servicio

- [x] 1.1 Agregar a `promptSistema` la regla del marcador `[VEHICULO:<id>]`: mencionar vehículos concretos del stock agregando el marcador al final, sin inventar ids fuera del contexto y sin usarlo cuando no se menciona ningún vehículo puntual
- [x] 1.2 Implementar `extraerMarcadoresVehiculo(respuesta) ([]uint, string)`: extrae ids únicos en orden de aparición (máx. 5), remueve los marcadores del texto visible
- [x] 1.3 En `Responder`: construir el set de ids servidos desde `construirContextoStock`, filtrar los marcadores contra ese set y poblar `RespuestaChat.VehiculosMencionados`
- [x] 1.4 Tests: extracción deduplica y limpia el texto; ids fuera del contexto se descartan; sin marcadores devuelve lista vacía

## 2. Backend — handler

- [x] 2.1 Exponer `vehiculosMencionados` (omitempty) en la respuesta JSON de `POST /api/chatbot/mensajes`

## 3. Frontend

- [x] 3.1 `types/chatbot.ts`: agregar `vehiculosMencionados?: number[]` a `RespuestaChatbot`
- [x] 3.2 `Chatbot.tsx`: renderizar chips "Ver ficha" bajo los mensajes del asistente cuando hay ids; navegar a `/catalogo/:id`; no enviar marcadores en el historial

## 4. Verificación y documentación

- [x] 4.1 `go build ./...`, `go vet ./...`, tests backend en verde; `npm run build` y tests frontend en verde
- [x] 4.2 Prueba E2E manual: preguntar por un vehículo del stock y verificar que aparece el chip que lleva a la ficha (verificada por API contra el stack Docker: respuesta sobre el Corolla con `vehiculosMencionados: [1]` sin marcadores visibles)
- [x] 4.3 Actualizar `docs/api.md` (campo nuevo) y AGENTS.md si corresponde
