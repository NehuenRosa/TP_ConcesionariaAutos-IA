import { useCallback, useEffect, useRef, useState } from 'react'
import { Link, useParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { Cotizacion } from '../types/cotizacion'
import { useAuth } from '../hooks/useAuth'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'

const INTERVALO_ACTUALIZACION_MS = 10000

// ChatCotizacionVendedor es la vista de atención de una cotización: el
// vendedor puede tomar la conversación (la IA queda silenciada) y responderle
// al cliente en su nombre.
export function ChatCotizacionVendedor() {
  const { id } = useParams<{ id: string }>()
  const cotizacionId = Number(id)
  const { usuario } = useAuth()
  const [cotizacion, setCotizacion] = useState<Cotizacion | null>(null)
  const [cargando, setCargando] = useState(true)
  const [tomando, setTomando] = useState(false)
  const [enviando, setEnviando] = useState(false)
  const [cerrando, setCerrando] = useState(false)
  const [entrada, setEntrada] = useState('')
  const [error, setError] = useState<string | null>(null)
  const mensajesRef = useRef<HTMLDivElement>(null)
  const ultimoIdRef = useRef(0)

  const cargarCotizacion = useCallback(async () => {
    try {
      const datos = await api.obtenerCotizacionPersonal(cotizacionId)
      setCotizacion(datos)
      ultimoIdRef.current = datos.mensajes[datos.mensajes.length - 1]?.id ?? 0
      setError(null)
      // Si es el vendedor asignado, abrir el hilo marca los mensajes como
      // leídos en el servidor: refresca el contador.
      window.dispatchEvent(new Event('mensajes-leidos'))
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'Error al cargar la cotización')
    } finally {
      setCargando(false)
    }
  }, [cotizacionId])

  useEffect(() => {
    if (!Number.isFinite(cotizacionId) || cotizacionId <= 0) return
    cargarCotizacion()
  }, [cargarCotizacion, cotizacionId])

  // Polling incremental (fetch por desdeId): el historial completo se carga
  // una sola vez y el temporizador trae solo los mensajes nuevos.
  useEffect(() => {
    if (!Number.isFinite(cotizacionId) || cotizacionId <= 0) return
    const temporizador = setInterval(async () => {
      try {
        const datos = await api.obtenerMensajesCotizacionPersonalDesde(
          cotizacionId,
          ultimoIdRef.current,
        )
        setCotizacion((prev) => {
          if (!prev) return prev
          const idsExistentes = new Set(prev.mensajes.map((m) => m.id))
          const aAgregar = datos.mensajes.filter((m) => !idsExistentes.has(m.id))
          const mensajes = aAgregar.length > 0 ? [...prev.mensajes, ...aAgregar] : prev.mensajes
          return { ...prev, mensajes, estado: datos.estado, vendedor: datos.vendedor, fechaToma: datos.fechaToma }
        })
        if (datos.mensajes.length > 0) {
          ultimoIdRef.current = datos.mensajes[datos.mensajes.length - 1].id
        }
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
  }, [cotizacion?.mensajes.length])

  const handleTomar = async () => {
    if (!cotizacion || tomando) return
    setTomando(true)
    setError(null)
    try {
      await api.tomarCotizacion(cotizacion.id)
      await cargarCotizacion()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo tomar la cotización')
    } finally {
      setTomando(false)
    }
  }

  const handleResponder = async () => {
    const contenido = entrada.trim()
    if (!contenido || !cotizacion || enviando) return
    setEnviando(true)
    setError(null)
    try {
      const actualizada = await api.responderCotizacionVendedor(cotizacion.id, contenido)
      setCotizacion(actualizada)
      setEntrada('')
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo enviar el mensaje')
    } finally {
      setEnviando(false)
    }
  }

  const handleCerrar = async () => {
    if (!cotizacion || cerrando) return
    setCerrando(true)
    setError(null)
    try {
      const cerrada = await api.cerrarCotizacionPersonal(cotizacion.id)
      setCotizacion(cerrada)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo cerrar la cotización')
    } finally {
      setCerrando(false)
    }
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando conversación…" />
  }

  if (!cotizacion) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-16 text-center">
        <p className="text-sm text-plata-400">
          La cotización no existe. Volvé a la bandeja para elegir una conversación.
        </p>
        <Boton variante="secundario" tamano="sm" className="mt-4">
          <Link to="/vendedor/cotizaciones">Volver a la bandeja</Link>
        </Boton>
      </div>
    )
  }

  const tomada = cotizacion.vendedor !== undefined && cotizacion.vendedor !== null
  const mia = Boolean(cotizacion.vendedor && usuario && cotizacion.vendedor.id === usuario.id)
  const abierta = cotizacion.estado === 'abierta'
  // Puede escribir solo si está abierta y la tomó este vendedor (el backend
  // revalida el id en cada pedido).
  const puedeEscribir = abierta && mia

  return (
    <div className="mx-auto flex h-[calc(100vh-10rem)] w-full max-w-3xl flex-col overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/60 shadow-luz backdrop-blur-sm">
      <div className="flex items-center justify-between gap-3 border-b border-white/8 px-5 py-3.5">
        <div className="min-w-0">
          <p className="font-display text-sm font-semibold text-plata-100">
            {cotizacion.cliente?.nombre ?? 'Cliente'} ·{' '}
            <span className="text-plata-300">
              {cotizacion.vehiculo.marca} {cotizacion.vehiculo.modelo}
            </span>
          </p>
          <p className="text-xs text-plata-500">
            {cotizacion.vehiculo.anio} · Consulta sobre precios y financiación ·{' '}
            <Link
              to={`/catalogo/${cotizacion.vehiculo.id}`}
              className="font-semibold text-acento-400 transition-colors hover:text-acento-300"
            >
              Ver ficha
            </Link>
          </p>
        </div>
        <div className="flex shrink-0 items-center gap-2">
          {!tomada && abierta && (
            <Boton tamano="sm" onClick={handleTomar} disabled={tomando}>
              {tomando ? '…' : 'Tomar conversación'}
            </Boton>
          )}
          {abierta && (
            <Boton variante="secundario" tamano="sm" onClick={handleCerrar} disabled={cerrando}>
              {cerrando ? '…' : 'Cerrar'}
            </Boton>
          )}
        </div>
      </div>

      {tomada && (
        <div className="border-b border-acento-400/20 bg-acento-500/10 px-5 py-2 text-xs text-acento-300">
          Conversación atendida por {cotizacion.vendedor?.nombre}. El asistente quedó en pausa.
        </div>
      )}

      <div ref={mensajesRef} className="flex-1 overflow-y-auto px-5 py-4">
        {!abierta ? (
          <div className="flex h-full items-center justify-center">
            <div className="text-center">
              <p className="font-display text-lg font-semibold text-plata-100">Cotización cerrada</p>
              <p className="mt-1 text-sm text-plata-500">
                La conversación ya no acepta mensajes. El vehículo sigue{' '}
                <Link
                  to={`/catalogo/${cotizacion.vehiculo.id}`}
                  className="font-semibold text-plata-300 underline-offset-4 hover:underline"
                >
                  disponible en el catálogo
                </Link>
                .
              </p>
            </div>
          </div>
        ) : (
          <>
            {error && (
              <div className="mb-4 rounded-xl border border-red-500/30 bg-red-500/10 p-3 text-sm text-red-300">
                {error}
              </div>
            )}

            <div className="space-y-4">
              {cotizacion.mensajes.map((mensaje) => {
                const delCliente = mensaje.remitente === 'cliente'
                const etiqueta =
                  mensaje.remitente === 'ia'
                    ? 'Asistente IA'
                    : mensaje.remitente === 'vendedor'
                      ? 'Asesor de ventas'
                      : null
                return (
                  <div key={mensaje.id} className={`flex ${delCliente ? 'justify-end' : 'justify-start'}`}>
                    <div
                      className={`max-w-[80%] rounded-2xl px-4 py-2.5 ${
                        delCliente
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
                        className={`mt-1 text-[11px] ${delCliente ? 'text-carbono-800/80' : 'text-plata-500'}`}
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
          </>
        )}
      </div>

      <div className="border-t border-white/8 p-4">
        {puedeEscribir ? (
          <div className="flex items-end gap-2">
            <textarea
              value={entrada}
              onChange={(e) => setEntrada(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                  e.preventDefault()
                  void handleResponder()
                }
              }}
              placeholder={`Respondé como asesor a ${cotizacion.cliente?.nombre ?? 'el cliente'}…`}
              className="campo resize-none"
              rows={2}
              disabled={enviando}
            />
            <Boton onClick={handleResponder} disabled={enviando || !entrada.trim()}>
              {enviando ? '…' : 'Enviar'}
            </Boton>
          </div>
        ) : (
          <p className="text-center text-sm text-plata-500">
            {abierta
              ? tomada
                ? 'Solo el vendedor que tomó la conversación puede responder.'
                : 'Tomá la conversación para responderle al cliente.'
              : 'La conversación está finalizada.'}
          </p>
        )}
      </div>
    </div>
  )
}
