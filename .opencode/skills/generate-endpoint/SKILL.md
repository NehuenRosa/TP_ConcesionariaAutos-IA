---
name: generate-endpoint
description: Crea handler + ruta + servicio para un nuevo endpoint REST en Go/Gin
compatibility: opencode
---

## Paso a paso para agregar un nuevo endpoint

1. **Modelo** — Si el recurso no existe, crear struct en `backend/internal/models/<modelo>.go` con tags GORM y JSON
2. **Repositorio** — Crear en `backend/internal/repositories/<modelo>_repository.go` encapsulando consultas GORM
3. **Servicio** — Crear en `backend/internal/services/<modelo>_service.go` con la lógica de negocio
4. **Handler** — Crear en `backend/internal/handlers/<modelo>_handler.go` con métodos HTTP (Create, Get, List, Update, Delete)
5. **Ruta** — Registrar en `backend/internal/routes/` (`Register<Modelo>Routes`) con middleware de auth según el rol
6. **Seed** (opcional) — Agregar datos de prueba en `backend/internal/seed/seed.go`

## Convenciones a respetar
- Handlers reciben servicios por inyección de dependencias
- Respuestas JSON con `gin.H{"data": ...}` o `gin.H{"error": "..."}`
- Status: 201 POST, 204 DELETE, 400 validación, 401 no auth, 403 forbiden, 404 no encontrado
- Errores genéricos al cliente, sin exponer detalles internos
