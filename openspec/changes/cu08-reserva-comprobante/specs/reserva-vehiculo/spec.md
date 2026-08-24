# Delta spec: reserva-vehiculo

## ADDED Requirements

### Requirement: Datos de transferencia para la seña

El sistema SHALL exponer el endpoint autenticado
`GET /api/reservas/datos-transferencia?vehiculoId=<id>` que devuelva el CBU y
el alias de la concesionaria (configurados vía variables de entorno) y el
monto de la seña, calculado como el **5 %** del precio del vehículo indicado.
El monto SHALL calcularse siempre en el backend a partir del precio vigente
del vehículo. Si el vehículo no existe o no está disponible, SHALL responder
`404`; si falta el parámetro `vehiculoId`, `400`.

#### Scenario: Consulta de datos para la seña

- **WHEN** un cliente autenticado consulta los datos de transferencia de un
  vehículo disponible con precio $10.000.000
- **THEN** el sistema responde `200` con el CBU, el alias y el monto de seña
  $500.000

#### Scenario: Vehículo inexistente o no disponible

- **WHEN** un cliente consulta los datos de transferencia de un vehículo que
  no existe o no está en estado `disponible`
- **THEN** el sistema responde con error `404` y un mensaje en español

### Requirement: Envío del comprobante de seña

El sistema SHALL exponer el endpoint autenticado
`POST /api/reservas/:id/comprobante` (multipart/form-data, campo
`comprobante`) para que el dueño de una reserva `activa` adjunte la imagen
del comprobante de la transferencia. El archivo SHALL ser JPG, PNG o WebP de
hasta 5 MB; ante archivo ausente, formato inválido o peso excedido SHALL
responder `400`. Si la reserva no existe o es ajena al cliente SHALL responder
`404`; si la reserva no está `activa` o ya venció sin comprobante SHALL
responder `409`. Al aceptarse, el sistema SHALL registrar fecha y hora del
envío, detener el cómputo del plazo de expiración y dejar la imagen
disponible para el vendedor. El envío SHALL permitirse más de una vez
mientras la reserva esté activa (reemplaza la imagen anterior). El sistema
SHALL exponer además `GET /api/reservas/:id/comprobante` que sirva la imagen
al dueño de la reserva y a vendedores/administradores (`404` si nunca se
envió, `403` a terceros).

#### Scenario: Envío válido dentro del plazo

- **WHEN** el dueño de una reserva activa dentro de las 2 horas sube una
  imagen JPG de 1 MB
- **THEN** el sistema responde `200` con la reserva actualizada, registra el
  envío y la reserva deja de estar sujeta a expiración automática

#### Scenario: Archivo inválido

- **WHEN** el cliente sube un archivo PDF o una imagen de 8 MB
- **THEN** el sistema responde con error `400` y un mensaje en español
  indicando el formato o tamaño admitido

#### Scenario: Reserva expirada

- **WHEN** el cliente intenta subir un comprobante después de que la reserva
  venció las 2 horas sin comprobante
- **THEN** el sistema responde con error `409` y un mensaje en español

#### Scenario: Visualización por el vendedor

- **WHEN** un vendedor solicita el comprobante de una reserva que lo tiene
- **THEN** el sistema devuelve la imagen; para un cliente distinto al dueño,
  responde `403`

### Requirement: Expiración automática sin comprobante

Una reserva `activa` sin comprobante enviado SHALL anularse automáticamente
al cumplirse **2 horas** desde su creación: pasa a estado `cancelada` y el
vehículo vuelve a `disponible` (vuelve a aparecer en el catálogo público).
La verificación SHALL ejecutarse periódicamente (job interno) y además
validarse perezosamente al operar sobre la reserva. La expiración SHALL ser
atómica y segura frente a carreras con la subida de un comprobante simultánea.

#### Scenario: Vencimiento del plazo

- **WHEN** transcurren 2 horas desde la creación de una reserva activa sin
  comprobante
