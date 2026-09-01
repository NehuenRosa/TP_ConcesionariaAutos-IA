# API REST

Referencia de la API del backend (Go + Gin + GORM + JWT). Base URL por defecto:
`http://localhost:8080/api`.

## Autenticación

La mayoría de los endpoints requieren un token JWT en el header:

```
Authorization: Bearer <token>
```

El token se obtiene de `POST /auth/login` (o se registra una cuenta en
`POST /auth/registro`). Los roles autorizados son:

| Rol | Descripción |
|-----|-------------|
| `cliente` | Catálogo, consultas, reservas y test drives propios. |
| `vendedor` | Bandeja de consultas, reservas y turnos. |
| `administrador` | ABM de vehículos, usuarios y métricas. |

Formato de error en toda la API: `{"error": "mensaje en español"}` con código
HTTP adecuado (`400`, `401`, `403`, `404`, `409`, `500`).

## Endpoints públicos (sin token)

### `GET /health`

Verificación de salud.

```
200 OK {"estado": "ok"}
```

### `POST /auth/registro`

Crea una cuenta con rol `cliente`.

Cuerpo:

```json
{
  "nombre": "María Pérez",
  "email": "maria@correo.com",
  "password": "MiContrasena123!"
}
```

Respuestas: `201` con el usuario creado (sin `password`), `400` datos inválidos,
`409` email en uso.

### `POST /auth/login`

Inicia sesión y devuelve el token JWT.

Cuerpo:

```json
{
  "email": "cliente@concesionaria.local",
  "password": "Cliente123!"
}
```

Respuesta `200`:

```json
{
  "token": "eyJhbGciOi...",
  "usuario": { "id": 3, "nombre": "...", "email": "...", "rol": "cliente" }
}
```

`401` si las credenciales son inválidas.

### `GET /vehiculos`

Catálogo público paginado de vehículos **disponibles**, con búsqueda, filtros
combinables y ordenamiento.

| Parámetro | Descripción |
|-----------|-------------|
| `pagina` | Número de página (por defecto `1`). |
| `tamano` | Cantidad por página (por defecto `12`, máx. `100`). |
| `busqueda` | Texto libre (coincide con marca/modelo). |
| `marca`, `modelo` | Filtro exacto. |
| `tipo` | Tipo de vehículo. |
| `combustible` | Combustible. |
| `condicion` | `nuevo` o `usado`. |
| `anio_min`, `anio_max` | Rango de años. |
| `precio_min`, `precio_max` | Rango de precio. |
| `orden_por` | `precio` o `anio`. |
| `orden_direccion` | `asc` o `desc`. |

Respuesta `200`:

```json
{
  "datos": [
    { "id": 1, "marca": "Fiat", "modelo": "Cronos", "anio": 2025,
      "precio": 20000000, "condicion": "nuevo", "tipo": "sedán", "imagen": "..." }
  ],
  "pagina": 1,
  "tamano": 12,
  "total": 42
}
```

`400` si un filtro numérico es inválido.

### `GET /vehiculos/:id`

Ficha técnica completa (incluye `kilometraje`, `combustible`, `transmision`,
`estado` e `imagenes[]`) de un vehículo **disponible**.

`200` con el detalle, `404` si no existe o no está disponible (los vehículos
reservados/vendidos no se muestran en el catálogo público).

### `GET /test-drives/franjas`

Catálogo de franjas horarias de una hora (horario comercial).

```json
[
  { "id": "09:00", "inicio": "09:00", "fin": "10:00" },
  { "id": "10:00", "inicio": "10:00", "fin": "11:00" }
]
```

### `POST /chatbot/mensajes`

Chat con historial. Responde sobre el stock real y orienta a consultar o pedir
un test drive.

Cuerpo:

```json
{
  "mensaje": "¿Hay Fiat Cronos disponibles?",
  "historial": [ { "rol": "user", "contenido": "..." } ]
}
```

`200` `{"respuesta": "...", "cotizacionId": null, "vehiculosMencionados": [3, 7]}`.
`400` si el mensaje está vacío o es muy largo.
Si el proveedor LLM está caído, responde `200` con un mensaje en español que
orienta al usuario (el error interno se loguea).

- `cotizacionId`: id de la cotización creada cuando el usuario pidió cotizar y
  está autenticado (si no, `null`).
- `vehiculosMencionados`: ids de vehículos del stock que la IA mencionó con el
  marcador interno `[VEHICULO:<id>]`, ya validados contra el contexto servido
  (máx. 5). El frontend los muestra como chips "Ver ficha" que llevan a
  `/catalogo/:id`; los marcadores nunca llegan al texto visible.

### `POST /chatbot/tasacion`

