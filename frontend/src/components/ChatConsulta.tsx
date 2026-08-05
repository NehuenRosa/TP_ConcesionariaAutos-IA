import { useEffect, useState, useRef } from 'react'
import { api, ErrorApi } from '../services/api'
import { useAuth } from '../hooks/useAuth'
import type { Mensaje, EstadoConsulta } from '../types/consulta'
import { Boton } from './ui/Boton'

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

  useEffect(() => {
    let cancelado = false

    const cargarMensajes = async () => {
      try {
        const datos = await api.obtenerMensajes(consultaId)
        if (!cancelado) {
          setMensajes(datos)
          if (datos.length > 0) {
            ultimoTimestampRef.current = datos[datos.length - 1].createdAt
          }
          setError(null)
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

  useEffect(() => {
    if (estado === 'cerrada') return

    const intervalo = setInterval(async () => {
      if (!ultimoTimestampRef.current) return

      try {
        const nuevos = await api.obtenerMensajesNuevos(consultaId, ultimoTimestampRef.current)
        if (nuevos.length > 0) {
          const mensajesDeOtros = nuevos.filter((m) => m.remitente.id !== usuario?.id)
          if (mensajesDeOtros.length > 0) {
            api.marcarComoLeidos(consultaId).then(() => {
              window.dispatchEvent(new Event('mensajes-leidos'))
            }).catch(() => {})
          }

          setMensajes((prev) => {
            const idsExistentes = new Set(prev.map((m) => m.id))
            const mensajesNuevos = nuevos.filter((m) => !idsExistentes.has(m.id))
            if (mensajesNuevos.length === 0) return prev
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
        <p className="text-sm text-plata-400">Cargando mensajes…</p>
      </div>
    )
  }

  const cerrada = estado === 'cerrada'

  return (
    <div className="flex h-full flex-col">
      <div className="border-b border-white/8 px-5 py-3.5">
        <p className="text-xs text-plata-400">
          {cerrada ? 'Consulta cerrada' : 'Escribí tu mensaje'}
        </p>
      </div>

      <div ref={mensajesRef} className="flex-1 overflow-y-auto px-5 py-4">
        {error && (
          <div className="mb-4 rounded-xl border border-red-500/30 bg-red-500/10 p-3 text-sm text-red-300">
            {error}
          </div>
        )}

        {mensajes.length === 0 ? (
          <div className="flex h-full items-center justify-center">
            <p className="text-sm text-plata-500">No hay mensajes aún</p>
          </div>
        ) : (
          <div className="space-y-4">
            {mensajes.map((mensaje) => {
              const esMio = mensaje.remitente.id === usuario?.id
              return (
                <div key={mensaje.id} className={`flex ${esMio ? 'justify-end' : 'justify-start'}`}>
                  <div
                    className={`max-w-[75%] rounded-2xl px-4 py-2.5 ${
                      esMio
                        ? 'rounded-br-sm bg-plata-100 text-carbono-900'
                        : 'rounded-bl-sm border border-white/8 bg-carbono-800 text-plata-100'
                    }`}
                  >
                    {!esMio && (
                      <p className="mb-1 text-xs font-semibold text-plata-400">
                        {mensaje.remitente.nombre}
                      </p>
                    )}
                    <p className="whitespace-pre-wrap text-sm">{mensaje.contenido}</p>
                    <p className={`mt-1 text-[11px] ${esMio ? 'text-carbono-500' : 'text-plata-500'}`}>
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

      {!cerrada && (
        <div className="border-t border-white/8 p-4">
          <div className="flex items-end gap-2">
            <textarea
              value={nuevoMensaje}
              onChange={(e) => setNuevoMensaje(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Escribí tu mensaje…"
              className="campo resize-none"
              rows={2}
              disabled={enviando}
            />
            <Boton onClick={handleEnviar} disabled={enviando || !nuevoMensaje.trim()}>
              {enviando ? '…' : 'Enviar'}
            </Boton>
          </div>
        </div>
      )}
    </div>
  )
}
