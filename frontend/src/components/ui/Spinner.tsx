export function Spinner({ className = '' }: { className?: string }) {
  return (
    <span
      role="status"
      aria-label="Cargando"
      className={`inline-block h-5 w-5 animate-spin rounded-full border-2 border-plata-400/30 border-t-plata-100 ${className}`}
    />
  )
}

export function ContenidoCargando({ etiqueta = 'Cargando…' }: { etiqueta?: string }) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-20 text-plata-400">
      <Spinner className="h-7 w-7" />
      <p className="text-sm">{etiqueta}</p>
    </div>
  )
}
