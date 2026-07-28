# Panel de Administración

## Requirements

### Requirement: Dashboard
- **WHEN** un administrador autenticado envía GET /api/admin/dashboard
- **THEN** el sistema retorna 200 con métricas: total_vehiculos, disponibles, reservados, vendidos, consultas_pendientes, test_drives_pendientes, reservas_activas

- **WHEN** un usuario no-admin envía GET /api/admin/dashboard
- **THEN** el sistema retorna 403

- **WHEN** un usuario no autenticado envía GET /api/admin/dashboard
- **THEN** el sistema retorna 401

### Requirement: Interfaz de dashboard
- **WHEN** un administrador accede a /admin/dashboard
- **THEN** el frontend muestra tarjetas con métricas principales, gráfico de barras (vehículos por estado) y gráfico de torta (distribución)
