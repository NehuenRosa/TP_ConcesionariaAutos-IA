import { useEffect, useState, useCallback } from 'react'
import { api, ErrorApi } from '../services/api'
import type { EstadoReserva, Reserva } from '../types/reserva'
import { Boton } from '../components/ui/Boton'
import { CampoSeleccion } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
import { formatearPrecio } from '../utils/formato'
import {
  estilosEstadoReserva,
  etiquetasEstadoReserva,
  EtiquetaEstado,
} from '../components/ui/EtiquetaEstado'

const ESTADOS: Array<{ valor: EstadoReserva | ''; etiqueta: string }> = [
  { valor: '', etiqueta: 'Todas' },
  { valor: 'activa', etiqueta: 'Activas' },
  { valor: 'vendida', etiqueta: 'Vendidas' },
  { valor: 'cancelada', etiqueta: 'Canceladas' },
]

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

export function GestionReservas() {
  const [reservas, setReservas] = useState<Reserva[]>([])
  const [filtro, setFiltro] = useState<EstadoReserva | ''>('')
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [accionId, setAccionId] = useState<number | null>(null)
  // Reserva en proceso de cancelación con su motivo obligatorio.
  const [cancelandoId, setCancelandoId] = useState<number | null>(null)
  const [motivo, setMotivo] = useState('')
  const [errorMotivo, setErrorMotivo] = useState<string | null>(null)

  // verComprobante descarga la imagen con el token y la abre en una pestaña.
  const verComprobante = async (id: number) => {
    setError(null)
    try {
      const imagen = await api.obtenerComprobanteReserva(id)
      const url = URL.createObjectURL(imagen)
      window.open(url, '_blank', 'noopener')
      window.setTimeout(() => URL.revokeObjectURL(url), 60_000)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo abrir el comprobante')
    }
  }

  const cargarReservas = useCallback(
    async (filtroSolicitado: EstadoReserva | '' = filtro) => {
      try {
        const datos = await api.listarReservas(filtroSolicitado || undefined)
        // Ignorar respuestas obsoletas si el filtro cambió mientras volvía.
        if (filtroSolicitado !== filtro) return
        setReservas(datos)
        setError(null)
      } catch (e: unknown) {
        if (filtroSolicitado !== filtro) return
        setError(e instanceof ErrorApi ? e.message : 'Error al cargar las reservas')
      } finally {
        if (filtroSolicitado === filtro) setCargando(false)
      }
    },
    [filtro],
  )

  useEffect(() => {
    setCargando(true)
    void cargarReservas(filtro)
  }, [cargarReservas, filtro])

  const ejecutarAccion = async (reserva: Reserva, accion: 'confirmar' | 'cancelar') => {
    if (accion === 'confirmar') {
      setAccionId(reserva.id)
      try {
        await api.confirmarReservaVenta(reserva.id)
        await cargarReservas()
      } catch (e: unknown) {
        setError(e instanceof ErrorApi ? e.message : 'No se pudo actualizar la reserva')
      } finally {
        setAccionId(null)
      }
      return
    }

    // Cancelación: exige el motivo que va a leer el cliente.
    const texto = motivo.trim()
    if (texto.length < 5) {
      setErrorMotivo('Explicá brevemente por qué se cancela la reserva (mínimo 5 caracteres)')
      return
    }
    setErrorMotivo(null)
    setAccionId(reserva.id)
    try {
      await api.cancelarReservaVendedor(reserva.id, texto)
      setCancelandoId(null)
      setMotivo('')
      await cargarReservas()
    } catch (e: unknown) {
      setErrorMotivo(e instanceof ErrorApi ? e.message : 'No se pudo cancelar la reserva')
    } finally {
      setAccionId(null)
    }
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando reservas…" />
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <EncabezadoPagina
        destacado="Vendedor"
        titulo="Reservas"
        descripcion="Confirmá la venta o cancelá las reservas de tus clientes."
        acciones={
          <div className="w-48">
            <CampoSeleccion
              value={filtro}
              onChange={(e) => setFiltro(e.target.value as EstadoReserva | '')}
            >
              {ESTADOS.map((estado) => (
                <option key={estado.valor || 'todas'} value={estado.valor}>
                  {estado.etiqueta}
                </option>
              ))}
            </CampoSeleccion>
          </div>
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
            titulo="No hay reservas para mostrar"
            descripcion="Cuando un cliente reserve una unidad, vas a poder gestionarla acá."
          />
        </div>
      ) : (
        <div className="mt-8 space-y-4">
          {reservas.map((reserva) => (
            <div
              key={reserva.id}
              className="flex flex-col gap-4 rounded-2xl border border-white/8 bg-carbono-850/60 p-5 shadow-luz backdrop-blur-sm lg:flex-row lg:items-center lg:justify-between"
            >
              <div className="space-y-2">
                <div className="flex flex-wrap items-center gap-3">
                  <p className="font-display text-lg font-semibold text-plata-100">
                    {reserva.vehiculo.marca} {reserva.vehiculo.modelo}
                  </p>
                  <EtiquetaEstado
                    estado={reserva.estado}
                    estilos={estilosEstadoReserva}
                    etiqueta={etiquetasEstadoReserva[reserva.estado]}
                  />
                  {reserva.comprobanteEnviadoAt ? (
                    <span className="rounded-full border border-emerald-400/30 bg-emerald-400/10 px-3 py-1 text-xs font-medium text-emerald-300">
                      Comprobante enviado {formatearHora(reserva.comprobanteEnviadoAt)}
                    </span>
                  ) : (
                    reserva.estado === 'activa' &&
                    reserva.vencimientoComprobante && (
                      <span className="rounded-full border border-amber-400/30 bg-amber-400/10 px-3 py-1 text-xs font-medium text-amber-300">
                        Pendiente de comprobante · vence {formatearHora(reserva.vencimientoComprobante)}
                      </span>
                    )
                  )}
                </div>
                <p className="text-sm text-plata-400">
                  Cliente: <span className="text-plata-300">{reserva.cliente.nombre}</span> · Seña:{' '}
                  <span className="texto-numerico text-plata-300">{formatearPrecio(reserva.montoSenia)}</span>
                </p>
                <p className="text-sm text-plata-400">Reservada el {formatearFecha(reserva.createdAt)}</p>
              </div>

              {reserva.estado === 'activa' && cancelandoId === reserva.id ? (
                <div className="w-full space-y-2 lg:w-96">
                  <label htmlFor={`motivo-${reserva.id}`} className="block text-xs font-medium text-plata-300">
                    Motivo de la cancelación (el cliente lo va a ver)
                  </label>
                  <textarea
                    id={`motivo-${reserva.id}`}
                    value={motivo}
                    onChange={(e) => setMotivo(e.target.value)}
                    rows={3}
                    placeholder="Ej.: el comprobante es ilegible o el monto no coincide con la seña…"
                    className="w-full resize-none rounded-lg border border-white/10 bg-carbono-900 px-3 py-2 text-sm text-plata-100 placeholder:text-plata-600 focus:border-acento-400 focus:outline-none"
                  />
                  {errorMotivo && <p className="text-xs text-red-300">{errorMotivo}</p>}
                  <div className="flex justify-end gap-2">
                    <Boton
                      variante="fantasma"
                      tamano="sm"
                      onClick={() => {
                        setCancelandoId(null)
                        setMotivo('')
                        setErrorMotivo(null)
                      }}
                      disabled={accionId === reserva.id}
                    >
                      Volver
                    </Boton>
                    <Boton
                      variante="peligro"
                      tamano="sm"
                      onClick={() => ejecutarAccion(reserva, 'cancelar')}
                      disabled={accionId === reserva.id || motivo.trim().length === 0}
                    >
                      {accionId === reserva.id ? 'Cancelando…' : 'Confirmar rechazo'}
                    </Boton>
                  </div>
                </div>
              ) : (
                <div className="flex shrink-0 flex-wrap gap-2">
                  {reserva.comprobanteEnviadoAt && (
                    <Boton variante="secundario" tamano="sm" onClick={() => verComprobante(reserva.id)}>
                      Ver comprobante
                    </Boton>
                  )}
                  {reserva.estado === 'activa' && (
                    <>
                      <Boton tamano="sm" onClick={() => ejecutarAccion(reserva, 'confirmar')} disabled={accionId === reserva.id}>
                        {accionId === reserva.id ? '…' : 'Confirmar venta'}
                      </Boton>
                      <Boton
                        variante="peligro"
                        tamano="sm"
                        onClick={() => {
                          setCancelandoId(reserva.id)
                          setMotivo('')
                          setErrorMotivo(null)
                        }}
                        disabled={accionId === reserva.id}
                      >
                        Rechazar / cancelar
                      </Boton>
                    </>
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
