import type { ReactNode } from 'react'

interface PropiedadesEncabezado {
  titulo: string
  descripcion?: string
  destacado?: string
  acciones?: ReactNode
}

export function EncabezadoPagina({ titulo, descripcion, destacado, acciones }: PropiedadesEncabezado) {
  return (
    <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div className="max-w-2xl">
        {destacado && (
          <p className="mb-2 font-display text-xs font-semibold tracking-[0.25em] text-acento-400 uppercase">
            {destacado}
          </p>
        )}
        <h1 className="font-display text-3xl font-bold tracking-tight text-plata-100 sm:text-4xl">
          {titulo}
        </h1>
        {descripcion && <p className="mt-2 text-sm text-plata-400 sm:text-base">{descripcion}</p>}
      </div>
      {acciones && <div className="flex shrink-0 flex-wrap items-center gap-3">{acciones}</div>}
    </div>
  )
}
