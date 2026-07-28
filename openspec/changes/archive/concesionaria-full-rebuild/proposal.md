# Propuesta: Reconstrucción completa del sistema de concesionaria

## Contexto

El proyecto actual tiene una implementación funcional pero desorganizada. Las specs existen
como documentación separada del código y no siguen el ciclo SDD. Se propone reconstruir
todo el sistema siguiendo OpenSpec: las specs como fuente de verdad, cambios atómicos
con proposal → design → specs → tasks → código → archive.

## Requisitos funcionales

- RF1: Autenticación con JWT (register, login, me) y roles (cliente, vendedor, administrador)
- RF2: Catálogo público de vehículos con filtros, paginación y ordenamiento
- RF3: CRUD de vehículos para administradores
- RF4: Consultas de clientes a vendedores con respuestas
- RF5: Solicitud y gestión de test drives con validación de superposición
- RF6: Reserva de vehículos con cambio de estado y confirmación/cancelación
- RF7: Dashboard administrativo con métricas y gráficos
- RF8: Chatbot IA con LangChain + OpenAI para consultas del inventario

## Requisitos no funcionales

- RNF1: Stack Go (Gin + GORM) + PostgreSQL + React (Vite + TypeScript + TailwindCSS)
- RNF2: API REST con endpoints protegidos por JWT y roles
- RNF3: Arquitectura backend en capas: Handler → Service → Repository → GORM
- RNF4: Frontend con React Router, componentes reutilizables, diseño responsive
- RNF5: Docker Compose para desarrollo local

## Criterios de aceptación

- CA1: Todos los endpoints responden correctamente según su spec
- CA2: El frontend compila sin errores de TypeScript
- CA3: El backend compila sin errores de Go
- CA4: Docker Compose levanta los 3 servicios (db, backend, frontend)
