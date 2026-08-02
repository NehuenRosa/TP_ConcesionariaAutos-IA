## 1. Backend: repository de vehículos

- [ ] 1.1 Crear `internal/repositories/vehiculos.go` con `ListarDisponibles(pagina, tamano)` y `ObtenerPorID(id)` sobre GORM.
- [ ] 1.2 Implementar la respuesta paginada con total de registros.

## 2. Backend: service de vehículos

- [ ] 2.1 Crear `internal/services/vehiculos.go` con la lógica de negocio (filtrar estado `disponible`, validar paginación).
- [ ] 2.2 Retornar `error` descriptivo en español cuando un vehículo no existe o no está disponible.

## 3. Backend: handlers y rutas

- [ ] 3.1 Crear `internal/handlers/vehiculos.go` con `Listar` y `ObtenerDetalle` (parsing de request, delegación al service, respuestas JSON).
- [ ] 3.2 Registrar en `internal/router/router.go` las rutas públicas `GET /api/vehiculos` y `GET /api/vehiculos/:id`.
- [ ] 3.3 Responder `404` con mensaje en español cuando el vehículo no existe o no está disponible.

## 4. Frontend: tipos y cliente HTTP

- [ ] 4.1 Crear `src/types/vehiculo.ts` con la ficha de vehículo y la respuesta paginada.
- [ ] 4.2 Agregar a `src/services/api.ts` los métodos `listarVehiculos(pagina, tamano)` y `obtenerVehiculo(id)`.

## 5. Frontend: páginas públicas

- [ ] 5.1 Crear `src/pages/Catalogo.tsx` con el listado paginado consumiendo el endpoint público.
- [ ] 5.2 Crear `src/pages/DetalleVehiculo.tsx` con la ficha técnica y la galería.
- [ ] 5.3 Registrar las rutas `/catalogo` y `/catalogo/:id` dentro del layout base en `src/routes`.
- [ ] 5.4 Manejar estados de carga, vacío y error con mensajes en español.

## 6. Verificación

- [ ] 6.1 `cd backend && go build ./...` compila sin errores.
- [ ] 6.2 `cd frontend && npm run build` compila sin errores.
- [ ] 6.3 Verificar con `docker compose up` que `/api/vehiculos` lista solo vehículos disponibles y que `/catalogo` funciona de punta a punta.
