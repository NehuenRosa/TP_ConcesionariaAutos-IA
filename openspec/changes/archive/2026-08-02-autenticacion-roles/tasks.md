## 1. Backend: modelo y migración

- [x] 1.1 Crear `internal/models/usuario.go` con `Usuario` (nombre, email único, password hasheado, rol) y constantes de rol (`cliente`, `vendedor`, `administrador`).
- [x] 1.2 Agregar `&models.Usuario{}` a `AutoMigrar` en `internal/database/database.go`.

## 2. Backend: dependencias y paquete token

- [x] 2.1 Agregar `github.com/golang-jwt/jwt/v5` y `golang.org/x/crypto/bcrypt` al módulo.
- [x] 2.2 Crear `internal/token/token.go` con `Generar(usuarioID, rol, secreto, duracion)` y `Validar(token, secreto)` (claims `sub`, `rol`, `exp`).

## 3. Backend: repository de usuarios

- [x] 3.1 Crear `internal/repositories/usuarios.go` con `Crear`, `ObtenerPorEmail` y `ObtenerPorID`.

## 4. Backend: service de autenticación

- [x] 4.1 Crear `internal/services/autenticacion.go` con `Registrar` (valida, hashea, asigna rol `cliente`), `IniciarSesion` (verifica credenciales y emite token) y `ObtenerPorID`.
- [x] 4.2 Definir errores de negocio en español (`ErrEmailEnUso`, `ErrCredencialesInvalidas`, `ErrDatosRegistroInvalidos`, `ErrUsuarioNoEncontrado`).

## 5. Backend: handlers y rutas de autenticación

- [x] 5.1 Crear `internal/handlers/autenticacion.go` con `Registrar`, `IniciarSesion` y `Perfil` (DTOs sin contraseña, códigos `400`/`401`/`409`/`500`).
- [x] 5.2 Registrar en `internal/router/router.go` las rutas `POST /api/auth/registro`, `POST /api/auth/login` y `GET /api/auth/perfil` (protegida).

## 6. Backend: middlewares reales y seed

- [x] 6.1 Reemplazar la implementación de `internal/middleware/auth.go`: `AutenticacionJWT(secreto)` valida el Bearer token y guarda `usuario_id` y `rol` en el contexto; `ExigirRol(rol)` responde `403` si no coincide (el admin accede a todo).
- [x] 6.2 Crear el seed de usuarios por defecto (administrador y vendedor) ejecutado al arrancar el backend.

## 7. Frontend: tipos, cliente HTTP y sesión

- [x] 7.1 Crear `src/types/usuario.ts` con `Rol`, `Usuario`, `DatosRegistro`, `DatosLogin` y `RespuestaLogin`.
- [x] 7.2 Agregar a `src/services/api.ts` los métodos `registrar`, `iniciarSesion` y `obtenerPerfil`, y el manejo del token (guardar/eliminar en localStorage y encabezado `Authorization`).
- [x] 7.3 Crear `src/hooks/useAuth.tsx` con `AuthProvider` y hook `useAuth` (carga el perfil, expone iniciarSesion/registrar/cerrarSesion/esAdministrador).

## 8. Frontend: páginas y rutas protegidas

- [x] 8.1 Implementar `src/pages/InicioSesion.tsx` con el formulario de login y mensajes de error en español.
- [x] 8.2 Crear `src/pages/Registro.tsx` con el formulario de registro.
- [x] 8.3 Crear un componente `RutaProtegida` que redirija al login sin sesión y bloquee roles insuficientes.
- [x] 8.4 Registrar `/registro` y envolver las rutas `/admin/*` con `RutaProtegida` (rol `administrador`) en `src/routes/Rutas.tsx`; actualizar la navegación del layout según la sesión.

## 9. Verificación

- [x] 9.1 `cd backend && go build ./...` y `go vet ./...` sin errores.
- [x] 9.2 `cd frontend && npm run build` sin errores.
- [x] 9.3 Verificar de punta a punta con Docker: registro de cliente, login de cliente/vendedor/administrador, acceso a `/api/admin/vehiculos` (401 sin token, 403 cliente, 200 administrador) y redirección en el frontend.

