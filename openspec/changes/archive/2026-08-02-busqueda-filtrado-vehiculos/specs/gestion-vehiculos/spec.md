## MODIFIED Requirements

### Requirement: Alta de vehículo

El sistema SHALL exponer un endpoint `POST /api/admin/vehiculos`, accesible solo
con rol `administrador`, que cree un vehículo nuevo con su ficha técnica
(`marca`, `modelo`, `anio`, `kilometraje`, `combustible`, `transmision`,
`precio`, `condicion`, `tipo`) y su lista de URLs de imágenes. El `estado`
inicial SHALL ser `disponible` salvo que el request lo especifique. El sistema
SHALL validar que los campos requeridos estén presentes, que `anio` esté en un
rango razonable, que `precio` sea positivo, que `condicion` sea `nuevo` o
`usado`, que `tipo` sea un valor no vacío y que `estado` sea uno de los estados
conocidos.

#### Scenario: Alta válida

- **WHEN** un administrador envía un vehículo con ficha técnica completa y
  lista de imágenes válida
- **THEN** el sistema crea el vehículo, persiste sus imágenes y responde con el
  detalle completo del vehículo creado, incluido su tipo

#### Scenario: Alta con estado por defecto

- **WHEN** un administrador envía un vehículo sin indicar estado
- **THEN** el sistema crea el vehículo con estado `disponible`

#### Scenario: Alta con datos inválidos

- **WHEN** un administrador envía un vehículo con campos requeridos faltantes,
  un `anio` fuera de rango, un `precio` no positivo, una `condicion`/`estado`
  desconocidos o un `tipo` vacío
- **THEN** el sistema responde con error `400` y un mensaje en español

### Requirement: Modificación de vehículo

El sistema SHALL exponer un endpoint `PUT /api/admin/vehiculos/:id`, accesible
solo con rol `administrador`, que actualice la ficha técnica (incluido el
campo `tipo`), el estado y la galería de imágenes de un vehículo existente. Las
imágenes SHALL reemplazarse por la lista enviada en el request. Se aplican las
mismas validaciones de campos que en el alta.

#### Scenario: Modificación válida

- **WHEN** un administrador modifica la ficha técnica, el estado o las imágenes
  de un vehículo existente
- **THEN** el sistema actualiza los datos, reemplaza la galería de imágenes y
  responde con el detalle actualizado, incluido su tipo

#### Scenario: Modificación de vehículo inexistente

- **WHEN** un administrador modifica un vehículo cuyo identificador no existe
- **THEN** el sistema responde con error `404` y un mensaje en español

#### Scenario: Modificación con datos inválidos

- **WHEN** un administrador envía una modificación con campos inválidos o un
  `estado` desconocido
- **THEN** el sistema responde con error `400` y un mensaje en español

### Requirement: Página administrativa de vehículos

El sistema SHALL ofrecer una página administrativa en la ruta
`/admin/vehiculos` con el listado paginado de todos los vehículos consumiendo
`GET /api/admin/vehiculos`, con filtro por estado, acciones para crear, editar y
dar de baja, y mensajes en español para los estados de carga, vacío y error.
Los formularios de alta y edición SHALL incluir el campo `tipo` de vehículo.

#### Scenario: Listado administrativo con datos

- **WHEN** un administrador accede a `/admin/vehiculos`
- **THEN** el sistema muestra la tabla de vehículos con su estado y los
  controles de paginación y filtro

#### Scenario: Creación desde la página administrativa

- **WHEN** un administrador accede al formulario de alta, completa la ficha
  técnica incluido el tipo, y guarda un vehículo
- **THEN** el sistema crea el vehículo y regresa al listado administrativo

#### Scenario: Edición desde la página administrativa

- **WHEN** un administrador accede al formulario de edición de un vehículo y
  guarda los cambios, incluido el tipo
- **THEN** el sistema actualiza el vehículo y regresa al listado administrativo

#### Scenario: Baja desde la página administrativa

- **WHEN** un administrador confirma la baja de un vehículo
- **THEN** el sistema da de baja el vehículo y actualiza el listado mostrando su
  nuevo estado
