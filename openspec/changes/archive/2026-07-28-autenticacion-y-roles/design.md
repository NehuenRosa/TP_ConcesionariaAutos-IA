# Diseño: Autenticación y Roles

## Arquitectura

```
┌─────────────────────────────────────────────────────────┐
│                     Frontend (React)                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │
│  │Login.tsx │  │Register  │  │   AuthContext.tsx     │  │
│  │          │  │ .tsx     │  │  user, token, login,  │  │
│  │          │  │          │  │  register, logout,    │  │
│  │          │  │          │  │  isAuthenticated,     │  │
│  │          │  │          │  │  isAdmin, isSeller    │  │
│  └────┬─────┘  └────┬─────┘  └──────────┬───────────┘  │
│       │             │                   │              │
│       └──────┬──────┘                   │              │
│              │                          │              │
│     ┌────────▼────────┐                │              │
│     │ authService.ts  │                │              │
│     │ (axios -> api)  │                │              │
│     └────────┬────────┘                │              │
│              │                          │              │
│     ┌────────▼────────┐                │              │
│     │  ProtectedRoute  │◄───────────────┘              │
│     │  .tsx            │  allowedRoles[]               │
│     └─────────────────┘                                │
└──────────────────────────┬──────────────────────────────┘
                           │ HTTP (JSON)
                           │ Authorization: Bearer <token>
┌──────────────────────────▼──────────────────────────────┐
│                   Backend (Go/Gin)                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │
│  │  Auth    │  │  Auth    │  │   UserRepository     │  │
│  │ Handler  │──▶ Service  │──▶  (GORM -> PostgreSQL)│  │
│  └──────────┘  └──────────┘  └──────────────────────┘  │
│                                                     │
│  ┌──────────────────────────────────────────────────┐ │
│  │  Middleware                                       │ │
│  │  AuthMiddleware: valida JWT, inyecta contexto     │ │
│  │  RoleMiddleware:  verifica rol contra lista       │ │
│  └──────────────────────────────────────────────────┘ │
│                                                     │
│  ┌──────────────────────────────────────────────────┐ │
│  │  Router (register.go)                             │ │
│  │  POST /api/auth/register  → público               │ │
│  │  POST /api/auth/login     → público               │ │
│  │  GET  /api/auth/me        → autenticado           │ │
│  └──────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Modelo de Datos

### User (models/user.go)

| Campo    | Tipo     | Tags                          |
|----------|----------|-------------------------------|
| ID       | uint     | primaryKey                    |
| Name     | string   | size:100; not null            |
| Email    | string   | size:100; uniqueIndex; not null |
| Password | string   | size:255; not null; json:"-"  |
| Role     | UserRole | size:20; not null; default:cliente |
| Phone    | string   | size:20                       |

### Constantes de Rol

```
RoleClient = "cliente"
RoleSeller = "vendedor"
RoleAdmin  = "administrador"
```

### Requests / Responses

```
RegisterRequest { name, email, password(min:6), phone? }
LoginRequest    { email, password(min:6) }
AuthResponse    { token: string, user: User }
```

## Endpoints

| Método | Ruta              | Auth     | Roles | Descripción         |
|--------|-------------------|----------|-------|---------------------|
| POST   | /api/auth/register | No       | -     | Registro de usuario |
| POST   | /api/auth/login    | No       | -     | Inicio de sesión    |
| GET    | /api/auth/me       | Sí       | todos | Perfil propio       |

### Flujo Login
1. Cliente envía `POST /api/auth/login` con `{email, password}`
2. `AuthService.Login` busca usuario por email
3. Compara password con bcrypt
4. Genera JWT con claims `{user_id, role, email, exp}`
5. Responde `{token, user}`

### Flujo Registro
1. Cliente envía `POST /api/auth/register` con `{name, email, password, phone?}`
2. `AuthService.Register` verifica email único
3. Hashea password con bcrypt
4. Crea usuario con `Role = cliente`
5. Genera JWT
6. Responde `{token, user}`

## Componentes Frontend

| Componente       | Ruta               | Rol requerido |
|------------------|--------------------|---------------|
| Login.tsx        | /login             | público       |
| Register.tsx     | /register          | público       |
| AuthContext.tsx   | - (provider)       | -             |
| ProtectedRoute   | - (wrapper)        | configurable  |
| Navbar.tsx       | - (layout)         | adaptativo    |

### Árbol de Componentes

```
<BrowserRouter>
  <AuthProvider>
    <Routes>
      <Route element={<Layout />}>
        <!-- Públicas -->
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <!-- Protegidas (cualquier rol autenticado) -->
        <Route element={<ProtectedRoute />}>
          <Route path="/consultar/:id" ... />
        </Route>

        <!-- Solo admin -->
        <Route element={<ProtectedRoute allowedRoles={['administrador']} />}>
          <Route path="/admin/*" ... />
        </Route>

        <!-- Vendedor + admin -->
        <Route element={<ProtectedRoute allowedRoles={['vendedor','administrador']} />}>
          <Route path="/seller/*" ... />
        </Route>
      </Route>
    </Routes>
  </AuthProvider>
</BrowserRouter>
```

## Flujo de Datos

### Login
```
Login.tsx → authService.login()
         → api.post('/auth/login')
         → AuthHandler.Login → AuthService.Login
         → UserRepository.FindByEmail → bcrypt.CompareHashAndPassword
         → generateToken() → JWT
         → AuthResponse → AuthContext (localStorage + estado)
         → ProtectedRoute verifica isAuthenticated + allowedRoles
         → api.ts interceptor agrega Bearer token a cada request
```

### Validación de token al recargar
```
App carga → AuthContext useEffect
         → localStorage tiene token?
           → Sí: GET /auth/me
             → 200: setUser(data)
             → 401: logout(), redirigir a /login
           → No: estado inicial (no autenticado)
```

## Middleware (autorización)

### AuthMiddleware
```
1. Extraer header "Authorization: Bearer <token>"
2. Parsear y validar JWT con jwtSecret
3. Si inválido/expirado → 401
4. Si válido → c.Set("user_id", claims.UserID)
              c.Set("role", claims.Role)
              c.Set("email", claims.Email)
              c.Next()
```

### RoleMiddleware
```
1. Leer "role" del contexto Gin
2. Comparar contra lista de roles permitidos
3. Si no coincide → 403 Forbidden
4. Si coincide → c.Next()
```
