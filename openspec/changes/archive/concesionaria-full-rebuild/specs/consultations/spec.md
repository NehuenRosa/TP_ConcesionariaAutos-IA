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
- **THEN** el sistema retorna 200 con sus consultas

- **WHEN** un vendedor/admin envía GET /api/consultations
- **THEN** el sistema retorna 200 con todas las consultas

### Requirement: Gestionar consultas
- **WHEN** un vendedor/admin envía PATCH /api/consultations/:id/status con un estado válido
- **THEN** el sistema retorna 200 con estado actualizado

- **WHEN** un usuario agrega una respuesta a una consulta pendiente
- **THEN** el sistema cambia automáticamente el estado a "en_conversacion"

### Requirement: Respuestas
- **WHEN** cualquier participante envía POST /api/consultations/:id/responses con message
- **THEN** el sistema retorna 200 con la consulta incluyendo la nueva respuesta
