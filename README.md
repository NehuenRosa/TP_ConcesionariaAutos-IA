# Concesionaria de Autos - Sistema de Gestión

Sistema web para administrar el stock de una concesionaria de autos, con catálogo
público, consultas y reservas de clientes, turnos de test drive, panel de
administración y un asistente conversacional (chatbot).

## Stack

| Capa | Tecnología |
|------|------------|
| Backend | Go + Gin + GORM + JWT |
| Frontend | React + Vite + TypeScript + React Router + TailwindCSS |
| Base de datos | PostgreSQL |
| Chatbot | LangChain (langchaingo) |
| Infra | Docker + Docker Compose |

## Estructura del repositorio

```
├── backend/     # API REST en Go (Gin + GORM, capas handler/service/repository)
├── frontend/    # Aplicación web en React (Vite + TypeScript + TailwindCSS)
├── openspec/    # Especificaciones y changes planificados (OpenSpec)
├── docs/        # Documentación (roadmap, decisiones)
├── docker-compose.yml
├── .env.example # Variables de entorno de ejemplo
└── AGENTS.md    # Guía de contexto para agentes de IA
```

## Cómo levantar el proyecto

### Con Docker (recomendado)

Levanta PostgreSQL, backend y frontend:

```powershell
# 1. Preparar variables de entorno
Copy-Item .env.example .env

# 2. Levantar los servicios
docker compose up --build
```

- Frontend: `http://localhost:5173`
- API: `http://localhost:8080/api`
- Health check: `http://localhost:8080/api/health`

Detener todo:

```powershell
docker compose down
# Para borrar también los datos de la base:
docker compose down -v
```

### Sin Docker

#### Backend (Go)

```powershell
cd backend
go mod tidy
go run ./cmd/api
```

Requisitos: Go 1.25+ y PostgreSQL corriendo localmente (ver variables `BD_*`).

#### Frontend (React)

```powershell
cd frontend
npm install
npm run dev
```

El frontend se sirve en `http://localhost:5173`. Apunta a la API mediante la
variable `VITE_API_URL` (por defecto `http://localhost:8080/api`).

## Variables de entorno

Ver `.env.example`. Las principales:

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `PUERTO_API` | Puerto del backend | `8080` |
| `HOST_API` | Host donde escucha la API | `0.0.0.0` |
| `BD_HOST` | Host de PostgreSQL | `localhost` |
| `BD_PUERTO` | Puerto de PostgreSQL | `5432` |
| `BD_USUARIO` | Usuario de la base | `concesionaria` |
| `BD_PASSWORD` | Contraseña de la base | `concesionaria` |
| `BD_NOMBRE` | Nombre de la base | `concesionaria` |
| `BD_SSL` | Modo SSL de la conexión | `disable` |
| `JWT_SECRETO` | Secreto para firmar tokens JWT | `cambiar-en-produccion` |
| `VITE_API_URL` | URL base de la API (frontend) | `http://localhost:8080/api` |

> En producción, cambiar siempre `JWT_SECRETO`.

## Comandos útiles

### Backend (Go)

```powershell
cd backend
go build ./...
go vet ./...
go run ./cmd/api
```

### Frontend (React)

```powershell
cd frontend
npm install
npm run dev
npm run build
```

### Docker

```powershell
docker compose up --build
docker compose down
```

### OpenSpec

```powershell
openspec status --change <nombre> --json
openspec validate --strict
openspec new change <nombre-en-kebab-case>
```

> Nota para Windows: si PowerShell bloquea `npm.ps1`/`openspec.ps1` por la
> política de ejecución, usar `npm.cmd` y `openspec.cmd`.

## Verificación de salud

- Backend: `GET /api/health` responde `200 OK` con `{"estado":"ok"}`.
- Frontend: abrir `http://localhost:5173`.

## Documentación

- `AGENTS.md`: contexto completo para agentes de IA (sistema, stack, convenciones, reglas de negocio).
- `docs/roadmap.md`: backlog de casos de uso y su estado.
- `openspec/`: especificaciones y changes planificados con OpenSpec.
