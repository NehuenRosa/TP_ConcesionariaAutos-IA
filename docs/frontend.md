# Frontend (React + Vite + TypeScript)

Aplicación web del catálogo y la gestión de la concesionaria. Framework: React
+ Vite + TypeScript, React Router para el enrutado y TailwindCSS para los
estilos.

## Estructura

```
frontend/src/
├── App.tsx               # Raíz: envuelve las rutas con el proveedor de auth
├── main.tsx              # Punto de entrada (mount de React)
├── index.css             # Estilos globales (Tailwind)
├── components/           # Componentes reutilizables
│   ├── Chatbot.tsx       # Widget flotante del asistente (chat + tasación)
│   ├── ChatConsulta.tsx  # Hilo de mensajes de una consulta (cliente/vendedor)
│   ├── RutaProtegida.tsx # Guarda de rutas por rol
│   ├── graficos/         # Gráficos del panel (GraficoBarras, TarjetaMetrica)
│   └── ui/               # Primitivas de UI (Boton, Campo, Paginacion, etc.)
├── hooks/
│   ├── useAuth.tsx       # Contexto de autenticación (login, logout, rol, token)
│   └── useNotificaciones.ts  # Contador de mensajes no leídos (polling)
├── layouts/
│   └── LayoutBase.tsx    # Header + footer + <Outlet/>; inyecta el Chatbot
├── pages/                # Una página por ruta
├── routes/
│   └── Rutas.tsx         # Definición de rutas y protección por rol
├── services/
│   └── api.ts            # Cliente HTTP centralizado (fetch + token JWT)
├── types/                # Tipos TypeScript (vehiculo, usuario, consulta, etc.)
├── utils/
│   └── formato.ts        # Formato de moneda y fechas en español
└── test/
    └── setup.ts          # Setup de Vitest + Testing Library
```

## Rutas

| Ruta | Página | Acceso |
|------|--------|--------|
| `/` | `Inicio` | Público |
| `/catalogo` | `Catalogo` | Público |
| `/catalogo/:id` | `DetalleVehiculo` | Público |
| `/catalogo/:id/test-drive` | `FormularioTestDrive` | `cliente` |
| `/catalogo/:id/reservar` | `FormularioReserva` | `cliente` |
| `/login` | `InicioSesion` | Público |
| `/registro` | `Registro` | Público |
| `/admin` | `PanelAdministracion` | `administrador` |
| `/admin/vehiculos` | `AdminVehiculos` | `administrador` |
| `/admin/vehiculos/nuevo` | `FormularioVehiculo` | `administrador` |
| `/admin/vehiculos/:id/editar` | `FormularioVehiculo` | `administrador` |
| `/admin/usuarios` | `GestionUsuarios` | `administrador` |
| `/mis-consultas` | `MisConsultas` | `cliente` |
| `/mis-consultas/:id` | `MisConsultas` | `cliente` |
| `/vendedor/bandeja` | `BandejaEntrada` | `vendedor` |
| `/vendedor/bandeja/:id` | `ChatVendedor` | `vendedor` |
| `/vendedor/test-drives` | `GestionTestDrives` | `vendedor` |
| `/vendedor/reservas` | `GestionReservas` | `vendedor` |
| `/mis-test-drives` | `MisTestDrives` | `cliente` |
| `/mis-reservas` | `MisReservas` | `cliente` |
| `*` | `NoEncontrada` | Público |

Las rutas privadas se protegen con `<RutaProtegida rol="...">` que redirige a
`/login` si no hay sesión o a `/` si el rol no coincide.

## Cliente HTTP (`services/api.ts`)

- Expone `api` con un método por endpoint del backend.
- Agrega automáticamente `Authorization: Bearer <token>` desde localStorage
  (clave `token_concesionaria`).
- Normaliza errores: lanza `ErrorApi` con el mensaje en español del backend
  (fallback a "Ocurrió un error inesperado...").
- Maneja `204 No Content` (eliminaciones) sin intentar parsear JSON.
- Timeout de 140 s para todas las peticiones (el chatbot del backend puede
  tardar hasta 120 s).
- `peticionMultipart` envía `multipart/form-data` con las fotos de la tasación
  (`fotos`) y la `descripcion` opcional, sin fijar `Content-Type` (el navegador
  agrega el boundary).

## Autenticación (`hooks/useAuth.tsx`)

- `ProveedorAutenticacion` envuelve la app (ver `App.tsx`).
- Expone el usuario actual, `iniciarSesion`, `registrar`, `cerrarSesion` y el
  helper `rolActual`.
- Al iniciar sesión guarda el token en localStorage y refetch del perfil.
- Se usa en el header del layout para mostrar menú según rol.

## Notificaciones (`hooks/useNotificaciones.ts`)

- Polling periódico a `GET /notificaciones/contador`.
- Muestra un badge con la cantidad de mensajes no leídos en el layout.

## Chatbot (`components/Chatbot.tsx`)

Widget flotante visible en todas las páginas (inyectado desde `LayoutBase`):

- **Chat**: `api.enviarMensajeChatbot` con `{ mensaje, historial }`; mantiene el
  historial de la conversación en el estado del componente.
- **Tasación**: sube hasta 5 fotos + descripción opcional con
  `api.enviarTasacion`; muestra el resultado en el hilo del chat.

Tipos en `types/chatbot.ts`. Ver `AGENTS.md` (sección CU-10) para el flujo de
valores reales y `docs/api.md` para el detalle de los endpoints.

## Tests

Framework: Vitest + React Testing Library + jsdom. Los tests viven junto a las
páginas/archivos (`*.test.tsx`) y en `services/api.test.ts`.

```powershell
cd frontend
npm run test          # una pasada
npm run test -- --watch   # modo watch
npm run build         # build de producción (typecheck incluido)
```

Tests existentes:

- `api.test.ts` — normalización de errores, header de auth, 204.
- `InicioSesion.test.tsx` — formulario de login (errores y submit).
- `FormularioTestDrive.test.tsx` — validación de fecha/franja.
- `ChatVendedor.test.tsx` — carga y envío de mensajes.
- `MisConsultas.test.tsx` — listado y apertura de consultas.

## Configuración

- `VITE_API_URL` (variable de entorno): URL base de la API. Por defecto
  `http://localhost:8080/api`. Se compila en tiempo de build (se pasa como
  `args` en `docker-compose.yml`).
- `vite.config.ts`: plugin React, alias y config de Vitest (jsdom).
- `Dockerfile`: build multi-etapa con nginx sirviendo el `dist` en el puerto 80.