Tasación por fotos (`multipart/form-data`). El modelo de visión solo
**identifica** el vehículo; el valor se compone en código con la Guía Oficial de
la CCA (ArgAutos, caché de 24 h). Nunca inventa montos.

| Campo | Descripción |
|-------|-------------|
| `fotos` | Archivos de imagen: hasta 5, JPG/PNG/WebP, máx. 5 MB c/u. |
| `descripcion` | Opcional: marca/modelo/año/estado/kilometraje del vehículo. |

`200` `{"respuesta": "..."}`. `400` si no hay fotos, son más de 5, pesan más de
5 MB o el formato no es válido. Si el proveedor LLM está caído o no se identifica
el vehículo, responde `200` con una respuesta honesta que orienta a la
concesionaria.

## Autenticado (requiere token)

### `GET /auth/perfil`

Devuelve el usuario autenticado.

## Vehículos (administrador)

### `GET /admin/vehiculos`

Listado de gestión con filtro opcional por estado y paginación.

| Parámetro | Descripción |
|-----------|-------------|
| `estado` | `disponible`, `reservado`, `vendido` o `dado_de_baja`. |
| `pagina`, `tamano` | Paginación (igual que el catálogo). |

`200` con `{ "datos": [...], "pagina", "tamano", "total" }`.

### `GET /admin/vehiculos/:id`

Ficha técnica completa de cualquier vehículo (incluye no disponibles).

### `POST /admin/vehiculos`

Da de alta un vehículo.

Cuerpo:

```json
{
  "marca": "Fiat",
  "modelo": "Cronos",
  "anio": 2025,
  "kilometraje": 0,
  "combustible": "nafta",
  "transmision": "manual",
  "tipo": "sedán",
  "precio": 20000000,
  "condicion": "nuevo",
  "estado": "disponible",
  "imagenes": ["https://.../foto1.jpg"]
}
```

`201` con el vehículo creado, `400` si hay datos inválidos o estado no permitido.

### `PUT /admin/vehiculos/:id`

Modifica ficha técnica, estado e imágenes.

### `DELETE /admin/vehiculos/:id`

Da de baja el vehículo (cambia su estado a `dado_de_baja`, no lo elimina).

## Usuarios (administrador)

### `GET /admin/usuarios`

Lista todos los usuarios.

### `POST /admin/usuarios`

Crea un usuario con rol explícito:

```json
{
  "nombre": "Juan",
  "email": "juan@correo.com",
  "password": "Secreta123!",
  "rol": "vendedor"
}
```

### `PUT /admin/usuarios/:id`

Actualiza nombre/rol (y contraseña si se envía).

### `DELETE /admin/usuarios/:id`

Elimina un usuario.

## Métricas (administrador)

### `GET /admin/metricas?periodo=30`

Métricas del panel. `periodo` es opcional (`7`, `30` o `90`; por defecto `30`).

```json
{
  "vehiculosPorEstado": [ { "estado": "disponible", "cantidad": 40 } ],
  "consultasPorPeriodo": [ { "fecha": "2026-08-16", "cantidad": 3 } ],
  "reservasActivas": 2,
  "reservasVendidas": 5,
  "testDrivesAgendados": 3,
  "testDrivesCompletados": 1,
  "consultasAbiertas": 4,
  "totalUsuarios": 12
}
```

`400` si el período no es válido.

## Consultas (cliente + vendedor)

Estados de consulta: `pendiente`, `en_conversacion`, `cerrada`.

### `POST /consultas`

El cliente crea una consulta asociada a un vehículo.

```json
{
  "vehiculoId": 1,
  "mensaje": "¿Cuál es el precio final con patentamiento?"
}
```

`201` con el resumen de la consulta, `400` mensaje vacío, `404` vehículo no
disponible.

### `GET /consultas/mis-consultas`

Lista las consultas del cliente autenticado.

### `GET /consultas/bandeja` *(vendedor)*

Bandeja de consultas de todos los clientes, con contador de mensajes nuevos.

### `PUT /consultas/:id/tomar` *(vendedor)*

Asigna la consulta al vendedor y la pasa a `en_conversacion`.

### `PUT /consultas/:id/cerrar` *(vendedor)*

Cierra la consulta.

### `DELETE /consultas/:id` *(vendedor)*

Elimina la consulta.

### `GET /consultas/:id/mensajes`

Historial de mensajes de la consulta (cliente participante o vendedor).

### `GET /consultas/:id/mensajes/nuevos?desdeId=`

Mensajes nuevos (con `id > desdeId`) para el polling incremental del chat.
`desdeId=0` devuelve el historial completo. Devuelve un array plano; usar la
última `id` de la respuesta como `desdeId` del siguiente pedido en lugar del
timestamp, para no saltearse mensajes de un mismo segundo.

