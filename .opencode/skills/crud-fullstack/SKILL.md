---
name: crud-fullstack
description: Crear un CRUD completo (alta, baja, modificación, listado y detalle) de una entidad de punta a punta. Usar cuando el usuario pida "armar un CRUD", "hacer un ABM", "CRUD de <algo>" o implementar la gestión completa de una entidad en este monorepo (backend Go + frontend React), integrado con el flujo OpenSpec.
license: MIT
metadata:
  author: proyecto-concesionaria
  version: "1.0"
---

# CRUD Full-stack (backend Go + frontend React + OpenSpec)

Implementa la gestión completa de una entidad respetando las convenciones del
repo (ver `AGENTS.md`): todo en español, backend en capas
`handler → service → repository → GORM`, frontend con tipos, cliente HTTP
centralizado y páginas, y planificación previa con OpenSpec.

## Pasos

### 1. Relevar la entidad

Antes de tocar código, confirmar con el usuario:

- Nombre de la entidad (en español, singular: `vehiculo`, `cliente`, `turno`).
- Campos y sus tipos. Guía de tipos según la ficha técnica del repo:
  - Textos: `string` (`gorm:"not null"`).
  - Números: `int`, `float64`.
  - Fechas: `time.Time`.
  - Relaciones (uno-a-muchos): slice del modelo hijo (ej. `Imagenes []Imagen`).
  - Estado/condición: `string` con constantes en español (como `EstadoDisponible`).
- Qué operaciones se necesitan. El CRUD típico es:
  - Listado paginado `GET /api/<plural>`
  - Detalle `GET /api/<plural>/:id`
  - Creación `POST /api/<plural>`
  - Actualización `PUT /api/<plural>/:id`
  - Baja `DELETE /api/<plural>/:id`
- ¿Alguna regla de negocio especial? (estados que limitan bajas, unicidad, etc.)

### 2. Planificar con OpenSpec

Seguir el ciclo del repo (una propuesta por entidad):

```powershell
openspec.cmd new change <nombre-en-kebab-case>
openspec.cmd status --change <nombre> --json
```

Completar los 4 artefactos en `openspec/changes/<nombre>/`:
- `proposal.md`: qué y por qué.
- `design.md`: decisiones (D1, D2…): nombres de endpoints, reglas de negocio,
  DTOs de respuesta, si el CRUD es público o autenticado.
- `specs/<capacidad>/spec.md`: requirements con escenarios `WHEN/THEN`
  (mirar `openspec/specs/catalogo-vehiculos/spec.md` como referencia).
- `tasks.md`: tareas numeradas por capa (backend, frontend, verificación).

Validar:

```powershell
openspec.cmd validate --all --strict
```

> Los cambios en `.opencode/` y `openspec/` ya están versionados en el repo.

### 3. Implementar backend

#### 3.1 Modelo (`backend/internal/models/<entidad>.go`)

- Tipo en PascalCase, archivo en minúsculas. `TableName()` en español.
- Tags `gorm` para el esquema y `json` en camelCase (marcar `CreatedAt`,
  `UpdatedAt` y las FK internas con `json:"-"`).
- Constantes de estado/condición en español cuando aplique.

#### 3.2 Repository (`backend/internal/repositories/<plural>.go`)

- Interfaz + implementación sobre `*gorm.DB`. Constructor
  `Nuevo<Entidad>Repository(base)`.
- Métodos: `Listar(ctx, estado, pagina, tamano)` (con `Count` del total y
  `Offset/Limit`), `ObtenerPorID(ctx, id)` (con `Preload` de relaciones),
  `Crear`, `Actualizar`, `Eliminar`.
- **Sin lógica de negocio**: validaciones y reglas viven en el service.

#### 3.3 Service (`backend/internal/services/<plural>.go`)

- Interfaz + implementación. Constructor `Nuevo<Entidad>Service(repositorio)`.
- Valida paginación y reglas de negocio; retorna `error` descriptivo en español
  (ej. `ErrXNoEncontrado`, `ErrYInvalido`). Mapea `gorm.ErrRecordNotFound`.

#### 3.4 Handler (`backend/internal/handlers/<plural>.go`)

- Struct `XHandler` con el servicio inyectado. Constructor
  `NuevoXHandler(servicio)`.
- DTOs de respuesta en camelCase (`XResumen` para listado, `X` para detalle).
- `parsearPaginacion` para el listado. Respuestas JSON con códigos correctos:
  `400` input inválido, `404` no encontrado, `409` conflicto de negocio,
  `500` error interno, mensajes siempre en español.

#### 3.5 Router (`backend/internal/router/router.go`)

- Instanciar repo → service → handler al inicio de `Nuevo`.
- Registrar las rutas del grupo. Para CRUD administrativo usar
  `middleware.AutenticacionJWT` + `middleware.ExigirRol("administrador")`;
  para CRUD público no.

#### 3.6 Migración

- Agregar el modelo a `AutoMigrar` en
  `backend/internal/database/database.go`.

### 4. Implementar frontend

#### 4.1 Tipos (`frontend/src/types/<entidad>.ts`)

- Tipos `Entidad`, `ResumenEntidad` (ficha del listado) y `PaginaEntidades`
  (`{ datos, pagina, tamano, total }`), en camelCase.

#### 4.2 Cliente HTTP (`frontend/src/services/api.ts`)

- Agregar métodos al objeto `api` usando la función `peticion` existente:
  `listarEntidades(pagina, tamano)`, `obtenerEntidad(id)`, `crearEntidad`,
  `actualizarEntidad(id, datos)`, `eliminarEntidad(id)`.

#### 4.3 Páginas (`frontend/src/pages/`)

- `Listado<Entidad>.tsx`: tabla o tarjetas, paginación, estados de carga /
  vacío / error con mensajes en español, acciones editar/eliminar.
- `Formulario<Entidad>.tsx`: alta y edición (reutilizando el mismo
  componente). Enviar a la API y navegar de vuelta al listado.
- Reutilizar patrones de `Catalogo.tsx` y `DetalleVehiculo.tsx` (uso de
  `api`, `ErrorApi`, `useEffect` con bandera `cancelado`).

#### 4.4 Rutas (`frontend/src/routes/Rutas.tsx`)

- Registrar las nuevas rutas dentro del `LayoutBase`.

### 5. Verificar

1. Backend: `cd backend && go build ./...` y `go vet ./...` sin errores.
2. Frontend: `cd frontend && npm run build` sin errores.
   (En PowerShell usar `npm.cmd`.)
3. De punta a punta con `docker compose up -d --build`:
   - Probar listado, detalle, creación, actualización y baja con
     `Invoke-WebRequest` contra `http://localhost:8080/api/<plural>`.
   - Verificar que el frontend responde en `http://localhost:5173`.
   - Si se siembran datos con acentos, NO pipear SQL por PowerShell (corrompe
     UTF-8); usar `docker cp` + `psql -f`.

### 6. Cerrar el change

- Marcar todas las tareas en `tasks.md` y `openspec.cmd validate --all --strict`.
- Al terminar: sincronizar specs (`openspec/specs/<capacidad>/spec.md`) y
  archivar con `openspec.cmd archive <nombre>`.

## Reglas de negocio del repo a respetar

- Los listados públicos exponen solo lo disponible; el detalle devuelve `404`
  para entidades no disponibles (no `403`).
- Prohibido saltarse capas: handlers no tocan GORM, repositories no validan
  negocio, services no escriben respuestas HTTP.
- Todo identificador, mensaje y texto de UI en español.
