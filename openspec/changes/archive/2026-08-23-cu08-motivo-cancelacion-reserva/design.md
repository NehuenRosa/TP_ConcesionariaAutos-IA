# Design: cu08-motivo-cancelacion-reserva

## Contexto

CU-08 ya cubre seña, comprobante y expiración. El único hueco es la
denegación: el vendedor cancela sin dejar explicación. Se resuelve con un
campo nuevo y validación en el camino existente, sin endpoints ni estados
nuevos.

## Decisiones

- **D1 — Campo único `MotivoCancelacion`.** Texto nullable en `reservas`.
  No se agrega estado `rechazada`: "cancelada + motivo" distingue por el
  motivo mismo (presente = decidió la concesionaria; ausente = decidió el
  cliente o expiró). Menos estados, misma información.
- **D2 — Obligatorio solo para el vendedor.** `CancelarComoVendedor(ctx, id,
  motivo)` exige texto no vacío tras recortar (`400` lógico → handler).
  La cancelación del cliente no cambia (sin motivo). El job de expiración
  tampoco escribe motivo.
- **D3 — Cuerpo JSON opcional en el endpoint existente.** `PUT /reservas/:id/
  cancelar` parsea `{motivo}` si viene; si no viene cuerpo se trata como
  vacío. El rol ya llega por JWT: el handler decide exigirlo según rol
  vendedor/administrador.
- **D4 — DTO.** `ReservaResumen.MotivoCancelacion string json:"motivoCancelacion,omitempty"`.
- **D5 — Frontend.** En `GestionReservas`, botón "Rechazar comprobante" /
  "Cancelar reserva" abre formulario inline que exige motivo (mínimo razonable,
  contador de caracteres) y confirma. En `MisReservas`, bloque destacado con
  el motivo cuando `estado === 'cancelada' && motivoCancelacion`.
- **D6 — Sin longitud máxima nueva:** reutiliza límites de texto libre;
  validación mínima de presencia en backend y frontend.

## Riesgos

- Clientes viejos del frontend (caché) cancelando como vendedor sin motivo:
  el backend responde `400` con mensaje claro; la UI pide el motivo.

## Migración

AutoMigrate agrega la columna; reservas históricas quedan con motivo vacío.
