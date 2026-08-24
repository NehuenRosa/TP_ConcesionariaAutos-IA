# Proposal: cu12-bandeja-cotizaciones

## Why

Cuando el asistente promete "un vendedor se va a poner en contacto" durante una
cotización, hoy esa promesa se incumple: las conversaciones de cotizaciones
(cliente ↔ IA, mensajes cifrados) son invisibles para el equipo comercial. El
vendedor necesita ver esos hilos, tomarlos y responderlos él mismo para
cerrar ventas reales.

## What Changes

- Modelo: `Cotizacion` suma `VendedorID *uint`, `Vendedor` y `FechaToma`;
  nuevo remitente de mensaje `"vendedor"` junto a `"cliente"` e `"ia"`.
- Backend (rutas bajo `/api/cotizaciones` con rol vendedor):
  - `GET /bandeja`: todas las cotizaciones con cliente, vehículo, estado de
    atención y preview del último mensaje (descifrado server-side).
  - `GET /:id/personal`: hilo completo descifrado para el personal.
  - `PUT /:id/tomar`: asigna el vendedor autenticado (409 si otro ya lo hizo).
  - `POST /:id/mensajes-vendedor`: respuesta manual del vendedor (cifrada,
    sin IA).
  - `PUT /:id/cerrar-personal`: cierre desde el lado del vendedor.
- **IA apagada al tomarla**: mientras haya vendedor asignado, los mensajes del
  cliente ya no generan respuesta automática; el backend responde con un
  indicador `atendidaPorVendedor` para que la UI lo comunique.
- Frontend:
  - Nueva página `/vendedor/cotizaciones` (BandejaCotizaciones) con listado y
    vista de atención: hilo completo, botón tomar, caja de respuesta y cierre.
  - `MisCotizaciones` del cliente muestra los mensajes del vendedor con estilo
    propio, avisa cuando hay vendedor atendiendo y refresca periódicamente.
- El contenido viaja cifrado en reposo como hoy: el servicio descifra al leer
  y cifra al escribir, también para los mensajes del vendedor.

## Capabilities

### New Capabilities

- `cotizaciones-ia`: gestión comercial de las conversaciones de cotización
  (bandeja, toma, respuesta del vendedor, apagado de la IA y cierre).

### Modified Capabilities

*(ninguna)*

## Impact

- Backend: `models/cotizacion.go`, repositorio, `services/cotizaciones.go`,
  `services/chatbot.go` (apagado condicional de IA), handler, router, tests.
- Frontend: páginas nuevas, `MisCotizaciones.tsx`, `api.ts`, `types/`,
  rutas y menú del vendedor.
- Migración automática segura: columnas nuevas nullable, sin tocar datos.
