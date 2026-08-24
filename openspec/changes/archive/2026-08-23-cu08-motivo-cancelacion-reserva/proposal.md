# Proposal: cu08-motivo-cancelacion-reserva

## Why

Hoy el vendedor puede ver el comprobante de una reserva y confirmar la venta o
cancelarla, pero al cancelar no deja constancia del porqué: el cliente ve su
reserva `cancelada` sin entender qué falló (¿comprobante ilegible?, ¿monto
incorrecto?, ¿transferencia nunca acreditada?). Registrar un motivo obligatorio
cierra el circuito de validación del comprobante con transparencia para el
cliente.

## What Changes

- Modelo `Reserva`: nuevo campo `MotivoCancelacion` (texto, opcional).
- `PUT /api/reservas/:id/cancelar` acepta cuerpo JSON opcional `{motivo}`.
  Cuando cancela un vendedor, el motivo es **obligatorio** (`400` si falta);
  la cancelación propia del cliente sigue sin motivo.
- El servicio valida el motivo (no vacío tras recortar) y lo persiste junto al
  cambio de estado atómico existente.
- `ReservaResumen` incluye `motivoCancelacion` (omitempty).
- Frontend: en `/vendedor/reservas`, cancelar abre un formulario que exige el
  motivo; en `/mis-reservas` el cliente ve el motivo bajo el estado
  `cancelada`.
- La pantalla de reserva ya muestra CBU/alias/monto antes de confirmar y el
  vendedor ya puede ver el comprobante: este change completa esa validación
  con la denegación explicada.

## Capabilities

### New Capabilities

*(ninguna)*

### Modified Capabilities

- `reserva-vehiculo`: cancelación del vendedor con motivo obligatorio visible.

## Impact

- Backend: `models/reserva.go`, `services/reservas.go` (+tests),
  `handlers/reservas.go`.
- Frontend: `types/reserva.ts`, `services/api.ts`, `GestionReservas.tsx`,
  `MisReservas.tsx`.
- Migración automática (columna nueva nullable), sin cambios de rutas nuevas.
