import { useEffect, useState, useCallback } from 'react'
import { api, ErrorApi } from '../services/api'
import type { EstadoReserva, Reserva } from '../types/reserva'
import { Boton } from '../components/ui/Boton'
import { CampoSeleccion } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
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

export function GestionReservas() {
  const [reservas, setReservas] = useState<Reserva[]>([])
  const [filtro, setFiltro] = useState<EstadoReserva | ''>('')
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [accionId, setAccionId] = useState<number | null>(null)

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
    if (accion === 'cancelar' && !window.confirm(`¿Cancelar la reserva de ${reserva.cliente.nombre}?`)) {
      return
    }

    setAccionId(reserva.id)
    try {
      if (accion === 'confirmar') {
        await api.confirmarReservaVenta(reserva.id)
      } else {
        await api.cancelarReservaVendedor(reserva.id)
      }
      await cargarReservas()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo actualizar la reserva')
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
                </div>
                <p className="text-sm text-plata-400">
                  Cliente: <span className="text-plata-300">{reserva.cliente.nombre}</span>
                </p>
                <p className="text-sm text-plata-400">Reservada el {formatearFecha(reserva.createdAt)}</p>
              </div>

              <div className="flex shrink-0 flex-wrap gap-2">
                {reserva.estado === 'activa' && (
                  <>
                    <Boton tamano="sm" onClick={() => ejecutarAccion(reserva, 'confirmar')} disabled={accionId === reserva.id}>
                      {accionId === reserva.id ? '…' : 'Confirmar venta'}
                    </Boton>
                    <Boton
                      variante="peligro"
                      tamano="sm"
                      onClick={() => ejecutarAccion(reserva, 'cancelar')}
                      disabled={accionId === reserva.id}
                    >
                      {accionId === reserva.id ? '…' : 'Cancelar'}
                    </Boton>
                  </>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
