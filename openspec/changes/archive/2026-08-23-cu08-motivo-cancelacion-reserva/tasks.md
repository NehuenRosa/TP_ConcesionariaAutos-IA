# Tasks: cu08-motivo-cancelacion-reserva

## 1. Backend

- [x] 1.1 `models/reserva.go`: campo `MotivoCancelacion string` (text) con json `motivoCancelacion,omitempty`
- [x] 1.2 `services/reservas.go`: `CancelarComoVendedor` exige motivo no vacío (error `ErrMotivoRequerido`); persistir el motivo en la transacción existente
- [x] 1.3 Handler `PUT /reservas/:id/cancelar`: parsear cuerpo opcional `{motivo}` y pasarlo al servicio; `400` si vendedor cancela sin motivo
- [x] 1.4 `ReservaResumen` + `aReservaResumen`: exponer `motivoCancelacion`
- [x] 1.5 Tests: cancelación de vendedor sin motivo → error; con motivo queda guardado; cliente sigue cancelando sin motivo

## 2. Frontend

- [x] 2.1 `types/reserva.ts` (`motivoCancelacion?: string`) y `api.ts` (`cancelarReservaVendedor(id, motivo)` con cuerpo JSON)
- [x] 2.2 `GestionReservas.tsx`: formulario inline obligatorio de motivo al rechazar/cancelar, con confirmación
- [x] 2.3 `MisReservas.tsx`: bloque con el motivo cuando la reserva está `cancelada` y tiene motivo

## 3. Verificación y documentación

- [x] 3.1 `go build ./...`, `go vet ./...`, tests backend en verde; `npm run build` y tests frontend en verde
- [x] 3.2 Prueba E2E manual: reservar → subir comprobante → vendedor lo ve, rechaza con motivo → cliente lee el motivo en Mis Reservas (verificada por API contra el stack Docker: 400 sin motivo, motivo visible en mis-reservas, unidad vuelve a disponible)
- [x] 3.3 Actualizar `docs/api.md` (cuerpo del endpoint) si corresponde
