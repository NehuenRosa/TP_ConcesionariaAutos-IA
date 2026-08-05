# Propuesta: CU-07 Turno de test drive

## Why

Los clientes interesados en una unidad no tienen forma de coordinar una prueba
de manejo: deben contactar al concesionario por otro medio. Hoy la aplicación
solo permite consultar sobre un vehículo, pero no agendar un turno para
manejarlo. Se necesita un flujo que permita al cliente reservar un horario para
un test drive de una unidad concreta y que el vendedor pueda gestionar esos
turnos, evitando que dos personas prueben la misma unidad a la misma hora.

## What Changes

- **Modelo de turno de test drive**: nueva entidad `TurnoTestDrive` con
  vehículo, cliente, fecha, franja horaria, estado y vendedor asignado.
- **Solicitud desde el catálogo**: el cliente autenticado puede pedir un turno
  desde el detalle de un vehículo disponible, eligiendo fecha y franja horaria.
- **Validación de superposición**: el sistema rechaza con error `409` cualquier
  turno que se solape con otro ya existente para la misma unidad en la misma
  fecha y franja horaria.
- **Gestión por el vendedor**: una página para el vendedor con el listado de
  turnos agendados, con filtro por estado, que permite confirmar o cancelar
  turnos.
- **Mis turnos (cliente)**: una página donde el cliente ve sus turnos y puede
  cancelar los que estén pendientes o confirmados.
- **Franjas horarias predefinidas**: catálogo de franjas (por ejemplo mañana y
  tarde) para acotar las opciones y facilitar la validación de solapamiento.
- **API REST** bajo `/api/test-drives` con endpoints de solicitud, listado
  (cliente y vendedor) y transiciones de estado (confirmar, cancelar).
- **Migración**: nueva tabla `turnos_test_drive` en la base de datos.

## Capabilities

### New Capabilities
- `turno-test-drive`: Solicitud, validación de disponibilidad (sin
  superposición para la misma unidad en la misma fecha y franja), gestión de
  estados por parte del vendedor y cancelación por parte del cliente.

### Modified Capabilities
- `catalogo-vehiculos`: el detalle de un vehículo disponible muestra la acción
  de "solicitar test drive" para clientes autenticados, manteniendo la
  referencia al vehículo del catálogo público.

## Impact

- **Backend**: nuevo modelo `TurnoTestDrive`, repository, service y handler de
  test drives; registro de rutas en el router; alta en la auto-migración.
- **Frontend**: nuevas páginas de solicitud, "Mis turnos" (cliente) y gestión
  de turnos (vendedor); modificaciones al detalle del vehículo y a la
  navegación.
- **Base de datos**: una tabla nueva (`turnos_test_drive`).
- **Autenticación**: se reutiliza el JWT existente (cliente y vendedor).
- **Dependencias**: ninguna nueva; se mantiene el stack (Go + Gin + GORM,
  React + Vite + TS).
