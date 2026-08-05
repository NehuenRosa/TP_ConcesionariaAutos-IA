## Purpose

Permite a los clientes solicitar turnos de prueba de manejo para un vehículo
disponible, validando que no haya superposición para la misma unidad en la
misma fecha y franja horaria, y permite al vendedor gestionar esos turnos.

## ADDED Requirements

### Requirement: Solicitud de turno de test drive

El sistema SHALL exponer un endpoint `POST /api/test-drives` que permita a un
cliente autenticado solicitar un turno de prueba de manejo para un vehículo
disponible, indicando `vehiculoId`, `fecha` (en formato `YYYY-MM-DD`) y
`franjaHoraria`. El sistema SHALL validar que el vehículo existe y está en
estado `disponible`, que la fecha no es anterior a hoy y que la franja horaria
es una de las franjas predefinidas. Si la solicitud es válida, SHALL crear el
turno con estado `solicitado` y responder con el turno creado. Si el vehículo
no existe o no está disponible, SHALL responder con error `404`; si la fecha o
la franja son inválidas, SHALL responder con error `400`; si el cliente no está
autenticado, SHALL responder con error `401`.

#### Scenario: Solicitud válida

- **WHEN** un cliente autenticado solicita un test drive para un vehículo
  disponible con una fecha futura y una franja horaria válida
- **THEN** el sistema crea el turno con estado `solicitado` y lo devuelve

#### Scenario: Vehículo no disponible

- **WHEN** un cliente solicita un test drive para un vehículo que no existe o
  no está en estado `disponible`
- **THEN** el sistema responde con error `404` y un mensaje en español

#### Scenario: Fecha inválida

- **WHEN** un cliente solicita un test drive con una fecha anterior a hoy
- **THEN** el sistema responde con error `400` y un mensaje en español

#### Scenario: Franja horaria inválida

- **WHEN** un cliente solicita un test drive con una franja horaria que no está
  en el catálogo de franjas predefinidas
- **THEN** el sistema responde con error `400` y un mensaje en español

#### Scenario: Cliente no autenticado

- **WHEN** un visitante no autenticado intenta solicitar un test drive
- **THEN** el sistema responde con error `401` y un mensaje en español

### Requirement: Prevención de superposición de turnos

El sistema SHALL evitar que exista más de un turno activo para la misma unidad
en la misma fecha y franja horaria. Un turno activo es aquel en estado
`solicitado` o `confirmado`. Cuando un cliente solicita un turno que se
superpone con un turno activo existente de la misma unidad, el sistema SHALL
responder con error `409` y un mensaje en español indicando que el turno ya
está ocupado.

#### Scenario: Turno ocupado

- **WHEN** un cliente solicita un test drive para una unidad, fecha y franja
  que ya tiene un turno activo
- **THEN** el sistema responde con error `409` y un mensaje en español

#### Scenario: Turnos en distinta franja

- **WHEN** un cliente solicita un test drive para una unidad y fecha que ya
  tiene un turno activo pero en otra franja horaria
- **THEN** el sistema crea el turno sin conflictos

#### Scenario: Turnos cancelados no bloquean

- **WHEN** un cliente solicita un test drive para una unidad, fecha y franja
  que tuvo un turno cancelado o completado
- **THEN** el sistema crea el turno sin conflictos

### Requirement: Catálogo de franjas horarias

El sistema SHALL exponer el listado de franjas horarias predefinidas para los
turnos de test drive. Cada franja SHALL tener un identificador, una hora de
inicio y una hora de fin. El listado SHALL ser accesible sin autenticación.

#### Scenario: Obtener franjas disponibles

- **WHEN** un visitante solicita las franjas horarias disponibles
- **THEN** el sistema responde con el listado de franjas predefinidas

### Requirement: Gestión de turnos por el vendedor

El sistema SHALL exponer endpoints bajo `/api/test-drives` para que un vendedor
gestione los turnos: listar turnos (con filtro opcional por estado), confirmar
un turno `solicitado`, cancelar un turno y marcar un turno `confirmado` como
`completado`. El sistema SHALL permitir confirmar solo turnos en estado
`solicitado`, cancelar turnos en estado `solicitado` o `confirmado`, y
completar solo turnos en estado `confirmado`. Si el usuario autenticado no es
vendedor o administrador, SHALL responder con error `403`; si el turno no
existe, SHALL responder con error `404`; si la transición de estado no es
válida, SHALL responder con error `409`.

