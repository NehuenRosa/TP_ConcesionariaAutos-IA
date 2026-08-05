# Design: CU-07 Turno de test drive

## Context

El sistema ya cuenta con catálogo público de vehículos disponibles (CU-03/CU-04)
y consultas entre cliente y vendedor (CU-05/CU-06). No existe un modelo de
turnos ni un flujo para coordinar pruebas de manejo. La regla de negocio clave
es: no puede existir más de un turno para la misma unidad en la misma fecha y
franja horaria. El stack es Go + Gin + GORM (backend por capas handler →
service → repository) y React + Vite + TS (frontend). Ver proposal.md para el
motivo y alcance.

## Goals / Non-Goals

**Goals:**
- Modelo `TurnoTestDrive` persistido con su tabla en español.
- Solicitud de turno por cliente autenticado desde el detalle del vehículo.
- Validación de superposición a nivel de servicio (unicidad de
  vehículo + fecha + franja para turnos activos).
- Gestión de turnos por vendedor (listar, confirmar, cancelar, completar).
- Cancelación de turno propio por el cliente.
- Sin dependencias nuevas.

**Non-Goals:**
- Notificaciones por email ni recordatorios de turnos.
- Calendario visual con arrastrar y soltar.
- Asignación automática de vendedores a turnos.
- Integración con CU-08 (reservas) ni CU-10 (chatbot).
- Bloqueo del vehículo: un test drive no cambia el estado de la unidad.

## Decisions

### D1: Modelo de datos

Entidad `TurnoTestDrive` en `backend/internal/models/turno_test_drive.go`:

```go
type TurnoTestDrive struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    VehiculoID uint      `gorm:"not null;index" json:"vehiculoId"`
    Vehiculo   Vehiculo  `gorm:"foreignKey:VehiculoID" json:"-"`
    ClienteID  uint      `gorm:"not null;index" json:"clienteId"`
    Cliente    Usuario   `gorm:"foreignKey:ClienteID" json:"-"`
    Fecha      string    `gorm:"type:date;not null;index" json:"fecha"`
    Franja     string    `gorm:"not null" json:"franja"`
    Estado     string    `gorm:"not null;index;default:solicitado" json:"estado"`
    CreatedAt  time.Time `json:"-"`
    UpdatedAt  time.Time `json:"-"`
}
```

- **Tabla**: `turnos_test_drive` (convención del repo: nombres en español).
- **Fecha como string** (`YYYY-MM-DD`): evita problemas de zona horaria al
  comparar fechas exactas con `time.Time` en PostgreSQL. El handler la valida y
  normaliza al formato ISO.
- **Alternativa descartada**: guardar `HoraInicio`/`HoraFin` como timestamps.
  Agrega complejidad de validación de solapamiento por rangos sin aportar
  flexibilidad real, ya que se usan franjas predefinidas.

### D2: Franjas horarias predefinidas

Catálogo fijo en el backend (constantes), expuesto por un endpoint público:

| Franja | Inicio | Fin |
|--------|--------|-----|
| `manana` | 09:00 | 12:00 |
| `tarde` | 14:00 | 18:00 |

- El turno persiste el identificador (`manana`/`tarde`).
- `GET /api/test-drives/franjas` devuelve el catálogo (id, inicio, fin).
- **Alternativa descartada**: tabla `franjas` en base de datos. Sin requisito de
  administrarlas, un catálogo en código evita una tabla extra y una migración
  más. Si en el futuro se necesitan franjas dinámicas, migrar es directo.

### D3: API REST

| Método | Ruta | Descripción | Auth |
|--------|------|-------------|------|
| `GET /api/test-drives/franjas` | Catálogo de franjas | Público |
| `POST /api/test-drives` | Solicitar turno | JWT (cliente) |
| `GET /api/test-drives/mis-turnos` | Turnos del cliente | JWT (cliente) |
| `DELETE /api/test-drives/:id` | Cancelar turno propio | JWT (cliente) |
| `GET /api/test-drives` | Listar turnos (filtro `estado` opcional) | JWT (vendedor) |
| `PUT /api/test-drives/:id/confirmar` | Confirmar turno | JWT (vendedor) |
| `PUT /api/test-drives/:id/cancelar` | Cancelar turno | JWT (vendedor) |
| `PUT /api/test-drives/:id/completar` | Completar turno | JWT (vendedor) |

- Las rutas de vendedor exigen `rol = vendedor` (o administrador, según
  `ExigirRol`).
- El `DELETE` del cliente es de "cancelación lógica" (cambia estado a
  `cancelado`), no borrado físico, para preservar historial.

### D4: Lógica de negocio (service)

**Crear turno (cliente):**
1. Validar que el vehículo existe y está `disponible` (reutiliza el repositorio
   de vehículos) → `404`.
2. Validar fecha (no anterior a hoy, formato ISO) y franja (pertenece al
   catálogo) → `400`.
3. Verificar superposición: no existe otro turno activo (`solicitado` o
   `confirmado`) con el mismo `vehiculoId`, `fecha` y `franja` → `409`.
