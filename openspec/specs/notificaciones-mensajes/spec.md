# notificaciones-mensajes Specification

## Purpose
TBD - created by syncing change cu13-notificaciones-cotizaciones. Update Purpose after archive.
## Requirements
### Requirement: Contador unificado de mensajes no leídos

El sistema SHALL exponer `GET /api/notificaciones/contador` (autenticado) que
devuelva `{"contador": <total>, "consultas": <n>, "cotizaciones": <m>}` donde
`consultas` son los mensajes no leídos en consultas (comportamiento existente)
y `cotizaciones` los mensajes no leídos en cotizaciones del usuario. Ante error
en un canal, SHALL devolver `0` para ese canal sin fallar el request.

#### Scenario: Conteo por canal

- **WHEN** el usuario autenticado consulta el contador
- **THEN** recibe el total junto con el desglose por canal (`consultas` y
  `cotizaciones`)

### Requirement: Mensajes de cotización no leídos para el cliente

El sistema SHALL contar como no leídos, para el cliente dueño, los mensajes con
remitente `vendedor` cuya marca `LeidoPorCliente` sea `false`. Los mensajes de
la IA SHALL NOT contar como no leídos (su respuesta llega sincrónica).

#### Scenario: Respuesta del vendedor sin leer

- **WHEN** un vendedor responde una cotización y el cliente aún no abrió el
  hilo después de esa respuesta
- **THEN** el contador de cotizaciones del cliente la incluye

#### Scenario: La IA no genera aviso

- **WHEN** la IA responde automáticamente una cotización
- **THEN** el contador del cliente no aumenta por ese mensaje

### Requirement: Mensajes de cotización no leídos para el personal

El sistema SHALL contar como no leídos, para el rol vendedor, los mensajes con
remitente `cliente` cuya marca `LeidoPorVendedor` sea `false`, en cotizaciones
abiertas sin asignar o asignadas a él. Las cotizaciones cerradas y las
asignadas a otro vendedor SHALL NOT contar.

#### Scenario: Cliente nuevo en la bandeja

- **WHEN** un cliente crea o retoma una cotización sin vendedor asignado
- **THEN** el contador de cotizaciones sube para todos los vendedores

#### Scenario: Cotización cerrada no avisa

- **WHEN** llega actividad sobre una cotización cerrada
- **THEN** ningún vendedor la cuenta como no leída

### Requirement: Marcado de lectura al abrir el hilo

Al abrir el hilo desde la vista cliente, el sistema SHALL marcar
`LeidoPorCliente = true` en los mensajes de remitentes `ia`/`vendedor`. Al
abrirlo desde la vista personal, SHALL marcar `LeidoPorVendedor = true` en los
mensajes de remitente `cliente`, únicamente cuando el solicitante es el
vendedor asignado.

#### Scenario: El cliente abre Mis Cotizaciones

- **WHEN** el cliente abre el detalle de su cotización
- **THEN** los mensajes del vendedor quedan leídos y el puntito desaparece

#### Scenario: El vendedor atiende la conversación

- **WHEN** el vendedor asignado abre `/cotizaciones/:id/personal`
- **THEN** los mensajes del cliente de esa cotización quedan leídos para él

### Requirement: Señales visuales por sección

La interfaz SHALL mostrar un puntito rojo independiente en cada sección con
mensajes no leídos (Mis Consultas y Mis Cotizaciones para el cliente; Bandeja
y Cotizaciones IA para el vendedor) y SHALL mostrar un toast cuando cualquiera
de los contadores aumente respecto del sondeo anterior.

#### Scenario: Puntito por sección

- **WHEN** un cliente tiene mensajes sin leer en consultas y otros en
  cotizaciones
- **THEN** ambos enlaces muestran su propio puntito rojo

#### Scenario: Toast al llegar un mensaje

- **WHEN** cualquier contador aumenta mientras la sesión está abierta
- **THEN** aparece el toast de mensaje nuevo durante algunos segundos
