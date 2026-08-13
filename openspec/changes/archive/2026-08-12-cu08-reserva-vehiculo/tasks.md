# Tasks: CU-08 Reserva de vehículo

## 1. Backend - Modelo

- [x] 1.1 Crear modelo `Reserva` en `backend/internal/models/reserva.go` con campos ID, VehiculoID, ClienteID, Estado, CreatedAt, UpdatedAt
- [x] 1.2 Definir constantes de estado: `EstadoReservaActiva`, `EstadoReservaVendida`, `EstadoReservaCancelada`
- [x] 1.3 Implementar `TableName()` devolviendo `reservas`
- [x] 1.4 Implementar `EsActiva()` (estado == `activa`)

## 2. Backend - Repository

- [x] 2.1 Crear `backend/internal/repositories/reservas.go` con interfaz `ReservaRepository`
- [x] 2.2 Implementar `CrearYReservar(ctx, reserva)` con transacción: crear reserva `activa` y setear el vehículo a `reservado`
- [x] 2.3 Implementar `ObtenerPorID(ctx, id)` con Preload de Vehiculo e Imagenes y Cliente
- [x] 2.4 Implementar `ListarPorCliente(ctx, clienteID)` con Preload de Vehiculo e Imagenes, ordenado por CreatedAt
- [x] 2.5 Implementar `Listar(ctx, estado)` con Preload de Vehiculo e Imagenes y Cliente, ordenado por CreatedAt
- [x] 2.6 Implementar `CancelarYLiberar(ctx, reserva)` con transacción: reserva a `cancelada` y vehículo a `disponible`
- [x] 2.7 Implementar `ConfirmarVentaYMarcarVendido(ctx, reserva)` con transacción: reserva a `vendida` y vehículo a `vendido`

## 3. Backend - Service

- [x] 3.1 Crear `backend/internal/services/reservas.go` con interfaz `ReservaService`
- [x] 3.2 Definir errores de negocio: `ErrReservaNoEncontrada`, `ErrVehiculoNoDisponible` (reutilizado), `ErrReservaEstadoInvalido`, `ErrFiltroEstadoReservaInvalido`
- [x] 3.3 Implementar `Crear`: validar vehículo disponible (404 si no existe/no disponible, 409 si reservado o vendido) y crear reserva `activa` con bloqueo del vehículo
- [x] 3.4 Implementar `ListarMisReservas(ctx, clienteID)`
- [x] 3.5 Implementar `Cancelar(ctx, id, clienteID)`: solo reservas propias y activas; ajena → no encontrada
- [x] 3.6 Implementar `Listar(ctx, estado)`: validar filtro de estado opcional
- [x] 3.7 Implementar `ConfirmarVenta(ctx, id)`: solo desde `activa`, vehículo a `vendido`
- [x] 3.8 Implementar `CancelarComoVendedor(ctx, id)`: solo desde `activa`, vehículo a `disponible`

## 4. Backend - Handler

- [x] 4.1 Crear `backend/internal/handlers/reservas.go` con struct `ReservaHandler` y `ReservaResumen`
- [x] 4.2 Implementar `Crear`: parsear body (vehiculoId), extraer cliente del JWT, responder 201 o error 400/404/409
- [x] 4.3 Implementar `ListarMisReservas`: extraer cliente del JWT
- [x] 4.4 Implementar `Cancelar`: extraer cliente del JWT, mapear errores 404/409
- [x] 4.5 Implementar `Listar`: extraer vendedor del JWT, parsear filtro `estado`
- [x] 4.6 Implementar `ConfirmarVenta` y `CancelarComoVendedor`: extraer vendedor del JWT, mapear errores 404/409
- [x] 4.7 Implementar `aReservaResumen`/`aReservasResumen` con vehículo resumido y cliente

## 5. Backend - Router y migración

- [x] 5.1 Registrar en `backend/internal/router/router.go` el grupo `/api/reservas` con rutas de cliente y subgrupo de vendedor
- [x] 5.2 Agregar `Reserva` a `AutoMigrate` en `backend/internal/database/database.go`

## 6. Frontend - Tipos y cliente HTTP

- [x] 6.1 Crear `frontend/src/types/reserva.ts` con `EstadoReserva`, `Reserva`, `CrearReserva`
- [x] 6.2 Agregar en `frontend/src/services/api.ts`: `crearReserva`, `listarMisReservas`, `cancelarReserva`, `listarReservas`, `confirmarReservaVenta`, `cancelarReservaVendedor`

## 7. Frontend - Páginas y rutas

- [x] 7.1 Crear `frontend/src/pages/FormularioReserva.tsx`: muestra la unidad a reservar, confirmación y manejo de errores 404/409
- [x] 7.2 Crear `frontend/src/pages/MisReservas.tsx`: listado de reservas propias con cancelar para activas
- [x] 7.3 Crear `frontend/src/pages/GestionReservas.tsx`: listado con filtro por estado y acciones confirmar venta/cancelar
- [x] 7.4 Modificar `frontend/src/pages/DetalleVehiculo.tsx`: botón "Reservar este vehículo" visible solo para clientes autenticados
- [x] 7.5 Registrar rutas en `frontend/src/routes/Rutas.tsx`: `/catalogo/:id/reservar` y `/mis-reservas` (cliente), `/vendedor/reservas` (vendedor)
- [x] 7.6 Agregar links de navegación en `frontend/src/layouts/LayoutBase.tsx` para cliente y vendedor

## 8. Verificación

- [x] 8.1 Ejecutar `go build ./...` y `go vet ./...` sin errores en `backend/`
- [x] 8.2 Ejecutar `npm run build` sin errores en `frontend/`
- [x] 8.3 Probar flujo de cliente: reservar, ver 409 de unidad no disponible, ver mis reservas, cancelar propia (vehículo vuelve a disponible)
- [x] 8.4 Probar flujo de vendedor: listar, confirmar venta (vehículo vendido), cancelar (vehículo disponible), transiciones inválidas
- [x] 8.5 Actualizar `docs/roadmap.md` (CU-08 → Implementado)

Verificación end-to-end contra la API (`docker compose up -d postgres backend`):

**Cliente (8.3):** reservar 201 (activa, vehículo reservado) · reservar la misma
unidad 409 · reservar vehículo inexistente 404 · mis reservas 200 · cancelar
propia 200 (cancelada, vehículo disponible) · cancelar ajena 404 · cancelar ya
cancelada 409 · sin token 401.

**Vendedor (8.4):** listar 200 · filtro por estado 200 · filtro inválido 400 ·
confirmar desde cancelada 409 · confirmar venta → vendida (vehículo vendido) ·
cancelar reserva 200 (vehículo disponible) · cliente en ruta de vendedor 403 ·
confirmar venta de vehículo ya vendido 409.
