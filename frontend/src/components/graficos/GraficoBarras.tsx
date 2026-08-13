export interface BarraDatos {
  etiqueta: string
  valor: number
  color?: string
}

interface PropiedadesGraficoBarras {
  datos: BarraDatos[]
  orientacion?: 'horizontal' | 'vertical'
  descripcion?: string
}

const COLOR_POR_DEFECTO = 'bg-acento-500'

// GraficoBarras es un gráfico simple de barras con CSS/Tailwind, sin
// librerías externas. En orientacion vertical muestra una barra por dato con
// scroll horizontal (útil para series de tiempo largas).
export function GraficoBarras({ datos, orientacion = 'horizontal', descripcion }: PropiedadesGraficoBarras) {
  if (datos.length === 0) {
    return <p className="py-6 text-center text-sm text-plata-500">Sin datos para mostrar.</p>
  }

  const maximo = Math.max(1, ...datos.map((dato) => dato.valor))

  if (orientacion === 'vertical') {
    return (
      <div>
        <div className="overflow-x-auto pb-2">
          <div className="flex items-end gap-1">
            {datos.map((dato) => {
              const alto = Math.max(4, Math.round((dato.valor / maximo) * 100))
              return (
                <div key={dato.etiqueta} className="flex w-7 flex-col items-center gap-1">
                  <div className="flex h-40 w-full items-end justify-center">
                    <div
                      title={`${dato.etiqueta}: ${dato.valor}`}
                      className={`w-4 rounded-t-md ${dato.color ?? COLOR_POR_DEFECTO}`}
                      style={{ height: `${alto}px` }}
                    />
                  </div>
                  <span className="w-full truncate text-center text-[9px] leading-tight text-plata-500">
                    {dato.etiqueta.slice(5)}
                  </span>
                </div>
              )
            })}
          </div>
        </div>
        {descripcion && <p className="pt-1 text-xs text-plata-500">{descripcion}</p>}
      </div>
    )
  }

  return (
    <div>
      <div className="space-y-3">
        {datos.map((dato) => {
          const ancho = Math.max(4, Math.round((dato.valor / maximo) * 100))
          return (
            <div key={dato.etiqueta} className="flex items-center gap-3">
              <span className="w-28 shrink-0 truncate text-sm text-plata-400">{dato.etiqueta}</span>
              <div className="flex h-4 flex-1 items-center">
                <div
                  title={`${dato.etiqueta}: ${dato.valor}`}
                  className={`h-full rounded-md ${dato.color ?? COLOR_POR_DEFECTO}`}
                  style={{ width: `${ancho}%` }}
                />
              </div>
              <span className="w-8 shrink-0 text-right font-display text-sm font-semibold text-plata-100">
                {dato.valor}
              </span>
            </div>
          )
        })}
      </div>
      {descripcion && <p className="pt-3 text-xs text-plata-500">{descripcion}</p>}
    </div>
  )
}
