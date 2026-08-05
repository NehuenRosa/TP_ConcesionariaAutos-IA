import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { Vehiculo } from '../types/vehiculo'
import type { FranjaHoraria } from '../types/testDrive'
import { Boton } from '../components/ui/Boton'
import { CampoTexto } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { formatearPrecio } from '../utils/formato'

function hoyISO(): string {
  return new Date().toISOString().slice(0, 10)
}

export function FormularioTestDrive() {
  const { id } = useParams<{ id: string }>()
  const [vehiculo, setVehiculo] = useState<Vehiculo | null>(null)
  const [franjas, setFranjas] = useState<FranjaHoraria[]>([])
  const [fecha, setFecha] = useState(hoyISO())
  const [franja, setFranja] = useState('')
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [enviando, setEnviando] = useState(false)
  const [errorEnvio, setErrorEnvio] = useState<string | null>(null)
  const [turnoCreado, setTurnoCreado] = useState(false)

  useEffect(() => {
    if (!id) return
    let cancelado = false

    Promise.all([api.obtenerVehiculo(Number(id)), api.obtenerFranjas()])
      .then(([dato, franjasDisponibles]) => {
        if (cancelado) return
        setVehiculo(dato)
        setFranjas(franjasDisponibles)
        if (franjasDisponibles.length > 0) {
          setFranja(franjasDisponibles[0].id)
        }
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'No se pudo cargar la información del vehículo.')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [id])

  const handleSolicitar = async () => {
    if (!vehiculo || !fecha || !franja) return

    setEnviando(true)
    setErrorEnvio(null)

    try {
      await api.solicitarTestDrive({
        vehiculoId: vehiculo.id,
        fecha,
        franja,
      })
      setTurnoCreado(true)
    } catch (e: unknown) {
      setErrorEnvio(e instanceof ErrorApi ? e.message : 'No se pudo solicitar el test drive')
    } finally {
      setEnviando(false)
    }
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando…" />
  }

  if (error || !vehiculo) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-24 text-center sm:px-6">
        <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">Error</p>
        <h1 className="font-display text-3xl font-bold text-plata-100">No se pudo solicitar el test drive</h1>
        <p className="mt-3 text-plata-400">{error ?? 'El vehículo solicitado no existe o no está disponible.'}</p>
        <div className="mt-8">
          <Boton variante="secundario">
            <Link to="/catalogo">Volver al catálogo</Link>
          </Boton>
        </div>
      </div>
    )
  }

  if (turnoCreado) {
    return (
      <div className="mx-auto max-w-xl px-4 py-20 sm:px-6">
        <div className="rounded-2xl border border-emerald-400/30 bg-emerald-400/10 p-8 text-center shadow-luz">
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full border border-emerald-400/40 bg-emerald-400/15 text-emerald-300">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="h-7 w-7">
              <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h1 className="font-display text-2xl font-bold text-plata-100">Test drive solicitado</h1>
          <p className="mt-3 text-sm leading-relaxed text-plata-300">
            Solicitaste el test drive de{' '}
            <span className="font-semibold text-plata-100">
              {vehiculo.marca} {vehiculo.modelo}
            </span>{' '}
            para el {fecha} en la franja de{' '}
            <span className="capitalize">{franja === 'manana' ? 'la mañana' : 'la tarde'}</span>. Un
            vendedor te confirmará el turno.
          </p>
          <div className="mt-6 flex flex-wrap justify-center gap-3">
            <Boton>
              <Link to="/mis-test-drives">Ver mis test drives</Link>
            </Boton>
            <Boton variante="secundario">
              <Link to="/catalogo">Seguir explorando</Link>
            </Boton>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-2xl px-4 py-12 sm:px-6">
      <Link
        to={`/catalogo/${vehiculo.id}`}
        className="inline-flex items-center gap-2 font-display text-sm font-medium text-plata-400 transition-colors hover:text-plata-100"
      >
        <span aria-hidden>←</span> Volver al detalle del vehículo
      </Link>

      <div className="mt-6 rounded-2xl border border-white/8 bg-carbono-850/60 p-6 shadow-luz backdrop-blur-sm sm:p-8">
        <h1 className="font-display text-2xl font-bold text-plata-100">Solicitar test drive</h1>

        <div className="mt-5 flex items-center gap-4 rounded-xl border border-white/8 bg-carbono-900/60 p-4">
          {vehiculo.imagenes[0] && (
            <img
              src={vehiculo.imagenes[0].url}
              alt={`${vehiculo.marca} ${vehiculo.modelo}`}
              className="h-20 w-28 rounded-lg object-cover"
            />
          )}
          <div>
            <p className="font-display text-lg font-semibold text-plata-100">
              {vehiculo.marca} {vehiculo.modelo}
            </p>
            <p className="text-sm text-plata-400">
              {vehiculo.anio} · {vehiculo.condicion === 'nuevo' ? 'Nuevo' : 'Usado'} ·{' '}
              <span className="texto-numerico text-plata-300">{formatearPrecio(vehiculo.precio)}</span>
            </p>
          </div>
        </div>

        <div className="mt-6 space-y-5">
          <CampoTexto
            id="fecha"
            etiqueta="Fecha"
            type="date"
            value={fecha}
            min={hoyISO()}
            onChange={(e) => setFecha(e.target.value)}
            disabled={enviando}
          />

          <div>
            <span className="etiqueta">Franja horaria</span>
            <div className="mt-2 space-y-2">
              {franjas.map((f) => (
                <label
                  key={f.id}
                  className={`flex cursor-pointer items-center justify-between rounded-xl border p-3.5 transition ${
                    franja === f.id
                      ? 'border-acento-400/60 bg-acento-400/10 shadow-resaltado'
                      : 'border-white/8 bg-carbono-900/50 hover:border-white/20'
                  }`}
                >
                  <span className="flex items-center gap-3">
                    <span
                      className={`flex h-4 w-4 items-center justify-center rounded-full border ${
                        franja === f.id ? 'border-acento-400' : 'border-plata-500'
                      }`}
                    >
                      {franja === f.id && <span className="h-2 w-2 rounded-full bg-acento-400" />}
                    </span>
                    <span className="font-display text-sm font-medium text-plata-100 capitalize">{f.id}</span>
                  </span>
                  <span className="texto-numerico text-sm text-plata-400">
                    {f.inicio} – {f.fin}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {errorEnvio && <p className="text-sm text-red-400">{errorEnvio}</p>}

          <Boton
            tamano="lg"
            className="w-full"
            onClick={handleSolicitar}
            disabled={enviando || !fecha || !franja}
          >
            {enviando ? 'Solicitando…' : 'Solicitar test drive'}
          </Boton>
        </div>
      </div>
    </div>
  )
}
