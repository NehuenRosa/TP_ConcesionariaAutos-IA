# Tasks: cu08-panel-transferencia-mis-reservas

## 1. Frontend

- [x] 1.1 Extraer `DatosSenia` a `components/reserva/DatosSenia.tsx` y usarlo
      desde `FormularioReserva`
- [x] 1.2 `MisReservas.tsx`: panel de datos de transferencia en reservas
      activas (caché por vehículo, fallback con monto/vencimiento propios)

## 2. Entorno

- [x] 2.1 `.env` con valores demo de `CBU_CONCESIONARIA`/`ALIAS_CONCESIONARIA`
      y `.env.example` documentado; recrear backend en Docker

## 3. Verificación y documentación

- [x] 3.1 `npm run build` + tests frontend en verde
- [x] 3.2 Prueba E2E manual: crear reserva → ver panel con CBU demo en
      FormularioReserva y en Mis Reservas; cancelar → desaparece
- [ ] 3.3 Actualizar roadmap y AGENTS.md si corresponde
