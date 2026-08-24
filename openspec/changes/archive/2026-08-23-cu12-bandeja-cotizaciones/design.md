# Design: cu12-bandeja-cotizaciones

## Contexto

Las cotizaciones ya existen: se crean desde el chatbot cuando el usuario pide
cotizar (marcador `[COTIZACION:<id>]`), los mensajes se guardan **cifrados**
(`cifrador.Cifrar/Descifrar`) y el cliente los ve en `/mis-cotizaciones/:id`.
Falta todo el lado comercial: ver los hilos, tomarlos y responderlos.

## Decisiones

- **D1 — Extender `Cotizacion`, no crear entidad nueva.** Columnas nuevas
  nullable (`VendedorID *uint`, `FechaToma *time.Time`) + relación `Vendedor`.
  AutoMigrate agrega las columnas sin tocar datos existentes. Nuevo remitente
  `"vendedor"` en `MensajeCotizacion` (varchar(20) ya lo permite).
- **D2 — Descifrado server-side, siempre.** El vendedor nunca accede al
  cifrado: `ListarBandeja` descifra solo el último mensaje (preview) y la
  vista de atención descifra el hilo completo. Los mensajes del vendedor se
  cifran igual que los del cliente.
- **D3 — IA condicionada.** `EnviarMensaje` (cliente) consulta si hay
  `VendedorID`: con vendedor guarda solo el mensaje del cliente y devuelve un
  resultado con `AtendidaPorVendedor = true` (sin llamada al LLM → ahorro de
  tokens/latencia); sin vendedor mantiene el comportamiento actual.
- **D4 — Rutas separadas para personal.** `/bandeja`, `/:id/personal`,
  `/:id/tomar`, `/:id/mensajes-vendedor`, `/:id/cerrar-personal` dentro del
  grupo autenticado de `/cotizaciones` con `middleware.ExigirRol("vendedor")`
  (el middleware de roles ya admite administrador). Evitamos sobrecargar los
  endpoints del cliente.
- **D5 — DTOs.** Bandeja: `{id, vehiculo{marca,modelo,anio,foto}, cliente
  {nombre,email}, estado, asignadaA, fechaToma, ultimoMensaje {remitente,
  contenido, fecha}}`. Atención: cotización completa + mensajes con remitente
  legible. El handler arma DTOs desde los modelos descifrados (los campos
  `Contenido` cifrados quedan ocultos con `json:"-"` como hoy).
- **D6 — Frontend vendedor.** Página única `/vendedor/cotizaciones`: listado a
  la izquierda (o arriba en mobile) y panel de atención al seleccionar; botón
  "Tomar" habilita la caja de respuesta; "Cerrar conversación" con
  confirmación. Link "Cotizaciones" en el menú del vendedor.
- **D7 — Frontend cliente.** `MisCotizaciones.tsx` distingue remitente
  `"vendedor"` (estilo neutro + nombre), muestra banner "Un vendedor está
  atendiendo esta conversación" cuando corresponde y hace polling cada 10 s
  mientras la página está abierta (patrón liviano, sin websockets).
- **D8 — Permisos de respuesta.** Solo el vendedor asignado responde; si el
  hilo no tiene asignación, el botón Tomar es previo obligatorio (409 en
  mensajes sin tomarla). Administrador puede cerrar pero no suplantar al
  vendedor asignado.

## Riesgos

- Cifrado/descifrado por mensaje en listados grandes: mitigado porque la
  bandeja solo descifra el último mensaje de cada hilo.
- Polling del cliente: 10 s por pestaña abierta es aceptable para desarrollo;
  no introducimos websockets fuera de stack.

## Migración

AutoMigrate agrega columnas nullable; cotizaciones viejas quedan "sin
atender" y aparecen en la bandeja normalmente.
