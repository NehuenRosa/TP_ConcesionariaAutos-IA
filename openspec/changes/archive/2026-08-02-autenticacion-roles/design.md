## Context

El monorepo tiene backend Go (Gin + GORM) y frontend React con cliente HTTP
centralizado. Los middlewares `AutenticacionJWT` y `ExigirRol` existen pero son
stubs que dejan pasar todo; las rutas administrativas de CU-02 ya los usan. No
hay modelo de usuario ni librería JWT todavía. La contraseña se guardará
hasheada con bcrypt y el token firmado con el secreto `JWT_SECRETO` ya
configurado en el entorno.

## Goals / Non-Goals

**Goals:**

- Registrar usuarios con rol `cliente` e iniciar sesión emitiendo un token JWT
  de 24 horas con claims de rol.
- Validar el token en cada petición protegida y autorizar según rol.
- Proteger las rutas administrativas existentes con rol `administrador`.
- Proveer usuarios por defecto (seed) de vendedor y administrador.
- Sesión en el frontend: guardar token, cargar perfil, proteger rutas admin y
  redirigir al login.

**Non-Goals:**

- Recuperación/restablecimiento de contraseña por email.
- Gestión de usuarios (alta/edición de vendedores por el admin): corresponde a
  un CU futuro; por ahora se crean vía seed.
- Roles dinámicos/permisos granulares: solo los tres roles fijos.
- Refresh tokens: un único token de acceso por sesión.

## Decisions

### D1: JWT con `golang-jwt/jwt/v5` y HS256

Se usa `github.com/golang-jwt/jwt/v5` (la librería JWT estándar de Go),
firmando con HS256 y el secreto `JWT_SECRETO`. Los claims son `usuario_id`
(con la reclamación estándar `sub`), `rol` y `exp` (24 h). Alternativa:
`lestrrat-go/jwx`; se descarta por ser más pesada para el caso de uso.

### D2: Hash de contraseñas con bcrypt

`golang.org/x/crypto/bcrypt` (cost por defecto). Alternativa: Argon2id; se
descarta por requerir manejo manual de salt/parámetros, mientras que bcrypt ya
está en el módulo `x/crypto` y es suficiente para este sistema. La contraseña
nunca se devuelve en las respuestas (`json:"-"` en el modelo).

### D3: El registro público siempre crea rol `cliente`

El endpoint `POST /api/auth/registro` ignora cualquier rol recibido y asigna
`cliente`. Vendedores y administradores se crean con el seed del arranque (D4).
Alternativa: permitir elegir rol en el registro. Se descarta por seguridad:
cualquiera podría registrarse como administrador.

### D4: Seed de usuarios por defecto

Al arrancar, se crean (si no existen por email) dos usuarios con contraseñas
documentadas: `administrador@concesionaria.local` (rol `administrador`) y
`vendedor@concesionaria.local` (rol `vendedor`). Esto permite operar el sistema
desde el primer arranque sin una pantalla de gestión de usuarios.

### D5: Middlewares reales que reemplazan los stubs

`AutenticacionJWT(secreto)` lee el encabezado `Authorization: Bearer <token>`,
valida la firma y expiración con el paquete `token`, y deja `usuario_id` y
`rol` en el contexto de Gin. `ExigirRol(rol)` lee el `rol` del contexto y
responde `403` si no coincide. El rol `administrador` tiene acceso a cualquier
ruta (superusuario), además del rol exacto solicitado.

### D6: Paquete `internal/token` para emitir y validar tokens

La emisión (`Generar`) y validación (`Validar`) de tokens viven en un paquete
dedicado `internal/token` para que la usen tanto el service de autenticación
(emitir en login) como el middleware (validar en cada request), sin acoplar
middleware ↔ services.

### D7: Sesión en el frontend con contexto y localStorage

Se guarda el token en `localStorage` y el perfil en un contexto React
(`AuthProvider` + hook `useAuth`). El cliente HTTP agrega el encabezado
`Authorization` automáticamente cuando hay token. Un componente
`RutaProtegida` redirige a `/login` si no hay sesión y a `/` si el rol no es
suficiente.

## Risks / Trade-offs

- [Token en localStorage expuesto a XSS] → Mitigación: es el patrón habitual en
  SPAs; el token dura 24 h y el secreto se mantiene en el backend.
- [Seed con contraseñas por defecto] → Mitigación: son credenciales de
  desarrollo documentadas y se cambian en producción; sirven para operar sin
  pantalla de gestión de usuarios.
- [Rol del usuario estático en el token] → Mitigación: si se cambia el rol de un
  usuario, el token vigente conserva el anterior hasta expirar; aceptable dado
  que no hay gestión de usuarios aún.
- [bcrypt cost por defecto lento en equipos débiles] → Mitigación: se usa el
  cost por defecto (10), balance razonable entre seguridad y latencia.
- [Rutas admin exigir rol real rompe flujos dev] → Mitigación: los usuarios
  seed permiten entrar al panel desde el primer arranque; el frontend redirige
  al login con mensaje claro.

## Migration Plan

- Cambio aditivo: nueva tabla `usuarios` (auto-migración) y seed idempotente.
  No hay datos que migrar. Los tokens emitidos antes de este change no existen
  (no había emisión), así que no hay tokens que invalidar.

## Open Questions

- Si en el futuro el administrador gestiona usuarios (alta de vendedores), se
  migrará la creación de usuarios del seed a un endpoint administrativo.
