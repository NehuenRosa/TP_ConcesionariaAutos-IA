# Autenticación y Autorización

## Requirements

### Requirement: Registro de usuarios
- **WHEN** un usuario envía POST /api/auth/register con name, email y password
- **THEN** el sistema crea un usuario con rol "cliente", password hasheado con bcrypt, y retorna 201 con token JWT + datos del usuario

- **WHEN** un usuario envía POST /api/auth/register con un email ya registrado
- **THEN** el sistema retorna 409 con mensaje de error

- **WHEN** un usuario envía POST /api/auth/register con datos inválidos (email sin @, password < 6 chars)
- **THEN** el sistema retorna 400 con error de validación

### Requirement: Inicio de sesión
- **WHEN** un usuario envía POST /api/auth/login con email y password correctos
- **THEN** el sistema retorna 200 con token JWT + datos del usuario

- **WHEN** un usuario envía POST /api/auth/login con credenciales incorrectas
- **THEN** el sistema retorna 401 con mensaje de error

### Requirement: Obtener usuario actual
- **WHEN** un usuario autenticado envía GET /api/auth/me con Bearer token válido
- **THEN** el sistema retorna 200 con datos del usuario autenticado

- **WHEN** un usuario no autenticado envía GET /api/auth/me
- **THEN** el sistema retorna 401

### Requirement: Autorización por roles
- **WHEN** un usuario sin rol administrador intenta crear/editar/eliminar vehículos
- **THEN** el sistema retorna 403

- **WHEN** un usuario sin rol vendedor/administrador intenta listar todas las consultas
- **THEN** el sistema retorna 403
