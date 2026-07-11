---
name: generate-component
description: Crea un componente React con TypeScript y TailwindCSS siguiendo los patrones del proyecto
compatibility: opencode
---

## Paso a paso para crear un componente React

1. **Ubicación**:
   - Páginas en `frontend/src/pages/` (organizadas por rol: `admin/`, `seller/`)
   - Componentes compartidos en `frontend/src/components/`
   - Naming: `NombreComponente.tsx`

2. **Estructura del componente**:
   ```tsx
   export function NombreComponente({ prop1, prop2 }: Props) {
     return (
       <div className="...">
         {/* contenido */}
       </div>
     )
   }
   ```

3. **Convenciones**:
   - Tipos TypeScript en `frontend/src/types/index.ts`
   - Llamadas HTTP en `frontend/src/services/<recurso>Service.ts`
   - Estados con `useState`, efectos con `useEffect`
   - Navegación con `useNavigate` de React Router
   - Auth con `useAuth()` hook
   - Estilos con TailwindCSS (sin CSS adicional)
   - Exportación nombrada (no default)
