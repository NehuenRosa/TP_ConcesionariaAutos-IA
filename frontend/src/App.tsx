import { ProveedorAutenticacion } from './hooks/useAuth'
import { Rutas } from './routes/Rutas'

export function App() {
  return (
    <ProveedorAutenticacion>
      <Rutas />
    </ProveedorAutenticacion>
  )
}
