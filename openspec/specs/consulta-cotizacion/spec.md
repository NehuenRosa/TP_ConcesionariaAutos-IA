# consulta-cotizacion Specification

## Purpose
TBD - created by archiving change cu05-consulta-cotizacion. Update Purpose after archive.
## Requirements
### Requirement: Crear consulta desde detalle de vehículo

El sistema SHALL permitir a un cliente autenticado crear una consulta asociada
a un vehículo específico desde la página de detalle del vehículo. La consulta
se crea con el estado `pendiente` y el primer mensaje contenido en la solicitud.

#### Scenario: Cliente crea consulta válida

- **WHEN** un cliente autenticado envía `POST /api/consultas` con `vehiculoId`
  y `mensaje`
- **THEN** el sistema crea la consulta con estado `pendiente`, registra el
  primer mensaje y retorna la consulta creada

#### Scenario: Cliente no autenticado intenta crear consulta

- **WHEN** un usuario no autenticado envía `POST /api/consultas`
- **THEN** el sistema responde con error `401`

#### Scenario: Vehículo inexistente o no disponible

- **WHEN** un cliente envía `POST /api/consultas` con un `vehiculoId` que no
  existe o no está en estado `disponible`
- **THEN** el sistema responde con error `404`

#### Scenario: Mensaje vacío

- **WHEN** un cliente envía `POST /api/consultas` con `mensaje` vacío
- **THEN** el sistema responde con error `400`

### Requirement: Bandeja de entrada del vendedor

El sistema SHALL exponer un endpoint `GET /api/consultas/bandeja` que devuelva
las consultas asignadas al vendedor autenticado, ordenadas por fecha de último
mensaje descendente. Cada consulta incluye un preview del último mensaje y
un indicador de mensajes nuevos.

#### Scenario: Vendedor con consultas

- **WHEN** un vendedor autenticado solicita su bandeja de entrada
- **THEN** el sistema responde con las consultas donde es el vendedor asignado,
  incluyendo datos del vehículo, cliente, estado, último mensaje y si tiene
  mensajes nuevos

#### Scenario: Vendedor sin consultas

- **WHEN** un vendedor autenticado solicita su bandeja de entrada y no tiene
  consultas asignadas
- **THEN** el sistema responde con una lista vacía

### Requirement: Tomar consulta pendiente

El sistema SHALL permitir a un vendedor autenticado tomar una consulta en
estado `pendiente`, asignándola a su usuario y cambiando el estado a
`en conversación`.

#### Scenario: Vendedor toma consulta pendiente

- **WHEN** un vendedor autenticado envía `PUT /api/consultas/:id/tomar` y la
  consulta está en estado `pendiente`
- **THEN** el sistema asigna el vendedor, cambia el estado a `en conversación`
  y retorna la consulta actualizada

#### Scenario: Consulta ya fue tomada

- **WHEN** un vendedor envía `PUT /api/consultas/:id/tomar` y la consulta no
  está en estado `pendiente`
- **THEN** el sistema responde con error `409`

#### Scenario: Consulta inexistente

- **WHEN** un vendedor envía `PUT /api/consultas/:id/tomar` con un ID inexistente
- **THEN** el sistema responde con error `404`

### Requirement: Cerrar consulta

El sistema SHALL permitir al vendedor asignado cerrar una consulta en estado
`en conversación`, cambiando el estado a `cerrada`.

#### Scenario: Vendedor cierra su consulta

- **WHEN** el vendedor asignado envía `PUT /api/consultas/:id/cerrar`
- **THEN** el sistema cambia el estado a `cerrada` y retorna la consulta

#### Scenario: Vendedor intenta cerrar consulta ajena

- **WHEN** un vendedor envía `PUT /api/consultas/:id/cerrar` y no es el
  vendedor asignado
- **THEN** el sistema responde con error `403`

#### Scenario: Consulta ya cerrada

- **WHEN** un vendedor envía `PUT /api/consultas/:id/cerrar` y la consulta
  ya está cerrada
- **THEN** el sistema responde con error `409`

### Requirement: Eliminar consulta cerrada

El sistema SHALL permitir al vendedor asignado eliminar una consulta en estado
`cerrada`, eliminándola permanentemente de la base de datos.

#### Scenario: Vendedor elimina consulta cerrada

- **WHEN** el vendedor asignado envía `DELETE /api/consultas/:id` y la consulta
  está cerrada
- **THEN** el sistema elimina la consulta y sus mensajes, retorna `204`

#### Scenario: Vendedor intenta eliminar consulta no cerrada

- **WHEN** un vendedor envía `DELETE /api/consultas/:id` y la consulta no
  está cerrada
- **THEN** el sistema responde con error `409`

### Requirement: Listar consultas del cliente

El sistema SHALL exponer un endpoint `GET /api/consultas/mis-consultas` que
devuelva las consultas del cliente autenticado, ordenadas por fecha de último
mensaje descendente.

#### Scenario: Cliente con consultas

- **WHEN** un cliente autenticado solicita sus consultas
- **THEN** el sistema responde con las consultas donde es el cliente,
  incluyendo datos del vehículo, vendedor (si existe), estado y último mensaje

#### Scenario: Cliente sin consultas

- **WHEN** un cliente autenticado solicita sus consultas y no tiene ninguna
- **THEN** el sistema responde con una lista vacía

### Requirement: Enviar mensaje en consulta

El sistema SHALL permitir a los participantes de una consulta (cliente o
vendedor asignado) enviar mensajes mientras la consulta no esté cerrada.

#### Scenario: Cliente envía mensaje

