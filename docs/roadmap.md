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
