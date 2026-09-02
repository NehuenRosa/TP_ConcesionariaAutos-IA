# Design: cu08-reserva-comprobante

## Context

CU-08 hoy: `POST /api/reservas` crea reserva `activa` y pasa el vehículo a
`reservado` de forma atómica (`CrearYReservar` con `RowsAffected`); el
vendedor confirma la venta o cancela liberando la unidad. No hay montos,
plazos ni archivos. El único manejo multipart del backend es la tasación del
chatbot (whitelist de extensiones, 5 MB, lectura en memoria). No existe
ningún job periódico; `main.go` es secuencial. Los contenedores corren
`read_only`, lo que descarta disco local sin volúmenes nuevos.

## Goals / Non-Goals

**Goals:**
- Seña obligatoria del 5 % con comprobante dentro de 2 h; expiración
  automática que libera la unidad.
- Vendedor revisa el comprobante y decide (confirmar venta / cancelar).
- CBU/alias por entorno; porcentaje fijo en código.

**Non-Goals:**
- Pasarelas de pago electrónicas o validación bancaria automática: el
  comprobante es una imagen verificada por personas.
- Estados de reserva adicionales (se reutilizan los tres existentes).
- Notificaciones push/mail al vendedor o al cliente.
- Múltiples cuotas/parciales de la seña; un solo comprobante vigente.

## Decisions

### D1: Sin estados nuevos — "pendiente de comprobante" derivado

`Reserva` suma dos campos:

```go
VencimientoComprobante time.Time  `json:"vencimientoComprobante"` // fijado en Crear: now + 2h
ComprobanteEnviadoAt   *time.Time `json:"comprobanteEnviadoAt,omitempty"`
```

"pendiente de comprobante" = `estado == activa && ComprobanteEnviadoAt == nil`;
"expirada" = pendiente && `now > VencimientoComprobante`. Así no cambian los
estados, los filtros de la bandeja, las métricas del panel admin
(`ContarReservasActivas/Vendidas`) ni el contrato TS existente (solo se
agregan campos).

### D2: Comprobante en PostgreSQL como bytes

Nueva entidad 1:1:

```go
type ComprobanteReserva struct {
    ID        uint      `gorm:"primaryKey"`
    ReservaID uint      `gorm:"not null;uniqueIndex"`
    MIME      string    `gorm:"not null"`
    Datos     []byte    `gorm:"not null"` // bytea
    CreatedAt time.Time `json:"-"`
}
```

- **Por qué BD**: contenedores read_only (evita volumen nuevo), backup único,
  cero dependencias. El tamaño acotado (1 imagen ≤ 5 MB por reserva) hace
  inocuo el costo en la tabla.
- Reenvío mientras esté activa = upsert (reemplaza bytes/MIME y actualiza
  `ComprobanteEnviadoAt`).
- La imagen viaja por `GET /api/reservas/:id/comprobante` con `Content-Type`
  del MIME guardado; acceso: dueño de la reserva o vendedor/administrador
  (`403` a otros clientes, `404` si nunca se envió).

### D3: Validación del archivo replicando el patrón del chatbot

Misma disciplina que `handlers/chatbot.go`/`services/chatbot.go`: campo
multipart `comprobante`, whitelist `.jpg/.jpeg/.png/.webp`, límite 5 MB,
`io.ReadAll` acotado. Constantes propias en `services/reservas.go`
(`MaximoPesoComprobanteBytes`) para no acoplar al chatbot.

### D4: Expiración — job periódico + chequeos perezosos

- Job: goroutine en `main.go` con `time.NewTicker(30 * time.Second)` que
  llama `servicioReservas.ExpirarVencidas(ctx)`.
- `ExpirarVencidas`: una transacción que cancela en lote
  `UPDATE reservas SET estado='cancelada' WHERE estado='activa' AND
  comprobante_enviado_at IS NULL AND vencimiento_comprobante < now()`
  y libera los vehículos afectados a `disponible`. Idempotente.
- Chequeo perezoso: `SubirComprobante`, `ConfirmarVenta` y
  `CancelarComoVendedor` evalúan si la reserva activa está vencida antes de
  operar; si lo está, aplican la expiración primero (y responden `409` en el
  caso de la subida tardía). Esto cubre el intervalo entre ticks y hace el
  comportamiento correcto aunque el job esté demorado.
