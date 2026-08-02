import { useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { useAuth } from '../hooks/useAuth'
import { ErrorApi } from '../services/api'

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

  const campo =
    'w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:border-gray-900 focus:outline-none'
  const etiqueta = 'block text-sm font-medium text-gray-700'

  return (
    <div className="mx-auto max-w-md space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Crear cuenta</h1>
        <p className="mt-1 text-gray-700">
          Registrate para consultar, reservar y pedir test drives.
        </p>
      </div>

      {error && <p className="text-red-600">{error}</p>}

      <form onSubmit={enviar} className="space-y-4 rounded-lg border border-gray-200 bg-white p-6">
        <div>
          <label htmlFor="nombre" className={etiqueta}>
            Nombre
          </label>
          <input
            id="nombre"
            type="text"
            required
            autoComplete="name"
            value={nombre}
            onChange={(e) => setNombre(e.target.value)}
            className={campo}
          />
        </div>
        <div>
          <label htmlFor="email" className={etiqueta}>
            Email
          </label>
          <input
            id="email"
            type="email"
            required
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className={campo}
          />
        </div>
        <div>
          <label htmlFor="password" className={etiqueta}>
            Contraseña
          </label>
          <input
            id="password"
            type="password"
            required
            minLength={8}
            autoComplete="new-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className={campo}
          />
          <p className="mt-1 text-xs text-gray-500">Debe tener al menos 8 caracteres.</p>
        </div>

        <button
          type="submit"
          disabled={enviando}
          className="w-full rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {enviando ? 'Creando cuenta…' : 'Registrarme'}
        </button>
      </form>

      <p className="text-center text-sm text-gray-700">
        ¿Ya tenés cuenta?{' '}
        <Link to="/login" className="text-gray-900 underline hover:text-gray-700">
          Iniciá sesión
        </Link>
      </p>
    </div>
  )
}
