# AGENTS.md

Guía de contexto para agentes de IA que trabajen sobre este repositorio.
Leé este archivo completo antes de tocar código.

## Descripción del sistema

Sistema web de gestión para una concesionaria de autos. Permite administrar el
stock de vehículos, exhibirlos en un catálogo público con búsqueda y filtros,
gestionar consultas y reservas de clientes, coordinar turnos de test drive y
ofrecer un asistente conversacional (chatbot) que responde preguntas sobre los
vehículos disponibles.

## Perfiles de usuario

| Perfil | Descripción |
|--------|-------------|
| **Cliente / visitante** | Navega el catálogo, consulta, reserva y solicita test drives. |
| **Vendedor** | Gestiona consultas, reservas y turnos. |
| **Administrador** | Administra stock de vehículos, usuarios y visualiza métricas. |

## Casos de uso

| ID | Caso de uso | Descripción |
|----|-------------|-------------|
| CU-01 | Autenticación y roles | Registro e inicio de sesión. Tres roles (cliente, vendedor, administrador) con autorización basada en JWT. |
| CU-02 | Gestión de vehículos (ABM) | El administrador da de alta, modifica y da de baja vehículos, con ficha técnica (marca, modelo, año, kilometraje, combustible, transmisión, precio) e imágenes. |
| CU-03 | Catálogo público | Listado paginado de vehículos disponibles y vista de detalle con ficha técnica y galería. |
| CU-04 | Búsqueda y filtrado avanzado | Búsqueda por texto libre y filtros combinables: marca, modelo, rango de años, rango de precio, tipo, combustible y condición (nuevo/usado), con ordenamiento por precio o año. |
| CU-05 | Consulta / cotización | El cliente envía una consulta asociada a un vehículo específico; queda vinculada al vehículo y al cliente. |
| CU-06 | Gestión de consultas | El vendedor ve la bandeja de consultas, responde y actualiza estado (pendiente, en conversación, cerrada), con historial. |
| CU-07 | Turno de test drive | El cliente solicita turno eligiendo fecha y franja horaria; el sistema valida disponibilidad y evita superposición. |
| CU-08 | Reserva de vehículo | El cliente reserva una unidad (disponible → reservado); el vendedor confirma la venta o cancela y libera la unidad. |
| CU-09 | Panel de administración | Dashboard con vehículos por estado, consultas por período, reservas activas y test drives agendados, con gráficos simples. |
| CU-10 | Chatbot asistente | Chat integrado con un asistente (LangChain + Gemini en la nube, con Ollama local de respaldo) que responde en lenguaje natural sobre el stock real, orienta a consultar o pedir un test drive, y tasa el auto del usuario por fotos usando valores reales de la Guía de la CCA (no inventa precios). **Resuelto** (ver `docs/roadmap.md`). |

El backlog completo con el estado de cada CU está en `docs/roadmap.md`.

## Stack (respetar exactamente)

| Capa | Tecnología |
|------|------------|
| Backend | Go + Gin (HTTP) + GORM (ORM) + JWT. API REST. |
| Frontend | React + Vite + TypeScript + React Router + TailwindCSS. |
| Base de datos | PostgreSQL. |
| Chatbot | LangChain (langchaingo en el backend Go), conectado a la API del sistema y a un LLM provisto por **Google AI Gemini en la nube** (`PROVEEDOR_LLM=googleai`, modelos `MODELO_CHATBOT` y `MODELO_VISION`) u **Ollama local** como respaldo (modelos `MODELO_CHATBOT`/`MODELO_VISION`). La tasación usa valores reales de la Guía Oficial de Precios de la CCA vía API de ArgAutos (`ARGAUTOS_URL`); nunca inventa montos. |
| Infra | Docker + Docker Compose para desarrollo local. Git. |

**Regla:** no agregar dependencias fuera de este stack sin autorización.

## Estructura del repositorio

