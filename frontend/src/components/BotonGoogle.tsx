import { useEffect, useRef, useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { api, ErrorApi } from '../services/api'
import { MensajeError } from './ui/MensajeError'

// Tipos mínimos de Google Identity Services (script oficial gsi/client).
interface RespuestaCredencialGoogle {
  credential: string
}

interface CuentaGoogleId {
  initialize(config: { client_id: string; callback: (respuesta: RespuestaCredencialGoogle) => void }): void
  renderButton(
    contenedor: HTMLElement,
    opciones: Record<string, string | number | boolean>,
  ): void
}

declare global {
  interface Window {
    google?: { accounts: { id: CuentaGoogleId } }
  }
}

const URL_SCRIPT_GIS = 'https://accounts.google.com/gsi/client'

// cargarScriptGIS inyecta el script oficial una única vez y resuelve cuando
// window.google está disponible.
let cargaScriptGIS: Promise<void> | null = null
function cargarScriptGIS(): Promise<void> {
  if (window.google?.accounts?.id) return Promise.resolve()
  cargaScriptGIS ??= new Promise<void>((resolver, rechazar) => {
    const script = document.createElement('script')
    script.src = URL_SCRIPT_GIS
    script.async = true
    script.defer = true
    script.onload = () => resolver()
    script.onerror = () => {
      cargaScriptGIS = null
      rechazar(new Error('No se pudo cargar Google Identity Services'))
    }
    document.head.appendChild(script)
  })
  return cargaScriptGIS
}

interface BotonGoogleProps {
  /** Se invoca cuando el ingreso con Google terminó con éxito. */
  alCompletar: () => void
}

// BotonGoogle muestra el botón oficial "Continuar con Google" si el backend
// tiene el acceso habilitado; en otro caso no renderiza nada.
export function BotonGoogle({ alCompletar }: BotonGoogleProps) {
  const { iniciarSesionConGoogle } = useAuth()
  const [habilitado, setHabilitado] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const contenedorRef = useRef<HTMLDivElement>(null)
  const alCompletarRef = useRef(alCompletar)

  useEffect(() => {
    alCompletarRef.current = alCompletar
  }, [alCompletar])

  useEffect(() => {
    let cancelado = false
    let clienteID: string | undefined

    function renderizar() {
      const gsi = window.google?.accounts?.id
      if (!gsi || !contenedorRef.current || !clienteID) return
      gsi.initialize({
        client_id: clienteID,
        callback: (respuesta) => {
          void manejarCredencial(respuesta.credential)
        },
      })
      gsi.renderButton(contenedorRef.current, {
        theme: 'filled_black',
        size: 'large',
        text: 'continue_with',
        shape: 'pill',
        locale: 'es',
        width: 320,
      })
    }

    async function manejarCredencial(credencial: string) {
      setError(null)
      try {
        await iniciarSesionConGoogle(credencial)
        alCompletarRef.current()
      } catch (e: unknown) {
        setError(e instanceof ErrorApi ? e.message : 'No se pudo iniciar sesión con Google.')
      }
    }

    api
      .obtenerProveedoresAuth()
      .then((proveedores) => {
        if (cancelado || !proveedores.google || !proveedores.client_id) return
        clienteID = proveedores.client_id
        setHabilitado(true)
        return cargarScriptGIS()
          .then(renderizar)
          .catch(() => {
            if (!cancelado) setHabilitado(false)
          })
      })
      .catch(() => {
        // Si el endpoint falla, el botón simplemente no aparece.
      })

    return () => {
      cancelado = true
    }
  }, [iniciarSesionConGoogle])

  if (!habilitado && !error) return null

  return (
    <div className="space-y-3">
      {error && <MensajeError titulo="No se pudo continuar con Google">{error}</MensajeError>}
      <div ref={contenedorRef} className={habilitado ? 'flex justify-center' : 'hidden'} />
    </div>
  )
}