4. Crear el turno con estado `solicitado`.

**Prevención de superposición:**
- Consulta en el repository: contar turnos con `vehiculo_id`, `fecha` y
  `franja` dados y `estado IN (solicitado, confirmado)`.
- La validación vive en el service (regla de negocio); el repository solo
  expone el conteo.
- **Mitigación de carrera**: la validación es a nivel aplicación. Para el
  volumen esperado es suficiente; en caso de necesitar garantía estricta, se
  podría agregar un índice único parcial SQL posteriormente (ver Riesgos).

**Transiciones de estado:**

| Desde | Acción | A |
|-------|--------|---|
| `solicitado` | confirmar (vendedor) | `confirmado` |
| `solicitado`/`confirmado` | cancelar (vendedor o cliente propio) | `cancelado` |
| `confirmado` | completar (vendedor) | `completado` |

- Cualquier otra transición → `409` "no se puede cambiar el estado".

**Permisos:**
- El cliente solo accede a sus propios turnos (`cliente_id` == usuario JWT).
  Un turno ajeno se trata como inexistente → `404` (no revela existencia).

### D5: Estructura de capas backend

- `models/turno_test_drive.go`: modelo + constantes de estado y franjas.
- `repositories/turnos_test_drive.go`: interfaz `TurnoTestDriveRepository` con
  `Crear`, `ObtenerPorID`, `ListarPorCliente`, `Listar`, `ExisteSuperposicion`,
  `Actualizar`; implementación GORM.
- `services/turnos_test_drive.go`: interfaz `TurnoTestDriveService` con
  `Solicitar`, `ListarMisTurnos`, `Cancelar`, `Listar`, `Confirmar`,
  `Cancelar`, `Completar`, `Franjas`.
- `handlers/turnos_test_drive.go`: parsea request/response y delega.
- `router/router.go`: registra el grupo `/api/test-drives`.
- `database/database.go`: agrega `TurnoTestDrive` a `AutoMigrate`.

### D6: Frontend

**Tipos** (`frontend/src/types/testDrive.ts`): `FranjaHoraria`, `TurnoTestDrive`
con vehículo resumido, `SolicitarTestDrive`, `EstadoTurno`.

**Cliente HTTP** (`services/api.ts`): `obtenerFranjas()`, `solicitarTestDrive()`,
`listarMisTestDrives()`, `cancelarTestDrive(id)`, `listarTestDrives(estado)`,
`confirmarTestDrive(id)`, `cancelarTestDriveVendedor(id)`,
`completarTestDrive(id)`.

**Páginas y rutas:**

| Ruta | Componente | Protección |
|------|------------|------------|
| `/catalogo/:id/test-drive` | `FormularioTestDrive` | JWT (cliente) |
| `/mis-test-drives` | `MisTestDrives` | JWT (cliente) |
| `/vendedor/test-drives` | `GestionTestDrives` | JWT (vendedor) |

- `DetalleVehiculo.tsx`: agrega botón "Solicitar test drive" solo para clientes
  autenticados, enlazando a `/catalogo/:id/test-drive`.
- `FormularioTestDrive.tsx`: fecha + selector de franja (desde
  `GET /api/test-drives/franjas`), maneja el `409` mostrando "turno ocupado".
- `MisTestDrives.tsx`: lista de turnos propios con cancelar para activos.
- `GestionTestDrives.tsx`: lista con filtro por estado y acciones
  confirmar/cancelar/completar.
- `LayoutBase.tsx`: agrega "Mis test drives" (cliente) y "Test drives"
  (vendedor) en el header.

### D7: Errores

| Código | Mensaje |
|--------|---------|
| 400 | Fecha o franja horaria inválida |
| 401 | No autenticado |
| 403 | No autorizado (no es vendedor) |
| 404 | Vehículo o turno no encontrado / turno ajeno |
| 409 | Turno ocupado o transición de estado inválida |
| 500 | Error interno del servidor |

## Risks / Trade-offs

- [Validación de superposición a nivel aplicación puede tener carreras] →
  Mitigación: volumen bajo de turnos; si se requiere, agregar índice único
  parcial en SQL (`vehiculo_id, fecha, franja WHERE estado IN ('solicitado',
  'confirmado')`).
- [Fecha como string en vez de `time.Time`] → Trade-off de tipado a cambio de
  comparaciones de fechas simples y sin ambigüedad de zona horaria. El formato
  se valida siempre en el service.
- [Franjas fijas en código] → Si el negocio requiere franjas dinámicas habrá
  que migrar a una tabla; hoy no es un requisito.

## Migration Plan

1. Auto-migración crea `turnos_test_drive` al arrancar el backend (GORM).
2. Sin datos existentes que migrar: la tabla es nueva.
3. Rollback: el cambio es aditivo; revertir el commit deja la tabla huérfana
   sin impacto funcional (GORM no la elimina).

## Open Questions

- Ninguna que bloquee la implementación; el alcance definido en los specs es
  suficiente.