```
├── backend/                # API REST en Go
│   ├── cmd/api/            # Punto de entrada (main)
│   ├── internal/
│   │   ├── config/         # Carga de configuración desde variables de entorno
│   │   ├── database/       # Conexión GORM + migraciones (auto-migración)
│   │   ├── models/         # Entidades / modelos de GORM
│   │   ├── handlers/       # Handlers HTTP (capa de presentación)
│   │   ├── services/       # Lógica de negocio (capa de aplicación)
│   │   ├── repositories/   # Acceso a datos (abstracción sobre GORM)
│   │   ├── middleware/     # Middlewares (CORS, autenticación JWT)
│   │   └── router/         # Registro de rutas
│   ├── go.mod
│   └── Dockerfile
├── frontend/               # Aplicación web en React
│   └── src/
│       ├── components/     # Componentes reutilizables
│       ├── layouts/        # Layouts (header/footer)
│       ├── pages/          # Páginas por ruta
│       ├── routes/         # Definición de rutas (React Router)
│       ├── services/       # Cliente HTTP centralizado (api.ts)
│       ├── types/          # Tipos TypeScript compartidos
│       └── hooks/          # Hooks personalizados
├── openspec/               # OpenSpec: specs + changes
│   ├── config.yaml         # Contexto del proyecto (se inyecta en toda planificación)
│   ├── specs/              # Especificaciones (fuente de verdad)
│   └── changes/            # Changes (propuestas de cambio, una por CU)
├── docs/                   # Documentación (roadmap, decisiones)
├── docker-compose.yml      # postgres + backend + frontend
├── .env.example            # Variables de entorno de ejemplo
└── AGENTS.md               # Este archivo
```

### Flujo de datos backend (capas)

`handler (HTTP) → service (negocio) → repository (GORM) → PostgreSQL`

- Los **handlers** solo parsean request/response y delegan en services. No
  contienen lógica de negocio.
- Los **services** contienen la lógica de negocio y reglas (estados, validaciones).
- Los **repositories** encapsulan el acceso a la base con GORM.
- Los **models** son las entidades de GORM (tablas).

## Convenciones de código

### Idioma

- **Todo en español**: código, variables, funciones, tipos, entidades, nombres de
  archivos, comentarios, commits, mensajes de error y textos de UI.
- Solo las palabras reservadas y bibliotecas del lenguaje quedan en su idioma
  original (ej. `func`, `if`, `package`, `import` en Go; `const`, `interface` en TS).

### Naming

- **Go**: `camelCase` para variables y funciones, `PascalCase` para tipos exportados,
  acrónimos en mayúsculas (`API`, `JWT`, `URL`).
- **TypeScript/React**: `camelCase` para variables, funciones y props; `PascalCase`
  para componentes y tipos. Archivos de componentes en `PascalCase.tsx`; el resto
  en `camelCase.ts`.
- Nombres de archivos/carpetas en español (ej. `vehiculos.go`, `api.ts`).
- Nombres descriptivos; evitar abreviaturas crípticas.

### Manejo de errores

- **Go**: los servicios retornan `error` y los handlers lo traducen a respuestas
  HTTP con código adecuado (`400`, `404`, `500`, etc.) y cuerpo JSON
  `{"error": "mensaje en español"}`. Loggear con `log`/`slog` del stdlib.
- **Frontend**: el cliente HTTP centralizado (`services/api.ts`) normaliza errores
  y expone mensajes en español al usuario.

### Capas

- Prohibido saltarse capas: handlers no tocan GORM directo; repositories no hacen
  validación de negocio; services no escriben respuestas HTTP.

## Reglas de negocio conocidas

### Estados del vehículo

- `disponible`, `reservado`, `vendido` y `dado_de_baja` (o equivalente).
- Solo los vehículos **disponibles** se muestran en el catálogo público (CU-03).
- Un vehículo reservado deja de estar disponible; la reserva confirmada como venta
  lo pasa a `vendido` (CU-08).

### Turnos de test drive

- El sistema valida disponibilidad y **evita superposición**: no puede existir más
  de un turno para la misma unidad en la misma fecha y franja horaria (CU-07).

### Roles y permisos

- **Cliente/visitante**: catálogo, consultas, reservas, test drives.
- **Vendedor**: gestiona consultas (CU-06), reservas y turnos.
- **Administrador**: ABM de vehículos (CU-02), gestión de usuarios y métricas (CU-09).
- Autorización basada en JWT con claims de rol.

### Chatbot asistente (CU-10)

- **Endpoints públicos**: `POST /api/chatbot/mensajes` (chat con historial) y
  `POST /api/chatbot/tasacion` (multipart con hasta 5 fotos JPG/PNG/WebP de máx.
  5 MB + `descripcion` opcional). No requieren autenticación.
- **Proveedores**: por defecto **Google AI Gemini en la nube** (`googleai`,
  clave `GOOGLE_API_KEY`, modelo `gemini-flash-lite-latest`: contexto 1M de
  tokens, texto + visión, free tier sin tarjeta). Si no hay clave o se setea
  `PROVEEDOR_LLM=ollama`, usa **Ollama local** (modelos `MODELO_CHATBOT`/
  `MODELO_VISION`). `NuevoChatbotService` en `chatbot.go` despacha según el
  proveedor; `PROVEEDOR_LLM` vacío auto-elige googleai si hay `GOOGLE_API_KEY`.
