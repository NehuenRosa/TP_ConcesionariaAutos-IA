import { useEffect, useState, useRef } from 'react'
import { api, ErrorApi } from '../services/api'
import { useAuth } from '../hooks/useAuth'
import type { Mensaje, EstadoConsulta } from '../types/consulta'

interface ChatConsultaProps {
  consultaId: number
  estado: EstadoConsulta
  onMensajeEnviado?: () => void
}

export function ChatConsulta({ consultaId, estado, onMensajeEnviado }: ChatConsultaProps) {
  const { usuario } = useAuth()
  const [mensajes, setMensajes] = useState<Mensaje[]>([])
  const [nuevoMensaje, setNuevoMensaje] = useState('')
  const [cargando, setCargando] = useState(true)
  const [enviando, setEnviando] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const mensajesRef = useRef<HTMLDivElement>(null)
  const ultimoTimestampRef = useRef<string>('')

  const scrollToBottom = () => {
    if (mensajesRef.current) {
      mensajesRef.current.scrollTop = mensajesRef.current.scrollHeight
    }
  }

  // Cargar mensajes iniciales y marcar como leídos
  useEffect(() => {
    let cancelado = false

    const cargarMensajes = async () => {
      try {
        const datos = await api.obtenerMensajes(consultaId)
        if (!cancelado) {
          setMensajes(datos)
          // Guardar timestamp del último mensaje
          if (datos.length > 0) {
            ultimoTimestampRef.current = datos[datos.length - 1].createdAt
          }
          setError(null)
          // Marcar mensajes como leídos al abrir el chat
          api.marcarComoLeidos(consultaId).then(() => {
            window.dispatchEvent(new Event('mensajes-leidos'))
          }).catch(() => {})
        }
      } catch (e: unknown) {
        if (!cancelado) {
          setError(e instanceof ErrorApi ? e.message : 'Error al cargar mensajes')
        }
      } finally {
        if (!cancelado) setCargando(false)
      }
    }

    cargarMensajes()

    return () => {
      cancelado = true
    }
  }, [consultaId])

  useEffect(() => {
    scrollToBottom()
  }, [mensajes])

  // Polling cada 5 segundos para mensajes nuevos (sin dependencia de mensajes)
  useEffect(() => {
    if (estado === 'cerrada') return

    const intervalo = setInterval(async () => {
      // Usar el ref en lugar del estado para evitar dependencia
      if (!ultimoTimestampRef.current) return
      
      try {
        const nuevos = await api.obtenerMensajesNuevos(consultaId, ultimoTimestampRef.current)
        if (nuevos.length > 0) {
          // Filtrar mensajes de otros (no míos) para marcar como leídos
          const mensajesDeOtros = nuevos.filter(m => m.remitente.id !== usuario?.id)
          if (mensajesDeOtros.length > 0) {
            api.marcarComoLeidos(consultaId).then(() => {
              window.dispatchEvent(new Event('mensajes-leidos'))
            }).catch(() => {})
          }

          setMensajes((prev) => {
            // Filtrar mensajes duplicados por ID
            const idsExistentes = new Set(prev.map(m => m.id))
            const mensajesNuevos = nuevos.filter(m => !idsExistentes.has(m.id))
            if (mensajesNuevos.length === 0) return prev
            // Actualizar timestamp del último mensaje
            ultimoTimestampRef.current = mensajesNuevos[mensajesNuevos.length - 1].createdAt
            return [...prev, ...mensajesNuevos]
          })
        }
      } catch {
        // Ignorar errores de polling
      }
    }, 5000)

    return () => clearInterval(intervalo)
  }, [consultaId, estado])

  const handleEnviar = async () => {
    if (!nuevoMensaje.trim() || enviando) return

    setEnviando(true)
    setError(null)

    try {
      const mensaje = await api.enviarMensaje(consultaId, nuevoMensaje.trim())
      setMensajes((prev) => [...prev, mensaje])
      // Actualizar timestamp con el mensaje enviado
      ultimoTimestampRef.current = mensaje.createdAt
      setNuevoMensaje('')
      onMensajeEnviado?.()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo enviar el mensaje')
    } finally {
      setEnviando(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleEnviar()
    }
  }

  if (cargando) {
    return (
      <div className="flex h-full items-center justify-center">
        <p className="text-gray-500">Cargando mensajes…</p>
      </div>
    )
  }

  const cerrada = estado === 'cerrada'

  return (
    <div className="flex h-full flex-col">
      {/* Encabezado */}
      <div className="border-b border-gray-200 p-4">
        <p className="text-sm text-gray-500">
          {cerrada ? 'Consulta cerrada' : 'Escribí tu mensaje'}
        </p>
      </div>

      {/* Mensajes */}
      <div ref={mensajesRef} className="flex-1 overflow-y-auto p-4">
        {error && (
          <div className="mb-4 rounded-lg bg-red-50 p-3">
            <p className="text-sm text-red-800">{error}</p>
          </div>
        )}

        {mensajes.length === 0 ? (
          <div className="flex h-full items-center justify-center">
            <p className="text-gray-500">No hay mensajes aún</p>
          </div>
        ) : (
          <div className="space-y-4">
            {mensajes.map((mensaje) => {
              const esMio = mensaje.remitente.id === usuario?.id
              return (
                <div
                  key={mensaje.id}
                  className={`flex ${esMio ? 'justify-end' : 'justify-start'}`}
                >
                  <div
                    className={`max-w-[70%] rounded-lg px-4 py-2 ${
                      esMio
                        ? 'bg-gray-900 text-white'
                        : 'bg-gray-100 text-gray-900'
                    }`}
                  >
                    {!esMio && (
                      <p className="mb-1 text-xs font-medium text-gray-500">
                        {mensaje.remitente.nombre}
                      </p>
                    )}
                    <p className="whitespace-pre-wrap">{mensaje.contenido}</p>
                    <p
                      className={`mt-1 text-xs ${
                        esMio ? 'text-gray-300' : 'text-gray-400'
                      }`}
                    >
                      {new Date(mensaje.createdAt).toLocaleTimeString('es-AR', {
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </p>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>

      {/* Input */}
      {!cerrada && (
        <div className="border-t border-gray-200 p-4">
          <div className="flex gap-2">
            <textarea
              value={nuevoMensaje}
              onChange={(e) => setNuevoMensaje(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Escribí tu mensaje..."
              className="flex-1 resize-none rounded-lg border border-gray-300 p-3 text-gray-900 focus:border-gray-500 focus:outline-none"
              rows={2}
              disabled={enviando}
            />
            <button
              type="button"
              onClick={handleEnviar}
              disabled={enviando || !nuevoMensaje.trim()}
              className="self-end rounded-lg bg-gray-900 px-4 py-3 text-white hover:bg-gray-700 disabled:opacity-50"
            >
              {enviando ? '...' : 'Enviar'}
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
