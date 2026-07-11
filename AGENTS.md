# AGENTS.md — Contexto para opencode

## Arquitectura general

- **Backend**: Go con framework Gin, GORM como ORM, autenticación JWT, API REST.
- **Frontend**: React con Vite + TypeScript, React Router, TailwindCSS.
- **Base de datos**: PostgreSQL.
- **Chatbot**: LangChain (langchaingo) para asistente conversacional integrado.
- **Infraestructura**: Docker Compose para desarrollo local.

## Convenciones de código

### Backend (Go)
- Usar `internal/` para código privado del paquete principal.
- Modelos en `internal/models/` con structs que incluyen tags GORM y JSON.
- Handlers en `internal/handlers/` reciben servicios por inyección de dependencias.
- Servicios en `internal/services/` contienen la lógica de negocio.
- Repositorios en `internal/repositories/` encapsulan consultas GORM.
- Middleware en `internal/middleware/`.
- Configuración vía variables de entorno en `internal/config/`.
- Rutas centralizadas en `internal/routes/router.go`.
- Seed data en `internal/seed/seed.go`.

### Frontend (React)
- Páginas en `src/pages/` organizadas por rol (admin/, seller/).
- Componentes compartidos en `src/components/`.
- Servicios HTTP en `src/services/`.
- Tipos TypeScript en `src/types/index.ts`.
- Contexto de autenticación en `src/context/AuthContext.tsx`.
- Hooks personalizados en `src/hooks/`.
- Estilos con TailwindCSS (sin archivos CSS adicionales salvo `index.css`).

## Reglas de negocio

### Roles
1. **Cliente**: puede ver catálogo, enviar consultas, solicitar test drives, reservar vehículos.
2. **Vendedor**: puede gestionar consultas (bandeja), turnos de test drive y reservas.
3. **Administrador**: ABM de vehículos, gestión de usuarios, dashboard con métricas.

### Vehículos
- Estados: `disponible`, `reservado`, `vendido`.
- Solo vehículos en estado `disponible` se muestran en el catálogo público.
- Imágenes: URLs almacenadas como arreglo de strings.
- Al reservar un vehículo, pasa a estado `reservado`.

### Consultas
- Estados: `pendiente`, `en_conversacion`, `cerrada`.
- Cada consulta se vincula a un vehículo y un cliente.
- El vendedor puede cambiar el estado y agregar respuestas.

### Test Drives
- Se validan contra superposición de horarios.
- Estados: `pendiente`, `confirmado`, `cancelado`, `completado`.

### Autenticación
- JWT con claims: user_id, role, email.
- Middleware verifica token y lo pasa al contexto de Gin.
- Contraseñas hasheadas con bcrypt.

## Custom agents (`.opencode/agents/`)

| `@` | Función |
|-----|---------|
| `@api-explorer` | Documenta endpoints con ejemplos curl |
| `@test-writer` | Genera tests Go/React siguiendo patrones existentes |
| `@migrate-db` | Planea migraciones de esquema con GORM |
| `@security` | Audita seguridad del proyecto |

## Skills de opencode (`.opencode/skills/`)

| Skill | Descripción |
|-------|-------------|
| `generate-endpoint` | Crea handler + ruta + servicio para un nuevo endpoint REST |
| `generate-component` | Crea componente React con TypeScript y TailwindCSS |
| `generate-migration` | Crea modelo GORM y migración automática |
| `run-verify` | Ejecuta compilación Go + TypeScript check + build frontend |

## MCP servers disponibles (`opencode.json`)

| Servidor | Uso |
|----------|-----|
| `sequential-thinking` | Razonamiento paso a paso para problemas complejos |
| `github` (remote) | Buscar código en GitHub vía grep.app |

## Metodología (SDD)

Cada funcionalidad sigue:
1. **Spec**: especificación escrita del caso de uso.
2. **Plan**: desglose en tareas técnicas.
3. **Implementación**: ejecución por el agente.
4. **Validación**: verificación contra la spec.

## Setup del entorno

### Base de datos (PostgreSQL)

```bash
# Levantar todos los servicios (db, backend, frontend)
docker compose up -d

# Verificar que la base esté lista
docker compose ps db

# Si la base "concesionaria" no se creó automáticamente:
docker exec tp_concesionariaautos-ia-db-1 psql -U postgres -c "CREATE DATABASE concesionaria;"

# Conectarse desde un cliente externo (DBeaver, pgAdmin, etc.):
#   Host:     localhost
#   Puerto:   5433
#   Usuario:  postgres
#   Password: postgres
#   Base:     concesionaria
```

### Credenciales de desarrollo

| Servicio  | URL                           |
|-----------|-------------------------------|
| Frontend  | http://localhost:5173         |
| Backend   | http://localhost:8080         |
| PostgreSQL| localhost:5433 (postgres/postgres) |

### Seed data (usuarios de prueba)

| Rol          | Email                      | Contraseña   |
|--------------|----------------------------|--------------|
| administrador| admin@concesionaria.com    | admin123     |
| vendedor     | vendedor@concesionaria.com | vendedor123  |
| cliente      | cliente@test.com           | cliente123   |

## Servidores MCP

- **GitHub MCP**: gestión de branches, issues y PRs.
- **Context7 MCP**: consulta de documentación de librerías.
- **PostgreSQL MCP**: inspección de esquema y debugging de base de datos.
