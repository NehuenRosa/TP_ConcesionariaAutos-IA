import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router'
import { useAuth } from '../hooks/useAuth'
import { ErrorApi } from '../services/api'
import { Boton } from '../components/ui/Boton'
import { CampoTexto } from '../components/ui/Campo'
import { MensajeError } from '../components/ui/MensajeError'

export function InicioSesion() {
  const navigate = useNavigate()
  const ubicacion = useLocation()
  const { iniciarSesion } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function enviar(e: React.FormEvent) {
    e.preventDefault()
    setEnviando(true)
    setError(null)

    try {
      await iniciarSesion({ email, password })
      const desde = (ubicacion.state as { desde?: { pathname?: string } } | null)?.desde
      navigate(desde?.pathname ?? '/')
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo iniciar sesión.')
      setEnviando(false)
    }
  }

  return (
    <div className="mx-auto max-w-md px-4 py-20 sm:px-6">
      <div className="text-center">
        <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">
          Bienvenido de nuevo
        </p>
        <h1 className="font-display text-3xl font-bold tracking-tight text-plata-100">Iniciar sesión</h1>
        <p className="mt-2 text-sm text-plata-400">Ingresá tus credenciales para continuar.</p>
      </div>

      <div className="mt-8 space-y-4">
        {error && <MensajeError>{error}</MensajeError>}

        <form
          onSubmit={enviar}
          className="space-y-5 rounded-2xl border border-white/8 bg-carbono-850/60 p-6 shadow-luz backdrop-blur-sm"
        >
          <CampoTexto
            id="email"
            etiqueta="Email"
            type="email"
            required
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <CampoTexto
            id="password"
            etiqueta="Contraseña"
            type="password"
            required
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          <Boton type="submit" tamano="lg" className="w-full" disabled={enviando}>
            {enviando ? 'Ingresando…' : 'Ingresar'}
          </Boton>
        </form>

        <p className="text-center text-sm text-plata-400">
          ¿No tenés cuenta?{' '}
          <Link
            to="/registro"
            className="font-semibold text-plata-100 underline-offset-4 hover:underline"
          >
            Registrate
          </Link>
        </p>
      </div>
    </div>
  )
}
