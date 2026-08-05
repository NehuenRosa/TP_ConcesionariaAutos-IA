## MODIFIED Requirements

### Requirement: Página pública de detalle de vehículo

El sistema SHALL ofrecer una página pública de detalle en la ruta
`/catalogo/:id` que muestre la ficha técnica del vehículo consumiendo el
endpoint público `GET /api/vehiculos/:id`, incluida la presentación del tipo de
vehículo. Cuando el usuario autenticado es cliente, la página SHALL mostrar una
acción para solicitar un test drive que enlace a la ruta `/catalogo/:id/test-drive`.

#### Scenario: Acceso al detalle de un vehículo disponible

- **WHEN** un cliente o visitante accede a `/catalogo/:id` de un vehículo
  disponible
- **THEN** el sistema muestra la ficha técnica, el tipo y la galería de imágenes
  del vehículo

#### Scenario: Detalle de vehículo inexistente o no disponible

- **WHEN** un cliente o visitante accede a `/catalogo/:id` de un vehículo
  inexistente o no disponible
- **THEN** el sistema muestra un mensaje de error en español con un enlace para
  volver al catálogo

#### Scenario: Acción de test drive para clientes

- **WHEN** un cliente autenticado accede al detalle de un vehículo disponible
- **THEN** el sistema muestra una acción de "solicitar test drive" que enlaza a
  `/catalogo/:id/test-drive`

#### Scenario: Sin acción de test drive para visitantes

- **WHEN** un visitante no autenticado accede al detalle de un vehículo
  disponible
- **THEN** el sistema no muestra la acción de solicitar test drive
