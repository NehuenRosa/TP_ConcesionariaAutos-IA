# Autenticación y Roles de Usuario

## Contexto

El sistema de concesionaria requiere un sistema de autenticación que permita
registro e inicio de sesión de usuarios, distinguiendo tres roles (cliente,
vendedor, administrador). Cada rol tiene acceso a funcionalidades específicas
que deben protegerse mediante autorización basada en JWT.

Actualmente no existe ningún sistema de autenticación implementado.

## Requisitos Funcionales

1. **Registro de usuarios**: formulario con nombre, email, contraseña y teléfono
   opcional. Al registrarse, el usuario obtiene rol `cliente` por defecto.
2. **Inicio de sesión**: autenticación por email y contraseña con bcrypt.
3. **JWT**: generación de token con claims `user_id`, `role`, `email` y expiración
   configurable.
4. **Middleware de autenticación**: validación del token JWT en cada request
   protegido, inyectando datos del usuario en el contexto.
5. **Middleware de roles**: restricción de endpoints según rol.
6. **Perfil propio (`/me`)**: endpoint que devuelve datos del usuario autenticado.
7. **Seed data**: 3 usuarios de prueba (admin, vendedor, cliente).
8. **Frontend Login/Register**: formularios con validación y manejo de errores.
9. **Protección de rutas frontend**: redirección si no autenticado o sin rol
   permitido.
10. **Persistencia de sesión**: token en localStorage, restauración al recargar.

## Requisitos No Funcionales

- Contraseñas hasheadas con bcrypt (costo por defecto).
- JWT firmado con HMAC-SHA256.
- El rol solo puede asignarse como `cliente` vía registro público.
- Los roles `vendedor` y `administrador` se asignan desde semilla o DB.
- Interceptor HTTP que agregue token automáticamente y maneje 401.
- Componentes protegidos con `ProtectedRoute` por rol.

## Criterios de Aceptación

- [ ] Un usuario puede registrarse y queda con rol `cliente`.
- [ ] Un usuario registrado puede iniciar sesión y recibe un JWT válido.
- [ ] El token JWT expira según lo configurado.
- [ ] Un endpoint público (`/catalogo`) no requiere autenticación.
- [ ] Un endpoint protegido (`/admin/*`) rechaza sin token o con rol incorrecto.
- [ ] El frontend redirige a `/login` si no hay token.
- [ ] El frontend redirige a `/` si el rol no está permitido.
- [ ] Al recargar la página, la sesión se restaura si el token sigue siendo
      válido.
- [ ] Los 3 usuarios de semilla pueden iniciar sesión correctamente.
