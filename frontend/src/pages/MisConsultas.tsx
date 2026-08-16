import { useEffect, useState, useCallback } from 'react'
import { Link, useParams, useNavigate } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { ConsultaResumen, EstadoConsulta } from '../types/consulta'
import { ChatConsulta } from '../components/ChatConsulta'
import { ContenidoCargando } from '../components/ui/Spinner'
import { Boton } from '../components/ui/Boton'
import { estilosEstadoConsulta, etiquetasEstadoConsulta, EtiquetaEstado } from '../components/ui/EtiquetaEstado'
import { formatearFechaHora } from '../utils/formato'

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

  useEffect(() => {
    const manejarLeidos = () => cargarConsultas()
    window.addEventListener('mensajes-leidos', manejarLeidos)

    const intervalo = setInterval(cargarConsultas, 5000)

    return () => {
      window.removeEventListener('mensajes-leidos', manejarLeidos)
      clearInterval(intervalo)
    }
  }, [cargarConsultas])

  const seleccionada = id ? consultas.find((c) => c.id === Number(id)) ?? null : null
  const idInvalido = id !== undefined && seleccionada === null && !cargando

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando consultas…" />
  }

  return (
    <div className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-10 sm:px-6 lg:h-[calc(100vh-8rem)] lg:flex-row lg:px-8">
      <aside className="flex w-full shrink-0 flex-col overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/60 backdrop-blur-sm lg:w-80">
        <div className="border-b border-white/8 px-5 py-4">
          <h1 className="font-display text-lg font-semibold text-plata-100">Mis consultas</h1>
        </div>

        {error && (
          <div className="border-b border-white/8 bg-red-500/10 px-5 py-3 text-sm text-red-300">
            {error}
          </div>
        )}

        {consultas.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-sm text-plata-400">No tenés consultas todavía.</p>
            <Boton variante="secundario" tamano="sm" className="mt-4">
              <Link to="/catalogo">Explorar catálogo</Link>
            </Boton>
          </div>
        ) : (
          <div className="flex-1 divide-y divide-white/6 overflow-y-auto">
            {consultas.map((consulta) => {
              const activa = seleccionada?.id === consulta.id
              const estado = consulta.estado as EstadoConsulta
              return (
                <button
                  key={consulta.id}
                  type="button"
                  onClick={() => navigate(`/mis-consultas/${consulta.id}`)}
                  className={`relative w-full p-4 text-left transition-colors ${
                    activa ? 'bg-white/6' : 'hover:bg-white/4'
                  }`}
                >
                  {consulta.mensajesNuevos > 0 && (
                    <span className="absolute top-3 right-3 h-2.5 w-2.5">
                      <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-acento-400 opacity-60" />
                      <span className="relative inline-flex h-2.5 w-2.5 rounded-full bg-acento-400" />
                    </span>
                  )}

                  <div className="flex items-center gap-2 pr-4">
                    <h3 className="truncate font-display text-sm font-semibold text-plata-100">
                      {consulta.vehiculo.marca} {consulta.vehiculo.modelo}
                    </h3>
                  </div>
                  <div className="mt-2">
                    <EtiquetaEstado estado={estado} estilos={estilosEstadoConsulta} etiqueta={etiquetasEstadoConsulta[estado]} />
                  </div>

                  {consulta.vendedor && (
                    <p className="mt-2 text-xs text-plata-400">
                      Vendedor: <span className="text-plata-300">{consulta.vendedor.nombre}</span>
                    </p>
                  )}

                  {consulta.ultimoMensaje && (
                    <p className="mt-1 truncate text-sm text-plata-400">
                      {consulta.ultimoMensaje.contenido}
                    </p>
                  )}

                  <p className="mt-1 text-xs text-plata-500">{formatearFechaHora(consulta.updatedAt)}</p>
                </button>
              )
            })}
          </div>
        )}
      </aside>

      <section className="flex min-h-[28rem] flex-1 flex-col overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/60 backdrop-blur-sm">
        {seleccionada ? (
          <ChatConsulta
            consultaId={seleccionada.id}
            estado={seleccionada.estado}
            onMensajeEnviado={cargarConsultas}
          />
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
              <p className="font-display text-base font-semibold text-plata-100">Consulta no encontrada</p>
              <p className="mt-1 max-w-xs text-sm text-plata-400">
                La consulta no existe o no te pertenece.
              </p>
            </div>
            <Boton variante="secundario" tamano="sm" onClick={() => navigate('/mis-consultas')}>
              Volver a mis consultas
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
                Seleccioná una consulta para ver la conversación con tu vendedor.
              </p>
            </div>
          </div>
        )}
      </section>
    </div>
  )
}
