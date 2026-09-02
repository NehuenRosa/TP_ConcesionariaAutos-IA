# Tasks: cu11-login-con-google

## 1. Backend — modelo y configuración

- [x] 1.1 Extender `internal/models/usuario.go`: agregar `Proveedor` (default `local`) y `GoogleSub *string` con `uniqueIndex`; ajustar tag de `Password` para permitir cadena vacía en cuentas solo-Google
- [x] 1.2 Agregar a `internal/config/config.go`: `GoogleClientID` (y `GoogleClientSecret` opcional) leído de `GOOGLE_CLIENT_ID`/`GOOGLE_CLIENT_SECRET`, sin defaults inseguros
- [x] 1.3 Actualizar `.env.example` y `docker-compose.yml` con las variables nuevas documentadas (flujo deshabilitado si faltan)

## 2. Backend — verificación del ID token

- [x] 2.1 Crear paquete `internal/googleid`: parseo RS256 del ID token con `github.com/golang-jwt/jwt/v5`, fetch y caché en memoria del JWKS de Google (`https://www.googleapis.com/oauth2/v3/certs`, TTL ~1 h, refresco ante `kid` desconocido)
- [x] 2.2 Validar claims en `internal/googleid`: emisor, audiencia == client ID, expiración con leeway de 2 min y `email_verified=true`; devolver identidad (`sub`, `email`, `nombre`)
- [x] 2.3 Escribir tests unitarios de `googleid` generando un par RSA propio y mockeando el JWKS: token válido, firma inválida, audiencia incorrecta, expirado y email sin verificar

## 3. Backend — servicio, repositorio y endpoint

- [x] 3.1 Agregar al repositorio de usuarios lo necesario para guardar cambios de vinculación (método genérico `Guardar` si no existe) *(ya existían `Crear` y `Actualizar`, suficientes)*
- [x] 3.2 Implementar `IniciarSesionConGoogle` en `services/autenticacion.go`: verificar credencial → buscar por email → vincular (conservando rol/password) o crear cliente nuevo con `proveedor=google` → emitir JWT propio vía `token.Generar`
- [x] 3.3 Agregar handler `IniciarSesionConGoogle` en `handlers/autenticacion.go` con mapeo de errores: `400` sin credencial, `401` inválida, `503` Google no configurado/indisponible; respuestas `{"error": "..."}` en español
- [x] 3.4 Registrar en `router.go`: ruta pública `POST /api/auth/google` y `GET /api/auth/proveedores` (`{"google": bool}`)

## 4. Frontend — botón e integración

- [x] 4.1 Agregar en `types/` el tipo del credential/respuesta y en `services/api.ts` las funciones `iniciarSesionGoogle(credencial)` y `obtenerProveedores()`
- [x] 4.2 Crear `src/components/BotonGoogle.tsx`: consulta `/api/auth/proveedores`, carga diferida del script GIS, renderiza el botón oficial y emite el credential al callback
- [x] 4.3 Integrar `<BotonGoogle />` en `InicioSesion.tsx` y `Registro.tsx`: guardar token, recargar perfil desde `useAuth`, redirigir a `state.desde` o `/`, mostrar error en español sin perder el formulario
- [x] 4.4 Verificar que el botón no aparece cuando Google está deshabilitado (respuesta `{"google": false}` o error del endpoint)

## 5. Verificación y documentación

- [x] 5.1 `go build ./...`, `go vet ./...` y tests de backend en verde
- [x] 5.2 `npm run build` y tests de frontend en verde
- [ ] 5.3 Prueba E2E manual con credenciales reales de Google (localhost): alta nueva como cliente, ingreso recurrente y vinculación con usuario creado por formulario; verificar JWT propio y rol en `/api/auth/perfil`
- [ ] 5.4 Prueba de degradación: sin `GOOGLE_CLIENT_ID`, `/api/auth/google` responde `503`, `/api/auth/proveedores` devuelve `{"google": false}` y el login tradicional sigue funcionando
- [x] 5.5 Agregar CU-11 a `docs/roadmap.md` (estado Implementado al cerrar) y actualizar AGENTS.md/`.env.example` si hace falta
