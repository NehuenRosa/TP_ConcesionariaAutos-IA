---
name: chatbot-tasacion
description: Mantener, diagnosticar o extender el chatbot de la concesionaria (CU-10) y la tasación por fotos con valores reales de la Guía de la CCA. Usar cuando se pida tocar "el chatbot", "la tasación", "tasacion por fotos", "precios de la guia", "ArgAutos", "modelo de vision" o los archivos chatbot.go / precios.go / Chatbot.tsx. Incluye el flujo de valores reales (identificar → ArgAutos → componer), la regla de oro de no inventar precios y cómo probar E2E en local.
license: MIT
metadata:
  author: proyecto-concesionaria
  version: "1.0"
---

# Chatbot asistente y tasación con valores reales (CU-10)

Documentación operativa del chatbot (CU-10, **Resuelto**). La fuente de verdad de
requisitos es `openspec/specs/chatbot-asistente/spec.md` y el contexto general
del proyecto está en `AGENTS.md`. Respetar el stack (Go + langchaingo + Ollama,
React + Vite + TS) y el idioma (todo en español).

## Arquitectura

| Capa | Archivo | Responsabilidad |
|------|---------|-----------------|
| Service | `backend/internal/services/chatbot.go` | Lógica del asistente: `Responder`, `Tasacion`, contexto de stock, generación con Ollama. |
| Service | `backend/internal/services/precios.go` | `ServicioPrecios`: valores reales de la Guía de la CCA vía API de ArgAutos, con caché. |
| Handler | `backend/internal/handlers/chatbot.go` | Parseo multipart/JSON y respuestas HTTP de los endpoints. |
| Config | `backend/internal/config/config.go` | `OLLAMA_URL`, `MODELO_CHATBOT`, `MODELO_VISION`, `ARGAUTOS_URL`. |
| Router | `backend/internal/router/router.go` | Registra `POST /api/chatbot/mensajes` y `POST /api/chatbot/tasacion` (públicos) e inyecta `ServicioPrecios`. |
| Frontend | `frontend/src/components/Chatbot.tsx` | Widget flotante (chat + fotos). |
| Frontend | `frontend/src/services/api.ts` | Cliente HTTP con timeout (`TIEMPO_MAXIMO_MILISEGUNDOS = 140000`) vía `AbortController`. |

## Endpoints (públicos, sin autenticación)

- `POST /api/chatbot/mensajes` — body `{"mensaje": "...", "historial": [{"rol": "usuario|asistente", "contenido": "..."}]}`. Máx. 1000 caracteres por mensaje y 20 turnos de historial. Respuesta `{"respuesta": "..."}`.
- `POST /api/chatbot/tasacion` — multipart/form-data: hasta 5 fotos (JPG/JPEG/PNG/WebP, máx. 5 MB cada una) + `descripcion` opcional. Requiere al menos una foto o una descripción (`400` si viene vacío). Respuesta `{"respuesta": "..."}`.

## Flujo de tasación (NO cambiar)

Regla de oro: el LLM **nunca** genera montos. El flujo es:

1. **Identificar** (`identificarVehiculo`): se envía `promptIdentificacion` +
   descripción + fotos al modelo de visión; el LLM devuelve un JSON estricto
   `{marca, modelo, anio, estado, kilometraje}` que se extrae con
   `parsearIdentificacion`. Si no hay marca/modelo → `respuestaNoIdentificado`
   (mensaje honesto, 200).
2. **Consultar valor real** (`ServicioPrecios.Buscar`): `precios.go` llama a
   ArgAutos (`ARGAUTOS_URL`, default `https://argautos.com/api/v1`):
   - `GET /search?q=<marca> <modelo>` → versión + precio USD (`price`, `price_year`), eligiendo la versión cuyo año más se acerque al buscado.
   - `GET /versions/{id}/valuations?currency=ars` → precio ARS del año más cercano.
   - Caché en memoria 24 h, clave `marca|modelo|anio`. Errores tipados: `ErrPrecioNoEncontrado` y `ErrPrecioNoDisponible`.
