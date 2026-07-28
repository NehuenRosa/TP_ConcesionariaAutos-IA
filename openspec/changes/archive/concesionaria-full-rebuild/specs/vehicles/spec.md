# Gestión de Vehículos

## Requirements

### Requirement: Catálogo público
- **WHEN** cualquier usuario envía GET /api/vehicles sin filtros
- **THEN** el sistema retorna 200 con lista paginada de vehículos con estado "disponible"

- **WHEN** cualquier usuario envía GET /api/vehicles con filtros (brand, fuel, price_from, price_to, year_from, year_to, condition, vehicle_type, search, sort_by, sort_order)
- **THEN** el sistema retorna 200 con resultados filtrados y paginados

- **WHEN** cualquier usuario envía GET /api/vehicles/brands
- **THEN** el sistema retorna 200 con array de marcas de vehículos disponibles

- **WHEN** cualquier usuario envía GET /api/vehicles/:id con ID existente
- **THEN** el sistema retorna 200 con detalle completo del vehículo

- **WHEN** cualquier usuario envía GET /api/vehicles/:id con ID inexistente
- **THEN** el sistema retorna 404

### Requirement: CRUD de administrador
- **WHEN** un administrador envía POST /api/vehicles con datos válidos
- **THEN** el sistema retorna 201 con el vehículo creado (status default: disponible)

- **WHEN** un administrador envía POST /api/vehicles sin token
- **THEN** el sistema retorna 401

- **WHEN** un usuario no-admin envía POST /api/vehicles
- **THEN** el sistema retorna 403

- **WHEN** un administrador envía POST /api/vehicles con datos inválidos (brand vacío, price negativo)
- **THEN** el sistema retorna 400

- **WHEN** un administrador envía PUT /api/vehicles/:id con datos válidos
- **THEN** el sistema retorna 200 con el vehículo actualizado

- **WHEN** un administrador envía PUT /api/vehicles/:id con ID inexistente
- **THEN** el sistema retorna 404

- **WHEN** un administrador envía DELETE /api/vehicles/:id con ID existente
- **THEN** el sistema retorna 200 con mensaje de confirmación

- **WHEN** un administrador envía DELETE /api/vehicles/:id con ID inexistente
- **THEN** el sistema retorna 404
