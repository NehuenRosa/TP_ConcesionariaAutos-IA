import { createContext, useContext, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { api, eliminarToken, guardarToken, obtenerToken } from '../services/api'
import type { DatosLogin, DatosRegistro, Usuario } from '../types/usuario'

interface ContextoSesion {
  usuario: Usuario | null
  cargando: boolean
  esAdministrador: boolean
  iniciarSesion: (datos: DatosLogin) => Promise<void>
  registrar: (datos: DatosRegistro) => Promise<void>
  cerrarSesion: () => void
}

const ContextoAutenticacion = createContext<ContextoSesion | undefined>(undefined)

export function ProveedorAutenticacion({ children }: { children: ReactNode }) {
  const [usuario, setUsuario] = useState<Usuario | null>(null)
  const [cargando, setCargando] = useState(true)

  useEffect(() => {
    let cancelado = false
    if (!obtenerToken()) {
      setCargando(false)
      return
    }

    api
      .obtenerPerfil()
      .then((perfil) => {
        if (!cancelado) setUsuario(perfil)
      })
      .catch(() => {
        eliminarToken()
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [])

  async function iniciarSesion(datos: DatosLogin) {
    const respuesta = await api.iniciarSesion(datos)
    guardarToken(respuesta.token)
    setUsuario(respuesta.usuario)
  }

  async function registrar(datos: DatosRegistro) {
    await api.registrar(datos)
    await iniciarSesion({ email: datos.email, password: datos.password })
  }

  function cerrarSesion() {
    eliminarToken()
    setUsuario(null)
  }

  return (
    <ContextoAutenticacion.Provider
      value={{
        usuario,
        cargando,
        esAdministrador: usuario?.rol === 'administrador',
        iniciarSesion,
        registrar,
        cerrarSesion,
      }}
    >
      {children}
    </ContextoAutenticacion.Provider>
  )
}

export function useAuth(): ContextoSesion {
  const contexto = useContext(ContextoAutenticacion)
  if (!contexto) {
    throw new Error('useAuth debe usarse dentro de ProveedorAutenticacion')
  }
  return contexto
}
