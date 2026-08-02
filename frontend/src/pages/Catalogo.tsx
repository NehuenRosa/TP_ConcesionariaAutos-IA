import { useEffect, useState } from 'react'
import { Link } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { ResumenVehiculo } from '../types/vehiculo'

const TAMANO_PAGINA = 12

function formatearPrecio(precio: number): string {
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
    maximumFractionDigits: 0,
  }).format(precio)
}

export function Catalogo() {
  const [vehiculos, setVehiculos] = useState<ResumenVehiculo[]>([])
  const [pagina, setPagina] = useState(1)
  const [total, setTotal] = useState(0)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelado = false
    setCargando(true)
    setError(null)

    api
      .listarVehiculos(pagina, TAMANO_PAGINA)
      .then((respuesta) => {
        if (cancelado) return
        setVehiculos(respuesta.datos)
        setTotal(respuesta.total)
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'Ocurrió un error inesperado al cargar el catálogo.')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [pagina])

  const totalPaginas = Math.max(1, Math.ceil(total / TAMANO_PAGINA))

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Catálogo de vehículos</h1>
        <p className="mt-1 text-gray-700">Descubrí las unidades disponibles en nuestro concesionario.</p>
      </div>

      {cargando && <p className="text-gray-700">Cargando vehículos…</p>}

      {!cargando && error && <p className="text-red-600">{error}</p>}

      {!cargando && !error && vehiculos.length === 0 && (
        <p className="text-gray-700">No hay vehículos disponibles en este momento.</p>
      )}

      {!cargando && !error && vehiculos.length > 0 && (
        <>
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {vehiculos.map((vehiculo) => (
              <Link
                key={vehiculo.id}
                to={`/catalogo/${vehiculo.id}`}
                className="group overflow-hidden rounded-lg border border-gray-200 bg-white transition hover:shadow-md"
              >
                <div className="flex h-48 items-center justify-center bg-gray-100">
                  {vehiculo.imagen ? (
                    <img
                      src={vehiculo.imagen}
                      alt={`${vehiculo.marca} ${vehiculo.modelo}`}
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <span className="text-sm text-gray-500">Sin imagen</span>
                  )}
                </div>
                <div className="space-y-1 p-4">
                  <h2 className="text-lg font-semibold text-gray-900">
                    {vehiculo.marca} {vehiculo.modelo}
                  </h2>
                  <p className="text-sm text-gray-600">
                    {vehiculo.anio} · {vehiculo.condicion === 'nuevo' ? 'Nuevo' : 'Usado'}
                  </p>
                  <p className="text-lg font-bold text-gray-900">{formatearPrecio(vehiculo.precio)}</p>
                </div>
              </Link>
            ))}
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
