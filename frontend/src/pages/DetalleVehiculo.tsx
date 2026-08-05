import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { Vehiculo } from '../types/vehiculo'
import { useAuth } from '../hooks/useAuth'
import { Boton } from '../components/ui/Boton'
import { CampoArea } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { etiquetaCondicion, formatearPrecio } from '../utils/formato'

function formatearKilometraje(kilometraje: number): string {
  return `${new Intl.NumberFormat('es-AR').format(kilometraje)} km`
}

export function DetalleVehiculo() {
  const { id } = useParams<{ id: string }>()
  const { usuario } = useAuth()
  const [vehiculo, setVehiculo] = useState<Vehiculo | null>(null)
  const [imagenActiva, setImagenActiva] = useState(0)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [mostrarFormulario, setMostrarFormulario] = useState(false)
  const [mensajeConsulta, setMensajeConsulta] = useState('')
  const [enviandoConsulta, setEnviandoConsulta] = useState(false)
  const [errorConsulta, setErrorConsulta] = useState<string | null>(null)
  const [exitoConsulta, setExitoConsulta] = useState(false)

  const esCliente = usuario?.rol === 'cliente'

  useEffect(() => {
    if (!id) return
    let cancelado = false
    setCargando(true)
    setError(null)
    setImagenActiva(0)

    api
      .obtenerVehiculo(Number(id))
      .then((dato) => {
        if (cancelado) return
        setVehiculo(dato)
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'Ocurrió un error inesperado al cargar el vehículo.')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [id])

  const handleEnviarConsulta = async () => {
    if (!vehiculo || !mensajeConsulta.trim()) return

    setEnviandoConsulta(true)
    setErrorConsulta(null)

    try {
      await api.crearConsulta({
        vehiculoId: vehiculo.id,
        mensaje: mensajeConsulta.trim(),
      })
      setExitoConsulta(true)
      setMensajeConsulta('')
      setMostrarFormulario(false)
    } catch (e: unknown) {
      setErrorConsulta(e instanceof ErrorApi ? e.message : 'No se pudo enviar la consulta')
    } finally {
      setEnviandoConsulta(false)
    }
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando vehículo…" />
  }

  if (error || !vehiculo) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-24 text-center sm:px-6">
        <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">
          Error
        </p>
        <h1 className="font-display text-3xl font-bold text-plata-100">Vehículo no encontrado</h1>
        <p className="mt-3 text-plata-400">{error ?? 'El vehículo solicitado no existe o no está disponible.'}</p>
        <div className="mt-8">
          <Boton variante="secundario">
            <Link to="/catalogo">Volver al catálogo</Link>
          </Boton>
        </div>
      </div>
    )
  }

  const imagenes = vehiculo.imagenes ?? []

  const ficha = [
    { nombre: 'Año', valor: String(vehiculo.anio) },
    { nombre: 'Kilometraje', valor: formatearKilometraje(vehiculo.kilometraje) },
    { nombre: 'Combustible', valor: vehiculo.combustible },
    { nombre: 'Transmisión', valor: vehiculo.transmision },
    {
      nombre: 'Tipo',
      valor: vehiculo.tipo ? vehiculo.tipo.charAt(0).toUpperCase() + vehiculo.tipo.slice(1) : '—',
    },
    { nombre: 'Condición', valor: etiquetaCondicion(vehiculo.condicion) },
  ]

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <Link
        to="/catalogo"
        className="inline-flex items-center gap-2 font-display text-sm font-medium text-plata-400 transition-colors hover:text-plata-100"
      >
        <span aria-hidden>←</span> Volver al catálogo
      </Link>

      <div className="mt-8 grid grid-cols-1 gap-10 lg:grid-cols-2">
        <div className="space-y-4">
          <div className="relative aspect-[4/3] overflow-hidden rounded-2xl border border-white/8 bg-carbono-850">
            {imagenes[imagenActiva] ? (
              <img
                src={imagenes[imagenActiva].url}
                alt={`${vehiculo.marca} ${vehiculo.modelo}`}
                className="h-full w-full object-cover"
              />
            ) : (
              <div className="flex h-full items-center justify-center text-plata-500">Sin imagen</div>
            )}
            <span className="absolute top-4 left-4 rounded-full border border-white/15 bg-carbono-950/70 px-3 py-1 font-display text-[11px] font-semibold tracking-[0.15em] text-plata-200 uppercase backdrop-blur-sm">
              {etiquetaCondicion(vehiculo.condicion)}
            </span>
          </div>

          {imagenes.length > 1 && (
            <div className="flex gap-3">
              {imagenes.map((imagen, indice) => (
                <button
                  key={imagen.id}
                  type="button"
                  onClick={() => setImagenActiva(indice)}
                  aria-label={`Ver imagen ${indice + 1}`}
                  className={`aspect-[4/3] w-24 overflow-hidden rounded-lg border transition ${
                    indice === imagenActiva
                      ? 'border-acento-400 shadow-resaltado'
                      : 'border-white/10 opacity-70 hover:opacity-100'
                  }`}
                >
                  <img src={imagen.url} alt="" className="h-full w-full object-cover" />
                </button>
              ))}
            </div>
          )}
        </div>

        <div className="space-y-8">
          <div>
            <p className="font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">
              {etiquetaCondicion(vehiculo.condicion)}
            </p>
            <h1 className="mt-1 font-display text-3xl font-bold tracking-tight text-plata-100 sm:text-4xl">
              {vehiculo.marca} {vehiculo.modelo}
            </h1>
            <p className="mt-2 text-plata-400">{vehiculo.anio}</p>
            <p className="texto-numerico mt-5 font-display text-4xl font-bold text-plata-100">
              {formatearPrecio(vehiculo.precio)}
            </p>
          </div>

          <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            {ficha.map((dato) => (
              <div key={dato.nombre} className="rounded-xl border border-white/8 bg-carbono-850/60 p-4">
                <dt className="font-display text-[11px] font-semibold tracking-[0.2em] text-plata-500 uppercase">
                  {dato.nombre}
                </dt>
                <dd className="mt-1 text-lg font-semibold text-plata-100">{dato.valor}</dd>
              </div>
            ))}
          </dl>

          <div className="space-y-3 border-t border-white/8 pt-6">
            {esCliente ? (
              <>
                <Boton tamano="lg" className="w-full">
                  <Link to={`/catalogo/${vehiculo.id}/test-drive`}>Solicitar test drive</Link>
                </Boton>

                {exitoConsulta ? (
                  <div className="rounded-xl border border-emerald-400/30 bg-emerald-400/10 p-4 text-sm text-emerald-300">
                    <p>Consulta enviada correctamente. Un vendedor te responderá pronto.</p>
                    <Link
                      to="/mis-consultas"
                      className="mt-2 inline-block font-semibold text-emerald-200 underline-offset-4 hover:underline"
                    >
                      Ver mis consultas →
                    </Link>
                  </div>
                ) : mostrarFormulario ? (
                  <div className="space-y-4">
                    <h3 className="font-display text-lg font-semibold text-plata-100">Enviar consulta</h3>
                    <CampoArea
                      value={mensajeConsulta}
                      onChange={(e) => setMensajeConsulta(e.target.value)}
                      placeholder="Escribí tu consulta sobre este vehículo…"
                      rows={4}
                      disabled={enviandoConsulta}
                      error={errorConsulta ?? undefined}
                    />
                    <div className="flex gap-3">
                      <Boton
                        variante="acento"
                        onClick={handleEnviarConsulta}
                        disabled={enviandoConsulta || !mensajeConsulta.trim()}
                      >
                        {enviandoConsulta ? 'Enviando…' : 'Enviar consulta'}
                      </Boton>
                      <Boton
                        variante="secundario"
                        onClick={() => {
                          setMostrarFormulario(false)
                          setErrorConsulta(null)
                        }}
                      >
                        Cancelar
                      </Boton>
                    </div>
                  </div>
                ) : (
                  <Boton variante="secundario" tamano="lg" className="w-full" onClick={() => setMostrarFormulario(true)}>
                    Consultar sobre este vehículo
                  </Boton>
                )}
              </>
            ) : (
              <div className="rounded-xl border border-white/8 bg-carbono-850/60 p-4 text-sm text-plata-400">
                <Link to="/login" className="font-semibold text-plata-100 underline-offset-4 hover:underline">
                  Iniciá sesión
                </Link>{' '}
                como cliente para consultar o reservar este vehículo.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
