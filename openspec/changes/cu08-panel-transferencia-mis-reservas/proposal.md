# Proposal: cu08-panel-transferencia-mis-reservas

## Why

Los datos para pagar la seña (CBU, alias, monto y plazo) solo se muestran en el
formulario de reserva. Una vez creada la reserva, el cliente no tiene dónde
volver a verlos desde Mis Reservas: si cierra esa pantalla o le avisan por
WhatsApp, queda sin referencia para hacer la transferencia dentro de las 2
horas. Además, con `CBU_CONCESIONARIA`/`ALIAS_CONCESIONARIA` vacíos el panel
muestra un aviso genérico y en el entorno local no había valores cargados.

## What Changes

- Mis Reservas: las reservas **activas** muestran un panel "Seña por
  transferencia (5 %)" con CBU, alias, monto a transferir y vencimiento del
  plazo (reutilizando el componente de datos que ya existe en el formulario);
  incluye el acceso directo a subir/reenviar comprobante.
- `.env.example`: documentadas las claves `CBU_CONCESIONARIA` /
  `ALIAS_CONCESIONARIA` con valores de ejemplo.
- Sin cambios de API: usa el endpoint existente
  `GET /api/reservas/datos-transferencia?vehiculoId=`.

## Capabilities

### Modified Capabilities

- `reserva-vehiculo`: la reserva activa mantiene visibles los datos de pago de
  la seña durante todo el plazo desde Mis Reservas.
