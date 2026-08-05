import type { ReactNode } from 'react'

interface PropiedadesError {
  titulo?: string
  children: ReactNode
}

export function MensajeError({ titulo = 'Ocurrió un error', children }: PropiedadesError) {
  return (
    <div
      role="alert"
      className="flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300"
    >
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="mt-0.5 h-5 w-5 shrink-0">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"
        />
      </svg>
      <div>
        <p className="font-display font-semibold">{titulo}</p>
        <p className="mt-0.5 text-red-300/80">{children}</p>
      </div>
    </div>
  )
}
