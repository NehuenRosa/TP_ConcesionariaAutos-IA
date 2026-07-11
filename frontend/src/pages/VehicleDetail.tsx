import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { ImageGallery } from '../components/ImageGallery'
import { vehicleService } from '../services/vehicleService'
import { useAuth } from '../hooks/useAuth'
import type { Vehicle } from '../types'

export function VehicleDetail() {
  const { id } = useParams()
  const [vehicle, setVehicle] = useState<Vehicle | null>(null)
  const [loading, setLoading] = useState(true)
  const { isAuthenticated } = useAuth()

  useEffect(() => {
    if (id) {
      vehicleService.getById(Number(id)).then(setVehicle).finally(() => setLoading(false))
    }
  }, [id])

  if (loading) return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-6">
        <div className="h-6 bg-gray-200 rounded-lg w-32" />
        <div className="grid md:grid-cols-2 gap-8">
          <div className="h-96 bg-gray-200 rounded-2xl" />
          <div className="space-y-4">
            <div className="h-8 bg-gray-200 rounded-lg w-3/4" />
            <div className="h-10 bg-gray-200 rounded-lg w-1/2" />
            <div className="grid grid-cols-2 gap-4">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="h-6 bg-gray-200 rounded-lg" />
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )

  if (!vehicle) return (
    <div className="text-center py-20">
      <div className="w-20 h-20 bg-gray-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
        <svg className="w-10 h-10 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <p className="text-gray-500 font-medium">Vehículo no encontrado</p>
    </div>
  )

  const formatPrice = (price: number) =>
    new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS' }).format(price)

  const statusStyles: Record<string, string> = {
    disponible: 'badge-green',
    reservado: 'badge-yellow',
    vendido: 'badge-red',
  }

  const specs = [
    { label: 'Año', value: vehicle.year },
    { label: 'Kilometraje', value: `${vehicle.mileage.toLocaleString()} km` },
    { label: 'Combustible', value: vehicle.fuel },
    { label: 'Transmisión', value: vehicle.transmission },
    { label: 'Condición', value: vehicle.condition },
    { label: 'Color', value: vehicle.color },
    { label: 'Tipo', value: vehicle.vehicle_type },
  ]

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-fade-in">
      <Link
        to="/catalogo"
        className="inline-flex items-center gap-1.5 text-sm font-medium text-gray-500 hover:text-brand-600 mb-6 transition-colors"
      >
        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        Volver al catálogo
      </Link>

      <div className="grid lg:grid-cols-5 gap-8">
        <div className="lg:col-span-3">
          <ImageGallery images={vehicle.images} alt={`${vehicle.brand} ${vehicle.model}`} />
        </div>

        <div className="lg:col-span-2">
          <div className="card p-6 lg:sticky lg:top-24">
            <div className="flex items-start justify-between mb-2">
              <div>
                <h1 className="text-2xl font-bold text-gray-900">
                  {vehicle.brand} {vehicle.model}
                </h1>
                <p className="text-gray-400 text-sm mt-0.5">{vehicle.year}</p>
              </div>
              <span className={statusStyles[vehicle.status] || 'badge-gray'}>
                {vehicle.status}
              </span>
            </div>

            <p className="text-3xl font-bold text-brand-600 mt-4 mb-6">
              {formatPrice(vehicle.price)}
            </p>

            <div className="space-y-3 border-t border-gray-100 pt-6">
              <div className="grid grid-cols-2 gap-x-4 gap-y-3">
                {specs.map((spec) => (
                  <div key={spec.label}>
                    <p className="text-xs text-gray-400 uppercase tracking-wider">{spec.label}</p>
                    <p className="text-sm font-medium text-gray-800 mt-0.5 capitalize">{spec.value}</p>
                  </div>
                ))}
              </div>
            </div>

            {vehicle.description && (
              <div className="mt-6 pt-6 border-t border-gray-100">
                <p className="text-sm text-gray-600 leading-relaxed">{vehicle.description}</p>
              </div>
            )}

            {isAuthenticated && vehicle.status === 'disponible' && (
              <div className="mt-6 pt-6 border-t border-gray-100 space-y-3">
                <Link to={`/consultar/${vehicle.id}`} className="btn-primary w-full flex items-center justify-center gap-2 text-sm">
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  Consultar al vendedor
                </Link>
                <div className="flex gap-3">
                  <Link to={`/test-drive/${vehicle.id}`} className="btn-secondary flex-1 flex items-center justify-center gap-2 text-sm">
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Test Drive
                  </Link>
                  <Link to={`/reservar/${vehicle.id}`} className="btn-warning flex-1 flex items-center justify-center gap-2 text-sm">
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
                    </svg>
                    Reservar
                  </Link>
                </div>
              </div>
            )}

            {!isAuthenticated && vehicle.status === 'disponible' && (
              <div className="mt-6 pt-6 border-t border-gray-100">
                <Link to="/login" className="text-sm text-brand-600 hover:text-brand-700 font-medium">
                  Iniciá sesión para consultar, reservar o solicitar un test drive
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
