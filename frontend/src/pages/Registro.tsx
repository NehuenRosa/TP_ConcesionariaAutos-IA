import { useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { useAuth } from '../hooks/useAuth'
import { ErrorApi } from '../services/api'
import { BotonGoogle } from '../components/BotonGoogle'
import { Boton } from '../components/ui/Boton'
import { CampoTexto } from '../components/ui/Campo'
import { MensajeError } from '../components/ui/MensajeError'

export function Registro() {
  const navigate = useNavigate()
  const { registrar } = useAuth()
  const [nombre, setNombre] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function enviar(e: React.FormEvent) {
    e.preventDefault()
    setEnviando(true)
    setError(null)

    try {
      await registrar({ nombre, email, password })
      navigate('/')
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo crear la cuenta.')
      setEnviando(false)
    }
  }

  return (
    <div className="mx-auto max-w-md px-4 py-20 sm:px-6">
      <div className="text-center">
        <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">
          Unite a Aurum Motors
        </p>
        <h1 className="font-display text-3xl font-bold tracking-tight text-plata-100">Crear cuenta</h1>
        <p className="mt-2 text-sm text-plata-400">Consultá, reservá y pedí test drives.</p>
      </div>

      <div className="mt-8 space-y-4">
        {error && <MensajeError>{error}</MensajeError>}

        <BotonGoogle alCompletar={() => navigate('/')} />

        <div className="flex items-center gap-3">
          <span className="h-px flex-1 bg-white/10" />
          <span className="text-xs uppercase tracking-widest text-plata-500">o registrate con tu email</span>
          <span className="h-px flex-1 bg-white/10" />
        </div>

        <form
          onSubmit={enviar}
          className="space-y-5 rounded-2xl border border-white/8 bg-carbono-850/60 p-6 shadow-luz backdrop-blur-sm"
        >
          <CampoTexto
            id="nombre"
            etiqueta="Nombre"
            type="text"
            required
            autoComplete="name"
            value={nombre}
            onChange={(e) => setNombre(e.target.value)}
          />
          <CampoTexto
            id="email"
            etiqueta="Email"
            type="email"
            required
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <div>
            <CampoTexto
              id="password"
              etiqueta="Contraseña"
              type="password"
              required
              minLength={8}
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <p className="mt-1 text-xs text-plata-500">Debe tener al menos 8 caracteres.</p>
          </div>

          <Boton type="submit" tamano="lg" className="w-full" disabled={enviando}>
            {enviando ? 'Creando cuenta…' : 'Registrarme'}
          </Boton>
        </form>

        <p className="text-center text-sm text-plata-400">
          ¿Ya tenés cuenta?{' '}
          <Link
            to="/login"
            className="font-semibold text-plata-100 underline-offset-4 hover:underline"
          >
            Iniciá sesión
          </Link>
        </p>
      </div>
    </div>
  )
}
