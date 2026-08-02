## Why

El catálogo público (CU-03) solo permite ver todas las unidades disponibles
paginas. El visitante no puede encontrar un vehículo puntual cuando el stock
crece: necesita buscar por texto y acotar resultados con filtros combinables
para comparar unidades por precio, año, tipo, combustible o condición.

## What Changes

- Ampliar `GET /api/vehiculos` (catálogo público) con **búsqueda por texto
  libre** (marca y modelo) y **filtros combinables**: marca, modelo, rango de
  años (`anio_min`/`anio_max`), rango de precio (`precio_min`/`precio_max`),
  tipo, combustible y condición (nuevo/usado).
- Agregar **ordenamiento** por precio o año, ascendente o descendente.
- Agregar el campo nuevo `tipo` al modelo `Vehiculo` (sedán, SUV, hatchback,
  pick-up, etc.), con migración automática. Esto habilita el filtro por tipo y
  la ficha del vehículo. **BREAKING**: la API de alta/modificación de vehículos
  y la ficha técnica pasan a incluir `tipo`.
- Incorporar `tipo` al ABM administrativo (CU-02): alta, modificación y
  formulario de vehículos.
- Rediseñar la página pública `/catalogo` con un panel de búsqueda y filtros
  (incluye selección de tipo y ordenamiento) que actualiza el listado paginado.

## Capabilities

### New Capabilities

- `busqueda-filtrado`: búsqueda por texto libre, filtros combinables
  (marca, modelo, años, precio, tipo, combustible, condición) y ordenamiento
  por precio o año en el catálogo público.

### Modified Capabilities

- `catalogo-vehiculos`: el endpoint público `GET /api/vehiculos` y la página
  `/catalogo` pasan a soportar búsqueda, filtros combinables y ordenamiento.
  La ficha técnica del vehículo incluye el nuevo campo `tipo`.
- `gestion-vehiculos`: el alta y la modificación de vehículos incorporan el
  campo nuevo `tipo` en la ficha técnica y en los formularios administrativos.

## Impact

- **Backend**: `models/vehiculo.go` (nuevo campo `tipo`), `services/vehiculos.go`
  (validación de filtros y orden), `repositories/vehiculos.go` (construcción de
  la consulta dinámica GORM), `handlers/vehiculos.go` (parseo de query params),
  `handlers/vehiculos_gestion.go` (DTO `tipo`), `database/database.go`
  (auto-migración del nuevo campo). No se agregan dependencias.
- **Frontend**: `pages/Catalogo.tsx` (panel de búsqueda/filtros/orden),
  `services/api.ts` (cliente con filtros), `types/vehiculo.ts` (campo `tipo`),
  `pages/FormularioVehiculo.tsx` (select de tipo), `pages/DetalleVehiculo.tsx`
  (muestra el tipo).
- **Datos**: la auto-migración de GORM agrega la columna `tipo` a `vehiculos`.
  Se sugiere sembrar tipos en el seed para que el filtro tenga datos de prueba.
