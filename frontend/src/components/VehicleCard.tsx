import { Link } from 'react-router-dom'
import type { Vehicle } from '../types'

interface Props {
  vehicle: Vehicle
}

const formatPrice = (price: number) =>
  new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS' }).format(price)

export function VehicleCard({ vehicle }: Props) {
  const hasImage = vehicle.images?.[0] && !vehicle.images[0].includes('placehold.co')

  return (
    <Link
      to={`/vehiculos/${vehicle.id}`}
      className="card-hover overflow-hidden flex flex-col group cursor-pointer"
    >
      <div className="relative aspect-[4/3] bg-[#DFE1ED] overflow-hidden">
        {hasImage ? (
          <img
            src={vehicle.images[0]}
            alt={`${vehicle.brand} ${vehicle.model}`}
            className="w-full h-full object-cover group-hover:scale-103 transition-transform duration-200 ease-out"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <svg className="w-12 h-12 text-[#B7BAD0]" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
        )}
        <div className="absolute top-3 right-3">
          <span
            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
              vehicle.condition === 'nuevo'
                ? 'bg-[#E1F5EE] text-[#085041]'
                : 'bg-[#FAEEDA] text-[#854F0B]'
            }`}
          >
            {vehicle.condition === 'nuevo' ? 'Nuevo' : 'Usado'}
          </span>
        </div>
      </div>
      <div className="p-4 flex-1 flex flex-col gap-1.5">
        <div>
          <h3 className="text-[15px] font-medium text-text-primary leading-tight">
            {vehicle.brand} {vehicle.model}
          </h3>
          <p className="text-xs text-text-secondary mt-0.5">{vehicle.year} · {vehicle.mileage.toLocaleString()} km</p>
        </div>
        <p className="text-[17px] font-semibold text-brand-500 mt-1">
          {formatPrice(vehicle.price)}
        </p>
        <div className="flex flex-wrap gap-1.5 mt-1">
          <span className="chip bg-surface text-surface-text capitalize">{vehicle.transmission}</span>
          <span className="chip bg-surface text-surface-text capitalize">{vehicle.fuel}</span>
          <span className="chip bg-surface text-surface-text capitalize">{vehicle.vehicle_type}</span>
        </div>
      </div>
    </Link>
  )
}
