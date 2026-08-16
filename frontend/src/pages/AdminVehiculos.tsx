import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { ResumenVehiculoGestion, EstadoVehiculo } from '../types/vehiculo'
import { Boton } from '../components/ui/Boton'
import { CampoSeleccion } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'
import { Paginacion } from '../components/ui/Paginacion'
import {
  estilosEstadoVehiculo,
  etiquetasEstadoVehiculo,
  EtiquetaEstado,
} from '../components/ui/EtiquetaEstado'
import { formatearKilometraje, formatearPrecio } from '../utils/formato'

const TAMANO_PAGINA = 10

const ESTADOS: EstadoVehiculo[] = ['disponible', 'reservado', 'vendido', 'dado_de_baja']

export function AdminVehiculos() {
  const navigate = useNavigate()
  const [vehiculos, setVehiculos] = useState<ResumenVehiculoGestion[]>([])
  const [filtroEstado, setFiltroEstado] = useState('')
  const [pagina, setPagina] = useState(1)
  const [total, setTotal] = useState(0)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelado = false
    setCargando(true)
    setError(null)

    api
      .listarVehiculosGestion(pagina, TAMANO_PAGINA, filtroEstado || undefined)
      .then((respuesta) => {
        if (cancelado) return
        setVehiculos(respuesta.datos)
        setTotal(respuesta.total)
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'Ocurrió un error inesperado al cargar el stock.')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [pagina, filtroEstado])

  const totalPaginas = Math.max(1, Math.ceil(total / TAMANO_PAGINA))

  async function darDeBaja(id: number) {
    const vehiculo = vehiculos.find((v) => v.id === id)
    const confirmacion = window.confirm(
      `¿Dar de baja ${vehiculo?.marca ?? ''} ${vehiculo?.modelo ?? ''}? Esta acción oculta el vehículo del catálogo público.`,
    )
    if (!confirmacion) return

    try {
      await api.darDeBajaVehiculo(id)
      setError(null)
      if (filtroEstado && filtroEstado !== 'dado_de_baja') {
        // Con un filtro activo que ya no matchea, la fila debe desaparecer.
        setVehiculos((actuales) => actuales.filter((v) => v.id !== id))
      } else {
        setVehiculos((actuales) =>
          actuales.map((v) => (v.id === id ? { ...v, estado: 'dado_de_baja' } : v)),
        )
      }
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo dar de baja el vehículo.')
    }
  }

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <EncabezadoPagina
        destacado="Administración"
        titulo="Gestión de vehículos"
        descripcion="Administrá el stock del concesionario."
        acciones={
          <Boton>
            <Link to="/admin/vehiculos/nuevo">+ Nuevo vehículo</Link>
          </Boton>
        }
      />

      <div className="mt-8 flex flex-wrap items-center gap-4">
        <div className="w-56">
          <CampoSeleccion
            id="filtro-estado"
            etiqueta="Filtrar por estado"
            value={filtroEstado}
            onChange={(e) => {
              setPagina(1)
              setFiltroEstado(e.target.value)
            }}
          >
            <option value="">Todos</option>
            {ESTADOS.map((estado) => (
              <option key={estado} value={estado}>
                {etiquetasEstadoVehiculo[estado]}
              </option>
            ))}
          </CampoSeleccion>
        </div>
      </div>

      {cargando && <ContenidoCargando etiqueta="Cargando vehículos…" />}

      {!cargando && error && (
        <div className="mt-6">
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-300">{error}</div>
        </div>
      )}

      {!cargando && !error && vehiculos.length === 0 && (
        <div className="mt-8">
          <EstadoVacio
            titulo="No hay vehículos para mostrar"
            descripcion="Agregá tu primera unidad al stock para comenzar a vender."
            accion={
              <Boton>
                <Link to="/admin/vehiculos/nuevo">+ Nuevo vehículo</Link>
              </Boton>
            }
          />
        </div>
      )}

      {!cargando && !error && vehiculos.length > 0 && (
        <>
          <div className="mt-6 overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/50">
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead className="border-b border-white/8 bg-carbono-900/70 text-plata-500">
                  <tr>
                    <th className="px-5 py-3.5 font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Vehículo</th>
                    <th className="px-5 py-3.5 font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Año</th>
                    <th className="px-5 py-3.5 font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Km</th>
                    <th className="px-5 py-3.5 font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Precio</th>
                    <th className="px-5 py-3.5 font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Estado</th>
                    <th className="px-5 py-3.5 text-right font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Acciones</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/6">
                  {vehiculos.map((vehiculo) => (
                    <tr key={vehiculo.id} className="transition-colors hover:bg-white/3">
                      <td className="px-5 py-3.5">
                        <div className="flex items-center gap-3">
                          <div className="h-12 w-16 shrink-0 overflow-hidden rounded-lg border border-white/8 bg-carbono-900">
                            {vehiculo.imagen ? (
                              <img
                                src={vehiculo.imagen}
                                alt={`${vehiculo.marca} ${vehiculo.modelo}`}
                                className="h-full w-full object-cover"
                              />
                            ) : (
                              <span className="flex h-full w-full items-center justify-center text-xs text-plata-500">
                                Sin imagen
                              </span>
                            )}
                          </div>
                          <div>
                            <p className="font-display font-semibold text-plata-100">
                              {vehiculo.marca} {vehiculo.modelo}
                            </p>
                            <p className="text-xs text-plata-500">
                              {vehiculo.condicion === 'nuevo' ? 'Nuevo' : 'Usado'}
                            </p>
                          </div>
                        </div>
                      </td>
                      <td className="px-5 py-3.5 text-plata-300">{vehiculo.anio}</td>
                      <td className="px-5 py-3.5 text-plata-300">{formatearKilometraje(vehiculo.kilometraje)}</td>
                      <td className="texto-numerico px-5 py-3.5 font-medium text-plata-100">
                        {formatearPrecio(vehiculo.precio)}
                      </td>
                      <td className="px-5 py-3.5">
                        <EtiquetaEstado
                          estado={vehiculo.estado}
                          estilos={estilosEstadoVehiculo}
                          etiqueta={etiquetasEstadoVehiculo[vehiculo.estado]}
                        />
                      </td>
                      <td className="px-5 py-3.5">
                        <div className="flex items-center justify-end gap-2">
                          <Boton
                            variante="secundario"
                            tamano="sm"
                            onClick={() => navigate(`/admin/vehiculos/${vehiculo.id}/editar`)}
                          >
                            Editar
                          </Boton>
                          {vehiculo.estado !== 'dado_de_baja' && (
                            <Boton variante="peligro" tamano="sm" onClick={() => darDeBaja(vehiculo.id)}>
                              Dar de baja
                            </Boton>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          <div className="mt-8">
            <Paginacion pagina={pagina} totalPaginas={totalPaginas} cambiarPagina={setPagina} />
          </div>
        </>
      )}
    </div>
  )
}
