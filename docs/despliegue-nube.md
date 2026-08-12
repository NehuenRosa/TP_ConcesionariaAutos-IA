# Despliegue en la nube (Render + RunPod)

Este documento explica cómo correr el sistema completo en la nube para no
consumir recursos de la PC local. El stack en nube queda así:

| Componente | Proveedor | Costo |
|------------|-----------|-------|
| Frontend (React) | Render (web service con Docker) | Plan free (duerme con inactividad) |
| Backend (Go + Gin) | Render (web service con Docker) | Plan free (duerme con inactividad) |
| PostgreSQL | Render (base de datos) | Plan free (expira a los 90 días) |
| Chatbot / visión (Ollama) | RunPod (pod con GPU) | Por segundo (~US$0,2/h, se apaga cuando no se usa) |

La configuración ya está en `render.yaml` (blueprint de Render).

## 1. Levantar Ollama en RunPod (GPU)

1. Crear cuenta en https://runpod.io y cargar saldo (US$5 alcanzan de sobra).
2. **Deploy** → buscar la plantilla **Ollama** → elegir una GPU con ≥ 16 GB de
   VRAM (p. ej. RTX 3090) → **Deploy On-Demand**.
3. Cuando el pod esté en **Running**, abrir la terminal web y descargar los
   modelos:
   ```
   ollama pull llama3
   ollama pull minicpm-v
   ```
4. En la pestaña del pod, exponer el puerto **11434** como puerto HTTP. RunPod
   muestra una URL pública con el formato:
   ```
   https://<TU-POD-ID>-11434.proxy.runpod.net
   ```
5. Apagar el pod cuando no se use (RunPod factura por segundo).

## 2. Desplegar en Render

1. Crear cuenta en https://render.com y conectarla con el repositorio de GitHub.
2. **New + → Blueprint** → seleccionar el repo → Render detecta `render.yaml`.
3. Reemplazar los placeholders en `render.yaml` antes de crear el blueprint:
   - `TU-POD-11434.proxy.runpod.net` → la URL del pod de RunPod (paso 1.4).
   - `TU-FRONTEND.onrender.com` → la URL que Render asigne al frontend.
   - `TU-BACKEND.onrender.com` → la URL que Render asigne al backend.
4. **Apply**. Render crea: base Postgres, backend y frontend.

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
