---
name: run-verify
description: Ejecuta verificación completa del backend Go y frontend React
compatibility: opencode
---

## Verificación completa post-cambio

Siempre que hagas cambios, ejecuta en orden:

### Backend (Go)
1. **Compilar** — `go build ./cmd/server` desde `backend/`
2. **Tests** — `go test ./...` (cuando existan)

### Frontend (React)
3. **TypeScript check** — `npx tsc --noEmit` desde `frontend/`
4. **Build** — `npm run build` desde `frontend/`

Si alguno falla:
- Detén la ejecución
- Explica qué archivo y línea causan el error
- Propón una corrección concreta
- No marques la tarea como completa hasta que pase todo
