# Design: cu11-login-con-google

## Context

CU-01 ya resuelve registro/login propios con JWT HS256 propio
(claims `usuario_id` + `rol`), contraseñas hasheadas con bcrypt y email
único en `usuarios`. El frontend guarda el token en localStorage y todas las
rutas protegidas consumen `Authorization: Bearer <token>`. No existe hoy
ningún proveedor externo de identidad. La sesión es API-pura: no hay cookies
ni redirects de servidor hacia el frontend, salvo el SPA de React.

Restricciones relevantes:
- Stack cerrado: Go + Gin + GORM + JWT / React + Vite + TS. Dependencias
  nuevas requieren justificación; `golang.org/x/oauth2` y
  `github.com/golang-jwt/jwt/v5` ya están en `go.mod`/`go.sum`.
- Convención de idioma español para todo el código nuevo.
- Capas: handler → service → repository; prohibido saltearlas.
- `GOOGLE_API_KEY` ya existe en config pero es del chatbot (Gemini): las
  claves nuevas no deben colisionar semánticamente con ella.

## Goals / Non-Goals

**Goals:**
- Permitir crear cuenta e iniciar sesión con una cuenta de Google desde
  `/login` y `/registro`.
- Emitir al final el mismo JWT propio del sistema: cero cambios en los
  demás endpoints ni en la sesión existente.
- Vincular automáticamente cuentas existentes por email, sin duplicados.
- Degradación clara si Google no está configurado o no disponible.

**Non-Goals:**
- No se implementa Authorization Code Flow ni refresh tokens de Google.
- No se agregan otros proveedores (Apple, Microsoft, etc.), aunque el campo
  `proveedor` deja la puerta abierta.
- No se modifica el flujo email/contraseña ni el seed de usuarios.
- No hay gestión administrativa de identidades vinculadas (desvincular,
  etc.).

## Decisions

### D1: Google Identity Services (ID token) en vez de Authorization Code Flow

El botón oficial de Google corre en el navegador, obtiene el `credential`
(ID token firmado) y el frontend lo manda por POST al backend. El backend
verifica firma/emisor/audiencia/expiración y responde con su JWT propio.

- **Por qué**: la sesión vive en localStorage y la API es stateless; el flujo
  por código exigiría `GOOGLE_CLIENT_SECRET`, un callback de redirect y
  estado temporal entre pasos. El ID token mantiene todo request/response y
  evita secretos en el servidor.
- **Alternativas descartadas**: Authorization Code Flow (más piezas para el
  mismo resultado aquí); OAuth implícito (deprecado).

### D2: Verificación del ID token con `jwt/v5` + JWKS de Google (sin dependencias nuevas)

Nuevo paquete interno `internal/googleid` que:
1. Parsea el ID token como JWS RS256 con `github.com/golang-jwt/jwt/v5`
   (ya usado para el JWT propio).
2. Obtiene las claves públicas de `https://www.googleapis.com/oauth2/v3/certs`
   (JWKS estándar), las cachea en memoria con TTL (~1 h, refresco ante `kid`
   desconocido) y verifica la firma convirtiendo `n`/`e` del JWK a
   `rsa.PublicKey` con stdlib (`encoding/base64`, `math/big`, `crypto/rsa`).
3. Valida claims: `iss` ∈ {`accounts.google.com`,
   `https://accounts.google.com`}, `aud == GOOGLE_CLIENT_ID`, `exp` vigente
   (leeway 2 min), `email_verified = true`. Extrae `sub`, `email`, `name`.

- **Por qué**: cumple la regla de no sumar dependencias fuera del stack;
  la alternativa oficial `google.golang.org/api/idtoken` arrastra muchas
  dependencias transitivas. El endpoint `tokeninfo` de Google se descarta:
  suma latencia por pedido y Google lo desaconseja en producción.

### D3: Modelo Usuario extendido

```go
type Usuario struct {
    // ...campos actuales...
    Password  string  `gorm:"not null" json:"-"`        // pasa a permitir "" para cuentas solo-Google
    Proveedor string  `gorm:"not null;default:local" json:"proveedor"` // local | google
    GoogleSub *string `gorm:"uniqueIndex" json:"-"`     // NULL permite múltiples cuentas locales
}
```

- `Proveedor` distingue el origen; los usuarios existentes quedan `local`
  por el default (AutoMigrate agrega columnas sin tocar datos).
- `GoogleSub` como `*string` con índice único: PostgreSQL admite múltiples
  NULL, así que los usuarios locales no chocan entre sí. La vinculación real
  es por email (único); `sub` queda como identificador estable de Google.
