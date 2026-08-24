import { useEffect, useRef, useState } from 'react'
import { Link, useParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { DatosTransferencia, Reserva } from '../types/reserva'
import type { Vehiculo } from '../types/vehiculo'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'
import { SubirComprobante } from '../components/SubirComprobante'
import { etiquetaCondicion, formatearPrecio } from '../utils/formato'

// formatearRestante convierte los milisegundos restantes en "1:23:45".
function formatearRestante(milisegundos: number): string {
  const totalSegundos = Math.max(0, Math.floor(milisegundos / 1000))
  const horas = Math.floor(totalSegundos / 3600)
  const minutos = Math.floor((totalSegundos % 3600) / 60)
  const segundos = totalSegundos % 60
  return `${horas}:${String(minutos).padStart(2, '0')}:${String(segundos).padStart(2, '0')}`
}

// DatosSenia vive en components/reserva: se comparte con Mis Reservas.
import { DatosSenia } from '../components/reserva/DatosSenia'

export function FormularioReserva() {
  const { id } = useParams<{ id: string }>()
  const [vehiculo, setVehiculo] = useState<Vehiculo | null>(null)
  const [datosSenia, setDatosSenia] = useState<DatosTransferencia | null>(null)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [enviando, setEnviando] = useState(false)
  const [errorEnvio, setErrorEnvio] = useState<string | null>(null)
  // Reserva creada: pasa a la vista de éxito con plazo y subida del comprobante.
  const [reserva, setReserva] = useState<Reserva | null>(null)
  const [ahora, setAhora] = useState(() => Date.now())
  const enviandoRef = useRef(false)

  useEffect(() => {
    if (!id) return
    let cancelado = false

    api
      .obtenerVehiculo(Number(id))
      .then((dato) => {
        if (cancelado) return
        setVehiculo(dato)
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'No se pudo cargar la información del vehículo.')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    api
      .obtenerDatosTransferencia(Number(id))
      .then((datos) => {
        if (!cancelado) setDatosSenia(datos)
      })
      .catch(() => {
        // Sin datos configurados o unidad no disponible: la UI muestra el
        // aviso genérico y el flujo sigue.
      })

    return () => {
      cancelado = true
    }
  }, [id])

  // Cuenta regresiva del plazo mientras haya una reserva activa pendiente.
  useEffect(() => {
    if (!reserva || reserva.comprobanteEnviadoAt || !reserva.vencimientoComprobante) return
    const intervalo = setInterval(() => setAhora(Date.now()), 1000)
    return () => clearInterval(intervalo)
  }, [reserva])

  const handleReservar = async () => {
    if (!vehiculo || enviandoRef.current) return

    // Guard síncrono: evita que dos clics rápidos disparen dos reservas.
    enviandoRef.current = true
    setEnviando(true)
    setErrorEnvio(null)

    try {
      const creada = await api.crearReserva({ vehiculoId: vehiculo.id })
      setReserva(creada)
      setAhora(Date.now())
    } catch (e: unknown) {
      setErrorEnvio(e instanceof ErrorApi ? e.message : 'No se pudo reservar el vehículo')
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
        <h1 className="font-display text-3xl font-bold text-plata-100">No se pudo reservar el vehículo</h1>
        <p className="mt-3 text-plata-400">{error ?? 'El vehículo solicitado no existe o no está disponible.'}</p>
        <div className="mt-8">
          <Boton variante="secundario">
            <Link to="/catalogo">Volver al catálogo</Link>
          </Boton>
        </div>
      </div>
    )
  }

  if (reserva) {
    const pendiente = !reserva.comprobanteEnviadoAt
    const vencimiento = reserva.vencimientoComprobante ? new Date(reserva.vencimientoComprobante).getTime() : null
    const restante = vencimiento !== null ? formatearRestante(vencimiento - ahora) : null
    const horaLimite =
      vencimiento !== null
        ? new Date(vencimiento).toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' })
        : null

    return (
      <div className="mx-auto max-w-xl px-4 py-20 sm:px-6">
        <div className="rounded-2xl border border-emerald-400/30 bg-emerald-400/10 p-8 shadow-luz">
          <div className="text-center">
            <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full border border-emerald-400/40 bg-emerald-400/15 text-emerald-300">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="h-7 w-7">
                <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h1 className="font-display text-2xl font-bold text-plata-100">Vehículo reservado</h1>
            <p className="mt-3 text-sm leading-relaxed text-plata-300">
              Reservaste la unidad{' '}
              <span className="font-semibold text-plata-100">
                {vehiculo.marca} {vehiculo.modelo}
              </span>
              . Ahora transferí la seña y subí el comprobante dentro del plazo.
            </p>
          </div>

          {pendiente && (
            <div className="mt-5 space-y-4">
              <div className="flex items-center justify-between rounded-xl border border-amber-400/40 bg-carbono-900/70 p-4">
                <div>
                  <p className="font-display text-sm font-semibold text-amber-200">Tiempo restante</p>
                  {horaLimite && (
                    <p className="text-xs text-plata-400">Subilo antes de las {horaLimite} o la reserva se anula.</p>
                  )}
                </div>
                {restante && (
                  <span className="texto-numerico font-display text-2xl font-bold text-amber-300">{restante}</span>
                )}
              </div>

              <DatosSenia datos={datosSenia} />

              <SubirComprobante reservaId={reserva.id} />
            </div>
          )}

          {!pendiente && (
            <div className="mt-5 space-y-4">
              <p className="flex items-center justify-center gap-2 rounded-xl border border-emerald-400/30 bg-carbono-900/60 p-4 text-sm text-emerald-300">
                Comprobante recibido. Un vendedor lo va a revisar y confirmar la venta.
              </p>
              <DatosSenia datos={datosSenia} />
            </div>
          )}

          <div className="mt-6 flex flex-wrap justify-center gap-3">
            <Boton>
              <Link to="/mis-reservas">Ver mis reservas</Link>
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
        <h1 className="font-display text-2xl font-bold text-plata-100">Reservar este vehículo</h1>
        <p className="mt-2 text-sm leading-relaxed text-plata-400">
          Al confirmar, la unidad queda reservada a tu nombre durante 2 horas: tenés ese plazo para transferir la
          seña (5 % del valor) y subir el comprobante. Si no lo hacés, la reserva se anula sola y la unidad vuelve al
          catálogo.
        </p>

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
              {vehiculo.anio} · {etiquetaCondicion(vehiculo.condicion)} ·{' '}
              <span className="texto-numerico text-plata-300">{formatearPrecio(vehiculo.precio)}</span>
            </p>
          </div>
        </div>

        <div className="mt-5">
          <DatosSenia datos={datosSenia} />
        </div>

        <div className="mt-6 space-y-5">
          {errorEnvio && <p className="text-sm text-red-400">{errorEnvio}</p>}

          <Boton tamano="lg" className="w-full" onClick={handleReservar} disabled={enviando}>
            {enviando ? 'Reservando…' : 'Confirmar reserva'}
          </Boton>
        </div>
      </div>
    </div>
  )
}
