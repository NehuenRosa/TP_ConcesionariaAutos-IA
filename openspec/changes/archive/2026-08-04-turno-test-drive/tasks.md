# Tasks: CU-07 Turno de test drive

## 1. Backend - Modelo

- [x] 1.1 Crear modelo `TurnoTestDrive` en `backend/internal/models/turno_test_drive.go` con campos ID, VehiculoID, ClienteID, Fecha (string), Franja, Estado, CreatedAt, UpdatedAt
- [x] 1.2 Definir constantes de estado: `EstadoSolicitado`, `EstadoConfirmado`, `EstadoCancelado`, `EstadoCompletado`
- [x] 1.3 Definir catálogo de franjas horarias (constantes `manana`, `tarde`) con hora de inicio y fin
- [x] 1.4 Implementar `TableName()` devolviendo `turnos_test_drive`

## 2. Backend - Repository

- [x] 2.1 Crear `backend/internal/repositories/turnos_test_drive.go` con interfaz `TurnoTestDriveRepository`
- [x] 2.2 Implementar `Crear(ctx, turno)` con GORM
- [x] 2.3 Implementar `ObtenerPorID(ctx, id)` con Preload de Vehiculo e Imagenes y Cliente
- [x] 2.4 Implementar `ListarPorCliente(ctx, clienteID)` con Preload de Vehiculo e Imagenes, ordenado por Fecha/Franja
- [x] 2.5 Implementar `Listar(ctx, estado)` con Preload de Vehiculo e Imagenes y Cliente, ordenado por Fecha/Franja
- [x] 2.6 Implementar `ExisteSuperposicion(ctx, vehiculoID, fecha, franja)` contando turnos activos (solicitado/confirmado)
- [x] 2.7 Implementar `Actualizar(ctx, turno)` para persistir cambios de estado

## 3. Backend - Service

- [x] 3.1 Crear `backend/internal/services/turnos_test_drive.go` con interfaz `TurnoTestDriveService`
- [x] 3.2 Definir errores de negocio: `ErrTurnoSuperpuesto`, `ErrTurnoNoEncontrado`, `ErrTurnoEstadoInvalido`, `ErrDatosTurnoInvalidos`
- [x] 3.3 Implementar `Solicitar`: validar vehículo disponible, fecha y franja, superposición (409) y crear turno `solicitado`
- [x] 3.4 Implementar `ListarMisTurnos(ctx, clienteID)`
- [x] 3.5 Implementar `Cancelar(ctx, id, clienteID)`: solo turnos propios y activos (solicitado/confirmado)
- [x] 3.6 Implementar `Listar(ctx, estado)`: validar filtro de estado opcional
- [x] 3.7 Implementar `Confirmar(ctx, id)`: solo desde `solicitado`
- [x] 3.8 Implementar `CancelarComoVendedor(ctx, id)`: solo desde `solicitado`/`confirmado`
- [x] 3.9 Implementar `Completar(ctx, id)`: solo desde `confirmado`
- [x] 3.10 Implementar `Franjas()` devolviendo el catálogo predefinido

## 4. Backend - Handler

- [x] 4.1 Crear `backend/internal/handlers/turnos_test_drive.go` con struct `TurnoTestDriveHandler`
- [x] 4.2 Implementar `Solicitar`: parsear body (vehiculoId, fecha, franja), extraer cliente del JWT, responder 201 o error 400/404/409
- [x] 4.3 Implementar `ListarMisTurnos`: extraer cliente del JWT
- [x] 4.4 Implementar `Cancelar`: extraer cliente del JWT, mapear errores 404/409
- [x] 4.5 Implementar `Listar`: extraer vendedor del JWT, parsear filtro `estado`
- [x] 4.6 Implementar `Confirmar`, `CancelarComoVendedor`, `Completar`: extraer vendedor del JWT, mapear errores 404/409
- [x] 4.7 Implementar `Franjas`: responder el catálogo público

## 5. Backend - Router y migración

- [x] 5.1 Registrar en `backend/internal/router/router.go` el grupo `/api/test-drives` con rutas públicas y protegidas
- [x] 5.2 Agregar `TurnoTestDrive` a `AutoMigrate` en `backend/internal/database/database.go`

## 6. Frontend - Tipos y cliente HTTP

- [x] 6.1 Crear `frontend/src/types/testDrive.ts` con `FranjaHoraria`, `EstadoTurnoTestDrive`, `TurnoTestDrive`, `SolicitarTestDrive`
- [x] 6.2 Agregar en `frontend/src/services/api.ts`: `obtenerFranjas`, `solicitarTestDrive`, `listarMisTestDrives`, `cancelarTestDrive`, `listarTestDrives`, `confirmarTestDrive`, `cancelarTestDriveVendedor`, `completarTestDrive`

## 7. Frontend - Páginas y rutas

- [x] 7.1 Crear `frontend/src/pages/FormularioTestDrive.tsx`: fecha + selector de franja, manejo de errores 400/404/409 y mensaje de éxito
- [x] 7.2 Crear `frontend/src/pages/MisTestDrives.tsx`: listado de turnos propios con cancelar para activos
- [x] 7.3 Crear `frontend/src/pages/GestionTestDrives.tsx`: listado con filtro por estado y acciones confirmar/cancelar/completar
- [x] 7.4 Modificar `frontend/src/pages/DetalleVehiculo.tsx`: botón "Solicitar test drive" visible solo para clientes autenticados
- [x] 7.5 Registrar rutas en `frontend/src/routes/Rutas.tsx`: `/catalogo/:id/test-drive` y `/mis-test-drives` (cliente), `/vendedor/test-drives` (vendedor)
- [x] 7.6 Agregar links de navegación en `frontend/src/layouts/LayoutBase.tsx` para cliente y vendedor

## 8. Verificación

- [x] 8.1 Ejecutar `go build ./...` y `go vet ./...` sin errores en `backend/`
- [x] 8.2 Ejecutar `npm run build` sin errores en `frontend/`
- [x] 8.3 Probar flujo de cliente: solicitar turno, ver superposición 409, ver mis turnos, cancelar
- [x] 8.4 Probar flujo de vendedor: listar, confirmar, completar y cancelar turnos

Verificación end-to-end contra la API (`docker compose up -d postgres backend`):

**Cliente (8.3):** franjas públicas 200 · solicitar turno 200 (solicitado) ·
superposición 409 · fecha pasada 400 · vehículo inexistente 404 · mis turnos
200 · cancelar turno propio 200 (cancelado) · cancelar turno ajeno 404 ·
cancelar turno ya cancelado 409.

**Vendedor (8.4):** listar 200 · filtro por estado 200 · filtro inválido 400 ·
completar desde solicitado 409 · confirmar → confirmado · completar →
completado · sin token 401 · cliente en ruta de vendedor 403 · cancelar turno
200 · cancelar turno ya cancelado 409.
