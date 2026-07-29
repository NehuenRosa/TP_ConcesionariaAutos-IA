# Diseño del Sistema

## Arquitectura

```
┌─────────────────────────────────────────────────────┐
│                  FRONTEND (React + Vite)              │
│  Auth | Catálogo | Cliente | Vendedor | Admin | Chat │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP (Axios) + JWT
┌──────────────────────┴──────────────────────────────┐
│                 BACKEND (Go + Gin)                    │
│  Handler → Service → Repository → GORM               │
└──────────────────────┬──────────────────────────────┘
                       │
            ┌──────────┴──────────┐
            │    PostgreSQL 16     │
            │   (concesionaria)    │
            └─────────────────────┘
```

## Modelos de datos

### User
| Campo | Tipo |
|-------|------|
| ID | uint pk |
| Name | string |
| Email | string unique |
| Password | string (bcrypt) |
| Role | enum: cliente, vendedor, administrador |
| Phone | string |

### Vehicle
| Campo | Tipo |
|-------|------|
| ID | uint pk |
| Brand | string |
| Model | string |
| Year | int |
| Price | float64 |
| Mileage | int |
| Fuel | enum: nafta, diesel, electrico, hibrido |
| Transmission | enum: manual, automatico |
| Condition | enum: nuevo, usado |
| Color | string |
| Description | text |
| Images | string[] |
| Status | enum: disponible, reservado, vendido |
| VehicleType | string |

### Consultation
| Campo | Tipo |
|-------|------|
| ID | uint pk |
| ClientID | uint FK→User |
| VehicleID | uint FK→Vehicle |
| Message | text |
| Status | enum: pendiente, en_conversacion, cerrada |
| AssignedTo | *uint FK→User |
| HasUnreadMessages | bool (default:false) |
| HasUnreadForClient | bool (default:false) |

### ConsultationResponse
| Campo | Tipo |
|-------|------|
| ID | uint pk |
| CreatedAt | timestamp |
| ConsultationID | uint FK |
| UserID | uint FK |
| Message | text |

### TestDrive
| Campo | Tipo |
|-------|------|
| ID | uint pk |
| ClientID | uint FK→User |
| VehicleID | uint FK→Vehicle |
| ScheduledAt | datetime |
| Status | enum: pendiente, confirmado, cancelado, completado |
| Notes | text |

### Reservation
| Campo | Tipo |
|-------|------|
| ID | uint pk |
| ClientID | uint FK→User |
| VehicleID | uint FK→Vehicle |
| Status | enum: activa, confirmada, cancelada |
| Notes | text |

## Endpoints

### Auth
- POST /api/auth/register → 201 {token, user}
- POST /api/auth/login → 200 {token, user}
- GET /api/auth/me → 200 {user}

### Vehicles (público)
- GET /api/vehicles → 200 paginado (solo disponibles)
- GET /api/vehicles/brands → 200 marcas
- GET /api/vehicles/:id → 200 detalle

### Vehicles (admin)
- POST /api/vehicles → 201 (JWT + admin)
- PUT /api/vehicles/:id → 200 (JWT + admin)
- DELETE /api/vehicles/:id → 200 (JWT + admin)

### Consultations
- POST /api/consultations → 201 (JWT)
- GET /api/consultations/mine → 200 (JWT)
- GET /api/consultations → 200 (JWT + seller/admin)
- GET /api/consultations/pending/count → 200 (JWT + seller/admin)
- GET /api/consultations/notifications/count → 200 (JWT) — total según rol
- GET /api/consultations/:id → 200 (JWT) — marca como leído según el rol
- PATCH /api/consultations/:id/status → 200 (JWT + seller/admin)
- POST /api/consultations/:id/responses → 200 (JWT)
- DELETE /api/consultations/:id → 200 (JWT) — dueño o seller/admin

### Test Drives
- POST /api/test-drives → 201 (JWT)
- GET /api/test-drives/mine → 200 (JWT)
- GET /api/test-drives → 200 (JWT + seller/admin)
- GET /api/test-drives/:id → 200 (JWT)
- PATCH /api/test-drives/:id/status → 200 (JWT + seller/admin)

### Reservations
- POST /api/reservations → 201 (JWT)
- GET /api/reservations/mine → 200 (JWT)
- GET /api/reservations → 200 (JWT + seller/admin)
- GET /api/reservations/:id → 200 (JWT)
- POST /api/reservations/:id/confirm → 200 (JWT + seller/admin)
- POST /api/reservations/:id/cancel → 200 (JWT + seller/admin)

### Admin
- GET /api/admin/dashboard → 200 (JWT + admin)

### Chatbot
- GET /api/chatbot/status → 200
- POST /api/chatbot/ask → 200

## Frontend Routes

| Ruta | Componente | Acceso |
|------|-----------|--------|
| /catalogo | Catalog | público |
| /vehiculos/:id | VehicleDetail | público |
| /login | Login | público |
| /register | Register | público |
| /consultar/:id | ContactSeller | cliente |
| /mis-consultas | MyConsultations | cliente |
| /test-drive/:id | TestDriveRequest | cliente |
| /reservar/:id | ReserveVehicle | cliente |
| /admin/dashboard | Dashboard | admin |
| /admin/vehiculos | VehicleManagement | admin |
| /admin/vehiculos/nuevo | VehicleForm | admin |
| /admin/vehiculos/:id | VehicleForm | admin |
| /seller/consultas | ConsultationInbox | seller/admin |
| /seller/test-drives | TestDriveManagement | seller/admin |
| /seller/reservas | ReservationManagement | seller/admin |
