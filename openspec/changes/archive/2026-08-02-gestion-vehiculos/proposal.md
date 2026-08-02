## Why

El sistema necesita que el administrador pueda gestionar el stock de vehículos
de forma integral: dar de alta unidades nuevas, corregir la ficha técnica de las
existentes y dar de baja las que ya no se comercializan. Hoy la entidad
`Vehiculo` solo se lee desde el catálogo público (CU-03); sin un ABM
administrativo no hay forma de cargar, mantener ni retirar unidades del stock,
y el catálogo no tendría datos que mostrar.

## What Changes

- Se agregan endpoints administrativos de vehículos bajo `/api/admin/vehiculos`,
  protegidos con autenticación JWT y rol `administrador`:
  - `GET /api/admin/vehiculos`: listado paginado de **todos** los estados, con
    filtro opcional por estado (a diferencia del catálogo público que solo
    muestra `disponible`).
  - `GET /api/admin/vehiculos/:id`: detalle completo de cualquier vehículo,
    incluidos los no disponibles.
  - `POST /api/admin/vehiculos`: alta de vehículo con ficha técnica e imágenes.
  - `PUT /api/admin/vehiculos/:id`: modificación de ficha técnica, estado e
    imágenes.
  - `DELETE /api/admin/vehiculos/:id`: baja lógica (estado → `dado_de_baja`),
    sin borrar el registro para preservar el historial.
- Se validan en el service los campos de la ficha técnica y los valores de
  `condicion` y `estado`.
- Se agregan las páginas de administración del frontend: listado con filtro por
  estado, formulario de alta/edición y acción de baja lógica.
- El catálogo público (CU-03) no cambia: sigue exponiendo solo `disponible`.
- La autenticación usa los middlewares existentes (stubs hasta CU-01); el
  control real de roles llega con CU-01.

## Capabilities

### New Capabilities
- `gestion-vehiculos`: ABM administrativo de vehículos (listar todos los
  estados, detalle, alta, modificación y baja lógica) restringido a usuarios con
  rol `administrador`, con ficha técnica e imágenes.

### Modified Capabilities
<!-- Sin cambios de requisitos a nivel spec: el catálogo público (catalogo-vehiculos) no altera su comportamiento. -->

## Impact

- **Backend**: se extienden repository, service y handlers de vehículos con las
  operaciones de escritura y el listado de gestión; se registran rutas
  administrativas en el router con middlewares de autenticación y rol. La
  auto-migración no cambia (las tablas `vehiculos` e `imagenes` ya existen).
- **API**: cinco endpoints nuevos bajo `/api/admin/vehiculos`; los endpoints
  públicos de lectura existentes no se modifican.
- **Frontend**: se agregan tipos de entrada/salida, métodos en el cliente HTTP
  centralizado y las páginas de gestión de vehículos en `/admin/vehiculos`.
- **Base de datos**: sin cambios de esquema; la baja se implementa como
  actualización de la columna `estado` a `dado_de_baja`.
