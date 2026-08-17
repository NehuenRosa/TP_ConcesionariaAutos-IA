import { useEffect, useState, useCallback } from 'react'
import { Link, useParams, useNavigate } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { CotizacionResumen, EstadoCotizacion } from '../types/cotizacion'
import { ChatCotizacion } from '../components/ChatCotizacion'
import { ContenidoCargando } from '../components/ui/Spinner'
import { Boton } from '../components/ui/Boton'
import { formatearFechaHora } from '../utils/formato'

const etiquetasEstadoCotizacion: Record<EstadoCotizacion, string> = {
  abierta: 'Abierta',
  cerrada: 'Cerrada',
}

function EtiquetaCotizacion({ estado }: { estado: EstadoCotizacion }) {
  const abierta = estado === 'abierta'
  return (
    <span
      className={`inline-flex rounded-full border px-2.5 py-1 text-[11px] font-semibold tracking-wide uppercase ${
        abierta
          ? 'border-emerald-400/30 bg-emerald-400/10 text-emerald-300'
          : 'border-white/10 bg-white/5 text-plata-400'
      }`}
    >
      {etiquetasEstadoCotizacion[estado]}
    </span>
  )
}

export function MisCotizaciones() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [cotizaciones, setCotizaciones] = useState<CotizacionResumen[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const cargarCotizaciones = useCallback(async () => {
    try {
      const datos = await api.listarMisCotizaciones()
      setCotizaciones(datos)
      setError(null)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'Error al cargar las cotizaciones')
    } finally {
      setCargando(false)
    }
  }, [])

  useEffect(() => {
    cargarCotizaciones()
  }, [cargarCotizaciones])

  const seleccionada = id ? cotizaciones.find((c) => c.id === Number(id)) ?? null : null
  const idInvalido = id !== undefined && seleccionada === null && !cargando

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando cotizaciones…" />
  }

  return (
    <div className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-10 sm:px-6 lg:h-[calc(100vh-8rem)] lg:flex-row lg:px-8">
      <aside className="flex w-full shrink-0 flex-col overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/60 backdrop-blur-sm lg:w-80">
        <div className="border-b border-white/8 px-5 py-4">
          <h1 className="font-display text-lg font-semibold text-plata-100">Mis cotizaciones</h1>
          <p className="mt-0.5 text-xs text-plata-500">Atendidas por el asistente IA</p>
        </div>

        {error && (
          <div className="border-b border-white/8 bg-red-500/10 px-5 py-3 text-sm text-red-300">
            {error}
          </div>
        )}

        {cotizaciones.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-sm text-plata-400">
              No tenés cotizaciones todavía. Pedí una por el chat o desde la ficha de un vehículo.
            </p>
            <Boton variante="secundario" tamano="sm" className="mt-4">
              <Link to="/catalogo">Explorar catálogo</Link>
            </Boton>
          </div>
        ) : (
          <div className="flex-1 divide-y divide-white/6 overflow-y-auto">
            {cotizaciones.map((cotizacion) => {
              const activa = seleccionada?.id === cotizacion.id
              return (
                <button
                  key={cotizacion.id}
                  type="button"
                  onClick={() => navigate(`/mis-cotizaciones/${cotizacion.id}`)}
                  className={`relative w-full p-4 text-left transition-colors ${
                    activa ? 'bg-white/6' : 'hover:bg-white/4'
                  }`}
                >
                  <div className="flex items-center gap-2 pr-4">
                    <h3 className="truncate font-display text-sm font-semibold text-plata-100">
                      {cotizacion.vehiculo.marca} {cotizacion.vehiculo.modelo}
                    </h3>
                  </div>
                  <div className="mt-2">
                    <EtiquetaCotizacion estado={cotizacion.estado} />
                  </div>

                  {cotizacion.ultimoMensaje && (
                    <p className="mt-1 truncate text-sm text-plata-400">
                      {cotizacion.ultimoMensaje.contenido}
                    </p>
                  )}

                  <p className="mt-1 text-xs text-plata-500">
                    {formatearFechaHora(cotizacion.updatedAt)}
                  </p>
                </button>
              )
            })}
          </div>
        )}
      </aside>

      <section className="flex min-h-[28rem] flex-1 flex-col overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/60 backdrop-blur-sm">
        {seleccionada ? (
          <ChatCotizacion cotizacionId={seleccionada.id} onCambio={cargarCotizaciones} />
        ) : idInvalido ? (
          <div className="flex h-full flex-col items-center justify-center gap-3 p-8 text-center">
            <div className="flex h-14 w-14 items-center justify-center rounded-full border border-red-500/30 bg-red-500/10 text-red-300">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-7 w-7">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"
                />
              </svg>
            </div>
            <div>
              <p className="font-display text-base font-semibold text-plata-100">
                Cotización no encontrada
              </p>
              <p className="mt-1 max-w-xs text-sm text-plata-400">
                La cotización no existe o no te pertenece.
              </p>
            </div>
            <Boton variante="secundario" tamano="sm" onClick={() => navigate('/mis-cotizaciones')}>
              Volver a mis cotizaciones
            </Boton>
          </div>
        ) : (
          <div className="flex h-full flex-col items-center justify-center gap-3 p-8 text-center">
            <div className="flex h-14 w-14 items-center justify-center rounded-full border border-white/10 bg-carbono-800 text-plata-400">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-7 w-7">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.86 9.86 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
            </div>
            <div>
              <p className="font-display text-base font-semibold text-plata-100">Conversación</p>
              <p className="mt-1 max-w-xs text-sm text-plata-400">
                Seleccioná una cotización para seguir charlando con el asistente sobre ese
                vehículo.
              </p>
            </div>
          </div>
        )}
      </section>
    </div>
  )
}