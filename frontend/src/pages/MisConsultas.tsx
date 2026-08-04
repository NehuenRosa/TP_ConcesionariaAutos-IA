import { useEffect, useState, useCallback } from 'react'
import { Link, useParams, useNavigate } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { ConsultaResumen, EstadoConsulta } from '../types/consulta'
import { ChatConsulta } from '../components/ChatConsulta'

function formatearFecha(fecha: string): string {
  return new Date(fecha).toLocaleDateString('es-AR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function estadoATexto(estado: EstadoConsulta): string {
  switch (estado) {
    case 'pendiente':
      return 'Pendiente'
    case 'en_conversacion':
      return 'En conversación'
    case 'cerrada':
      return 'Cerrada'
    default:
      return estado
  }
}

function estadoColor(estado: EstadoConsulta): string {
  switch (estado) {
    case 'pendiente':
      return 'bg-yellow-100 text-yellow-800'
    case 'en_conversacion':
      return 'bg-green-100 text-green-800'
    case 'cerrada':
      return 'bg-gray-100 text-gray-600'
    default:
      return 'bg-gray-100 text-gray-600'
  }
}

export function MisConsultas() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [consultas, setConsultas] = useState<ConsultaResumen[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const cargarConsultas = useCallback(async () => {
    try {
      const datos = await api.listarMisConsultas()
      setConsultas(datos)
      setError(null)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'Error al cargar las consultas')
    } finally {
      setCargando(false)
    }
  }, [])

  useEffect(() => {
    cargarConsultas()
  }, [cargarConsultas])

  // Recargar al volver del chat (se marcaron leídos) y cada 5 segundos para
  // detectar mensajes nuevos.
  useEffect(() => {
    const manejarLeidos = () => cargarConsultas()
    window.addEventListener('mensajes-leidos', manejarLeidos)

    const intervalo = setInterval(cargarConsultas, 5000)

    return () => {
      window.removeEventListener('mensajes-leidos', manejarLeidos)
      clearInterval(intervalo)
    }
  }, [cargarConsultas])

  const seleccionada = id ? consultas.find((c) => c.id === Number(id)) : null

  if (cargando) {
    return <p className="text-gray-700">Cargando consultas…</p>
  }

  return (
    <div className="flex h-[calc(100vh-200px)] gap-4">
      {/* Lista de consultas */}
      <div className="w-80 flex-shrink-0 overflow-y-auto rounded-lg border border-gray-200 bg-white">
        <div className="border-b border-gray-200 p-4">
          <h1 className="text-lg font-semibold text-gray-900">Mis Consultas</h1>
        </div>

        {error && (
          <div className="border-b border-gray-200 bg-red-50 p-3">
            <p className="text-sm text-red-800">{error}</p>
          </div>
        )}

        {consultas.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-gray-500">No tenés consultas</p>
            <Link
              to="/catalogo"
              className="mt-2 inline-block text-sm text-gray-900 hover:underline"
            >
              Ver catálogo
            </Link>
          </div>
        ) : (
          <div className="divide-y divide-gray-200">
            {consultas.map((consulta) => (
              <button
                key={consulta.id}
                type="button"
                onClick={() => navigate(`/mis-consultas/${consulta.id}`)}
                className={`relative w-full p-4 text-left transition-colors hover:bg-gray-50 ${
                  seleccionada?.id === consulta.id ? 'bg-gray-100' : ''
                }`}
              >
                {consulta.mensajesNuevos > 0 && (
                  <span className="absolute right-3 top-3 h-2.5 w-2.5 rounded-full bg-red-500" />
                )}

                <div className="flex items-center gap-2">
                  <h3 className="font-medium text-gray-900">
                    {consulta.vehiculo.marca} {consulta.vehiculo.modelo}
                  </h3>
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${estadoColor(consulta.estado)}`}>
                    {estadoATexto(consulta.estado)}
                  </span>
                </div>

                {consulta.vendedor && (
                  <p className="mt-1 text-xs text-gray-500">
                    Vendedor: {consulta.vendedor.nombre}
                  </p>
                )}

                {consulta.ultimoMensaje && (
                  <p className="mt-1 truncate text-sm text-gray-500">
                    {consulta.ultimoMensaje.contenido}
                  </p>
                )}

                <p className="mt-1 text-xs text-gray-400">
                  {formatearFecha(consulta.updatedAt)}
                </p>
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Área de chat */}
      <div className="flex-1 rounded-lg border border-gray-200 bg-white">
        {seleccionada ? (
          <ChatConsulta
            consultaId={seleccionada.id}
            estado={seleccionada.estado}
            onMensajeEnviado={cargarConsultas}
          />
        ) : (
          <div className="flex h-full items-center justify-center">
            <p className="text-gray-500">Seleccioná una consulta para ver la conversación</p>
          </div>
        )}
      </div>
    </div>
  )
}