- **THEN** el sistema la cambia a `cancelada` y el vehículo vuelve a
  `disponible`, reapareciendo en el catálogo

#### Scenario: Comprobante enviado a tiempo evita la expiración

- **WHEN** el cliente sube el comprobante a 1 hora y 50 minutos de creada la
  reserva
- **THEN** la reserva permanece `activa` aunque pase el plazo de las 2 horas,
  a la espera de la revisión del vendedor

## MODIFIED Requirements

### Requirement: Creación de reserva

El sistema SHALL exponer un endpoint `POST /api/reservas` que permita a un
cliente autenticado reservar un vehículo disponible indicando `vehiculoId`. Al
crear la reserva, el sistema SHALL cambiar el estado del vehículo a `reservado`,
de modo que deje de aparecer en el catálogo público, registrar el monto de la
seña (5 % del precio del vehículo) y fijar el vencimiento del comprobante en
**2 horas**. Si el vehículo no existe o
no está en estado `disponible`, SHALL responder con error `404`; si el vehículo
no está disponible para reservar (por ejemplo, ya está reservado o vendido),
SHALL responder con error `409`; si el cliente no está autenticado, SHALL
responder con error `401`. La creación de la reserva y el cambio de estado del
vehículo SHALL persistirse de forma atómica. La respuesta SHALL incluir el
monto de la seña y el momento exacto de vencimiento.

#### Scenario: Reserva válida

- **WHEN** un cliente autenticado reserva un vehículo en estado `disponible`
- **THEN** el sistema crea la reserva con estado `activa`, cambia el vehículo a
  `reservado`, fija el vencimiento del comprobante a 2 horas y devuelve la
  reserva con el monto de la seña

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
devuelva las reservas del cliente autenticado con su vehículo, estado, monto de
la seña, vencimiento del comprobante y si ya fue enviado, y un
endpoint `DELETE /api/reservas/:id` que permita al cliente cancelar una reserva
propia en estado `activa`. Al cancelar, el sistema SHALL liberar el vehículo
volviéndolo a estado `disponible`. Si la reserva no pertenece al cliente o no
existe, SHALL responder con error `404`; si la reserva no está activa, SHALL
responder con error `409`.

#### Scenario: Listado de reservas del cliente

- **WHEN** un cliente solicita sus reservas
- **THEN** el sistema responde con las reservas del cliente, incluyendo el
  vehículo, el estado, el monto de la seña, el vencimiento y si el comprobante
  fue enviado, ordenadas por fecha de creación

#### Scenario: Cancelación de reserva propia

- **WHEN** un cliente cancela una reserva propia en estado `activa`
- **THEN** el sistema cambia la reserva a `cancelada` y el vehículo vuelve a
  `disponible`

#### Scenario: Cancelación de reserva ajena

- **WHEN** un cliente intenta cancelar una reserva que pertenece a otro cliente
- **THEN** el sistema responde con error `404` y un mensaje en español

### Requirement: Gestión de reservas por el vendedor

El sistema SHALL exponer endpoints bajo `/api/reservas` para que un vendedor
gestione las reservas: listar reservas con filtro opcional por estado e
indicación de si cada reserva tiene comprobante enviado, confirmar
la venta de una reserva `activa` (cambiando la reserva a `vendida` y el vehículo
a `vendido`) y cancelar una reserva `activa` (cambiando la reserva a `cancelada`
y el vehículo a `disponible`). Antes de operar sobre una reserva activa sin
comprobante cuyo plazo venció, el sistema SHALL aplicarle la expiración
automática en lugar de permitir la transición. Si el usuario autenticado no es
vendedor o administrador, SHALL responder con error `403`; si la reserva no
existe, SHALL responder con error `404`; si la reserva no está `activa`, SHALL
responder con error `409`. La transición de estado de la reserva y del vehículo
SHALL persistirse de forma atómica.

#### Scenario: Listado de reservas del vendedor

