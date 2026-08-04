# Design: CU-05 Consulta / cotización

## Decisiones de diseño

### D1: Modelo de datos

#### Entidad Consulta

```go
type Consulta struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    VehiculoID    uint      `gorm:"not null" json:"vehiculoId"`
    Vehiculo      Vehiculo  `gorm:"foreignKey:VehiculoID" json:"-"`
    ClienteID     uint      `gorm:"not null" json:"clienteId"`
    Cliente       Usuario   `gorm:"foreignKey:ClienteID" json:"-"`
    VendedorID    *uint     `json:"vendedorId"` // nil si no fue tomada
    Vendedor      *Usuario  `gorm:"foreignKey:VendedorID" json:"-"`
    Estado        string    `gorm:"type:varchar(20);not null;default:'pendiente'" json:"estado"`
    Mensajes      []Mensaje `gorm:"foreignKey:ConsultaID" json:"mensajes,omitempty"`
    CreatedAt     time.Time `json:"createdAt"`
    UpdatedAt     time.Time `json:"updatedAt"`
}
```

**Estados posibles:**
- `pendiente`: Consulta creada, sin vendedor asignado
- `en_conversacion`: Vendedor tomó la consulta, intercambiando mensajes
- `cerrada`: Vendedor cerró la consulta

#### Entidad Mensaje

```go
type Mensaje struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    ConsultaID   uint      `gorm:"not null" json:"consultaId"`
    RemitenteID  uint      `gorm:"not null" json:"remitenteId"`
    Remitente    Usuario   `gorm:"foreignKey:RemitenteID" json:"-"`
    Contenido    string    `gorm:"type:text;not null" json:"contenido"`
    Leido        bool      `gorm:"default:false" json:"leido"`
    CreatedAt    time.Time `json:"createdAt"`
}
```

### D2: API REST

#### Consultas

| Método | Ruta | Descripción | Auth |
|--------|------|-------------|------|
| `POST /api/consultas` | Crear consulta | JWT (cliente) |
| `GET /api/consultas/mis-consultas` | Listar consultas del cliente | JWT (cliente) |
| `GET /api/consultas/bandeja` | Listar consultas para vendedor | JWT (vendedor) |
| `PUT /api/consultas/:id/tomar` | Vendedor toma la consulta | JWT (vendedor) |
| `PUT /api/consultas/:id/cerrar` | Vendedor cierra la consulta | JWT (vendedor) |
| `DELETE /api/consultas/:id` | Eliminar consulta cerrada | JWT (vendedor) |

#### Mensajes

| Método | Ruta | Descripción | Auth |
|--------|------|-------------|------|
| `POST /api/consultas/:id/mensajes` | Enviar mensaje | JWT (cliente o vendedor) |
| `GET /api/consultas/:id/mensajes` | Obtener mensajes | JWT (participantes) |
| `GET /api/consultas/:id/mensajes/nuevos?desde=<timestamp>` | Mensajes nuevos desde timestamp | JWT (participantes) |
| `PUT /api/consultas/:id/mensajes/leidos` | Marcar mensajes del otro como leídos | JWT (participantes) |

#### Notificaciones

| Método | Ruta | Descripción | Auth |
|--------|------|-------------|------|
| `GET /api/notificaciones/contador` | Total de mensajes no leídos | JWT (cliente o vendedor) |

### D3: Lógica de negocio

#### Crear consulta (cliente)
1. Validar que el vehículo existe y está disponible
2. Validar que el cliente está autenticado
3. Crear consulta con estado `pendiente`
4. Crear primer mensaje con el contenido del cliente
5. Retornar la consulta creada

#### Tomar consulta (vendedor)
1. Validar que la consulta existe y está en estado `pendiente`
2. Asignar el vendedor actual
3. Cambiar estado a `en_conversacion`
4. Retornar la consulta actualizada

#### Cerrar consulta (vendedor)
1. Validar que la consulta existe y está en estado `en_conversacion`
2. Validar que el vendedor actual es el asignado
3. Cambiar estado a `cerrada`
4. Retornar la consulta actualizada

#### Enviar mensaje
1. Validar que la consulta existe y no está cerrada
2. Validar que el remitente es el cliente o vendedor de la consulta
3. Crear mensaje con el contenido
4. Si es el vendedor y la consulta estaba pendiente, cambiar a `en_conversacion`
5. Retornar el mensaje creado