3. **Componer** (`componerTasacionConReferencia` / `componerTasacionSinReferencia`):
   la respuesta final se arma en código con los valores formateados (US$ y ARS
   con separadores de miles), la versión, el año de referencia y la fuente
   ("Guía Oficial de Precios de la CCA"). Sin referencia → mensaje honesto.

**Regla de oro**: si alguien propone "pedir al LLM que estime el precio" o
"hacer que el modelo devuelva el monto", rechazarlo: los valores se componen en
código con la referencia oficial. Un cambio de prompt NO debe alterar esto.

## Fallback y degradación

- Ollama caído → chat y tasación devuelven `200` con mensaje en español que
  orienta (catálogo / concesionaria). El error interno se loguea con `slog`.
- Vehículo no identificado o sin valor oficial → respuesta honesta, sin inventar.

## Límites y constantes clave

`LargoMaximoMensaje = 1000`, `MaximoTurnosHistorial = 20`,
`MaximoImagenesTasacion = 5`, `MaximoPesoImagenBytes = 5 MiB`,
`TimeoutChatbot = TimeoutVision = 120 s`, `NumCtxChatbot = 4096`,
`NumCtxVision = 2048`. En el frontend: `TIEMPO_MAXIMO_MILISEGUNDOS = 140000` y
`EXTENSIONES_PERMITIDAS = ['jpg', 'jpeg', 'png', 'webp']` (se acepta por
extensión aunque el MIME venga vacío).

## Cómo probar en local (Windows/PowerShell)

Backend con envs explícitos (importante: `.env` puede declarar un
`MODELO_VISION` que no esté descargado; setear los modelos reales):

```powershell
# detener y levantar el backend en segundo plano
$pidApi = (Get-NetTCPConnection -LocalPort 8080).OwningProcess; Stop-Process -Id $pidApi -Force
$env:OLLAMA_URL='http://localhost:11434'; $env:MODELO_CHATBOT='llama3'; $env:MODELO_VISION='minicpm-v'
Start-Process -FilePath "$env:TEMP\concesionaria-api.exe" -WorkingDirectory (Get-Location) -RedirectStandardOutput "$env:TEMP\concesionaria-backend.log" -RedirectStandardError "$env:TEMP\concesionaria-backend-err.log" -WindowStyle Hidden
```

Verificar modelos disponibles en Ollama: `Invoke-RestMethod http://localhost:11434/api/tags`.

E2E de tasación con valores reales:

```powershell
curl.exe -s -X POST http://localhost:8080/api/chatbot/tasacion -F "fotos=@$env:TEMP\auto-prueba.png;type=image/png" -F "descripcion=Fiat Cronos 2020, bueno, 60000 km"
# → respuesta en UTF-8 con US$ y ARS reales (fuente CCA)
```

El widget se prueba en `http://localhost:5173` (frontend Vite).

## Gotchas

- **"US$ US$" duplicado** (bug ya corregido): al componer con el precio USD se
  formateaba la moneda dos veces; usar `formatearUSD`/`formatearARS` que ya
  incluyen el símbolo y NO anteponer el símbolo manualmente.
- **Modelo de visión**: `MODELO_VISION` del `.env` (ej. `llama3.2-vision`) puede
  no estar descargado; verificar con `/api/tags` y setear el modelo real
  (`minicpm-v`) al correr. `ollama` CLI no está en PATH en Windows.
- **Fotos sin MIME**: algunos navegadores reportan `file.type === ''`; el
  frontend acepta por extensión para no bloquear la tasación.
- **Colgado en "pensando"**: si una petición tarda demasiado, el cliente aborta
  a los 140 s (`AbortError` → mensaje en español); la visión puede tardar ~1 min.
- **Tasación lenta**: reconstruir el backend con `go build` + reinicio encadenado
  toma tiempo en Windows (recompilación + Windows Defender); esperar o setear
  timeouts amplios.
- **npm/openspec bloqueados**: la política de ejecución bloquea los `.ps1`; usar
  `npm.cmd` y `openspec.cmd`.

## Skills de OpenSpec

Cambios sobre el chatbot se planifican con OpenSpec (un change por cambio de
comportamiento). Ver `openspec/specs/chatbot-asistente/spec.md` (spec archivada,
fuente de verdad) y los skills `openspec-propose` / `openspec-apply-change` /
`openspec-archive-change`.
