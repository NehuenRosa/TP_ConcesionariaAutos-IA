export const estilosEstadoConsulta: Record<string, string> = {
  pendiente: 'border-amber-400/40 bg-amber-400/10 text-amber-300',
  en_conversacion: 'border-sky-400/40 bg-sky-400/10 text-sky-300',
  cerrada: 'border-plata-400/25 bg-plata-400/10 text-plata-300',
}

export const etiquetasEstadoConsulta: Record<string, string> = {
  pendiente: 'Pendiente',
  en_conversacion: 'En conversación',
  cerrada: 'Cerrada',
}

export const estilosEstadoTestDrive: Record<string, string> = {
  solicitado: 'border-amber-400/40 bg-amber-400/10 text-amber-300',
  confirmado: 'border-sky-400/40 bg-sky-400/10 text-sky-300',
  cancelado: 'border-red-400/40 bg-red-400/10 text-red-300',
  completado: 'border-emerald-400/40 bg-emerald-400/10 text-emerald-300',
}

export const etiquetasEstadoTestDrive: Record<string, string> = {
  solicitado: 'Solicitado',
  confirmado: 'Confirmado',
  cancelado: 'Cancelado',
  completado: 'Completado',
}

export const estilosEstadoVehiculo: Record<string, string> = {
  disponible: 'border-emerald-400/40 bg-emerald-400/10 text-emerald-300',
  reservado: 'border-amber-400/40 bg-amber-400/10 text-amber-300',
  vendido: 'border-sky-400/40 bg-sky-400/10 text-sky-300',
  dado_de_baja: 'border-red-400/40 bg-red-400/10 text-red-300',
}

export const etiquetasEstadoVehiculo: Record<string, string> = {
  disponible: 'Disponible',
  reservado: 'Reservado',
  vendido: 'Vendido',
  dado_de_baja: 'Dado de baja',
}

export function EtiquetaEstado({ estado, estilos, etiqueta }: { estado: string; estilos: Record<string, string>; etiqueta?: string }) {
  return (
    <span
      className={`inline-flex items-center rounded-full border px-2.5 py-0.5 font-display text-[11px] font-semibold tracking-wide ${estilos[estado] ?? 'border-plata-400/25 bg-plata-400/10 text-plata-300'}`}
    >
      {etiqueta ?? estado}
    </span>
  )
}
