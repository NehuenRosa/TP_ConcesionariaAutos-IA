import type { ButtonHTMLAttributes } from 'react'

export type VarianteBoton = 'primario' | 'secundario' | 'fantasma' | 'peligro' | 'acento'
export type TamanoBoton = 'sm' | 'md' | 'lg'

const VARIANTES: Record<VarianteBoton, string> = {
  primario: 'bg-plata-100 text-carbono-900 hover:bg-white focus-visible:ring-plata-300/60',
  secundario:
    'border border-plata-400/25 bg-carbono-800/50 text-plata-100 hover:border-plata-300/50 hover:bg-carbono-700/60 focus-visible:ring-plata-400/40',
  fantasma: 'text-plata-300 hover:bg-white/5 hover:text-plata-100 focus-visible:ring-plata-400/40',
  peligro:
    'border border-red-500/40 bg-red-500/10 text-red-300 hover:bg-red-500/20 focus-visible:ring-red-500/40',
  acento: 'bg-acento-500 text-carbono-900 hover:bg-acento-400 focus-visible:ring-acento-400/50',
}

const TAMANOS: Record<TamanoBoton, string> = {
  sm: 'gap-1.5 px-3 py-1.5 text-xs',
  md: 'gap-2 px-5 py-2.5 text-sm',
  lg: 'gap-2.5 px-8 py-3.5 text-sm',
}

interface PropiedadesBoton extends ButtonHTMLAttributes<HTMLButtonElement> {
  variante?: VarianteBoton
  tamano?: TamanoBoton
}

export function Boton({
  variante = 'primario',
  tamano = 'md',
  className = '',
  type = 'button',
  ...resto
}: PropiedadesBoton) {
  return (
    <button
      type={type}
      className={`inline-flex cursor-pointer items-center justify-center rounded-full font-display font-semibold tracking-wide transition-all duration-200 focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-offset-carbono-950 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50 ${VARIANTES[variante]} ${TAMANOS[tamano]} ${className}`}
      {...resto}
    />
  )
}
