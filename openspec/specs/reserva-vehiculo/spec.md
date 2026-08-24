# reserva-vehiculo Specification

## Purpose
Permite a los clientes reservar una unidad disponible del catálogo, bloqueándola
mientras la reserva esté activa, y al vendedor confirmar la venta o cancelar la
reserva liberando la unidad.
## Requirements
### Requirement: Creación de reserva

El sistema SHALL exponer un endpoint `POST /api/reservas` que permita a un
cliente autenticado reservar un vehículo disponible indicando `vehiculoId`. Al
crear la reserva, el sistema SHALL cambiar el estado del vehículo a `reservado`,
de modo que deje de aparecer en el catálogo público. Si el vehículo no existe o
no está en estado `disponible`, SHALL responder con error `404`; si el vehículo
no está disponible para reservar (por ejemplo, ya está reservado o vendido),
SHALL responder con error `409`; si el cliente no está autenticado, SHALL
responder con error `401`. La creación de la reserva y el cambio de estado del
vehículo SHALL persistirse de forma atómica.

#### Scenario: Reserva válida

- **WHEN** un cliente autenticado reserva un vehículo en estado `disponible`
- **THEN** el sistema crea la reserva con estado `activa`, cambia el vehículo a
  `reservado` y devuelve la reserva creada

#### Scenario: Vehículo no disponible

- **WHEN** un cliente reserva un vehículo que no existe o no está en estado
  `disponible`
- **THEN** el sistema responde con error `404` y un mensaje en español

#### Scenario: Vehículo ya reservado o vendido

- **WHEN** un cliente intenta reservar un vehículo que ya está reservado o
  vendido
- **THEN** el sistema responde con error `409` y un mensaje en español

#### Scenario: Cliente no autenticado

- **WHEN** un visitante no autenticado intenta reservar un vehículo
- **THEN** el sistema responde con error `401` y un mensaje en español

### Requirement: Mis reservas del cliente

El sistema SHALL exponer un endpoint `GET /api/reservas/mis-reservas` que
devuelva las reservas del cliente autenticado con su vehículo y estado, y un
endpoint `DELETE /api/reservas/:id` que permita al cliente cancelar una reserva
propia en estado `activa`. Al cancelar, el sistema SHALL liberar el vehículo
volviéndolo a estado `disponible`. Si la reserva no pertenece al cliente o no
existe, SHALL responder con error `404`; si la reserva no está activa, SHALL
responder con error `409`.

#### Scenario: Listado de reservas del cliente

- **WHEN** un cliente solicita sus reservas
- **THEN** el sistema responde con las reservas del cliente, incluyendo el
  vehículo y el estado, ordenadas por fecha de creación

#### Scenario: Cancelación de reserva propia

- **WHEN** un cliente cancela una reserva propia en estado `activa`
- **THEN** el sistema cambia la reserva a `cancelada` y el vehículo vuelve a
  `disponible`

#### Scenario: Cancelación de reserva ajena

- **WHEN** un cliente intenta cancelar una reserva que pertenece a otro cliente
- **THEN** el sistema responde con error `404` y un mensaje en español

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

### Requirement: Página de reserva de vehículo

El sistema SHALL ofrecer una página de reserva en la ruta `/catalogo/:id/reservar`
accesible para clientes autenticados que muestre la unidad a reservar y un botón
para confirmar la reserva consumiendo el endpoint `POST /api/reservas`. La
página SHALL mostrar un mensaje en español cuando la reserva se crea con éxito y
un mensaje de error en español cuando falla (incluido el caso de unidad ya no
disponible).

#### Scenario: Acceso a la página de reserva

- **WHEN** un cliente autenticado navega a `/catalogo/:id/reservar` de un
  vehículo disponible
- **THEN** el sistema muestra la unidad a reservar y la acción de confirmar la
  reserva

#### Scenario: Reserva exitosa

- **WHEN** un cliente confirma la reserva de una unidad disponible
- **THEN** el sistema crea la reserva y muestra un mensaje de éxito en español

#### Scenario: Unidad ya no disponible

- **WHEN** un cliente intenta reservar una unidad que dejó de estar disponible
- **THEN** el sistema muestra un mensaje de error en español indicando que la
  unidad ya no está disponible

### Requirement: Página de mis reservas del cliente

El sistema SHALL ofrecer una página en la ruta `/mis-reservas` accesible para
clientes que liste sus reservas con el vehículo asociado y permita cancelar las
que estén en estado `activa`.

#### Scenario: Listado de reservas propias

- **WHEN** un cliente navega a `/mis-reservas`
- **THEN** el sistema muestra sus reservas con el vehículo y la posibilidad de
  cancelar las activas

### Requirement: Página de gestión de reservas del vendedor

El sistema SHALL ofrecer una página de gestión de reservas en la ruta
`/vendedor/reservas` accesible para vendedores y administradores que liste las
reservas, permita filtrar por estado y ejecutar las acciones de confirmar la
venta o cancelar la reserva.

#### Scenario: Listado con acciones

- **WHEN** un vendedor navega a `/vendedor/reservas`
- **THEN** el sistema muestra las reservas con las acciones de confirmar venta y
  cancelar según su estado

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

