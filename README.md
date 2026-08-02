# Concesionaria de Autos - Sistema de Gestión

Sistema web para administrar el stock de una concesionaria de autos, con catálogo
público, consultas y reservas de clientes, turnos de test drive, panel de
administración y un asistente conversacional (chatbot).

> **Nota:** este README se completa en el Entregable 6 con la guía definitiva de
> levante (Docker y sin Docker). Mientras tanto, mirá `docs/roadmap.md` y `AGENTS.md`.

## Stack

| Capa | Tecnología |
|------|------------|
| Backend | Go + Gin + GORM + JWT |
| Frontend | React + Vite + TypeScript + React Router + TailwindCSS |
| Base de datos | PostgreSQL |
| Chatbot | LangChain (langchaingo) |
| Infra | Docker + Docker Compose |

## Estructura del repositorio

```
├── backend/     # API REST en Go
├── frontend/    # Aplicación web en React
├── openspec/    # Especificaciones y cambios planificados
├── docs/        # Documentación (roadmap, decisiones)
├── docker-compose.yml
└── AGENTS.md    # Guía de contexto para agentes de IA
```
