import { useEffect, useRef, useState } from 'react'
import { Link, useParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { Vehiculo } from '../types/vehiculo'
import type { FranjaHoraria, TurnoTestDrive } from '../types/testDrive'
import { Boton } from '../components/ui/Boton'
import { CampoTexto } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { formatearPrecio } from '../utils/formato'

function hoyISO(): string {
  const ahora = new Date()
  const mes = String(ahora.getMonth() + 1).padStart(2, '0')
  const dia = String(ahora.getDate()).padStart(2, '0')
  return `${ahora.getFullYear()}-${mes}-${dia}`
}

// franjaEnPasado indica si la franja ya comenzó cuando se pide para hoy.
function franjaEnPasado(fecha: string, franja: FranjaHoraria): boolean {
  if (fecha !== hoyISO()) return false
  const [hora, minuto] = franja.inicio.split(':').map(Number)
  const comienzo = new Date()
  comienzo.setHours(hora, minuto, 0, 0)
  return comienzo.getTime() <= Date.now()
}

// turnoPendienteDelVehiculo encuentra un turno activo (solicitado o
// confirmado) del cliente para el vehículo dado.
function turnoPendienteDelVehiculo(turnos: TurnoTestDrive[], vehiculoId: number): TurnoTestDrive | null {
  return (
    turnos.find(
      (t) =>
        t.vehiculo.id === vehiculoId &&
        (t.estado === 'solicitado' || t.estado === 'confirmado'),
    ) ?? null
  )
}

export function FormularioTestDrive() {
  const { id } = useParams<{ id: string }>()
  const [vehiculo, setVehiculo] = useState<Vehiculo | null>(null)
  const [franjas, setFranjas] = useState<FranjaHoraria[]>([])
  const [fecha, setFecha] = useState(hoyISO())
  const [franja, setFranja] = useState('')
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [errorFranjas, setErrorFranjas] = useState<string | null>(null)
  const [turnoPendiente, setTurnoPendiente] = useState<TurnoTestDrive | null>(null)
  const [enviando, setEnviando] = useState(false)
  const [errorEnvio, setErrorEnvio] = useState<string | null>(null)
  const [turnoCreado, setTurnoCreado] = useState(false)
  const enviandoRef = useRef(false)

  useEffect(() => {
    if (!id) return
    let cancelado = false

    Promise.all([
      api.obtenerVehiculo(Number(id)),
      // Si el pedido de turnos propios falla (ej. visitante sin sesión), se
      // ignora: la validación definitiva la hace el backend al enviar.
      api.listarMisTestDrives().catch(() => [] as TurnoTestDrive[]),
    ])
      .then(([dato, misTurnos]) => {
        if (cancelado) return
        setVehiculo(dato)
        setTurnoPendiente(turnoPendienteDelVehiculo(misTurnos, dato.id))
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

  useEffect(() => {
    if (!id || !fecha) return
    let cancelado = false

    api
      .obtenerFranjas(Number(id), fecha)
      .then((disponibles) => {
        if (cancelado) return
        setFranjas(disponibles)
        setErrorFranjas(null)
        setFranja((actual) => (actual && disponibles.find((f) => f.id === actual)?.ocupada ? '' : actual))
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setErrorFranjas(e instanceof ErrorApi ? e.message : 'No se pudieron obtener los horarios del test drive.')
      })

    return () => {
      cancelado = true
    }
  }, [id, fecha])

  const handleSolicitar = async () => {
    if (!vehiculo || !fecha || !franja || enviandoRef.current) return
    enviandoRef.current = true
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
      enviandoRef.current = false
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
            para el {fecha} a las <span className="font-semibold text-plata-100">{franja} hs</span>. Un
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

          {turnoPendiente && (
            <div
              role="alert"
              className="rounded-xl border border-amber-500/40 bg-amber-500/10 p-3 text-sm text-amber-200"
            >
              Ya tenés un test drive pendiente para este vehículo (
              {turnoPendiente.fecha} a las {turnoPendiente.franja} hs). Podés cancelarlo desde{' '}
              <Link to="/mis-test-drives" className="font-medium underline">
                Mis test drives
              </Link>
              .
            </div>
          )}

          <div>
            <span className="etiqueta">Hora del turno</span>
            <p className="mt-1 text-xs text-plata-500">
              Los turnos duran media hora; los horarios ya reservados aparecen en rojo.
            </p>
            <div className="mt-2 grid grid-cols-2 gap-2 sm:grid-cols-4">
              {franjas.map((f) => {
                const enPasado = franjaEnPasado(fecha, f)
                const ocupada = f.ocupada === true
                const seleccionada = franja === f.id
                return (
                  <button
                    key={f.id}
                    type="button"
                    onClick={() => {
                      setFranja(f.id)
                      setErrorEnvio(null)
                    }}
                    disabled={enPasado || ocupada || enviando}
                    aria-pressed={seleccionada}
                    aria-label={`${f.inicio} a ${f.fin}${ocupada ? ' (ocupado)' : ''}`}
                    className={`rounded-xl border px-3 py-3 text-center transition ${
                      enPasado
                        ? 'cursor-not-allowed border-white/5 bg-carbono-900/30 text-plata-600 opacity-50'
                        : ocupada
                          ? 'cursor-not-allowed border-rose-500/40 bg-rose-500/10 text-plata-500'
                          : seleccionada
                            ? 'border-acento-400/60 bg-acento-400/10 text-acento-300 shadow-resaltado'
                            : 'border-white/8 bg-carbono-900/50 text-plata-100 hover:border-white/20'
                    }`}
                  >
                    <span className="texto-numerico block font-display text-sm font-semibold">
                      {f.inicio}
                    </span>
                    {ocupada ? (
                      <span className="mt-0.5 block text-[11px] font-medium text-rose-400">Ocupado</span>
                    ) : (
                      <span className="mt-0.5 block text-[11px] text-plata-500">
                        {f.inicio} – {f.fin} hs
                      </span>
                    )}
                  </button>
                )
              })}
            </div>
            {franjas.some((f) => f.ocupada) && (
              <p className="mt-2 text-xs text-rose-300/80">
                Los horarios en rojo ya tienen un turno reservado para esa fecha.
              </p>
            )}
            {errorFranjas && <p className="mt-2 text-xs text-red-300">{errorFranjas}</p>}
            {franjas.length === 0 && (
              <p className="mt-2 text-sm text-plata-500">No hay horarios disponibles por el momento.</p>
            )}
          </div>

          {errorEnvio && (
            <p role="alert" className="rounded-xl border border-red-500/30 bg-red-500/10 p-3 text-sm text-red-300">
              {errorEnvio}
            </p>
          )}

          <Boton
            tamano="lg"
            className="w-full"
            onClick={handleSolicitar}
            disabled={enviando || !fecha || !franja || !!turnoPendiente}
          >
            {enviando ? 'Solicitando…' : turnoPendiente ? 'Ya tenés un turno pendiente' : 'Solicitar test drive'}
          </Boton>
          {turnoPendiente && (
            <p className="text-center text-xs text-plata-500">
              Esperá a que completen o cancelen tu turno vigente para pedir otro.
            </p>
          )}
          {!franja && (
            <p className="text-center text-xs text-plata-500">
              Seleccioná una hora para habilitar la solicitud.
            </p>
          )}
        </div>
      </div>
    </div>
  )
}
