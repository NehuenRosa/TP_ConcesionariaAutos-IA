import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { CotizacionResumen, EstadoCotizacion } from '../types/cotizacion'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
import { formatearPrecio } from '../utils/formato'

const INTERVALO_ACTUALIZACION_MS = 10000

const estilosEstadoCotizacion: Record<EstadoCotizacion, string> = {
  abierta: 'border-emerald-400/30 bg-emerald-400/10 text-emerald-300',
  cerrada: 'border-white/10 bg-white/5 text-plata-400',
}

const etiquetasEstadoCotizacion: Record<EstadoCotizacion, string> = {
  abierta: 'Abierta',
  cerrada: 'Cerrada',
}

function EtiquetaEstadoCotizacion({ estado }: { estado: EstadoCotizacion }) {
  return (
    <span
      className={`inline-flex rounded-full border px-2.5 py-1 text-[11px] font-semibold tracking-wide uppercase ${estilosEstadoCotizacion[estado]}`}
    >
      {etiquetasEstadoCotizacion[estado]}
    </span>
  )
}

// etiquetaAtencion describe quién atiende la conversación hoy.
function etiquetaAtencion(cotizacion: CotizacionResumen): string {
  if (cotizacion.vendedor) return `Atendida por ${cotizacion.vendedor.nombre}`
  return 'Con la IA'
}

export function BandejaCotizaciones() {
  const [cotizaciones, setCotizaciones] = useState<CotizacionResumen[]>([])
  const [filtro, setFiltro] = useState<EstadoCotizacion | ''>('')
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const cargarBandeja = useCallback(async () => {
    try {
      const datos = await api.listarBandejaCotizaciones()
      setCotizaciones(datos)
      setError(null)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'Error al cargar la bandeja')
    } finally {
      setCargando(false)
    }
  }, [])

  useEffect(() => {
    cargarBandeja()
    // La bandeja se refresca sola para ver las conversaciones nuevas.
    const temporizador = setInterval(cargarBandeja, INTERVALO_ACTUALIZACION_MS)
    return () => clearInterval(temporizador)
  }, [cargarBandeja])

  const filtradas =
    filtro === '' ? cotizaciones : cotizaciones.filter((cotizacion) => cotizacion.estado === filtro)

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando cotizaciones…" />
  }

  return (
    <div className="mx-auto max-w-5xl px-4 py-12 sm:px-6">
      <EncabezadoPagina
        destacado="Vendedor"
        titulo="Cotizaciones con IA"
        descripcion="Conversaciones de los clientes con el asistente sobre precios y financiación. Tomá una conversación para atenderla personalmente."
      />

      <div className="mt-6 flex flex-wrap gap-2">
        <Boton
          tamano="sm"
          variante={filtro === '' ? 'acento' : 'secundario'}
          onClick={() => setFiltro('')}
        >
          Todas
        </Boton>
        <Boton
          tamano="sm"
          variante={filtro === 'abierta' ? 'acento' : 'secundario'}
          onClick={() => setFiltro('abierta')}
        >
          Abiertas
        </Boton>
        <Boton
          tamano="sm"
          variante={filtro === 'cerrada' ? 'acento' : 'secundario'}
          onClick={() => setFiltro('cerrada')}
        >
          Cerradas
        </Boton>
      </div>

      {error && (
        <div className="mt-6 rounded-xl border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-300">
          {error}
        </div>
      )}

      {filtradas.length === 0 ? (
        <div className="mt-8">
          <EstadoVacio
            titulo="No hay cotizaciones"
            descripcion="Cuando un cliente hable con el asistente sobre un vehículo, la conversación va a aparecer acá."
          />
        </div>
      ) : (
        <div className="mt-8 space-y-3">
          {filtradas.map((cotizacion) => (
            <Link
              key={cotizacion.id}
              to={`/vendedor/cotizaciones/${cotizacion.id}`}
              className="flex items-center justify-between gap-4 rounded-xl border border-white/8 bg-carbono-850/60 p-4 shadow-luz backdrop-blur-sm transition-colors hover:border-acento-400/40 hover:bg-carbono-800/60"
            >
              <div className="min-w-0 flex-1">
                <div className="flex flex-wrap items-center gap-2">
                  <p className="font-medium text-plata-100">
                    {cotizacion.cliente?.nombre ?? 'Cliente'}
                    <span className="text-plata-500"> · </span>
                    <span className="text-plata-300">
                      {cotizacion.vehiculo.marca} {cotizacion.vehiculo.modelo} {cotizacion.vehiculo.anio}
                    </span>
                  </p>
                  <EtiquetaEstadoCotizacion estado={cotizacion.estado} />
                </div>
                <p className="mt-1 truncate text-sm text-plata-400">
                  {cotizacion.ultimoMensaje?.contenido ?? 'Sin mensajes'}
                </p>
                <p className="mt-1 text-xs text-acento-400">{etiquetaAtencion(cotizacion)}</p>
              </div>
              <p className="shrink-0 texto-numerico text-sm font-medium text-plata-200">
                {formatearPrecio(cotizacion.vehiculo.precio)}
              </p>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
