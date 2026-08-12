# Concesionaria de Autos - Sistema de Gestión

Sistema web para administrar el stock de una concesionaria de autos, con catálogo
público, consultas y reservas de clientes, turnos de test drive, panel de
administración y un asistente conversacional (chatbot) que responde sobre el
stock real y tasa el auto del usuario por fotos con valores de la Guía de la CCA.

## Stack

| Capa | Tecnología |
|------|------------|
| Backend | Go + Gin + GORM + JWT |
| Frontend | React + Vite + TypeScript + React Router + TailwindCSS |
| Base de datos | PostgreSQL |
| Chatbot | LangChain (langchaingo) + Ollama (chat y visión); tasación con valores reales de la Guía de la CCA vía API de ArgAutos |
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

Levanta PostgreSQL, backend, frontend y Ollama (chatbot):

```powershell
# 1. Preparar variables de entorno
Copy-Item .env.example .env

# 2. Levantar los servicios
docker compose up --build
```

Para el chatbot, descargar los modelos dentro del contenedor de Ollama (una vez):

```powershell
docker compose exec ollama ollama pull llama3
docker compose exec ollama ollama pull minicpm-v
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
| `BD_URL` | URL completa de Postgres (p. ej. Render; pisa `BD_*`) | *(vacío)* |
| `BD_HOST` | Host de PostgreSQL | `localhost` |
| `BD_PUERTO` | Puerto de PostgreSQL | `5432` |
| `BD_USUARIO` | Usuario de la base | `concesionaria` |
| `BD_PASSWORD` | Contraseña de la base | `concesionaria` |
| `BD_NOMBRE` | Nombre de la base | `concesionaria` |
| `BD_SSL` | Modo SSL de la conexión | `disable` |
| `JWT_SECRETO` | Secreto para firmar tokens JWT | `cambiar-en-produccion` |
| `VITE_API_URL` | URL base de la API (frontend) | `http://localhost:8080/api` |
| `OLLAMA_URL` | URL de Ollama (en Docker: `http://ollama:11434`) | `http://localhost:11434` |
| `MODELO_CHATBOT` | Modelo de chat / comparación | `llama3` |
| `MODELO_VISION` | Modelo de visión para la tasación por fotos | `minicpm-v` |
| `ARGAUTOS_URL` | API de valores de la Guía de la CCA (tasación) | `https://argautos.com/api/v1` |

> En producción, cambiar siempre `JWT_SECRETO`. Los modelos de Ollama se
> descargan con `ollama pull <modelo>`.

## Perfiles de prueba

Al arrancar, el backend crea automáticamente tres cuentas por defecto (sembradas)
para poder probar los distintos roles. Son credenciales de **desarrollo**; las
reales se definirán más adelante.

| Rol | Email | Contraseña |
|-----|-------|------------|
| Administrador | `administrador@concesionaria.local` | `Admin123!` |
| Vendedor | `vendedor@concesionaria.local` | `Vendedor123!` |
| Cliente | `cliente@concesionaria.local` | `Cliente123!` |

> El rol `cliente` también se puede crear desde el registro público en
> `http://localhost:5173/registro` con cualquier email y una contraseña de al
> menos 8 caracteres.

## Chatbot (CU-10)

Asistente integrado en las páginas públicas (widget flotante) con dos endpoints:

- `POST /api/chatbot/mensajes` — responde sobre el stock disponible y orienta a
  consultar o pedir un test drive.
- `POST /api/chatbot/tasacion` — tasación por fotos (hasta 5, JPG/PNG/WebP) +
  `descripcion` opcional. El modelo de visión **identifica** el vehículo y el
  valor se compone en código con la **Guía Oficial de Precios de la CCA** (vía
  API de ArgAutos): nunca se inventan precios. Si no se identifica el vehículo o
  no hay valor de referencia, la respuesta es honesta y orienta al usuario.

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