- **WHEN** un vendedor solicita el listado de reservas
- **THEN** el sistema responde con las reservas, incluyendo vehículo, cliente,
  monto de seña y si el comprobante fue enviado, ordenadas por fecha de
  creación

#### Scenario: Confirmación de venta

- **WHEN** un vendedor verifica el comprobante y confirma la venta de una
  reserva en estado `activa`
- **THEN** el sistema cambia la reserva a `vendida` y el vehículo a `vendido`

#### Scenario: Cancelación de reserva por el vendedor

- **WHEN** un vendedor cancela una reserva en estado `activa` (por ejemplo, por
  un comprobante inválido)
- **THEN** el sistema cambia la reserva a `cancelada` y el vehículo a
  `disponible`

#### Scenario: Transición inválida

- **WHEN** un vendedor intenta confirmar o cancelar una reserva que no está en
  estado `activa`
- **THEN** el sistema responde con error `409` y un mensaje en español

#### Scenario: Usuario sin permisos

- **WHEN** un cliente intenta usar un endpoint de gestión de reservas
- **THEN** el sistema responde con error `403` y un mensaje en español

### Requirement: Página de reserva de vehículo

El sistema SHALL ofrecer una página de reserva en la ruta `/catalogo/:id/reservar`
accesible para clientes autenticados que muestre la unidad a reservar, los
datos de transferencia (CBU, alias y monto de la seña del 5 %) y un botón
para confirmar la reserva consumiendo el endpoint `POST /api/reservas`. Tras
confirmar, la página SHALL mostrar la confirmación con cuenta regresiva del
plazo de 2 horas y la acción para subir el comprobante. La
página SHALL mostrar un mensaje en español cuando la reserva se crea con éxito y
un mensaje de error en español cuando falla (incluido el caso de unidad ya no
disponible).

#### Scenario: Acceso a la página de reserva

- **WHEN** un cliente autenticado navega a `/catalogo/:id/reservar` de un
  vehículo disponible
- **THEN** el sistema muestra la unidad a reservar, el CBU/alias de la
  concesionaria y el monto de la seña antes de confirmar

#### Scenario: Reserva exitosa

- **WHEN** un cliente confirma la reserva de una unidad disponible
- **THEN** el sistema crea la reserva, muestra el éxito con el plazo de 2 horas
  y habilita la subida del comprobante

#### Scenario: Unidad ya no disponible

- **WHEN** un cliente intenta reservar una unidad que dejó de estar disponible
- **THEN** el sistema muestra un mensaje de error en español indicando que la
  unidad ya no está disponible

### Requirement: Página de mis reservas del cliente

El sistema SHALL ofrecer una página en la ruta `/mis-reservas` accesible para
clientes que liste sus reservas con el vehículo asociado, el estado, el monto
de la seña, el vencimiento del comprobante y su situación (pendiente de envío
o enviado), permita subir el comprobante de las activas pendientes y cancelar
las que estén en estado `activa`.

#### Scenario: Listado de reservas propias

- **WHEN** un cliente navega a `/mis-reservas`
- **THEN** el sistema muestra sus reservas con el vehículo, la seña, el plazo
  restante de las pendientes de comprobante, la acción de subirlo y la
  posibilidad de cancelar las activas

### Requirement: Página de gestión de reservas del vendedor

El sistema SHALL ofrecer una página de gestión de reservas en la ruta
`/vendedor/reservas` accesible para vendedores y administradores que liste las
reservas, indique si cada una tiene comprobante enviado y si está pendiente de
comprobante (con el plazo), permita visualizar el comprobante y ejecutar las
acciones de confirmar la
venta o cancelar la reserva.

#### Scenario: Listado con acciones

- **WHEN** un vendedor navega a `/vendedor/reservas`
- **THEN** el sistema muestra las reservas con el estado del comprobante, la
  acción de verlo y las acciones de confirmar venta y cancelar según su estado
