# Spec: autenticacion-roles

## Purpose

Permitir que los usuarios se registren e inicien sesión en el sistema, y
restringir el acceso a las funcionalidades según su rol (cliente, vendedor,
administrador) mediante tokens JWT con claims de rol.

## ADDED Requirements

### Requirement: Registro de usuario

El sistema SHALL exponer un endpoint público `POST /api/auth/registro` que cree
un usuario con rol `cliente`, dados `nombre`, `email` y `password`. El email
SHALL ser único, el password SHALL tener una longitud mínima y la contraseña
SHALL guardarse hasheada. El sistema SHALL ignorar cualquier rol enviado en el
request y asignar siempre `cliente`. La respuesta SHALL incluir el usuario sin
la contraseña.

#### Scenario: Registro exitoso

- **WHEN** un visitante envía nombre, email y password válidos
- **THEN** el sistema crea el usuario con rol `cliente` y responde con su
  perfil sin la contraseña

#### Scenario: Email ya registrado

- **WHEN** un visitante se registra con un email que ya existe
- **THEN** el sistema responde con error `409` y un mensaje en español

#### Scenario: Datos de registro inválidos

- **WHEN** un visitante envía nombre vacío, email mal formado o password muy
  corto
- **THEN** el sistema responde con error `400` y un mensaje en español

### Requirement: Inicio de sesión

El sistema SHALL exponer un endpoint público `POST /api/auth/login` que, dados
`email` y `password` correctos, responda con un token JWT y el perfil del
usuario. El token SHALL incluir el rol del usuario y SHALL expirar a las 24
horas.

#### Scenario: Credenciales correctas

- **WHEN** un usuario envía un email y una contraseña válidos
- **THEN** el sistema responde con un token JWT y el perfil del usuario

#### Scenario: Credenciales incorrectas

- **WHEN** un usuario envía un email inexistente o una contraseña equivocada
- **THEN** el sistema responde con error `401` y un mensaje en español sin
  revelar cuál dato fue incorrecto

### Requirement: Perfil del usuario autenticado

El sistema SHALL exponer un endpoint `GET /api/auth/perfil`, protegido con
autenticación JWT, que devuelva el perfil del usuario autenticado.

#### Scenario: Perfil con token válido

- **WHEN** un usuario autenticado solicita su perfil
- **THEN** el sistema responde con el perfil del usuario del token

#### Scenario: Perfil sin token o con token inválido

- **WHEN** un visitante solicita el perfil sin token, con token vencido o con
  token inválido
- **THEN** el sistema responde con error `401` y un mensaje en español

### Requirement: Autorización de rutas por rol

El sistema SHALL validar el token JWT en toda petición protegida y SHALL
restringir las rutas administrativas al rol `administrador`. Las peticiones sin
token o con token inválido SHALL recibir `401`; las de un rol no autorizado
SHALL recibir `403`. El rol `administrador` SHALL tener acceso a cualquier ruta
protegida.

#### Scenario: Acceso administrativo autorizado

- **WHEN** un administrador autenticado accede a una ruta administrativa
- **THEN** el sistema permite la petición

#### Scenario: Acceso administrativo sin token

- **WHEN** un visitante accede a una ruta administrativa sin token
- **THEN** el sistema responde con error `401` y un mensaje en español

#### Scenario: Acceso administrativo con token vencido

- **WHEN** un usuario envía un token vencido a una ruta administrativa
- **THEN** el sistema responde con error `401` y un mensaje en español

#### Scenario: Acceso administrativo con rol insuficiente

- **WHEN** un cliente autenticado accede a una ruta administrativa
- **THEN** el sistema responde con error `403` y un mensaje en español

### Requirement: Usuarios por defecto

El sistema SHALL crear al arranque, si no existen, un usuario con rol
`administrador` y un usuario con rol `vendedor`, con credenciales de desarrollo
documentadas, para permitir operar el sistema desde el primer arranque.

#### Scenario: Seed en el primer arranque

- **WHEN** el backend arranca con la tabla `usuarios` vacía
- **THEN** el sistema crea los usuarios por defecto de administrador y vendedor

#### Scenario: Seed idempotente

- **WHEN** el backend arranca y los usuarios por defecto ya existen
- **THEN** el sistema no los duplica

### Requirement: Sesión en el frontend

El sistema SHALL ofrecer páginas de registro (`/registro`) e inicio de sesión
(`/login`) que consuman los endpoints de autenticación, guarden el token
localmente y carguen el perfil del usuario. El cliente HTTP SHALL adjuntar el
token en las peticiones autenticadas.

#### Scenario: Registro desde el frontend

- **WHEN** un visitante completa el formulario de registro con datos válidos
- **THEN** el sistema crea la cuenta y el visitante queda identificado en la
  aplicación

#### Scenario: Inicio de sesión desde el frontend

- **WHEN** un usuario completa el formulario de login con credenciales válidas
- **THEN** el sistema guarda el token y muestra la sesión del usuario

#### Scenario: Cierre de sesión

- **WHEN** un usuario cierra sesión
- **THEN** el sistema descarta el token y el usuario queda como visitante

### Requirement: Protección de rutas administrativas en el frontend

El sistema SHALL redirigir al login a los visitantes que intenten acceder a una
ruta administrativa sin sesión, y SHALL impedir el acceso a los usuarios cuyo
rol no sea suficiente.

#### Scenario: Visitante sin sesión en ruta administrativa

- **WHEN** un visitante accede a `/admin/vehiculos` sin sesión
- **THEN** el sistema redirige a `/login` con un mensaje que pide iniciar sesión

#### Scenario: Cliente en ruta administrativa

- **WHEN** un usuario con rol `cliente` accede a una ruta administrativa
- **THEN** el sistema no muestra el panel y orienta al usuario
