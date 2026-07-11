---
name: generate-migration
description: Crea un modelo GORM y lo registra en AutoMigrate
compatibility: opencode
---

## Paso a paso para agregar un nuevo modelo

1. **Modelo** — Crear struct en `backend/internal/models/<modelo>.go`:
   - `gorm:"primaryKey"` para ID
   - `gorm:"size:N;not null"` para strings
   - `gorm:"uniqueIndex"` para campos únicos
   - `json:"nombre"` para serialización
   - `json:"-"` para ocultar campos sensibles (password)

2. **Migración** — Agregar el modelo en `backend/internal/routes/router.go` dentro de `db.AutoMigrate(...)`

3. **Seed** (opcional) — Agregar datos de prueba en `backend/internal/seed/seed.go`

4. **Verificación** — El modelo se crea automáticamente al reiniciar el backend gracias a GORM AutoMigrate
