# Tasks: cu13-notificaciones-cotizaciones

## 1. Backend — modelo y lectura

- [x] 1.1 `models/cotizacion.go`: `MensajeCotizacion.LeidoPorCliente bool` y
      `LeidoPorVendedor bool` (default false)
- [x] 1.2 Repositorio: contar no leídos por lado (cliente: remitente `vendedor`
      en sus cotizaciones; vendedor: remitente `cliente` en abiertas sin
      asignar o asignadas a él) y marcar leídos por hilo/lado

## 2. Backend — servicio y handler

- [x] 2.1 Servicio: `ContarNoLeidos(ctx, usuarioID)` y
      `MarcarLeidasCotizacion(ctx, cotizacionID, lado)` con validaciones de
      propiedad/participación
- [x] 2.2 Marcar lectura al abrir hilos: vista cliente (`obtener`) marca lado
      cliente; `ObtenerPersonal` marca lado personal solo para el asignado
- [x] 2.3 Handler notificaciones: extender `Contador` a
      `{contador, consultas, cotizaciones}` con degradación graciosa
- [x] 2.4 Tests: conteo por rol, marcado al abrir, IA no cuenta, cerradas no
      pitan al vendedor

## 3. Frontend

- [x] 3.1 `useNotificaciones`: consumir `{consultas, cotizaciones}` y exponer
      ambos contadores + aviso de incremento
- [x] 3.2 LayoutBase: puntito por enlace vía clave (`consultas` /
      `cotizaciones`) según tabla de secciones; toast cuando sube cualquiera
- [x] 3.3 Chats de cotización (cliente y personal): marcar leídos al cargar y
      en cada refresh

## 4. Verificación y documentación

- [x] 4.1 `go build ./...`, `go vet ./...`, tests backend; `npm run build` +
      tests frontend en verde
- [x] 4.2 Prueba E2E manual: vendedor responde → cliente ve puntito/toast;
      cliente responde → vendedor ve puntito en "Cotizaciones IA"; abrir el
      hilo limpia el puntito
- [ ] 4.3 Actualizar `docs/api.md`, roadmap y AGENTS.md
