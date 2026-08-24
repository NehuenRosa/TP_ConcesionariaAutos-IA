# Tasks: cu12-bandeja-cotizaciones

## 1. Backend — modelo y repositorio

- [x] 1.1 `models/cotizacion.go`: `VendedorID *uint`, `Vendedor *Usuario` (foreignKey), `FechaToma *time.Time`; constante `RemitenteVendedor = "vendedor"`
- [x] 1.2 Repositorio: `ListarBandeja(ctx)` con preload de Vehiculo/Cliente/Vendedor/último mensaje; `Actualizar` persiste estado, vendedor y fecha de toma

## 2. Backend — servicio

- [x] 2.1 `CotizacionService`: métodos `ListarBandeja`, `ObtenerPersonal`, `Tomar`, `ResponderComoVendedor`, `CerrarPersonal`; errores `ErrCotizacionYaAtendida` y `ErrCotizacionNoTomada`
- [x] 2.2 Preview: descifrar solo el último mensaje por cotización en la bandeja; hilo completo descifrado en la atención
- [x] 2.3 Apagado de IA en `EnviarMensaje`: si `VendedorID != nil`, guardar solo el mensaje del cliente sin llamar al generador
- [x] 2.4 Tests con fakes: tomar asigna/idempotente/otro vendedor, respuesta vendedor cifra sin IA, mensaje cliente en hilo atendido no llama al generador, cierre del personal

## 3. Backend — handler y rutas

- [x] 3.1 DTOs JSON en español (`cliente`, `vendedor`, `fechaToma`) y handlers nuevos
- [x] 3.2 Rutas bajo `/api/cotizaciones` con `ExigirRol("vendedor")`: `GET /bandeja`, `GET /:id/personal`, `PUT /:id/tomar`, `POST /:id/mensajes-vendedor`, `PUT /:id/cerrar-personal`

## 4. Frontend — vendedor

- [x] 4.1 `types/cotizacion.ts`: tipos de bandeja/atención + funciones en `api.ts`
- [x] 4.2 Página `/vendedor/cotizaciones` (bandeja con filtro) y `/vendedor/cotizaciones/:id` (tomar/responder/cerrar) + link "Cotizaciones IA" en menú vendedor

## 5. Frontend — cliente

- [x] 5.1 `ChatCotizacion.tsx`: estilo diferenciado para remitente `"vendedor"`, banner de atención humana, polling cada 10 s

## 6. Verificación y documentación

- [x] 6.1 `go build ./...`, `go vet ./...`, tests backend en verde; `npm run build` y tests frontend en verde
- [x] 6.2 Prueba E2E manual: pedir cotización por chat → tomarla como vendedor → responder → ver la respuesta del vendedor desde el cliente y que la IA no conteste más (verificada por API contra el stack Docker: bandeja, tomar, silencio de IA, mensaje `vendedor`, cierre)
- [x] 6.3 Actualizar `docs/api.md`, `docs/frontend.md`, roadmap y AGENTS.md
