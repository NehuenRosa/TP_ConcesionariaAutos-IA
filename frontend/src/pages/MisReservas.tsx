import { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { Reserva } from '../types/reserva'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
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

export function MisReservas() {
  const [reservas, setReservas] = useState<Reserva[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [cancelandoId, setCancelandoId] = useState<number | null>(null)

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
              className="flex flex-col gap-4 rounded-2xl border border-white/8 bg-carbono-850/60 p-5 shadow-luz backdrop-blur-sm sm:flex-row sm:items-center sm:justify-between"
            >
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
                  <p className="mt-1 text-sm text-plata-400">Reservada el {formatearFecha(reserva.createdAt)}</p>
                  <div className="mt-2">
                    <EtiquetaEstado
                      estado={reserva.estado}
                      estilos={estilosEstadoReserva}
                      etiqueta={etiquetasEstadoReserva[reserva.estado]}
                    />
                  </div>
                </div>
              </div>

              {reserva.estado === 'activa' && (
                <Boton
                  variante="peligro"
                  tamano="sm"
                  onClick={() => handleCancelar(reserva)}
                  disabled={cancelandoId === reserva.id}
                >
                  {cancelandoId === reserva.id ? 'Cancelando…' : 'Cancelar reserva'}
                </Boton>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
