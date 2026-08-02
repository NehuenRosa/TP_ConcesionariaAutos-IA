import { useDeferredValue, useEffect, useState } from 'react'
import { Link } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { FiltrosVehiculos, ResumenVehiculo } from '../types/vehiculo'

const TAMANO_PAGINA = 12

const TIPOS = ['sedán', 'SUV', 'hatchback', 'pick-up', 'coupe']

const COMBUSTIBLES = ['Nafta', 'Diésel', 'Eléctrico', 'Híbrido', 'GNC']

const MARCAS_SUGERIDAS = ['Toyota', 'Ford', 'Volkswagen', 'Chevrolet', 'Fiat', 'Renault', 'Peugeot']

interface DatosFiltros {
  busqueda: string
  marca: string
  modelo: string
  anioMin: string
  anioMax: string
  precioMin: string
  precioMax: string
  tipo: string
  combustible: string
  condicion: string
  ordenPor: FiltrosVehiculos['ordenPor']
  ordenDireccion: FiltrosVehiculos['ordenDireccion']
}

function filtrosVacios(): DatosFiltros {
  return {
    busqueda: '',
    marca: '',
    modelo: '',
    anioMin: '',
    anioMax: '',
    precioMin: '',
    precioMax: '',
    tipo: '',
    combustible: '',
    condicion: '',
    ordenPor: 'anio',
    ordenDireccion: 'desc',
  }
}

function aFiltrosConsulta(datos: DatosFiltros): FiltrosVehiculos {
  const filtros: FiltrosVehiculos = {
    busqueda: datos.busqueda.trim() || undefined,
    marca: datos.marca.trim() || undefined,
    modelo: datos.modelo.trim() || undefined,
    anioMin: datos.anioMin ? Number(datos.anioMin) : undefined,
    anioMax: datos.anioMax ? Number(datos.anioMax) : undefined,
    precioMin: datos.precioMin ? Number(datos.precioMin) : undefined,
    precioMax: datos.precioMax ? Number(datos.precioMax) : undefined,
    tipo: datos.tipo || undefined,
    combustible: datos.combustible || undefined,
    condicion: (datos.condicion as FiltrosVehiculos['condicion']) || undefined,
    ordenPor: datos.ordenPor,
    ordenDireccion: datos.ordenDireccion,
  }
  return filtros
}

function formatearPrecio(precio: number): string {
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
    maximumFractionDigits: 0,
  }).format(precio)
}

