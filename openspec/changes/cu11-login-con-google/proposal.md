# Proposal: cu11-login-con-google

## Why

Hoy el único modo de crear cuenta y entrar es el formulario propio de
email/contraseña (CU-01). Los clientes que ya tienen una cuenta de Google
tienen que registrar otra credencial para operar con la concesionaria, lo
que agrega fricción en el primer contacto y abandono de registros. Agregar
"Continuar con Google" reduce esa fricción reutilizando toda la
infraestructura existente de roles y JWT.

## What Changes

- Nuevo caso de uso **CU-11 Login con Google** (documentado en
  `docs/roadmap.md`): inicio de sesión y registro mediante Google Identity
  Services.
- El frontend incorpora el botón oficial de Google ("Continuar con Google")
  en las páginas `/login` y `/registro`; obtiene un ID token (credential) y
  lo envía al backend.
- Nuevo endpoint público `POST /api/auth/google` en el backend: verifica el
  ID token contra Google (firma, emisor, audiencia y expiración), busca o
  crea el usuario y responde con el mismo JWT propio que usa el resto del
  sistema (claims `usuario_id` + `rol`, 24 h).
- Vinculación automática de cuentas: si ya existe un usuario con ese email
  (creado por formulario), se vincula la identidad de Google al usuario
  existente sin duplicarlo.
- Modelo `Usuario` extendido: campo `proveedor` (`local` | `google`) y
  `google_sub` (identificador estable de Google); la contraseña deja de ser
  obligatoria para cuentas creadas vía Google.
- Rol asignado siempre `cliente` para usuarios nuevos vía Google; los
  usuarios vinculados conservan su rol actual.
- Configuración nueva: `GOOGLE_CLIENT_ID` (y opcionalmente
  `GOOGLE_CLIENT_SECRET`) sin tocar `GOOGLE_API_KEY` (que sigue siendo del
  chatbot/Gemini).
- Degradación controlada: si Google no está configurado o falla la
  verificación, el endpoint responde con error claro y el botón puede
  ocultarse; el flujo email/contraseña queda intacto.

## Capabilities

### New Capabilities

- `autenticacion-google`: inicio de sesión y registro federado con Google
  Identity Services. Cubre la verificación del ID token en el backend, el
  alta automática de clientes, la vinculación automática de cuentas por
  email, la emisión del JWT propio del sistema y la integración en el
  frontend (botón en login/registro).

### Modified Capabilities

_(ninguna: el flujo email/contraseña y los requirements vigentes de
`autenticacion-roles` no cambian; las cuentas solo-Google simplemente no
pueden iniciar sesión con contraseña, comportamiento ya cubierto por el
401 genérico)_

## Impact

- **Backend**
  - `internal/models/usuario.go`: campos `proveedor` y `google_sub`;
    `password` nullable (auto-migración GORM).
  - `internal/services/` nuevo servicio de Google (verificación de ID token,
    búsqueda/vinculación/alta de usuario) + método en repositorio de usuarios.
  - `internal/handlers/autenticacion.go`: nuevo endpoint
    `POST /api/auth/google`.
  - `internal/router/router.go`: ruta pública en el grupo `/auth`.
  - `internal/config/config.go`: claves `GOOGLE_CLIENT_ID` /
    `GOOGLE_CLIENT_SECRET`.
  - Dependencias Go: se usará `golang.org/x/oauth2` (ya presente como
    dependencia indirecta en `go.sum`) para obtener la identidad de Google;
    la verificación criptográfica del ID token se hace con
    `github.com/golang-jwt/jwt/v5` (ya en uso). Sin dependencias nuevas
    fuera del stack.
- **Frontend**
  - `src/pages/InicioSesion.tsx` y `src/pages/Registro.tsx`: botón
    "Continuar con Google".
  - Carga del script de Google Identity Services y envío del credential al
    nuevo endpoint; sesión idéntica a la actual (mismo token en
    localStorage).
- **Infra / config**: `.env.example` y `docker-compose.yml` con las claves
  nuevas (sin defaults inseguros).
- **API**: un endpoint público nuevo; ningún contrato existente cambia.
- **Docs**: `docs/roadmap.md` suma CU-11.
