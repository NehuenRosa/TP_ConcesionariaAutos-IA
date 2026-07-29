# Concesionaria de Autos - TP

Sistema de gestión de concesionaria de autos con backend Go + PostgreSQL y frontend React + TypeScript.

## Requisitos

- Docker y Docker Compose

## Levantar el proyecto

```bash
docker compose up -d
```

Esto inicia 3 servicios:
- **PostgreSQL** en `localhost:5433` (user: `postgres`, password: `postgres`, db: `concesionaria`)
- **Backend** (Go/Gin) en `http://localhost:8080`
- **Frontend** (React/Vite) en `http://localhost:5173` (dev) / `http://localhost:5174` (docker)

## Credenciales de prueba

| Rol          | Email                      | Contraseña   |
|--------------|----------------------------|--------------|
| administrador| admin@concesionaria.com    | admin123     |
| vendedor     | vendedor@concesionaria.com | vendedor123  |
| cliente      | cliente@test.com           | cliente123   |

## Conectarse a la base de datos

Desde DBeaver, pgAdmin o cualquier cliente PostgreSQL:
- Host: `localhost`
- Puerto: `5433`
- Base de datos: `concesionaria`
- Usuario: `postgres`
- Contraseña: `postgres`

## Comandos útiles

```bash
# Ver logs
docker compose logs -f

# Reconstruir un servicio
docker compose up -d --build backend

# Detener todo
docker compose down

# Acceder a la base con psql
docker compose exec db psql -U postgres -d concesionaria
```