- **Tasación con valores reales**: el modelo de visión solo *identifica* el
  vehículo (JSON `{marca, modelo, anio, estado, kilometraje}`); el monto lo
  compone el código (`precios.go`) con el valor oficial de la **Guía de la CCA**
  vía ArgAutos (`ServicioPrecios`, caché 24 h, clave `marca|modelo|anio`). El
  LLM **nunca** genera montos. Si no se identifica o no hay referencia, la
  respuesta es honesta (orienta a la concesionaria) — no inventa valores.
- **Regla de oro**: no cambiar el flujo de la tasación por un prompt que pida al
  LLM "estimar" precios; siempre componer en código con la referencia oficial.
- **Fallback**: si el proveedor LLM está caído (Gemini u Ollama), chat y
  tasación devuelven `200` con un mensaje en español que orienta al usuario (el
  error interno se loguea con `slog`).

## Flujo de trabajo con OpenSpec

OpenSpec se usa para planificar cada caso de uso antes de implementarlo. Schema:
`spec-driven` (artefactos: proposal → specs → design → tasks).

### Comandos principales (PowerShell)

```powershell
# Inicializar / actualizar instrucciones
openspec init
openspec update

# Crear un change nuevo (una por CU)
openspec new change <nombre-en-kebab-case>

# Estado de un change (artefactos pendientes)
openspec status --change <nombre> --json

# Instrucciones para generar un artefacto
openspec instructions <artefacto> --change <nombre> --json

# Validar (obligatorio antes de considerar listo un change)
openspec validate --strict

# Archivar un change completado (funde los specs en los principales)
openspec archive <nombre>
```

### Ciclo recomendado

1. **Planificar**: crear el change y completar los 4 artefactos.
2. **Validar**: `openspec validate --strict` sin errores.
3. **Implementar**: tareas del `tasks.md`, verificando con build/tests.
4. **Archivar**: `openspec archive <nombre>` al finalizar.

> Los commands `/opsx:*` (propose, apply, sync, archive) también están disponibles
> desde el asistente y siguen el mismo flujo.

## Comandos habituales

### Backend (Go)

```powershell
cd backend
go mod tidy
go build ./...
go run ./cmd/api
go vet ./...
```

### Frontend (React)

```powershell
cd frontend
npm install
npm run dev
npm run build
```

### Docker (todo el sistema)

```powershell
docker compose up --build
docker compose down
```

### Verificación de salud

- Backend: `GET /api/health` → `200 OK`.
- Frontend: abrir `http://localhost:<puerto>`.

## Variables de entorno

Ver `.env.example`. Para levantar local, copiar a `.env` y ajustar.
El backend requiere: conexión a PostgreSQL (`BD_*`), puerto (`PUERTO_API`) y
secreto JWT (`JWT_SECRETO`). El frontend usa `VITE_API_URL` para apuntar a la API.
Para el chatbot: `PROVEEDOR_LLM`, `GOOGLE_API_KEY` (Gemini en la nube, gratis en
Google AI Studio), `OLLAMA_URL` (respaldo local), `MODELO_CHATBOT`,
`MODELO_VISION` (modelos a descargar con `ollama pull` si usás local) y
`ARGAUTOS_URL` (fuente de precios de la CCA).

`docker compose` aplica fail-fast: **falla si `BD_PASSWORD` o `JWT_SECRETO` no
están definidas** en `.env` (no hay defaults inseguros). Los puertos publicados
quedan en loopback (127.0.0.1) y los contenedores corren no-root con
`cap_drop: [ALL]` y `read_only` donde la imagen lo permite (ver
`docker-compose.yml`). El backend usa el **Ollama nativo del host**
(`OLLAMA_URL=http://host.docker.internal:11434`) por defecto: no duplica
modelos ni procesos. Si se prefiere Ollama en contenedor (con `gpus: all`,
requiere NVIDIA Container Toolkit; sin toolkit quitar ese bloque para CPU),
levantarlo con `docker compose --profile ollama-contenedor up -d`.

El MCP de GitHub (configurado en `.opencode/opencode.json`) autentica con un
PAT vía `{env:GITHUB_TOKEN}`. Como opencode resuelve las variables de entorno
desde el proceso al iniciar, `GITHUB_TOKEN` debe estar seteado en la sesión
(no basta con `.env`); se puede obtener con `gh auth token`.

## Notas de la plataforma

- Trabajamos en **Windows (PowerShell 5.1)**. Los comandos de este archivo son
  compatibles con PowerShell.
- El script `npm.ps1`/`openspec.ps1` puede estar bloqueado por la política de
  ejecución; usar `npm.cmd` y `openspec.cmd` como alternativa.
