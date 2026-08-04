# Roadmap del proyecto

Backlog de casos de uso pendientes de planificar con OpenSpec.
Cada caso de uso se planifica con un change propio cuando se vaya a implementar.

| ID | Caso de uso | Estado |
|----|-------------|--------|
| CU-01 | Autenticación y roles (registro, login, JWT, 3 roles) | Implementado |
| CU-02 | Gestión de vehículos (ABM) | Implementado |
| CU-03 | Catálogo público (listado paginado + detalle) | Implementado |
| CU-04 | Búsqueda y filtrado avanzado | Implementado |
| CU-05 | Consulta / cotización | Implementado |
| CU-06 | Gestión de consultas (bandeja del vendedor) | Implementado (junto con CU-05) |
| CU-07 | Turno de test drive | Pendiente |
| CU-08 | Reserva de vehículo | Pendiente |
| CU-09 | Panel de administración (dashboard) | Pendiente |
| CU-10 | Chatbot asistente (LangChain) | Pendiente |

## Orden sugerido de implementación

1. **CU-01** Autenticación y roles (base para autorizar el resto).
2. **CU-02** Gestión de vehículos (alimenta el catálogo).
3. **CU-03** Catálogo público (depende de CU-02).
4. **CU-04** Búsqueda y filtrado (amplía CU-03).
5. **CU-05 / CU-06** Consultas y su gestión.
6. **CU-07** Test drives.
7. **CU-08** Reservas.
8. **CU-09** Panel de administración.
9. **CU-10** Chatbot.

## Proceso para planificar un CU

1. Crear el change: `openspec new change <nombre-en-kebab-case>`
2. Completar los artefactos: proposal → specs → design → tasks.
3. Validar: `openspec validate --strict`
4. Implementar y archivar al finalizar.
