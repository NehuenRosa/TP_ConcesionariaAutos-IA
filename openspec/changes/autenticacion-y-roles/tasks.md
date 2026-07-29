# Tareas: Autenticación y Roles

## Tarea 1: Modelo User y migración

- [ ] Crear modelo `User` en `internal/models/user.go` con campos ID, Name, Email,
      Password, Role, Phone, timestamps
- [ ] Definir constantes `RoleClient`, `RoleSeller`, `RoleAdmin`
- [ ] Definir structs `LoginRequest`, `RegisterRequest`, `AuthResponse`
- [ ] Registrar en `AutoMigrate` en `router.go`

**Verificación**: `go build ./...` compila sin errores.

## Tarea 2: UserRepository

- [ ] Crear `internal/repositories/user_repo.go`
- [ ] Implementar `Create`, `FindByEmail`, `FindByID`, `List`, `Update`, `Delete`

**Verificación**: `go build ./...` compila sin errores.

## Tarea 3: AuthService

- [ ] Crear `internal/services/auth_service.go`
- [ ] Implementar `Register`: verificar email único, hashear con bcrypt, crear
      usuario con rol `cliente`, generar JWT
- [ ] Implementar `Login`: buscar por email, comparar bcrypt, generar JWT
- [ ] Implementar `generateToken`: crear JWT con claims user_id, role, email, exp

**Verificación**: `go build ./...` compila sin errores.

## Tarea 4: AuthHandler

- [ ] Crear `internal/handlers/auth_handler.go`
- [ ] Implementar `Register` (POST /api/auth/register)
- [ ] Implementar `Login` (POST /api/auth/login)
- [ ] Implementar `Me` (GET /api/auth/me) — devuelve datos del token

**Verificación**: `go build ./...` compila sin errores.

## Tarea 5: Middleware de autenticación y roles

- [ ] Crear `internal/middleware/auth.go`
- [ ] Implementar `AuthMiddleware`: extraer Bearer token, parsear JWT, validar,
      inyectar user_id/role/email en contexto Gin
- [ ] Implementar `RoleMiddleware`: leer role del contexto, comparar contra roles
      permitidos, 403 si no coincide
- [ ] Crear `internal/middleware/cors.go` si no existe

**Verificación**: `go build ./...` compila sin errores.

## Tarea 6: Rutas y seed

- [ ] Registrar rutas de auth en `internal/routes/register.go`
- [ ] Agregar rutas públicas: POST /auth/register, POST /auth/login
- [ ] Agregar ruta protegida: GET /auth/me (requiere AuthMiddleware)
- [ ] Crear 3 usuarios de prueba en `internal/seed/seed.go` con bcrypt:
      admin, vendedor, cliente

**Verificación**: `go build ./...` compila sin errores.

## Tarea 7: Servicios HTTP frontend (api.ts + authService.ts)

- [ ] Crear `src/services/api.ts` con axios, baseURL configurable, interceptor
      que agregue Bearer token automáticamente
- [ ] Agregar interceptor de respuesta que limpie sesión en 401
- [ ] Crear `src/services/authService.ts` con métodos `login`, `register`, `me`
- [ ] Definir tipos `User`, `AuthResponse` en `src/types/index.ts`

**Verificación**: `npx tsc --noEmit` = 0 errores.

## Tarea 8: AuthContext y useAuth

- [ ] Crear `src/context/AuthContext.tsx` con provider que exponga:
      user, token, login, register, logout, isAuthenticated, isAdmin, isSeller
- [ ] En `useEffect`, si hay token en localStorage, llamar a `me()` para
      verificar vigencia y restaurar sesión
- [ ] Crear `src/hooks/useAuth.ts`

**Verificación**: `npx tsc --noEmit` = 0 errores.

## Tarea 9: ProtectedRoute

- [ ] Crear `src/components/ProtectedRoute.tsx`
- [ ] Si no autenticado → redirigir a /login
- [ ] Si `allowedRoles` definido y rol no incluido → redirigir a /
- [ ] Si pasa ambas → renderizar `<Outlet />`

**Verificación**: `npx tsc --noEmit` = 0 errores.

## Tarea 10: Páginas Login y Register

- [ ] Crear `src/pages/Login.tsx` con formulario email/password, llamar a
      `useAuth().login()`, mostrar errores, redirigir en éxito
- [ ] Crear `src/pages/Register.tsx` con formulario
      nombre/email/password/teléfono, llamar a `useAuth().register()`
- [ ] Manejar estados: loading, error, éxito

**Verificación**: `npx tsc --noEmit` = 0 errores.

## Tarea 11: Integrar en App.tsx y Layout

- [ ] Envolver rutas con `<AuthProvider>`
- [ ] Agregar rutas `/login` y `/register` como públicas
- [ ] Aplicar `ProtectedRoute` en rutas que requieran autenticación/roles
- [ ] Actualizar `Navbar.tsx` para mostrar enlaces según rol:
      Todos: Catálogo
      Autenticado: Mis consultas, Mis reservas
      Admin: Dashboard, Gestionar vehículos
      Vendedor: Bandeja consultas, Turnos, Reservas

**Verificación**: `npm run build` = 0 errores.

## Tarea 12: Verificación final

- [ ] Ejecutar `go build ./...` en backend
- [ ] Ejecutar `npx tsc --noEmit` en frontend
- [ ] Ejecutar `npm run build` en frontend

**Verificación**: Los 3 comandos devuelven 0 errores.
