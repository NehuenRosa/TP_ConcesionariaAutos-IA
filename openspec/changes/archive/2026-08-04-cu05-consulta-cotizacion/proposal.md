# Propuesta: CU-05 Consulta / cotización sobre un vehículo

## Resumen

Implementar un sistema de consultas/cotizaciones que permita a los clientes
iniciar conversaciones con vendedores sobre vehículos específicos. El cliente
puede enviar una consulta desde el detalle del vehículo, y el vendedor puede
tomar esa consulta desde una bandeja de entrada para responder y gestionar la
conversación.

## Problema

Los clientes necesitan poder consultar información específica sobre vehículos
interesantes y establecer una comunicación directa con un vendedor para
cotizaciones, dudas o negociaciones. Actualmente no existe un canal estructurado
para esta interacción.

## Solución

Implementar un sistema de chat integrado donde:

1. **El cliente** puede iniciar una consulta desde el detalle de cualquier
   vehículo disponible, escribiendo un mensaje libre
2. **El vendedor** tiene una bandeja de entrada con tarjetas que muestran un
   preview de las consultas pendientes, y puede tomarlas para responder
3. **Ambos** pueden intercambiar mensajes simples dentro de la consulta
4. **Las notificaciones** de mensajes nuevos se muestran con un punto rojo.
   El backend calcula `mensajesNuevos` por consulta y expone un endpoint
   liviano `GET /api/notificaciones/contador`; el navbar consulta el contador
   con polling cada 3 segundos y las listas se actualizan al marcar mensajes
   como leídos

## Alcance

### Incluido

- Modelo de Consulta (conversación) con estados: pendiente, en conversación, cerrada
- Modelo de Mensaje con remitente, contenido y timestamp
- API REST para crear consultas, listar por cliente/vendedor, enviar mensajes
- Página de detalle del vehículo con botón "Consultar"
- Bandeja de entrada del vendedor con tarjetas y preview
- Vista tipo chat para el cliente ("Mis Consultas")
- Sistema de notificaciones con punto rojo: contador liviano para el navbar y campo `mensajesNuevos` por consulta en las listas
- Eliminación de consultas cerradas por parte del vendedor

### No incluido

- Chatbot con LangChain (CU-10 separado)
- Archivos adjuntos en mensajes
- Notificaciones push o WebSocket
- Historial de mensajes eliminados

## Decisiones clave

| Decisión | Elección | Razón |
|----------|----------|-------|
| Asignación de vendedor | Manual (el vendedor toma la consulta) | Flexibilidad para que el vendedor elija qué atender |
| Tipo de mensajes | Simples (solo texto) | Complejidad acorde al alcance |
| Notificaciones | Contador liviano en backend + polling corto en navbar + marcado como leído al abrir el chat | El conteo lo calcula el backend (fiabilidad) y el frontend solo muestra/oculta el punto |
| API | Endpoints separados para consultas y mensajes | Más RESTful y mantenible |
| Consultas cerradas | Se mantienen, vendedor puede eliminarlas | Preserva historial con opción de limpieza |

## Impacto

- **Backend**: Nuevos modelos, repositories, services y handlers para consultas y mensajes
- **Frontend**: Nuevas páginas (bandeja vendedor, chat cliente), modificaciones al detalle del vehículo
- **Base de datos**: Dos tablas nuevas (consultas, mensajes)
- **Autenticación**: Se reutiliza el sistema JWT existente