#### Obtener mensajes nuevos (polling)
1. Validar que la consulta existe
2. Validar que el usuario es participante
3. Filtrar mensajes con `created_at > desde`
4. Retornar solo los mensajes no leídos del otro participante
5. Marcar mensajes como leídos

### D4: Notificaciones (punto rojo)

- **Backend**: `GET /api/consultas/mis-consultas` y `GET /api/consultas/bandeja`
  incluyen el campo `mensajesNuevos` por consulta, calculado contando los
  mensajes del otro participante con `leido = false`
- **Backend**: endpoint liviano `GET /api/notificaciones/contador` devuelve el
  total de mensajes no leídos del usuario autenticado
- **Navbar**: `useNotificaciones` consulta el contador cada 3 segundos y también
  re-verifica al recibir el evento `mensajes-leidos` disparado por el chat
- **Listas**: muestran punto rojo cuando `mensajesNuevos > 0` y se recargan con
  el evento `mensajes-leidos` + polling cada 5 segundos
- **Chat**: al abrir la conversación marca los mensajes del otro como leídos
  (`PUT /api/consultas/:id/mensajes/leidos`) y dispara el evento
  `mensajes-leidos`

### D5: DTOs de respuesta

#### ConsultaResumen (para listados)

```json
{
  "id": 1,
  "vehiculo": {
    "id": 5,
    "marca": "Toyota",
    "modelo": "Corolla",
    "anio": 2023,
    "imagen": "url.jpg"
  },
  "cliente": {
    "id": 10,
    "nombre": "Juan Pérez"
  },
  "vendedor": {
    "id": 20,
    "nombre": "María López"
  },
  "estado": "pendiente",
  "ultimoMensaje": {
    "contenido": "Hola, me interesa este vehículo...",
    "createdAt": "2024-01-15T10:30:00Z"
  },
  "tieneMensajesNuevos": true,
  "createdAt": "2024-01-15T10:30:00Z"
}
```

#### ConsultaDetalle (para chat)

```json
{
  "id": 1,
  "vehiculo": { ... },
  "cliente": { ... },
  "vendedor": { ... },
  "estado": "en_conversacion",
  "mensajes": [ ... ],
  "createdAt": "...",
  "updatedAt": "..."
}
```

#### Mensaje

```json
{
  "id": 100,
  "consultaId": 1,
  "remitenteId": 10,
  "contenido": "Hola, me interesa este vehículo...",
  "leido": false,
  "createdAt": "2024-01-15T10:30:00Z"
}
```

### D6: Frontend - Páginas y componentes

#### Rutas

| Ruta | Componente | Protección |
|------|------------|------------|
| `/catalogo/:id` | DetalleVehiculo (modificado) | Pública |
| `/mis-consultas` | MisConsultas | JWT (cliente) |
| `/mis-consultas/:id` | ChatConsulta | JWT (cliente) |
| `/vendedor/bandeja` | BandejaEntrada | JWT (vendedor) |
| `/vendedor/bandeja/:id` | ChatVendedor | JWT (vendedor) |

#### Componentes clave

1. **BotonConsultar** en DetalleVehiculo: Abre modal o redirige a formulario
2. **BandejaEntrada**: Lista de tarjetas con preview del último mensaje
3. **MisConsultas**: Vista tipo chat (lista izquierda, conversación derecha)
4. **ChatConsulta / ChatVendedor**: Componente de chat reutilizable
5. **BadgeNotificacion**: Punto rojo que se muestra cuando hay mensajes nuevos

### D7: Navegación

- **Header**: Agregar link "Mis Consultas" para clientes autenticados
- **Header**: Agregar link "Bandeja" para vendedores autenticados
- **DetalleVehiculo**: Botón "Consultar" visible solo para clientes autenticados

### D8: Errores

| Código | Mensaje |
|--------|---------|
| 400 | Datos inválidos / faltan campos |
| 401 | No autenticado |
| 403 | No autorizado (no es participante de la consulta) |
| 404 | Consulta o vehículo no encontrado |
| 409 | La consulta ya fue tomada / ya está cerrada |
| 500 | Error interno del servidor |