#### Scenario: Listado de turnos del vendedor

- **WHEN** un vendedor solicita el listado de turnos
- **THEN** el sistema responde con los turnos agendados, incluyendo vehículo y
  cliente, ordenados por fecha y franja

#### Scenario: Confirmación de turno

- **WHEN** un vendedor confirma un turno en estado `solicitado`
- **THEN** el sistema cambia el turno a estado `confirmado`

#### Scenario: Cancelación de turno

- **WHEN** un vendedor cancela un turno en estado `solicitado` o `confirmado`
- **THEN** el sistema cambia el turno a estado `cancelado`

#### Scenario: Completar turno

- **WHEN** un vendedor marca como completado un turno en estado `confirmado`
- **THEN** el sistema cambia el turno a estado `completado`

#### Scenario: Transición inválida

- **WHEN** un vendedor intenta confirmar un turno que no está en estado
  `solicitado` o completar un turno que no está en estado `confirmado`
- **THEN** el sistema responde con error `409` y un mensaje en español

#### Scenario: Usuario sin permisos

- **WHEN** un cliente intenta usar un endpoint de gestión de turnos
- **THEN** el sistema responde con error `403` y un mensaje en español

### Requirement: Mis turnos del cliente

El sistema SHALL exponer un endpoint `GET /api/test-drives/mis-turnos` que
devuelva los turnos del cliente autenticado, y un endpoint
`DELETE /api/test-drives/:id` que permita al cliente cancelar un turno propio
en estado `solicitado` o `confirmado`. Si el turno no pertenece al cliente o no
existe, SHALL responder con error `404`; si el turno no se puede cancelar por
su estado, SHALL responder con error `409`.

#### Scenario: Listado de los turnos del cliente

- **WHEN** un cliente solicita sus turnos
- **THEN** el sistema responde con los turnos del cliente, incluyendo el
  vehículo y el estado, ordenados por fecha

#### Scenario: Cancelación de turno propio

- **WHEN** un cliente cancela un turno propio en estado `solicitado` o
  `confirmado`
- **THEN** el sistema cambia el turno a estado `cancelado`

#### Scenario: Cancelación de turno ajeno

- **WHEN** un cliente intenta cancelar un turno que pertenece a otro cliente
- **THEN** el sistema responde con error `404` y un mensaje en español

### Requirement: Página de solicitud de test drive

El sistema SHALL ofrecer una página de solicitud de test drive en la ruta
`/catalogo/:id/test-drive` que muestre un formulario con la fecha y la franja
horaria para la unidad, consumiendo el endpoint `POST /api/test-drives`. La
página SHALL mostrar un mensaje en español cuando el turno se crea con éxito y
un mensaje de error en español cuando la solicitud falla (incluido el caso de
turno ya ocupado).

#### Scenario: Acceso a la página de solicitud

- **WHEN** un cliente autenticado navega a `/catalogo/:id/test-drive` de un
  vehículo disponible
- **THEN** el sistema muestra el formulario de fecha y franja horaria para ese
  vehículo

#### Scenario: Solicitud con turno ocupado

- **WHEN** un cliente envía el formulario con una unidad, fecha y franja que ya
  tienen un turno activo
- **THEN** el sistema muestra el error `409` en español indicando que el turno
  está ocupado

#### Scenario: Solicitud exitosa

- **WHEN** un cliente envía un formulario válido
- **THEN** el sistema crea el turno y muestra un mensaje de éxito en español

### Requirement: Página de gestión de turnos del vendedor

El sistema SHALL ofrecer una página de gestión de turnos en la ruta
`/vendedor/test-drives` accesible para vendedores y administradores que liste
los turnos agendados, permita filtrar por estado y ejecutar las acciones de
confirmar, cancelar y completar.

#### Scenario: Listado con acciones

- **WHEN** un vendedor navega a `/vendedor/test-drives`
- **THEN** el sistema muestra los turnos agendados con acciones según su estado

### Requirement: Página de mis turnos del cliente

El sistema SHALL ofrecer una página en la ruta `/mis-test-drives` accesible
para clientes que liste sus turnos de test drive y permita cancelar los que
estén en estado `solicitado` o `confirmado`.

#### Scenario: Listado de turnos propios

- **WHEN** un cliente navega a `/mis-test-drives`
- **THEN** el sistema muestra sus turnos con el vehículo y la posibilidad de
  cancelar los activos
