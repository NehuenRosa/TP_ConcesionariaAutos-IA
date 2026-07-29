# Tareas: Autenticación y Roles

## Tarea 1: Modelo User y migración

- [x] Crear modelo `User` en `internal/models/user.go` con campos ID, Name, Email,
      Password, Role, Phone, timestamps
- [x] Definir constantes `RoleClient`, `RoleSeller`, `RoleAdmin`
- [x] Definir structs `LoginRequest`, `RegisterRequest`, `AuthResponse`
- [x] Registrar en `AutoMigrate` en `router.go`

**Verificación**: `go build ./...` compila sin errores.

## Tarea 2: UserRepository

- [x] Crear `internal/repositories/user_repo.go`
- [x] Implementar `Create`, `FindByEmail`, `FindByID`, `List`, `Update`, `Delete`

**Verificación**: `go build ./...` compila sin errores.

## Tarea 3: AuthService

- [x] Crear `internal/services/auth_service.go`
- [x] Implementar `Register`: verificar email único, hashear con bcrypt, crear
      usuario con rol `cliente`, generar JWT
- [x] Implementar `Login`: buscar por email, comparar bcrypt, generar JWT
- [x] Implementar `generateToken`: crear JWT con claims user_id, role, email, exp

**Verificación**: `go build ./...` compila sin errores.

## Tarea 4: AuthHandler

- [x] Crear `internal/handlers/auth_handler.go`
- [x] Implementar `Register` (POST /api/auth/register)
- [x] Implementar `Login` (POST /api/auth/login)
- [x] Implementar `Me` (GET /api/auth/me) — devuelve datos del token

**Verificación**: `go build ./...` compila sin errores.

## Tarea 5: Middleware de autenticación y roles

- [x] Crear `internal/middleware/auth.go`
- [x] Implementar `AuthMiddleware`: extraer Bearer token, parsear JWT, validar,
      inyectar user_id/role/email en contexto Gin
- [x] Implementar `RoleMiddleware`: leer role del contexto, comparar contra roles
      permitidos, 403 si no coincide
- [x] Crear `internal/middleware/cors.go` si no existe

**Verificación**: `go build ./...` compila sin errores.

## Tarea 6: Rutas y seed

- [x] Registrar rutas de auth en `internal/routes/register.go`
- [x] Agregar rutas públicas: POST /auth/register, POST /auth/login
- [x] Agregar ruta protegida: GET /auth/me (requiere AuthMiddleware)
- [x] Crear 3 usuarios de prueba en `internal/seed/seed.go` con bcrypt:
      admin, vendedor, cliente

**Verificación**: `go build ./...` compila sin errores.

## Tarea 7: Servicios HTTP frontend (api.ts + authService.ts)

- [x] Crear `src/services/api.ts` con axios, baseURL configurable, interceptor
      que agregue Bearer token automáticamente
- [x] Agregar interceptor de respuesta que limpie sesión en 401
- [x] Crear `src/services/authService.ts` con métodos `login`, `register`, `me`
- [x] Definir tipos `User`, `AuthResponse` en `src/types/index.ts`

**Verificación**: `npx tsc --noEmit` = 0 errores.

## Tarea 8: AuthContext y useAuth

- [x] Crear `src/context/AuthContext.tsx` con provider que exponga:
      user, token, login, register, logout, isAuthenticated, isAdmin, isSeller
- [x] En `useEffect`, si hay token en localStorage, llamar a `me()` para
      verificar vigencia y restaurar sesión
- [x] Crear `src/hooks/useAuth.ts`

**Verificación**: `npx tsc --noEmit` = 0 errores.

## Tarea 9: ProtectedRoute

- [x] Crear `src/components/ProtectedRoute.tsx`
- [x] Si no autenticado → redirigir a /login
- [x] Si `allowedRoles` definido y rol no incluido → redirigir a /
- [x] Si pasa ambas → renderizar `<Outlet />`

**Verificación**: `npx tsc --noEmit` = 0 errores.

## Tarea 10: Páginas Login y Register

- [x] Crear `src/pages/Login.tsx` con formulario email/password, llamar a
      `useAuth().login()`, mostrar errores, redirigir en éxito
- [x] Crear `src/pages/Register.tsx` con formulario
      nombre/email/password/teléfono, llamar a `useAuth().register()`
- [x] Manejar estados: loading, error, éxito

**Verificación**: `npx tsc --noEmit` = 0 errores.

## Tarea 11: Integrar en App.tsx y Layout

- [x] Envolver rutas con `<AuthProvider>`
- [x] Agregar rutas `/login` y `/register` como públicas
- [x] Aplicar `ProtectedRoute` en rutas que requieran autenticación/roles
- [x] Actualizar `Navbar.tsx` para mostrar enlaces según rol:
      Todos: Catálogo
      Autenticado: Mis consultas, Mis reservas
      Admin: Dashboard, Gestionar vehículos
      Vendedor: Bandeja consultas, Turnos, Reservas

**Verificación**: `npm run build` = 0 errores.

## Tarea 12: Verificación final

- [x] Ejecutar `go build ./...` en backend
- [x] Ejecutar `npx tsc --noEmit` en frontend
- [x] Ejecutar `npm run build` en frontend

**Verificación**: Los 3 comandos devuelven 0 errores.
