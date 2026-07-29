# Consultas al Vendedor

## Requirements

### Requirement: Crear consulta
- **WHEN** un cliente autenticado envía POST /api/consultations con vehicle_id y message
- **THEN** el sistema retorna 201 con la consulta creada en estado "pendiente"

- **WHEN** un cliente envía POST /api/consultations sin token
- **THEN** el sistema retorna 401

- **WHEN** un cliente envía POST /api/consultations con vehicle_id de vehículo no disponible
- **THEN** el sistema retorna 400

### Requirement: Listar consultas
- **WHEN** un cliente autenticado envía GET /api/consultations/mine
- **THEN** el sistema retorna 200 con sus consultas (ordenadas: no leídas primero, luego pendientes, luego por fecha descendente)

- **WHEN** un vendedor/admin envía GET /api/consultations
- **THEN** el sistema retorna 200 con todas las consultas (ordenadas: pendientes primero, luego con mensajes no leídos, luego por fecha descendente)

### Requirement: Gestionar consultas
- **WHEN** un vendedor/admin envía PATCH /api/consultations/:id/status con un estado válido ("en_conversacion" o "cerrada")
- **THEN** el sistema retorna 200 con estado actualizado y asigna el vendedor como responsable

- **WHEN** un usuario agrega una respuesta a una consulta pendiente
- **THEN** el sistema cambia automáticamente el estado a "en_conversacion"

- **WHEN** un vendedor/admin cambia el estado a "cerrada"
- **THEN** la consulta se marca como finalizada y ya no se permiten más respuestas

### Requirement: Respuestas
- **WHEN** cualquier participante envía POST /api/consultations/:id/responses con message
- **THEN** el sistema retorna 200 con la consulta incluyendo la nueva respuesta

- **WHEN** el cliente responde en una consulta "en_conversacion"
- **THEN** el sistema marca has_unread_messages=true para notificar al vendedor

- **WHEN** el vendedor responde en una consulta "en_conversacion"
- **THEN** el sistema marca has_unread_for_client=true para notificar al cliente

### Requirement: Notificaciones en navbar
- **WHEN** un vendedor está autenticado y hay consultas pendientes o con mensajes no leídos
- **THEN** el navbar muestra un badge rojo con el total en el botón "Bandeja"

- **WHEN** un cliente está autenticado y tiene respuestas del vendedor sin leer
- **THEN** el navbar muestra un badge rojo con el total en el botón "Mis Consultas"

- **WHEN** un vendedor abre el detalle de una consulta (GET /:id)
- **THEN** el sistema marca has_unread_messages=false para esa consulta

- **WHEN** un cliente abre el detalle de una consulta (GET /:id)
- **THEN** el sistema marca has_unread_for_client=false para esa consulta

### Requirement: Indicador de no leídos
- **WHEN** un vendedor ve la lista de consultas
- **THEN** las consultas con has_unread_messages=true muestran un punto rojo en el avatar

- **WHEN** un cliente ve la lista de sus consultas
- **THEN** las consultas con has_unread_for_client=true muestran un punto rojo en el avatar

### Requirement: Eliminar consultas
- **WHEN** el cliente dueño de una consulta envía DELETE /api/consultations/:id
- **THEN** el sistema retorna 200 y elimina la consulta y sus respuestas

- **WHEN** un vendedor/admin envía DELETE /api/consultations/:id
- **THEN** el sistema retorna 200 y elimina cualquier consulta

- **WHEN** un cliente intenta eliminar una consulta de otro cliente
- **THEN** el sistema retorna 403

### Requirement: Contadores de notificación
- **WHEN** un vendedor/admin autenticado envía GET /api/consultations/notifications/count
- **THEN** el sistema retorna { pending, unread, total } con pendientes + no leídas para el vendedor

- **WHEN** un cliente autenticado envía GET /api/consultations/notifications/count
- **THEN** el sistema retorna { unread, total } con consultas con respuesta del vendedor sin leer

- **WHEN** un usuario no autenticado envía GET /api/consultations/notifications/count
- **THEN** el sistema retorna 401

### Requirement: Visualización de estados
- **WHEN** una consulta está en estado "pendiente"
- **THEN** se muestra un badge con círculo ámbar, fondo amarillo claro, borde amarillo y texto "Pendiente"

- **WHEN** una consulta está en estado "en_conversacion"
- **THEN** se muestra un badge con círculo azul, fondo azul claro, borde azul y texto "En conversación"

- **WHEN** una consulta está en estado "cerrada"
- **THEN** se muestra un badge con círculo rojo, fondo rojo claro, borde rojo y texto "Cerrada"
