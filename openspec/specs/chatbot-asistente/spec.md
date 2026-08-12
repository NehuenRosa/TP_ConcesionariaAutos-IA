# chatbot-asistente Specification

## Purpose
Asistente conversacional público que responde preguntas del usuario en lenguaje
natural sobre el stock real de la concesionaria y lo orienta a consultar un
vehículo o solicitar un test drive.
## Requirements
### Requirement: Responder mensajes del usuario

El sistema SHALL exponer un endpoint público `POST /api/chatbot/mensajes` que
reciba el mensaje del usuario y devuelva la respuesta del asistente en lenguaje
natural. El request puede incluir un campo `historial` opcional (lista de turnos
previos con `rol` y `contenido`) para mantener el contexto de la conversación.
El endpoint no requiere autenticación.

#### Scenario: Mensaje válido

- **WHEN** un usuario envía `POST /api/chatbot/mensajes` con un `mensaje` no vacío
- **THEN** el sistema responde con `200` y la respuesta del asistente en un campo `respuesta`

#### Scenario: Mensaje con historial

- **WHEN** un usuario envía `POST /api/chatbot/mensajes` con `mensaje` y un
  `historial` con turnos previos
- **THEN** el sistema usa el historial como contexto y responde de forma coherente
  con la conversación previa

#### Scenario: Mensaje vacío

- **WHEN** un usuario envía `POST /api/chatbot/mensajes` con `mensaje` vacío o ausente
- **THEN** el sistema responde con error `400` y un mensaje en español

#### Scenario: Modelo local no disponible

- **WHEN** el usuario envía un mensaje y el modelo de lenguaje local no está disponible
- **THEN** el sistema responde con `200` y una respuesta en español que orienta al usuario a
  consultar o pedir un test drive por los canales habituales

### Requirement: Respuesta basada en el stock real

El sistema SHALL construir la respuesta del asistente usando como contexto la
lista de vehículos en estado `disponible` del catálogo, con su ficha técnica
(marca, modelo, año, precio, tipo, combustible, transmisión, condición y
kilometraje). El asistente NO debe inventar vehículos o precios que no existan
en el stock.

#### Scenario: Pregunta sobre un vehículo del stock

- **WHEN** el usuario pregunta por una marca o modelo que existe en el stock disponible
- **THEN** la respuesta del asistente refleja los datos reales de ese vehículo (marca,
  modelo, año y precio)

#### Scenario: Stock disponible vacío

- **WHEN** el usuario envía un mensaje y no hay vehículos en estado `disponible`
- **THEN** la respuesta del asistente indica que no hay vehículos disponibles en el momento

#### Scenario: Pregunta sobre un vehículo inexistente

- **WHEN** el usuario pregunta por una marca o modelo que no está en el stock
- **THEN** la respuesta del asistente indica que no hay stock de ese vehículo, sin inventar datos

### Requirement: Orientación hacia consultas y test drives

El sistema SHALL configurar al asistente para orientar al usuario hacia las
acciones disponibles del sistema: crear una consulta/cotización sobre un
vehículo de interés o solicitar un turno de test drive.

#### Scenario: Usuario muestra interés en un vehículo

- **WHEN** el usuario expresa interés o pide cotización de un vehículo disponible
- **THEN** la respuesta del asistente lo orienta a consultar el vehículo o solicitar un test drive

### Requirement: Widget de chat en el frontend

El sistema SHALL mostrar un widget de chat flotante en las páginas públicas del
frontend que permita al usuario enviar mensajes y ver la conversación.

#### Scenario: Abrir el chat

- **WHEN** el usuario hace clic en el botón flotante del chat
- **THEN** se muestra la ventana de chat con el historial de la sesión

#### Scenario: Enviar mensaje desde el widget

- **WHEN** el usuario escribe un mensaje y lo envía en el widget
- **THEN** el widget muestra el mensaje del usuario, un indicador de "escribiendo…"
  y luego la respuesta del asistente

#### Scenario: Respuesta fallida en el widget

- **WHEN** el envío del mensaje al backend falla (error de red o respuesta no `2xx`)
- **THEN** el widget muestra un mensaje de error en español y permite reintentar

#### Scenario: Mensaje muy largo

- **WHEN** el usuario envía un mensaje de más de 1000 caracteres
- **THEN** el sistema responde con error `400` y un mensaje en español indicando el límite

#### Scenario: Foto con tipo MIME vacío

- **WHEN** el usuario adjunta una foto sin tipo MIME pero con extensión válida (JPG/JPEG/PNG/WebP)
- **THEN** el widget la acepta para la tasación

#### Scenario: Tasación que supera el tiempo máximo

- **WHEN** la tasación no responde dentro del tiempo máximo configurado en el cliente
- **THEN** el widget corta la petición y muestra un mensaje de error en español
  (no queda en "pensando" de forma indefinida) y permite reintentar

### Requirement: Tasación del auto del usuario por fotos

