# Design: CU-08 Reserva de vehículo

## Context

El sistema ya cuenta con catálogo público de vehículos disponibles (CU-03/CU-04),
consultas (CU-05/CU-06) y turnos de test drive (CU-07). No existe un modelo de
reservas: el estado `reservado` está definido en `models/vehiculo.go` pero nunca
se asigna por un flujo. La regla de negocio clave: reservar una unidad disponible
la bloquea (disponible → reservado) y solo el vendedor puede confirmar la venta
(vendido) o liberarla (disponible). El stack es Go + Gin + GORM (backend por
capas handler → service → repository) y React + Vite + TS (frontend). Ver
proposal.md para el motivo y alcance.

## Goals / Non-Goals

**Goals:**
- Modelo `Reserva` persistido con su tabla en español (`reservas`).
- Reserva de un vehículo disponible por cliente autenticado, que lo bloquea de
  forma atómica (reserva + cambio de estado del vehículo en una transacción).
- Cancelación de reserva propia por el cliente (libera la unidad).
- Gestión de reservas por vendedor (listar con filtro, confirmar venta, cancelar).
- Sin dependencias nuevas.

**Non-Goals:**
- Vencimiento automático de reservas (por plazo de retención).
- Pagos ni seña en línea.
- Notificaciones por email.
- Reservas con turnos o franjas horarias (eso es CU-07).
- Integración con CU-09 (dashboard) ni CU-10 (chatbot).

## Decisions

### D1: Modelo de datos

Entidad `Reserva` en `backend/internal/models/reserva.go`:

```go
type Reserva struct {
    ID         uint     `gorm:"primaryKey" json:"id"`
    VehiculoID uint     `gorm:"not null;index" json:"vehiculoId"`
    Vehiculo   Vehiculo `gorm:"foreignKey:VehiculoID" json:"-"`
    ClienteID  uint     `gorm:"not null;index" json:"clienteId"`
    Cliente    Usuario  `gorm:"foreignKey:ClienteID" json:"-"`
    Estado     string   `gorm:"not null;index;default:activa" json:"estado"`
    CreatedAt  time.Time `json:"-"`
    UpdatedAt  time.Time `json:"-"`
}
```

- **Tabla**: `reservas` (convención del repo: nombres en español).
- **Estados** (`EstadoReservaActiva`, `EstadoReservaVendida`,
  `EstadoReservaCancelada`): la reserva nace `activa` (bloquea la unidad), pasa a
  `vendida` cuando el vendedor confirma la venta, o a `cancelada` si se libera.
- **Alternativa descartada**: reservas con `pendiente`/`confirmada` (estilo test
  drive). En las reservas el cliente ya se compromete a la unidad al crearla, por
  lo que el bloqueo es inmediato; no hay paso previo de confirmación.
- No se persisten fecha/hora de entrega ni precio: el precio queda en el
  vehículo (CU-02). Si en el futuro se necesita retención de precio, se agrega un
  campo `precio_fijado`.

### D2: Cambio de estado atómico

A diferencia de CU-07 (el test drive no cambia el vehículo), la reserva modifica
el estado de la unidad. Para no dejar estados inconsistentes, el repository de
reservas ejecuta las operaciones compuestas dentro de una transacción GORM:

| Operación | Reserva | Vehículo |
|-----------|---------|----------|
| `CrearYReservar` | crea `activa` | → `reservado` |
| `CancelarYLiberar` | → `cancelada` | → `disponible` |
| `ConfirmarVentaYMarcarVendido` | → `vendida` | → `vendido` |

- El repository recibe el estado destino del vehículo y ejecuta
  `base.WithContext(ctx).Transaction(...)` (mismo patrón que
  `repositories/vehiculos.go` en `Crear`/`Actualizar`).
- La decisión de *qué* estados asignar es del service (regla de negocio); el
  repository solo persiste.
- **Alternativa descartada**: el service coordina reservas y vehículos con dos
  repos distintos fuera de transacción. Ante un fallo intermedio dejaría una
  reserva sin bloqueo o una unidad bloqueada sin reserva.

### D3: API REST

| Método | Ruta | Descripción | Auth |
|--------|------|-------------|------|
| `POST /api/reservas` | Crear reserva de un vehículo disponible | JWT (cliente) |
| `GET /api/reservas/mis-reservas` | Reservas del cliente | JWT (cliente) |
| `DELETE /api/reservas/:id` | Cancelar reserva propia | JWT (cliente) |
| `GET /api/reservas` | Listar reservas (filtro `estado` opcional) | JWT (vendedor) |
| `PUT /api/reservas/:id/confirmar` | Confirmar venta | JWT (vendedor) |
| `PUT /api/reservas/:id/cancelar` | Cancelar reserva | JWT (vendedor) |

- Las rutas de vendedor exigen `rol = vendedor` (o administrador, según
  `ExigirRol`).
- El `DELETE` del cliente es una "cancelación lógica" (cambia a `cancelada` y
  libera la unidad), no borrado físico, para preservar historial.

