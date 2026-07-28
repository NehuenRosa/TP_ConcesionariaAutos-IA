# Reservas y Ventas

## Requirements

### Requirement: Crear reserva
- **WHEN** un cliente autenticado envía POST /api/reservations con vehicle_id
- **THEN** el sistema cambia el vehículo a "reservado" y retorna 201 con la reserva en estado "activa"

- **WHEN** un cliente envía POST /api/reservations con vehículo no disponible
- **THEN** el sistema retorna 400

- **WHEN** un cliente envía POST /api/reservations sin token
- **THEN** el sistema retorna 401

### Requirement: Confirmar venta
- **WHEN** un vendedor/admin envía POST /api/reservations/:id/confirm
- **THEN** el sistema cambia la reserva a "confirmada" y el vehículo a "vendido"

### Requirement: Cancelar reserva
- **WHEN** un vendedor/admin envía POST /api/reservations/:id/cancel
- **THEN** el sistema cambia la reserva a "cancelada" y el vehículo vuelve a "disponible"

### Requirement: Listar reservas
- **WHEN** un cliente autenticado envía GET /api/reservations/mine
- **THEN** el sistema retorna 200 con sus reservas

- **WHEN** un vendedor/admin envía GET /api/reservations
- **THEN** el sistema retorna 200 con todas las reservas