### `POST /consultas/:id/mensajes`

Envía un mensaje. Cuerpo: `{ "contenido": "..." }`. Responde `201` con el
mensaje creado.

### `PUT /consultas/:id/mensajes/leidos`

Marca todos los mensajes de la consulta como leídos.

### `GET /notificaciones/contador`

Cantidad de mensajes no leídos del usuario autenticado, desglosada por canal:
`consultas` (bandeja clásica) y `cotizaciones` (cotizaciones iniciadas con la
IA, CU-12). Para el cliente cuentan las respuestas de vendedor; para el
vendedor, mensajes de cliente en cotizaciones abiertas sin asignar o asignadas
a él. Los mensajes de la IA no notifican (llegan en la misma respuesta del
chat). Abrir el hilo correspondiente marca sus mensajes como leídos.

```json
{ "contador": 3, "consultas": 1, "cotizaciones": 2 }
```

## Test drives (cliente + vendedor)

Estados de turno: `solicitado`, `confirmado`, `cancelado`, `completado`.

### `POST /test-drives`

El cliente solicita un turno.

```json
{
  "vehiculoId": 1,
  "fecha": "2026-08-20",
  "franja": "10:00"
}
```

`201` con el resumen del turno. `400` datos inválidos o turno en el pasado,
`404` vehículo no disponible, `409` turno superpuesto (no puede haber más de un
turno para la misma unidad en la misma fecha y franja).

### `GET /test-drives/mis-turnos`

Turnos del cliente autenticado.

### `GET /test-drives` *(vendedor)*

Todos los turnos.

### `PUT /test-drives/:id/confirmar` *(vendedor)*

Confirma el turno.

### `PUT /test-drives/:id/cancelar` *(vendedor)*

Cancela el turno.

### `PUT /test-drives/:id/completar` *(vendedor)*

Marca el turno como completado.

### `DELETE /test-drives/:id`

El cliente cancela su propio turno.

## Reservas (cliente + vendedor)

Estados de reserva: `activa`, `vendida`, `cancelada`. Mientras una reserva está
`activa`, el vehículo pasa a estado `reservado` y deja de mostrarse en el
catálogo público.

**Seña (CU-08)**: al reservar, el backend fija un plazo de **2 horas** para que
el cliente transfiera el **5 % del precio** del vehículo y suba la imagen del
comprobante. Vencido el plazo sin comprobante, la reserva se anula sola
(`cancelada`) y la unidad vuelve a `disponible`; lo aplica un job interno cada
30 s más chequeos perezosos al operar sobre la reserva. Las respuestas de
reserva incluyen `montoSenia` (calculado por el backend), `vencimientoComprobante`
y `comprobanteEnviadoAt`.

### `GET /reservas/datos-transferencia?vehiculoId=`

CBU/alias de la concesionaria (variables `CBU_CONCESIONARIA` /
`ALIAS_CONCESIONARIA`) y monto de la seña del vehículo indicado, solo si está
disponible. `404` si no existe o no está disponible.

### `POST /reservas`

El cliente reserva una unidad.

```json
{ "vehiculoId": 1 }
```

`201` con el resumen de la reserva (incluye monto de seña y vencimiento). `404`
vehículo no disponible, `409` el vehículo ya tiene una reserva activa.

### `POST /reservas/:id/comprobante`

El dueño sube la imagen del comprobante (multipart, campo `comprobante`;
JPG/PNG/WebP ≤ 5 MB) mientras la reserva esté activa y dentro del plazo; puede
reenviarse (reemplaza la anterior). `400` archivo inválido, `404` reserva ajena
o inexistente, `409` reserva no activa o plazo vencido.

### `GET /reservas/:id/comprobante`

Sirve la imagen del comprobante al dueño de la reserva o a
vendedores/administradores (`403` a otros clientes, `404` si nunca se envió).

### `GET /reservas/mis-reservas`

Reservas del cliente autenticado.

### `GET /reservas` *(vendedor)*

Todas las reservas.

### `PUT /reservas/:id/confirmar` *(vendedor)*

Confirma la venta: la reserva pasa a `vendida` y el vehículo a `vendido`.

### `PUT /reservas/:id/cancelar` *(vendedor)*

Cancela la reserva y libera la unidad (el vehículo vuelve a `disponible`). El
motivo es **obligatorio** y queda visible para el cliente.

Cuerpo:

```json
{ "motivo": "El comprobante no corresponde a la seña solicitada." }
```