### D4: Lógica de negocio (service)

**Crear reserva (cliente):**
1. Validar que el vehículo existe y está `disponible` → `404` si no existe o no
   es comercializable; `409` si está `reservado` o `vendido` (la validación
   previa al insert evita depender de un error de constraint).
2. Crear la reserva `activa` y cambiar el vehículo a `reservado` en una sola
   transacción.
3. Verificación en el service del estado del vehículo justo antes de la
   transacción para el mensaje de error adecuado (404 vs 409).

**Cancelación (cliente):**
- Solo reservas propias y `activas`. Una reserva ajena se trata como inexistente
  → `404` (no revela existencia). Estado no activo → `409`.

**Transiciones de estado (vendedor):**

| Desde | Acción | A |
|-------|--------|---|
| `activa` | confirmar venta (vendedor) | `vendida` (vehículo → `vendido`) |
| `activa` | cancelar (vendedor o cliente propio) | `cancelada` (vehículo → `disponible`) |

- Cualquier otra transición → `409` "no se puede cambiar el estado".

**Permisos:**
- El cliente solo accede a sus propias reservas (`cliente_id` == usuario JWT).

### D5: Estructura de capas backend

- `models/reserva.go`: modelo + constantes de estado + `EsActiva()`.
- `repositories/reservas.go`: interfaz `ReservaRepository` con `CrearYReservar`,
  `ObtenerPorID`, `ListarPorCliente`, `Listar`, `CancelarYLiberar`,
  `ConfirmarVentaYMarcarVendido`; implementación GORM con transacciones.
- `services/reservas.go`: interfaz `ReservaService` con `Crear`,
  `ListarMisReservas`, `Cancelar`, `Listar`, `ConfirmarVenta`, `CancelarComoVendedor`.
- `handlers/reservas.go`: parsea request/response y delega.
- `router/router.go`: registra el grupo `/api/reservas`.
- `database/database.go`: agrega `Reserva` a `AutoMigrate`.

### D6: Frontend

**Tipos** (`frontend/src/types/reserva.ts`): `EstadoReserva`, `Reserva` con
vehículo resumido y cliente, `CrearReserva`.

**Cliente HTTP** (`services/api.ts`): `crearReserva(vehiculoId)`,
`listarMisReservas()`, `cancelarReserva(id)`, `listarReservas(estado)`,
`confirmarReservaVenta(id)`, `cancelarReservaVendedor(id)`.

**Páginas y rutas:**

| Ruta | Componente | Protección |
|------|------------|------------|
| `/catalogo/:id/reservar` | `FormularioReserva` | JWT (cliente) |
| `/mis-reservas` | `MisReservas` | JWT (cliente) |
| `/vendedor/reservas` | `GestionReservas` | JWT (vendedor) |

- `DetalleVehiculo.tsx`: agrega botón "Reservar este vehículo" junto a "Solicitar
  test drive", solo para clientes autenticados, enlazando a `/catalogo/:id/reservar`.
- `FormularioReserva.tsx`: muestra la unidad a reservar y un botón de confirmar;
  maneja el `409` mostrando "la unidad ya no está disponible".
- `MisReservas.tsx`: lista de reservas propias con cancelar para las activas.
- `GestionReservas.tsx`: lista con filtro por estado y acciones confirmar venta /
  cancelar.
- `LayoutBase.tsx`: agrega "Reservas" (cliente) y "Reservas" (vendedor) en el
  header.

### D7: Errores

| Código | Mensaje |
|--------|---------|
| 400 | Datos inválidos / filtro de estado inválido |
| 401 | No autenticado |
| 403 | No autorizado (no es vendedor) |
| 404 | Vehículo o reserva no encontrado / reserva ajena |
| 409 | Vehículo no disponible para reservar o transición de estado inválida |
| 500 | Error interno del servidor |

## Risks / Trade-offs

- [Carrera entre la validación de disponibilidad y el insert (dos clientes
  reservan la misma unidad)] → Mitigación: el service vuelve a verificar el
  estado del vehículo dentro de la transacción y aborta si cambió; para
  garantía estricta se podría agregar un índice único sobre
  `vehiculo_id WHERE estado = 'activa'` (SQL) en el futuro.
- [El bloqueo del vehículo depende del correcto flujo de cancelación] → El
  vendedor y el cliente pueden liberar la unidad; la auto-migración mantiene el
  modelo sin datos iniciales.
- [Sin vencimiento automático de reservas] → Una reserva activa bloquea la
  unidad indefinidamente hasta que se confirme o cancele; hoy no es un requisito
  y se documenta como no-goal.

## Migration Plan

1. Auto-migración crea `reservas` al arrancar el backend (GORM).
2. Sin datos existentes que migrar: la tabla es nueva.
3. Rollback: el cambio es aditivo; revertir el commit deja la tabla huérfana sin
   impacto funcional (GORM no la elimina).

## Open Questions

- Ninguna que bloquee la implementación; el alcance definido en los specs es
  suficiente.
