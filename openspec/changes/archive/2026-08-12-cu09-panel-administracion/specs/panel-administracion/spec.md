## Purpose

Provee al administrador un panel de métricas agregadas del negocio (stock por
estado, consultas por período, reservas activas y test drives agendados)
expuestas por una API autenticada y presentadas con gráficos simples.

## ADDED Requirements

### Requirement: Consultar métricas agregadas del negocio

El sistema SHALL exponer un endpoint `GET /api/admin/metricas` accesible solo
para usuarios autenticados con rol `administrador`, que devuelve en una única
respuesta las métricas agregadas de las tablas existentes:

- **Vehículos por estado**: cantidad de vehículos por cada estado
  (`disponible`, `reservado`, `vendido`, `dado_de_baja`).
- **Consultas por período**: cantidad de consultas creadas por día dentro del
  período solicitado (por defecto últimos 30 días), incluidos los días sin
  consultas.
- **Reservas activas y vendidas**: total de reservas en estado `activa` y en
  estado `vendida`.
- **Test drives agendados y completados**: total de turnos en estados
  `solicitado`/`confirmado` y en estado `completado`.
- **Consultas abiertas**: total de consultas en estados `pendiente`/
  `en_conversacion`.
- **Usuarios registrados**: total de cuentas de usuario.

El endpoint SHALL aceptar un parámetro de consulta `periodo` (cantidad de días
a considerar para las consultas por período), permitiendo los valores `7`, `30`
y `90`; el valor por defecto SHALL ser `30`. Si el valor no está en la lista, el
sistema SHALL responder con error HTTP `400` y un mensaje en español.

#### Scenario: Administrador consulta métricas con período por defecto

- **WHEN** un usuario autenticado con rol `administrador` envía
  `GET /api/admin/metricas`
- **THEN** el sistema responde `200 OK` con un objeto JSON que incluye
  `vehiculosPorEstado`, `consultasPorPeriodo` (30 días), `reservasActivas`,
  `reservasVendidas`, `testDrivesAgendados`, `testDrivesCompletados`,
  `consultasAbiertas` y `totalUsuarios`

#### Scenario: Administrador consulta métricas con período válido

- **WHEN** un administrador envía `GET /api/admin/metricas?periodo=7`
- **THEN** el sistema responde `200 OK` y `consultasPorPeriodo` contiene un
  registro por cada uno de los últimos 7 días, incluso los que tienen cero
  consultas

#### Scenario: Administrador envía período inválido

- **WHEN** un administrador envía `GET /api/admin/metricas?periodo=15`
- **THEN** el sistema responde HTTP `400` con cuerpo JSON
  `{"error": "mensaje en español"}` indicando que el período no es válido

#### Scenario: Usuario sin rol administrador consulta métricas

- **WHEN** un usuario autenticado con rol `cliente` o `vendedor` envía
  `GET /api/admin/metricas`
- **THEN** el sistema responde HTTP `403` y no devuelve datos de métricas

#### Scenario: Usuario no autenticado consulta métricas

- **WHEN** una solicitud sin token JWT envía `GET /api/admin/metricas`
- **THEN** el sistema responde HTTP `401` y no devuelve datos de métricas

### Requirement: Presentar el panel con resumen y gráficos simples

La página `/admin` SHALL mostrar un dashboard que:

- Presenta tarjetas de resumen con: vehículos disponibles, reservas activas,
  test drives agendados y consultas abiertas.
- Presenta un gráfico de barras con la cantidad de vehículos por estado.
- Presenta un gráfico de barras con la evolución de consultas por día en el
  período seleccionado.
- Permite al administrador elegir el período del gráfico de consultas entre
  los últimos 7, 30 o 90 días, y recarga las métricas al cambiar la selección.
- Conserva los accesos rápidos a las secciones de gestión (vehículos y
  usuarios).
- Muestra el estado de carga mientras la respuesta está pendiente y, si la
  solicitud falla, muestra un mensaje de error en español sin romper la página.

Los gráficos SHALL implementarse con componentes propios de CSS/SVG sin
agregar librerías externas de visualización.

#### Scenario: El panel muestra las métricas del negocio

- **WHEN** un administrador navega a `/admin` con la API respondiendo
  correctamente
- **THEN** se muestran las tarjetas de resumen, el gráfico de vehículos por
  estado y el gráfico de consultas por período con los valores recibidos

#### Scenario: El administrador cambia el período del gráfico de consultas

- **WHEN** un administrador selecciona "Últimos 7 días" en el selector de
  período
- **THEN** el sistema consulta `GET /api/admin/metricas?periodo=7` y el gráfico
  se actualiza con los datos del nuevo período

#### Scenario: La carga de métricas falla

- **WHEN** la solicitud de métricas devuelve un error
- **THEN** el panel muestra un mensaje de error en español indicando que no se
  pudieron cargar las métricas y conserva el resto de la navegación
