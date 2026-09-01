import { useCallback, useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router'
import { ChatConsulta } from '../components/ChatConsulta'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'
import { api, ErrorApi } from '../services/api'
import type { ConsultaResumen } from '../types/consulta'
import {
  estilosEstadoConsulta,
  etiquetasEstadoConsulta,
  EtiquetaEstado,
} from '../components/ui/EtiquetaEstado'

export function ChatVendedor() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [consulta, setConsulta] = useState<ConsultaResumen | null>(null)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [tomando, setTomando] = useState(false)

  const cargarConsulta = useCallback(async () => {
    try {
      const consultas = await api.listarBandeja()
      const encontrada = consultas.find((c) => c.id === Number(id)) ?? null
      if (encontrada) {
        setConsulta(encontrada)
        setError(null)
      } else {
        setError('La consulta no existe o no está disponible en tu bandeja.')
      }
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo cargar la consulta.')
    } finally {
      setCargando(false)
    }
  }, [id])

  useEffect(() => {
    if (!id) return
    setCargando(true)
    setError(null)
    void cargarConsulta()
  }, [id, cargarConsulta])

  const handleTomar = async () => {
    if (!consulta || tomando) return
    setTomando(true)
    setError(null)
    try {
      await api.tomarConsulta(consulta.id)
      await cargarConsulta()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo tomar la consulta')
    } finally {
      setTomando(false)
    }
  }

  if (!id) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-24 text-center sm:px-6">
        <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">Error</p>
        <h1 className="font-display text-3xl font-bold text-plata-100">Consulta no encontrada</h1>
        <p className="mt-3 text-plata-400">La consulta solicitada no existe.</p>
        <div className="mt-8">
          <Boton onClick={() => navigate('/vendedor/bandeja')}>Volver a la bandeja</Boton>
        </div>
      </div>
    )
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando consulta…" />
  }

  if (error || !consulta) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-24 text-center sm:px-6">
        <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">Error</p>
        <h1 className="font-display text-3xl font-bold text-plata-100">Consulta no encontrada</h1>
        <p className="mt-3 text-plata-400">{error ?? 'La consulta solicitada no existe.'}</p>
        <div className="mt-8">
          <Boton onClick={() => navigate('/vendedor/bandeja')}>Volver a la bandeja</Boton>
        </div>
      </div>
    )
  }

  const estado = consulta.estado
  const pendiente = estado === 'pendiente'

  return (
    <div className="mx-auto flex h-[calc(100vh-6.5rem)] max-w-4xl flex-col px-4 py-6 sm:px-6">
      <div className="mb-3 flex items-center justify-between rounded-2xl border border-white/8 bg-carbono-850/60 px-5 py-4 backdrop-blur-sm">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => navigate('/vendedor/bandeja')}
            className="inline-flex items-center gap-1.5 font-display text-sm font-medium text-plata-400 transition-colors hover:text-plata-100"
          >
            <span aria-hidden>←</span> Bandeja
          </button>
          <span className="h-4 w-px bg-white/10" />
          <h1 className="font-display text-lg font-semibold text-plata-100">Consulta #{id}</h1>
        </div>
        <div className="flex shrink-0 items-center gap-3">
          <Link
            to={`/catalogo/${consulta.vehiculo.id}`}
            className="text-xs font-semibold text-acento-400 transition-colors hover:text-acento-300"
          >
            Ver ficha
          </Link>
          <EtiquetaEstado estado={estado} estilos={estilosEstadoConsulta} etiqueta={etiquetasEstadoConsulta[estado]} />
        </div>
      </div>

      {error && (
        <div className="mb-3 rounded-xl border border-red-500/30 bg-red-500/10 p-3 text-sm text-red-300">
          {error}
        </div>
      )}

      {pendiente ? (
        <div className="flex flex-1 flex-col items-center justify-center gap-4 rounded-2xl border border-white/8 bg-carbono-850/60 p-8 text-center backdrop-blur-sm">
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
            <p className="font-display text-base font-semibold text-plata-100">Consulta pendiente</p>
            <p className="mt-1 max-w-sm text-sm text-plata-400">
              Tomá la consulta para ver la conversación y responderle al cliente.
            </p>
          </div>
          <Boton onClick={handleTomar} disabled={tomando}>
            {tomando ? 'Tomando…' : 'Tomar consulta'}
          </Boton>
        </div>
      ) : (
        <div className="min-h-0 flex-1 overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/60 backdrop-blur-sm">
          <ChatConsulta consultaId={consulta.id} estado={estado} vehiculoId={consulta.vehiculo.id} />
        </div>
      )}
    </div>
  )
}
