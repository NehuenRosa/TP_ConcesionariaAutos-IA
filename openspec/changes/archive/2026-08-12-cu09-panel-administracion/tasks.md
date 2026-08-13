## 1. Backend: endpoint de métricas

- [x] 1.1 Crear `backend/internal/repositories/metricas.go` con el repository `MetricasRepository` que implementa las agregaciones con GORM: `ContarVehiculosPorEstado`, `ContarConsultasPorDia(desde time.Time)`, `ContarReservasPorEstado`, `ContarTestDrivesAgendados`, `ContarTestDrivesCompletados`, `ContarConsultasAbiertas` y `ContarUsuarios`
- [x] 1.2 Crear `backend/internal/services/metricas.go` con `MetricasService` que valida `periodo ∈ {7, 30, 90}` (default 30), arma las fechas del rango, rellena los días sin consultas con `cantidad 0` y compone el payload de respuesta
- [x] 1.3 Crear `backend/internal/handlers/metricas.go` con `MetricasHandler` y el método `ObtenerMetricas` (parsea `periodo`, delega en el service, traduce errores a `400`/`500` con `{"error": "mensaje en español"}`)
- [x] 1.4 Registrar en `backend/internal/router/router.go` la ruta `GET /api/admin/metricas` dentro de un grupo con `middleware.AutenticacionJWT` + `middleware.ExigirRol("administrador")` e inyectar las dependencias del handler

## 2. Frontend: dashboard del panel

- [x] 2.1 Crear `frontend/src/types/metricas.ts` con los tipos `VehiculoPorEstado`, `ConsultaPorDia` y `Metricas` que reflejan el payload del endpoint
- [x] 2.2 Agregar en `frontend/src/services/api.ts` la función `obtenerMetricas(periodo)` que llama a `GET /api/admin/metricas?periodo=...` con el token JWT
- [x] 2.3 Crear los componentes presentacionales `frontend/src/components/graficos/GraficoBarras.tsx` y `frontend/src/components/graficos/TarjetaMetrica.tsx` (barras con CSS/Tailwind, sin librerías)
- [x] 2.4 Reescribir `frontend/src/pages/PanelAdministracion.tsx`: tarjetas de resumen, gráfico de vehículos por estado, gráfico de consultas por período, selector de período (7/30/90 días) que recarga las métricas, estados de carga/error y accesos rápidos existentes a vehículos y usuarios

## 3. Verificación

- [x] 3.1 Build backend: `go build ./...` y `go vet ./...` en `backend/` sin errores
- [x] 3.2 Build frontend: `npm.cmd run build` en `frontend/` sin errores
- [x] 3.3 E2E con PostgreSQL local (Docker + `go run ./cmd/api`): login de administrador, `200` en `GET /api/admin/metricas` con payload completo, `400` con `periodo` inválido, `403` con rol vendedor/cliente y `401` sin token
- [x] 3.4 Actualizar `docs/roadmap.md` marcando CU-09 como Implementado
