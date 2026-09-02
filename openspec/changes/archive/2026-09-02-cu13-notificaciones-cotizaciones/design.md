# Design: cu13-notificaciones-cotizaciones

## Contexto

`/api/notificaciones/contador` (handlers/notificaciones.go) hoy cuenta mensajes
no leídos solo de consultas: `ConsultaService.ListarPorUsuario` +
`MensajeService.ContarNoLeidosPorConsultas`, que usa `Mensaje.Leido` y
`remitente_id != yo`. El hook `useNotificaciones` hace polling cada 3 s y
LayoutBase pinta un único puntito en el enlace marcado con
`notificacion: true`, además de un toast cuando el total sube.

Las cotizaciones (`MensajeCotizacion`) no tienen marcas de lectura y su
remitente es una cadena (`cliente` / `ia` / `vendedor`), no un ID de usuario.
El contenido viaja cifrado, pero los contadores solo usan metadatos, así que el
cifrado no se ve afectado.

## Objetivos

- Un único endpoint de notificaciones para ambos canales.
- Puntito rojo por sección (no uno global) según dónde están los mensajes.
- Marcar como leídos al abrir el hilo, igual que ya pasa en consultas.

## Decisiones

### D1: marcas de lectura por columna booleana

`LeidoPorCliente bool` y `LeidoPorVendedor bool` en `mensaje_cotizaciones`
(default false). Alternativa descartada: tabla aparte de lecturas por usuario —
innecesaria porque las partes humanas son siempre dos fijas (cliente vs
personal) y GORM auto-migra las columnas sin script.

### D2: qué cuenta como no leído

- **Cliente**: mensajes con remitente `"vendedor"` y `LeidoPorCliente = false`
  en sus cotizaciones (abiertas o cerradas; si cerró con mensaje del vendedor
  sin leer, sigue mereciendo aviso).
- **Vendedor**: mensajes con remitente `"cliente"` y `LeidoPorVendedor = false`
  en cotizaciones **abiertas**, sin asignar o asignadas a él. Las cerradas no
  pitan; la IA nunca cuenta (su respuesta llega sincrónica en el mismo POST).

### D3: contrato del contador

`GET /api/notificaciones/contador` responde `{ "contador": consultas +
cotizaciones, "consultas": n, "cotizaciones": m }`. Mantener `contador` evita
romper el hook actual durante la migración del frontend. Si falla la parte de
cotizaciones, se devuelve `cotizaciones: 0` con log interno (degradación
graciosa igual que hoy con consultas).

### D4: marcado de lectura

Nuevo método de servicio `MarcarLeidasCotizacion(ctx, cotizacionID, lado)`:
- lado `cliente`: marca `LeidoPorCliente = true` en mensajes de remitentes
  `ia`/`vendedor`; lo llama la vista de detalle del cliente al cargar hilos.
- lado `personal`: marca `LeidoPorVendedor = true` en mensajes de remitente
  `cliente`; lo llama `ObtenerPersonal` (solo si hay vendedor asignado y es él;
  en bandeja sin asignar nadie marca).

Se ejecuta tras la lectura exitosa del hilo, dentro del handler correspondiente
(no duplica consultas en el listado de bandeja).

### D5: frontend por sección

- `useNotificaciones` pasa a devolver `{ cantidadConsultas,
  cantidadCotizaciones, nuevoAviso }` consumiendo los tres campos nuevos.
- Los enlaces del menú llevan `claveNotificacion: "consultas" |
  "cotizaciones"` y LayoutBase resuelve qué contador mostrarle a cada enlace;
  así cada sección tiene su propio puntito (y varios pueden picar a la vez).
- Toast: se dispara cuando sube cualquiera de los contadores respecto del
  poll anterior (mantiene el comportamiento actual de 8 s).

### Riesgos

- Polling cada 3 s ya existe; el conteo nuevo agrega una consulta liviana
  indexada por estado/remitente (columnas nuevas default false, índice sobre
  `cotizacion_id` ya presente).
- Hilos históricos: todos los mensajes viejos quedan "no leídos" al migrar →
  primer conteo alto. Se acepta: abrir cada hilo lo limpia; alternativa
  (backfill marcando todo leído) escondría mensajes reales pendientes.
