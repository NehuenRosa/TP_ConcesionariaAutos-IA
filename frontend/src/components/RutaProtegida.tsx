import type { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router'
import { useAuth } from '../hooks/useAuth'
import type { Rol } from '../types/usuario'

interface Props {
  rol?: Rol
  children: ReactNode
}

export function RutaProtegida({ rol, children }: Props) {
  const { usuario, cargando } = useAuth()
  const ubicacion = useLocation()

  if (cargando) {
    return <p className="text-gray-700">Cargando sesión…</p>
  }

  if (!usuario) {
    return <Navigate to="/login" replace state={{ desde: ubicacion }} />
  }

  if (rol && usuario.rol !== rol && usuario.rol !== 'administrador') {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}
