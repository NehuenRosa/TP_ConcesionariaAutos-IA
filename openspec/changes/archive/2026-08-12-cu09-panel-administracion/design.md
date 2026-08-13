## Context

- La ruta `/admin` ya existe y está protegida por `RutaProtegida rol="administrador"`, pero `PanelAdministracion.tsx` es solo un menú de enlaces a `/admin/vehiculos` y `/admin/usuarios`.
- No hay librería de gráficos en `frontend/package.json` (solo react, react-dom, react-router), y AGENTS.md prohíbe agregar dependencias fuera del stack sin autorización → gráficos con CSS/SVG propios.
- Los modelos ya existen (`vehiculos`, `consultas`, `reservas`, `turnos_test_drive`, `usuarios`) con los estados documentados en `models/*.go`. No hay migraciones nuevas.
- Arquitectura por capas: `handler → service → repository (GORM)`. La ruta se registra en `router.go` con `middleware.AutenticacionJWT` + `middleware.ExigirRol("administrador")`.
- Consultas con `Fecha string` y `CreatedAt time.Time` (ver `models/consulta.go:24`, `models/turno_test_drive.go:46`); `Reserva` solo tiene `CreatedAt`. Las métricas de "por día" usan la zona horaria del servidor.

## Goals / Non-Goals

**Goals:**
- Un solo endpoint `GET /api/admin/metricas` que responda todo el payload del dashboard (evita N round-trips).
- Gráficos simples (barras) sin dependencias nuevas.
- Cero cambios de esquema: solo lecturas/agregaciones.

**Non-Goals:**
- No se agrega autenticación por período de sesión ni exportación de datos.
- No se meten otras entidades al dashboard (p. ej. mensajes por vendedor).
- No se modifican los estados ni las tablas existentes.

## Decisions

### 1. Un único endpoint agregado por backend (no múltiples endpoints)

Se exponen `GET /api/admin/metricas` (con `periodo` opcional) en vez de
`/api/admin/vehiculos/por-estado`, `/api/admin/consultas/por-periodo`, etc.
Razón: el dashboard consume el payload completo en una sola llamada; es más
simple de cachear, probar y documentar. Alternativa descartada: endpoints
separados → más latencia y más rutas por mantener para el mismo caso de uso.

### 2. Construcción de las series en el repository con GORM

- `vehiculosPorEstado`: `GROUP BY estado` sobre `vehiculos`.
- `consultasPorPeriodo`: `GROUP BY` sobre la fecha de creación filtrado por
  `created_at >= now() - interval '<n> days'` y rellenado en el service con los
  días sin consultas en `0` (garantiza un registro por día, requisito del spec).
  Para que la fecha del día coincida con la zona horaria local del servidor (la
  sesión de Postgres suele estar en UTC), se aplica el desplazamiento de zona
  en SQL: `date(created_at + make_interval(secs => <offset>))::text AS fecha`,
  donde `<offset>` es `time.Now().Zone()` del servidor.
- `reservasActivas` / `reservasVendidas`: `COUNT` por estado.
- `testDrivesAgendados`: `COUNT` con `estado IN (solicitado, confirmado)`.
- `testDrivesCompletados`: `COUNT` con `estado = completado`.
- `consultasAbiertas`: `COUNT` con `estado IN (pendiente, en_conversacion)`.
- `totalUsuarios`: `COUNT` de `usuarios`.

El service valida `periodo ∈ {7, 30, 90}` (default `30`) y devuelve error para
valores inválidos; el handler lo traduce a `400`.

### 3. Respuesta JSON plana con los campos que pide el spec

Payload de `GET /api/admin/metricas`:

```json
{
  "vehiculosPorEstado": [{ "estado": "disponible", "cantidad": 8 }],
  "consultasPorPeriodo": [{ "fecha": "2026-07-13", "cantidad": 2 }],
  "reservasActivas": 1,
  "reservasVendidas": 0,
  "testDrivesAgendados": 2,
  "testDrivesCompletados": 0,
  "consultasAbiertas": 3,
  "totalUsuarios": 4
}
```

Los nombres de campos van en camelCase como el resto de la API.

### 4. Gráficos con CSS/SVG propios en el frontend

Nuevos componentes presentacionales en `frontend/src/components/`:
`GraficoBarras.tsx` (genérico: recibe etiquetas + valores + colores y renderiza
barras con `div`/Tailwind) y `TarjetaMetrica.tsx`. Sin librerías de charts.
Alternativa descartada: recharts/chart.js → fuera del stack permitido.

### 5. Estado de la página con hook propio

`PanelAdministracion.tsx` usa un `useEffect` + estado local
(`datos`, `cargando`, `error`, `periodo`) y el método
`obtenerMetricas(periodo)` de `api.ts`. Se conservan los enlaces de secciones
existentes. Alternativa descartada: React Query → no está en el stack.

## Risks / Trade-offs

- **Zona horaria de las series por día** → Se resuelve en el repository con el
  desplazamiento de zona del servidor (`make_interval(secs => ?)`); el día de la
  serie coincide con la fecha local. Es consistente con el resto del sistema
  (mismo servidor para todas las agregaciones).
- **Volumen grande de datos a futuro** → El dashboard usa `COUNT/GROUP BY` con
  índices ya existentes en `estado` y `created_at`; si crece, se puede cachear
  la respuesta, pero no es necesario hoy.
- **Cambio de período dispara nueva request** → Impacto despreciable con el
  volumen actual; se muestra indicador de carga.

## Migration Plan

- Sin migraciones de base de datos (solo AutoMigrate existente).
- Rollback: quitar la ruta en `router.go` y revertir el cambio de la página;
  no hay datos nuevos que limpiar.

## Open Questions

Ninguna.
