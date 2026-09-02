# Proposal: cu13-notificaciones-cotizaciones

## Why

El sistema ya avisa mensajes nuevos de **consultas** (puntito rojo en el menú y
toast, vía `/api/notificaciones/contador`), pero las **cotizaciones** quedan
afuera: cuando el vendedor responde una cotización el cliente no se entera, y
cuando un cliente escribe o retoma una cotización ningún vendedor recibe aviso
para tomarla. El equipo comercial pide "un puntito rojo donde tenga mensajes"
también para este canal.

## What Changes

- Modelo: `MensajeCotizacion` suma dos marcas de lectura (`LeidoPorCliente`,
  `LeidoPorVendedor`) que la auto-migración agrega como columnas.
- Semántica de no leídos:
  - **Cliente**: mensajes con remitente `vendedor` sin marcar como leídos por
    él (las respuestas de la IA llegan sincrónicas en el mismo request, así
    que no cuentan).
  - **Vendedor**: mensajes con remitente `cliente` en cotizaciones **abiertas**
    sin asignar o asignadas a él.
- Backend:
  - Repositorio/servicio de cotizaciones: contar no leídos por usuario y
    marcar como leídos al abrir el hilo (vista cliente y vista personal).
  - `GET /api/notificaciones/contador` pasa a devolver
    `{ "contador": total, "consultas": n, "cotizaciones": m }` (el campo
    `contador` se mantiene por compatibilidad).
- Frontend:
  - `useNotificaciones` expone los dos contadores y LayoutBase pinta el
    puntito rojo **por sección**: Mis Consultas / Bandeja (vendedor) /
    Mis Cotizaciones (cliente) / Cotizaciones IA (vendedor).
  - El toast de "mensaje nuevo" se dispara cuando sube cualquiera de los dos
    contadores.
  - Los chats de cotización marcan los mensajes como leídos al cargarlos y en
    cada refresh periódico.

## Capabilities

### New Capabilities

- `notificaciones-mensajes`: contador unificado de mensajes no leídos
  (consultas + cotizaciones), marcado de lectura en cotizaciones y señales
  visuales (puntito por sección y toast) en el menú.