export function Catalogo() {
  const [datos, setDatos] = useState<DatosFiltros>(filtrosVacios)
  const busquedaDiferida = useDeferredValue(datos.busqueda)
  const [vehiculos, setVehiculos] = useState<ResumenVehiculo[]>([])
  const [pagina, setPagina] = useState(1)
  const [total, setTotal] = useState(0)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelado = false
    setCargando(true)
    setError(null)

    const filtros: FiltrosVehiculos = aFiltrosConsulta({ ...datos, busqueda: busquedaDiferida })

    api
      .listarVehiculos(pagina, TAMANO_PAGINA, filtros)
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
  }, [pagina, busquedaDiferida, datos])

  const totalPaginas = Math.max(1, Math.ceil(total / TAMANO_PAGINA))
  const hayFiltros = Object.values(datos).some(
    (valor) => valor !== '' && valor !== 'anio' && valor !== 'desc',
  )

  function actualizar(campo: keyof DatosFiltros, valor: string) {
    setPagina(1)
    setDatos((actuales) => ({ ...actuales, [campo]: valor }))
  }

  function limpiar() {
    setPagina(1)
    setDatos(filtrosVacios())
  }

  const campo =
    'w-full rounded-md border border-gray-300 px-3 py-1.5 text-sm text-gray-700 focus:border-gray-900 focus:outline-none'
  const etiqueta = 'block text-xs font-medium text-gray-600'

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Catálogo de vehículos</h1>
        <p className="mt-1 text-gray-700">
          Buscá y filtrá las unidades disponibles en nuestro concesionario.
        </p>
      </div>

      <form
        onSubmit={(e) => e.preventDefault()}
        className="space-y-4 rounded-lg border border-gray-200 bg-white p-4"
      >
        <div>
          <label htmlFor="busqueda" className={etiqueta}>
            Búsqueda por marca o modelo
          </label>
          <input
            id="busqueda"
            type="search"
            value={datos.busqueda}
            onChange={(e) => actualizar('busqueda', e.target.value)}
            placeholder="Ej.: Corolla, Toyota…"
            className={campo}
          />
        </div>

        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
          <div>
            <label htmlFor="marca" className={etiqueta}>
              Marca
            </label>
            <input
              id="marca"
              type="text"
              list="marcas-sugeridas"
              value={datos.marca}
              onChange={(e) => actualizar('marca', e.target.value)}
              placeholder="Todas"
              className={campo}
            />
            <datalist id="marcas-sugeridas">
              {MARCAS_SUGERIDAS.map((marca) => (
                <option key={marca} value={marca} />
              ))}
            </datalist>
          </div>
          <div>
            <label htmlFor="modelo" className={etiqueta}>
              Modelo
            </label>
            <input
              id="modelo"
              type="text"
              value={datos.modelo}
              onChange={(e) => actualizar('modelo', e.target.value)}
              placeholder="Todos"
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="tipo" className={etiqueta}>
              Tipo
            </label>
            <select
              id="tipo"
              value={datos.tipo}
              onChange={(e) => actualizar('tipo', e.target.value)}
              className={campo}
            >
              <option value="">Todos</option>
              {TIPOS.map((tipo) => (
                <option key={tipo} value={tipo}>
                  {tipo.charAt(0).toUpperCase() + tipo.slice(1)}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="combustible" className={etiqueta}>
              Combustible
            </label>
            <select
              id="combustible"
              value={datos.combustible}
              onChange={(e) => actualizar('combustible', e.target.value)}
              className={campo}
            >
              <option value="">Todos</option>
              {COMBUSTIBLES.map((combustible) => (
                <option key={combustible} value={combustible}>
                  {combustible}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="condicion" className={etiqueta}>
              Condición
            </label>
            <select
              id="condicion"
              value={datos.condicion}
              onChange={(e) => actualizar('condicion', e.target.value)}
              className={campo}
            >
              <option value="">Todas</option>
              <option value="nuevo">Nuevo</option>
              <option value="usado">Usado</option>
            </select>
          </div>
          <div>
            <label htmlFor="anio-min" className={etiqueta}>
              Año desde
            </label>
            <input
              id="anio-min"
              type="number"
              min={1900}
              max={new Date().getFullYear() + 1}
              value={datos.anioMin}
              onChange={(e) => actualizar('anioMin', e.target.value)}
              placeholder="Ej.: 2018"
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="anio-max" className={etiqueta}>
              Año hasta
            </label>
            <input
              id="anio-max"
              type="number"
              min={1900}
              max={new Date().getFullYear() + 1}
              value={datos.anioMax}
              onChange={(e) => actualizar('anioMax', e.target.value)}
              placeholder="Ej.: 2022"
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="precio-min" className={etiqueta}>
              Precio mínimo
            </label>
            <input
              id="precio-min"
              type="number"
              min={0}
              value={datos.precioMin}
              onChange={(e) => actualizar('precioMin', e.target.value)}
              placeholder="Ej.: 10000000"
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="precio-max" className={etiqueta}>
              Precio máximo
            </label>
            <input
              id="precio-max"
              type="number"
              min={0}
              value={datos.precioMax}
              onChange={(e) => actualizar('precioMax', e.target.value)}
              placeholder="Ej.: 30000000"
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="orden-por" className={etiqueta}>
              Ordenar por
            </label>
            <select
              id="orden-por"
              value={datos.ordenPor}
              onChange={(e) => actualizar('ordenPor', e.target.value)}
              className={campo}
            >
              <option value="anio">Año</option>
              <option value="precio">Precio</option>
            </select>
          </div>
          <div>
            <label htmlFor="orden-direccion" className={etiqueta}>
              Dirección
            </label>
            <select
              id="orden-direccion"
              value={datos.ordenDireccion}
              onChange={(e) => actualizar('ordenDireccion', e.target.value)}
              className={campo}
            >
              <option value="desc">Descendente</option>
              <option value="asc">Ascendente</option>
            </select>
          </div>
        </div>

        {hayFiltros && (
          <button
            type="button"
            onClick={limpiar}
            className="rounded-md border border-gray-300 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
          >
            Limpiar filtros
          </button>
        )}
      </form>

      {cargando && <p className="text-gray-700">Cargando vehículos…</p>}

      {!cargando && error && <p className="text-red-600">{error}</p>}

      {!cargando && !error && vehiculos.length === 0 && (
        <p className="text-gray-700">
          {hayFiltros
            ? 'No hay vehículos que coincidan con los filtros aplicados.'
            : 'No hay vehículos disponibles en este momento.'}
        </p>
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
                    {vehiculo.tipo ? ` · ${vehiculo.tipo.charAt(0).toUpperCase() + vehiculo.tipo.slice(1)}` : ''}
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