- **WHEN** el cliente de la consulta envía `POST /api/consultas/:id/mensajes`
  con `contenido`
- **THEN** el sistema crea el mensaje, lo asocia a la consulta y retorna el
  mensaje creado

#### Scenario: Vendedor envía mensaje

- **WHEN** el vendedor asignado envía `POST /api/consultas/:id/mensajes`
  con `contenido`
- **THEN** el sistema crea el mensaje y lo retorna

#### Scenario: Usuario no participante envía mensaje

- **WHEN** un usuario que no es cliente ni vendedor de la consulta envía
  `POST /api/consultas/:id/mensajes`
- **THEN** el sistema responde con error `403`

#### Scenario: Mensaje en consulta cerrada

- **WHEN** se envía un mensaje a una consulta cerrada
- **THEN** el sistema responde con error `409`

### Requirement: Obtener mensajes de una consulta

El sistema SHALL exponer un endpoint `GET /api/consultas/:id/mensajes` que
devuelva todos los mensajes de una consulta para sus participantes.

#### Scenario: Participante obtiene mensajes

- **WHEN** el cliente o vendedor asignado solicita los mensajes de una consulta
- **THEN** el sistema responde con el listado de mensajes ordenados por fecha

#### Scenario: Usuario no participante solicita mensajes

- **WHEN** un usuario que no es participante solicita los mensajes
- **THEN** el sistema responde con error `403`

### Requirement: Obtener mensajes nuevos (polling)

El sistema SHALL exponer un endpoint `GET /api/consultas/:id/mensajes/nuevos?desde=<timestamp>`
que devuelva los mensajes recibidos después del timestamp indicado, marcándolos
como leídos.

#### Scenario: Hay mensajes nuevos

- **WHEN** un participante solicita mensajes nuevos desde un timestamp
- **THEN** el sistema responde con los mensajes posteriores al timestamp,
  marcados como leídos

#### Scenario: No hay mensajes nuevos

- **WHEN** un participante solicita mensajes nuevos desde un timestamp y no
  hay mensajes posteriores
- **THEN** el sistema responde con una lista vacía

### Requirement: Página de detalle con botón consultar

El sistema SHALL mostrar un botón "Consultar" en la página de detalle del
vehículo (`/catalogo/:id`) visible solo para clientes autenticados. Al hacer
clic, se muestra un formulario para escribir el mensaje inicial.

#### Scenario: Cliente autenticado en detalle de vehículo

- **WHEN** un cliente autenticado accede al detalle de un vehículo disponible
- **THEN** se muestra el botón "Consultar"

#### Scenario: Visitante o no autenticado en detalle

- **WHEN** un visitante o usuario no autenticado accede al detalle
- **THEN** NO se muestra el botón "Consultar"

### Requirement: Bandeja de entrada del vendedor (página)

El sistema SHALL ofrecer una página en `/vendedor/bandeja` que muestre las
consultas del vendedor en formato de tarjetas con preview del último mensaje.

#### Scenario: Vendedor accede a su bandeja

- **WHEN** un vendedor autenticado accede a `/vendedor/bandeja`
- **THEN** se muestran las tarjetas de sus consultas con vehicle info,
  cliente, estado, preview del último mensaje y badge de mensajes nuevos

### Requirement: Vista de chat del cliente

El sistema SHALL ofrecer una página en `/mis-consultas` que muestre una vista
tipo chat: lista de consultas a la izquierda, conversación a la derecha.

#### Scenario: Cliente accede a sus consultas

- **WHEN** un cliente autenticado accede a `/mis-consultas`
- **THEN** se muestra la lista de sus consultas y al seleccionar una se
  muestra la conversación completa

### Requirement: Notificaciones con punto rojo

El sistema SHALL mostrar un punto rojo en las tarjetas/listas de consultas
cuando existan mensajes nuevos no leídos del otro participante. El conteo se
basa en el campo `mensajesNuevos` devuelto por el backend en
`GET /api/consultas/mis-consultas` y `GET /api/consultas/bandeja`, que cuenta
solo los mensajes del otro participante con `leido = false`.

#### Scenario: Mensaje nuevo recibido

- **WHEN** un participante recibe un mensaje nuevo del otro participante
- **THEN** el campo `mensajesNuevos` de la consulta es mayor a 0 y se muestra
  un punto rojo en la tarjeta/lista correspondiente

#### Scenario: Mensaje nuevo leído

- **WHEN** el usuario abre la conversación
- **THEN** el sistema marca los mensajes del otro participante como leídos
  (`PUT /api/consultas/:id/mensajes/leidos`), `mensajesNuevos` pasa a 0 y se
  oculta el punto rojo

### Requirement: Contador de notificaciones para el navbar

El sistema SHALL exponer un endpoint liviano `GET /api/notificaciones/contador`
que devuelva la cantidad total de mensajes no leídos del usuario autenticado,
sin cargar el listado completo de consultas.

#### Scenario: Usuario autenticado consulta el contador

- **WHEN** un usuario autenticado (cliente o vendedor) solicita
  `GET /api/notificaciones/contador`
- **THEN** el sistema responde con el total de mensajes no leídos en todas sus
  consultas activas

#### Scenario: Usuario sin mensajes no leídos

- **WHEN** el usuario no tiene mensajes no leídos
- **THEN** el sistema responde con `{"contador": 0}`

#### Scenario: Navbar con notificaciones

- **WHEN** el navbar consulta el contador cada 3 segundos o recibe el evento
  `mensajes-leidos` disparado por el chat
- **THEN** se muestra el punto rojo si el contador es mayor a 0 y se oculta al
  llegar a 0

