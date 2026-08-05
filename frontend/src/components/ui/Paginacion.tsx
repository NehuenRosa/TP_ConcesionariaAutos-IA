interface PropiedadesPaginacion {
  pagina: number
  totalPaginas: number
  cambiarPagina: (pagina: number) => void
}

export function Paginacion({ pagina, totalPaginas, cambiarPagina }: PropiedadesPaginacion) {
  const deshabilitado = 'disabled:cursor-not-allowed disabled:opacity-40'
  const base = 'inline-flex items-center rounded-full border px-4 py-1.5 text-sm transition'
  return (
    <nav className="flex items-center justify-center gap-3" aria-label="Paginación">
      <button
        type="button"
        onClick={() => cambiarPagina(Math.max(1, pagina - 1))}
        disabled={pagina <= 1}
        className={`${base} ${deshabilitado} border-plata-400/25 text-plata-300 hover:border-plata-300/60 hover:text-plata-100`}
      >
        ← Anterior
      </button>
      <span className="font-display text-sm tracking-wide text-plata-400">
        {pagina} <span className="text-plata-500">/</span> {totalPaginas}
      </span>
      <button
        type="button"
        onClick={() => cambiarPagina(Math.min(totalPaginas, pagina + 1))}
        disabled={pagina >= totalPaginas}
        className={`${base} ${deshabilitado} border-plata-400/25 text-plata-300 hover:border-plata-300/60 hover:text-plata-100`}
      >
        Siguiente →
      </button>
    </nav>
  )
}
