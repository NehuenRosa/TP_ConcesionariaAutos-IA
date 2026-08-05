import type { ReactNode } from 'react'

interface PropiedadesEstadoVacio {
  titulo: string
  descripcion: string
  accion?: ReactNode
}

export function EstadoVacio({ titulo, descripcion, accion }: PropiedadesEstadoVacio) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 rounded-xl border border-dashed border-white/10 bg-carbono-850/40 px-6 py-16 text-center">
      <div className="flex h-12 w-12 items-center justify-center rounded-full border border-white/10 bg-carbono-800 text-plata-400">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
          <path strokeLinecap="round" strokeLinejoin="round" d="M9 10h.01M15 10h.01M9.5 14.5a4 4 0 005 0" />
        </svg>
      </div>
      <div>
        <h3 className="font-display text-lg font-semibold text-plata-100">{titulo}</h3>
        <p className="mt-1 max-w-sm text-sm text-plata-400">{descripcion}</p>
      </div>
      {accion}
    </div>
  )
}
