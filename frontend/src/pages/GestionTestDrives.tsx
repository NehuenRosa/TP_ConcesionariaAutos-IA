import { useEffect, useState, useCallback } from 'react'
import { api, ErrorApi } from '../services/api'
import type { EstadoTurnoTestDrive, TurnoTestDrive } from '../types/testDrive'
import { Boton } from '../components/ui/Boton'
import { CampoSeleccion } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
import {
  estilosEstadoTestDrive,
  etiquetasEstadoTestDrive,
  EtiquetaEstado,
} from '../components/ui/EtiquetaEstado'

const ESTADOS: Array<{ valor: EstadoTurnoTestDrive | ''; etiqueta: string }> = [
  { valor: '', etiqueta: 'Todos' },
  { valor: 'solicitado', etiqueta: 'Solicitados' },
  { valor: 'confirmado', etiqueta: 'Confirmados' },
  { valor: 'cancelado', etiqueta: 'Cancelados' },
  { valor: 'completado', etiqueta: 'Completados' },
]

function formatearFecha(fecha: string): string {
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-AR', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

export function GestionTestDrives() {
  const [turnos, setTurnos] = useState<TurnoTestDrive[]>([])
  const [filtro, setFiltro] = useState<EstadoTurnoTestDrive | ''>('')
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [accionId, setAccionId] = useState<number | null>(null)

  const cargarTurnos = useCallback(async () => {
    try {
      const datos = await api.listarTestDrives(filtro || undefined)
      setTurnos(datos)
      setError(null)
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'Error al cargar los test drives')
    } finally {
      setCargando(false)
    }
  }, [filtro])

  useEffect(() => {
    setCargando(true)
    cargarTurnos()
  }, [cargarTurnos])

  const ejecutarAccion = async (turno: TurnoTestDrive, accion: 'confirmar' | 'cancelar' | 'completar') => {
    if (accion === 'cancelar' && !window.confirm(`¿Cancelar el turno de ${turno.cliente.nombre}?`)) {
      return
    }

    setAccionId(turno.id)
    try {
      if (accion === 'confirmar') {
        await api.confirmarTestDrive(turno.id)
      } else if (accion === 'cancelar') {
        await api.cancelarTestDriveVendedor(turno.id)
      } else {
        await api.completarTestDrive(turno.id)
      }
      await cargarTurnos()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo actualizar el turno')
    } finally {
      setAccionId(null)
    }
  }

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando test drives…" />
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <EncabezadoPagina
        destacado="Vendedor"
        titulo="Test drives"
        descripcion="Confirmá, completá o cancelá los turnos de prueba de manejo."
        acciones={
          <div className="w-48">
            <CampoSeleccion
              value={filtro}
              onChange={(e) => setFiltro(e.target.value as EstadoTurnoTestDrive | '')}
            >
              {ESTADOS.map((estado) => (
                <option key={estado.valor || 'todos'} value={estado.valor}>
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

      {turnos.length === 0 ? (
        <div className="mt-8">
          <EstadoVacio
            titulo="No hay test drives para mostrar"
            descripcion="Cuando un cliente solicite un turno, vas a poder gestionarlo acá."
          />
        </div>
      ) : (
        <div className="mt-8 space-y-4">
          {turnos.map((turno) => (
            <div
              key={turno.id}
              className="flex flex-col gap-4 rounded-2xl border border-white/8 bg-carbono-850/60 p-5 shadow-luz backdrop-blur-sm lg:flex-row lg:items-center lg:justify-between"
            >
              <div className="space-y-2">
                <div className="flex flex-wrap items-center gap-3">
                  <p className="font-display text-lg font-semibold text-plata-100">
                    {turno.vehiculo.marca} {turno.vehiculo.modelo}
                  </p>
                  <EtiquetaEstado
                    estado={turno.estado}
                    estilos={estilosEstadoTestDrive}
                    etiqueta={etiquetasEstadoTestDrive[turno.estado]}
                  />
                </div>
                <p className="text-sm text-plata-400">
                  Cliente: <span className="text-plata-300">{turno.cliente.nombre}</span>
                </p>
                <p className="text-sm text-plata-400">
                  {formatearFecha(turno.fecha)} ·{' '}
                  <span className="text-plata-300 capitalize">{turno.franja}</span>
                </p>
              </div>

              <div className="flex shrink-0 flex-wrap gap-2">
                {turno.estado === 'solicitado' && (
                  <Boton tamano="sm" onClick={() => ejecutarAccion(turno, 'confirmar')} disabled={accionId === turno.id}>
                    {accionId === turno.id ? '…' : 'Confirmar'}
                  </Boton>
                )}
                {turno.estado === 'confirmado' && (
                  <Boton tamano="sm" onClick={() => ejecutarAccion(turno, 'completar')} disabled={accionId === turno.id}>
                    {accionId === turno.id ? '…' : 'Completar'}
                  </Boton>
                )}
                {(turno.estado === 'solicitado' || turno.estado === 'confirmado') && (
                  <Boton
                    variante="peligro"
                    tamano="sm"
                    onClick={() => ejecutarAccion(turno, 'cancelar')}
                    disabled={accionId === turno.id}
                  >
                    {accionId === turno.id ? '…' : 'Cancelar'}
                  </Boton>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
