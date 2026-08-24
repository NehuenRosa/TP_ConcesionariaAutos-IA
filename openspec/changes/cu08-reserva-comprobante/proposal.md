# Proposal: cu08-reserva-comprobante

## Why

Hoy cualquier cliente puede reservar una unidad sin ningún compromiso
económico, lo que genera reservas fantasma: unidades bloqueadas del catálogo
por clientes sin intención real de compra. Exigir una seña del 5 % del valor
del vehículo vía transferencia bancaria, con comprobante obligatorio dentro
de las 2 horas, filtra esos casos y libera automáticamente la unidad cuando
el cliente no completa el pago.

## What Changes

- Al crear una reserva (`POST /api/reservas`) el sistema registra un plazo de
  **2 horas** para que el cliente envíe el comprobante de la transferencia.
- Nuevo endpoint público-autenticado `GET /api/reservas/datos-transferencia`
  que devuelve el CBU/alias de la concesionaria (variables de entorno) y el
  **monto de la seña: 5 % del precio del vehículo** (constante en código).
- Nuevo endpoint `POST /api/reservas/:id/comprobante` (multipart): el cliente
  sube una imagen del comprobante (JPG/PNG/WebP, máx. 5 MB) dentro del plazo;
  se guarda en la base de datos y queda disponible para el vendedor.
- Nuevo endpoint `GET /api/reservas/:id/comprobante`: el dueño de la reserva
  y los vendedores pueden visualizar la imagen.
- **Expiración automática**: si pasan 2 horas sin comprobante, la reserva se
  anula sola (estado `cancelada`) y el vehículo vuelve a `disponible`
  (reaparece en el catálogo). Verificación periódica (job interno cada 30 s)
  más chequeos perezosos al operar sobre la reserva.
- El vendedor verifica el comprobante manualmente desde su bandeja de
  reservas y recién entonces confirma la venta o cancela la reserva.
- Los estados de reserva no cambian (`activa`, `vendida`, `cancelada`); el
  "pendiente de comprobante" se deriva de reserva activa sin comprobante.
- Frontend: la página de reserva muestra CBU/alias/monto antes de confirmar;
  después de reservar muestra cuenta regresiva y subida de comprobante;
  `/mis-reservas` avisa el plazo y permite subirlo; `/vendedor/reservas`
  permite ver el comprobante antes de confirmar la venta.

## Capabilities

### New Capabilities

_(ninguna: es una extensión del flujo existente de reservas)_

### Modified Capabilities

- `reserva-vehiculo`: la creación de reserva incorpora plazo de 2 h y monto
  de seña (5 %); nuevos requisitos de datos de transferencia, envío y
  consulta del comprobante, y expiración automática sin comprobante; las
  páginas de reserva, mis reservas y gestión del vendedor se adaptan al nuevo
  flujo.

## Impact

- **Backend**
  - `internal/models/reserva.go`: campos `VencimientoComprobante` y
    `ComprobanteEnviadoAt`; nuevo modelo `ComprobanteReserva` (bytes + MIME,
    1:1 con reserva).
  - `internal/services/reservas.go`: constantes `PorcentajeSena` (5 %) y
    `PlazoComprobante` (2 h); métodos nuevos `ObtenerDatosTransferencia`,
    `SubirComprobante`, `ObtenerComprobante`, `ExpirarVencidas`.
  - `internal/repositories/reservas.go`: persistencia del comprobante y
    expiración en lote (transacción).
  - `internal/handlers/reservas.go`: endpoints nuevos (multipart siguiendo el
    patrón del chatbot).
  - `cmd/api/main.go`: job periódico de expiración (goroutine + ticker).
  - `internal/config/config.go`: `CbuConcesionaria`, `AliasConcesionaria`.
- **Frontend**
  - `types/reserva.ts`, `services/api.ts` (multipart autenticado genérico),
    `FormularioReserva.tsx`, `MisReservas.tsx`, `GestionReservas.tsx`.
- **Infra / config**: `.env.example` y `docker-compose.yml` con
  `CBU_CONCESIONARIA` / `ALIAS_CONCESIONARIA`. Sin dependencias nuevas.
- **API**: dos endpoints nuevos + campos agregados a respuestas existentes;
  ningún contrato vigente cambia de forma incompatible.
- **Docs**: CU-08 en `docs/roadmap.md` pasa a "Implementado (con seña)".
