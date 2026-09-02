# Despliegue en la nube (Render)

Este documento explica cómo correr el sistema completo en la nube para no
consumir recursos de la PC local. El stack en nube queda así:

| Componente | Proveedor | Costo |
|------------|-----------|-------|
| Frontend (React) | Render (web service con Docker) | Plan free (duerme con inactividad) |
| Backend (Go + Gin) | Render (web service con Docker) | Plan free (duerme con inactividad) |
| PostgreSQL | Render (base de datos) | Plan free (expira a los 90 días) |
| Chatbot / visión (Gemini) | Google AI (nube) | Free tier sin tarjeta (~1000 requests/día) |

La configuración ya está en `render.yaml` (blueprint de Render).

## 1. Obtener la API key de Gemini

1. Crear la clave en https://aistudio.google.com → **Get API key** (requiere una
   cuenta de Google, sin tarjeta de crédito).
2. Guardarla en la variable de entorno `GOOGLE_API_KEY` del backend en Render
   (y en el `.env` local si querés la nube también en dev).

> El backend usa **Google AI Gemini en la nube** como proveedor del LLM (chat y
> visión) con la clave `GOOGLE_API_KEY`. El free tier de
> `gemini-3.5-flash-lite` (1M de contexto, texto + visión) alcanza para una
> demo sin costo. Este proyecto ya no usa Ollama.

## 2. Desplegar en Render

1. Crear cuenta en https://render.com y conectarla con el repositorio de GitHub.
2. **New + → Blueprint** → seleccionar el repo → Render detecta `render.yaml`.
3. Reemplazar los placeholders en `render.yaml` antes de crear el blueprint:
   - `TU-FRONTEND.onrender.com` → la URL que Render asigne al frontend.
   - `TU-BACKEND.onrender.com` → la URL que Render asigne al backend.
4. En el panel del backend, agregar la variable `GOOGLE_API_KEY` (paso 1).
5. **Apply**. Render crea: base Postgres, backend y frontend.

> Las URLs se conocen después del primer deploy. Un truco: dejar el placeholder,
> desplegar, y luego actualizar `CORS_ORIGENES` y `VITE_API_URL` en el panel de
> Render (o editar `render.yaml` y volver a aplicar) y redeployar.

## 3. Verificar

- Backend: `https://<TU-BACKEND>.onrender.com/api/health` → `200 OK`.
- Frontend: abrir `https://<TU-FRONTEND>.onrender.com` y probar el catálogo y el
  chatbot (chat y tasación con fotos).
- El backend crea las tablas y los datos iniciales automáticamente al arrancar
  (auto-migración + seed).

## Notas del plan free

- Los web services duermen tras **15 minutos** de inactividad; el primer
  request tarda unos segundos en "despertarlos".
- El Postgres free **expira a los 90 días** (solo región US). Para una demo es
  suficiente; para más tiempo se puede migrar a Neon (también gratis) usando la
  variable `BD_URL`.
- `JWT_SECRETO` lo genera Render automáticamente (`generateValue: true`).

## Local vs. nube

- Local (Docker Compose): sigue funcionando igual con las variables de
  `docker-compose.yml`. La nube usa `BD_URL` (URL completa) y `PORT` que inyecta
  Render; el backend soporta ambos modos (`backend/internal/config/config.go` y
  `backend/internal/database/database.go`).
- Chatbot: en ambos entornos `PROVEEDOR_LLM` vacío auto-elige Gemini en la nube
  (requiere `GOOGLE_API_KEY`).
