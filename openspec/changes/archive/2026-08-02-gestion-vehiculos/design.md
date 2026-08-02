## Context

El proyecto es un monorepo con backend en Go (Gin + GORM + JWT), frontend en
React (Vite + TypeScript + React Router + TailwindCSS) y PostgreSQL. La entidad
`Vehiculo` (con su relación `Imagen`) ya existe y se auto-migra, y el catálogo
público (CU-03) la lee filtrando `estado = disponible`.

El backend ya tiene las capas `handler → service → repository → GORM`. El
frontend tiene el cliente HTTP centralizado (`services/api.ts`) y páginas
públicas. La autenticación con JWT y roles (CU-01) aún es un stub que deja pasar
todas las peticiones; los middlewares `AutenticacionJWT` y `ExigirRol` ya
existen y se pueden enganchar desde ahora.

Este change (CU-02) agrega la gestión administrativa del stock: alta,
modificación, baja lógica y listado de todos los estados.

## Goals / Non-Goals

**Goals:**

- Exponer un CRUD administrativo completo de vehículos, restringido a rol
  `administrador`.
- Listar todos los estados en la gestión (a diferencia del catálogo público).
- Implementar la baja como **baja lógica** (`estado = dado_de_baja`), sin
  borrar filas ni perder el historial.
- Mantener la separación de capas del backend y las convenciones del repo.
- No alterar el comportamiento del catálogo público existente.

**Non-Goals:**

- Autenticación real con JWT y roles (CU-01): solo se conectan los middlewares
  stub.
- Carga de imágenes como archivos: se gestionan URLs de imagen.
- Búsqueda/filtrado avanzado del catálogo público (CU-04).
- Gestión de reservas/ventas que cambien el estado (CU-08).
- Panel de administración con métricas (CU-09).

## Decisions

### D1: Grupo administrativo propio `/api/admin/vehiculos`

El CRUD administrativo vive en un grupo de rutas separado (`/api/admin/
vehiculos`) protegido con `middleware.AutenticacionJWT` +
`middleware.ExigirRol("administrador")`. Alternativa considerada: reutilizar
`/api/vehiculos` con los mismos paths y métodos HTTP distintos. Se descarta
porque el listado administrativo tiene un contrato diferente (todos los estados,
ficha completa) al público (solo `disponible`, ficha básica); separar los
grupos evita ambigüedad de permisos y de DTOs.

### D2: Baja lógica mediante actualización de `estado`

`DELETE /api/admin/vehiculos/:id` no borra el registro: el service actualiza
`estado` a `dado_de_baja`. Alternativa: `gorm.DeletedAt` (soft delete de GORM) o
baja física. Se prefiere la columna `estado` porque ya es la semántica de
negocio definida en el modelo y no introduce columnas nuevas; además CU-08
reservas/ventas también mutará `estado`. La operación es idempotente: si ya está
`dado_de_baja`, responde `200` sin cambios.

### D3: Reemplazo completo de imágenes en alta y modificación

El formulario envía siempre la lista completa de URLs de imágenes. En `Crear` se
insertan las filas de `Imagen`; en `Actualizar` se borran las existentes y se
insertan las nuevas. Alternativa: agregar/remover por diffs. Se prefiere el
reemplazo por simplicidad y porque el ABM es la única fuente de escritura de la
galería en este momento.

### D4: Validaciones de negocio en el service

El service valida: campos requeridos (`marca`, `modelo`, `anio`, `precio`,
`condicion`), `anio` en un rango razonable, `precio` positivo, `condicion` en
`{nuevo, usado}`, `estado` en `{disponible, reservado, vendido, dado_de_baja}`
y el filtro de estado del listado. El repository no valida; los handlers solo
parsean y responden. Se siguen las reglas del repo (prohibido saltarse capas).

### D5: Un solo service de vehículos ampliado

Se amplía la interfaz `VehiculoService` existente con `ListarParaGestion`,
`ObtenerParaGestion`, `Crear`, `Actualizar` y `DarDeBaja`, en lugar de crear un
service separado. Alternativa: `VehiculoGestionService` aparte. Se prefiere un
solo service porque la entidad es la misma y así se comparte el mapeo de errores
y la validación de paginación. En el handler se agrega un struct dedicado para
gestión (`VehiculoGestionHandler`) con sus DTOs, sin tocar los DTOs públicos.

### D6: Los middlewares de rol se conectan desde ya (stubs)

Aunque CU-01 no esté implementado, las rutas administrativas se registran con
`AutenticacionJWT(configuracion.JWTSecreto)` y `ExigirRol("administrador")`.
Como hoy los middlewares son pasantes, no bloquean el desarrollo; cuando CU-01
los implemente, la protección queda activa sin tocar el router de vehículos.

## Risks / Trade-offs

- [Rutas admin abiertas hasta CU-01] → Mitigación: los middlewares stub son
  temporales y ya están enganchados en las rutas; la protección real llega con
  CU-01 sin cambios en este módulo.
- [Baja lógica no libera datos de historial] → Mitigación: es una decisión
  deliberada para preservar consultas y reservas futuras; el catálogo público ya
  excluye `dado_de_baja`.
- [Reemplazo de imágenes puede perder URLs no enviadas] → Mitigación: el
  formulario siempre envía la lista completa; el detalle de gestión expone las
  imágenes actuales para que el admin las vea al editar.
- [Ficha de gestión con más campos que la pública] → Mitigación: se reutiliza el
  modelo completo en los DTOs de gestión; el catálogo público mantiene su DTO
  resumido.
- [Estado editado libremente por el admin] → Mitigación: el admin puede fijar
  `estado` al alta y en la edición (necesario hasta que CU-08/otras mutaciones
  existan); se valida que sea uno de los cuatro valores conocidos.

## Migration Plan

- Cambio aditivo: no requiere migración de datos ni de esquema. Las tablas
  `vehiculos` e `imagenes` ya están migradas. El despliegue se hace levantando
  backend y frontend sin pasos extra.

## Open Questions

- Si CU-08 (reservas/ventas) restringe qué estados puede modificar el admin
  (p. ej. no dar de baja un vehículo reservado), se ajustará la regla en el
  service en ese change.
