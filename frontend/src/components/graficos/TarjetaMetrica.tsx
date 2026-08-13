import type { ReactNode } from 'react'

interface PropiedadesTarjetaMetrica {
  etiqueta: string
  valor: number
  icono?: ReactNode
  detalle?: string
}

export function TarjetaMetrica({ etiqueta, valor, icono, detalle }: PropiedadesTarjetaMetrica) {
  return (
    <div className="rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm">
      <div className="flex items-center justify-between gap-2">
        <p className="text-sm text-plata-400">{etiqueta}</p>
        {icono}
      </div>
      <p className="mt-2 font-display text-3xl font-bold text-plata-100">{valor}</p>
      {detalle && <p className="mt-1 text-xs text-plata-500">{detalle}</p>}
    </div>
  )
}
