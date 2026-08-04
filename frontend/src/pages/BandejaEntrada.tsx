import { useEffect, useState, useCallback } from 'react'
import { useNavigate } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { ConsultaResumen, EstadoConsulta } from '../types/consulta'

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

export function BandejaEntrada() {
  const navigate = useNavigate()
  const [consultas, setConsultas] = useState<ConsultaResumen[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const cargarBandeja = useCallback(async () => {
    try {
      const datos = await api.listarBandeja()
      setConsultas(datos)
      setError(null)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'Error al cargar la bandeja')
    } finally {
      setCargando(false)
    }
  }, [])

  useEffect(() => {
    cargarBandeja()
  }, [cargarBandeja])

  // Recargar al volver del chat (se marcaron leídos) y cada 5 segundos para
  // detectar mensajes nuevos.
  useEffect(() => {
    const manejarLeidos = () => cargarBandeja()
    window.addEventListener('mensajes-leidos', manejarLeidos)

    const intervalo = setInterval(cargarBandeja, 5000)

    return () => {
      window.removeEventListener('mensajes-leidos', manejarLeidos)
      clearInterval(intervalo)
    }
  }, [cargarBandeja])

  const handleTomar = async (id: number) => {
    try {
      await api.tomarConsulta(id)
      await cargarBandeja()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo tomar la consulta')
    }
  }

  const handleCerrar = async (id: number) => {
    try {
      await api.cerrarConsulta(id)
      await cargarBandeja()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo cerrar la consulta')
    }
  }

  const handleEliminar = async (id: number) => {
    try {
      await api.eliminarConsulta(id)
      await cargarBandeja()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo eliminar la consulta')
    }
  }

  if (cargando) {
    return <p className="text-gray-700">Cargando bandeja de entrada…</p>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Bandeja de entrada</h1>

      {error && (
        <div className="rounded-lg bg-red-50 p-4">
          <p className="text-red-800">{error}</p>
        </div>
      )}

      {consultas.length === 0 ? (
        <div className="rounded-lg border border-gray-200 p-8 text-center">
          <p className="text-gray-500">No tenés consultas asignadas</p>
        </div>
      ) : (
        <div className="grid gap-4">
          {consultas.map((consulta) => (
            <div
              key={consulta.id}
              className={`relative rounded-lg border p-4 transition-all hover:shadow-md ${
                consulta.mensajesNuevos > 0
                  ? 'border-blue-300 bg-blue-50'
                  : 'border-gray-200 bg-white'
              }`}
            >
              {consulta.mensajesNuevos > 0 && (
                <span className="absolute right-4 top-4 h-3 w-3 rounded-full bg-red-500" />
              )}

              <div className="flex items-start justify-between">
                <div
                  className="flex-1 cursor-pointer"
                  onClick={() => navigate(`/vendedor/bandeja/${consulta.id}`)}
                >
                  <div className="flex items-center gap-3">
                    <h3 className="font-semibold text-gray-900">
                      {consulta.vehiculo.marca} {consulta.vehiculo.modelo} {consulta.vehiculo.anio}
                    </h3>
                    <span className={`rounded-full px-2 py-1 text-xs font-medium ${estadoColor(consulta.estado)}`}>
                      {estadoATexto(consulta.estado)}
                    </span>
                  </div>

                  <p className="mt-1 text-sm text-gray-600">
                    Cliente: {consulta.cliente.nombre}
                  </p>

                  {consulta.ultimoMensaje && (
                    <p className="mt-2 truncate text-sm text-gray-500">
                      {consulta.ultimoMensaje.contenido}
                    </p>
                  )}

                  <p className="mt-1 text-xs text-gray-400">
                    {formatearFecha(consulta.updatedAt)}
                  </p>
                </div>

                <div className="ml-4 flex gap-2">
                  {consulta.estado === 'pendiente' && (
                    <button
                      type="button"
                      onClick={() => handleTomar(consulta.id)}
                      className="rounded-md bg-gray-900 px-3 py-1 text-sm text-white hover:bg-gray-700"
                    >
                      Tomar
                    </button>
                  )}
                  {consulta.estado === 'en_conversacion' && (
                    <button
                      type="button"
                      onClick={() => handleCerrar(consulta.id)}
                      className="rounded-md border border-gray-300 px-3 py-1 text-sm text-gray-700 hover:bg-gray-50"
                    >
                      Cerrar
                    </button>
                  )}
                  {consulta.estado === 'cerrada' && (
                    <button
                      type="button"
                      onClick={() => handleEliminar(consulta.id)}
                      className="rounded-md border border-red-200 px-3 py-1 text-sm text-red-600 hover:bg-red-50"
                    >
                      Eliminar
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
