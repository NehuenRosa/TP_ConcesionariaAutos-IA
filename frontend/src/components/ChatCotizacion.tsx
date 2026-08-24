import { useEffect, useState, useRef } from 'react'
import { Link } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { Cotizacion, MensajeCotizacion } from '../types/cotizacion'
import { Boton } from './ui/Boton'
import { formatearPrecio } from '../utils/formato'

interface ChatCotizacionProps {
  cotizacionId: number
  onCambio: () => void
}

const INTERVALO_ACTUALIZACION_MS = 10000

export function ChatCotizacion({ cotizacionId, onCambio }: ChatCotizacionProps) {
  const [mensajes, setMensajes] = useState<MensajeCotizacion[]>([])
  const [estado, setEstado] = useState<'abierta' | 'cerrada'>('abierta')
  const [vehiculo, setVehiculo] = useState<Cotizacion['vehiculo'] | null>(null)
  const [vendedor, setVendedor] = useState<Cotizacion['vendedor']>()
  const [nuevoMensaje, setNuevoMensaje] = useState('')
  const [cargando, setCargando] = useState(true)
  const [enviando, setEnviando] = useState(false)
  const [cerrando, setCerrando] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const mensajesRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    let cancelado = false
    setCargando(true)

    api
      .obtenerCotizacion(cotizacionId)
      .then((datos) => {
        if (cancelado) return
        setMensajes(datos.mensajes)
        setEstado(datos.estado)
        setVehiculo(datos.vehiculo)
        setVendedor(datos.vendedor)
        setError(null)
        // Abrir el hilo lo marca como leído en el servidor: refresca el contador.
        window.dispatchEvent(new Event('mensajes-leidos'))
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'Error al cargar la cotización')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [cotizacionId])

  // La conversación se refresca sola para ver las respuestas del vendedor.
  useEffect(() => {
    const temporizador = setInterval(async () => {
      try {
        const datos = await api.obtenerCotizacion(cotizacionId)
        setMensajes(datos.mensajes)
        setEstado(datos.estado)
        setVehiculo(datos.vehiculo)
        setVendedor(datos.vendedor)
      } catch {
        // Se ignora: el próximo intento lo vuelve a traer.
      }
    }, INTERVALO_ACTUALIZACION_MS)
    return () => clearInterval(temporizador)
  }, [cotizacionId])

  useEffect(() => {
    if (mensajesRef.current) {
      mensajesRef.current.scrollTop = mensajesRef.current.scrollHeight
    }
  }, [mensajes])

  const handleEnviar = async () => {
    const contenido = nuevoMensaje.trim()
    if (!contenido || enviando) return

    setEnviando(true)
    setError(null)

    // Mensaje del cliente en vuelo
    setMensajes((prev) => [
      ...prev,
      { id: -Date.now(), remitente: 'cliente', contenido, createdAt: new Date().toISOString() },
    ])
    setNuevoMensaje('')

    try {
      const resultado = await api.enviarMensajeCotizacion(cotizacionId, contenido)
      // Si un vendedor tomó la conversación, la IA queda silenciada y no
      // responde: el mensaje queda esperando la respuesta del personal.
      if (resultado.respuesta !== '') {
        setMensajes((prev) => [
          ...prev,
          { id: -Date.now() - 1, remitente: 'ia', contenido: resultado.respuesta, createdAt: new Date().toISOString() },
        ])
      }
    } catch (e: unknown) {
      setMensajes((prev) => prev.slice(0, -1))
      setNuevoMensaje(contenido)
      setError(e instanceof ErrorApi ? e.message : 'No se pudo enviar el mensaje')
    } finally {
      setEnviando(false)
    }
  }

  const handleCerrar = async () => {
    if (cerrando) return
    setCerrando(true)
    setError(null)
    try {
      await api.cerrarCotizacion(cotizacionId)
      setEstado('cerrada')
      onCambio()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo cerrar la cotización')
    } finally {
      setCerrando(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      void handleEnviar()
    }
  }

  if (cargando) {
    return (
      <div className="flex h-full items-center justify-center">
        <p className="text-sm text-plata-400">Cargando conversación…</p>
      </div>
    )
  }

  const cerrada = estado === 'cerrada'

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between gap-3 border-b border-white/8 px-5 py-3.5">
        <div className="min-w-0">
          <p className="font-display text-sm font-semibold text-plata-100">
            {vehiculo ? `${vehiculo.marca} ${vehiculo.modelo}` : 'Cotización'}
          </p>
          {vehiculo && (
            <p className="text-xs text-plata-500">
              {vehiculo.anio} · {formatearPrecio(vehiculo.precio)}
            </p>
          )}
        </div>
        <div className="flex shrink-0 items-center gap-2">
          <span
            className={`rounded-full px-2.5 py-1 text-[11px] font-semibold tracking-wide uppercase ${
              cerrada ? 'bg-white/10 text-plata-400' : 'bg-emerald-400/15 text-emerald-300'
            }`}
          >
            {cerrada ? 'Cerrada' : 'Abierta'}
          </span>
          {!cerrada && (
            <Boton variante="secundario" tamano="sm" onClick={handleCerrar} disabled={cerrando}>
              {cerrando ? '…' : 'Cerrar'}
            </Boton>
          )}
        </div>
      </div>

      {vendedor && (
        <div className="border-b border-acento-400/20 bg-acento-500/10 px-5 py-2 text-xs text-acento-300">
          Estás siendo atendido por {vendedor.nombre}. El asistente quedó en pausa.
        </div>
      )}

      <div ref={mensajesRef} className="flex-1 overflow-y-auto px-5 py-4">
        {error && (
          <div className="mb-4 rounded-xl border border-red-500/30 bg-red-500/10 p-3 text-sm text-red-300">
            {error}
          </div>
        )}

        <div className="space-y-4">
          {mensajes.length === 0 && (
            <p className="text-sm text-plata-500">
              {cerrada ? 'Esta cotización está cerrada.' : 'Escribile a la IA para seguir cotizando este vehículo.'}
            </p>
          )}
          {mensajes.map((mensaje) => {
            const esCliente = mensaje.remitente === 'cliente'
            const etiqueta =
              mensaje.remitente === 'ia'
                ? 'Asistente IA'
                : mensaje.remitente === 'vendedor'
                  ? 'Asesor de ventas'
                  : null
            return (
              <div key={mensaje.id} className={`flex ${esCliente ? 'justify-end' : 'justify-start'}`}>
                <div
                  className={`max-w-[80%] rounded-2xl px-4 py-2.5 ${
                    esCliente
                      ? 'rounded-br-sm bg-acento-500 text-carbono-900'
                      : mensaje.remitente === 'vendedor'
                        ? 'rounded-bl-sm border border-acento-400/30 bg-acento-500/10 text-plata-100'
                        : 'rounded-bl-sm border border-white/8 bg-carbono-800 text-plata-100'
                  }`}
                >
                  {etiqueta && (
                    <p
                      className={`mb-1 text-xs font-semibold ${
                        mensaje.remitente === 'vendedor' ? 'text-acento-300' : 'text-plata-400'
                      }`}
                    >
                      {etiqueta}
                    </p>
                  )}
                  <p className="text-sm whitespace-pre-wrap">{mensaje.contenido}</p>
                  <p
                    className={`mt-1 text-[11px] ${
                      esCliente ? 'text-carbono-800/80' : 'text-plata-500'
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
          {enviando && (
            <div className="flex justify-start">
              <div className="flex items-center gap-1.5 rounded-2xl rounded-bl-sm border border-white/8 bg-carbono-800/70 px-4 py-3">
                <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-plata-400" />
                <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-plata-400 [animation-delay:0.15s]" />
                <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-plata-400 [animation-delay:0.3s]" />
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="border-t border-white/8 p-4">
        {cerrada ? (
          <p className="text-center text-sm text-plata-500">
            Cotización finalizada.{' '}
            <Link to="/catalogo" className="font-semibold text-plata-300 underline-offset-4 hover:underline">
              Seguir explorando el catálogo
            </Link>
          </p>
        ) : (
          <div className="flex items-end gap-2">
            <textarea
              value={nuevoMensaje}
              onChange={(e) => setNuevoMensaje(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Preguntá sobre precios, financiación o planes…"
              className="campo resize-none"
              rows={2}
              disabled={enviando}
            />
            <Boton onClick={handleEnviar} disabled={enviando || !nuevoMensaje.trim()}>
              {enviando ? '…' : 'Enviar'}
            </Boton>
          </div>
        )}
      </div>
    </div>
  )
}