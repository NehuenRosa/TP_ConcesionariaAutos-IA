import { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { EstadoTurnoTestDrive, TurnoTestDrive } from '../types/testDrive'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
import {
  estilosEstadoTestDrive,
  etiquetasEstadoTestDrive,
  EtiquetaEstado,
} from '../components/ui/EtiquetaEstado'
import { formatearFranja } from '../utils/formato'

function formatearFecha(fecha: string): string {
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-AR', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

function esActivo(estado: EstadoTurnoTestDrive): boolean {
  return estado === 'solicitado' || estado === 'confirmado'
}

export function MisTestDrives() {
  const [turnos, setTurnos] = useState<TurnoTestDrive[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [cancelandoId, setCancelandoId] = useState<number | null>(null)

  const cargarTurnos = useCallback(async () => {
    try {
      const datos = await api.listarMisTestDrives()
      setTurnos(datos)
      setError(null)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'Error al cargar los test drives')
    } finally {
      setCargando(false)
    }
  }, [])

  useEffect(() => {
    cargarTurnos()
  }, [cargarTurnos])

  const handleCancelar = async (turno: TurnoTestDrive) => {
    if (!window.confirm(`¿Cancelar el test drive de ${turno.vehiculo.marca} ${turno.vehiculo.modelo}?`)) {
      return
    }

    setCancelandoId(turno.id)
    try {
      await api.cancelarTestDrive(turno.id)
      await cargarTurnos()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo cancelar el turno')
    } finally {
      setCancelandoId(null)
    }
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando turnos…" />
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <EncabezadoPagina
        destacado="Cliente"
        titulo="Mis test drives"
        descripcion="Tus turnos de prueba de manejo y su estado."
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

      {turnos.length === 0 ? (
        <div className="mt-8">
          <EstadoVacio
            titulo="No tenés test drives solicitados"
            descripcion="Explorá el catálogo y elegí un vehículo para probarlo en la ruta."
            accion={
              <Boton>
                <Link to="/catalogo">Explorar catálogo</Link>
              </Boton>
            }
          />
        </div>
      ) : (
        <div className="mt-8 space-y-4">
          {turnos.map((turno) => (
            <div
              key={turno.id}
              className="flex flex-col gap-4 rounded-2xl border border-white/8 bg-carbono-850/60 p-5 shadow-luz backdrop-blur-sm sm:flex-row sm:items-center sm:justify-between"
            >
              <div className="flex items-center gap-4">
                {turno.vehiculo.imagen && (
                  <img
                    src={turno.vehiculo.imagen}
                    alt={`${turno.vehiculo.marca} ${turno.vehiculo.modelo}`}
                    className="h-16 w-24 rounded-lg object-cover"
                  />
                )}
                <div>
                  <Link
                    to={`/catalogo/${turno.vehiculo.id}`}
                    className="font-display text-lg font-semibold text-plata-100 underline-offset-4 hover:underline"
                  >
                    {turno.vehiculo.marca} {turno.vehiculo.modelo}
                  </Link>
                  <p className="mt-1 text-sm text-plata-400">
                    {formatearFecha(turno.fecha)} ·{' '}
                    <span className="text-plata-300">{formatearFranja(turno.franja)}</span>
                  </p>
                  <div className="mt-2">
                    <EtiquetaEstado
                      estado={turno.estado}
                      estilos={estilosEstadoTestDrive}
                      etiqueta={etiquetasEstadoTestDrive[turno.estado]}
                    />
                  </div>
                </div>
              </div>

              {esActivo(turno.estado) && (
                <Boton
                  variante="peligro"
                  tamano="sm"
                  onClick={() => handleCancelar(turno)}
                  disabled={cancelandoId === turno.id}
                >
                  {cancelandoId === turno.id ? 'Cancelando…' : 'Cancelar turno'}
                </Boton>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
