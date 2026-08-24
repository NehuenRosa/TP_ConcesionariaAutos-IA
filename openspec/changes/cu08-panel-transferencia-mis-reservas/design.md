# Design: cu08-panel-transferencia-mis-reservas

## Contexto

`FormularioReserva.tsx` ya resuelve todo: llama
`GET /reservas/datos-transferencia?vehiculoId=` y renderiza el componente local
`DatosSenia` (CBU/alias/monto, o aviso ámbar si la concesionaria no cargó sus
datos bancarios). `MisReservas.tsx` no muestra nada de eso hoy. El backend no
requiere cambios.

## Objetivos

- Que el cliente vea siempre dónde transferir mientras su reserva esté activa,
  sin depender de haber guardado la pantalla del formulario.

## Decisiones

### D1: reutilizar `DatosSenia`

Extraer el componente a `components/reserva/DatosSenia.tsx` e importarlo desde
el formulario y desde Mis Reservas (misma fuente visual, cero duplicación).

### D2: fetch perezoso por reserva activa

Mis Reservas obtiene los datos una sola vez por vehículo con reserva activa
(clave por `vehiculoId`, caché en un `Map` en estado) al montar; evita N
requests idénticos cuando hay varias reservas del mismo vehículo. Si falla el
fetch, el panel muestra igualmente monto/vencimiento calculados de la reserva y
el aviso de que el personal pasa los datos.

### D3: dónde ubicarlo

Dentro de la tarjeta de cada reserva activa, debajo del estado y antes de la
zona de comprobante: primero el "dónde pago", después el "ya pagué" (subir
comprobante). Las reservas canceladas/vendidas no lo muestran.

### Riesgos

- Ninguno funcional: endpoint ya existente y probado; solo cambia la vista.
