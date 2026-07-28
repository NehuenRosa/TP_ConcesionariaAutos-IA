# Tareas — Reconstrucción completa

## Backend

- [x] **B1: Config + Middleware** — Config desde env, CORS, Auth JWT, Role middleware
- [x] **B2: Modelos** — User, Vehicle, Consultation, ConsultationResponse, TestDrive, Reservation
- [x] **B3: Repositorios** — CRUD para cada modelo con queries específicas
- [x] **B4: Servicios** — Lógica de negocio para cada dominio
- [x] **B5: Handlers** — Handlers HTTP con validación y respuestas
- [x] **B6: Rutas** — Registro de rutas, migraciones, seed data
- [x] **B7: Main** — Entry point con conexión DB, migrations, server start

## Frontend

- [x] **F1: Config** — Vite, Tailwind, TypeScript, postcss, index.html, index.css
- [x] **F2: Types + Services** — Interfaces TypeScript, servicios HTTP con Axios
- [x] **F3: Auth Context + Hooks** — AuthContext con login/register/logout, useAuth hook
- [x] **F4: Shared Components** — Layout, Navbar, ProtectedRoute, Pagination, Chatbot
- [x] **F5: Catalog Components** — VehicleCard, FilterPanel, ImageGallery
- [x] **F6: Auth Pages** — Login, Register
- [x] **F7: Public Pages** — Catalog, VehicleDetail, ContactSeller, TestDriveRequest, ReserveVehicle
- [x] **F8: Admin Pages** — Dashboard, VehicleManagement, VehicleForm
- [x] **F9: Seller Pages** — ConsultationInbox, TestDriveManagement, ReservationManagement
- [x] **F10: App + Main** — Router configuration, entry point

## Infraestructura

- [x] **I1: Docker Compose** — db, backend, frontend servicios
- [x] **I2: Dockerfiles** — backend (multi-stage Go), frontend (Node)
- [x] **I3: .env.example** — Variables de entorno documentadas