El sistema SHALL exponer un endpoint público `POST /api/chatbot/tasacion` que
reciba fotos del auto del usuario (multipart/form-data) y una `descripcion`
opcional. El sistema identifica el vehículo (marca, modelo, año, estado,
kilometraje) usando un modelo con visión sobre las imágenes y la descripción, y
devuelve en lenguaje natural el resultado de la tasación. El endpoint no
requiere autenticación.

#### Scenario: Tasación con fotos válidas

- **WHEN** un usuario envía `POST /api/chatbot/tasacion` con una o más fotos
- **THEN** el sistema identifica el vehículo, consulta su valor de referencia y
  responde con `200` y una respuesta en lenguaje natural con el valor, la
  identificación del vehículo y la fuente del valor

#### Scenario: Sin fotos adjuntas ni descripción

- **WHEN** un usuario envía `POST /api/chatbot/tasacion` sin archivos ni `descripcion`
- **THEN** el sistema responde con error `400` y un mensaje en español

#### Scenario: Formato de archivo no soportado

- **WHEN** un usuario envía archivos que no son imágenes (JPG/PNG/WebP)
- **THEN** el sistema responde con error `400` y un mensaje en español

#### Scenario: Modelo de visión no disponible

- **WHEN** el usuario envía fotos y el modelo de visión no está disponible
- **THEN** el sistema responde con `200` y un mensaje en español que orienta al
  usuario a acercarse a la concesionaria o contactar a un vendedor para la tasación

#### Scenario: Vehículo no identificado con certeza

- **WHEN** el usuario envía fotos y el sistema no puede identificar marca y modelo
- **THEN** el sistema responde con `200` y un mensaje en español que aclara que
  no puede identificar el vehículo, que no va a inventar un valor, y pide
  completar la descripción o acercarse a la concesionaria

### Requirement: Valores de tasación oficiales y no inventados

El sistema SHALL devolver en la tasación el valor de mercado de la **Guía
Oficial de Precios de la CCA** (consultada a través de la API pública de
ArgAutos) para el vehículo identificado, expresado en dólares (US$) y pesos
(ARS), con la versión y el año de referencia y la fuente. El sistema NO debe
generar montos inventados: el valor lo compone el sistema con la referencia
oficial; si no hay referencia, se responde con honestidad.

#### Scenario: Valor de referencia disponible

- **WHEN** el sistema identifica el vehículo y la guía tiene un valor de referencia
- **THEN** la respuesta incluye el valor real en US$ y ARS, la versión y año de
  referencia y menciona la Guía de la CCA como fuente

#### Scenario: Valor de referencia no disponible

- **WHEN** el sistema identifica el vehículo pero la guía no tiene valor de referencia
- **THEN** la respuesta indica que no hay valor oficial para ese vehículo y
  orienta a consultar en la concesionaria, sin inventar un monto

### Requirement: Comparación entre el auto del usuario y el catálogo

El sistema SHALL detectar en la conversación cuándo el usuario quiere comparar
su auto actual con uno de nuestro catálogo. Cuando lo detecta, el asistente
ofrece compararlos y pregunta qué aspectos quiere comparar (precio, consumo,
potencia, seguridad, etc.). La comparación usa la ficha técnica real del
vehículo del stock (disponible) y el conocimiento del modelo para el auto del
usuario. Si el vehículo del catálogo no está en stock disponible, el asistente
lo indica sin inventar datos.

#### Scenario: Usuario pide comparar y no indica aspectos

- **WHEN** el usuario dice qué auto tiene y qué auto quiere comprar del catálogo
  sin especificar aspectos
- **THEN** el asistente ofrece la comparación y pregunta qué aspectos quiere comparar

#### Scenario: Usuario indica los aspectos a comparar

- **WHEN** el usuario responde qué aspectos quiere comparar (p. ej. precio,
  consumo y potencia)
- **THEN** el asistente compara ambos vehículos usando los datos reales del
  catálogo para el auto de la concesionaria y el conocimiento del modelo para el
  auto del usuario

#### Scenario: Vehículo del catálogo no encontrado

- **WHEN** el usuario pide comparar con un vehículo que no está en el stock disponible
- **THEN** el asistente indica que ese vehículo no está disponible, sin inventar precios ni fichas

### Requirement: Tasación con descripción opcional

El sistema SHALL aceptar en `POST /api/chatbot/tasacion` un campo de texto
opcional `descripcion` que el usuario puede adjuntar junto con las fotos (año,
kilometraje, estado, versión) para afinar la identificación del vehículo.

#### Scenario: Tasación con fotos y descripción

- **WHEN** un usuario envía fotos junto con una `descripcion` (año, kilometraje, estado)
- **THEN** la identificación considera la descripción provista además de las imágenes

#### Scenario: Tasación solo con descripción

- **WHEN** un usuario envía `POST /api/chatbot/tasacion` con `descripcion` pero sin fotos
- **THEN** el sistema responde con `200` y un resultado de tasación basado en la
  identificación desde la descripción (con valor oficial si existe, o un mensaje
  honesto si no)

