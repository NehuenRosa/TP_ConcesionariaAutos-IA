import { useEffect, useState, useCallback } from 'react'
import { useNavigate } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { ConsultaResumen, EstadoConsulta } from '../types/consulta'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
import { estilosEstadoConsulta, etiquetasEstadoConsulta, EtiquetaEstado } from '../components/ui/EtiquetaEstado'
import { formatearFechaHora } from '../utils/formato'

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
    if (!window.confirm('¿Eliminar esta consulta y su conversación? Esta acción no se puede deshacer.')) {
      return
    }
    try {
      await api.eliminarConsulta(id)
      await cargarBandeja()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo eliminar la consulta')
    }
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando bandeja de entrada…" />
  }

  return (
    <div className="mx-auto max-w-5xl px-4 py-12 sm:px-6">
      <EncabezadoPagina
        destacado="Vendedor"
        titulo="Bandeja de entrada"
        descripcion="Consultas asignadas y disponibles. Tomá, respondé y gestioná cada conversación."
      />

      {error && (
        <div className="mt-6 rounded-xl border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-300">
          {error}
        </div>
      )}

      {consultas.length === 0 ? (
        <div className="mt-8">
          <EstadoVacio
            titulo="No tenés consultas asignadas"
            descripcion="Cuando un cliente inicie una consulta o cotización, va a aparecer acá."
          />
        </div>
      ) : (
        <div className="mt-8 grid gap-4">
          {consultas.map((consulta) => {
            const estado = consulta.estado as EstadoConsulta
            const nueva = consulta.mensajesNuevos > 0
            return (
              <article
                key={consulta.id}
                className={`relative rounded-2xl border p-5 shadow-luz backdrop-blur-sm transition-all ${
                  nueva
                    ? 'border-acento-400/40 bg-acento-400/5'
                    : 'border-white/8 bg-carbono-850/60'
                }`}
              >
                {nueva && (
                  <span className="absolute top-4 right-4 flex h-2.5 w-2.5">
                    <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-acento-400 opacity-60" />
                    <span className="relative inline-flex h-2.5 w-2.5 rounded-full bg-acento-400" />
                  </span>
                )}

                <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                  <div
                    className="min-w-0 flex-1 cursor-pointer"
                    onClick={() => navigate(`/vendedor/bandeja/${consulta.id}`)}
                  >
                    <div className="flex flex-wrap items-center gap-3">
                      <h3 className="font-display text-lg font-semibold text-plata-100">
                        {consulta.vehiculo.marca} {consulta.vehiculo.modelo}{' '}
                        <span className="font-normal text-plata-400">· {consulta.vehiculo.anio}</span>
                      </h3>
                      <EtiquetaEstado estado={estado} estilos={estilosEstadoConsulta} etiqueta={etiquetasEstadoConsulta[estado]} />
                    </div>

                    <p className="mt-2 text-sm text-plata-400">
                      Cliente: <span className="text-plata-300">{consulta.cliente.nombre}</span>
                    </p>

                    {consulta.ultimoMensaje && (
                      <p className="mt-2 truncate text-sm text-plata-400">
                        {consulta.ultimoMensaje.contenido}
                      </p>
                    )}

                    <p className="mt-1 text-xs text-plata-500">{formatearFechaHora(consulta.updatedAt)}</p>
                  </div>

                  <div className="flex shrink-0 flex-wrap gap-2">
                    {estado === 'pendiente' && (
                      <Boton tamano="sm" onClick={() => handleTomar(consulta.id)}>
                        Tomar consulta
                      </Boton>
                    )}
                    {estado === 'en_conversacion' && (
                      <Boton variante="secundario" tamano="sm" onClick={() => handleCerrar(consulta.id)}>
                        Cerrar
                      </Boton>
                    )}
                    {estado === 'cerrada' && (
                      <Boton variante="peligro" tamano="sm" onClick={() => handleEliminar(consulta.id)}>
                        Eliminar
                      </Boton>
                    )}
                    <Boton
                      variante="fantasma"
                      tamano="sm"
                      onClick={() => navigate(`/vendedor/bandeja/${consulta.id}`)}
                    >
                      Abrir chat
                    </Boton>
                    <Boton
                      variante="fantasma"
                      tamano="sm"
                      onClick={() => navigate(`/catalogo/${consulta.vehiculo.id}`)}
                    >
                      Ver ficha
                    </Boton>
                  </div>
                </div>
              </article>
            )
          })}
        </div>
      )}
    </div>
  )
}
