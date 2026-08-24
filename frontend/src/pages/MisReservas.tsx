import { useEffect, useState, useCallback, useRef } from 'react'
import { Link } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { DatosTransferencia, Reserva } from '../types/reserva'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
import { SubirComprobante } from '../components/SubirComprobante'
import { DatosSenia } from '../components/reserva/DatosSenia'
import { formatearPrecio } from '../utils/formato'
import {
  estilosEstadoReserva,
  etiquetasEstadoReserva,
  EtiquetaEstado,
} from '../components/ui/EtiquetaEstado'

function formatearFecha(fecha: string): string {
  return new Date(fecha).toLocaleDateString('es-AR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

function formatearHora(fecha: string): string {
  return new Date(fecha).toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' })
}

// EstadoComprobanteResume muestra el plazo pendiente o el envío hecho.
function EstadoComprobanteResume({ reserva }: { reserva: Reserva }) {
  if (reserva.comprobanteEnviadoAt) {
    return (
      <p className="text-sm text-emerald-300">
        Comprobante enviado a las {formatearHora(reserva.comprobanteEnviadoAt)}
      </p>
    )
  }
  const vencimiento = reserva.vencimientoComprobante
  return (
    <p className="text-sm text-amber-300">
      Seña pendiente{vencimiento ? ` · subí el comprobante antes de las ${formatearHora(vencimiento)}` : ''}
    </p>
  )
}

export function MisReservas() {
  const [reservas, setReservas] = useState<Reserva[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [cancelandoId, setCancelandoId] = useState<number | null>(null)
  // Datos de transferencia por vehículo (para las reservas activas).
  const [datosPorVehiculo, setDatosPorVehiculo] = useState<Record<number, DatosTransferencia | null>>({})
  const vehiculosPedidosRef = useRef(new Set<number>())

  const cargarReservas = useCallback(async () => {
    try {
      const datos = await api.listarMisReservas()
      setReservas(datos)
      setError(null)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'Error al cargar las reservas')
    } finally {
      setCargando(false)
    }
  }, [])

  useEffect(() => {
    cargarReservas()
  }, [cargarReservas])

  // Para cada reserva activa trae los datos bancarios una sola vez por
  // vehículo: así el cliente siempre ve dónde transferir la seña.
  useEffect(() => {
    reservas
      .filter((reserva) => reserva.estado === 'activa')
      .forEach(({ vehiculo }) => {
        if (vehiculosPedidosRef.current.has(vehiculo.id)) return
        vehiculosPedidosRef.current.add(vehiculo.id)

        api
          .obtenerDatosTransferencia(vehiculo.id)
          .then((datos) => setDatosPorVehiculo((previo) => ({ ...previo, [vehiculo.id]: datos })))
          .catch(() => setDatosPorVehiculo((previo) => ({ ...previo, [vehiculo.id]: null })))
      })
  }, [reservas])

  const handleCancelar = async (reserva: Reserva) => {
    if (!window.confirm(`¿Cancelar la reserva de ${reserva.vehiculo.marca} ${reserva.vehiculo.modelo}?`)) {
      return
    }

    setCancelandoId(reserva.id)
    try {
      await api.cancelarReserva(reserva.id)
      await cargarReservas()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo cancelar la reserva')
    } finally {
      setCancelandoId(null)
    }
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando reservas…" />
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <EncabezadoPagina
        destacado="Cliente"
        titulo="Mis reservas"
        descripcion="Las unidades que reservaste y su estado."
        acciones={
          <Boton variante="secundario" tamano="sm">
            <Link to="/catalogo">Explorar catálogo</Link>
          </Boton>
        }
      />

      {error && (
        <div className="mt-6 rounded-xl border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-300">
          {error}
        </div>
      )}

      {reservas.length === 0 ? (
        <div className="mt-8">
          <EstadoVacio
            titulo="No tenés reservas"
            descripcion="Explorá el catálogo y asegurá la unidad que te gusta reservándola."
            accion={
              <Boton>
                <Link to="/catalogo">Explorar catálogo</Link>
              </Boton>
            }
          />
        </div>
      ) : (
        <div className="mt-8 space-y-4">
          {reservas.map((reserva) => (
            <div
              key={reserva.id}
              className="flex flex-col gap-4 rounded-2xl border border-white/8 bg-carbono-850/60 p-5 shadow-luz backdrop-blur-sm"
            >
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex items-center gap-4">
                {reserva.vehiculo.imagen && (
                  <img
                    src={reserva.vehiculo.imagen}
                    alt={`${reserva.vehiculo.marca} ${reserva.vehiculo.modelo}`}
                    className="h-16 w-24 rounded-lg object-cover"
                  />
                )}
                <div>
                  <Link
                    to={`/catalogo/${reserva.vehiculo.id}`}
                    className="font-display text-lg font-semibold text-plata-100 underline-offset-4 hover:underline"
                  >
                    {reserva.vehiculo.marca} {reserva.vehiculo.modelo}
                  </Link>
                  <p className="mt-1 text-sm text-plata-400">
                    Reservada el {formatearFecha(reserva.createdAt)} · Seña:{' '}
                    <span className="texto-numerico text-plata-300">{formatearPrecio(reserva.montoSenia)}</span>
                  </p>
                  <div className="mt-2 flex flex-wrap items-center gap-3">
                    <EtiquetaEstado
                      estado={reserva.estado}
                      estilos={estilosEstadoReserva}
                      etiqueta={etiquetasEstadoReserva[reserva.estado]}
                    />
                    {reserva.estado === 'activa' && <EstadoComprobanteResume reserva={reserva} />}
                  </div>
                  {reserva.estado === 'cancelada' && reserva.motivoCancelacion && (
                    <p className="mt-2 max-w-xl rounded-lg border border-red-500/25 bg-red-500/10 px-3 py-2 text-sm text-red-200">
                      La concesionaria canceló esta reserva: {reserva.motivoCancelacion}
                    </p>
                  )}
                </div>
              </div>

              {reserva.estado === 'activa' && (
                <div className="flex flex-col items-start gap-3 sm:items-end">
                  {!reserva.comprobanteEnviadoAt && (
                    <SubirComprobante reservaId={reserva.id} alEnviar={cargarReservas} />
                  )}
                  <Boton
                    variante="peligro"
                    tamano="sm"
                    onClick={() => handleCancelar(reserva)}
                    disabled={cancelandoId === reserva.id}
                  >
                    {cancelandoId === reserva.id ? 'Cancelando…' : 'Cancelar reserva'}
                  </Boton>
                </div>
              )}
              </div>

              {reserva.estado === 'activa' && (
                <DatosSenia datos={datosPorVehiculo[reserva.vehiculo.id] ?? null} />
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
