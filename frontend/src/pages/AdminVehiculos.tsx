import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { ResumenVehiculoGestion, EstadoVehiculo } from '../types/vehiculo'

const TAMANO_PAGINA = 10

const ESTADOS: EstadoVehiculo[] = ['disponible', 'reservado', 'vendido', 'dado_de_baja']

const ETIQUETAS_ESTADO: Record<EstadoVehiculo, string> = {
  disponible: 'Disponible',
  reservado: 'Reservado',
  vendido: 'Vendido',
  dado_de_baja: 'Dado de baja',
}

function formatearPrecio(precio: number): string {
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
    maximumFractionDigits: 0,
  }).format(precio)
}

function colorEstado(estado: EstadoVehiculo): string {
  switch (estado) {
    case 'disponible':
      return 'bg-green-100 text-green-800'
    case 'reservado':
      return 'bg-amber-100 text-amber-800'
    case 'vendido':
      return 'bg-blue-100 text-blue-800'
    case 'dado_de_baja':
      return 'bg-gray-200 text-gray-600'
  }
}

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
      setVehiculos((actuales) =>
        actuales.map((v) => (v.id === id ? { ...v, estado: 'dado_de_baja' } : v)),
      )
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo dar de baja el vehículo.')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Gestión de vehículos</h1>
          <p className="mt-1 text-gray-700">Administrá el stock del concesionario.</p>
        </div>
        <Link
          to="/admin/vehiculos/nuevo"
          className="rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700"
        >
          + Nuevo vehículo
        </Link>
      </div>

      <div className="flex items-center gap-2">
        <label htmlFor="filtro-estado" className="text-sm text-gray-700">
          Estado:
        </label>
        <select
          id="filtro-estado"
          value={filtroEstado}
          onChange={(e) => {
            setPagina(1)
            setFiltroEstado(e.target.value)
          }}
          className="rounded-md border border-gray-300 px-3 py-1.5 text-sm text-gray-700"
        >
          <option value="">Todos</option>
          {ESTADOS.map((estado) => (
            <option key={estado} value={estado}>
              {ETIQUETAS_ESTADO[estado]}
            </option>
          ))}
        </select>
      </div>

      {cargando && <p className="text-gray-700">Cargando vehículos…</p>}

      {!cargando && error && <p className="text-red-600">{error}</p>}

      {!cargando && !error && vehiculos.length === 0 && (
        <p className="text-gray-700">No hay vehículos para mostrar.</p>
      )}

      {!cargando && !error && vehiculos.length > 0 && (
        <>
          <div className="overflow-x-auto rounded-lg border border-gray-200">
            <table className="w-full text-left text-sm text-gray-700">
              <thead className="border-b border-gray-200 bg-gray-50 text-gray-500">
                <tr>
                  <th className="px-4 py-3">Vehículo</th>
                  <th className="px-4 py-3">Año</th>
                  <th className="px-4 py-3">Km</th>
                  <th className="px-4 py-3">Precio</th>
                  <th className="px-4 py-3">Estado</th>
                  <th className="px-4 py-3 text-right">Acciones</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {vehiculos.map((vehiculo) => (
                  <tr key={vehiculo.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-3">
                        <div className="h-12 w-16 shrink-0 overflow-hidden rounded bg-gray-100">
                          {vehiculo.imagen ? (
                            <img
                              src={vehiculo.imagen}
                              alt={`${vehiculo.marca} ${vehiculo.modelo}`}
                              className="h-full w-full object-cover"
                            />
                          ) : (
                            <span className="flex h-full w-full items-center justify-center text-xs text-gray-400">
                              Sin imagen
                            </span>
                          )}
                        </div>
                        <div>
                          <p className="font-semibold text-gray-900">
                            {vehiculo.marca} {vehiculo.modelo}
                          </p>
                          <p className="text-xs text-gray-500">
                            {vehiculo.condicion === 'nuevo' ? 'Nuevo' : 'Usado'}
                          </p>
                        </div>
                      </div>
                    </td>
                    <td className="px-4 py-3">{vehiculo.anio}</td>
                    <td className="px-4 py-3">{new Intl.NumberFormat('es-AR').format(vehiculo.kilometraje)}</td>
                    <td className="px-4 py-3">{formatearPrecio(vehiculo.precio)}</td>
                    <td className="px-4 py-3">
                      <span
                        className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${colorEstado(vehiculo.estado)}`}
                      >
                        {ETIQUETAS_ESTADO[vehiculo.estado]}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <button
                          type="button"
                          onClick={() => navigate(`/admin/vehiculos/${vehiculo.id}/editar`)}
                          className="rounded-md border border-gray-300 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
                        >
                          Editar
                        </button>
                        {vehiculo.estado !== 'dado_de_baja' && (
                          <button
                            type="button"
                            onClick={() => darDeBaja(vehiculo.id)}
                            className="rounded-md border border-red-300 px-3 py-1.5 text-sm text-red-700 hover:bg-red-50"
                          >
                            Dar de baja
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <nav className="flex items-center justify-center gap-2" aria-label="Paginación">
            <button
              type="button"
              onClick={() => setPagina((p) => Math.max(1, p - 1))}
              disabled={pagina <= 1}
              className="rounded-md border border-gray-300 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
            >
              Anterior
            </button>
            <span className="px-2 text-sm text-gray-700">
              Página {pagina} de {totalPaginas}
            </span>
            <button
              type="button"
              onClick={() => setPagina((p) => Math.min(totalPaginas, p + 1))}
              disabled={pagina >= totalPaginas}
              className="rounded-md border border-gray-300 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
            >
              Siguiente
            </button>
          </nav>
        </>
      )}
    </div>
  )
}
