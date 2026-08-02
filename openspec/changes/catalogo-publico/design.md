## Context

El proyecto es un monorepo con backend en Go (Gin + GORM + JWT), frontend en
React (Vite + TypeScript + React Router + TailwindCSS) y PostgreSQL como base de
datos. Este change agrega la primera funcionalidad de solo lectura del sistema:
el catálogo público de vehículos disponibles (CU-03).

El backend ya cuenta con un esqueleto en capas: `handler → service →
repository → GORM`, con auto-migración de la entidad `Vehicle`. El frontend
tiene un layout base y un cliente HTTP centralizado en `services/api.ts`. No
existen endpoints públicos de catálogo todavía.

## Goals / Non-Goals

**Goals:**

- Exponer dos endpoints públicos de solo lectura: listado paginado y detalle.
- Restringir el catálogo a vehículos con estado `disponible`.
- Mostrar el catálogo y el detalle en páginas públicas del frontend.
- Respetar la separación de capas del backend.

**Non-Goals:**

- Búsqueda y filtrado avanzado (CU-04).
- Autenticación ni autorización (CU-01): el catálogo es público por diseño.
- ABM de vehículos (CU-02): la entidad `Vehicle` ya existe como esqueleto.
- Gestión de imágenes como sistema de archivos: se usan URLs de imagen.

## Decisions

### D1: Paginación por `página` y `tamaño` con metadatos

Se usa paginación clásica con parámetros `página` (1-indexada) y `tamaño`
(limit), con respuesta `{ datos, pagina, tamano, total }`. Alternativa
considerada: paginación por cursor. Se descarta por simplicidad: el catálogo
no es un feed y los requisitos no exigen consistencia bajo escrituras
concurrentes.

### D2: Filtro de estado en la capa de service

La condición `estado = disponible` se aplica en el service de vehículos (lógica
de negocio), nunca en el repository. El repository recibe la condición como
parámetro, manteniendo las capas desacopladas. Alternativa: filtrar en el
repository con un método dedicado `listarDisponibles`. Se prefiere parametrizar
para que CU-04 (filtros) pueda reutilizarlo sin nuevos métodos.

### D3: El detalle devuelve `404` para vehículos no disponibles

Un vehículo reservado, vendido o dado de baja no se expone públicamente. Devolver
`404` (en lugar de `403`) evita revelar la existencia de unidades que no están
en venta. Alternativa considerada: `404` solo para inexistente y `403` para no
disponible; se descarta por filtración de información.

### D4: El frontend consume los endpoints con el cliente HTTP centralizado

Las páginas de catálogo y detalle usan `services/api.ts` con la URL base desde
`VITE_API_URL`. No se crea lógica de red en los componentes. Los tipos
compartidos de vehículo viven en `types/`.

### D5: Rutas públicas en el router del frontend

`/catalogo` y `/catalogo/:id` se registran fuera de cualquier layout protegido,
ya que el catálogo es público. El layout base (header/footer) envuelve ambas
páginas.

## Risks / Trade-offs

- [Paginación desactualizada si cambia el stock] → Mitigación: el catálogo es
  de solo lectura para el cliente y las ventas son poco frecuentes; se acepta
  una vista eventualmente consistente.
- [Exponer demasiados campos del vehículo en el listado] → Mitigación: el DTO
  del listado incluye solo la ficha básica; el detalle expone la ficha completa.
- [Names de estado como strings crudos] → Mitigación: los estados se definen
  como constantes en español (`disponible`, `reservado`, `vendido`,
  `dado_de_baja`) para evitar errores de tipeo.
- [Dependencia del esquema `Vehicle` existente] → Mitigación: este change solo
  lee; no altera la migración ni la estructura de la tabla.

## Migration Plan

- No requiere migración de datos: el change es aditivo (endpoints y páginas
  nuevos). El despliegue se hace levantando backend y frontend sin cambios de
  esquema.

## Open Questions

- Formato exacto de la galería de imágenes en el detalle: el modelo `Vehicle`
  usa URLs; se resuelve al implementar CU-02. Para este change basta con
  devolver las URLs tal como están persistidas.
