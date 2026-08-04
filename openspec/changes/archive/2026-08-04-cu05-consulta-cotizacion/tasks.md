# Tasks: CU-05 Consulta / cotización

## Backend

### T1: Modelo de Consulta y Mensaje

**Archivo:** `backend/internal/models/consulta.go`

- [x] T1.1: Crear modelo `Consulta` con campos: ID, VehiculoID, ClienteID, VendedorID (puntero), Estado, Mensajes, CreatedAt, UpdatedAt
- [x] T1.2: Crear modelo `Mensaje` con campos: ID, ConsultaID, RemitenteID, Contenido, Leido, CreatedAt
- [x] T1.3: Definir constantes de estado: `EstadoPendiente`, `EstadoEnConversacion`, `EstadoCerrada`
- [x] T1.4: Implementar `TableName()` en español para ambos modelos

### T2: Repository de Consultas

**Archivo:** `backend/internal/repositories/consultas.go`

- [x] T2.1: Crear interfaz `ConsultaRepository` con métodos: Crear, ObtenerPorID, ListarPorCliente, ListarPorVendedor, Actualizar, Eliminar
- [x] T2.2: Implementar `consultaRepository` con GORM
- [x] T2.3: Implementar `ListarPorCliente` con Preload de Vehiculo, Vendedor y último Mensaje, ordenado por UpdatedAt descendente
- [x] T2.4: Implementar `ListarPorVendedor` con Preload de Vehiculo, Cliente y último Mensaje, ordenado por UpdatedAt descendente
- [x] T2.5: Implementar `ObtenerPorID` con Preload completo (Vehiculo, Cliente, Vendedor, Mensajes)

### T3: Repository de Mensajes

**Archivo:** `backend/internal/repositories/mensajes.go`

- [x] T3.1: Crear interfaz `MensajeRepository` con métodos: Crear, ListarPorConsulta, ObtenerNuevos, MarcarComoLeidos
- [x] T3.2: Implementar `mensajeRepository` con GORM
- [x] T3.3: Implementar `ListarPorConsulta` ordenado por CreatedAt ascendente
- [x] T3.4: Implementar `ObtenerNuevos` con filtro por timestamp y remitente
- [x] T3.5: Implementar `MarcarComoLeidos` para actualizar estado leído

### T4: Service de Consultas

**Archivo:** `backend/internal/services/consultas.go`

- [x] T4.1: Crear interfaz `ConsultaService` con métodos: Crear, ObtenerPorID, ListarPorCliente, ListarPorVendedor, Tomar, Cerrar, Eliminar
- [x] T4.2: Implementar `consultaService` con inyección de dependencias
- [x] T4.3: Implementar `Crear`: validar vehículo disponible, crear consulta y primer mensaje
- [x] T4.4: Implementar `Tomar`: validar estado pendiente, asignar vendedor, cambiar estado
- [x] T4.5: Implementar `Cerrar`: validar que es el vendedor asignado, cambiar estado
- [x] T4.6: Implementar `Eliminar`: validar que está cerrada, eliminar consulta y mensajes

### T5: Service de Mensajes

**Archivo:** `backend/internal/services/mensajes.go`

- [x] T5.1: Crear interfaz `MensajeService` con métodos: Enviar, ObtenerPorConsulta, ObtenerNuevos
- [x] T5.2: Implementar `mensajeService` con inyección de dependencias
- [x] T5.3: Implementar `Enviar`: validar participante, validar consulta no cerrada, crear mensaje
- [x] T5.4: Implementar `ObtenerNuevos`: obtener mensajes desde timestamp, marcar como leídos

### T6: Handler de Consultas

**Archivo:** `backend/internal/handlers/consultas.go`

- [x] T6.1: Crear struct `ConsultaHandler` con servicios inyectados
- [x] T6.2: Implementar `Crear`: parsear body, llamar service, retornar 201
- [x] T6.3: Implementar `ListarMisConsultas`: extraer cliente del JWT, llamar service
- [x] T6.4: Implementar `ListarBandeja`: extraer vendedor del JWT, llamar service
- [x] T6.5: Implementar `Tomar`: extraer vendedor del JWT, llamar service
- [x] T6.6: Implementar `Cerrar`: extraer vendedor del JWT, llamar service
- [x] T6.7: Implementar `Eliminar`: extraer vendedor del JWT, llamar service

### T7: Handler de Mensajes

**Archivo:** `backend/internal/handlers/mensajes.go`

- [x] T7.1: Crear struct `MensajeHandler` con servicios inyectados
- [x] T7.2: Implementar `Enviar`: parsear body, validar participante, llamar service
- [x] T7.3: Implementar `ObtenerMensajes`: validar participante, llamar service
- [x] T7.4: Implementar `ObtenerNuevos`: parsear query `desde`, llamar service

### T8: Router de Consultas

**Archivo:** `backend/internal/router/router.go`

- [x] T8.1: Registrar grupo `/api/consultas` con middleware de autenticación
- [x] T8.2: Registrar rutas POST, GET para crear y listar
- [x] T8.3: Registrar rutas PUT para tomar, cerrar
- [x] T8.4: Registrar ruta DELETE para eliminar
- [x] T8.5: Registrar sub-rutas `/mensajes` y `/mensajes/nuevos`

### T9: Migración

**Archivo:** `backend/internal/database/database.go`

