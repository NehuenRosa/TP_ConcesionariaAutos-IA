import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { VehiculoEntrada, EstadoVehiculo } from '../types/vehiculo'

const ESTADOS: EstadoVehiculo[] = ['disponible', 'reservado', 'vendido', 'dado_de_baja']

const TIPOS = ['sedán', 'SUV', 'hatchback', 'pick-up', 'coupe']

const ETIQUETAS_ESTADO: Record<EstadoVehiculo, string> = {
  disponible: 'Disponible',
  reservado: 'Reservado',
  vendido: 'Vendido',
  dado_de_baja: 'Dado de baja',
}

interface DatosFormulario {
  marca: string
  modelo: string
  anio: string
  kilometraje: string
  combustible: string
  transmision: string
  tipo: string
  precio: string
  condicion: string
  estado: EstadoVehiculo
  imagenes: string
}

function estadoVacio(): DatosFormulario {
  return {
    marca: '',
    modelo: '',
    anio: '',
    kilometraje: '0',
    combustible: '',
    transmision: '',
    tipo: '',
    precio: '',
    condicion: 'nuevo',
    estado: 'disponible',
    imagenes: '',
  }
}

export function FormularioVehiculo() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const esEdicion = Boolean(id)
  const [datos, setDatos] = useState<DatosFormulario>(estadoVacio)
  const [cargando, setCargando] = useState(esEdicion)
  const [guardando, setGuardando] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    let cancelado = false
    setCargando(true)
    setError(null)

    api
      .obtenerVehiculoGestion(Number(id))
      .then((vehiculo) => {
        if (cancelado) return
        setDatos({
          marca: vehiculo.marca,
          modelo: vehiculo.modelo,
          anio: String(vehiculo.anio),
          kilometraje: String(vehiculo.kilometraje),
          combustible: vehiculo.combustible,
          transmision: vehiculo.transmision,
          tipo: vehiculo.tipo,
          precio: String(vehiculo.precio),
          condicion: vehiculo.condicion,
          estado: vehiculo.estado,
          imagenes: (vehiculo.imagenes ?? []).map((imagen) => imagen.url).join('\n'),
        })
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'No se pudo cargar el vehículo.')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [id])

  function actualizarCampo(campo: keyof DatosFormulario, valor: string) {
    setDatos((actuales) => ({ ...actuales, [campo]: valor }))
  }

  async function guardar(e: React.FormEvent) {
    e.preventDefault()
    setGuardando(true)
    setError(null)

    const entrada: VehiculoEntrada = {
      marca: datos.marca.trim(),
      modelo: datos.modelo.trim(),
      anio: Number(datos.anio),
      kilometraje: Number(datos.kilometraje) || 0,
      combustible: datos.combustible.trim(),
      transmision: datos.transmision.trim(),
      tipo: datos.tipo,
      precio: Number(datos.precio),
      condicion: datos.condicion as VehiculoEntrada['condicion'],
      estado: datos.estado,
      imagenes: datos.imagenes
        .split('\n')
        .map((url) => url.trim())
        .filter((url) => url !== ''),
    }

    try {
      if (id) {
        await api.actualizarVehiculo(Number(id), entrada)
      } else {
        await api.crearVehiculo(entrada)
      }
      navigate('/admin/vehiculos')
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo guardar el vehículo.')
      setGuardando(false)
    }
  }

  const campo =
    'w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:border-gray-900 focus:outline-none'
  const etiqueta = 'block text-sm font-medium text-gray-700'

  if (cargando) {
    return <p className="text-gray-700">Cargando vehículo…</p>
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">
          {esEdicion ? 'Editar vehículo' : 'Nuevo vehículo'}
        </h1>
        <Link to="/admin/vehiculos" className="text-sm text-gray-700 hover:text-gray-900">
          ← Volver a la gestión
        </Link>
      </div>

      {error && <p className="text-red-600">{error}</p>}

      <form onSubmit={guardar} className="space-y-6 rounded-lg border border-gray-200 bg-white p-6">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label htmlFor="marca" className={etiqueta}>
              Marca
            </label>
            <input
              id="marca"
              type="text"
              required
              value={datos.marca}
              onChange={(e) => actualizarCampo('marca', e.target.value)}
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="modelo" className={etiqueta}>
              Modelo
            </label>
            <input
              id="modelo"
              type="text"
              required
              value={datos.modelo}
              onChange={(e) => actualizarCampo('modelo', e.target.value)}
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="anio" className={etiqueta}>
              Año
            </label>
            <input
              id="anio"
              type="number"
              required
              min={1900}
              max={new Date().getFullYear() + 1}
              value={datos.anio}
              onChange={(e) => actualizarCampo('anio', e.target.value)}
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="kilometraje" className={etiqueta}>
              Kilometraje
            </label>
            <input
              id="kilometraje"
              type="number"
              min={0}
              value={datos.kilometraje}
              onChange={(e) => actualizarCampo('kilometraje', e.target.value)}
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="combustible" className={etiqueta}>
              Combustible
            </label>
            <input
              id="combustible"
              type="text"
              placeholder="Ej.: Nafta, Diésel, Eléctrico"
              value={datos.combustible}
              onChange={(e) => actualizarCampo('combustible', e.target.value)}
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="transmision" className={etiqueta}>
              Transmisión
            </label>
            <input
              id="transmision"
              type="text"
              placeholder="Ej.: Manual, Automática"
              value={datos.transmision}
              onChange={(e) => actualizarCampo('transmision', e.target.value)}
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="tipo" className={etiqueta}>
              Tipo de vehículo
            </label>
            <select
              id="tipo"
              required
              value={datos.tipo}
              onChange={(e) => actualizarCampo('tipo', e.target.value)}
              className={campo}
            >
              <option value="">Seleccionar tipo</option>
              {TIPOS.map((tipo) => (
                <option key={tipo} value={tipo}>
                  {tipo.charAt(0).toUpperCase() + tipo.slice(1)}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="precio" className={etiqueta}>
              Precio
            </label>
            <input
              id="precio"
              type="number"
              required
              min={1}
              step="any"
              value={datos.precio}
              onChange={(e) => actualizarCampo('precio', e.target.value)}
              className={campo}
            />
          </div>
          <div>
            <label htmlFor="condicion" className={etiqueta}>
              Condición
            </label>
            <select
              id="condicion"
              value={datos.condicion}
              onChange={(e) => actualizarCampo('condicion', e.target.value)}
              className={campo}
            >
              <option value="nuevo">Nuevo</option>
              <option value="usado">Usado</option>
            </select>
          </div>
          <div>
            <label htmlFor="estado" className={etiqueta}>
              Estado
            </label>
            <select
              id="estado"
              value={datos.estado}
              onChange={(e) => actualizarCampo('estado', e.target.value)}
              className={campo}
            >
              {ESTADOS.map((estado) => (
                <option key={estado} value={estado}>
                  {ETIQUETAS_ESTADO[estado]}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div>
          <label htmlFor="imagenes" className={etiqueta}>
            Imágenes (una URL por línea)
          </label>
          <textarea
            id="imagenes"
            rows={4}
            placeholder={'https://ejemplo.com/imagen1.jpg\nhttps://ejemplo.com/imagen2.jpg'}
            value={datos.imagenes}
            onChange={(e) => actualizarCampo('imagenes', e.target.value)}
            className={campo}
          />
        </div>

        <div className="flex items-center justify-end gap-3">
          <Link
            to="/admin/vehiculos"
            className="rounded-md border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
          >
            Cancelar
          </Link>
          <button
            type="submit"
            disabled={guardando}
            className="rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {guardando ? 'Guardando…' : esEdicion ? 'Guardar cambios' : 'Crear vehículo'}
          </button>
        </div>
      </form>
    </div>
  )
}
