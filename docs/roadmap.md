# Roadmap del proyecto

Backlog de casos de uso pendientes de planificar con OpenSpec.
Cada caso de uso se planifica con un change propio cuando se vaya a implementar.

| ID | Caso de uso | Estado |
|----|-------------|--------|
| CU-01 | Autenticación y roles (registro, login, JWT, 3 roles) | Implementado |
| CU-02 | Gestión de vehículos (ABM) | Implementado |
| CU-03 | Catálogo público (listado paginado + detalle) | Implementado |
| CU-04 | Búsqueda y filtrado avanzado | Implementado |
| CU-05 | Consulta / cotización | Implementado |
| CU-06 | Gestión de consultas (bandeja del vendedor) | Implementado (junto con CU-05) |
| CU-07 | Turno de test drive | Implementado |
| CU-08 | Reserva de vehículo | Implementado (con seña y comprobante) |
| CU-09 | Panel de administración (dashboard) | Implementado |
| CU-10 | Chatbot asistente (LangChain + Gemini/Ollama) | Resuelto |
| CU-11 | Login con Google (Google Identity Services + JWT propio) | Implementado |
| CU-12 | Bandeja de cotizaciones con IA para el vendedor | Implementado |

> CU-10 resuelto: incluye tasación con valores reales de la Guía de Precios de
> la CCA (vía API pública de ArgAutos), chat con contexto del stock disponible,
> comparación en vivo y widget flotante en el frontend. Detalle en
> `openspec/specs/chatbot-asistente/`. Extensión: la IA marca los vehículos que
> menciona y el widget muestra chips "Ver ficha" hacia el catálogo.
>
> CU-11 implementado: alta y vinculación automática de cuentas por email,
> verificación del ID token en el backend (JWKS de Google) y emisión del JWT
> propio del sistema. Requiere `GOOGLE_CLIENT_ID` para habilitarse. Detalle en
> `openspec/specs/autenticacion-google/`.
>
> CU-08 con seña: al reservar, el cliente tiene 2 horas para transferir el
> 5 % del valor del vehículo (CBU/alias configurables por entorno) y subir el
> comprobante; vencido el plazo sin comprobante, la reserva se anula sola y la
> unidad vuelve al catálogo. El vendedor revisa el comprobante antes de
> confirmar la venta; si la cancela, debe indicar un **motivo** que el cliente
> ve en Mis Reservas. Detalle en `openspec/specs/reserva-vehiculo/`.
>
> CU-12 implementado: bandeja `/vendedor/cotizaciones` donde el vendedor toma
> conversaciones de cotización iniciadas con la IA; al tomarla, la IA queda
> silenciada y responde el vendedor en su nombre. Detalle en
> `openspec/specs/cotizaciones-ia/`.
>
> Mejoras: las respuestas de la IA en una cotización cuentan como mensaje no
> leído para el cliente, así el aviso global (polling + toast) se dispara igual
> que cuando responde un vendedor aunque el cliente haya salido de la pestaña.
> El aviso dirige al canal correcto (Mis cotizaciones / Mis consultas según el
> rol y qué bandeja subió). Además, todas las pantallas de cotizaciones y
> consultas (cliente y vendedor) tienen un enlace "Ver ficha" hacia el vehículo
> en el catálogo (`/catalogo/:id`).
>
> CU-07: además de cancelar, el cliente puede **eliminar** un turno con baja
> lógica (`borrado_por_cliente`): `DELETE /api/test-drives/:id/eliminar` lo
> cancela si estaba activo (libera la franja) y lo oculta de Mis test drives.
> El vendedor sigue viéndolo con su estado real.
>
> Cotizaciones: al crearlas, cerrarlas o tomarlas, la API devuelve los mensajes
> ya descifrados (no expone el texto cifrado en las respuestas). En el
> frontend, una cotización cerrada se muestra con un panel "Cotización cerrada"
> y oculta la conversación.
>
> UX: al navegar entre páginas, la barra de scroll vuelve automáticamente al
> inicio (scroll to top en `LayoutBase`).
>
> Escalabilidad de conversaciones (cotizaciones y consultas):
> 1. **Fetch incremental por `desdeId`**: abrir un hilo baja el historial
>    completo; el polling (5 s consultas, 10 s cotizaciones) trae solo los
>    mensajes con `id > desdeId` (`GET /consultas/:id/mensajes/nuevos?desdeId=`,
>    `GET /cotizaciones/:id/mensajes`, `GET /cotizaciones/:id/mensajes/personal`).
>    Se usa el id y no el timestamp para no saltearse mensajes de un mismo
>    segundo.
> 2. **Recorte del historial del LLM**: `MaximoTurnosHistorial = 10` (antes 20);
>    `aTurnosChat` recorta y el widget envía `.slice(-10)`.
> 3. **Retención**: env `RETENCION_CONVERSACIONES_DIAS` (default 180). Un job
>    interno (al arrancar y cada 1 h) purga consultas y cotizaciones cerradas
>    con `updated_at` anterior al corte, junto con sus mensajes, en transacción.
> 4. **Índices compuestos**: `idx_cotizacion_mensajes_hilo (cotizacion_id,
>    created_at)` y `idx_mensajes_consulta_hilo (consulta_id, created_at)` para
>    no tocar el índice de PK al leer por id del hilo.

## Orden sugerido de implementación

1. **CU-01** Autenticación y roles (base para autorizar el resto).
2. **CU-02** Gestión de vehículos (alimenta el catálogo).
3. **CU-03** Catálogo público (depende de CU-02).
4. **CU-04** Búsqueda y filtrado (amplía CU-03).
5. **CU-05 / CU-06** Consultas y su gestión.
6. **CU-07** Test drives.
7. **CU-08** Reservas.
8. **CU-09** Panel de administración.
9. **CU-10** Chatbot.

## Proceso para planificar un CU

1. Crear el change: `openspec new change <nombre-en-kebab-case>`
2. Completar los artefactos: proposal → specs → design → tasks.
3. Validar: `openspec validate --strict`
4. Implementar y archivar al finalizar.
