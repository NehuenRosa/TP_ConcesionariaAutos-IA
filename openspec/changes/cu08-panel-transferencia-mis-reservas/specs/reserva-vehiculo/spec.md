# Delta spec: reserva-vehiculo

## MODIFIED Requirements

### Requirement: Datos de transferencia visibles durante la reserva

Mientras una reserva está activa, el cliente SHALL poder ver en Mis Reservas
los datos para pagar la seña (CBU y alias de la concesionaria, monto del 5 %
y vencimiento del plazo) para el vehículo reservado, además de los datos que
ya muestra el formulario de reserva. Si la concesionaria no configuró
`CBU_CONCESIONARIA`/`ALIAS_CONCESIONARIA`, la vista SHALL indicar que el
personal le hará llegar los datos, sin bloquear la subida del comprobante.

#### Scenario: Cliente vuelve a Mis Reservas dentro del plazo

- **WHEN** el cliente abre Mis Reservas con una reserva activa
- **THEN** ve CBU/alias (o el aviso si no están configurados), el monto de la
  seña y el tiempo restante, junto con el acceso para subir el comprobante

#### Scenario: Reserva no activa

- **WHEN** la reserva está cancelada o vendida
- **THEN** no se muestra el panel de datos de transferencia
