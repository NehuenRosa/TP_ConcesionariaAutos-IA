# Spec: gestion-vehiculos

## Purpose

Permitir al administrador gestionar el stock de vehículos del concesionario:
listar todas las unidades (en cualquier estado), ver su detalle, darlas de alta,
modificarlas y darlas de baja lógicamente, cargando la ficha técnica (marca,
modelo, año, kilometraje, combustible, transmisión, precio, condición) e
imágenes. Es la contraparte administrativa del catálogo público (CU-03).

## ADDED Requirements

### Requirement: Listado administrativo de vehículos

El sistema SHALL exponer un endpoint `GET /api/admin/vehiculos`, accesible solo
con rol `administrador`, que devuelva el listado paginado de **todos** los
vehículos del stock, sin importar su estado. La respuesta SHALL incluir los
metadatos de paginación (página, tamaño, total) y, por cada vehículo, su ficha
básica con estado. El endpoint SHALL aceptar un filtro opcional `estado` que
restrinja el resultado a un estado válido (`disponible`, `reservado`,
`vendido` o `dado_de_baja`).

#### Scenario: Listado con todos los estados

- **WHEN** un administrador solicita el listado administrativo sin filtro
- **THEN** el sistema responde con un listado paginado que incluye vehículos en
  cualquier estado

#### Scenario: Listado filtrado por estado

- **WHEN** un administrador solicita el listado con `estado=dado_de_baja`
- **THEN** el sistema responde únicamente con vehículos en estado
  `dado_de_baja`

#### Scenario: Filtro por estado inválido

- **WHEN** un administrador solicita el listado con un valor de `estado` que no
  es uno de los estados conocidos
- **THEN** el sistema responde con error `400` y un mensaje en español

#### Scenario: Listado sin resultados

- **WHEN** no existen vehículos que cumplan el filtro solicitado
- **THEN** el sistema responde con una lista vacía y el total en cero

### Requirement: Detalle administrativo de vehículo

El sistema SHALL exponer un endpoint `GET /api/admin/vehiculos/:id`, accesible
solo con rol `administrador`, que devuelva el detalle completo de un vehículo
**en cualquier estado**, incluyendo ficha técnica y galería de imágenes.

#### Scenario: Detalle de vehículo en cualquier estado

- **WHEN** un administrador solicita el detalle de un vehículo existente
- **THEN** el sistema responde con la ficha técnica completa, el estado y la
  galería de imágenes del vehículo

#### Scenario: Detalle de vehículo inexistente

- **WHEN** un administrador solicita el detalle de un vehículo cuyo
  identificador no existe
- **THEN** el sistema responde con error `404` y un mensaje en español

#### Scenario: Identificador inválido

- **WHEN** un administrador solicita el detalle con un identificador que no es
  numérico
- **THEN** el sistema responde con error `400` y un mensaje en español

### Requirement: Alta de vehículo

El sistema SHALL exponer un endpoint `POST /api/admin/vehiculos`, accesible solo
con rol `administrador`, que cree un vehículo nuevo con su ficha técnica
(`marca`, `modelo`, `anio`, `kilometraje`, `combustible`, `transmision`,
`precio`, `condicion`) y su lista de URLs de imágenes. El `estado` inicial SHALL
ser `disponible` salvo que el request lo especifique. El sistema SHALL validar
que los campos requeridos estén presentes, que `anio` esté en un rango
razonable, que `precio` sea positivo, que `condicion` sea `nuevo` o `usado` y
que `estado` sea uno de los estados conocidos.

#### Scenario: Alta válida

- **WHEN** un administrador envía un vehículo con ficha técnica completa y
  lista de imágenes válida
- **THEN** el sistema crea el vehículo, persiste sus imágenes y responde con el
  detalle completo del vehículo creado

#### Scenario: Alta con estado por defecto

- **WHEN** un administrador envía un vehículo sin indicar estado
- **THEN** el sistema crea el vehículo con estado `disponible`

#### Scenario: Alta con datos inválidos

- **WHEN** un administrador envía un vehículo con campos requeridos faltantes,
  un `anio` fuera de rango, un `precio` no positivo o una `condicion`/`estado`
  desconocidos
- **THEN** el sistema responde con error `400` y un mensaje en español

### Requirement: Modificación de vehículo

El sistema SHALL exponer un endpoint `PUT /api/admin/vehiculos/:id`, accesible
solo con rol `administrador`, que actualice la ficha técnica, el estado y la
galería de imágenes de un vehículo existente. Las imágenes SHALL reemplazarse
por la lista enviada en el request. Se aplican las mismas validaciones de campos
que en el alta.

#### Scenario: Modificación válida

- **WHEN** un administrador modifica la ficha técnica, el estado o las imágenes
  de un vehículo existente
- **THEN** el sistema actualiza los datos, reemplaza la galería de imágenes y
  responde con el detalle actualizado

#### Scenario: Modificación de vehículo inexistente

- **WHEN** un administrador modifica un vehículo cuyo identificador no existe
- **THEN** el sistema responde con error `404` y un mensaje en español

#### Scenario: Modificación con datos inválidos

- **WHEN** un administrador envía una modificación con campos inválidos o un
  `estado` desconocido
- **THEN** el sistema responde con error `400` y un mensaje en español

### Requirement: Baja lógica de vehículo

El sistema SHALL exponer un endpoint `DELETE /api/admin/vehiculos/:id`,
accesible solo con rol `administrador`, que dé de baja un vehículo de forma
lógica: cambia su `estado` a `dado_de_baja` **sin eliminar el registro** ni sus
imágenes. La operación SHALL ser idempotente: si el vehículo ya está
`dado_de_baja`, responde `200` sin cambios.

#### Scenario: Baja de un vehículo disponible

- **WHEN** un administrador da de baja un vehículo existente
- **THEN** el sistema cambia su estado a `dado_de_baja` y responde con el
  detalle actualizado

#### Scenario: Baja idempotente

- **WHEN** un administrador da de baja un vehículo que ya está `dado_de_baja`
- **THEN** el sistema responde `200` sin modificar el registro

#### Scenario: Baja de vehículo inexistente

- **WHEN** un administrador da de baja un vehículo cuyo identificador no existe
- **THEN** el sistema responde con error `404` y un mensaje en español

#### Scenario: Vehículo dado de baja no aparece en el catálogo público

- **WHEN** un vehículo pasa a estado `dado_de_baja`
- **THEN** el catálogo público (`GET /api/vehiculos`) deja de mostrarlo

### Requirement: Página administrativa de vehículos

El sistema SHALL ofrecer una página administrativa en la ruta
`/admin/vehiculos` con el listado paginado de todos los vehículos consumiendo
`GET /api/admin/vehiculos`, con filtro por estado, acciones para crear, editar y
dar de baja, y mensajes en español para los estados de carga, vacío y error.

#### Scenario: Listado administrativo con datos

- **WHEN** un administrador accede a `/admin/vehiculos`
- **THEN** el sistema muestra la tabla de vehículos con su estado y los
  controles de paginación y filtro

#### Scenario: Creación desde la página administrativa

- **WHEN** un administrador accede al formulario de alta y guarda un vehículo
- **THEN** el sistema crea el vehículo y regresa al listado administrativo

#### Scenario: Edición desde la página administrativa

- **WHEN** un administrador accede al formulario de edición de un vehículo y
  guarda los cambios
- **THEN** el sistema actualiza el vehículo y regresa al listado administrativo

#### Scenario: Baja desde la página administrativa

- **WHEN** un administrador confirma la baja de un vehículo
- **THEN** el sistema da de baja el vehículo y actualiza el listado mostrando su
  nuevo estado
