## 1. Backend: repository de vehículos (gestión)

- [x] 1.1 Extender `internal/repositories/vehiculos.go` con `ListarParaGestion(ctx, estado, pagina, tamano)` donde `estado` vacío significa "todos los estados", con `Count` del total y `Preload("Imagenes")`.
- [x] 1.2 Agregar `Crear(ctx, vehiculo)` que persista el vehículo y sus imágenes.
- [x] 1.3 Agregar `Actualizar(ctx, vehiculo)` que actualice la ficha/estado y reemplace las imágenes (borrar las existentes e insertar las nuevas).
- [x] 1.4 Agregar `DarDeBaja(ctx, id)` que actualice `estado` a `dado_de_baja`.

## 2. Backend: service de vehículos (gestión)

- [x] 2.1 Ampliar `VehiculoService` en `internal/services/vehiculos.go` con `ListarParaGestion`, `ObtenerParaGestion`, `Crear`, `Actualizar` y `DarDeBaja`.
- [x] 2.2 Implementar validaciones: campos requeridos, `anio` en rango, `precio` positivo, `condicion` en `{nuevo, usado}`, `estado` en los cuatro conocidos y filtro de estado válido.
- [x] 2.3 Mapear `gorm.ErrRecordNotFound` a `ErrVehiculoNoEncontrado` en las operaciones de gestión.

## 3. Backend: handler y rutas de gestión

- [x] 3.1 Crear `internal/handlers/vehiculos_gestion.go` con `VehiculoGestionHandler` y DTOs (`VehiculoEntrada`, `VehiculoGestionResumen`, `RespuestaPaginadaGestion`).
- [x] 3.2 Implementar `Listar`, `ObtenerDetalle`, `Crear`, `Actualizar` y `DarDeBaja` delegando en el service y respondiendo `400`/`404`/`500` con mensajes en español.
- [x] 3.3 Registrar en `internal/router/router.go` el grupo `/api/admin/vehiculos` con `middleware.AutenticacionJWT(configuracion.JWTSecreto)` y `middleware.ExigirRol("administrador")`.

## 4. Frontend: tipos y cliente HTTP

- [x] 4.1 Agregar a `src/types/vehiculo.ts` los tipos `VehiculoEntrada` (payload de alta/edición) y `PaginaVehiculosGestion`.
- [x] 4.2 Agregar a `src/services/api.ts` los métodos `listarVehiculosGestion(pagina, tamano, estado)`, `obtenerVehiculoGestion(id)`, `crearVehiculo(datos)`, `actualizarVehiculo(id, datos)` y `darDeBajaVehiculo(id)`.

## 5. Frontend: páginas de gestión

- [x] 5.1 Crear `src/pages/AdminVehiculos.tsx` con listado paginado, filtro por estado, estados de carga/vacío/error y acciones editar/baja en español.
- [x] 5.2 Crear `src/pages/FormularioVehiculo.tsx` reutilizable para alta y edición (ficha técnica + URLs de imágenes), con envío a la API y retorno al listado.
- [x] 5.3 Registrar las rutas `/admin/vehiculos`, `/admin/vehiculos/nuevo` y `/admin/vehiculos/:id/editar` en `src/routes/Rutas.tsx` dentro del layout base y enlazarlas desde el panel de administración.

## 6. Verificación

- [x] 6.1 `cd backend && go build ./...` compila sin errores.
- [x] 6.2 `cd backend && go vet ./...` sin errores.
- [x] 6.3 `cd frontend && npm run build` compila sin errores.
- [x] 6.4 Verificar de punta a punta con `docker compose up` que el ABM funciona (alta, listado, detalle, modificación y baja lógica) y que un vehículo `dado_de_baja` desaparece del catálogo público.