- [x] T9.1: Agregar `Consulta` y `Mensaje` a AutoMigrate

## Frontend

### T10: Tipos TypeScript

**Archivo:** `frontend/src/types/consulta.ts`

- [x] T10.1: Crear tipo `Consulta` con todos los campos
- [x] T10.2: Crear tipo `Mensaje` con todos los campos
- [x] T10.3: Crear tipo `ConsultaResumen` para listados
- [x] T10.4: Crear tipo `CrearConsulta` para el formulario

### T11: Cliente HTTP

**Archivo:** `frontend/src/services/api.ts`

- [x] T11.1: Agregar método `crearConsulta(datos)`
- [x] T11.2: Agregar método `listarMisConsultas()`
- [x] T11.3: Agregar método `listarBandeja()`
- [x] T11.4: Agregar método `tomarConsulta(id)`
- [x] T11.5: Agregar método `cerrarConsulta(id)`
- [x] T11.6: Agregar método `eliminarConsulta(id)`
- [x] T11.7: Agregar método `enviarMensaje(consultaId, contenido)`
- [x] T11.8: Agregar método `obtenerMensajes(consultaId)`
- [x] T11.9: Agregar método `obtenerMensajesNuevos(consultaId, desde)`

### T12: Modificar DetalleVehiculo

**Archivo:** `frontend/src/pages/DetalleVehiculo.tsx`

- [x] T12.1: Agregar botón "Consultar" visible solo para clientes autenticados
- [x] T12.2: Crear modal o sección con formulario de mensaje
- [x] T12.3: Implementar envío de consulta a la API
- [x] T12.4: Mostrar feedback de éxito/error

### T13: Página Bandeja de Entrada

**Archivo:** `frontend/src/pages/BandejaEntrada.tsx`

- [x] T13.1: Crear página con lista de tarjetas de consultas
- [x] T13.2: Mostrar info del vehículo, cliente, estado, preview del último mensaje
- [x] T13.3: Agregar badge de mensajes nuevos (punto rojo)
- [x] T13.4: Implementar badge de mensajes nuevos usando `mensajesNuevos` del backend (recarga con evento `mensajes-leidos` + polling cada 5 segundos)
- [x] T13.5: Agregar botón "Tomar" para consultas pendientes
- [x] T13.6: Agregar botón "Cerrar" para consultas en conversación
- [x] T13.7: Agregar botón "Eliminar" para consultas cerradas
- [x] T13.8: Navegar al chat al hacer clic en una tarjeta

### T14: Página Mis Consultas (Cliente)

**Archivo:** `frontend/src/pages/MisConsultas.tsx`

- [x] T14.1: Crear página con vista tipo chat (lista izquierda, conversación derecha)
- [x] T14.2: Lista de consultas con info del vehículo y vendedor
- [x] T14.3: Badge de mensajes nuevos en cada consulta de la lista
- [x] T14.4: Implementar badge de mensajes nuevos usando `mensajesNuevos` del backend (recarga con evento `mensajes-leidos` + polling cada 5 segundos)
- [x] T14.5: Al seleccionar una consulta, mostrar conversación completa

### T15: Componente de Chat

**Archivo:** `frontend/src/components/ChatConsulta.tsx`

- [x] T15.1: Crear componente reutilizable de chat
- [x] T15.2: Mostrar mensajes con estilo de burbujas (izquierda/derecha según remitente)
- [x] T15.3: Input para escribir mensaje y botón enviar
- [x] T15.4: Auto-scroll al último mensaje
- [x] T15.5: Cargar mensajes al montar y actualizar periódicamente

### T16: Rutas

**Archivo:** `frontend/src/routes/Rutas.tsx`

- [x] T16.1: Agregar ruta `/mis-consultas` protegida (cliente)
- [x] T16.2: Agregar ruta `/mis-consultas/:id` protegida (cliente)
- [x] T16.3: Agregar ruta `/vendedor/bandeja` protegida (vendedor)
- [x] T16.4: Agregar ruta `/vendedor/bandeja/:id` protegida (vendedor)

### T17: Navegación

**Archivo:** `frontend/src/layouts/LayoutBase.tsx`

- [x] T17.1: Agregar link "Mis Consultas" en header para clientes
- [x] T17.2: Agregar link "Bandeja" en header para vendedores
- [x] T17.3: Agregar badge de notificaciones en los links si hay mensajes nuevos

## Verificación

### T18: Verificación Backend

- [x] T18.1: Ejecutar `go build ./...` sin errores
- [x] T18.2: Ejecutar `go vet ./...` sin errores
- [x] T18.3: Probar endpoints con `Invoke-WebRequest`:
  - Crear consulta
  - Listar bandeja del vendedor
  - Tomar consulta
  - Enviar mensajes
  - Obtener mensajes nuevos
  - Cerrar consulta
  - Eliminar consulta

### T19: Verificación Frontend

- [x] T19.1: Ejecutar `npm run build` sin errores
- [x] T19.2: Probar flujo completo:
  - Login como cliente
  - Ir a detalle de vehículo
  - Crear consulta
  - Login como vendedor
  - Ver bandeja de entrada
  - Tomar consulta
  - Responder mensaje
  - Login como cliente
  - Ver respuesta en Mis Consultas
  - Cerrar consulta como vendedor
  - Eliminar consulta cerrada
