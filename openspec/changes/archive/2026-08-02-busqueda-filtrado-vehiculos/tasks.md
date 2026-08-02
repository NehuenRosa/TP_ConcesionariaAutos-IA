## 1. Base de datos y modelo

- [x] 1.1 Agregar el campo `Tipo` a `models.Vehiculo` con `json:"tipo"` y
      validarlo en el seed de la base de datos.
- [x] 1.2 Actualizar `backend/internal/database/database.go` (auto-migración) si
      corresponde y poblar `tipo` en los vehículos de prueba del seed.

## 2. Backend: filtros, búsqueda y ordenamiento

- [x] 2.1 Definir en `services/vehiculos.go` el struct `FiltrosBusqueda` con los
      criterios opcionales: busqueda, marca, modelo, anioMin, anioMax,
      precioMin, precioMax, tipo, combustible, condicion, ordenPor y
      ordenDireccion.
- [x] 2.2 Agregar el error `ErrFiltroInvalido` y la validación semántica en el
      service (rangos válidos, `condicion` en nuevo/usado, `ordenPor` en
      precio/anio y `ordenDireccion` en asc/desc).
- [x] 2.3 Cambiar la firma de `ListarDisponibles` para recibir
      `FiltrosBusqueda` y propagar el error `ErrFiltroInvalido`.
- [x] 2.4 Actualizar `VehiculoRepository` y `repositories/vehiculos.go` para
      construir la consulta GORM dinámica: `Where` encadenados por criterio
      presente, búsqueda `ILIKE` sobre marca o modelo con escape de comodines,
      rangos inclusivos y `ORDER BY` con whitelist de columnas/direcciones.
- [x] 2.5 Parsear en `handlers/vehiculos.go` los query params de búsqueda,
      filtros y ordenamiento en `Listar` y traducir `ErrFiltroInvalido` a
      `400` con mensaje en español.
- [x] 2.6 Verificar con `go build ./...` y `go vet ./...` en `backend/`.

## 3. Backend: campo tipo en la ficha y el ABM

- [x] 3.1 Incluir `Tipo` en `EntradaVehiculo`, `validarEntrada` (no vacío) y
      `aModelo` en `services/vehiculos.go`.
- [x] 3.2 Agregar `Tipo` a `VehiculoEntrada`, `aEntrada` y, si aplica, al resumen
      de gestión en `handlers/vehiculos_gestion.go`.
- [x] 3.3 Confirmar que el detalle público (`GET /api/vehiculos/:id`) expone
      `tipo` en la respuesta JSON.

## 4. Frontend: tipos y cliente API

- [x] 4.1 Agregar `tipo: string` a `Vehiculo`, `VehiculoEntrada` y
      `ResumenVehiculo` en `frontend/src/types/vehiculo.ts`.
- [x] 4.2 Definir el tipo `FiltrosVehiculos` (con los mismos campos opcionales
      que el backend) y el tipo `OrdenVehiculos` en `types/vehiculo.ts`.
- [x] 4.3 Actualizar `api.listarVehiculos` en `services/api.ts` para aceptar
      filtros y construir la query string con `URLSearchParams`, omitiendo los
      parámetros vacíos.

## 5. Frontend: panel de búsqueda y filtros en /catalogo

- [x] 5.1 Agregar en `pages/Catalogo.tsx` un estado `filtros` y un formulario
      con: campo de texto libre, selects de marca, modelo, tipo, combustible y
      condición, inputs de rango de año y precio, y control de ordenamiento.
- [x] 5.2 Al aplicar o modificar filtros, volver a la página 1 y refrescar el
      listado llamando a `api.listarVehiculos` con los filtros.
- [x] 5.3 Usar debounce (o `useDeferredValue`) en el campo de búsqueda de texto
      para evitar una petición por tecla.
- [x] 5.4 Mostrar mensaje en español cuando no hay resultados para los filtros
      aplicados y un control para limpiar todos los filtros.
- [x] 5.5 Opcional: mostrar el tipo como etiqueta en las tarjetas del listado.

## 6. Frontend: tipo en formulario administrativo y detalle

- [x] 6.1 Agregar el campo `tipo` (select con valores sugeridos: sedán, SUV,
      hatchback, pick-up, coupe) en `pages/FormularioVehiculo.tsx` para alta y
      edición.
- [x] 6.2 Mostrar el tipo en `pages/DetalleVehiculo.tsx`.
- [x] 6.3 Verificar con `npm run build` (o `npm run typecheck` si existe) en
      `frontend/`.

## 7. Verificación integral

- [x] 7.1 Probar manualmente con el backend y la base levantados: búsqueda por
      texto, filtros combinados, rangos inválidos (400), orden por precio/año y
      combinación de búsqueda + filtros + orden + paginación.
- [x] 7.2 Confirmar que el alta y la edición administrativa guardan y muestran
      el `tipo`, y que la ficha pública lo incluye.
- [x] 7.3 Ejecutar `openspec validate --strict` sin errores.
