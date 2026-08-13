# Propuesta: CU-08 Reserva de vehículo

## Why

Los clientes no tienen forma de asegurar una unidad del catálogo: hoy solo pueden
consultar y pedir test drives, pero un vehículo puede venderse mientras tanto. El
chatbot de CU-10 ya orienta a "reservar", pero esa acción no existe en el sistema.
Se necesita un flujo que permita al cliente reservar una unidad disponible
(bloqueándola) y que el vendedor confirme la venta o libere la unidad.

## What Changes

- **Modelo de reserva**: nueva entidad `Reserva` con vehículo, cliente y estado
  (`activa`, `vendida`, `cancelada`).
- **Solicitud desde el catálogo**: el cliente autenticado puede reservar una
  unidad disponible desde el detalle del vehículo. Al crear la reserva, la unidad
  pasa de `disponible` a `reservado` y deja de mostrarse en el catálogo público.
- **Cambio de estado transaccional**: la creación de la reserva y el cambio de
  estado del vehículo se persisten en una única transacción (igual que la
  confirmación de venta y la cancelación).
- **Gestión por el vendedor**: una página con el listado de reservas, con filtro
  por estado, que permite confirmar la venta (vehículo → `vendido`) o cancelar y
  liberar la unidad (vehículo → `disponible`).
- **Mis reservas (cliente)**: una página donde el cliente ve sus reservas y puede
  cancelar las que estén activas, liberando la unidad.
- **API REST** bajo `/api/reservas` con endpoints de creación, listado (cliente y
  vendedor) y transiciones de estado (confirmar venta, cancelar).
- **Migración**: nueva tabla `reservas` en la base de datos.

## Capabilities

### New Capabilities
- `reserva-vehiculo`: Creación de una reserva que bloquea un vehículo disponible,
  cancelación por parte del cliente, gestión por parte del vendedor (confirmación
  de venta o cancelación que libera la unidad).

### Modified Capabilities
- `catalogo-vehiculos`: el detalle de un vehículo disponible muestra la acción de
  "reservar" para clientes autenticados, y las unidades reservadas dejan de
  aparecer en el catálogo público (ya se cumple con la regla de solo `disponible`).
- `gestion-vehiculos`: los estados `reservado` y `vendido` ya existen; la gestión
  de stock refleja las transiciones producidas por el flujo de reservas.

## Impact

- **Backend**: nuevo modelo `Reserva`, repository (con operaciones
  transaccionales que actualizan la reserva y el estado del vehículo), service y
  handler de reservas; registro de rutas en el router; alta en la auto-migración.
- **Frontend**: nuevas páginas de reserva, "Mis reservas" (cliente) y gestión de
  reservas (vendedor); modificaciones al detalle del vehículo y a la navegación.
- **Base de datos**: una tabla nueva (`reservas`).
- **Autenticación**: se reutiliza el JWT existente (cliente y vendedor).
- **Dependencias**: ninguna nueva; se mantiene el stack (Go + Gin + GORM,
  React + Vite + TS).