- Carrera job vs subida simultánea: la cancelación en lote excluye filas con
  comprobante ya registrado y la subida valida estado+vencimiento sobre la
  misma fila; ante conflicto gana quien persista primero y el otro ve `409`.
- **Alternativa descartada**: `UPDATE ... RETURNING` con CTE atómico por
  reserva (más hermético, pero mucho menos legible en GORM y el patrón
  transaccional existente ya resuelve la condición de carrera).

### D5: Monto y datos de transferencia calculados en el backend

- Constantes en `services/reservas.go`: `PorcentajeSena = 0.05`,
  `PlazoComprobante = 2 * time.Hour`.
- Config nueva: `CBU_CONCESIONARIA`, `ALIAS_CONCESIONARIA` (vacías =
  endpoint devuelve cadenas vacías y el frontend muestra "la concesionaria te
  va a pasar los datos"; el flujo sigue funcionando).
- `GET /api/reservas/datos-transferencia?vehiculoId=` (autenticado) responde
  `{cbu, alias, monto}`; monto = precio actual × PorcentajeSena, redondeo a
  2 decimales. El frontend NUNCA calcula el monto por su cuenta.

### D6: Frontend

- `api.ts`: variante genérica `peticionMultipart<T>(ruta, campos)` (token
  siempre si existe) junto a la específica de tasación, que queda delegando
  en ella; funciones nuevas `obtenerDatosTransferencia(vehiculoId)`,
  `subirComprobanteReserva(id, archivo)`, `obtenerComprobanteReserva(id)`
  (devuelve Blob para `<img>` con `URL.createObjectURL`).
- `FormularioReserva.tsx`: antes de confirmar muestra CBU/alias/monto
  (endpoint D5) e instrucciones; después de confirmar muestra éxito +
  cuenta regresiva (interval de 1 s contra `vencimientoComprobante`) y subida
  inmediata del comprobante.
- `MisReservas.tsx`: badge "comprobante enviado"/"pendiente" con plazo
  restante; botón subir cuando falta; el countdown se recalcula al renderizar
  (sin interval global, suficiente para esta vista).
- `GestionReservas.tsx`: columna/indicador de comprobante (pendiente/enviado
  con hora), botón "Ver comprobante" (abre la imagen) visible para activas
  con envío; acciones existentes intactas.

## Risks / Trade-offs

- [Job caído o backend reiniciado durante el tick] → los chequeos perezosos
  garantizan la expiración lógica al primer uso; el job solo la adelanta. Al
  arrancar, el primer tick barre rezagos.
- [`bytea` de hasta 5 MB por fila] → volumen bajo (una reserva por unidad);
  si creciera, migrar a object storage queda contenido detrás del repositorio.
- [Cliente sube comprobante ilegible] → el vendedor puede pedir reenvío
  (el cliente vuelve a subir mientras esté activa) o cancelar.
- [`now` del servidor vs cliente] → el plazo lo fija y juzga exclusivamente
  el backend; el countdown del frontend es informativo.
- [CBU vacío en algún despliegue] → degradación explícita (texto indicativo),
  sin romper el flujo de reserva.
- [Confirmar venta sin comprobante dentro de las 2 h] → permitido
  deliberadamente (pago en efectivo en el local); la UI destaca el faltante
  pero la API no bloquea: decisión de negocio documentada.

## Migration Plan

1. AutoMigrate agrega columnas y tabla nueva; reservas previas quedan con
   `vencimiento_comprobante` en cero ⇒ se tratan como "ya con comprobante"
   (`ComprobanteEnviadoAt IS NULL AND vencimiento zero` no debe expirar:
   el barrido excluye `vencimiento_comprobante IS NULL/zero`). Así ninguna
   reserva histórica se anula sola tras el deploy.
2. Deploy backend + frontend juntos; endpoints viejos siguen respondiendo.
3. Rollback: revertir deploy; columnas aditivas no afectan el código anterior.

## Open Questions

_(ninguno pendiente: verificación manual, storage en BD y configuración por
entorno fueron decididos con el usuario)_
