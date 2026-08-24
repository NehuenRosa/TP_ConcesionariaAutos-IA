# Tasks: cu08-reserva-comprobante

## 1. Backend — modelo y configuración

- [x] 1.1 Extender `models/reserva.go`: `VencimientoComprobante time.Time` y `ComprobanteEnviadoAt *time.Time` en `Reserva`; nuevo modelo `ComprobanteReserva` (ReservaID uniqueIndex, MIME, Datos []byte)
- [x] 1.2 Agregar `ComprobanteReserva` al AutoMigrate de `internal/database/database.go`
- [x] 1.3 Agregar a `config.go`: `CbuConcesionaria` / `AliasConcesionaria` desde `CBU_CONCESIONARIA`/`ALIAS_CONCESIONARIA`; actualizar `.env.example` y `docker-compose.yml`

## 2. Backend — repositorio

- [x] 2.1 Agregar a `repositories/reservas.go`: `GuardarComprobante` (upsert por ReservaID), `ObtenerComprobante(reservaID)` y `ExpirarVencidas(ctx)` transaccional que cancela activas sin comprobante vencidas y libera sus vehículos a `disponible`, excluyendo reservas históricas con vencimiento cero
- [x] 2.2 Incluir el comprobante en los preloads necesarios o exponer existencia vía join para saber si la reserva tiene comprobante enviado

## 3. Backend — servicio

- [x] 3.1 Constantes en `services/reservas.go`: `PorcentajeSena = 0.05`, `PlazoComprobante = 2 * time.Hour`, peso/formatos máximos del comprobante
- [x] 3.2 Extender `Crear` para fijar `VencimientoComprobante` y devolver monto de seña; método `ObtenerDatosTransferencia(ctx, clienteID, vehiculoID)` con validación de vehículo disponible
- [x] 3.3 Métodos `SubirComprobante(ctx, reservaID, clienteID, archivo)` (valida dueño, estado activa, no vencida-sin-comprobante; aplica expiración perezosa si corresponde) y `ObtenerComprobante(ctx, reservaID, usuario)` con permisos dueño/vendedor/administrador
- [x] 3.4 Chequeo perezoso de expiración en `ConfirmarVenta` y `CancelarComoVendedor`; método `ExpirarVencidas(ctx)` delegando en el repositorio
- [x] 3.5 Tests del service con fakes: creación fija vencimiento a 2 h, subida válida registra envío, subida tardía → error 409-lógico, expiración en lote libera vehículos, confirmación sobre vencida primero expira

## 4. Backend — handler y rutas

- [x] 4.1 Handler `DatosTransferencia` (`GET /api/reservas/datos-transferencia`) con `400` sin vehiculoId y `404` vehículo no disponible
- [x] 4.2 Handler `SubirComprobante` (multipart campo `comprobante`, misma disciplina de validación que el chatbot: JPG/PNG/WebP ≤ 5 MB; errores `400`/`404`/`409`) y handler que sirva la imagen con el MIME correcto (`GET /api/reservas/:id/comprobante`)
- [x] 4.3 Registrar rutas nuevas en `router.go` respetando autenticación JWT y rol vendedor para la imagen cuando no sea el dueño
- [x] 4.4 Agregar a las respuestas existentes (`ReservaResumen`) los campos `montoSenia`, `vencimientoComprobante` y `comprobanteEnviadoAt`
- [x] 4.5 Job periódico en `cmd/api/main.go`: goroutine con ticker de 30 s llamando `ExpirarVencidas`

## 5. Frontend

- [x] 5.1 Actualizar `types/reserva.ts` (campos nuevos + `DatosTransferencia`) y generalizar `peticionMultipart` en `api.ts` agregando `obtenerDatosTransferencia`, `subirComprobanteReserva` y `obtenerComprobanteReserva` (blob)
- [x] 5.2 `FormularioReserva.tsx`: mostrar CBU/alias/monto antes de confirmar; tras confirmar, éxito con cuenta regresiva de 2 h y subida del comprobante con feedback en español
- [x] 5.3 `MisReservas.tsx`: indicador pendiente/enviado con plazo restante, botón subir comprobante en activas pendientes
- [x] 5.4 `GestionReservas.tsx`: indicador de comprobante (pendiente/enviado + hora), acción "Ver comprobante" para activas con envío

## 6. Verificación y documentación

- [x] 6.1 `go build ./...`, `go vet ./...`, tests backend en verde; `npm run build` y tests frontend en verde
- [ ] 6.2 Prueba E2E manual: reservar → ver CBU/monto → subir comprobante dentro del plazo → vendedor lo ve y confirma venta; y flujo alternativo: reservar sin subir → forzar vencimiento → unidad vuelve a `disponible` en catálogo y la reserva queda `cancelada`
- [x] 6.3 Actualizar `docs/roadmap.md` (CU-08 con seña) y AGENTS.md (reglas de negocio de la seña)
