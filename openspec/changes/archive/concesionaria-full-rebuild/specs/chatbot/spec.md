# Chatbot IA

## Requirements

### Requirement: Estado del chatbot
- **WHEN** cualquier usuario envía GET /api/chatbot/status
- **THEN** el sistema retorna 200 con "enabled": true/false según disponibilidad de API key

### Requirement: Consulta al chatbot
- **WHEN** cualquier usuario envía POST /api/chatbot/ask con question
- **THEN** el sistema retorna 200 con una respuesta generada por IA basada en el inventario actual de vehículos disponibles

- **WHEN** un usuario envía POST /api/chatbot/ask sin question
- **THEN** el sistema retorna 400

### Requirement: Widget de chatbot
- **WHEN** un usuario navega por el frontend
- **THEN** el sistema muestra un widget flotante en la esquina inferior derecha que permite abrir/cerrar y conversar con el asistente