`400` si el motivo está vacío (menos de 5 caracteres tras recortar espacios),
`404` si la reserva no existe y `409` si ya no está activa. El cliente lo ve en
Mis Reservas como `motivoCancelacion`.

### `DELETE /reservas/:id`

El cliente cancela su propia reserva (libera la unidad).

## Cotizaciones con IA (cliente + vendedor)

Las cotizaciones son conversaciones del cliente con el asistente sobre precios y
financiación de un vehículo; los mensajes se guardan cifrados. Endpoints del
cliente: `POST /cotizaciones`, `GET /cotizaciones/mis-cotizaciones`,
`GET /cotizaciones/:id`, `POST /cotizaciones/:id/mensajes`,
`PUT /cotizaciones/:id/cerrar`.

Un **vendedor** puede tomar una conversación para atenderla personalmente:
desde ese momento la IA queda silenciada (los mensajes del cliente se guardan
sin respuesta automática) y solo el vendedor que la tomó puede responder.

| Endpoint | Descripción |
|----------|-------------|
| `GET /cotizaciones/bandeja` *(vendedor)* | Todas las cotizaciones ordenadas por actividad, con cliente, vendedor asignado, estado y último mensaje descifrado. |
| `GET /cotizaciones/:id/personal` *(vendedor)* | Detalle completo con el hilo descifrado. |
| `PUT /cotizaciones/:id/tomar` *(vendedor)* | Asigna la conversación al vendedor autenticado y silencia la IA. Idempotente para el mismo vendedor; `409` si otro vendedor la tomó o está cerrada. |
| `POST /cotizaciones/:id/mensajes-vendedor` *(vendedor)* | Guarda el mensaje del vendedor (remitente `vendedor`) sin pasar por la IA. `400` sin haberla tomado o mensaje vacío/muy largo; `403` si la tomó otro vendedor; `409` si está cerrada. |
| `PUT /cotizaciones/:id/cerrar-personal` *(vendedor)* | Cierra la cotización desde el personal. |
| `GET /cotizaciones/:id/mensajes?desdeId=` | Fetch incremental del hilo (cliente). Devuelve `{ mensajes, total, estado, vendedor, fechaToma }` con los mensajes de `id > desdeId` ya descifrados; `desdeId=0` devuelve el historial completo. |
| `GET /cotizaciones/:id/mensajes/personal?desdeId=` *(vendedor)* | Igual que el anterior pero desde la vista del personal (marca los mensajes como leídos). |

El cliente ve en su panel quién lo atiende (`vendedor`) y los mensajes del
personal con remitente `vendedor`; ambos lados refrescan cada 10 s trayendo
solo lo nuevo (fetch incremental por `desdeId`): el historial completo se baja
una sola vez al abrir el hilo.

## Roles en los endpoints

| Endpoint | Público | Cliente | Vendedor | Administrador |
|----------|:-------:|:-------:|:--------:|:-------------:|
| `/health` | ✔ | ✔ | ✔ | ✔ |
| `/auth/registro`, `/auth/login` | ✔ | | | |
| `/vehiculos*` (catálogo) | ✔ | ✔ | ✔ | ✔ |
| `/test-drives/franjas` | ✔ | ✔ | ✔ | ✔ |
| `/chatbot/*` | ✔ | ✔ | ✔ | ✔ |
| `/auth/perfil` | | ✔ | ✔ | ✔ |
| `/consultas/mis-consultas`, `/consultas/*` (POST, mensajes) | | ✔ | ✔* | |
| `/consultas/bandeja`, `/consultas/:id/tomar\|cerrar\|DELETE` | | | ✔ | |
| `/test-drives` (POST, mis-turnos, DELETE) | | ✔ | ✔* | |
| `/test-drives` (gestión: listar, confirmar, cancelar, completar) | | | ✔ | |
| `/reservas` (POST, mis-reservas, DELETE, comprobante propio) | | ✔ | ✔* | |
| `/reservas` (gestión: listar, confirmar, cancelar) | | | ✔ | |
| `/reservas/:id/comprobante` (GET) | | ✔* | ✔ | ✔ |
| `/cotizaciones` (cliente: crear, mis-cotizaciones, hilo propio, cerrar) | | ✔ | ✔* | |
| `/cotizaciones/bandeja`, `/:id/personal`, `/:id/tomar`, `/:id/mensajes-vendedor`, `/:id/cerrar-personal` | | | ✔ | |
| `/admin/vehiculos*` | | | | ✔ |
| `/admin/usuarios*` | | | | ✔ |
| `/admin/metricas*` | | | | ✔ |

\* Consulta/mensajes y test drives/reservas propios según participante: los
mensajes de una consulta son accesibles solo por el cliente dueño, el vendedor
asignado y el administrador.
