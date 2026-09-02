# autenticacion-google Specification

## Purpose

Autenticación federada con Google Identity Services: el frontend envía el
credential (ID token), el backend lo verifica contra el JWKS de Google
(firma, emisor, audiencia, `email_verified`), crea el cliente o vincula la
cuenta existente por email y emite el JWT propio del sistema.

## Requirements

### Requirement: Endpoint de inicio de sesión con Google

El sistema SHALL exponer el endpoint público `POST /api/auth/google` que
recibe el ID token (credential) emitido por Google Identity Services y,
si es válido, responde `200` con el mismo formato del inicio de sesión
existente (`{token, usuario}`), donde `token` es un JWT propio del sistema
con claims `usuario_id` y `rol`, vigencia de 24 horas, utilizable en todos
los endpoints protegidos existentes sin cambios.

#### Scenario: Primer ingreso con una cuenta de Google nueva

- **WHEN** se envía a `POST /api/auth/google` el credential de una cuenta de
  Google verificada que no corresponde a ningún usuario registrado
- **THEN** el sistema responde `200` con un JWT propio y el usuario creado
  con rol `cliente`

#### Scenario: Ingreso recurrente con la misma cuenta de Google

- **WHEN** se envía el credential de una cuenta de Google ya vinculada a un
  usuario existente
- **THEN** el sistema responde `200` con un JWT propio para ese usuario,
  conservando su rol actual, sin duplicar usuarios

### Requirement: Verificación del ID token en el backend

El backend SHALL validar criptográficamente cada ID token recibido antes de
autenticar al usuario: verificar la firma contra los certificados públicos
de Google, que el emisor sea `accounts.google.com` o `https://accounts.google.com`,
que la audiencia coincida exactamente con `GOOGLE_CLIENT_ID`, que no esté
expirado y que el email esté verificado por Google. El backend MUST NOT
confiar en datos de identidad sin verificar (ni en campos enviados por el
cliente aparte del credential).

#### Scenario: Credential con firma inválida

- **WHEN** se envía un ID token cuya firma no proviene de los certificados
  de Google
- **THEN** el sistema responde `401` con cuerpo `{"error": "..."}` en español
  y no crea ni autentica ningún usuario

#### Scenario: Credential emitido para otra aplicación

- **WHEN** el ID token tiene una audiencia distinta de `GOOGLE_CLIENT_ID`
- **THEN** el sistema responde `401` y no inicia sesión

#### Scenario: Credential expirado

- **WHEN** el ID token está vencido según su claim `exp`
- **THEN** el sistema responde `401`

#### Scenario: Email sin verificar en Google

- **WHEN** el ID token indica `email_verified=false`
- **THEN** el sistema responde `401` y no crea ni vincula la cuenta

### Requirement: Alta automática como cliente

Cuando el credential es válido y no existe un usuario con ese email, el
sistema SHALL crear automáticamente un nuevo usuario con rol `cliente`, el
nombre y email informados por Google, `proveedor=google` y el identificador
estable de la cuenta (`google_sub`) persistido. La contraseña queda vacía
para estas cuentas y el campo no se expone nunca en las respuestas.

#### Scenario: Alta de cliente vía Google

- **WHEN** un visitante completa el flujo de Google con una cuenta sin
  usuario asociado
- **THEN** queda creado un usuario activo de rol `cliente`, con
  `proveedor=google` y su `google_sub` guardado, y puede operar como
  cualquier otro cliente (consultas, reservas, test drives)

### Requirement: Vinculación automática de cuentas por email

Si el email del credential válido coincide con un usuario existente creado
por formulario (proveedor `local`), el sistema SHALL vincular la identidad
de Google a ese usuario (setear `proveedor`/`google_sub`) y responder con el
JWT de ese usuario, conservando nombre, email, rol y contraseña actuales,
sin crear duplicados.

#### Scenario: Cliente existente ingresa por primera vez con Google

- **WHEN** el email del credential coincide con un usuario creado
  previamente por registro con contraseña
- **THEN** el sistema vincula la identidad de Google a ese mismo usuario y
  responde `200` con su JWT y su rol original

#### Scenario: Sin duplicados

- **WHEN** dos ingresos consecutivos con Google corresponden al mismo email
  ya vinculado
- **THEN** el sistema mantiene un único usuario con ese email

### Requirement: Manejo de errores del flujo Google

El endpoint `POST /api/auth/google` SHALL manejar de forma diferenciada:
credencial ausente o malformada → `400`; credencial inválida (firma,
audiencia, expiración o email no verificado) → `401`; Google no configurado
(`GOOGLE_CLIENT_ID` vacío) o servicio de certificados de Google
indisponible → `503` con mensaje claro en español. Ninguna respuesta de
error MUST revelar detalles internos ni stack traces.

#### Scenario: Petición sin credencial

- **WHEN** se llama a `POST /api/auth/google` sin el campo del credential
- **THEN** el sistema responde `400` con mensaje en español

#### Scenario: Google no configurado

- **WHEN** el backend arrancó sin `GOOGLE_CLIENT_ID` y llega una petición al
  endpoint de Google
- **THEN** el sistema responde `503` indicando que el acceso con Google no
  está disponible

### Requirement: Botón Continuar con Google en el frontend

Las páginas `/login` y `/registro` SHALL mostrar el botón oficial
"Continuar con Google" cuando el acceso con Google esté habilitado. Al
completarse el flujo de Google con éxito, el frontend SHALL almacenar el
token devuelto igual que en el inicio de sesión tradicional, cargar el
perfil y redirigir a la página de origen (o inicio), quedando el usuario
con sesión iniciada. Si ocurre un error, se muestra un mensaje en español
sin perder el formulario actual.

#### Scenario: Login completo desde el botón

- **WHEN** el usuario presiona "Continuar con Google" en `/login` y
  selecciona su cuenta
- **THEN** queda con sesión iniciada, redirigido a la página desde la que
  venía, y las peticiones siguientes incluyen el JWT propio del sistema

#### Scenario: Registro desde el botón

- **WHEN** el usuario presiona "Continuar con Google" en `/registro` con
  una cuenta de Google sin usuario asociado
- **THEN** se crea la cuenta como cliente y queda con sesión iniciada sin
  completar el formulario manual

#### Scenario: Error durante el flujo de Google

- **WHEN** el popup de Google se cierra o el backend rechaza el credential
- **THEN** el usuario permanece en la página y ve un mensaje de error claro
  en español
