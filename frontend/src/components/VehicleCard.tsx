import { Link } from 'react-router-dom'
import type { Vehicle } from '../types'

interface Props {
  vehicle: Vehicle
}

export function VehicleCard({ vehicle }: Props) {
  const formatPrice = (price: number) =>
    new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS' }).format(price)

  return (
    <Link
      to={`/vehiculos/${vehicle.id}`}
      className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition"
    >
      <img
        src={vehicle.images?.[0] || 'https://placehold.co/600x400?text=Sin+Imagen'}
        alt={`${vehicle.brand} ${vehicle.model}`}
        className="w-full h-48 object-cover"
      />
      <div className="p-4">
        <h3 className="text-lg font-semibold">
          {vehicle.brand} {vehicle.model}
        </h3>
        <p className="text-sm text-gray-500">{vehicle.year}</p>
        <p className="text-xl font-bold text-blue-600 mt-2">{formatPrice(vehicle.price)}</p>
        <div className="flex gap-2 mt-2 text-xs text-gray-500">
          <span className="bg-gray-100 px-2 py-1 rounded">{vehicle.fuel}</span>
          <span className="bg-gray-100 px-2 py-1 rounded">{vehicle.transmission}</span>
          <span className="bg-gray-100 px-2 py-1 rounded">{vehicle.condition}</span>
        </div>
      </div>
    </Link>
  )
}