- Contraseña vacía para cuentas Google: `bcrypt.CompareHashAndPassword("",
  ...)` falla, así que `POST /auth/login` sigue respondiendo `401` genérico
  sin cambios de lógica.

### D4: Lógica de negocio en el servicio de autenticación

Nuevo método `IniciarSesionConGoogle(ctx, credencial string)
(models.Usuario, string, error)` en `AutenticacionService`:

1. Verifica el ID token (D2). Error → `ErrCredencialInvalida`.
2. Busca usuario por email:
   - Existe y `proveedor=google`: actualiza `GoogleSub` si cambió.
   - Existe y `proveedor=local`: vincula (setea `proveedor=google` +
     `GoogleSub`), conserva rol/password/nombre. Un solo UPDATE.
   - No existe: crea usuario rol `cliente` con nombre/email de Google,
     `password=""`.
3. Genera el JWT propio reutilizando `token.Generar` y devuelve
   `RespuestaLogin` igual que el login tradicional.

El repositorio de usuarios suma `Guardar(usuario)` genérico si aún no existe
un método adecuado. El handler nuevo solo parsea el body, delega y mapea
errores a códigos HTTP (`400` sin credencial, `401` inválida, `503` no
configurado/caída de Google).

### D5: Configuración

- `GOOGLE_CLIENT_ID`: habilita el flujo cuando está definido. Vacío ⇒ el
  backend responde `503` en `/api/auth/google` y el frontend oculta el botón.
- Se expone además `GET /api/auth/proveedores` (público) que devuelve qué
  métodos están habilitados (`{"google": true/false}`) para que el frontend
  no dependa de duplicar la variable en build. Alternativa descartada:
  `VITE_GOOGLE_CLIENT_ID` (duplicaría configuración y puede desincronizarse
  del backend).
- `GOOGLE_CLIENT_SECRET` no es necesario para este flujo; se documenta como
  opcional/reservado.
- `.env.example` y `docker-compose.yml`: claves nuevas documentadas; sin
  defaults inseguros (el flujo simplemente queda deshabilitado si faltan).

### D6: Frontend

- Nuevo componente `src/components/BotonGoogle.tsx`: consulta
  `/api/auth/proveedores`; si Google está habilitado, carga el script oficial
  `https://accounts.google.com/gsi/client` dinámicamente, inicializa GIS con
  el client ID recibido del backend y renderiza el botón oficial.
- Al recibir el `credential`, llama `api.iniciarSesionGoogle(credencial)` →
  `guardarToken(token)` + recarga de perfil vía el contexto `useAuth`
  (misma mecánica que `registrar()`), redirige a `state.desde` o `/`.
- `<BotonGoogle />` se incluye en `InicioSesion.tsx` y `Registro.tsx`,
  separado del formulario por el divisor "o continuá con". Errores → mismo
  patrón de mensaje en español que ya usan esas páginas.
- Sin rutas/callbacks nuevos: todo ocurre en las páginas existentes.

## Risks / Trade-offs

- [Google inaccesible al validar un token (fetch de JWKS)] → caché de claves
  en memoria con TTL y reintento ante `kid` desconocido; error `503` con
  mensaje claro si no hay claves disponibles.
- [Deriva de reloj rechaza tokens válidos] → leeway de 2 minutos en `exp`.
- [Confusión de identidad por emails reutilizados] → solo se aceptan
  credentials con `email_verified=true` emitidos a nuestro `aud`.
- [Usuario vinculado pierde acceso si su cuenta Google se compromete] →
  fuera de alcance; mitigable a futuro con verificación de contraseña al
  vincular (queda como mejora posterior).
- [`Password=""` rompe supuestos futuros de "todo usuario tiene hash"] →
  decisión explícita documentada (D3); cualquier código nuevo debe tratar
  cadena vacía como "sin contraseña local".
- [Acoplamiento del frontend a la respuesta de `/proveedores`] → contrato
  mínimo y versionado implícito: si el endpoint falla, el botón no se
  muestra y el formulario tradicional sigue funcionando.
- [Costo/latencia del script GIS de terceros] → carga diferida solo cuando
  Google está habilitado; el formulario propio nunca depende de él.

## Migration Plan

1. AutoMigrate agrega `proveedor` (default `local`) y `google_sub` (NULL);
   ningún dato existente cambia de comportamiento.
2. Deploy backend (con o sin `GOOGLE_CLIENT_ID`) y luego frontend; ambos son
   compatibles en cualquier orden porque el botón se oculta si
   `/api/auth/proveedores` dice que Google está off.
3. Rollback: quitar `GOOGLE_CLIENT_ID` del entorno deshabilita el flujo sin
   deploy; las columnas nuevas son aditivas y no afectan el resto.
