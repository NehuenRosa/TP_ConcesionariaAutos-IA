import { Link } from 'react-router'
import type { ResumenVehiculo } from '../../types/vehiculo'
import { capitalizar, etiquetaCondicion, formatearPrecio } from '../../utils/formato'

interface PropiedadesTarjeta {
  vehiculo: ResumenVehiculo
  ruta?: string
}

export function TarjetaVehiculo({ vehiculo, ruta = `/catalogo/${vehiculo.id}` }: PropiedadesTarjeta) {
  return (
    <Link
      to={ruta}
      className="group overflow-hidden rounded-xl border border-white/8 bg-carbono-850/60 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-white/20 hover:shadow-luz"
    >
      <div className="relative aspect-[16/10] overflow-hidden bg-carbono-800">
        {vehiculo.imagen ? (
          <img
            src={vehiculo.imagen}
            alt={`${vehiculo.marca} ${vehiculo.modelo}`}
            loading="lazy"
            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full items-center justify-center text-sm text-plata-500">Sin imagen</div>
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-carbono-950/70 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100" />
        <span className="absolute top-3 left-3 rounded-full border border-white/15 bg-carbono-950/70 px-2.5 py-1 font-display text-[11px] font-semibold tracking-[0.15em] text-plata-200 uppercase backdrop-blur-sm">
          {etiquetaCondicion(vehiculo.condicion)}
        </span>
      </div>
      <div className="space-y-2 p-5">
        <div className="flex items-start justify-between gap-3">
          <div>
            <h3 className="font-display text-lg font-semibold text-plata-100 transition-colors group-hover:text-white">
              {vehiculo.marca} {vehiculo.modelo}
            </h3>
            <p className="text-xs text-plata-400">
              {vehiculo.anio}
              {vehiculo.tipo ? ` · ${capitalizar(vehiculo.tipo)}` : ''}
            </p>
          </div>
        </div>
        <p className="texto-numerico font-display text-xl font-bold text-plata-100">
          {formatearPrecio(vehiculo.precio)}
        </p>
      </div>
    </Link>
  )
}
