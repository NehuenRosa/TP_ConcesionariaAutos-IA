# Proposal: Panel de administración (CU-09)

## Why

El panel de administración actual (`/admin`) es solo una página de enlaces a
secciones (vehículos y usuarios): no muestra ninguna métrica del negocio. El
administrador no puede responder preguntas básicas —cuántas unidades hay por
estado, cuántas consultas se recibieron esta semana, cuántas reservas están
activas o cuántos test drives hay agendados— sin inspeccionar cada módulo por
separado.

## What Changes

- Nuevo endpoint `GET /api/admin/metricas` accesible solo con rol
  `administrador` que agrupa los datos de las tablas existentes
  (`vehiculos`, `consultas`, `reservas`, `turnos_test_drive`, `usuarios`) en un
  único payload: vehículos por estado, consultas por período (últimos 7/30/90
  días), reservas activas/vendidas, test drives agendados/completados,
  consultas abiertas y total de usuarios.
- La página `/admin` (`PanelAdministracion.tsx`) deja de ser solo un menú y se
  convierte en un dashboard: tarjetas de resumen, gráficos simples (barras con
  CSS puro, sin librerías nuevas) para vehículos por estado y consultas por
  período, un selector de período, y conserva los accesos rápidos a las
  secciones de gestión.
- Sin cambios en la base de datos: se agregan solo agregaciones de solo lectura.

## Capabilities

### New Capabilities

- `panel-administracion`: métricas agregadas para el administrador, expuestas
  por un endpoint autenticado con rol `administrador` y presentadas en un
  dashboard con tarjetas de resumen y gráficos simples.

### Modified Capabilities

Ninguna: el comportamiento existente de las capacidades ya implementadas no
cambia; el dashboard agrega una vista nueva sobre los datos actuales.

## Impact

- **Backend**: nuevo `repositories/metricas.go`, `services/metricas.go`,
  `handlers/metricas.go`; registro de ruta en `router.go` dentro del grupo
  protegido por rol `administrador`.
- **Frontend**: reescritura de `frontend/src/pages/PanelAdministracion.tsx`,
  nuevo tipo `frontend/src/types/metricas.ts`, métodos en
  `frontend/src/services/api.ts`; sin dependencias nuevas (gráficos con CSS).
- **Docs**: `docs/roadmap.md` (CU-09 → Implementado).
