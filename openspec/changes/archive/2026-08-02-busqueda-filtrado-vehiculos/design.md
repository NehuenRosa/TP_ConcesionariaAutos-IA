## Context

El catálogo público (`GET /api/vehiculos`) hoy devuelve todas las unidades con
estado `disponible`, paginadas y sin filtros. El modelo `Vehiculo` no tiene un
campo `tipo` (sedán, SUV, hatchback, pick-up, etc.), por lo que el filtro de
tipo del CU-04 requiere agregarlo de punta a punta: modelo GORM, DTOs de
handlers, servicios, formulario administrativo, ficha pública y seed.

La arquitectura es de capas estrictas: handler (HTTP) → service (negocio) →
repository (GORM) → PostgreSQL. Todo el código, mensajes y textos en español.

## Goals / Non-Goals

**Goals:**
- Ampliar `GET /api/vehiculos` con búsqueda por texto libre, filtros
  combinables y ordenamiento, manteniendo el contrato de respuesta paginada
  actual (`datos`, `pagina`, `tamano`, `total`).
- Agregar el campo `tipo` a la ficha técnica del vehículo en catálogo, detalle
  y ABM administrativo.
- Rediseñar la página `/catalogo` con panel de búsqueda/filtros/orden que
  actualiza el listado paginado.

**Non-Goals:**
- No se cambia la paginación ni el formato de respuesta del catálogo.
- No se implementa guardado de filtros, historial de búsquedas ni filtros sobre
  el listado administrativo (solo se agrega el campo `tipo` a su ABM).
- No se indexa la búsqueda full-text de PostgreSQL; se usa `ILIKE` sobre
  marca/modelo, suficiente para el volumen de un concesionario.

## Decisions

### 1. Campo nuevo `tipo` en el modelo Vehiculo

Se agrega `Tipo string` a `models.Vehiculo` con la misma lógica que `condicion`:
valor libre en backend, sin enum estricto (se valida solo que no esté vacío),
pero con valores sugeridos en el frontend (sedán, SUV, hatchback, pick-up,
coupe, etc.).

- **Alternativa considerada:** usar un enum en base de datos. Se descarta por
  simplicidad y coherencia con el tratamiento actual de `condicion` y
  `combustible` (strings libres).
- **Impacto:** auto-migración de GORM agrega la columna `tipo`; el seed puede
  asignar tipos a los vehículos para que el filtro tenga datos de prueba.

### 2. Filtros y búsqueda como query params opcionales en GET /api/vehiculos

Se mantiene el mismo endpoint público (sin versión nueva). Los parámetros son
todos opcionales y combinables:

- `busqueda`: texto libre, coincide (case-insensitive) contra marca o modelo.
- `marca`, `modelo`: igualdad exacta.
- `anio_min`/`anio_max`, `precio_min`/`precio_max`: rangos inclusivos.
- `tipo`, `combustible`, `condicion`: igualdad exacta.
- `orden_por` (`precio`|`anio`, default `anio`) y `orden_direccion`
  (`asc`|`desc`, default `desc`).

La capa de service valida la semántica (rangos válidos, valores de
condición/orden conocidos) y el repository construye la consulta GORM
dinámicamente. Los handlers solo parsean strings y delegan.

- **Alternativa considerada:** crear un endpoint dedicado `/api/vehiculos/buscar`.
  Se descarta: el CU-03 ya usa `/api/vehiculos` con paginación y los filtros
  amplían naturalmente ese listado sin duplicar rutas.

### 3. Búsqueda por texto con ILIKE sobre marca o modelo

En PostgreSQL, `WHERE (marca ILIKE '%busqueda%' OR modelo ILIKE '%busqueda%')`
con escape de caracteres comodín (`%`, `_`) en el texto de entrada. No se usan
índices full-text: el volumen de vehículos de un concesionario es pequeño y la
búsqueda por substring cubre el caso de uso.

- **Riesgo:** performance en tablas muy grandes → mitigado por el tamaño real
  del dominio; si crece, se migra a índices GIN de `tsvector`.

### 4. Query builder en el repository con GORM

`Listar` recibe un struct de criterios `FiltrosBusqueda` y encadena `Where`
según cuáles estén presentes, reutilizando la misma consulta base para el
`Count` y el `Find`, incluyendo el `ORDER BY` dinámico con cláusula permitida
(whitelist de columnas `precio`/`anio`) para evitar inyección por orden.

### 5. Frontend: estado de filtros en la página /catalogo

La página mantiene un estado `filtros` (busqueda, marca, modelo, anioMin,
anioMax, precioMin, precioMax, tipo, combustible, condicion, ordenPor,
ordenDireccion). Al cambiar filtros se vuelve a la página 1 y se llama a
`api.listarVehiculos(pagina, tamano, filtros)`. El panel usa selects para los
valores acotados (tipo, combustible, condición, orden) e inputs para los
rangos. Se usa `useDeferredValue` o debounce simple en la búsqueda de texto
para no disparar una petición por tecla.

- **Alternativa considerada:** filtrar del lado del cliente con el listado
  completo. Se descarta: el catálogo está paginado en el servidor y los filtros
  deben operar sobre todo el stock, no solo la página visible.

### 6. Tipos TypeScript actualizados

Se agrega `tipo: string` a `Vehiculo`, `VehiculoEntrada` y `ResumenVehiculo`
(para mostrar el tipo en tarjetas del catálogo). El cliente `api.ts` construye
la query string con `URLSearchParams` para mantener el código legible al
combinar filtros.

## Risks / Trade-offs

- [Cambio de contrato en la ficha técnica] → El alta/modificación y el detalle
  ahora incluyen `tipo`. Es compatible hacia atrás en lectura (campo nuevo) pero
  el ABM exige `tipo` no vacío al guardar; se actualizan formularios y seed.
- [Inyección SQL en el ordenamiento] → Se usa una whitelist de columnas y
  direcciones (`precio`/`anio`, `asc`/`desc`) en el repository; nunca se
  interpola el valor crudo.
- [Búsqueda ILIKE con `%`/`_` en la entrada] → Se escapan los comodines para que
  el usuario busque literales y no expanda la consulta.
- [Demasiadas peticiones al escribir en el campo de búsqueda] → Debounce en el
  frontend para evitar una llamada por tecla.
- [Auto-migración GORM con columna nueva en tablas existentes] → GORM agrega la
  columna sin pérdida de datos; se recomienda verificar el arranque y, si hay
  filas sembradas, actualizar el seed para poblarla.

## Migration Plan

1. Agregar `Tipo` al modelo `models.Vehiculo` (la auto-migración en
   `database.go` crea la columna).
2. Actualizar seed para poblar `tipo` en los vehículos de prueba.
3. Desplegar backend (nuevo campo y filtros) y frontend (panel y formulario).
4. Rollback: revertir el commit del cambio; la columna extra no rompe el
   arranque si el código nuevo no la usa.

## Open Questions

- ¿Qué valores de `tipo` se ofrecen en el select del formulario y del filtro?
  Propuesta inicial: sedán, SUV, hatchback, pick-up, coupe. Confirmar con el
  usuario durante la implementación.
- ¿Se muestra el tipo en las tarjetas del listado o solo en el detalle?
  Propuesta: mostrarlo en el detalle y opcionalmente como etiqueta en la
  tarjeta.
