## Why

El sistema necesita identificar a los usuarios y restringir el acceso a las
funcionalidades según su perfil (cliente, vendedor o administrador). Hoy los
middlewares de autenticación y rol son stubs que dejan pasar todas las
peticiones, y el ABM de vehículos (CU-02) ya quedó registrado tras ellos pero
sin control real. Sin CU-01, cualquier persona podría manipular el stock y no
hay base para los casos de uso posteriores (consultas del vendedor, reservas,
panel de administración).

## What Changes

- Se agrega la entidad `Usuario` con email único, contraseña hasheada y rol.
- Se agregan los endpoints públicos de autenticación:
  - `POST /api/auth/registro`: crea un usuario con rol `cliente`.
  - `POST /api/auth/login`: valida credenciales y emite un token JWT de 24 h.
  - `GET /api/auth/perfil`: devuelve el usuario autenticado (protegido).
- Se implementan los middlewares reales `AutenticacionJWT` (valida el Bearer
  token) y `ExigirRol` (verifica el rol del claim), reemplazando los stubs.
- Las rutas administrativas (`/api/admin/vehiculos`) quedan protegidas con rol
  `administrador`.
- Se agrega seed de usuarios por defecto (vendedor y administrador) para poder
  operar desde el primer arranque.
- En el frontend: páginas de registro e inicio de sesión, contexto de sesión que
  guarda el token y el perfil, cliente HTTP que adjunta el token, y protección
  de rutas administrativas con redirección al login.
- **BREAKING**: los middlewares cambian su contrato: `AutenticacionJWT` ahora
  valida el token y rechaza peticiones sin credenciales válidas.

## Capabilities

### New Capabilities
- `autenticacion-roles`: registro e inicio de sesión de usuarios, emisión y
  validación de tokens JWT con claims de rol, y autorización de rutas por rol
  (cliente, vendedor, administrador).

### Modified Capabilities
<!-- Sin cambios de requisitos a nivel spec: gestion-vehiculos ya preveía la protección por rol al registrarse las rutas administrativas. -->

## Impact

- **Backend**: se agregan `models/usuario.go`, `repositories/usuarios.go`,
  `services/autenticacion.go`, `handlers/autenticacion.go` y un paquete `token`
  (o funciones en services) para emitir/validar JWT. Se reemplaza la
  implementación de `middleware/auth.go`. Se agrega `Usuario` a la
  auto-migración y un seed de usuarios por defecto.
- **Dependencias**: se agrega `github.com/golang-jwt/jwt/v5` para tokens JWT y
  `golang.org/x/crypto/bcrypt` para el hash de contraseñas (ambas dentro del
  stack de autenticación del proyecto).
- **API**: tres endpoints nuevos; los endpoints administrativos existentes
  pasan a exigir autenticación y rol reales.
- **Frontend**: se agregan `types/usuario.ts`, métodos de auth en `api.ts`,
  contexto de sesión (`hooks/useAuth.tsx` o similar), páginas de registro e
  inicio de sesión, y un wrapper de rutas protegidas.
- **Base de datos**: nueva tabla `usuarios` (auto-migración).
