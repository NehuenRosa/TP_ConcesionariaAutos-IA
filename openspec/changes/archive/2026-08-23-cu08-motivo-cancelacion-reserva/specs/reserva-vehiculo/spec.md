# Delta spec: reserva-vehiculo

## MODIFIED Requirements

### Requirement: Gestión de reservas por el vendedor

El sistema SHALL exponer endpoints bajo `/api/reservas` para que un vendedor
gestione las reservas: listar reservas con filtro opcional por estado, confirmar
la venta de una reserva `activa` (cambiando la reserva a `vendida` y el vehículo
a `vendido`) y cancelar una reserva `activa` (cambiando la reserva a `cancelada`
y el vehículo a `disponible`). Al cancelar desde la gestión del vendedor, el
cuerpo SHALL incluir un **motivo** explicativo no vacío que el sistema SHALL
persistir junto a la reserva y mostrar al cliente; sin motivo SHALL responder
`400`. Si el usuario autenticado no es vendedor o administrador, SHALL responder
con error `403`; si la reserva no existe, SHALL responder con error `404`; si la
reserva no está `activa`, SHALL responder con error `409`. La transición de
estado de la reserva y del vehículo SHALL persistirse de forma atómica.

#### Scenario: Listado de reservas del vendedor

- **WHEN** un vendedor solicita el listado de reservas
- **THEN** el sistema responde con las reservas, incluyendo vehículo y cliente,
  ordenadas por fecha de creación

#### Scenario: Confirmación de venta

- **WHEN** un vendedor confirma la venta de una reserva en estado `activa`
- **THEN** el sistema cambia la reserva a `vendida` y el vehículo a `vendido`

#### Scenario: Cancelación de reserva por el vendedor

- **WHEN** un vendedor cancela una reserva en estado `activa` indicando un
  motivo (por ejemplo, comprobante ilegible o monto incorrecto)
- **THEN** el sistema cambia la reserva a `cancelada`, guarda el motivo y el
  vehículo pasa a `disponible`

#### Scenario: Cancelación sin motivo

- **WHEN** un vendedor cancela una reserva sin indicar motivo o con texto vacío
- **THEN** el sistema responde `400` con un mensaje en español y la reserva
  permanece `activa`

#### Scenario: Transición inválida

- **WHEN** un vendedor intenta confirmar o cancelar una reserva que no está en
  estado `activa`
- **THEN** el sistema responde con error `409` y un mensaje en español

#### Scenario: Usuario sin permisos

- **WHEN** un cliente intenta usar un endpoint de gestión de reservas
- **THEN** el sistema responde con error `403` y un mensaje en español

## ADDED Requirements

### Requirement: Motivo de cancelación visible para el cliente

Las respuestas que describen una reserva SHALL incluir el campo
`motivoCancelacion` cuando exista. El panel de reservas del cliente SHALL
mostrar ese motivo junto al estado `cancelada` para que entienda por qué su
reserva fue anulada por la concesionaria. La cancelación iniciada por el
propio cliente no registra motivo.

#### Scenario: Cliente lee el motivo

- **WHEN** el vendedor canceló una reserva con motivo y el dueño abre sus
  reservas
- **THEN** ve la reserva `cancelada` acompañada del texto del motivo

#### Scenario: Cancelación propia sin motivo

- **WHEN** el cliente canceló él mismo su reserva
- **THEN** la reserva se muestra `cancelada` sin bloque de motivo
