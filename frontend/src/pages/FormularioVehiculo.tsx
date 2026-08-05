import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { VehiculoEntrada, EstadoVehiculo } from '../types/vehiculo'
import { Boton } from '../components/ui/Boton'
import { CampoArea, CampoSeleccion, CampoTexto } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { MensajeError } from '../components/ui/MensajeError'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'

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

  if (cargando) {
    return <ContenidoCargando etiqueta="Cargando vehículo…" />
  }

  return (
    <div className="mx-auto max-w-3xl px-4 py-12 sm:px-6">
      <EncabezadoPagina
        destacado="Administración"
        titulo={esEdicion ? 'Editar vehículo' : 'Nuevo vehículo'}
        descripcion="Completá la ficha técnica y las imágenes de la unidad."
        acciones={
          <Link
            to="/admin/vehiculos"
            className="inline-flex items-center gap-2 font-display text-sm font-medium text-plata-400 transition-colors hover:text-plata-100"
          >
            <span aria-hidden>←</span> Volver a la gestión
          </Link>
        }
      />

      <div className="mt-8 space-y-4">
        {error && <MensajeError>{error}</MensajeError>}

        <form
          onSubmit={guardar}
          className="space-y-6 rounded-2xl border border-white/8 bg-carbono-850/60 p-6 shadow-luz backdrop-blur-sm sm:p-8"
        >
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-2">
            <CampoTexto
              id="marca"
              etiqueta="Marca"
              required
              value={datos.marca}
              onChange={(e) => actualizarCampo('marca', e.target.value)}
            />
            <CampoTexto
              id="modelo"
              etiqueta="Modelo"
              required
              value={datos.modelo}
              onChange={(e) => actualizarCampo('modelo', e.target.value)}
            />
            <CampoTexto
              id="anio"
              etiqueta="Año"
              type="number"
              required
              min={1900}
              max={new Date().getFullYear() + 1}
              value={datos.anio}
              onChange={(e) => actualizarCampo('anio', e.target.value)}
            />
            <CampoTexto
              id="kilometraje"
              etiqueta="Kilometraje"
              type="number"
              min={0}
              value={datos.kilometraje}
              onChange={(e) => actualizarCampo('kilometraje', e.target.value)}
            />
            <CampoTexto
              id="combustible"
              etiqueta="Combustible"
              placeholder="Ej.: Nafta, Diésel, Eléctrico"
              value={datos.combustible}
              onChange={(e) => actualizarCampo('combustible', e.target.value)}
            />
            <CampoTexto
              id="transmision"
              etiqueta="Transmisión"
              placeholder="Ej.: Manual, Automática"
              value={datos.transmision}
              onChange={(e) => actualizarCampo('transmision', e.target.value)}
            />
            <CampoSeleccion
              id="tipo"
              etiqueta="Tipo de vehículo"
              required
              value={datos.tipo}
              onChange={(e) => actualizarCampo('tipo', e.target.value)}
            >
              <option value="">Seleccionar tipo</option>
              {TIPOS.map((tipo) => (
                <option key={tipo} value={tipo}>
                  {tipo.charAt(0).toUpperCase() + tipo.slice(1)}
                </option>
              ))}
            </CampoSeleccion>
            <CampoTexto
              id="precio"
              etiqueta="Precio"
              type="number"
              required
              min={1}
              step="any"
              value={datos.precio}
              onChange={(e) => actualizarCampo('precio', e.target.value)}
            />
            <CampoSeleccion
              id="condicion"
              etiqueta="Condición"
              value={datos.condicion}
              onChange={(e) => actualizarCampo('condicion', e.target.value)}
            >
              <option value="nuevo">Nuevo</option>
              <option value="usado">Usado</option>
            </CampoSeleccion>
            <CampoSeleccion
              id="estado"
              etiqueta="Estado"
              value={datos.estado}
              onChange={(e) => actualizarCampo('estado', e.target.value)}
            >
              {ESTADOS.map((estado) => (
                <option key={estado} value={estado}>
                  {ETIQUETAS_ESTADO[estado]}
                </option>
              ))}
            </CampoSeleccion>
          </div>

          <CampoArea
            id="imagenes"
            etiqueta="Imágenes (una URL por línea)"
            rows={4}
            placeholder={'https://ejemplo.com/imagen1.jpg\nhttps://ejemplo.com/imagen2.jpg'}
            value={datos.imagenes}
            onChange={(e) => actualizarCampo('imagenes', e.target.value)}
          />

          <div className="flex items-center justify-end gap-3 border-t border-white/8 pt-5">
            <Boton variante="secundario">
              <Link to="/admin/vehiculos">Cancelar</Link>
            </Boton>
            <Boton type="submit" disabled={guardando}>
              {guardando ? 'Guardando…' : esEdicion ? 'Guardar cambios' : 'Crear vehículo'}
            </Boton>
          </div>
        </form>
      </div>
    </div>
  )
}
