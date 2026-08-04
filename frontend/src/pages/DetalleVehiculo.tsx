import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { Vehiculo } from '../types/vehiculo'
import { useAuth } from '../hooks/useAuth'

function formatearPrecio(precio: number): string {
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
    maximumFractionDigits: 0,
  }).format(precio)
}

function formatearKilometraje(kilometraje: number): string {
  return new Intl.NumberFormat('es-AR').format(kilometraje)
}

export function DetalleVehiculo() {
  const { id } = useParams<{ id: string }>()
  const { usuario } = useAuth()
  const [vehiculo, setVehiculo] = useState<Vehiculo | null>(null)
  const [imagenActiva, setImagenActiva] = useState(0)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [mostrarFormulario, setMostrarFormulario] = useState(false)
  const [mensajeConsulta, setMensajeConsulta] = useState('')
  const [enviandoConsulta, setEnviandoConsulta] = useState(false)
  const [errorConsulta, setErrorConsulta] = useState<string | null>(null)
  const [exitoConsulta, setExitoConsulta] = useState(false)

  const esCliente = usuario?.rol === 'cliente'

  useEffect(() => {
    if (!id) return
    let cancelado = false
    setCargando(true)
    setError(null)
    setImagenActiva(0)

    api
      .obtenerVehiculo(Number(id))
      .then((dato) => {
        if (cancelado) return
        setVehiculo(dato)
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'Ocurrió un error inesperado al cargar el vehículo.')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [id])

  const handleEnviarConsulta = async () => {
    if (!vehiculo || !mensajeConsulta.trim()) return

    setEnviandoConsulta(true)
    setErrorConsulta(null)

    try {
      await api.crearConsulta({
        vehiculoId: vehiculo.id,
        mensaje: mensajeConsulta.trim(),
      })
      setExitoConsulta(true)
      setMensajeConsulta('')
      setMostrarFormulario(false)
    } catch (e: unknown) {
      setErrorConsulta(e instanceof ErrorApi ? e.message : 'No se pudo enviar la consulta')
    } finally {
      setEnviandoConsulta(false)
    }
  }

  if (cargando) {
    return <p className="text-gray-700">Cargando vehículo…</p>
  }

  if (error || !vehiculo) {
    return (
      <div className="space-y-4">
        <h1 className="text-2xl font-bold text-gray-900">Vehículo no encontrado</h1>
        <p className="text-gray-700">{error ?? 'El vehículo solicitado no existe o no está disponible.'}</p>
        <Link
          to="/catalogo"
          className="inline-block rounded-md bg-gray-900 px-4 py-2 text-white hover:bg-gray-700"
        >
          Volver al catálogo
        </Link>
      </div>
    )
  }

  const imagenes = vehiculo.imagenes ?? []

  return (
    <div className="space-y-8">
      <Link to="/catalogo" className="text-sm text-gray-700 hover:text-gray-900">
        ← Volver al catálogo
      </Link>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
        <div className="space-y-4">
          <div className="flex h-80 items-center justify-center overflow-hidden rounded-lg bg-gray-100">
            {imagenes[imagenActiva] ? (
              <img
                src={imagenes[imagenActiva].url}
                alt={`${vehiculo.marca} ${vehiculo.modelo}`}
                className="h-full w-full object-cover"
              />
            ) : (
              <span className="text-gray-500">Sin imagen</span>
            )}
          </div>

          {imagenes.length > 1 && (
            <div className="flex gap-2">
              {imagenes.map((imagen, indice) => (
                <button
                  key={imagen.id}
                  type="button"
                  onClick={() => setImagenActiva(indice)}
                  className={`h-16 w-24 overflow-hidden rounded-md border ${
                    indice === imagenActiva ? 'border-gray-900' : 'border-gray-200'
                  }`}
                  aria-label={`Ver imagen ${indice + 1}`}
                >
                  <img src={imagen.url} alt="" className="h-full w-full object-cover" />
                </button>
              ))}
            </div>
          )}
        </div>

        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">
              {vehiculo.marca} {vehiculo.modelo}
            </h1>
            <p className="mt-1 text-sm text-gray-600">
              {vehiculo.anio} · {vehiculo.condicion === 'nuevo' ? 'Nuevo' : 'Usado'}
            </p>
            <p className="mt-4 text-3xl font-bold text-gray-900">{formatearPrecio(vehiculo.precio)}</p>
          </div>

          <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="rounded-lg border border-gray-200 p-4">
              <dt className="text-sm text-gray-500">Año</dt>
              <dd className="text-lg font-semibold text-gray-900">{vehiculo.anio}</dd>
            </div>
            <div className="rounded-lg border border-gray-200 p-4">
              <dt className="text-sm text-gray-500">Kilometraje</dt>
              <dd className="text-lg font-semibold text-gray-900">
                {formatearKilometraje(vehiculo.kilometraje)} km
              </dd>
            </div>
            <div className="rounded-lg border border-gray-200 p-4">
              <dt className="text-sm text-gray-500">Combustible</dt>
              <dd className="text-lg font-semibold text-gray-900">{vehiculo.combustible}</dd>
            </div>
            <div className="rounded-lg border border-gray-200 p-4">
              <dt className="text-sm text-gray-500">Transmisión</dt>
              <dd className="text-lg font-semibold text-gray-900">{vehiculo.transmision}</dd>
            </div>
            <div className="rounded-lg border border-gray-200 p-4">
              <dt className="text-sm text-gray-500">Tipo</dt>
              <dd className="text-lg font-semibold text-gray-900">
                {vehiculo.tipo ? vehiculo.tipo.charAt(0).toUpperCase() + vehiculo.tipo.slice(1) : '—'}
              </dd>
            </div>
            <div className="rounded-lg border border-gray-200 p-4">
              <dt className="text-sm text-gray-500">Condición</dt>
              <dd className="text-lg font-semibold text-gray-900">
                {vehiculo.condicion === 'nuevo' ? 'Nuevo' : 'Usado'}
              </dd>
            </div>
          </dl>

          {esCliente && (
            <div className="border-t border-gray-200 pt-6">
              {exitoConsulta ? (
                <div className="rounded-lg bg-green-50 p-4">
                  <p className="text-green-800">Consulta enviada correctamente. Un vendedor te responderá pronto.</p>
                  <Link
                    to="/mis-consultas"
                    className="mt-2 inline-block text-sm font-medium text-green-700 hover:text-green-900"
                  >
                    Ver mis consultas →
                  </Link>
                </div>
              ) : mostrarFormulario ? (
                <div className="space-y-4">
                  <h3 className="text-lg font-semibold text-gray-900">Enviar consulta</h3>
                  <textarea
                    value={mensajeConsulta}
                    onChange={(e) => setMensajeConsulta(e.target.value)}
                    placeholder="Escribí tu consulta sobre este vehículo..."
                    className="w-full rounded-lg border border-gray-300 p-3 text-gray-900 focus:border-gray-500 focus:outline-none"
                    rows={4}
                    disabled={enviandoConsulta}
                  />
                  {errorConsulta && (
                    <p className="text-sm text-red-600">{errorConsulta}</p>
                  )}
                  <div className="flex gap-2">
                    <button
                      type="button"
                      onClick={handleEnviarConsulta}
                      disabled={enviandoConsulta || !mensajeConsulta.trim()}
                      className="rounded-md bg-gray-900 px-4 py-2 text-white hover:bg-gray-700 disabled:opacity-50"
                    >
                      {enviandoConsulta ? 'Enviando...' : 'Enviar consulta'}
                    </button>
                    <button
                      type="button"
                      onClick={() => {
                        setMostrarFormulario(false)
                        setErrorConsulta(null)
                      }}
                      className="rounded-md border border-gray-300 px-4 py-2 text-gray-700 hover:bg-gray-50"
                    >
                      Cancelar
                    </button>
                  </div>
                </div>
              ) : (
                <button
                  type="button"
                  onClick={() => setMostrarFormulario(true)}
                  className="w-full rounded-md bg-gray-900 px-4 py-3 text-white hover:bg-gray-700"
                >
                  Consultar sobre este vehículo
                </button>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
