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

  if (loading) return <div className="text-center py-20 text-gray-500">Cargando...</div>
  if (!vehicle) return <div className="text-center py-20 text-gray-500">Vehículo no encontrado</div>

  const formatPrice = (price: number) =>
    new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS' }).format(price)

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <Link to="/catalogo" className="text-blue-600 hover:underline mb-4 block">&larr; Volver al catálogo</Link>
      <div className="grid md:grid-cols-2 gap-8">
        <ImageGallery images={vehicle.images} alt={`${vehicle.brand} ${vehicle.model}`} />
        <div>
          <h1 className="text-3xl font-bold">
            {vehicle.brand} {vehicle.model}
          </h1>
          <p className="text-4xl font-bold text-blue-600 mt-4">{formatPrice(vehicle.price)}</p>
          <div className="mt-6 space-y-3">
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div><span className="text-gray-500">Año:</span> {vehicle.year}</div>
              <div><span className="text-gray-500">Kilometraje:</span> {vehicle.mileage.toLocaleString()} km</div>
              <div><span className="text-gray-500">Combustible:</span> {vehicle.fuel}</div>
              <div><span className="text-gray-500">Transmisión:</span> {vehicle.transmission}</div>
              <div><span className="text-gray-500">Condición:</span> {vehicle.condition}</div>
              <div><span className="text-gray-500">Color:</span> {vehicle.color}</div>
              <div><span className="text-gray-500">Tipo:</span> {vehicle.vehicle_type}</div>
              <div><span className="text-gray-500">Estado:</span>
                <span className={`ml-1 px-2 py-0.5 rounded text-xs ${
                  vehicle.status === 'disponible' ? 'bg-green-100 text-green-700' :
                  vehicle.status === 'reservado' ? 'bg-yellow-100 text-yellow-700' : 'bg-red-100 text-red-700'
                }`}>{vehicle.status}</span>
              </div>
            </div>
            {vehicle.description && (
              <p className="text-gray-600 mt-4">{vehicle.description}</p>
            )}
          </div>
          {isAuthenticated && vehicle.status === 'disponible' && (
            <div className="mt-8 flex gap-3">
              <Link to={`/consultar/${vehicle.id}`} className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700">
                Consultar
              </Link>
              <Link to={`/test-drive/${vehicle.id}`} className="bg-green-600 text-white px-6 py-2 rounded hover:bg-green-700">
                Test Drive
              </Link>
              <Link to={`/reservar/${vehicle.id}`} className="bg-yellow-600 text-white px-6 py-2 rounded hover:bg-yellow-700">
                Reservar
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
