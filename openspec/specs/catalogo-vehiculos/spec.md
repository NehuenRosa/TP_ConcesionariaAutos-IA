# Spec: catalogo-vehiculos

## Purpose

Exponer el stock de vehículos disponibles del concesionario al público, con un
listado paginado y el detalle de cada unidad. Es la base de los casos de uso
posteriores (consulta, reserva, test drive) que referencian un vehículo.

## Requirements

### Requirement: Listado paginado de vehículos disponibles

El sistema SHALL exponer un endpoint público `GET /api/vehiculos` que devuelva
el listado de vehículos del concesionario paginado, mostrando únicamente los
vehículos con estado `disponible`. La respuesta SHALL incluir los datos
paginados (página, tamaño de página, total) y la ficha básica de cada vehículo:
marca, modelo, año, precio, imagen y condición.

#### Scenario: Listado con vehículos disponibles

- **WHEN** un cliente o visitante solicita la primera página del catálogo
- **THEN** el sistema responde con un listado paginado que solo contiene
  vehículos con estado `disponible`

#### Scenario: Listado sin vehículos disponibles

- **WHEN** no existen vehículos con estado `disponible`
- **THEN** el sistema responde con una lista vacía y el total en cero

#### Scenario: Paginación fuera de rango

- **WHEN** un cliente o visitante solicita una página más allá de la última
- **THEN** el sistema responde con una lista vacía sin errores

### Requirement: Detalle de vehículo por identificador

El sistema SHALL exponer un endpoint público `GET /api/vehiculos/:id` que
devuelva el detalle completo de un vehículo disponible, incluyendo la ficha
técnica: marca, modelo, año, kilometraje, combustible, transmisión, precio,
condición y galería de imágenes.

#### Scenario: Vehículo disponible existente

- **WHEN** un cliente o visitante solicita el detalle de un vehículo que existe
  y está disponible
- **THEN** el sistema responde con la ficha técnica completa del vehículo

#### Scenario: Vehículo inexistente

- **WHEN** un cliente o visitante solicita el detalle de un vehículo cuyo
  identificador no existe
- **THEN** el sistema responde con error `404` y un mensaje en español

#### Scenario: Vehículo no disponible

- **WHEN** un cliente o visitante solicita el detalle de un vehículo que existe
  pero no está en estado `disponible`
- **THEN** el sistema responde con error `404` para no revelar información de
  unidades no comercializables

### Requirement: Página pública de catálogo

El sistema SHALL ofrecer una página pública de catálogo en la ruta `/catalogo`
que muestre el listado paginado de vehículos disponibles consumiendo el
endpoint público `GET /api/vehiculos`, con controles para navegar entre páginas.

#### Scenario: Navegación al catálogo

- **WHEN** un cliente o visitante accede a `/catalogo`
- **THEN** el sistema muestra las tarjetas de los vehículos disponibles y los
  controles de paginación

#### Scenario: Página de catálogo sin resultados

- **WHEN** no hay vehículos disponibles
- **THEN** el sistema muestra un mensaje en español indicando que no hay
  vehículos disponibles

### Requirement: Página pública de detalle de vehículo

El sistema SHALL ofrecer una página pública de detalle en la ruta
`/catalogo/:id` que muestre la ficha técnica del vehículo consumiendo el
endpoint público `GET /api/vehiculos/:id`.

#### Scenario: Acceso al detalle de un vehículo disponible

- **WHEN** un cliente o visitante accede a `/catalogo/:id` de un vehículo
  disponible
- **THEN** el sistema muestra la ficha técnica y la galería de imágenes del
  vehículo

#### Scenario: Detalle de vehículo inexistente o no disponible

- **WHEN** un cliente o visitante accede a `/catalogo/:id` de un vehículo
  inexistente o no disponible
- **THEN** el sistema muestra un mensaje de error en español con un enlace para
  volver al catálogo
