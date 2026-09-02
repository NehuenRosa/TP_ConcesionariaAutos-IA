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
| Chatbot | LangChain (langchaingo) conectado a **Google AI Gemini en la nube** (`googleai`); tasación con valores reales de la Guía de la CCA vía API de ArgAutos |
| Infra | Docker + Docker Compose |

## Estructura del repositorio

```
├── backend/                 # API REST en Go
│   ├── cmd/api/             # Punto de entrada (main)
│   └── internal/
│       ├── config/          # Carga de configuración desde variables de entorno
│       ├── database/        # Conexión GORM + auto-migración + seed
│       ├── models/          # Entidades / modelos de GORM
│       ├── handlers/        # Handlers HTTP (parsean request/response)
│       ├── services/        # Lógica de negocio
│       ├── repositories/    # Acceso a datos (GORM)
│       ├── middleware/      # CORS, autenticación JWT, roles
│       ├── router/          # Registro de rutas
│       └── token/           # Utilidades de JWT
├── frontend/                # Aplicación web en React
│   └── src/
│       ├── components/      # Componentes reutilizables (Chatbot, UI)
│       ├── layouts/         # Layout base (header/footer)
│       ├── pages/           # Páginas por ruta
│       ├── routes/          # Definición de rutas (React Router)
│       ├── services/        # Cliente HTTP centralizado (api.ts)
│       ├── types/           # Tipos TypeScript compartidos
│       ├── hooks/           # Hooks personalizados (useAuth, useNotificaciones)
│       ├── utils/           # Utilidades de formato
│       └── test/            # Setup de tests (Vitest + Testing Library)
├── openspec/                # OpenSpec: specs + changes (uno por CU)
├── docs/                    # Documentación (roadmap, despliegue, API, frontend)
├── docker-compose.yml       # postgres + backend + frontend
├── .env.example             # Variables de entorno de ejemplo
├── AGENTS.md                # Guía de contexto para agentes de IA
└── render.yaml              # Blueprint de despliegue en Render (nube)
```

## Cómo levantar el proyecto

### Con Docker (recomendado)

Levanta PostgreSQL, backend y frontend. El chatbot usa **Gemini en la nube**
(configurá `GOOGLE_API_KEY`).

```powershell
# 1. Preparar variables de entorno (editar y completar BD_PASSWORD y JWT_SECRETO)
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

> **Ollama**: este proyecto ya no usa Ollama. El LLM es exclusivamente Gemini en
> la nube.

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
| `BD_PASSWORD` | Contraseña de la base (**obligatoria** en compose) | `concesionaria` |
| `BD_NOMBRE` | Nombre de la base | `concesionaria` |
| `BD_SSL` | Modo SSL de la conexión | `disable` |
| `JWT_SECRETO` | Secreto para firmar tokens JWT (**obligatorio** en compose) | `cambiar-en-produccion` |
| `CORS_ORIGENES` | Orígenes permitidos por CORS | `*` |
| `PROVEEDOR_LLM` | Proveedor del LLM (único soportado: `googleai`) | *(vacío → `googleai`)* |
| `GOOGLE_API_KEY` | API key de Gemini (gratis en Google AI Studio) | *(vacío)* |
| `MODELO_CHATBOT` | Modelo de chat / comparación (vacío = default) | `gemini-3.5-flash-lite` |
| `MODELO_VISION` | Modelo de visión para la tasación (vacío = default) | `gemini-3.5-flash-lite` |
| `ARGAUTOS_URL` | API de valores de la Guía de la CCA (tasación) | `https://argautos.com/api/v1` |
| `GOOGLE_CLIENT_ID` | Client ID OAuth de Google para "Continuar con Google" (CU-11); vacío = deshabilitado | *(vacío)* |
| `CBU_CONCESIONARIA` | CBU mostrado al cliente para transferir la seña de la reserva (CU-08); vacío = la UI avisa que el personal pasa los datos | *(vacío)* |
| `ALIAS_CONCESIONARIA` | Alias CBU para la transferencia de la seña (CU-08) | *(vacío)* |
| `VITE_API_URL` | URL base de la API (frontend) | `http://localhost:8080/api` |

> En producción, cambiar siempre `JWT_SECRETO` y no exponer `GOOGLE_API_KEY` en
> el repositorio (vive solo en `.env` local o en variables de Render). El
> backend usa por defecto el modelo `gemini-3.5-flash-lite` (1M de contexto,
> texto + visión, gratis para cuentas nuevas).

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

Asistente integrado en las páginas públicas (widget flotante) con dos endpoints
públicos (sin autenticación):

- `POST /api/chatbot/mensajes` — chat con historial; responde sobre el stock
  disponible y orienta a consultar o pedir un test drive.
- `POST /api/chatbot/tasacion` — tasación por fotos (hasta 5, JPG/PNG/WebP) +
  `descripcion` opcional. El modelo de visión solo **identifica** el vehículo
  (`{marca, modelo, anio, estado, kilometraje}`) y el valor se compone en código
  con la **Guía Oficial de Precios de la CCA** (vía API de ArgAutos, caché de
  24 h): el LLM nunca genera montos. Si no se identifica el vehículo o no hay
  valor de referencia, la respuesta es honesta y orienta al usuario.

**Proveedor**: **Gemini en la nube** (`googleai`, clave `GOOGLE_API_KEY`,
modelo `gemini-3.5-flash-lite`) — único soportado. Los errores
transitorios de Gemini (503 por alta demanda, 429 de cuota, otros 5xx) se
reintentan con espera exponencial; si el proveedor sigue caído, chat y tasación
devuelven `200` con un mensaje en español que orienta al usuario (el error
interno se loguea).

## Documentación detallada

- `docs/api.md` — referencia completa de la API REST (endpoints, métodos, roles, payloads).
- `docs/frontend.md` — rutas, páginas, componentes y estructura del frontend.
- `docs/despliegue-nube.md` — cómo desplegar en Render.
- `AGENTS.md` — contexto completo para agentes de IA (sistema, stack, convenciones, reglas de negocio).
- `docs/roadmap.md` — backlog de casos de uso y su estado.
- `openspec/` — especificaciones y changes planificados con OpenSpec.

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
