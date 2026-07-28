# Test Drives

## Requirements

### Requirement: Solicitar turno
- **WHEN** un cliente autenticado envía POST /api/test-drives con vehicle_id y scheduled_at
- **THEN** el sistema retorna 201 con el turno creado en estado "pendiente"

- **WHEN** un cliente envía POST /api/test-drives con fecha pasada
- **THEN** el sistema retorna 400

- **WHEN** un cliente envía POST /api/test-drives con vehículo no disponible
- **THEN** el sistema retorna 400

- **WHEN** un cliente envía POST /api/test-drives con horario ya ocupado (mismo vehículo + horario)
- **THEN** el sistema retorna 400

- **WHEN** un cliente envía POST /api/test-drives sin token
- **THEN** el sistema retorna 401

### Requirement: Listar turnos
- **WHEN** un cliente autenticado envía GET /api/test-drives/mine
- **THEN** el sistema retorna 200 con sus turnos ordenados por fecha descendente

- **WHEN** un vendedor/admin envía GET /api/test-drives
- **THEN** el sistema retorna 200 con todos los turnos

### Requirement: Gestionar estados
- **WHEN** un vendedor/admin envía PATCH /api/test-drives/:id/status con transición válida
- **THEN** el sistema retorna 200 con estado actualizado

- **WHEN** un vendedor/admin envía PATCH con transición inválida (pendiente→completado)
- **THEN** el sistema retorna 400

Transiciones válidas: pendiente→confirmado, pendiente→cancelado, confirmado→completado, confirmado→cancelado
