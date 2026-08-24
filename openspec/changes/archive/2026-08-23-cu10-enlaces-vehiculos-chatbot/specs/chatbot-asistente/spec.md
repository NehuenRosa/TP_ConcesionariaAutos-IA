# Delta spec: chatbot-asistente

## ADDED Requirements

### Requirement: Enlaces a fichas de vehículos

El sistema SHALL permitir que el asistente señale los vehículos concretos del
stock que menciona mediante el marcador interno `[VEHICULO:<id>]`, usando el id
que figura en el contexto del stock. El backend SHALL remover los marcadores de
la respuesta visible y SHALL devolver en `POST /api/chatbot/mensajes` el campo
`vehiculosMencionados` con los ids únicos mencionados, ordenados y limitados a
vehículos presentes en el contexto servido. El widget SHALL renderizar, debajo
del mensaje del asistente, un enlace por vehículo hacia `/catalogo/:id`. Si el
modelo no incluye marcadores, la respuesta SHALL mostrarse igual que hoy, sin
enlaces.

#### Scenario: Respuesta con enlaces a vehículos

- **WHEN** el asistente menciona un vehículo disponible (id 3) y agrega el
  marcador `[VEHICULO:3]`
- **THEN** la respuesta visible no contiene el marcador, el campo
  `vehiculosMencionados` es `[3]` y el widget muestra un enlace
  "Ver ficha" que navega a `/catalogo/3`

#### Scenario: Vehículo fuera del contexto

- **WHEN** el modelo incluye un marcador con un id que no figuraba en el
  contexto del stock (ej. `[VEHICULO:99]`)
- **THEN** el backend descarta ese id y no lo incluye en
  `vehiculosMencionados`

#### Scenario: Respuesta sin vehículos mencionados

- **WHEN** el asistente responde una pregunta general sin mencionar vehículos
  concretos
- **THEN** la respuesta no incluye marcadores y `vehiculosMencionados` está
  vacío o ausente
