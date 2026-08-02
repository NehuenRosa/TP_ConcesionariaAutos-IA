## Why

El sistema necesita exhibir el stock de vehículos disponibles a clientes y
visitantes de forma pública. Sin un catálogo, los clientes no pueden descubrir
los vehículos en venta ni acceder a su ficha técnica, y los casos de uso
posteriores (consulta, reserva, test drive) dependen de poder referenciar un
vehículo desde el catálogo.

## What Changes

- Se agrega un endpoint público para listar vehículos **disponibles** con
  paginación (`GET /api/vehiculos`).
- Se agrega un endpoint público para obtener el detalle de un vehículo por ID
  (`GET /api/vehiculos/:id`).
- Se agrega la página de catálogo con el listado paginado de vehículos
  disponibles (ruta `/catalogo`).
- Se agrega la página de detalle de vehículo con ficha técnica (ruta
  `/catalogo/:id`).
- Solo los vehículos con estado `disponible` se muestran en el catálogo público.
- No se implementan búsqueda ni filtros (corresponden a CU-04).

## Capabilities

### New Capabilities
- `catalogo-vehiculos`: exposición pública, paginada y de solo lectura, de los
  vehículos disponibles del concesionario, incluyendo el detalle con ficha
  técnica de cada unidad.

### Modified Capabilities
<!-- Sin capacidades existentes modificadas: es el primer change del proyecto. -->

## Impact

- **Backend**: se agregan handlers, services y repositories de vehículos; se
  exponen rutas públicas de lectura en `/api/vehiculos`. Se usa la entidad
  `Vehicle` (GORM) ya esqueletada, con auto-migración.
- **Frontend**: se agregan páginas `/catalogo` y `/catalogo/:id`, el layout
  base, el cliente HTTP centralizado (`services/api.ts`) y tipos compartidos
  para vehículos.
- **API**: dos endpoints nuevos de solo lectura; sin cambios en los existentes.
- **Base de datos**: sin cambios de esquema; solo se consulta la tabla de
  vehículos.
