import { useDeferredValue, useEffect, useState } from 'react'
import { useSearchParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { FiltrosVehiculos, ResumenVehiculo } from '../types/vehiculo'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { TarjetaVehiculo } from '../components/ui/TarjetaVehiculo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { MensajeError } from '../components/ui/MensajeError'
import { Paginacion } from '../components/ui/Paginacion'
import { EstadoVacio } from '../components/ui/EstadoVacio'
import { CampoTexto, CampoSeleccion } from '../components/ui/Campo'
import { Boton } from '../components/ui/Boton'

const TAMANO_PAGINA = 12

const TIPOS = [
  { valor: 'sedán', etiqueta: 'Sedán' },
  { valor: 'suv', etiqueta: 'SUV' },
  { valor: 'hatchback', etiqueta: 'Hatchback' },
  { valor: 'pick-up', etiqueta: 'Pick-up' },
  { valor: 'coupe', etiqueta: 'Coupé' },
]

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
  return {
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
}

export function Catalogo() {
  const [searchParams] = useSearchParams()
  const [datos, setDatos] = useState<DatosFiltros>(() => {
    const base = filtrosVacios()
    const tipo = searchParams.get('tipo')
    if (tipo) base.tipo = tipo.toLowerCase()
    return base
  })
  const busquedaDiferida = useDeferredValue(datos.busqueda)
  const [vehiculos, setVehiculos] = useState<ResumenVehiculo[]>([])
  const [pagina, setPagina] = useState(1)
  const [total, setTotal] = useState(0)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Refleja el tipo elegido desde las tarjetas de la portada (p. ej.
  // /catalogo?tipo=suv) aunque ya estés dentro del catálogo.
  useEffect(() => {
    const tipo = (searchParams.get('tipo') ?? '').toLowerCase()
    setPagina(1)
    setDatos((actuales) => (actuales.tipo === tipo ? actuales : { ...actuales, tipo }))
  }, [searchParams])

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

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <EncabezadoPagina
        destacado="Colección 2026"
        titulo="Catálogo de vehículos"
        descripcion="Buscá y filtrá las unidades disponibles en nuestro concesionario."
      />

      <form
        onSubmit={(e) => e.preventDefault()}
        className="mt-8 space-y-5 rounded-2xl border border-white/8 bg-carbono-850/50 p-6 backdrop-blur-sm"
      >
        <CampoTexto
          id="busqueda"
          etiqueta="Búsqueda por marca o modelo"
          type="search"
          value={datos.busqueda}
          onChange={(e) => actualizar('busqueda', e.target.value)}
          placeholder="Ej.: Corolla, Toyota…"
        />

        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
          <div>
            <label htmlFor="marca" className="etiqueta">
              Marca
            </label>
            <input
              id="marca"
              type="text"
              list="marcas-sugeridas"
              value={datos.marca}
              onChange={(e) => actualizar('marca', e.target.value)}
              placeholder="Todas"
              className="campo"
            />
            <datalist id="marcas-sugeridas">
              {MARCAS_SUGERIDAS.map((marca) => (
                <option key={marca} value={marca} />
              ))}
            </datalist>
          </div>
          <CampoTexto
            id="modelo"
            etiqueta="Modelo"
            value={datos.modelo}
            onChange={(e) => actualizar('modelo', e.target.value)}
            placeholder="Todos"
          />
          <CampoSeleccion
            id="tipo"
            etiqueta="Tipo"
            value={datos.tipo}
            onChange={(e) => actualizar('tipo', e.target.value)}
          >
            <option value="">Todos</option>
            {TIPOS.map((tipo) => (
              <option key={tipo.valor} value={tipo.valor}>
                {tipo.etiqueta}
              </option>
            ))}
          </CampoSeleccion>
          <CampoSeleccion
            id="combustible"
            etiqueta="Combustible"
            value={datos.combustible}
            onChange={(e) => actualizar('combustible', e.target.value)}
          >
            <option value="">Todos</option>
            {COMBUSTIBLES.map((combustible) => (
              <option key={combustible} value={combustible}>
                {combustible}
              </option>
            ))}
          </CampoSeleccion>
          <CampoSeleccion
            id="condicion"
            etiqueta="Condición"
            value={datos.condicion}
            onChange={(e) => actualizar('condicion', e.target.value)}
          >
            <option value="">Todas</option>
            <option value="nuevo">Nuevo</option>
            <option value="usado">Usado</option>
          </CampoSeleccion>
          <CampoTexto
            id="anio-min"
            etiqueta="Año desde"
            type="number"
            min={1900}
            max={new Date().getFullYear() + 1}
            value={datos.anioMin}
            onChange={(e) => actualizar('anioMin', e.target.value)}
            placeholder="Ej.: 2018"
          />
          <CampoTexto
            id="anio-max"
            etiqueta="Año hasta"
            type="number"
            min={1900}
            max={new Date().getFullYear() + 1}
            value={datos.anioMax}
            onChange={(e) => actualizar('anioMax', e.target.value)}
            placeholder="Ej.: 2022"
          />
          <CampoTexto
            id="precio-min"
            etiqueta="Precio mínimo"
            type="number"
            min={0}
            value={datos.precioMin}
            onChange={(e) => actualizar('precioMin', e.target.value)}
            placeholder="Ej.: 10000000"
          />
          <CampoTexto
            id="precio-max"
            etiqueta="Precio máximo"
            type="number"
            min={0}
            value={datos.precioMax}
            onChange={(e) => actualizar('precioMax', e.target.value)}
            placeholder="Ej.: 30000000"
          />
          <CampoSeleccion
            id="orden-por"
            etiqueta="Ordenar por"
            value={datos.ordenPor}
            onChange={(e) => actualizar('ordenPor', e.target.value)}
          >
            <option value="anio">Año</option>
            <option value="precio">Precio</option>
          </CampoSeleccion>
          <CampoSeleccion
            id="orden-direccion"
            etiqueta="Dirección"
            value={datos.ordenDireccion}
            onChange={(e) => actualizar('ordenDireccion', e.target.value)}
          >
            <option value="desc">Descendente</option>
            <option value="asc">Ascendente</option>
          </CampoSeleccion>
        </div>

        <div className="flex flex-wrap items-center justify-between gap-3 pt-1">
          <p className="text-xs text-plata-500">{total} vehículo{total === 1 ? '' : 's'}</p>
          {hayFiltros && (
            <Boton variante="fantasma" tamano="sm" onClick={limpiar}>
              Limpiar filtros
            </Boton>
          )}
        </div>
      </form>

      {cargando && <ContenidoCargando etiqueta="Cargando vehículos…" />}

      {!cargando && error && (
        <div className="mt-6">
          <MensajeError>{error}</MensajeError>
        </div>
      )}

      {!cargando && !error && vehiculos.length === 0 && (
        <div className="mt-8">
          <EstadoVacio
            titulo={hayFiltros ? 'Sin coincidencias' : 'Sin vehículos disponibles'}
            descripcion={
              hayFiltros
                ? 'No hay vehículos que coincidan con los filtros aplicados. Probá ajustarlos.'
                : 'No hay vehículos disponibles en este momento. Volvé pronto.'
            }
            accion={
              hayFiltros ? (
                <Boton variante="secundario" tamano="sm" onClick={limpiar}>
                  Limpiar filtros
                </Boton>
              ) : undefined
            }
          />
        </div>
      )}

      {!cargando && !error && vehiculos.length > 0 && (
        <>
          <div className="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {vehiculos.map((vehiculo) => (
              <TarjetaVehiculo key={vehiculo.id} vehiculo={vehiculo} />
            ))}
          </div>

          <div className="mt-12">
            <Paginacion pagina={pagina} totalPaginas={totalPaginas} cambiarPagina={setPagina} />
          </div>
        </>
      )}
    </div>
  )
}
