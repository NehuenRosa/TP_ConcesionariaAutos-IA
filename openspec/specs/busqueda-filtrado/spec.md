# Spec: busqueda-filtrado

## Purpose

Permitir al visitante del catálogo público encontrar vehículos por texto libre
y acotar los resultados con filtros combinables (marca, modelo, rango de años,
rango de precio, tipo, combustible y condición) y ordenamiento por precio o
año. Complementa al catálogo paginado (CU-03).

## Requirements

### Requirement: Búsqueda por texto libre

El sistema SHALL permitir buscar vehículos por texto libre en el catálogo
público mediante el parámetro de consulta `busqueda` del endpoint
`GET /api/vehiculos`. La búsqueda SHALL hacer coincidir el texto, sin distinguir
mayúsculas y minúsculas, contra la marca o el modelo del vehículo. Los
resultados SHALL limitarse siempre a vehículos con estado `disponible`.

#### Scenario: Búsqueda por marca o modelo

- **WHEN** un visitante busca "corolla"
- **THEN** el sistema responde únicamente con vehículos disponibles cuya marca o
  modelo contienen "corolla"

#### Scenario: Búsqueda sin resultados

- **WHEN** un visitante busca un texto que no coincide con ninguna marca o modelo
- **THEN** el sistema responde con una lista vacía y el total en cero

#### Scenario: Búsqueda sin texto

- **WHEN** un visitante solicita el catálogo sin el parámetro `busqueda`
- **THEN** el sistema responde con todos los vehículos disponibles

### Requirement: Filtros combinables por ficha técnica

El sistema SHALL permitir combinar filtros en el endpoint `GET /api/vehiculos`
usando los parámetros de consulta `marca`, `modelo`, `anio_min`, `anio_max`,
`precio_min`, `precio_max`, `tipo`, `combustible` y `condicion`. Los filtros
SHALL aplicarse en forma conjunta (AND) y SHALL restringir el resultado a
vehículos con estado `disponible`. El sistema SHALL responder con error `400` y
un mensaje en español cuando un rango es inválido (mínimo mayor que máximo), un
valor numérico no es válido o `condicion` no es `nuevo`/`usado`.

#### Scenario: Filtros combinados

- **WHEN** un visitante filtra por marca "Toyota", condición "usado" y un rango
  de precio de 10000000 a 20000000
- **THEN** el sistema responde únicamente con vehículos disponibles de marca
  Toyota, usados y dentro del rango de precio indicado

#### Scenario: Filtro por rango de años y precio

- **WHEN** un visitante filtra con `anio_min=2018` y `anio_max=2022`
- **THEN** el sistema responde únicamente con vehículos disponibles cuyo año
  está entre 2018 y 2022 inclusive

#### Scenario: Filtro por tipo, combustible y condición

- **WHEN** un visitante filtra por tipo "suv", combustible "nafta" y condición
  "nuevo"
- **THEN** el sistema responde únicamente con vehículos disponibles que
  coinciden con los tres valores

#### Scenario: Rango inválido

- **WHEN** un visitante envía un rango con mínimo mayor que máximo
  (por ejemplo `precio_min=30000000` y `precio_max=10000000`)
- **THEN** el sistema responde con error `400` y un mensaje en español

#### Scenario: Condición inválida

- **WHEN** un visitante envía `condicion` con un valor que no es `nuevo` ni `usado`
- **THEN** el sistema responde con error `400` y un mensaje en español

#### Scenario: Filtros sin resultados

- **WHEN** un visitante aplica filtros que no coinciden con ningún vehículo
- **THEN** el sistema responde con una lista vacía y el total en cero

### Requirement: Ordenamiento por precio o año

El sistema SHALL permitir ordenar los resultados del catálogo público mediante
los parámetros de consulta `orden_por` (con valores `precio` o `anio`) y
`orden_direccion` (con valores `asc` o `desc`). El valor por defecto SHALL ser
orden por año descendente. El sistema SHALL responder con error `400` y un
mensaje en español cuando `orden_por` o `orden_direccion` no son válidos.

#### Scenario: Ordenamiento por precio ascendente

- **WHEN** un visitante ordena por `orden_por=precio` y `orden_direccion=asc`
- **THEN** el sistema responde los vehículos disponibles ordenados de menor a
  mayor precio

#### Scenario: Ordenamiento por año descendente

- **WHEN** un visitante ordena por `orden_por=anio` y `orden_direccion=desc`
- **THEN** el sistema responde los vehículos disponibles ordenados del año más
  reciente al más antiguo

#### Scenario: Ordenamiento por defecto

- **WHEN** un visitante solicita el catálogo sin parámetros de ordenamiento
- **THEN** el sistema responde los vehículos disponibles ordenados por año
  descendente

#### Scenario: Ordenamiento inválido

- **WHEN** un visitante envía `orden_por=kilometraje` u otro valor no soportado
- **THEN** el sistema responde con error `400` y un mensaje en español

### Requirement: Búsqueda y filtros en la página de catálogo

El sistema SHALL ofrecer en la página pública `/catalogo` un panel de búsqueda y
filtros con: campo de texto libre, selección de marca, modelo, tipo,
combustible, condición, rangos de año y precio, y control de ordenamiento por
precio o año. Al aplicar o modificar los filtros, el sistema SHALL refrescar el
listado paginado consumiendo `GET /api/vehiculos` con los parámetros
correspondientes y SHALL volver a la primera página. La página SHALL mostrar un
mensaje en español cuando no hay resultados y un control para limpiar todos los
filtros.

#### Scenario: Aplicación de filtros desde la página

- **WHEN** un visitante completa búsqueda y filtros y los aplica
- **THEN** el sistema refresca el listado mostrando solo los vehículos que
  cumplen todos los criterios, en la primera página

#### Scenario: Filtros sin resultados en la página

- **WHEN** un visitante aplica filtros que no coinciden con ningún vehículo
- **THEN** el sistema muestra un mensaje en español indicando que no hay
  resultados para los filtros aplicados

#### Scenario: Limpiar filtros

- **WHEN** un visitante presiona el control para limpiar filtros
- **THEN** el sistema restablece el panel a sus valores por defecto y muestra
  nuevamente el catálogo completo
