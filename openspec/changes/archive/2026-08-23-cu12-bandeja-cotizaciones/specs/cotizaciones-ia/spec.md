# Delta spec: cotizaciones-ia

## ADDED Requirements

### Requirement: Bandeja de cotizaciones del vendedor

El sistema SHALL exponer `GET /api/cotizaciones/bandeja` (rol vendedor o
administrador) que liste todas las cotizaciones ordenadas por la fecha de su
último mensaje en orden descendente, incluyendo datos del cliente y del
vehículo, el estado (`abierta`/`cerrada`), si tiene vendedor asignado y un
preview del último mensaje descifrado. Si no es vendedor ni administrador,
SHALL responder `403`.

#### Scenario: Listado con cotizaciones

- **WHEN** un vendedor solicita la bandeja
- **THEN** recibe las cotizaciones con cliente, vehículo, estado de atención y
  preview del último mensaje, ordenadas por actividad reciente

#### Scenario: Usuario sin permisos

- **WHEN** un cliente solicita la bandeja
- **THEN** el sistema responde `403`

### Requirement: Atención del hilo completo

El sistema SHALL exponer `GET /api/cotizaciones/:id/personal` (rol vendedor o
administrador) que devuelva la cotización con todos sus mensajes descifrados,
el vehículo y los datos de contacto del cliente. Si no existe, `404`.

#### Scenario: Vendedor lee la conversación

- **WHEN** un vendedor abre una cotización desde su bandeja
- **THEN** ve todo el hilo cliente ↔ IA descifrado junto a la ficha del
  vehículo y los datos del cliente

### Requirement: Tomar una cotización

El sistema SHALL exponer `PUT /api/cotizaciones/:id/tomar` para que un
vendedor autenticado se asigne una cotización `abierta` sin vendedor,
registrando su usuario y la fecha de toma. Si otra persona ya la tomó, SHALL
responder `409`; si está cerrada, `409`. Volver a tomarla el mismo vendedor
SHALL ser idempotente.

#### Scenario: Toma exitosa

- **WHEN** un vendedor toma una cotización abierta sin asignar
- **THEN** queda asignada a él con fecha de toma y pasa al estado
  "en atención"

#### Scenario: Ya tomada por otro

- **WHEN** un vendedor intenta tomar una cotización asignada a otro usuario
- **THEN** el sistema responde `409` con un mensaje en español

### Requirement: Respuesta manual del vendedor

El sistema SHALL exponer `POST /api/cotizaciones/:id/mensajes-vendedor` para
que el vendedor asignado envíe mensajes dentro del hilo. El contenido SHALL
persistirse cifrado con remitente `"vendedor"`. Mensaje vacío → `400`;
cotización inexistente → `404`; cotización cerrada → `409`. No SHALL invocarse
al LLM en este flujo.

#### Scenario: El vendedor responde

- **WHEN** el vendedor asignado envía un texto por el endpoint de mensajes
- **THEN** el mensaje queda en el hilo como remitente `"vendedor"` y el
  cliente lo ve en su panel

### Requirement: Apagado de la IA cuando hay vendedor

Mientras una cotización tenga vendedor asignado, el envío de mensajes del
cliente SHALL seguir guardando el mensaje pero SHALL NOT generar respuesta del
asistente; la respuesta SHALL indicar `atendidaPorVendedor: true` para que la
interfaz avise que el vendedor responde personalmente.

#### Scenario: Cliente escribe en hilo atendido

- **WHEN** el cliente envía un mensaje en una cotización con vendedor asignado
- **THEN** su mensaje se guarda, no aparece respuesta automática y la UI le
  informa que un vendedor está atendiendo la conversación

### Requirement: Cierre desde el vendedor

El sistema SHALL permitir cerrar una cotización `abierta` mediante
`PUT /api/cotizaciones/:id/cerrar-personal` (rol vendedor o administrador).
Si ya está cerrada, SHALL responder `409`.

#### Scenario: Cierre comercial

- **WHEN** el vendedor cierra la conversación tras atenderla
- **THEN** la cotización pasa a `cerrada` y deja de recibir mensajes

### Requirement: Panel del cliente con atención humana

El panel de cotizaciones del cliente SHALL diferenciar visualmente los
mensajes del vendedor de los del asistente, SHALL mostrar un aviso cuando hay
un vendedor atendiendo el hilo y SHALL refrescar periódicamente para mostrar
las respuestas nuevas.

#### Scenario: Respuesta del vendedor en el panel del cliente

- **WHEN** el vendedor respondió dentro del hilo y el cliente tiene abierto su
  panel de cotizaciones
- **THEN** la respuesta aparece diferenciada como mensaje del vendedor, junto
  al aviso de que hay atención humana
